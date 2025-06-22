package logstore

import (
	"proxy-gateway/internal/db"
)

func StoreLog(timeUnix int64, ip, path, port, protocol string, authOK bool, action string) {
	ipID := db.GetOrInsertID("ips", "ip", ip)
	pathID := db.GetOrInsertID("paths", "path", path)
	actionID := db.GetOrInsertID("actions", "action", action)

	db.Mutex.Lock()
	defer db.Mutex.Unlock()

	_, _ = db.DB.Exec(
		`INSERT INTO logs (time, ip_id, path_id, port, protocol, auth_ok, action_id)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		timeUnix, ipID, pathID, port, protocol, authOK, actionID,
	)
}
