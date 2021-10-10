package main

import (
	"encoding/json"
	"net/http"
	"os"
)

type todo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

func httpError(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func handleListTodo(w http.ResponseWriter, _ *http.Request) {
	var todoList []todo
	if err := getTodoList(&todoList); err != nil {
		httpError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&todoList); err != nil {
		httpError(w)
		return
	}
}

func handleUpdateTodo(w http.ResponseWriter, r *http.Request) {
	var u todo
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		httpError(w)
		return
	}

	var todoList []todo
	if err := getTodoList(&todoList); err != nil {
		httpError(w)
		return
	}

	found := false
	for i := range todoList {
		if todoList[i].ID == u.ID {
			todoList[i].Text = u.Text
			todoList[i].Done = u.Done
			found = true
			break
		}
	}

	if !found {
		todoList = append(todoList, todo{
			ID:   len(todoList) + 1,
			Text: u.Text,
			Done: u.Done,
		})
	}

	if err := writeTodoList(&todoList); err != nil {
		httpError(w)
		return
	}
}

func getTodoList(v *[]todo) error {
	f, err := os.OpenFile("todo.json", os.O_RDONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	if err = json.NewDecoder(f).Decode(v); err != nil {
		return err
	}
	return nil
}

func writeTodoList(list *[]todo) error {
	f, err := os.OpenFile("todo.json", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(list); err != nil {
		return err
	}
	return nil
}
