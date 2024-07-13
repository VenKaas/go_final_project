package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/VenKaas/go_final_project/servicetask"
)

type Service struct {
	SrvService servicetask.TaskStore
}

func NewService(SrvService servicetask.TaskStore) Service {
	return Service{SrvService: SrvService}
}

// проверяем валидность запросов
func (s Service) ReqValidate(t *servicetask.Task) (servicetask.TaskResp, error) {
	// проверяем что все поля date и title в task валидные
	var tr servicetask.TaskResp
	err := checkFieldsTask(t)
	if err != nil {
		tr.Err = "ошибка в формате поля date или title"
	}
	return tr, nil
}

func (s Service) Response(t any, w http.ResponseWriter) {
	resp, err := json.Marshal(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (s Service) RequestUpd(r *http.Request) (servicetask.Task, error) {
	var buf bytes.Buffer
	var task servicetask.Task
	// получаем данные из веб-интерфейса
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		return servicetask.Task{}, err
	}
	//переводим данные в стркутуру task
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		return servicetask.Task{}, err
	}
	return task, nil
}

func (s Service) RequestId(r *http.Request) int {
	id := r.FormValue("id")
	idInt, _ := strconv.Atoi(id)
	return idInt

}
