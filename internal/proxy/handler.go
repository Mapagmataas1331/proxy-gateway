package proxy

import (
	"fmt"
	"html/template"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"proxy-gateway/internal/config"
	"proxy-gateway/internal/logstore"
)

var (
	authToken        = config.AppConfig.AuthToken
	homeIP           = config.AppConfig.HomeIP
	protectedPorts   = config.AppConfig.ProtectedPorts
	unprotectedPorts = config.AppConfig.UnprotectedPorts
	timeout          = config.AppConfig.ConnectionTimeout
)

type PageData struct {
	Title   string
	Message string
}

func renderMessage(w http.ResponseWriter, title, message string) {
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/message.html"))
	tmpl.ExecuteTemplate(w, "layout", PageData{Title: title, Message: message})
}

func isPortOpen(ip, port string) bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, port), timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func Handler(w http.ResponseWriter, r *http.Request) {
	port := r.URL.Port()
	if port == "" {
		if r.TLS != nil {
			port = "443"
		} else {
			port = "80"
		}
	}

	ip := strings.Split(r.RemoteAddr, ":")[0]
	protocol := "http"
	if r.TLS != nil {
		protocol = "https"
	}

	authOK := true
	action := "forwarded"

	if protectedPorts[port] {
		if r.Header.Get("X-Auth") != authToken {
			authOK = false
			action = "denied: bad token"
			renderMessage(w, "Authentication Required", "You must provide a valid token in the X-Auth header.")
			logstore.StoreLog(time.Now().Unix(), ip, r.URL.Path, port, protocol, authOK, action)
			return
		}
	} else if !unprotectedPorts[port] {
		action = "denied: invalid port"
		renderMessage(w, "Invalid Port", fmt.Sprintf("Port %s is not open for external access.", port))
		logstore.StoreLog(time.Now().Unix(), ip, r.URL.Path, port, protocol, authOK, action)
		return
	}

	if !isPortOpen(homeIP, port) {
		action = "denied: port closed"
		renderMessage(w, "Server Offline", fmt.Sprintf("Port %s appears to be closed or your PC is offline.", port))
		logstore.StoreLog(time.Now().Unix(), ip, r.URL.Path, port, protocol, authOK, action)
		return
	}

	logstore.StoreLog(time.Now().Unix(), ip, r.URL.Path, port, protocol, authOK, action)

	target := fmt.Sprintf("%s://%s:%s", protocol, homeIP, port)
	u, err := url.Parse(target)
	if err != nil {
		renderMessage(w, "Proxy Error", "Failed to parse backend address.")
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		renderMessage(w, "Backend Error", "Your app is not currently hosted on this port.")
	}
	proxy.ServeHTTP(w, r)
}
