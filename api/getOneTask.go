package api

import (
	"net/http"
)

func (srv Server) GetOneTask(rw http.ResponseWriter, rq *http.Request) {

	id := srv.Server.RequestId(rq)

	task, tr, err := srv.Server.SrvService.GetOneTask(id)
	checkErr(err)

	if tr.Err != "" {
		srv.Server.Response(tr, rw)
		return
	}

	srv.Server.Response(task, rw)

}
