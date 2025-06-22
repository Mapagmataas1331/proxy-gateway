package main

import (
	"log"
	"net/http"
	"os"
	"proxy-gateway/internal/auth"
	"proxy-gateway/internal/config"
	"proxy-gateway/internal/dashboard"
	"proxy-gateway/internal/db"
	"proxy-gateway/internal/proxy"
	"proxy-gateway/internal/watchdog"
)

func main() {
	config.Load()
	db.Init()
	go watchdog.StartQuotaMonitor()

	port := config.AppConfig.AppPort
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat("static" + r.URL.Path); err == nil {
			http.ServeFile(w, r, "static"+r.URL.Path)
			return
		}

		if r.URL.Path == "/dash" {
			auth.AdminOnly(dashboard.LogsHandler)(w, r)
			return
		}

		proxy.Handler(w, r)
	})

	log.Printf("Server running on port :%s", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}
