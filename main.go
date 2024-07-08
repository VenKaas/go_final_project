package main

import (
	"log"

	"github.com/VenKaas/go_final_project/api"
	"github.com/VenKaas/go_final_project/db"
	"github.com/VenKaas/go_final_project/env"
)

func main() {
	env.SetFlagParams()

	err := db.DbExistance()
	if err != nil {
		log.Println("Ошибка при подключении к базе:", err)
		return
	}
	api.StartWebServer()
}
