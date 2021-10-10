package main

import (
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

func main() {
	http.Handle("/", new(hello).Handler("/"))
	http.Handle("/api", NewAPIHandler())

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

type hello struct {
	app.Compo
}

func (h *hello) Render() app.UI {
	return app.H1().Text("Hello World!")
}

func (h *hello) Handler(path string) http.Handler {
	app.Route(path, h)
	app.RunWhenOnBrowser()
	return &app.Handler{
		Description: "An Hello World! example",
		Name:        "Hello",
	}
}
