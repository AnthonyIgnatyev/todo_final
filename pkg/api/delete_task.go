package api

import (
	"net/http"
	"strconv"

	db "todo_final/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("id") {
		writeJsonError(w, "The ID parameter can't be empty", http.StatusBadRequest)
		return
	}

	id := r.URL.Query().Get("id")
	strID, err := strconv.Atoi(id)
	if err != nil {
		writeJsonError(w, "Wrong ID format: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = db.DeleteTask(strID)
	if err != nil {
		writeJsonError(w, "Can't delete task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}
