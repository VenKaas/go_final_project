package env

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// SetFlagParams - настраивает параметры запуска приложения
func SetFlagParams() {
	pass := flag.String("password", "12345", "Пароль для приложения")
	port := flag.String("port", "7540", "Порт для запуска веб сервера")
	dbPath := flag.String("dbpath", "", "Путь к базе данных")

	flag.Parse()
	os.Setenv("TODO_PASSWORD", *pass)
	os.Setenv("TODO_PORT", *port)
	os.Setenv("TODO_DBFILE", *dbPath)
}

func DbName() string {
	dbFile := "scheduler.db"
	envFile := os.Getenv("TODO_DBFILE")
	if len(envFile) > 0 {
		dbFile = filepath.Join(envFile, "scheduler.db")
	}
	log.Println("путь к БД:", dbFile)
	return dbFile
}

var rightPassword string

var readPassOnce = sync.OnceFunc(func() {
	rightPassword = os.Getenv("TODO_PASSWORD")
})

func SetPass() string {
	readPassOnce()
	return rightPassword
}

func SetPort() string {
	return os.Getenv("TODO_PORT")
}
