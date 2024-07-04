package web

import (
	"fmt"
	"html/template"
	"net/http"
	"wb/internal/cache"
)

type Server struct {
	CacheSM *cache.CacheSM
}

func New(csm *cache.CacheSM) *Server {
	return &Server{
		CacheSM: csm,
	}
}

func (s *Server) HandlerOrders(w http.ResponseWriter, r *http.Request) {
	orderUid := r.URL.Query().Get("order_uid")
	order, err := s.CacheSM.GetFromCache(orderUid)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	tmpl, err := template.ParseFiles("/home/hello/wb/internal/web/order.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
