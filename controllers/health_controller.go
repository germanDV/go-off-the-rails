package controllers

import "net/http"

type HealthController struct {
	mdw *MiddlewareChain
}

func NewHealthController(mdw *MiddlewareChain) *HealthController {
	return &HealthController{
		mdw: mdw,
	}
}

func (c *HealthController) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /health", c.mdw.RBACNone(c.Health))
}

func (c *HealthController) Health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
