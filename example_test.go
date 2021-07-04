package decor_test

import (
	"log"
	"os"

	"github.com/halimath/decor"
	"github.com/halimath/decor/text"
)

func Example_layout() {
	tpls := decor.Templates{
		Loader: text.NewFilesLoader(text.FilesConfig{
			TemplatesPattern: "%s.txt",
			BasePath:         "testtemplates",
			IncludePaths: []string{
				"layouts/base.txt",
			},
		}),
	}

	if err := tpls.ExecuteTemplate(os.Stdout, "a", "world"); err != nil {
		log.Fatal(err)
	}

	// Output:
	// Hello, world!
}
