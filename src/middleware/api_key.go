package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/jaxkodex/go-gin-template/src/config"
)

type apiKeyVerifier struct {
	validKeys map[string]struct{}
}

func newAPIKeyVerifier(cfg *config.Config) *apiKeyVerifier {
	keys := make(map[string]struct{}, len(cfg.APIKeys))
	for _, k := range cfg.APIKeys {
		keys[k] = struct{}{}
	}
	return &apiKeyVerifier{validKeys: keys}
}

// verify checks the X-API-Key header against the valid keys set.
// Returns ok=true and the key itself as keyID on success.
func (v *apiKeyVerifier) verify(c *gin.Context) (ok bool, keyID string, err error) {
	key := c.GetHeader("X-API-Key")
	if key == "" {
		return false, "", errors.New("unauthorized")
	}
	if _, found := v.validKeys[key]; !found {
		return false, "", errors.New("unauthorized")
	}
	return true, key, nil
}
