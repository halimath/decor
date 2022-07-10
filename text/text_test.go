package text

import (
	"strings"
	"testing"

	"github.com/halimath/decor"
)

func TestNewFilesLoader(t *testing.T) {
	tpls := decor.Templates{
		Includes: []string{"layouts/base"},
		Loader:   NewFilesLoader("%s.txt", "../testtemplates"),
	}

	var w strings.Builder

	if err := tpls.ExecuteTemplate(&w, "a", "world"); err != nil {
		t.Fatal(err)
	}

	if w.String() != "Hello, world!" {
		t.Errorf("unexpected output: '%s'", w.String())
	}
}
