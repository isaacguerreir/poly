(window.webpackJsonp=window.webpackJsonp||[]).push([[11],{110:function(e,t,n){"use strict";n.r(t),n.d(t,"frontMatter",(function(){return o})),n.d(t,"metadata",(function(){return l})),n.d(t,"rightToc",(function(){return c})),n.d(t,"default",(function(){return p}));var a=n(2),r=n(6),i=(n(0),n(136)),o={id:"installation",title:"Installing Poly"},l={id:"installation",isDocsHomePage:!0,title:"Installing Poly",description:"Poly can be used in two ways.",source:"@site/docs/installation.md",permalink:"/polydocs/docs/",editUrl:"https://github.com/timothystiles/poly/edit/prime/docs/installation.md",sidebar:"someSidebar",next:{title:"Converting Sequence Files",permalink:"/polydocs/docs/cli-converting"}},c=[{value:"Installing Poly as a Go library and Command Line Utility",id:"installing-poly-as-a-go-library-and-command-line-utility",children:[]},{value:"Installing Poly as a Command Line Utility",id:"installing-poly-as-a-command-line-utility",children:[{value:"Mac OS",id:"mac-os",children:[]},{value:"Linux",id:"linux",children:[]},{value:"Windows",id:"windows",children:[]}]},{value:"Building Poly from Scratch",id:"building-poly-from-scratch",children:[]}],s={rightToc:c};function p(e){var t=e.components,n=Object(r.a)(e,["components"]);return Object(i.b)("wrapper",Object(a.a)({},s,n,{components:t,mdxType:"MDXLayout"}),Object(i.b)("p",null,"Poly can be used in two ways."),Object(i.b)("ol",null,Object(i.b)("li",{parentName:"ol"},"As a Go library where you have finer control and can make magical things happen."),Object(i.b)("li",{parentName:"ol"},"As a command line utility where you can bash script your way to greatness and make DNA go brrrrrrrr.")),Object(i.b)("h2",{id:"installing-poly-as-a-go-library-and-command-line-utility"},"Installing Poly as a Go library and Command Line Utility"),Object(i.b)("p",null,"This assumes you already have a working Go environment, if not please see\n",Object(i.b)("a",Object(a.a)({parentName:"p"},{href:"https://golang.org/doc/install"}),"this page")," first."),Object(i.b)("p",null,Object(i.b)("inlineCode",{parentName:"p"},"go get")," ",Object(i.b)("em",{parentName:"p"},"will always pull the latest released version from the prime branch.")),Object(i.b)("pre",null,Object(i.b)("code",Object(a.a)({parentName:"pre"},{className:"language-bash"}),"go get github.com/TimothyStiles/poly\n")),Object(i.b)("h2",{id:"installing-poly-as-a-command-line-utility"},"Installing Poly as a Command Line Utility"),Object(i.b)("p",null,"Poly ships many binaries for many different operating systems and package managers thanks to the wonderful work of the ",Object(i.b)("a",Object(a.a)({parentName:"p"},{href:"https://goreleaser.com/"}),"go releaser")," team. You can check out our ",Object(i.b)("a",Object(a.a)({parentName:"p"},{href:"https://github.com/TimothyStiles/poly/releases"}),"releases page")," on github or install via package manager for your OS with the instructions below."),Object(i.b)("h3",{id:"mac-os"},"Mac OS"),Object(i.b)("pre",null,Object(i.b)("code",Object(a.a)({parentName:"pre"},{className:"language-bash"}),"brew install timothystiles/poly/poly\n")),Object(i.b)("h3",{id:"linux"},"Linux"),Object(i.b)("pre",null,Object(i.b)("code",Object(a.a)({parentName:"pre"},{className:"language-bash"}),"sudo snap install --classic gopoly\n")),Object(i.b)("h3",{id:"windows"},"Windows"),Object(i.b)("p",null,Object(i.b)("a",Object(a.a)({parentName:"p"},{href:"https://github.com/TimothyStiles/poly/issues/16"}),"Coming soon...")),Object(i.b)("h2",{id:"building-poly-from-scratch"},"Building Poly from Scratch"),Object(i.b)("p",null,"This assumes you already have a working Go environment, if not please see\n",Object(i.b)("a",Object(a.a)({parentName:"p"},{href:"https://golang.org/doc/install"}),"this page")," first."),Object(i.b)("pre",null,Object(i.b)("code",Object(a.a)({parentName:"pre"},{className:"language-bash"}),"git clone https://github.com/TimothyStiles/poly.git && cd poly && go build && go install\n")))}p.isMDXComponent=!0},136:function(e,t,n){"use strict";n.d(t,"a",(function(){return b})),n.d(t,"b",(function(){return d}));var a=n(0),r=n.n(a);function i(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function o(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);t&&(a=a.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,a)}return n}function l(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?o(Object(n),!0).forEach((function(t){i(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):o(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function c(e,t){if(null==e)return{};var n,a,r=function(e,t){if(null==e)return{};var n,a,r={},i=Object.keys(e);for(a=0;a<i.length;a++)n=i[a],t.indexOf(n)>=0||(r[n]=e[n]);return r}(e,t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(a=0;a<i.length;a++)n=i[a],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(r[n]=e[n])}return r}var s=r.a.createContext({}),p=function(e){var t=r.a.useContext(s),n=t;return e&&(n="function"==typeof e?e(t):l(l({},t),e)),n},b=function(e){var t=p(e.components);return r.a.createElement(s.Provider,{value:t},e.children)},u={inlineCode:"code",wrapper:function(e){var t=e.children;return r.a.createElement(r.a.Fragment,{},t)}},m=r.a.forwardRef((function(e,t){var n=e.components,a=e.mdxType,i=e.originalType,o=e.parentName,s=c(e,["components","mdxType","originalType","parentName"]),b=p(n),m=a,d=b["".concat(o,".").concat(m)]||b[m]||u[m]||i;return n?r.a.createElement(d,l(l({ref:t},s),{},{components:n})):r.a.createElement(d,l({ref:t},s))}));function d(e,t){var n=arguments,a=t&&t.mdxType;if("string"==typeof e||a){var i=n.length,o=new Array(i);o[0]=m;var l={};for(var c in t)hasOwnProperty.call(t,c)&&(l[c]=t[c]);l.originalType=e,l.mdxType="string"==typeof e?e:a,o[1]=l;for(var s=2;s<i;s++)o[s]=n[s];return r.a.createElement.apply(null,o)}return r.a.createElement.apply(null,n)}m.displayName="MDXCreateElement"}}]);