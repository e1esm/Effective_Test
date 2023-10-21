package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/e1esm/Effective_Test/internal/models/users"
	"github.com/e1esm/Effective_Test/pkg/utils/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
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
