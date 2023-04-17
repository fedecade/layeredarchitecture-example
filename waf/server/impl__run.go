package server

import "net/http"

func (i *impl) Run() error {
	return http.ListenAndServe(i.port, i.router.Handler())
}
