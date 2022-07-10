package decor_test

import (
	"log"
	"os"

	"github.com/halimath/decor"
	"github.com/halimath/decor/text"
)

func Example_layout() {
	tpls := decor.Templates{
		Includes: []string{
			"layouts/base",
		},
		Loader: text.NewFilesLoader("%s.txt", "testtemplates"),
	}

	if err := tpls.ExecuteTemplate(os.Stdout, "a", "world"); err != nil {
		log.Fatal(err)
	}

	// Output:
	// Hello, world!
}
