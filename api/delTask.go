package api

import (
	"net/http"
)

func (srv Server) DelTask(rw http.ResponseWriter, rq *http.Request) {

	id := srv.Server.RequestId(rq)

	tr, err := srv.Server.SrvService.Delete(id)
	checkErr(err)

	srv.Server.Response(tr, rw)

}
