package identifier

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/e1esm/Effective_Test/internal/models/users"
	"io"
	"log"
	"net/http"
)

type nationalityResponse struct {
	Count         int           `json:"count"`
	Name          string        `json:"name"`
	Nationalities []nationality `json:"country"`
}

type nationality struct {
	ID          string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

func (is *IdentityService) requestNationality(ctx context.Context, user *users.ProtectedUser) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://api.nationalize.io/?name=%s", user.GetUser().Name), nil)
	if err != nil {
		log.Println(err.Error())
	}
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err.Error())
	}

	content, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println(err.Error())
	}

	var natData nationalityResponse
	if err := json.Unmarshal(content, &natData); err != nil {
		log.Println(err.Error())
	}

	countryCodes := make([]string, len(natData.Nationalities))

	for i := 0; i < len(natData.Nationalities); i++ {
		countryCodes[i] = natData.Nationalities[i].ID
	}

	user.SetNationality(countryCodes)

}
