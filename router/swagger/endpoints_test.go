package docs

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/blueprint/blueprint/router"
)

func TestEndpoints(t *testing.T) {

	e := []router.Endpoint{{Type: "GET", Path: "/"}, {Type: "GET", Path: "/weather/:city"}}

	if err := ioutil.WriteFile("endpoints.tmpl", []byte(tmpl), 0777); err != nil {
		t.Error(err)
	}

	if err := Endpoints("endpoints.tmpl", "endpoints.html", e); err != nil {
		t.Error(err)
	}

	bytes, err := ioutil.ReadFile("endpoints.tmpl")
	if err != nil {
		t.Error(err)
	}

	if string(bytes) != html {
		t.Error("HTML rendered does not match HTML output expected")
	}

	if err := os.Remove("endpoints.tmpl"); err != nil {
		t.Error(err)
	}

	if err := os.Remove("endpoints.html"); err != nil {
		t.Error(err)
	}
}

var tmpl = `<html>
	<head>
		<title>API Docs</title>
		<link href='//fonts.googleapis.com/css?family=Lato:100' rel='stylesheet' type='text/css'>
		<style>
			body {
				margin: 0;
				padding: 0;
				width: 100%;
				color: #333;
				display: table;
				font-weight: 100;
				font-family: 'Lato';
			}
			.container {
				text-align: center;
        padding: 10px;
			}
		</style>
	</head>
	<body>
		<div class="container">
				<h1>API Endpoint Documentation</h1>
        <table align="center" style="margin: 0px auto;">
          <tr>
            <th>Request Type</th>
            <th>Endpoint Path</th>
          </tr>
          {{ range $endpoint := .Endpoints }}
          <tr>
            <td>{{ $endpoint.Type }}</td>
            <td>{{ $endpoint.Path }}</td>
          </tr>
        {{end}}
      </table>
		</div>
	</body>
</html>
`

var html = `<html>
	<head>
		<title>API Docs</title>
		<link href='//fonts.googleapis.com/css?family=Lato:100' rel='stylesheet' type='text/css'>
		<style>
			body {
				margin: 0;
				padding: 0;
				width: 100%;
				color: #333;
				display: table;
				font-weight: 100;
				font-family: 'Lato';
			}
			.container {
				text-align: center;
        padding: 10px;
			}
		</style>
	</head>
	<body>
		<div class="container">
				<h1>API Endpoint Documentation</h1>
        <table align="center" style="margin: 0px auto;">
          <tr>
            <th>Request Type</th>
            <th>Endpoint Path</th>
          </tr>
          {{ range $endpoint := .Endpoints }}
          <tr>
            <td>{{ $endpoint.Type }}</td>
            <td>{{ $endpoint.Path }}</td>
          </tr>
        {{end}}
      </table>
		</div>
	</body>
</html>
`
