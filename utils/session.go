package utils

import (
	"net/http"
	"github.com/astaxie/session"
	_ "github.com/astaxie/session/providers/memory"
)

var globalSessions *session.Manager

func init() {
	globalSessions, _ = session.NewManager("memory", "gosessionid", 3600)
	go globalSessions.GC()
}

func GetSession(w http.ResponseWriter, r *http.Request) session.Session {
	session := globalSessions.SessionStart(w, r)
	return session
}

