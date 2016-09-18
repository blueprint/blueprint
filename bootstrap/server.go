package bootstrap

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/blueprint/blueprint/dba"
	"github.com/blueprint/blueprint/docs"
	"github.com/blueprint/blueprint/http/middleware"
	"github.com/blueprint/blueprint/router"
	"github.com/spf13/viper"
)

// Server takes a register function and bootstraps a server.
func Server(r router.Router, register func(r router.Router)) error {

	// bootstrap environment and configuration settings
	if err := Config(); err != nil {
		return err
	}

	// bootstrap database admin
	if err := DBA(); err != nil {
		return err
	}

	// open all db connections
	for _, db := range dba.All() {
		if err := db.Dial(db.Name()); err != nil {
			return err
		}
	}

	// defer closing dba connections
	for _, db := range dba.All() {
		defer db.Close()
	}

	// set api keys as middleware
	keys := viper.GetStringMapStringSlice("keys")
	if len(keys) > 0 {
		km := middleware.NewKeys(keys)
		r.Middleware(km.Do)
	}

	register(r)

	port := viper.GetString("port")
	static := viper.GetString("static")
	tls := viper.GetStringMapString("tls")

	// if a static directory path is provided, register it
	if len(static) > 0 {
		r.Static("/public", static)
	} else {
		r.Static("/public", string(os.PathSeparator)+"public")
	}

	// generate docs
	if err := Docs(static, router.Endpoints()); err != nil {
		return err
	}

	// prep port and green output
	portStr := fmt.Sprintf(":%s", port)

	// let the humans know we are serving...
	if tls["cert"] != "" && tls["key"] != "" {
		fmt.Printf("https listening and serving with TLS on port :%s\n", port)
		return http.ListenAndServeTLS(portStr, tls["cert"], tls["key"], r)
	}

	fmt.Printf("http listening and serving on port :%s\n", port)
	return http.ListenAndServe(portStr, r)
}

// Docs renders all the endpoint docs for the API application service.
func Docs(static string, endpoints []router.Endpoint) error {
	tmpl := filepath.Join(filepath.Dir(static), "templates", "endpoints.tmpl")
	html := filepath.Join(static, "docs", "api", "index.html")
	if err := docs.Endpoints(tmpl, html, endpoints); err != nil {
		return err
	}
	return nil
}
