// Generated from text/text_test.go; DO NOT EDIT

package html

import (
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/halimath/decor"
)

func TestFilesConfig_TemplatePath(t *testing.T) {
	c := FilesConfig{
		TemplatesPattern: "%s.tpl",
		BasePath:         "foo/bar",
	}

	act := c.TemplatePath("spam")

	if act != "foo/bar/spam.tpl" {
		t.Errorf("unexpected path: '%s'", act)
	}
}

func TestFilesConfig_TemplatePaths(t *testing.T) {
	c := FilesConfig{
		TemplatesPattern: "%s.tpl",
		BasePath:         "foo/bar",
		IncludePaths: []string{
			"inc/1.tpl",
			"inc/2.tpl",
		},
	}

	act := c.TemplatePaths("spam")

	if diff := deep.Equal(act, []string{
		"foo/bar/spam.tpl",
		"foo/bar/inc/1.tpl",
		"foo/bar/inc/2.tpl",
	}); diff != nil {
		t.Error(diff)
	}
}

func TestNewFilesLoader(t *testing.T) {
	tpls := decor.Templates{
		Loader: NewFilesLoader(FilesConfig{
			TemplatesPattern: "%s.txt",
			BasePath:         "../testtemplates",
			IncludePaths: []string{
				"layouts/base.txt",
			},
		}),
	}

	var w strings.Builder

	if err := tpls.ExecuteTemplate(&w, "a", "world"); err != nil {
		t.Fatal(err)
	}

	if w.String() != "Hello, world!" {
		t.Errorf("unexpected output: '%s'", w.String())
	}
}
