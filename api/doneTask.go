package api

import (
	"net/http"
)

func (srv Server) DoneTask(rw http.ResponseWriter, rq *http.Request) {

	id := srv.Server.RequestId(rq)

	task, err := srv.Server.SrvService.Done(id)
	checkErr(err)

	srv.Server.Response(task, rw)

}
