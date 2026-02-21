package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jaxkodex/go-gin-template/src/config"
)

const (
	ctxKeyUserID  = "auth_user_id"
	ctxKeyMethod  = "auth_method"
	ctxKeyClaims  = "auth_claims"

	authMethodFirebase = "firebase"
	authMethodAPIKey   = "api_key"
)

// AuthMiddleware orchestrates authentication verifiers.
type AuthMiddleware struct {
	firebaseVerifier *firebaseVerifier
	apiKeyVerifier   *apiKeyVerifier
	enabled          bool
}

// NewAuth initializes enabled verifiers based on cfg.
// Returns an error (fail-fast) if Firebase is enabled but credentials are invalid.
func NewAuth(cfg *config.Config) (*AuthMiddleware, error) {
	m := &AuthMiddleware{}

	if cfg.AuthFirebaseEnabled {
		fv, err := newFirebaseVerifier(cfg)
		if err != nil {
			return nil, err
		}
		m.firebaseVerifier = fv
		m.enabled = true
	}

	if cfg.AuthAPIKeyEnabled {
		m.apiKeyVerifier = newAPIKeyVerifier(cfg)
		m.enabled = true
	}

	return m, nil
}

// Authenticate returns a Gin handler that enforces authentication.
// If both verifiers are disabled it returns a no-op pass-through handler.
// Firebase is tried first; API key is tried second.
// Any successful verification allows the request through.
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	if !m.enabled {
		return func(c *gin.Context) { c.Next() }
	}

	return func(c *gin.Context) {
		// Try Firebase first.
		if m.firebaseVerifier != nil {
			ok, uid, claims, _ := m.firebaseVerifier.verify(c)
			if ok {
				c.Set(ctxKeyUserID, uid)
				c.Set(ctxKeyMethod, authMethodFirebase)
				c.Set(ctxKeyClaims, claims)
				c.Next()
				return
			}
		}

		// Try API key second.
		if m.apiKeyVerifier != nil {
			ok, keyID, _ := m.apiKeyVerifier.verify(c)
			if ok {
				c.Set(ctxKeyUserID, keyID)
				c.Set(ctxKeyMethod, authMethodAPIKey)
				c.Next()
				return
			}
		}

		// All verifiers failed.
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
			"code":  "AUTH_REQUIRED",
		})
	}
}

// GetAuthUserID returns the authenticated user ID stored in the context.
func GetAuthUserID(c *gin.Context) (string, bool) {
	v, ok := c.Get(ctxKeyUserID)
	if !ok {
		return "", false
	}
	id, ok := v.(string)
	return id, ok
}

// GetAuthMethod returns the authentication method used ("firebase" or "api_key").
func GetAuthMethod(c *gin.Context) (string, bool) {
	v, ok := c.Get(ctxKeyMethod)
	if !ok {
		return "", false
	}
	method, ok := v.(string)
	return method, ok
}

// GetAuthClaims returns the Firebase token claims stored in the context.
// Only populated when the Firebase verifier was used.
func GetAuthClaims(c *gin.Context) (map[string]interface{}, bool) {
	v, ok := c.Get(ctxKeyClaims)
	if !ok {
		return nil, false
	}
	claims, ok := v.(map[string]interface{})
	return claims, ok
}
