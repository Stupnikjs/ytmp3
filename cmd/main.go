package main

import (
	"fmt"
	"net/http"

	"github.com/Stupnikjs/skeleton/api"
	_ "google.golang.org/api/option"
)

func main() {

	app := api.Application{
		Port: 8080,
	}
	http.ListenAndServe(fmt.Sprintf(":%d", app.Port), app.Routes())

}
