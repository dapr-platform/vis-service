package api

import "github.com/go-chi/chi/v5"

func InitRoute(r chi.Router) {
	InitTopologyRoute(r)
	InitDashboardDataRoute(r)
	InitDashboardRoute(r)
}
