package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	source = flag.String("s", path.Join(".", "templates"), "Location of templates")
	output = flag.String("o", "", "Output file")
)

func main() {
	flag.Parse()

	buf := new(bytes.Buffer)
	fmt.Fprint(buf, `package templates

	import (
		"text/template"
		"strings"
	)

  	var templates = map[string]string{`)

	if err := filepath.Walk(*source, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Ignore non-templates files
		if filepath.Ext(path) != ".tmpl" {
			return nil
		}

		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		s := strings.ReplaceAll(string(b), "`", "__TEMP_BACKTICK__")
		fmt.Fprintf(buf, "\"%s\": `%s`,\n", filepath.Base(path), s)

		return nil
	}); err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(buf, `}

	func Parse(t *template.Template) (*template.Template, error) {
		for name, s := range templates {
			var tmpl *template.Template
			if t == nil {
				t = template.New(name)
			}
			if name == t.Name() {
				tmpl = t
			} else {
				tmpl = t.New(name)
			}
			ss := strings.ReplaceAll(s, "__TEMP_BACKTICK__", "%s")
			if _, err := tmpl.Parse(ss); err != nil {
				return nil, err
			}
		}
		return t, nil
	}`, "`")

	clean, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	file := os.Stdout
	if *output != "" {
		file, err = os.Create(*output)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Fprintln(file, string(clean))
}
