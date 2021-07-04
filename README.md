# decor

![CI Status][ci-img-url] [![Go Report Card][go-report-card-img-url]][go-report-card-url] [![Package Doc][package-doc-img-url]][package-doc-url] [![Releases][release-img-url]][release-url]

`decor` contains a simple-to-user yet powerful template loader
for [Golang](https://golang.org) mainly focussed on the built-in 
packages `text/template` and `html/template`.

While these packages provide powerful templating facilities with
great support for rendering textual content, using these packages
to build a complete layout-based template solution involves a lot
of boilerplate code. decor tries to fill this gap by providing just
what's needed to reduce the boilerplate and enable a great developer
experience.

## Installation

```
$ go get github.com/halimath/decor
```

Decor requires go >= 1.14.

## Usage

The core components provided by `decor` is the `Templates` struct which
provides the rendering functionality of templates. When creating an instance
you need to provide a `Loader` which is responsible for loading a named template.

`decor` provides a files loader for both `text` and `html` templates that includes
capabilities to load layout files as well.

```go
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
```

The above snippet configures a `Templates` value and effectively executes something like.

```go
t, _ := "text/template".ParseFiles("testtemplates/a.txt", "testtemplates/layouts/base.txt")
t.Execute(os.Stdout, "world")
```

Note that error handling as well as template caching is omitted from the example; the two lines above
only sketch what's happening.

### Caching Templates

By default templates loaded are put into a cache the first time they are loaded. The `Templates` struct
is safe to use across multiple goroutines, so loading and rendering template in parallel works out of the
box. Internally, `Templates` uses an efficient `sync.RWLock` to guard access to the thread.

You can put `Templates` into _development mode_ by setting its `DevelMode`-Field to `true`. This 
completely disables the cache and templates are loaded using the loader _every time_ they are executed. 
This is great during development when changes made to templates become immediately effective.

### Rendering HTTP Response

Its a common case to render a template producing HTML as part of a web server application using the 
`net/http` package. `decor` provides a utility function to handle this case as a one liner:
`Templates.SendHTML`. The method accepts a `http.ResponseWriter`, a template name and template data
and directly writes any output rendered from template to the HTTP response. The method also sets the
content type to `text/html` as well as the content length.

When in development mode (see above), any errors generated during rendering will be send to the HTTP 
client for debugging. When not in development mode, errors will be logged and the response will be empty.

## Development

The files in package `html` are generated from the corresponding files in package `text`. The code 
generation tool can be found in `cmd/gen-html`. In order to run this tool, you need Go >= 1.16. The
resulting code will work with Go >= 1.14.

## License

Copyright 2021 Alexander Metzner
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

[ci-img-url]: https://github.com/halimath/decor/workflows/CI/badge.svg
[go-report-card-img-url]: https://goreportcard.com/badge/github.com/halimath/decor
[go-report-card-url]: https://goreportcard.com/report/github.com/halimath/decor
[package-doc-img-url]: https://img.shields.io/badge/GoDoc-Reference-blue.svg
[package-doc-url]: https://pkg.go.dev/github.com/halimath/decor
[release-img-url]: https://img.shields.io/github/v/release/halimath/decor.svg
[release-url]: https://github.com/halimath/decor/releases
