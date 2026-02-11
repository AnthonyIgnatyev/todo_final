package api

import (
	"net/http"
	"strconv"
	"time"

	db "todo_final/pkg/db"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {

		writeJsonError(w, "Required POST method", http.StatusMethodNotAllowed)
		return

	}

	if !r.URL.Query().Has("id") {
		writeJsonError(w, "ID is required", http.StatusBadRequest)
		return

	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeJsonError(w, "Invalid ID format", http.StatusBadRequest)
		return

	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if task.Repeat == "" {
		err := db.DeleteTask(id)
		if err != nil {
			writeJsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
		return

	}

	now := time.Now()

	newDate, err := NextDate(now, task.Date, task.Repeat)

	if err != nil {
		writeJsonError(w, err.Error(), http.StatusInternalServerError)
		return

	}

	err = db.UpdateDate(newDate, idStr)
	if err != nil {
		writeJsonError(w, err.Error(), http.StatusInternalServerError)
		return

	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))

}
