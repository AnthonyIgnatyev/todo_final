package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"todo_final/pkg/config"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {

	if date == "" {
		date = time.Now().Format(config.DateFormat)
	}

	dateTime, err := time.Parse(config.DateFormat, date)
	if err != nil {
		return "", fmt.Errorf("Invalid date format: %v", err)
	}

	rule := strings.Split(repeat, " ")

	switch rule[0] {
	case "d":
		if len(rule) < 2 {
			return "", fmt.Errorf("Missing interval for 'd' rule")
		}

		interval, err := strconv.Atoi(rule[1])
		if err != nil {
			return "", fmt.Errorf("Invalid interval for 'd' rule: %v", err)
		}

		if interval > 400 {
			return "", fmt.Errorf("Interval exceeds maximum allowed value (400)")
		}

		if interval == 1 {
			dateTime = dateTime.AddDate(0, 0, 1)
		} else {
			for {
				dateTime = dateTime.AddDate(0, 0, interval)
				if dateTime.After(now) {
					break
				}
			}
		}

	case "y":
		for {
			dateTime = dateTime.AddDate(1, 0, 0)
			if dateTime.After(now) {
				break
			}
		}
	case "":
		dateTime = time.Now()

	default:
		return "", fmt.Errorf("Unsupported repeat rule: %s", rule[0])
	}

	return dateTime.Format(config.DateFormat), nil
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Query().Get("repeat") == "" {
		fmt.Println("Repeat rule is empty")
		return
	}

	var now time.Time

	if r.URL.Query().Has("now") {
		if r.URL.Query().Get("now") == "today" {
			now = time.Now()
		} else {
			parsedNow, err := time.Parse(config.DateFormat, r.URL.Query().Get("now"))
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid 'now' date format. Expected:[%s]", config.DateFormat), http.StatusBadRequest)
				return
			}
			now = parsedNow
		}
	} else {
		now = time.Now()
	}

	result, err := NextDate(now, r.FormValue("date"), r.FormValue("repeat"))

	if err != nil {
		fmt.Println(err)
		return
	}

	w.Write([]byte(result))
}
