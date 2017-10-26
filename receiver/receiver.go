package receiver

import (
	"carbon-go-relay/receiver/socket"
	"carbon-go-relay/utils"
)

//Start receiver
func Start() {
	utils.Zlog.Info("Start receiver")
	go socket.StartSocket()
}
