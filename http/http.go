package http

import (
	"encoding/json"
	"net/http"

	"carbon-go-relay/global"
	"carbon-go-relay/utils"
)

//Dto define htttp response data struct
type Dto struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

//Start to start a http server
func Start() {
	go startHTTPServer()
}

func startHTTPServer() {
	if !global.Config().HTTP.Enabled {
		return
	}

	addr := global.Config().HTTP.Listen
	if addr == "" {
		return
	}

	configCommonRoutes()

	s := &http.Server{
		Addr:           addr,
		MaxHeaderBytes: 1 << 30,
	}

	utils.Zlog.Info("http.startHttpServer ok, listening", addr)
	utils.Zlog.Fatal(s.ListenAndServe())
}

//RenderJSON render data to json format
func RenderJSON(w http.ResponseWriter, v interface{}) {
	bs, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bs)
}

//RenderDataJSON give json response
func RenderDataJSON(w http.ResponseWriter, data interface{}) {
	RenderJSON(w, Dto{Msg: "success", Data: data})
}
