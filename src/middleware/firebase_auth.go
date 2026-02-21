package middleware

import (
	"context"
	"errors"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"github.com/jaxkodex/go-gin-template/src/config"
	"google.golang.org/api/option"
)

type firebaseVerifier struct {
	client *auth.Client
}

// newFirebaseVerifier initializes a Firebase app and auth client using credentials in this priority:
//  1. FIREBASE_CREDENTIALS_FILE (path to service account JSON)
//  2. FIREBASE_CREDENTIALS_JSON (inline JSON)
//  3. Application Default Credentials (GOOGLE_APPLICATION_CREDENTIALS or GCP metadata)
func newFirebaseVerifier(cfg *config.Config) (*firebaseVerifier, error) {
	ctx := context.Background()

	var opts []option.ClientOption

	switch {
	case cfg.FirebaseCredentialsFile != "":
		opts = append(opts, option.WithCredentialsFile(cfg.FirebaseCredentialsFile))
	case cfg.FirebaseCredentialsJSON != "":
		opts = append(opts, option.WithCredentialsJSON([]byte(cfg.FirebaseCredentialsJSON)))
	}
	// No opts → Firebase will use Application Default Credentials automatically.

	app, err := firebase.NewApp(ctx, nil, opts...)
	if err != nil {
		return nil, err
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	return &firebaseVerifier{client: client}, nil
}

// verify extracts and validates a Firebase ID token from the Authorization header.
// Returns ok=true, the Firebase UID, and the token claims on success.
// Does not call c.Abort — the orchestrator decides.
func (v *firebaseVerifier) verify(c *gin.Context) (ok bool, userID string, claims map[string]interface{}, err error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return false, "", nil, errors.New("missing Authorization header")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return false, "", nil, errors.New("invalid Authorization header format")
	}

	token, err := v.client.VerifyIDToken(c.Request.Context(), parts[1])
	if err != nil {
		return false, "", nil, err
	}

	return true, token.UID, token.Claims, nil
}
