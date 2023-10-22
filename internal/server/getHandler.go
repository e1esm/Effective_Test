package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/e1esm/Effective_Test/internal/models/options"
	"github.com/e1esm/Effective_Test/pkg/utils/logger"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
)

func (hs *HttpServer) Get(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	sex := r.URL.Query().Get("sex")
	age, err := strconv.Atoi(r.URL.Query().Get("age"))
	if err != nil {
		if r.URL.Query().Get("age") != "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	name := r.URL.Query().Get("name")
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	nations := strings.Split(r.URL.Query().Get("nations"), ",")

	options := options.NewQueryOptions(options.UserOptions{
		sex, age, name, limit, offset,
	},
		options.NewNationalityOptions(nations),
	)
	users, err := hs.userService.Get(context.Background(), options)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err = w.Write([]byte(fmt.Sprintf("Eerror while querying: %v", err.Error()))); err != nil {
			logger.GetLogger().Error("Error while writing", zap.String("err", err.Error()))
		}
		return
	}

	bytes, err := json.Marshal(&users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte("Error while marshalling fetched data")); err != nil {
			logger.GetLogger().Error("Marshaling error")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(bytes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.GetLogger().Error(err.Error())
	}
}
