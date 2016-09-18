package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/target/gophersaurus/bootstrap"
	"github.com/target/gophersaurus/http"
	"github.com/target/gophersaurus/router"
)

// Define the root and serve commands.
var (
	RootCmd  = &cobra.Command{}
	ServeCmd = &cobra.Command{
		Use:     "serve",
		Aliases: []string{"server", "s"},
		Short:   "Start HTTP Server",
		Long:    "Start An HTTP Server To Listen and Serve over HTTP",
		Run: func(cmd *cobra.Command, args []string) {
			log.Fatal(bootstrap.Server(router.NewMux(), register))
		},
	}
)

// register takes a router and assocates HTTP endpoints to methods.
func register(r router.Router) {
	r.GET("/", home)
}

// home is an index method.
func home(resp http.ResponseWriter, req *http.Request) {
	result := struct {
		Status        int    `json:"status" xml:"status"`
		Message       string `json:"message" xml:"message"`
		PublicPage    string `json:"public_page" xml:"public_page"`
		PublicAPIDocs string `json:"public_api_docs" xml:"public_api_docs"`
	}{
		200,
		"Welcome fellow gopher.",
		"http://" + req.Host + "/public",
		"http://" + req.Host + "/public/docs/api",
	}
	resp.WriteFormat(req, result)
}

func main() {
	RootCmd.AddCommand(ServeCmd)
	RootCmd.Execute()
}
