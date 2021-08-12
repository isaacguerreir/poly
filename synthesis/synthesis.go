/*
Package synthesis fixes synthetic DNA molecules in preparation for synthesis.

Many synthesis companies have restrictions on the DNA they can synthesize. This
synthesis fixer takes advantage of synonymous codons in protein coding
sequences (CDS) to remove problematic sequences that either users don't want
(like restriction enzymes sites) or that would cause DNA synthesis companies to
reject a synthesis project.

This synthesis fixer is meant to cover the majority of use cases for DNA
fixing. It is not intended to cover all possible use cases, since the majority
of DNA design does not actually have these edge cases.

FixCds does not guarantee that all requested features will be removed. If you
have use case that FixCds cannot properly fix, please put an issue in the poly
github.
*/
package synthesis

import (
	"errors"
	"regexp"
	"strings"
	"sync"

	"github.com/TimothyStiles/poly/transform"
	"github.com/TimothyStiles/poly/transform/codon"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite" // imports CGO-less sqlite
)

// FixIterations is the standard default for the number of rounds FixCds will
// iterate for while trying to find a fixed solution for a sequence.
var FixIterations = 100

// DnaSuggestion is a suggestion of a fixer, generated by a problematicSequenceFunc.
type DnaSuggestion struct {
	Start          int    `db:"start"`
	End            int    `db:"end"`
	Bias           string `db:"gcbias"`
	QuantityFixes  int    `db:"quantityfixes"`
	SuggestionType string `db:"suggestiontype"`
}

// A change is an change to given DNA sequence
type Change struct {
	Position int    `db:"position"`
	Step     int    `db:"step"`
	From     string `db:"codonfrom"`
	To       string `db:"codonto"`
	Reason   string `db:"reason"`
}

type dbDnaSuggestion struct {
	Start          int    `db:"start"`
	End            int    `db:"end"`
	Bias           string `db:"gcbias"`
	QuantityFixes  int    `db:"quantityfixes"`
	SuggestionType string `db:"suggestiontype"`
	Step           int    `db:"step"`
	ID             int    `db:"id"`
}

// RemoveSequence is a generator for a problematicSequenceFuncs for specific sequences.
func RemoveSequence(sequencesToRemove []string, reason string) func(string, chan DnaSuggestion, *sync.WaitGroup) {
	return func(sequence string, c chan DnaSuggestion, wg *sync.WaitGroup) {
		var enzymes []string
		for _, enzyme := range sequencesToRemove {
			enzymes = []string{enzyme, transform.ReverseComplement(enzyme)}
			for _, site := range enzymes {
				re := regexp.MustCompile(site)
				locs := re.FindAllStringIndex(sequence, -1)
				for _, loc := range locs {
					position := loc[0] / 3
					leftover := loc[0] % 3
					switch {
					case leftover == 0:
						c <- DnaSuggestion{position, (loc[1] / 3), "NA", 1, reason}
					case leftover != 0:
						c <- DnaSuggestion{position, (loc[1] / 3) - 1, "NA", 1, reason}
					}
				}
			}
		}
		wg.Done()
	}
}

// RemoveRepeat is a generator to make a problematicSequenceFunc for repeats.
func RemoveRepeat(repeatLen int) func(string, chan DnaSuggestion, *sync.WaitGroup) {
	return func(sequence string, c chan DnaSuggestion, wg *sync.WaitGroup) {
		// Get a kmer list
		kmers := make(map[string]bool)
		for i := 0; i < len(sequence)-repeatLen; i++ {
			_, alreadyFound := kmers[sequence[i:i+repeatLen]]
			if alreadyFound {
				position := i / 3
				leftover := i % 3
				switch {
				case leftover == 0:
					c <- DnaSuggestion{position, ((i + repeatLen) / 3), "NA", 1, "Repeat sequence"}
				case leftover != 0:
					c <- DnaSuggestion{position, ((i + repeatLen) / 3) - 1, "NA", 1, "Repeat sequence"}
				}
			}
			kmers[sequence[i:i+repeatLen]] = true
		}
		wg.Done()
	}
}

func findProblems(sequence string, problematicSequenceFuncs []func(string, chan DnaSuggestion, *sync.WaitGroup)) []DnaSuggestion {
	// Run functions to get suggestions
	suggestions := make(chan DnaSuggestion, 100)
	var wg sync.WaitGroup
	for _, f := range problematicSequenceFuncs {
		wg.Add(1)
		go f(sequence, suggestions, &wg)
	}
	wg.Wait()
	close(suggestions)

	var suggestionsList []DnaSuggestion
	for suggestion := range suggestions {
		suggestionsList = append(suggestionsList, suggestion)
	}
	return suggestionsList
}

// FixCds fixes a CDS given the CDS sequence, a codon table, and a list of functions to solve for.
func FixCds(sqlitePath string, sequence string, codontable codon.Table, problematicSequenceFuncs []func(string, chan DnaSuggestion, *sync.WaitGroup)) (string, []Change, error) {
	db := sqlx.MustConnect("sqlite", sqlitePath)
	createMemoryDbSQL := `
	CREATE TABLE codon (
		codon TEXT PRIMARY KEY,
		aa TEXT
	);

	CREATE TABLE seq (
		pos INT PRIMARY KEY
	);

	CREATE TABLE history (
		pos INTEGER REFERENCES seq(pos),
		codon TEXT NOT NULL REFERENCES codon(codon),
		step INT,
		suggestedfix INT REFERENCES suggestedfix(id)
	);

	-- Weights are set on a per position basis for codon harmonization at a later point
	CREATE TABLE weights (
		pos INTEGER REFERENCES seq(pos),
		codon TEXT NOT NULL REFERENCES codon(codon),
		weight INTEGER
	);

	CREATE TABLE codonbias (
		fromcodon TEXT REFERENCES codon(codon),
		tocodon TEXT REFERENCES codon(codon),
		gcbias TEXT
	);

	CREATE TABLE suggestedfix (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		step INTEGER,
		start INTEGER REFERENCES seq(pos),
		end INTEGER REFERENCES seq(pos),
		gcbias TEXT,
		quantityfixes INTEGER,
		suggestiontype TEXT
	);
`
	db.MustExec(createMemoryDbSQL)
	// Insert codons
	weightTable := make(map[string]int)
	codonInsert := `INSERT INTO codon(codon, aa) VALUES (?, ?)`
	for _, aminoAcid := range codontable.AminoAcids {
		for _, codon := range aminoAcid.Codons {
			db.MustExec(codonInsert, codon.Triplet, aminoAcid.Letter)
			weightTable[codon.Triplet] = codon.Weight

			codonBias := strings.Count(codon.Triplet, "G") + strings.Count(codon.Triplet, "C")
			for _, toCodon := range aminoAcid.Codons {
				if codon.Triplet != toCodon.Triplet {
					toCodonBias := strings.Count(toCodon.Triplet, "G") + strings.Count(toCodon.Triplet, "C")
					switch {
					case codonBias == toCodonBias:
						db.MustExec(`INSERT INTO codonbias(fromcodon, tocodon, gcbias) VALUES (?, ?, ?)`, codon.Triplet, toCodon.Triplet, "NA")
					case codonBias > toCodonBias:
						db.MustExec(`INSERT INTO codonbias(fromcodon, tocodon, gcbias) VALUES (?, ?, ?)`, codon.Triplet, toCodon.Triplet, "AT")
					case codonBias < toCodonBias:
						db.MustExec(`INSERT INTO codonbias(fromcodon, tocodon, gcbias) VALUES (?, ?, ?)`, codon.Triplet, toCodon.Triplet, "GC")
					}
				}
			}
		}
	}

	// Insert seq and history
	pos := 0
	for i := 0; i < len(sequence); i = i + 3 {
		codon := sequence[i : i+3]
		db.MustExec(`INSERT INTO seq(pos) VALUES (?)`, pos)
		db.MustExec(`INSERT INTO history(pos, codon, step) VALUES (?, ?, 0)`, pos, codon)
		db.MustExec(`INSERT INTO weights(pos, codon, weight) VALUES (?,?,?)`, pos, codon, weightTable[codon])
		pos++
	}

	var err error
	// For a maximum of 100 iterations, see if we can do better. Usually sequences will be solved within 1-3 rounds,
	// so 100 just effectively acts as the max cap for iterations. Once you get to 100, you pretty much know that
	// we cannot fix the sequence.
	for i := 1; i < FixIterations; i++ {
		suggestions := findProblems(sequence, problematicSequenceFuncs)
		// If there are no suggestions, break the iteration!
		if len(suggestions) == 0 {
			// Add a historical log of changes
			var changes []Change
			_ = db.Select(&changes, `SELECT h.pos AS position, h.step AS step, (SELECT codon FROM history WHERE pos = h.pos AND step = h.step-1 LIMIT 1) AS codonfrom, h.codon AS codonto, sf.suggestiontype AS reason FROM history AS h JOIN suggestedfix AS sf ON sf.id = h.suggestedfix WHERE h.suggestedfix IS NOT NULL ORDER BY sf.id`)
			return sequence, changes, nil
		}
		for _, suggestion := range suggestions { // if you want to add overlaps, add suggestionIndex
			// First, let's insert the suggestions that we found using our problematicSequenceFuncs
			_, err = db.Exec(`INSERT INTO suggestedfix(step, start, end, gcbias, quantityfixes, suggestiontype) VALUES (?, ?, ?, ?, ?, ?)`, i, suggestion.Start, suggestion.End, suggestion.Bias, suggestion.QuantityFixes, suggestion.SuggestionType)
			if err != nil {
				return sequence, []Change{}, err
			}
		}

		// The following statements are the magic sauce that makes this all worthwhile.
		// Parameters: step, gcbias, start, end, quantityfix
		sqlFix1 := `INSERT INTO history
		            (codon,
		             pos,
		             step,
			     suggestedfix)
		SELECT t.codon,
		       t.pos,
		       ? AS step,
		       ? AS suggestedfix
		FROM   (SELECT cb.tocodon AS codon,
		               s.pos      AS pos
		        FROM   seq AS s
		               JOIN history AS h
		                 ON h.pos = s.pos
		               JOIN weights AS w
		                 ON w.pos = s.pos
		               JOIN codon AS c
		                 ON h.codon = c.codon
		               JOIN codonbias AS cb
		                 ON cb.fromcodon = c.codon
		        WHERE `
		sqlFix2 := ` s.pos >= ?
		               AND s.pos <= ?
		               AND h.codon != cb.tocodon
		        ORDER  BY w.weight) AS t
		GROUP  BY t.pos
		LIMIT  ?; `

		independentSuggestions := []dbDnaSuggestion{}
		_ = db.Select(&independentSuggestions, `SELECT * FROM suggestedfix WHERE step = ?`, i)

		for _, independentSuggestion := range independentSuggestions {
			switch independentSuggestion.Bias {
			case "NA":
				db.MustExec(sqlFix1+sqlFix2, i, independentSuggestion.ID, independentSuggestion.Start, independentSuggestion.End, independentSuggestion.QuantityFixes)
			case "GC":
				db.MustExec(sqlFix1+`cb.gcbias = 'GC' AND `+sqlFix2, i, independentSuggestion.ID, independentSuggestion.Start, independentSuggestion.End, independentSuggestion.QuantityFixes)
			case "AT":
				db.MustExec(sqlFix1+`cb.gcbias = 'AT' AND `+sqlFix2, i, independentSuggestion.ID, independentSuggestion.Start, independentSuggestion.End, independentSuggestion.QuantityFixes)
			}
		}
		var codons []string
		_ = db.Select(&codons, `SELECT codon FROM (SELECT codon, pos FROM history ORDER BY step DESC) GROUP BY pos`)
		sequence = strings.Join(codons, "")
	}
	return sequence, []Change{}, errors.New("Could not find a solution to sequence space")
}

// FixCdsSimple is FixCds with some defaults for normal usage, including
// removing of homopolymers and removing any repeat larger than 18 base pairs.
func FixCdsSimple(sequence string, codontable codon.Table, sequencesToRemove []string) (string, []Change, error) {
	var functions []func(string, chan DnaSuggestion, *sync.WaitGroup)
	// Remove homopolymers
	functions = append(functions, RemoveSequence([]string{"AAAAAAAA", "GGGGGGGG"}, "Homopolymers"))

	// Remove user defined sequences
	functions = append(functions, RemoveSequence(sequencesToRemove, "Removal requested by user"))

	// Remove repeats
	functions = append(functions, RemoveRepeat(18))

	// Ensure normal GC range

	return FixCds(":memory:", sequence, codontable, functions)
}
