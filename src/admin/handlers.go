package admin

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"runtime"
	"time"
)

// Handler manages admin routes and authentication
type Handler struct {
	auth      *AuthManager
	version   string
	commit    string
	buildDate string
	startTime time.Time
}

// NewHandler creates a new admin handler
func NewHandler(username, password, apiToken string, sessionTimeout int, sslEnabled bool, version, commit, buildDate string) *Handler {
	return &Handler{
		auth:      NewAuthManager(username, password, apiToken, sessionTimeout, sslEnabled),
		version:   version,
		commit:    commit,
		buildDate: buildDate,
		startTime: time.Now(),
	}
}

// RegisterRoutes registers admin routes on http.ServeMux
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// Admin web interface (session auth)
	mux.HandleFunc("/admin", h.handleAdminLogin)
	mux.HandleFunc("/admin/login", h.handleAdminLoginPost)
	mux.HandleFunc("/admin/logout", h.handleAdminLogout)
	mux.HandleFunc("/admin/dashboard", h.requireSession(h.handleAdminDashboard))
	mux.HandleFunc("/admin/settings", h.requireSession(h.handleAdminSettings))

	// Admin API (bearer token auth)
	mux.HandleFunc("/api/v1/admin/status", h.requireToken(h.handleAPIStatus))
	mux.HandleFunc("/api/v1/admin/config", h.handleAPIConfig)
	mux.HandleFunc("/api/v1/admin/reload", h.requireToken(h.handleAPIReload))
}

// Middleware for session authentication
func (h *Handler) requireSession(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, ok := h.auth.GetSessionFromRequest(r)
		if !ok {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			return
		}
		// Refresh session on activity
		h.auth.RefreshSession(session.ID)
		handler(w, r)
	}
}

// Middleware for bearer token authentication
func (h *Handler) requireToken(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := GetTokenFromRequest(r)
		if token == "" || !h.auth.ValidateAPIToken(token) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
			return
		}
		handler(w, r)
	}
}

// handleAdminLogin shows the login page
func (h *Handler) handleAdminLogin(w http.ResponseWriter, r *http.Request) {
	// Check if already logged in
	if _, ok := h.auth.GetSessionFromRequest(r); ok {
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	h.renderLoginPage(w, "")
}

// handleAdminLoginPost processes login form
func (h *Handler) handleAdminLoginPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.renderLoginPage(w, "Invalid request")
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if h.auth.Authenticate(username, password) {
		session := h.auth.CreateSession(username, r.RemoteAddr)
		h.auth.SetSessionCookie(w, session)
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	h.renderLoginPage(w, "Invalid username or password")
}

// handleAdminLogout logs out the user
func (h *Handler) handleAdminLogout(w http.ResponseWriter, r *http.Request) {
	if session, ok := h.auth.GetSessionFromRequest(r); ok {
		h.auth.DeleteSession(session.ID)
	}
	h.auth.ClearSessionCookie(w)
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// handleAdminDashboard shows the admin dashboard
func (h *Handler) handleAdminDashboard(w http.ResponseWriter, r *http.Request) {
	h.renderDashboardPage(w)
}

// handleAdminSettings shows the settings page
func (h *Handler) handleAdminSettings(w http.ResponseWriter, r *http.Request) {
	message := ""
	if r.Method == http.MethodPost {
		message = "Settings updated successfully"
	}
	h.renderSettingsPage(w, message)
}

// API Handlers

func (h *Handler) handleAPIStatus(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "healthy",
		"version":   h.version,
		"commit":    h.commit,
		"buildDate": h.buildDate,
		"uptime":    time.Since(h.startTime).String(),
		"memory": map[string]interface{}{
			"alloc":      m.Alloc,
			"totalAlloc": m.TotalAlloc,
			"sys":        m.Sys,
			"numGC":      m.NumGC,
		},
		"runtime": map[string]interface{}{
			"goroutines": runtime.NumGoroutine(),
			"cpus":       runtime.NumCPU(),
			"goVersion":  runtime.Version(),
		},
	})
}

func (h *Handler) handleAPIConfig(w http.ResponseWriter, r *http.Request) {
	token := GetTokenFromRequest(r)
	if token == "" || !h.auth.ValidateAPIToken(token) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodPut {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"version": h.version,
	})
}

func (h *Handler) handleAPIReload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "reloaded"})
}

// HTML Templates

func (h *Handler) renderLoginPage(w http.ResponseWriter, errorMsg string) {
	tmpl := template.Must(template.New("login").Parse(loginTemplate))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, map[string]interface{}{
		"Error": errorMsg,
	})
}

func (h *Handler) renderDashboardPage(w http.ResponseWriter) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	tmpl := template.Must(template.New("dashboard").Parse(dashboardTemplate))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, map[string]interface{}{
		"Version":    h.version,
		"Commit":     h.commit,
		"BuildDate":  h.buildDate,
		"MemAlloc":   fmt.Sprintf("%.2f MB", float64(m.Alloc)/1024/1024),
		"Goroutines": runtime.NumGoroutine(),
		"Uptime":     time.Since(h.startTime).Round(time.Second).String(),
	})
}

func (h *Handler) renderSettingsPage(w http.ResponseWriter, message string) {
	tmpl := template.Must(template.New("settings").Parse(settingsTemplate))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, map[string]interface{}{
		"Message": message,
	})
}

const loginTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Admin Login - Quotes API</title>
    <style>
        :root {
            --bg-color: #282a36;
            --fg-color: #f8f8f2;
            --accent: #bd93f9;
            --red: #ff5555;
            --green: #50fa7b;
            --input-bg: #44475a;
        }
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
            background: var(--bg-color);
            color: var(--fg-color);
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        .login-container {
            background: var(--input-bg);
            padding: 2rem;
            border-radius: 8px;
            width: 100%;
            max-width: 400px;
            box-shadow: 0 4px 6px rgba(0,0,0,0.3);
        }
        h1 { text-align: center; margin-bottom: 1.5rem; color: var(--accent); }
        .error { background: var(--red); color: #fff; padding: 0.75rem; border-radius: 4px; margin-bottom: 1rem; }
        label { display: block; margin-bottom: 0.5rem; font-weight: 500; }
        input[type="text"], input[type="password"] {
            width: 100%;
            padding: 0.75rem;
            border: none;
            border-radius: 4px;
            background: var(--bg-color);
            color: var(--fg-color);
            margin-bottom: 1rem;
            font-size: 1rem;
        }
        input:focus { outline: 2px solid var(--accent); }
        button {
            width: 100%;
            padding: 0.75rem;
            border: none;
            border-radius: 4px;
            background: var(--accent);
            color: var(--bg-color);
            font-size: 1rem;
            font-weight: 600;
            cursor: pointer;
            transition: opacity 0.2s;
        }
        button:hover { opacity: 0.9; }
    </style>
</head>
<body>
    <div class="login-container">
        <h1>Admin Login</h1>
        {{if .Error}}<div class="error">{{.Error}}</div>{{end}}
        <form method="POST" action="/admin/login">
            <label for="username">Username</label>
            <input type="text" id="username" name="username" required autofocus>
            <label for="password">Password</label>
            <input type="password" id="password" name="password" required>
            <button type="submit">Login</button>
        </form>
    </div>
</body>
</html>`

const dashboardTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Admin Dashboard - Quotes API</title>
    <style>
        :root {
            --bg-color: #282a36;
            --fg-color: #f8f8f2;
            --accent: #bd93f9;
            --card-bg: #44475a;
        }
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
            background: var(--bg-color);
            color: var(--fg-color);
            min-height: 100vh;
        }
        .navbar {
            background: var(--card-bg);
            padding: 1rem 2rem;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .navbar h1 { color: var(--accent); font-size: 1.5rem; }
        .navbar a { color: var(--fg-color); text-decoration: none; margin-left: 1rem; }
        .navbar a:hover { color: var(--accent); }
        .container { max-width: 1200px; margin: 2rem auto; padding: 0 1rem; }
        .cards { display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 1rem; }
        .card {
            background: var(--card-bg);
            padding: 1.5rem;
            border-radius: 8px;
        }
        .card h3 { color: var(--accent); margin-bottom: 0.5rem; }
        .card p { font-size: 1.5rem; font-weight: bold; }
    </style>
</head>
<body>
    <nav class="navbar">
        <h1>Quotes API Admin</h1>
        <div>
            <a href="/admin/dashboard">Dashboard</a>
            <a href="/admin/settings">Settings</a>
            <a href="/admin/logout">Logout</a>
        </div>
    </nav>
    <div class="container">
        <div class="cards">
            <div class="card">
                <h3>Version</h3>
                <p>{{.Version}}</p>
            </div>
            <div class="card">
                <h3>Memory Usage</h3>
                <p>{{.MemAlloc}}</p>
            </div>
            <div class="card">
                <h3>Goroutines</h3>
                <p>{{.Goroutines}}</p>
            </div>
            <div class="card">
                <h3>Uptime</h3>
                <p>{{.Uptime}}</p>
            </div>
        </div>
    </div>
</body>
</html>`

const settingsTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Settings - Quotes API</title>
    <style>
        :root {
            --bg-color: #282a36;
            --fg-color: #f8f8f2;
            --accent: #bd93f9;
            --card-bg: #44475a;
            --green: #50fa7b;
        }
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
            background: var(--bg-color);
            color: var(--fg-color);
            min-height: 100vh;
        }
        .navbar {
            background: var(--card-bg);
            padding: 1rem 2rem;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .navbar h1 { color: var(--accent); font-size: 1.5rem; }
        .navbar a { color: var(--fg-color); text-decoration: none; margin-left: 1rem; }
        .navbar a:hover { color: var(--accent); }
        .container { max-width: 800px; margin: 2rem auto; padding: 0 1rem; }
        .message { background: var(--green); color: #000; padding: 1rem; border-radius: 4px; margin-bottom: 1rem; }
        .card {
            background: var(--card-bg);
            padding: 1.5rem;
            border-radius: 8px;
        }
        .card h2 { color: var(--accent); margin-bottom: 1rem; }
    </style>
</head>
<body>
    <nav class="navbar">
        <h1>Quotes API Admin</h1>
        <div>
            <a href="/admin/dashboard">Dashboard</a>
            <a href="/admin/settings">Settings</a>
            <a href="/admin/logout">Logout</a>
        </div>
    </nav>
    <div class="container">
        {{if .Message}}<div class="message">{{.Message}}</div>{{end}}
        <div class="card">
            <h2>Settings</h2>
            <p>Settings configuration coming soon.</p>
        </div>
    </div>
</body>
</html>`
