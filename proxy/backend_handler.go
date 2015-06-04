package proxy

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
)

type BackendHandler struct {
	proxyManager proxyManager
}

func (h *BackendHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	log.Info("A Connection!")
	hostKey := req.Header.Get("X-Cattle-HostId")
	if hostKey == "" {
		log.Errorf("No hostKey provided in request.")
		http.Error(rw, "Missing X-Cattle-HostId Header", 400)
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	ws, err := upgrader.Upgrade(rw, req, nil)
	if err != nil {
		log.Errorf("Couldn't upgrade connnection for backend [%v]. Error: [%v]", hostKey, err)
		http.Error(rw, "Failed to upgrade connection.", 500)
		return
	}

	log.Info("Adding backend for [%v]", hostKey)
	h.proxyManager.addBackend(hostKey, ws)
}
