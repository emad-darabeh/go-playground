package itemshttp

import (
	"fmt"
	"net/http"

	"modern-api-service/internal/httptransport"
)

func parseID(r *http.Request) (uint64, error) {
	raw := r.PathValue("id")
	var id uint64
	_, err := fmt.Sscanf(raw, "%d", &id)
	if err != nil {
		return 0, httptransport.BadRequestError("invalid id")
	}
	return id, nil
}
