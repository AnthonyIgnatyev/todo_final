package api

import (
	"net/http"
	"time"

	"todo_final/pkg/scheduler"
)

func init() {
	http.HandleFunc("/api/nextdate", func(w http.ResponseWriter, r *http.Request) {
		nowStr := r.URL.Query().Get("now")
		dateStr := r.URL.Query().Get("date")
		repeatStr := r.URL.Query().Get("repeat")

		var now time.Time
		if nowStr == "" {
			now = time.Now()
		} else {
			var err error
			now, err = time.Parse("20060102", nowStr)
			if err != nil {
				http.Error(w, "invalid date format", http.StatusBadRequest)
				return
			}
		}

		result, err := scheduler.NextDate(now, dateStr, repeatStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Write([]byte(result))
	})
}
