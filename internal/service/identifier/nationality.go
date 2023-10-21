package identifier

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/e1esm/Effective_Test/internal/models/nationalities"
	"github.com/e1esm/Effective_Test/internal/models/users"
	"github.com/e1esm/Effective_Test/pkg/utils/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func (is *IdentityService) requestNationality(ctx context.Context, user *users.ProtectedUser) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://api.nationalize.io/?name=%s", user.GetUser().Name), nil)
	if err != nil {
		logger.GetLogger().Error("Error while building request",
			zap.String("err", err.Error()))
		return
	}
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		logger.GetLogger().Error("Failed to send request",
			zap.String("err", err.Error()), zap.String("url", request.URL.RequestURI()))
		return
	}

	content, err := io.ReadAll(resp.Body)

	if err != nil {
		logger.GetLogger().Error("Invalid response body",
			zap.String("err", err.Error()))
		return
	}

	var natData nationalities.NationalityResponse
	if err := json.Unmarshal(content, &natData); err != nil {
		logger.GetLogger().Error("Marshalling error",
			zap.String("err", err.Error()))
		return
	}

	user.SetNationality(natData.Nationalities)

}
