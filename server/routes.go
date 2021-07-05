package server

import (
	"fmt"
	"net/http"

	"github.com/EricNeid/go-webserver/database"
	"github.com/EricNeid/go-webserver/model"
)

func (srv ApplicationServer) user(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		user, err := model.NewUser(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		id, err := database.AddUser(srv.Logger, srv.db, user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Could not add user to datbase: %v", err)))
			return
		}
		res := model.ResponseUserId{UserId: id}
		w.WriteHeader(http.StatusOK)
		w.Write(res.ToJson())
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
