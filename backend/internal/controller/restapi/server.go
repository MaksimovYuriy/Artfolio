package restapi

import (
	"net"
	"net/http"

	"github.com/maksimovyuriy/artfolio/backend/internal/config"
)

func NewServer(serverCfg config.HTTPConfig, router http.Handler) *http.Server {
	server := &http.Server{
		Addr:              net.JoinHostPort(serverCfg.Address, serverCfg.Port),
		Handler:           router,
		ReadHeaderTimeout: serverCfg.ReadHeaderTimeout,
		ReadTimeout:       serverCfg.ReadTimeout,
		WriteTimeout:      serverCfg.WriteTimeout,
		IdleTimeout:       serverCfg.IdleTimeout,
	}
	return server
}
