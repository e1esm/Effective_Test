package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/e1esm/Effective_Test/internal/models/users"
	"io"
	"log"
	"net/http"
)

var (
	methodErr      = errors.New("invalid method for the URL: %v")
	invalidReq     = errors.New("invalid request was sent to the URL: %v")
	marshallingErr = errors.New("error while operating over the request input data")
	identityErr    = errors.New("error occurred while identifying the user: %v")
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

	r.WriteHeader(http.StatusOK)
}
