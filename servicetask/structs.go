package servicetask

import "database/sql"

type TaskStore struct {
	Db *sql.DB
}

type Task struct {
	Id      string `json:"id,omitempty"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title,omitempty"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

type TaskResp struct {
	Id  string `json:"id,omitempty"`
	Err string `json:"error,omitempty"`
}
