package api

import (
	"net/http"
	"strconv"

	db "todo_final/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if !r.URL.Query().Has("id") {
		writeJsonError(w, "The ID parametr can't be empty", http.StatusBadRequest)
	}

	id := r.URL.Query().Get("id")
	strID, err := strconv.Atoi(id)
	if err != nil {
		writeJsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = db.DeleteTask(strID)
	if err != nil {
		writeJsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}
