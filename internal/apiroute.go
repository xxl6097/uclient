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

	router.HandleFunc("/api/work/update", this.restApi.UpdatetWorkTime).Methods(http.MethodPost)
	router.HandleFunc("/api/work/add", this.restApi.AddWorkTime).Methods(http.MethodPost)
	router.HandleFunc("/api/work/del", this.restApi.DelWorkTime).Methods(http.MethodPost)
	router.HandleFunc("/api/work/get", this.restApi.GetWorkTime).Methods(http.MethodPost)
	router.HandleFunc("/api/work/tigger", this.restApi.TiggerSignCardEvent).Methods(http.MethodPost)

	router.HandleFunc("/api/network/reset", this.restApi.ResetNetwork).Methods(http.MethodPost)

	router.HandleFunc("/api/clients/get", this.restApi.GetClients).Methods(http.MethodGet)
	router.HandleFunc("/api/clients/reset", this.restApi.ResetClients).Methods(http.MethodPost)

	router.HandleFunc("/api/client/offline", this.restApi.OfflineDevice).Methods(http.MethodPost)

	router.HandleFunc("/api/ntfy/set", this.restApi.SetNtfy).Methods(http.MethodPost)

	router.HandleFunc("/api/webhook/set", this.restApi.SetWebhook).Methods(http.MethodPost)

	router.HandleFunc("/api/staticip/set", this.restApi.AddStaticIp).Methods(http.MethodPost)
	router.HandleFunc("/api/staticip/delete", this.restApi.DeleteStaticIp).Methods(http.MethodDelete)
	router.HandleFunc("/api/staticip/list", this.restApi.GetStaticIps).Methods(http.MethodGet)

	router.HandleFunc("/api/checkversion", this.restApi.ApiCheckVersion).Methods("GET")
	router.HandleFunc("/api/upgrade", this.restApi.ApiUpdate).Methods("POST")
	router.HandleFunc("/api/upgrade", this.restApi.ApiUpdate).Methods("PUT")
	router.HandleFunc("/api/version", this.restApi.ApiVersion).Methods("GET")

	router.Handle("/api/client/sse", this.restApi.GetSSE().Handler())
	//subRouter.Handle("/api/client/sse", this.sseApi)
	//httpserver.RouterUtil.AddHandleFunc(router, ihttpserver.ApiModel{
	//	Method: http.MethodPost,
	//	Path:   "/frp",
	//	Fun:    this.controller.Frp,
	//	NoAuth: false,
	//})
}
