package decor

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

// FuncMap defines a map of template functions.
type FuncMap map[string]interface{}

// Template defines an internal abstraction for templates provided either by package text/template
// or html/template or a custom template rendering engine.
type Template interface {
	// Execute executes the template rendering content using data and
	// writing the output to wr. It returns an error in case of an error
	// or nil to indicate success.
	Execute(wr io.Writer, data interface{}) error
}

// Loader defines the interface for implementations that load Template values
// from a source. Templates are loaded by name.
type Loader interface {
	// Loads the template named name and returns it. If funcs is a non-nil, non-empty map
	// the contained functions must be made available to the template.
	// Returns a non-nil error to indicate that loading failed.
	Load(name string, funcs FuncMap) (Template, error)
}

// Templates implements loading, caching and execution of named templates.
type Templates struct {
	// When set to true DevelMode puts these templates into development mode, which
	// causes errors during rendering being written to the target writer as well
	// as templates being reload on every execution.
	DevelMode bool

	// Funcs contains a map of functions to be used in templates.
	// These are passed to every template loaded using the given
	// Templates.
	Funcs FuncMap

	// Loader is used to load Templates when requested.
	Loader Loader

	lock  sync.RWMutex
	cache map[string]Template
}

// ExecuteTemplate executes the template named templateName with data and writes the output to w.
// It returns any non-nil error returned from the loader or from Template.Execute.
func (t *Templates) ExecuteTemplate(w io.Writer, templateName string, data interface{}) error {
	tpl, err := t.loadTemplate(templateName)
	if err != nil {
		return err
	}

	return tpl.Execute(w, data)
}

// SendHTML executes the template named templateName with the given data (using Repository.ExecuteTemplate).
// It writes the rendered output to w after setting the responsive HTTP content headers.
// This method pays attention to the r's DevelMode: When activated any error produced from rendering the
// template will be send as the http response. When DevelopMode is false, any error will be suppressed
// possibly resulting in an empty response.
func (t *Templates) SendHTML(w http.ResponseWriter, templateName string, data interface{}) {
	var buf bytes.Buffer
	err := t.ExecuteTemplate(&buf, templateName, data)

	if err != nil {
		if t.DevelMode {
			w.Header().Set("content-type", "text/plain")
			fmt.Fprintf(w, "Failed to render template '%s': %s", templateName, err)
			return
		}

		log.Printf("Error rendering template %s: %s", templateName, err)
	}

	header := w.Header()
	header.Set("Content-Type", "text/html")
	header.Set("Content-Length", fmt.Sprintf("%d", buf.Len()))

	w.Write(buf.Bytes())
}

// loadTemplate loades the named template using the configured loader.
// When DevelMode is set to true, this method always delegates to r's
// template loader. When set to false, templates are loaded on first
// invocation and are then stored and reused from a synchronized cache.
func (t *Templates) loadTemplate(name string) (tpl Template, err error) {
	if t.DevelMode {
		tpl, err = t.Loader.Load(name, t.Funcs)
		return
	}

	t.lock.RLock()

	var ok bool
	if t.cache != nil {
		tpl, ok = t.cache[name]
	}

	t.lock.RUnlock()

	if ok {
		return
	}

	t.lock.Lock()
	defer t.lock.Unlock()

	if t.cache == nil {
		t.cache = make(map[string]Template)
	} else {
		if tpl, ok = t.cache[name]; ok {
			return
		}
	}

	tpl, err = t.Loader.Load(name, t.Funcs)
	if err != nil {
		return
	}

	t.cache[name] = tpl

	return
}
