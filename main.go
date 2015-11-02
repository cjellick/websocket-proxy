package main

import (
	log "github.com/Sirupsen/logrus"

	"github.com/rancherio/websocket-proxy/proxy"
)

func main() {

	conf, err := proxy.GetConfig()
	if err != nil {
		log.WithField("error", err).Fatal("Error getting config.")
	}

	p := &proxy.ProxyStarter{
		BackendPaths:     []string{"/v1/connectbackend"},
		FrontendPaths:    []string{"/v1/{logs:logs}/", "/v1/{stats:stats}", "/v1/{stats:stats}/{statsid}", "/v1/exec/", "/v1/dockersocket/"},
		StatsPaths:       []string{"/v1/{hoststats:hoststats(\\/project)?(\\/)?}", "/v1/{containerstats:containerstats(\\/service)?(\\/)?}", "/v1/{containerstats:containerstats}/{containerid}"},
		CattleProxyPaths: []string{"/{cattle-proxy:.*}"},
		Config:           conf,
	}

	log.Infof("Starting websocket proxy. Listening on [%s], Proxying to cattle API at [%s], Monitoring parent pid [%v].",
		conf.ListenAddr, conf.CattleAddr, conf.ParentPid)

	err = p.StartProxy()

	log.WithFields(log.Fields{"error": err}).Info("Exiting proxy.")
}
