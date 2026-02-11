package db

import (
	"database/sql"
	"fmt"
	"time"

	"todo_final/pkg/config"

	_ "modernc.org/sqlite"
)

type Task struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Date    string `json:"date"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		fmt.Println("Error when tryng to exec AddTask")
		fmt.Println(err)
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		fmt.Println("Error LastInsertID in AddTask")
		return 0, err
	}
	return id, nil
}

func Tasks(limit int) ([]*Task, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("Incorrect number of tasks: %d", limit)
	}

	rows, err := db.Query("SELECT id, title, date, comment, repeat FROM scheduler ORDER BY date ASC LIMIT :limit;",
		sql.Named("limit", limit))
	if err != nil {
		return nil, err
	}

	var tasks []*Task

	for rows.Next() {
		var task Task

		if err := rows.Scan(&task.ID, &task.Title, &task.Date, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}

		tasks = append(tasks, &task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return []*Task{}, nil
	}

	return tasks, nil
}

func GetTask(id int) (*Task, error) {
	if id < 0 {
		return nil, fmt.Errorf("Incorrect ID")
	}

	var task Task
	err := db.QueryRow("SELECT id, title, date, comment, repeat FROM scheduler WHERE id = :id",
		sql.Named("id", id)).Scan(
		&task.ID,
		&task.Title,
		&task.Date,
		&task.Comment,
		&task.Repeat)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("Task not found")
		}
		return nil, fmt.Errorf("Error to get task: %w", err)
	}

	return &task, nil
}

func UpdateTask(task *Task) error {

	query := `UPDATE scheduler 
              SET date = :date, 
                  title = :title, 
                  comment = :comment, 
                  repeat = :repeat 
              WHERE id = :id`

	res, err := db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID),
	)

	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf(`Incorrect ID for updating task`)
	}
	return nil
}

func DeleteTask(id int) error {

	query := `DELETE FROM scheduler WHERE id = :id`
	_, err := db.Exec(query, sql.Named("id", id))
	if err != nil {
		return err
	}

	return nil
}

func UpdateDate(next string, id string) error {

	_, err := time.Parse(config.DateFormat, next)
	if err != nil {
		return err
	}

	query := `UPDATE scheduler 
              SET date = :date 
              WHERE id = :id`

	_, err = db.Exec(query,
		sql.Named("date", next),
		sql.Named("id", id))

	if err != nil {
		return err
	}

	return nil
}
