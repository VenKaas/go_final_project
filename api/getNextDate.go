package api

import (
	"log"
	"net/http"
	"time"

	"github.com/VenKaas/go_final_project/dformat"
	"github.com/VenKaas/go_final_project/nextdate"
)

func GetNextDate(rw http.ResponseWriter, rq *http.Request) {
	// парсим форму
	nowInString := rq.FormValue("now")
	now, err := time.Parse(dformat.DFormat, nowInString)
	if err != nil {
		log.Println("ошибка парсинга формата заданной даты:", err)
	}
	date := rq.FormValue("date")
	repeat := rq.FormValue("repeat")

	// получаем новую дату
	st, err := nextdate.NextDate(now, date, repeat)
	if err != nil {
		log.Println("ошибка при переносе даты:", err)
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(st))
}
