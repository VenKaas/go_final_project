package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"

	"github.com/VenKaas/go_final_project/env"
)

// Existance - TODO: описание
func Existance() error {

	dbFile := env.DbName()

	_, err := os.Stat(dbFile)
	if err != nil {
		log.Println("Создаем новую базу данных с таблицей scheduler")
		err = dbCreate(dbFile)
		if err != nil {
			return fmt.Errorf("ошибка создания новой базы: %w", err)
		}
	}
	log.Println("Успешное подключение к БД")
	return nil
}

func dbCreate(dbFile string) error {
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("ошибка при подключении к БД: %w", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE scheduler (id INTEGER PRIMARY KEY AUTOINCREMENT,	
											date CHAR(8), 
											title VARCHAR(128) NOT NULL DEFAULT "", 
											comment VARCHAR(256) ,
											repeat VARCHAR(128))`)
	if err != nil {
		return fmt.Errorf("ошибка при создании таблицы в БД: %w", err)
	}

	_, err = db.Exec("CREATE INDEX dateindex ON scheduler (date)")
	if err != nil {
		return fmt.Errorf("ошибка при создании индекса dateindex в БД: %w", err)
	}
	return nil
}
