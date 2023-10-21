package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/e1esm/Effective_Test/internal/models/users"
	"github.com/e1esm/Effective_Test/internal/repository/postgres"
	"github.com/e1esm/Effective_Test/pkg/utils/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
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

func (hs *HttpServer) New(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		if _, err := w.Write([]byte(fmt.Sprintf(methodErr.Error(), r.RequestURI))); err != nil {
			logger.GetLogger().Error(err.Error(), zap.String("URL", r.RequestURI))
		}
		return
	}
	content, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte(fmt.Sprintf(invalidReq.Error(), r.RequestURI))); err != nil {
			logger.GetLogger().Error(err.Error(), zap.String("URL", r.RequestURI))
		}
		return
	}

	request := &users.User{}
	if err := json.Unmarshal(content, request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte(fmt.Sprintf(marshallingErr.Error()))); err != nil {
			logger.GetLogger().Error(err.Error(), zap.String("URL", r.RequestURI))
		}
		return
	}
	user := hs.identityService.Identify(*request)
	if user == nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(fmt.Sprintf(identityErr.Error(), request.Name))); err != nil {
			logger.GetLogger().Error(err.Error(), zap.String("URL", r.RequestURI))
		}
		return
	}

	id, err := hs.userService.Save(context.Background(), user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(fmt.Sprintf(saveErr.Error(), user.User))); err != nil {
			logger.GetLogger().Error(err.Error(), zap.String("URL", r.RequestURI))
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(id.String())); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.GetLogger().Error(err.Error(), zap.String("URL", r.RequestURI))
	}
}

func (hs *HttpServer) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		if _, err := w.Write([]byte(fmt.Sprintf(methodErr.Error(), r.RequestURI))); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.GetLogger().Error(err.Error(), zap.String("URL", r.RequestURI))
		}
		return
	}

	type DeleteRequest struct {
		ID string `json:"id"`
	}

	var toBeDeleted DeleteRequest
	content, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte(fmt.Sprintf(invalidReq.Error(), r.RequestURI))); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.GetLogger().Error(err.Error(), zap.String("URL", r.RequestURI))
		}
		return
	}
	if err := json.Unmarshal(content, &toBeDeleted); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		if _, err := w.Write([]byte(marshallingErr.Error())); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.GetLogger().Error(err.Error(), zap.String("URL", r.RequestURI))
		}
		return
	}

	if _, err := hs.userService.Delete(context.Background(), uuid.MustParse(toBeDeleted.ID)); err != nil {
		switch err {
		case postgres.NoRecordsFound:
			w.WriteHeader(http.StatusNotFound)
			if _, err := w.Write([]byte(err.Error())); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logger.GetLogger().Error(err.Error(), zap.String("URL", r.RequestURI))
			}
		default:
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte(err.Error())); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logger.GetLogger().Error(err.Error(), zap.String("URL", r.RequestURI))
			}
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write([]byte(toBeDeleted.ID)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.GetLogger().Error(err.Error(), zap.String("URL", r.RequestURI))
	}
}
