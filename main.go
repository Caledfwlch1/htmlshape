package main

import (
	"flag"
	"golang.org/x/net/html"
	"os"
	"strings"
	"log"
	"io"
	"fmt"
	"bufio"
)

var f_out = flag.String("o","out.gv.txt","output graphviz file")
var f_in = flag.String("i","input.html","input html file")
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
	//if label != "" {
	//	fmt.Fprintf(w, ` [ label = %s ]`, str(label))
	//}
	fmt.Fprint(w, `;`)
}

func walk(w io.Writer, n *html.Node) {
	if n == nil {
		return
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode && strings.TrimSpace(c.Data) == "" {
			continue
		}
		link(w, n, c, "child")
		walk(w, c)
	}
}

func main() {
	flag.Parse()
	inp := `<html>
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
	inp, err := readFile(*f_in)
	if err != io.EOF {
		log.Fatal("Error reading file ", *f_in, " - ", err)
		return
	}
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

func readFile(in string) (sl string, err error) {

	_, err = os.Stat(in)
	if err != nil {
		return "", err
	}
	fi, err := os.Open(in)
	if err != nil {
		return "", err
	}
	fiReader := bufio.NewReader(fi)
	var s []byte
	for err == nil {
		s, _, err = fiReader.ReadLine()
		sl += string(s)
	}

	return sl, err
}