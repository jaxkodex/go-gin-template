package config

import (
	"os"
	"strings"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	ServerPort              string
	AuthFirebaseEnabled     bool
	FirebaseCredentialsFile string // path to service account JSON
	FirebaseCredentialsJSON string // inline JSON (alternative to file)
	AuthAPIKeyEnabled       bool
	APIKeys                 []string // parsed from comma-separated AUTH_API_KEYS
}

// Load reads environment variables and returns a populated Config.
func Load() *Config {
	return &Config{
		ServerPort:              getEnv("SERVER_PORT", "8000"),
		AuthFirebaseEnabled:     getEnvBool("AUTH_FIREBASE_ENABLED", false),
		FirebaseCredentialsFile: getEnv("FIREBASE_CREDENTIALS_FILE", ""),
		FirebaseCredentialsJSON: getEnv("FIREBASE_CREDENTIALS_JSON", ""),
		AuthAPIKeyEnabled:       getEnvBool("AUTH_API_KEY_ENABLED", false),
		APIKeys:                 parseCommaSeparated(getEnv("AUTH_API_KEYS", "")),
	}
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	v = strings.ToLower(strings.TrimSpace(v))
	return v == "true" || v == "1" || v == "yes"
}

func parseCommaSeparated(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
