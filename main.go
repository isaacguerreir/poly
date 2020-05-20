package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type Meta struct {
	// shared
	Name        string
	GffVersion  string
	RegionStart int
	RegionEnd   int
	// genbank specific
	Size            int
	Type            string
	GenbankDivision string
	Date            string
	Definition      string
	Accession       string
	Version         string
	Keywords        string
	Organism        string
	Source          string
	Origin          string
	Locus           Locus
	References      []Reference
	Primaries       []Primary
}

type Primary struct {
	RefSeq, PrimaryIdentifier, Primary_Span, Comp string
}

// genbank specific
// type Reference struct {
// 	Authors []string
// 	Title   string
// 	Journal string
// 	PubMed  string
// }

type Reference struct {
	Index, Authors, Title, Journal, PubMed, Remark string
}

type Locus struct {
	Name, SequenceLength, MoleculeType, GenBankDivision, ModDate string
	Circular                                                     bool
}

// from https://github.com/blachlylab/gff3/blob/master/gff3.go
type Feature struct {
	Name string //Seqid in gff, name in gbk
	//gff specific
	Source     string
	Type       string
	Start      int
	End        int
	Score      float64
	Strand     byte
	Phase      int
	Attributes map[string]string // Known as "qualifiers" for gbk, "attributes" for gff.
	//gbk specific
	Location string
	Sequence string
}

type Sequence struct {
	Description string
	Sequence    string
}

type AnnotatedSequence struct {
	Meta     Meta
	Features []Feature
	Sequence Sequence
}

func parseGbk(path string) {

	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}

	bio := bufio.NewReader(f)

	var lines []string
	i := 0

	// Read all lines of the file into buffer
	for {
		line, err := bio.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		sline := strings.TrimRight(string(line), `\n`)
		lines = append(lines, sline)
		i++
	}
	// End read of file into buffer

	// Create meta struct
	meta := Meta{}

	// Create features struct
	// features := []Feature{}

	// Create sequence struct
	// sequence := Sequence{}

	for numLine, line := range lines {
		// fmt.Print / ln(numLine)
		splitLine := strings.Split(line, " ")
		subLines := lines[numLine+1:]

		switch splitLine[0] {

		case "":
			continue
		case "LOCUS":
			meta.Locus = parseLocus(line)
		case "DEFINITION":
			meta.Definition = joinSubLines(splitLine, subLines)
		case "ACCESSION":
			meta.Accession = joinSubLines(splitLine, subLines)
		case "VERSION":
			meta.Version = joinSubLines(splitLine, subLines)
		case "KEYWORDS":
			meta.Keywords = joinSubLines(splitLine, subLines)
		case "SOURCE":
			meta.Source, meta.Organism = getSourceOrganism(splitLine, subLines)
		case "REFERENCE":
			meta.References = append(meta.References, getReference(splitLine, subLines))
			continue
		case "FEATURES":
			continue
		case "ORIGIN":
			continue
		default:
			continue
		}

	}
	file, _ := json.MarshalIndent(meta, "", " ")

	_ = ioutil.WriteFile("test.json", file, 0644)

}

func main() {

	// fmt.Println(parseGff("data/ecoli-mg1655.gff"))
	// parseGbk("data/addgene-plasmid-50005-sequence-74677.gbk")
	parseGbk("data/test.gbk")
}
