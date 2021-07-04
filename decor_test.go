package decor

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

type tpl struct {
	name string
}

func (t *tpl) Execute(w io.Writer, data interface{}) error {
	if t.name == "error" {
		return errors.New("error")
	}

	_, err := fmt.Fprintf(w, "%s: %#v", t.name, data)
	return err
}

var _ Template = &tpl{}

// --

type ldr struct {
	loadCount map[string]int
}

func (l *ldr) Load(name string, funcs FuncMap) (Template, error) {
	l.loadCount[name] = l.loadCount[name] + 1

	return &tpl{name}, nil
}

var _ Loader = &ldr{}

func TestRepository_develMode(t *testing.T) {
	l := &ldr{
		loadCount: make(map[string]int),
	}
	tpls := Templates{
		Loader:    l,
		DevelMode: true,
	}

	var w strings.Builder

	for i := 0; i < 10; i++ {
		tpls.ExecuteTemplate(&w, "foo", nil)
	}

	if l.loadCount["foo"] != 10 {
		t.Errorf("expected load count of 10 but got %d", l.loadCount["foo"])
	}
}

func TestRepository_nonDevelMode(t *testing.T) {
	l := &ldr{
		loadCount: make(map[string]int),
	}
	tpls := Templates{
		Loader: l,
	}

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var w strings.Builder
			name := "foo"
			if i%2 == 0 {
				name = "bar"
			}
			tpls.ExecuteTemplate(&w, name, nil)
		}(i)
	}

	wg.Wait()

	if l.loadCount["foo"] != 1 {
		t.Errorf("expected load count of 1 but got %d", l.loadCount["foo"])
	}

	if l.loadCount["bar"] != 1 {
		t.Errorf("expected load count of 1 but got %d", l.loadCount["bar"])
	}
}

func TestRepository_SendHTML_develMode(t *testing.T) {
	l := &ldr{
		loadCount: make(map[string]int),
	}
	tpls := Templates{
		Loader:    l,
		DevelMode: true,
	}

	var body bytes.Buffer
	w := httptest.ResponseRecorder{
		Body: &body,
	}

	tpls.SendHTML(&w, "foo", "hello, world")

	if body.String() != `foo: "hello, world"` {
		t.Errorf("unexpected body: '%s'", body.String())
	}

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("unexpected status: %d", w.Result().StatusCode)
	}

	if w.Result().Header.Get("Content-Type") != "text/html" {
		t.Errorf("unexpected content type: %s", w.Result().Header.Get("Content-Type"))
	}
}

func TestRepository_SendHTML_develMode_error(t *testing.T) {
	l := &ldr{
		loadCount: make(map[string]int),
	}
	tpls := Templates{
		Loader:    l,
		DevelMode: true,
	}

	var body bytes.Buffer
	w := httptest.ResponseRecorder{
		Body: &body,
	}

	tpls.SendHTML(&w, "error", nil)

	if body.String() != `Failed to render template 'error': error` {
		t.Errorf("unexpected body: '%s'", body.String())
	}

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("unexpected status: %d", w.Result().StatusCode)
	}

	if w.Result().Header.Get("Content-Type") != "text/plain" {
		t.Errorf("unexpected content type: %s", w.Result().Header.Get("Content-Type"))
	}
}

func TestRepository_SendHTML_nonDevelMode(t *testing.T) {
	l := &ldr{
		loadCount: make(map[string]int),
	}
	tpls := Templates{
		Loader: l,
	}

	var body bytes.Buffer
	w := httptest.ResponseRecorder{
		Body: &body,
	}

	tpls.SendHTML(&w, "foo", "hello, world")

	if body.String() != `foo: "hello, world"` {
		t.Errorf("unexpected body: '%s'", body.String())
	}

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("unexpected status: %d", w.Result().StatusCode)
	}

	if w.Result().Header.Get("Content-Type") != "text/html" {
		t.Errorf("unexpected content type: %s", w.Result().Header.Get("Content-Type"))
	}
}

func TestRepository_SendHTML_nonDevelMode_error(t *testing.T) {
	l := &ldr{
		loadCount: make(map[string]int),
	}
	tpls := Templates{
		Loader: l,
	}

	var body bytes.Buffer
	w := httptest.ResponseRecorder{
		Body: &body,
	}

	tpls.SendHTML(&w, "error", nil)

	if len(body.Bytes()) != 0 {
		t.Errorf("unexpected body: '%s'", body.String())
	}

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("unexpected status: %d", w.Result().StatusCode)
	}

	if w.Result().Header.Get("Content-Type") != "text/html" {
		t.Errorf("unexpected content type: %s", w.Result().Header.Get("Content-Type"))
	}
}
