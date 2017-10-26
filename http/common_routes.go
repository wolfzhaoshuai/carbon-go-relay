package http

import (
	"fmt"
	"net/http"
	"strings"

	"carbon-go-relay/global"
)

func configCommonRoutes() {
	http.HandleFunc("/check_health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok\n"))
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("%s\n", global.VERSION)))
	})

	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJSON(w, global.Config())
	})

	http.HandleFunc("/config/reload", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RemoteAddr, "127.0.0.1") {
			global.ParseConfig(global.ConfigFile)
			RenderDataJSON(w, "ok")
		} else {
			RenderDataJSON(w, "no privilege")
		}
	})
}
