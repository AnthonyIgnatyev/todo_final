package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	db "todo_final/pkg/db"
)

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var task db.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeJsonError(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		writeJsonError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		writeJsonError(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	if _, err := strconv.Atoi(task.ID); err != nil {
		writeJsonError(w, "Task ID must be numeric", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeJsonError(w, "Task title is required", http.StatusBadRequest)
		return
	}

	if err = checkDate(&task); err != nil {
		writeJsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		if err.Error() == "task not found" {
			writeJsonError(w, "Task not found", http.StatusNotFound)
		} else {
			writeJsonError(w, "Failed to update task", http.StatusInternalServerError)
		}
		return

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("{}"))
	if err != nil {
		writeJsonError(w, "Failed to write to responseWriter", http.StatusInternalServerError)
		return
	}
}
