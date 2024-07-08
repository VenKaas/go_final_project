package api

import (
	"net/http"
)

func (srv Server) UpdTask(rw http.ResponseWriter, rq *http.Request) {

	task, err := srv.Server.RequestUpd(rq)
	checkErr(err)

	tr, err := srv.Server.ReqValidate(&task)
	checkErr(err)

	if tr.Err != "" {
		srv.Server.Response(tr, rw)
		return
	}

	tr, err = srv.Server.SrvService.Update(task)
	checkErr(err)

	if tr.Err != "" {
		srv.Server.Response(tr, rw)
		return
	}
	srv.Server.Response(task, rw)
}
