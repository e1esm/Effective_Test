package identifier

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/e1esm/Effective_Test/internal/models/nationalities"
	"github.com/e1esm/Effective_Test/internal/models/users"
	"io"
	"log"
	"net/http"
)

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

	var natData nationalities.NationalityResponse
	if err := json.Unmarshal(content, &natData); err != nil {
		log.Println(err.Error())
	}

	user.SetNationality(natData.Nationalities)

}
