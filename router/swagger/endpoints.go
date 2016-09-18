// Package docs generates endpoint documentation.
package docs

import (
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/blueprint/blueprint/router"
)

// Endpoints generates API documentation for API HTTP route endpoints.
// The documentation is written to /public/docs/api/index.html.
func Endpoints(tmpl, file string, endpoints []router.Endpoint) error {

	// ensure directories exist
	if err := os.MkdirAll(path.Dir(file), 0777); err != nil {
		return err
	}

	// create the file
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	// defer close
	defer f.Close()

	// create a new template
	t := template.New("endpoints")

	// read/parse the template
	t, err = t.ParseFiles(tmpl)
	if err != nil {
		return err
	}

	data := struct {
		Endpoints []router.Endpoint
	}{
		Endpoints: endpoints,
	}

	// get filename
	tmplbase := filepath.Base(tmpl)
	filename := tmplbase[:]

	// write file to disk
	if err = t.ExecuteTemplate(f, filename, data); err != nil {
		return err
	}

	return nil
}
