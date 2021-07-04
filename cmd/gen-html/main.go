// Main contains a helper application to generate the html/template based implementation by
// copying the text/template based sources.
package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

func main() {
	if err := os.RemoveAll("html"); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll("html", 0755); err != nil {
		log.Fatal(err)
	}

	files, err := os.ReadDir("text")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".go") {
			return
		}

		b, err := os.ReadFile(path.Join("text", f.Name()))
		if err != nil {
			log.Fatal(err)
		}

		src := string(b)
		src = strings.ReplaceAll(src, `"text/template"`, `"html/template"`)
		src = strings.ReplaceAll(src, `package text`, `package html`)
		src = fmt.Sprintf("// Generated from text/%s; DO NOT EDIT\n\n%s", f.Name(), src)

		if err := os.WriteFile(path.Join("html", strings.ReplaceAll(f.Name(), "text", "html")), []byte(src), 0644); err != nil {
			log.Fatal(err)
		}
	}
}
