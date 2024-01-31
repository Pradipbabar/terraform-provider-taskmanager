package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// Task represents a simple task structure.
type Task struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	IsDone bool   `json:"is_done"`
}

// TaskStore is an in-memory data store for tasks.
type TaskStore struct {
	mu    sync.RWMutex
	tasks map[int]Task
}

// NewTaskStore creates a new TaskStore instance.
func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks: make(map[int]Task),
	}
}

// CreateTask adds a new task to the store.
func (s *TaskStore) CreateTask(task Task) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	task.ID = len(s.tasks) + 1
	s.tasks[task.ID] = task
	return task.ID
}

// GetTask retrieves a task by ID from the store.
func (s *TaskStore) GetTask(id int) (Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.tasks[id]
	return task, ok
}

// UpdateTask updates an existing task in the store.
func (s *TaskStore) UpdateTask(id int, task Task) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.tasks[id]
	if !ok {
		return false
	}

	task.ID = id
	s.tasks[id] = task
	return true
}

// DeleteTask removes a task from the store by ID.
func (s *TaskStore) DeleteTask(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.tasks[id]
	if !ok {
		return false
	}

	delete(s.tasks, id)
	return true
}

// Handler function for handling task-related requests.
func handleTasks(w http.ResponseWriter, r *http.Request, store *TaskStore) {
	switch r.Method {
	case http.MethodPost:
		handleCreateTask(w, r, store)
	case http.MethodGet:
		handleGetTasks(w, r, store)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Handler function for creating a new task.
func handleCreateTask(w http.ResponseWriter, r *http.Request, store *TaskStore) {
	var newTask Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	id := store.CreateTask(newTask)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"id": %d}`, id)
}

// Handler function for retrieving tasks.
func handleGetTasks(w http.ResponseWriter, r *http.Request, store *TaskStore) {
	idStr := r.URL.Query().Get("id")
	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
			return
		}

		task, ok := store.GetTask(id)
		if !ok {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)
		return
	}

	// If no specific ID is provided, return all tasks.
	tasks := make([]Task, 0, len(store.tasks))
	store.mu.RLock()
	defer store.mu.RUnlock()
	for _, task := range store.tasks {
		tasks = append(tasks, task)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func main() {
	store := NewTaskStore()

	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		handleTasks(w, r, store)
	})

	port := 8080
	fmt.Printf("Server is running on :%d...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
