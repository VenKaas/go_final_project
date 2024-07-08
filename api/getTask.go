package api

import (
	"net/http"
	"time"

	"github.com/VenKaas/go_final_project/servicetask"
)

func (srv Server) GetTask(rw http.ResponseWriter, rq *http.Request) {
	var tasks = map[string][]servicetask.Task{}
	var err error
	var tr servicetask.TaskResp

	searchString := rq.FormValue("search")

	switch searchString {
	//если нажат поиск по дате, то выбираем записи согласно дате
	case "":
		tasks, tr, err = srv.Server.SrvService.GetAll()
		checkErr(err)

	//если нажат поиск, то выбираем записи исходя из строки поиска
	default:
		searchDate, errParse := time.Parse("02.01.2006", searchString)
		//если в поиске дата
		if errParse == nil {
			tasks, err = srv.Server.SrvService.GetSearchDate(searchDate)
			checkErr(err)

			//если в поиске строка
		} else {
			tasks, err = srv.Server.SrvService.GetSearch(searchString)
			checkErr(err)
		}
	}

	if tr.Err != "" {
		srv.Server.Response(tr, rw)
		return
	}

	srv.Server.Response(tasks, rw)

}
