package main

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"
)

type todo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

func httpError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

type Repo interface {
	getTodoList() ([]todo, error)
	writeTodoList(list *[]todo) error
}

type APIHandler struct {
	repo Repo
}

func NewAPIHandler() *APIHandler {
	return &APIHandler{fileJSON}
}

func (api *APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		api.handleListTodo(w, r)
	case http.MethodPost:
		api.handleUpdateTodo(w, r)
	default:
		httpError(w, http.StatusMethodNotAllowed)
	}
}

func (api *APIHandler) handleListTodo(w http.ResponseWriter, _ *http.Request) {
	todoList, err := api.repo.getTodoList()
	if err != nil {
		httpError(w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&todoList); err != nil {
		httpError(w, http.StatusInternalServerError)
		return
	}
}

func (api *APIHandler) handleUpdateTodo(w http.ResponseWriter, r *http.Request) {
	var u todo
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		httpError(w, http.StatusInternalServerError)
		return
	}

	todoList, err := api.repo.getTodoList()
	if err != nil {
		httpError(w, http.StatusInternalServerError)
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

	if err := api.repo.writeTodoList(&todoList); err != nil {
		httpError(w, http.StatusInternalServerError)
		return
	}
}

type FileJSON struct {
	mu sync.Mutex
}

var fileJSON = new(FileJSON)

func (j *FileJSON) getTodoList() ([]todo, error) {
	j.mu.Lock()
	defer j.mu.Unlock()

	var v []todo
	f, err := os.OpenFile("todo.json", os.O_RDONLY|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if err = json.NewDecoder(f).Decode(&v); err != nil {
		return nil, err
	}
	return v, nil
}

func (j *FileJSON) writeTodoList(list *[]todo) error {
	j.mu.Lock()
	defer j.mu.Unlock()

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
