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

func (is *IdentityService) requestAge(ctx context.Context, user *users.ProtectedUser) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://api.agify.io/?name=%s", user.GetUser().Name), nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	client := http.Client{}
	resp, err := client.Do(request)

	if err != nil {
		log.Println(err.Error())
		return
	}

	content, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println(err.Error())
		return
	}

	var data map[string]interface{}

	if err := json.Unmarshal(content, &data); err != nil {
		log.Println(err.Error())
	}

	user.SetAge(int(data["age"].(float64)))
}
