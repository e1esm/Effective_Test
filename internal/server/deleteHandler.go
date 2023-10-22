package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/e1esm/Effective_Test/internal/repository/postgres"
	"github.com/e1esm/Effective_Test/pkg/utils/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
)

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
	defer r.Body.Close()
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
