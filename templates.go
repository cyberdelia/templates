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
)

var source = flag.String("source", path.Join(".", "templates"), "Location of templates")

func main() {
	flag.Parse()

	buf := new(bytes.Buffer)
	fmt.Fprint(buf, `package templates

	import "text/template"

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

		fmt.Fprintf(buf, "\"%s\": `%s`,\n", filepath.Base(path), b)

		return nil
	}); err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(buf, `}

	// Parse parses declared templates.
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
	  		if _, err := tmpl.Parse(s); err != nil {
  				return nil, err
  			}
  		}
  		return t, nil
  	}`)

	clean, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(clean))
}
