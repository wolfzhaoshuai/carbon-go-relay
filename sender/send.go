package sender

import (
	"carbon-go-relay/global"
	"carbon-go-relay/utils"
)

var (
	configBrubeckMaxSize int
	sendBatchSize        int
)

//Start start send job
func Start() {
	utils.Zlog.Info("Sender Start")
	cfg := global.Config()
	configBrubeckMaxSize = cfg.MaxBrubeckLength
	sendBatchSize = cfg.SendBatchSize
	getConnPatterns()
	getStatsGroup()
	checkBrubeckQueue()
	sendToRelay()
}
