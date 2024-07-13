package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/VenKaas/go_final_project/env"
	"github.com/VenKaas/go_final_project/service"
	"github.com/VenKaas/go_final_project/servicetask"
	"github.com/go-chi/chi"
)

type Server struct {
	Server service.Service
}

func NewServer(server service.Service) Server {
	return Server{Server: server}
}

func StartWebServer() {

	// подключаемся к БД
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		log.Println("ошибка при подключении к БД:", err)
	}
	defer db.Close()

	store := servicetask.NewTaskStore(db)
	service := service.NewService(store)
	server := NewServer(service)

	todoPort := env.SetPort()
	webDir := "web"

	r := chi.NewRouter()
	r.Get("/api/nextdate", GetNextDate)
	r.Handle("/*", http.FileServer(http.Dir(webDir)))
	r.Post("/api/task", Auth(server.AddTask))
	r.Post("/api/task/done", Auth(server.DoneTask))
	r.Get("/api/task", Auth(server.GetOneTask))
	r.Delete("/api/task", Auth(server.DelTask))
	r.Put("/api/task", Auth(server.UpdTask))
	r.Get("/api/tasks", Auth(server.GetTask))
	r.Post("/api/signin", server.CheckPass)

	log.Println("Запускаем веб сервер")
	err = http.ListenAndServe(":"+todoPort, r)
	if err != nil {
		log.Println("ошибка при запуске веб сервера:", err)
	}
}
