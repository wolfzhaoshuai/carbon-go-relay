package sender

import (
	"fmt"
	"sync/atomic"
	"time"

	"carbon-go-relay/utils"

	"carbon-go-relay/global"
)

func cronCheckHealth() {
	ticker := time.NewTicker(time.Second * time.Duration(10))
	go func() {
		for {
			select {
			case <-ticker.C:
				checkHealth()
			}
		}
	}()
}

func checkHealth() {
	for _, line := range GlobalConnPatterns {
		atomic.AddInt32(&global.GlobalOriginSendItems, int32(1))
		prefix := "carbon-c-relay.check_health"
		metricName := fmt.Sprintf("%s.%s.%s.count", prefix, global.Hostname, line.AliasName)
		currentTS := int(time.Now().Unix())
		msg := fmt.Sprintf("%s %.3f %d", metricName, float64(1), currentTS)
		_, err := line.Conn.Write(utils.StringToBytes(msg + "\n"))
		if err != nil {
			utils.Zlog.Critical(line.Patterns, " connection write error")
			continue
		}
		if global.Config().Debug {
			utils.Zlog.Debug("check health send relay: ", line.Patterns, msg)
		}

		atomic.AddInt32(&global.GlobalSendItems, int32(1))
	}
}
