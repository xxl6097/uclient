package internal

import (
	"github.com/gorilla/mux"
	"github.com/xxl6097/go-http/pkg/ihttpserver"
	"net/http"
)

type ApiRoute struct {
	restApi *Api
}

func NewRoute(ctl *Api) ihttpserver.IRoute {
	opt := &ApiRoute{
		restApi: ctl,
	}
	return opt
}

func (this *ApiRoute) Setup(router *mux.Router) {

	router.HandleFunc("/api/get/status", this.restApi.GetStatus).Methods(http.MethodGet)
	router.HandleFunc("/api/clear", this.restApi.Clear).Methods(http.MethodDelete)
	router.HandleFunc("/api/nick/set", this.restApi.SetNick).Methods(http.MethodPost)

	router.HandleFunc("/api/network/reset", this.restApi.ResetNetwork).Methods(http.MethodPost)

	router.HandleFunc("/api/clients/get", this.restApi.GetClients).Methods(http.MethodGet)
	router.HandleFunc("/api/clients/reset", this.restApi.ResetClients).Methods(http.MethodPost)

	router.HandleFunc("/api/staticip/set", this.restApi.SetStaticIp).Methods(http.MethodPost)
	router.HandleFunc("/api/staticip/delete", this.restApi.DeleteStaticIp).Methods(http.MethodDelete)
	router.HandleFunc("/api/staticip/list", this.restApi.GetStaticIps).Methods(http.MethodGet)

	router.HandleFunc("/api/checkversion", this.restApi.ApiCheckVersion).Methods("GET")
	router.HandleFunc("/api/upgrade", this.restApi.ApiUpdate).Methods("POST")
	router.HandleFunc("/api/upgrade", this.restApi.ApiUpdate).Methods("PUT")
	router.HandleFunc("/api/version", this.restApi.ApiVersion).Methods("GET")

	router.Handle("/api/client/sse", this.restApi.GetSSE())
	//subRouter.Handle("/api/client/sse", this.sseApi)
	//httpserver.RouterUtil.AddHandleFunc(router, ihttpserver.ApiModel{
	//	Method: http.MethodPost,
	//	Path:   "/frp",
	//	Fun:    this.controller.Frp,
	//	NoAuth: false,
	//})
}
