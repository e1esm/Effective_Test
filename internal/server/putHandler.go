package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/e1esm/Effective_Test/internal/models/users"
	"github.com/e1esm/Effective_Test/pkg/utils/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func (hs *HttpServer) Change(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		if _, err := w.Write([]byte(fmt.Sprintf(methodErr.Error(), r.RequestURI))); err != nil {
			logger.GetLogger().Error(err.Error(), zap.String("URL", r.RequestURI))
		}
		return
	}
	type UpdateRequest struct {
		users.User
		ID uuid.UUID `json:"id"`
	}

	content, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte(fmt.Sprintf(invalidReq.Error(), r.RequestURI))); err != nil {
			logger.GetLogger().Error(err.Error(), zap.String("URL", r.RequestURI))
		}
		return
	}

	var updateRequest UpdateRequest

	if err := json.Unmarshal(content, &updateRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte(fmt.Sprintf(invalidReq.Error(), r.RequestURI))); err != nil {
			logger.GetLogger().Error(err.Error(), zap.String("URL", r.RequestURI))
		}
		return
	}

	user := hs.identityService.Identify(updateRequest.User)
	id, err := hs.userService.Update(context.WithValue(context.Background(), "id", updateRequest.ID), user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(fmt.Sprintf(updateErr.Error(), updateRequest.Name))); err != nil {
			logger.GetLogger().Error(err.Error(), zap.String("URL", r.RequestURI))
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(id.String())); err != nil {
		logger.GetLogger().Error(err.Error(), zap.String("URL", r.RequestURI))
	}
}
