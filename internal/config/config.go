package config

import (
	"log"
	"os"
	"time"
)

type Config struct {
	AppPort           string
	ProtectedPorts    map[string]bool
	UnprotectedPorts  map[string]bool
	HomeIP            string
	AuthToken         string
	DashboardPass     string
	ConnectionTimeout time.Duration
	FlyApiToken       string
}

var AppConfig Config

func Load() {
	AppConfig = Config{
		AppPort: getEnv("PORT", "8080"),
		ProtectedPorts: map[string]bool{
			"80":  true,
			"443": true,
		},
		UnprotectedPorts: map[string]bool{
			"25565": true,
			"24454": true,
		},
		HomeIP:            getEnv("HOME_IP", "0.0.0.0"),
		AuthToken:         getEnv("AUTH_TOKEN", "securetoken"),
		DashboardPass:     getEnv("DASHBOARD_PASSWORD", "securepassword"),
		ConnectionTimeout: 3 * time.Second,
		FlyApiToken:       getEnv("FLY_API_TOKEN", ""),
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Printf("[config] %s not set, using default: %s", key, fallback)
		return fallback
	}
	return val
}
