package dashboard

import (
	"html/template"
	"net/http"
	"proxy-gateway/internal/db"
	"time"
)

type LogEntry struct {
	Time     string
	IP       string
	Path     string
	Port     string
	Protocol string
	AuthOK   bool
	Action   string
}

func LogsHandler(w http.ResponseWriter, r *http.Request) {
	db.Mutex.Lock()
	rows, err := db.DB.Query(`
		SELECT logs.time, ips.ip, paths.path, logs.port, logs.protocol, logs.auth_ok, actions.action
		FROM logs
		JOIN ips ON logs.ip_id = ips.id
		JOIN paths ON logs.path_id = paths.id
		JOIN actions ON logs.action_id = actions.id
		ORDER BY logs.time DESC
	`)
	db.Mutex.Unlock()

	if err != nil {
		http.Error(w, "Failed to load logs", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var entries []LogEntry
	for rows.Next() {
		var tUnix int64
		var entry LogEntry
		if err := rows.Scan(&tUnix, &entry.IP, &entry.Path, &entry.Port, &entry.Protocol, &entry.AuthOK, &entry.Action); err == nil {
			entry.Time = time.Unix(tUnix, 0).Format("2006-01-02 15:04:05")
			entries = append(entries, entry)
		}
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/dashboard.html"))
	tmpl.ExecuteTemplate(w, "layout", entries)
}
