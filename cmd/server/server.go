package main

import (
	"contestive/config"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/go-chi/chi/v5"
)

func runServer() {
	cfg, err := config.ReadConfig("config.json")
	if err != nil {
		logger.Fatalf("failed to load configuration: %v", err)
	}

	router := chi.NewRouter()
	router.Mount("/api", HandleAPI(cfg))

	if cfg.FrontEndProxy != "" {
		url, err := url.Parse(cfg.FrontEndProxy)
		if err != nil {
			logger.Fatalf("FrontEndProxy url parse error: %v", err)
		} else {
			router.Mount("/", httputil.NewSingleHostReverseProxy(url))
		}
	} else if cfg.FrontEndDir != "" {
		router.Mount("/", FrontEndServer(cfg.FrontEndDir))
	}

	logger.Printf("Listening on %v", cfg.Address)
	if err := http.ListenAndServe(cfg.Address, router); err != nil {
		logger.Fatalf("listen and serve failed: %v", err)
	}
}
