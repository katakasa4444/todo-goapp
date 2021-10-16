package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

func main() {
	http.Handle("/", new(todoList).Handler("/"))
	http.Handle("/api", NewAPIHandler())

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

type todoList struct {
	app.Compo

	list []todo
}

func (h *todoList) OnMount(ctx app.Context) {
	ctx.Async(func() {
		h.updateTodo()
	})
}

func (h *todoList) updateTodo() {
	resp, err := api.Get("http://localhost:8080/api")
	if err != nil {
		log.Printf("%v", err)
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&h.list); err != nil {
		log.Printf("%v", err)
	}
	h.Update()
	log.Println("update todo")

}

func (h *todoList) Render() app.UI {
	log.Println("this is render")
	return app.Div().Body(
		app.H1().Body(
			app.Text("Todo List"),
		),
		app.Input().
			Value("").
			OnChange(h.OnInputChange),
		app.P().Body(
			app.Range(h.list).Slice(func(i int) app.UI {
				v := h.list[i]

				return app.Div().Body(
					app.Input().
						ID(fmt.Sprintf("%d", v.ID)).
						Type("checkbox").
						Checked(v.Done).
						OnChange(h.OnDoneChange),
					app.Div().Text(v.Text),
					app.Input().
						ID(fmt.Sprintf("%d", v.ID)).
						Type("button").
						Value("x").
						OnClick(h.OnClieckDelete),
				)
			})),
	)
}

func (h *todoList) Handler(path string) http.Handler {
	app.Route(path, h)
	app.RunWhenOnBrowser()
	return &app.Handler{
		Description: "An Hello World! example",
		Name:        "Hello",
	}
}

func (h *todoList) OnInputChange(ctx app.Context, e app.Event) {
	log.Println("input change event")
	text := ctx.JSSrc().Get("value").String()

	b, err := json.Marshal(todo{Text: text})
	if err != nil {
		log.Println(err)
		return
	}

	resp, err := api.Post("/api", "application/json", bytes.NewReader(b))
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	h.updateTodo()
	ctx.JSSrc().Set("value", "")
}

func (h *todoList) OnDoneChange(ctx app.Context, e app.Event) {
	s := ctx.JSSrc().Get("id").String()
	id, _ := strconv.Atoi(s)

	done := ctx.JSSrc().Get("checked").Bool()

	b, err := json.Marshal(todo{ID: id, Done: done})
	if err != nil {
		log.Println(err)
		return
	}

	resp, err := api.Post("/api", "application/json", bytes.NewReader(b))
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
}

func (h *todoList) OnClieckDelete(ctx app.Context, e app.Event) {
	s := ctx.JSSrc().Get("id").String()
	id, _ := strconv.Atoi(s)
	b, err := json.Marshal(todo{ID: id})
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("delete todo %s\n", s)
	req, err := http.NewRequest(http.MethodDelete, "/api", bytes.NewReader(b))
	if err != nil {
		log.Println(err)
		return
	}
	resp, err := api.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	h.updateTodo()
}

var api = &http.Client{}
