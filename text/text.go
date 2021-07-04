// Package text contains a layout based loader for text/template.
package text

import (
	"fmt"
	"path"
	"text/template"

	"github.com/halimath/decor"
)

// FilesConfig implements the configuration used by the FSLoader to load templates
// from files.
type FilesConfig struct {
	// TemplatesPattern is a string template containing a single %s placeholder to
	// be replaced with the template name to load and should resolve to the file
	// name for loading the named template.
	TemplatesPattern string

	// BasePath is used as the root for all templates to load, both by pattern
	// and by includes.
	BasePath string

	// IncludePaths defines a list of static paths to include with every template
	// to load.
	IncludePaths []string
}

// TemplatePath resolves the templateName to a path to load the corresponding template from.
func (c *FilesConfig) TemplatePath(templateName string) string {
	return path.Join(c.BasePath, fmt.Sprintf(c.TemplatesPattern, templateName))
}

// TemplatePaths resolves the templateName to a path to load the corresponding template from.
func (c *FilesConfig) TemplatePaths(templateName string) []string {
	r := make([]string, len(c.IncludePaths)+1)
	r[0] = c.TemplatePath(templateName)

	for i, p := range c.IncludePaths {
		r[i+1] = path.Join(c.BasePath, p)
	}

	return r
}

type filesLoader struct {
	c FilesConfig
}

// NewFilesLoader creates a new loader loading templates from files with c as the configuration.
func NewFilesLoader(c FilesConfig) decor.Loader {
	return &filesLoader{
		c: c,
	}
}

var _ decor.Loader = &filesLoader{}

func (l *filesLoader) Load(templateName string, f decor.FuncMap) (decor.Template, error) {
	t, err := template.ParseFiles(l.c.TemplatePaths(templateName)...)
	if err != nil {
		return nil, err
	}

	if f != nil {
		t = t.Funcs(template.FuncMap(f))
	}

	return t, nil
}
