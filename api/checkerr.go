package api

import "log"

func checkErr(err error) {
	if err != nil {
		log.Println(err)

	}
}
