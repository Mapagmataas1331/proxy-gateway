package watchdog

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"proxy-gateway/internal/config"
	"time"
)

type FlyVMInfo struct {
	Current float64 `json:"current"`
	Limit   float64 `json:"limit"`
}

func CheckQuota(threshold float64) (shouldShutdown bool, usage FlyVMInfo, err error) {
	token := config.AppConfig.FlyApiToken
	if token == "" {
		return false, usage, fmt.Errorf("missing FLY_API_TOKEN")
	}

	req, err := http.NewRequest("GET", "https://api.fly.io/api/v1/metrics/vm", nil)
	if err != nil {
		return false, usage, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, usage, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, usage, fmt.Errorf("fly API returned status: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&usage); err != nil {
		return false, usage, err
	}

	if usage.Limit == 0 {
		return false, usage, fmt.Errorf("invalid quota limit")
	}

	percentUsed := usage.Current / usage.Limit
	return percentUsed >= threshold, usage, nil
}

func StartQuotaMonitor() {
	go func() {
		for {
			shouldShutdown, usage, err := CheckQuota(0.9)
			if err != nil {
				log.Printf("Quota check failed: %v", err)
			} else if shouldShutdown {
				log.Printf("Usage exceeded: %.2f%%", usage.Current/usage.Limit*100)
				ShutdownGracefully(10 * time.Second)
				return
			}
			time.Sleep(5 * time.Minute)
		}
	}()
}
