package main

import (
	"flag"
	"golang.org/x/net/html"
	"os"
	"strings"
	"log"
	"io"
	"fmt"
)

var f_out = flag.String("o","out.gv.txt","output graphviz file")
var escaper = strings.NewReplacer(`"`,`\"`)

func str(s string) string {
	return `"`+escaper.Replace(s)+`"`
}

var nodes = make(map[*html.Node]int)

func nodeStr(n *html.Node) string {
	id, ok := nodes[n]
	if !ok {
		id = len(nodes)
		nodes[n] = id
	}
	return str(fmt.Sprintf("%d: %s", id, n.Data))
}

func link(w io.Writer, n1, n2 *html.Node, label string) {
	if n1 == nil || n2 == nil {
		return
	}
	fmt.Fprintf(w, `%s -> %s`, nodeStr(n1), nodeStr(n2))
	if label != "" {
		fmt.Fprintf(w, ` [ label = %s ]`, str(label))
	}
	fmt.Fprint(w, `;`)
}

func walk(w io.Writer, n *html.Node) {
	if n == nil {
		return
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		link(w, n, c, "child")
		walk(w, c)
	}
}

func main() {
	flag.Parse()
	const inp = `<html>
	<head>
	<title>
		Title of the document
	</title>
	</head>
	<body>
	<title>Title</title>
	<p>bla</p>
	<div>
	<a href="#">link</a>
	</div>
	</body>
	</html>`
	n, err := html.Parse(strings.NewReader(inp))
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Create(*f_out)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fmt.Fprint(file,"digraph html {\n")
	walk(file, n)
	fmt.Fprint(file,"}\n")
}
