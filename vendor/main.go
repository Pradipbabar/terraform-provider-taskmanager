package main

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

// Task struct represents a task in the todo list
type Task struct {
	ID     int
	Name   string
	Status bool
}

var tasksMap map[int]Task
var taskIDCounter int

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var taskList []Task
	for _, task := range tasksMap {
		taskList = append(taskList, task)
	}

	err = tmpl.Execute(w, taskList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		taskName := r.FormValue("task")
		taskID := taskIDCounter
		task := Task{ID: taskID, Name: taskName, Status: false}
		tasksMap[taskID] = task
		taskIDCounter++
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		taskIDStr := strings.TrimPrefix(r.URL.Path, "/delete/")
		taskID, err := strconv.Atoi(taskIDStr)
		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		if _, exists := tasksMap[taskID]; exists {
			delete(tasksMap, taskID)
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		taskIDStr := strings.TrimPrefix(r.URL.Path, "/update/")
		taskID, err := strconv.Atoi(taskIDStr)
		if err != nil {
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		taskName := r.FormValue("task")
		if task, exists := tasksMap[taskID]; exists {
			task.Name = taskName
			tasksMap[taskID] = task
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func main() {
	tasksMap = make(map[int]Task)
	taskIDCounter = 1

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/delete/", deleteHandler)
	http.HandleFunc("/update/", updateHandler)

	http.ListenAndServe(":8080", nil)
}
