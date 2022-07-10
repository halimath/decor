// Package text contains a layout based loader for text/template.
package text

import (
	"fmt"
	"io/fs"
	"os"
	"text/template"

	"github.com/halimath/decor"
)

type fsLoader struct {
	pattern string
	fsys    fs.FS
}

func NewFSLoader(templateNamePattern string, fsys fs.FS) decor.Loader {
	return &fsLoader{
		pattern: templateNamePattern,
		fsys:    fsys,
	}
}

func NewFilesLoader(templateNamePattern, rootDir string) decor.Loader {
	return NewFSLoader(templateNamePattern, os.DirFS(rootDir))
}

func (l *fsLoader) Load(names []string, funcs decor.FuncMap) (decor.Template, error) {
	paths := make([]string, len(names))
	for i, n := range names {
		paths[i] = fmt.Sprintf(l.pattern, n)
	}

	t, err := template.ParseFS(l.fsys, paths...)
	if err != nil {
		return nil, err
	}

	if funcs != nil {
		t = t.Funcs(template.FuncMap(funcs))
	}

	return t, nil
}
