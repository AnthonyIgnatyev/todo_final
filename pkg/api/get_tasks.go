package api

import (
	"net/http"
	"strconv"

	db "todo_final/pkg/db"

	_ "modernc.org/sqlite"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50)
	if err != nil {
		writeJson(w, err)
		return
	}
	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if !r.URL.Query().Has("id") {
		writeJsonError(w, "The ID parameter can't be empty", http.StatusBadRequest)
		return
	}

	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJsonError(w, "Wrong ID format: "+err.Error(), http.StatusBadRequest)
		return
	}

	if id < 0 {
		writeJsonError(w, "Wrong ID", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJsonError(w, "Can't get task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJson(w, task)
}
