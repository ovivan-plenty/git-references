package version

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

// Auth authentication details
type Auth struct {
	PlentyID int
	UserID   int
	Username string
}

// GetAuthenticatedUser extracts user from auth header
func GetAuthenticatedUser(r *http.Request) (*Auth, error) {

	header := r.Header.Get("X-Authenticated-Userid")

	if len(header) < 3 {
		return nil, errors.New("no authentication header provided")
	}

	parts := strings.Split(header, "::")

	if len(parts) != 3 {
		return nil, errors.New("malformed authentication header")
	}

	au := Auth{Username: parts[0]}
	au.UserID, _ = strconv.Atoi(parts[1])
	au.PlentyID, _ = strconv.Atoi(parts[2])

	return &au, nil
}
