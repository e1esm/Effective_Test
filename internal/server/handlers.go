package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/e1esm/Effective_Test/internal/models/users"
	"github.com/e1esm/Effective_Test/internal/repository/postgres"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
)

var (
	methodErr      = errors.New("invalid method for the URL: %v")
	invalidReq     = errors.New("invalid request was sent to the URL: %v")
	marshallingErr = errors.New("error while operating over the request input data")
	identityErr    = errors.New("error occurred while identifying the user: %v")
	saveErr        = errors.New("error while inserting user: %v")
	deleteErr      = errors.New("error while deleting user with this ID: %v")
)

func (hs *HttpServer) New(r http.ResponseWriter, h *http.Request) {
	if h.Method != http.MethodPost {
		r.WriteHeader(http.StatusMethodNotAllowed)
		if _, err := r.Write([]byte(fmt.Sprintf(methodErr.Error(), h.RequestURI))); err != nil {
			log.Println(err.Error())
		}
		return
	}
	content, err := io.ReadAll(h.Body)
	if err != nil {
		r.WriteHeader(http.StatusBadRequest)
		if _, err := r.Write([]byte(fmt.Sprintf(invalidReq.Error(), h.RequestURI))); err != nil {
			log.Println(err.Error())
		}
		return
	}

	request := &users.User{}
	if err := json.Unmarshal(content, request); err != nil {
		r.WriteHeader(http.StatusBadRequest)
		if _, err := r.Write([]byte(fmt.Sprintf(marshallingErr.Error()))); err != nil {
			log.Println(err.Error())
		}
		return
	}
	user := hs.identityService.Identify(*request)
	if user == nil {
		r.WriteHeader(http.StatusInternalServerError)
		if _, err := r.Write([]byte(fmt.Sprintf(identityErr.Error(), request.Name))); err != nil {
			log.Println(err.Error())
		}
		return
	}

	id, err := hs.userService.Save(context.Background(), user)
	if err != nil {
		r.WriteHeader(http.StatusInternalServerError)
		if _, err := r.Write([]byte(fmt.Sprintf(saveErr.Error(), user.User))); err != nil {
			log.Println(err.Error())
		}
		return
	}

	r.WriteHeader(http.StatusOK)
	r.Write([]byte(id.String()))
}

func (hs *HttpServer) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf(methodErr.Error(), r.RequestURI)))
		return
	}

	type DeleteRequest struct {
		ID string `json:"id"`
	}

	var toBeDeleted DeleteRequest
	content, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(invalidReq.Error(), r.RequestURI)))
		return
	}
	if err := json.Unmarshal(content, &toBeDeleted); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		w.Write([]byte(marshallingErr.Error()))
		return
	}

	if _, err := hs.userService.Delete(context.Background(), uuid.MustParse(toBeDeleted.ID)); err != nil {
		switch err {
		case postgres.NoRecordsFound:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(deleteErr.Error(), toBeDeleted)))
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(toBeDeleted.ID))
}
