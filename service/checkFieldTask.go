package service

import (
	"fmt"
	"time"

	"github.com/VenKaas/go_final_project/nextdate"
	"github.com/VenKaas/go_final_project/servicetask"
)

func checkFieldsTask(task *servicetask.Task) error {
	if task.Title == "" {
		return fmt.Errorf("не указан заголовок задачи")
	}

	if task.Date == "" {
		task.Date = time.Now().Format(nextdate.DFormat)
		return nil
	}
	_, err := time.Parse(nextdate.DFormat, task.Date)
	if err != nil {
		return fmt.Errorf("дата неверного формата")
	}

	newDate := time.Now().Format(nextdate.DFormat)
	err = nil
	if task.Repeat != "" {
		newDate, err = nextdate.NextDate(time.Now(), task.Date, task.Repeat)
	}

	if task.Date < time.Now().Format(nextdate.DFormat) {
		task.Date = newDate
	}

	return err
}
