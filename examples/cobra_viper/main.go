package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

func init() {
	// bind port flag
	ServeCmd.Flags().IntP("port", "p", 8080, "Port to run Application server on")
	viper.BindPFlag("port", ServeCmd.Flags().Lookup("port"))

	// bind static flag
	ServeCmd.Flags().StringP("static", "s", "public", "Where the public static files are")
	viper.BindPFlag("static", ServeCmd.Flags().Lookup("static"))

	// bind env flag
	ServeCmd.Flags().StringP("env", "e", "dev", "The environment that we are running")
	viper.BindPFlag("env", ServeCmd.Flags().Lookup("env"))
}

func main() {
	RootCmd.AddCommand(ServeCmd)
	RootCmd.Execute()
}
