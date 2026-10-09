package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"todoapp/todo"

	"github.com/gorilla/mux"
)

type HTTPHandlers struct {
	todoList *todo.List
}

func NewHTTPHandlers(todoList *todo.List) *HTTPHandlers {
	return &HTTPHandlers{
		todoList: todoList,
	}
}

/*
pattern/endpoint: /tasks
method: POST
info: JSON in HTTP request body

succeed:
  - status code: 201 Created
  - response body: JSON represent created task

failed:
  - status code: 400, 409, 500 ...
  - response body: JSON with error + time
*/
func (h *HTTPHandlers) HandlerCreateTask(w http.ResponseWriter, r *http.Request) {
	var taskDTO TaskDTO
	if err := json.NewDecoder(r.Body).Decode(&taskDTO); err != nil {
		errDTO := NewErrorDTO(err)

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	if err := taskDTO.ValidateForCreate(); err != nil {
		errDTO := NewErrorDTO(err)

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)

		return
	}

	todoTask := todo.NewTask(taskDTO.Title, taskDTO.Description)
	if err := h.todoList.AddTask(todoTask); err != nil {
		errDTO := NewErrorDTO(err)

		if errors.Is(err, todo.ErrTaskAlreadyExists) {
			http.Error(w, errDTO.ToString(), http.StatusConflict)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}

		return
	}

	b, err := json.MarshalIndent(todoTask, "", "	")

	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response:", err)
		return
	}
}

/*
pattern: /tasks/{title}
method: GET
info: pattern

succeed:
  - status code: 200 OK
  - response body: JSON represent created task

failed:
  - status code: 400, 404, 409, 500, ...
  - response body: JSON with error + time
*/
func (h *HTTPHandlers) HandlerGetTask(w http.ResponseWriter, r *http.Request) {
	title := mux.Vars(r)["title"]

	task, err := h.todoList.GetTask(title)
	if err != nil {
		errDTO := NewErrorDTO(err)

		if errors.Is(err, todo.ErrTaskNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}

		return
	}

	b, err := json.MarshalIndent(task, "", "	")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response:", err)
		return
	}
}

/*
pattern: /tasks
method: GET
info: -

succeed:
  - status code: 200 OK
  - response body: JSON represent found tasks

failed:
  - status code: 400, 404, 409, 500, ...
  - response body: JSON with error + time
*/
func (h *HTTPHandlers) HandlerGetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.todoList.ListTasks()

	b, err := json.MarshalIndent(tasks, "", "    ")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response:", err)
		return
	}
}

/*
pattern: /tasks?completed=false
method: GET
info: -

succeed:
  - status code: 200 OK
  - response body: JSON represent found tasks

failed:
  - status code: 400, 404, 409, 500, ...
  - response body: JSON with error + time
*/
func (h *HTTPHandlers) HandlerGetAllUncompletedTasks(w http.ResponseWriter, r *http.Request) {
	notCompletedTasks := h.todoList.ListNotCompletedTask()

	b, err := json.MarshalIndent(notCompletedTasks, "", "    ")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response:", err)
		return
	}
}

/*
pattern: /tasks/{title}
method: PATCH
info: pattern + JSON in request body

succeed:
  - status code: 200 OK
  - response body: JSON represents changed task

failed:
  - status code: 400, 404, 409, 500, ...
  - response body: JSON with error + time
*/
func (h *HTTPHandlers) HandlerCompleteTask(w http.ResponseWriter, r *http.Request) {
	var completeDTO CompleteTaskDTO
	if err := json.NewDecoder(r.Body).Decode(&completeDTO); err != nil {
		errDTO := NewErrorDTO(err)

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	title := mux.Vars(r)["title"]

	var (
		changedTask todo.Task
		err         error
	)

	if completeDTO.Complete {
		changedTask, err = h.todoList.CompleteTask(title)
	} else {
		changedTask, err = h.todoList.UncompleteTask(title)
	}

	if err != nil {
		errDTO := NewErrorDTO(err)

		if errors.Is(err, todo.ErrTaskNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}

		return
	}

	b, err := json.MarshalIndent(changedTask, "", "    ")

	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response:", err)
		return
	}
}

/*
pattern: /tasks/{title}
method: DELETE
info: pattern

succeed:
  - status code: 204 No Content
  - response body: -

failed:
  - status code: 400, 404, 500, ...
  - response body: JSON with error + time
*/
func (h *HTTPHandlers) HandlerDeleteTask(w http.ResponseWriter, r *http.Request) {
	title := mux.Vars(r)["title"]

	if err := h.todoList.DeleteTask(title); err != nil {
		errDTO := NewErrorDTO(err)
		if errors.Is(err, todo.ErrTaskNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
