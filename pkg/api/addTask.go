package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"todo_final/pkg/config"
	db "todo_final/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	var task db.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeJson(w, err)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		writeJson(w, err)
		return
	}

	if task.Title == "" {
		writeJsonError(w, "Empty title", http.StatusBadRequest)
		return
	}

	_, err = NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeJsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = checkDate(&task); err != nil {
		writeJsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var resp = make(map[string]int64)
	resp["id"], err = db.AddTask(&task)

	if err != nil {
		writeJson(w, err)
		return
	}

	respJson, err := json.Marshal(resp)
	if err != nil {
		writeJson(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(respJson)

}

func checkDate(task *db.Task) error {
	now := time.Now()
	today := now.Format(config.DateFormat)

	if task.Date == today {
		return nil
	}

	if task.Date == "" {
		task.Date = today
		return nil
	}

	if task.Repeat == "d 1" {
		task.Date = today
		return nil
	}

	t, err := time.Parse(config.DateFormat, task.Date)
	if err != nil {
		return err
	}

	if now.After(t) && task.Repeat != "" {
		newDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
		task.Date = newDate
		return nil
	}

	if now.After(t) && task.Repeat == "" {
		newDate := today
		task.Date = newDate
		return nil
	}

	return nil
}

func writeJson(w http.ResponseWriter, data any) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	jsonData, err := json.Marshal(data)

	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)

		errorResponse := map[string]string{"error": "Failed to marshal JSON"}

		json.NewEncoder(w).Encode(errorResponse)

		return

	}

	w.Write(jsonData)

}

func writeJsonError(w http.ResponseWriter, msg string, statusCode int) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(map[string]string{

		"error": msg,
	})

}
