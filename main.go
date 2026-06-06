package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Done      bool   `json:"done"`
	CreatedAt string `json:"created_at"`
}

type CreateTaskRequest struct {
	Title string `json:"title"`
}

var (
	tasks     []Task
	tasksMu   sync.Mutex
	nextID    = 1
	startTime time.Time
	logger    *slog.Logger
)

func init() {
	startTime = time.Now()
	logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	tasks = []Task{
		{ID: 1, Title: "Изучить Kubernetes", Done: false, CreatedAt: time.Now().Format(time.RFC3339)},
		{ID: 2, Title: "Написать манифесты", Done: true, CreatedAt: time.Now().Format(time.RFC3339)},
		{ID: 3, Title: "Задеплоить приложение", Done: false, CreatedAt: time.Now().Format(time.RFC3339)},
	}
	nextID = 4
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}

func main() {
	port := getEnv("PORT", "5000")
	podName := getEnv("HOSTNAME", "unknown")

	logger.Info("Starting task backend", "port", port, "pod", podName)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", corsMiddleware(healthHandler))
	mux.HandleFunc("/api/tasks", corsMiddleware(tasksHandler))
	mux.HandleFunc("/api/tasks/", corsMiddleware(taskByIDHandler))

	logger.Info("Server ready", "address", fmt.Sprintf(":%s", port))
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		logger.Error("Server failed", "error", err)
		os.Exit(1)
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error("JSON encode error", "error", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]any{
		"status":  "healthy",
		"service": "task-backend",
		"uptime":  time.Since(startTime).Seconds(),
		"pod":     getEnv("HOSTNAME", "unknown"),
	})
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tasksMu.Lock()
		result := make([]Task, len(tasks))
		copy(result, tasks)
		tasksMu.Unlock()
		logger.Info("GET /api/tasks", "count", len(result))
		respondJSON(w, http.StatusOK, result)

	case http.MethodPost:
		var req CreateTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Title) == "" {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "title is required"})
			return
		}
		tasksMu.Lock()
		t := Task{ID: nextID, Title: strings.TrimSpace(req.Title), Done: false, CreatedAt: time.Now().Format(time.RFC3339)}
		nextID++
		tasks = append(tasks, t)
		tasksMu.Unlock()
		logger.Info("POST /api/tasks", "id", t.ID, "title", t.Title)
		respondJSON(w, http.StatusCreated, t)

	default:
		respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

func taskByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	switch r.Method {
	case http.MethodPut:
		tasksMu.Lock()
		defer tasksMu.Unlock()
		for i, t := range tasks {
			if t.ID == id {
				tasks[i].Done = !tasks[i].Done
				logger.Info("PUT /api/tasks/:id", "id", id, "done", tasks[i].Done)
				respondJSON(w, http.StatusOK, tasks[i])
				return
			}
		}
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})

	case http.MethodDelete:
		tasksMu.Lock()
		defer tasksMu.Unlock()
		for i, t := range tasks {
			if t.ID == id {
				tasks = append(tasks[:i], tasks[i+1:]...)
				logger.Info("DELETE /api/tasks/:id", "id", id)
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})

	default:
		respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}
