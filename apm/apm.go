package apm

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"carbon-go-relay/global"
	"carbon-go-relay/sender"
	"carbon-go-relay/utils"
)

var metricPrefix = "carbon-c-relay.portal"

//Start application performance metrics
func Start() {
	utils.Zlog.Info("APM statistic start")
	ticker := time.NewTicker(time.Second * 10)
	go func() {
		for {
			select {
			case <-ticker.C:
				handleReceiveStatsMetric()
				handleRelayQueueMetric()
				handleGlobalStatsMetric()
			}
		}
	}()
}

func handleReceiveStatsMetric() {
	statsGroup := global.Config().StatsGroup
	for index := range sender.GlobalStatsMap {
		metricMiddle := strings.Replace(statsGroup[index], ".", "_", -1)
		metricName := fmt.Sprintf("%s.%s.traffic_stats.%s.recieve.count", metricPrefix, global.Hostname, metricMiddle)
		currentTS := int(time.Now().Unix())
		msg := fmt.Sprintf("%s %.3f %d", metricName, float64(atomic.LoadInt32(&sender.GlobalStatsMap[index])), currentTS)
		sender.FindMatchedPattern(msg)
		atomic.StoreInt32(&sender.GlobalStatsMap[index], 0)
	}
}

func handleRelayQueueMetric() {
	for _, line := range sender.GlobalConnPatterns {
		metricMiddle := strings.Replace(line.AliasName, ".", "_", -1)
		metricName := fmt.Sprintf("%s.%s.pattern_stats.%s.queue.free_usage.percent", metricPrefix, global.Hostname, metricMiddle)
		currentTS := int(time.Now().Unix())
		metricValue := float64(global.Config().MaxBrubeckLength-line.Data.Len()) / float64(global.Config().MaxBrubeckLength)
		msg := fmt.Sprintf("%s %.3f %d", metricName, metricValue, currentTS)
		sender.FindMatchedPattern(msg)

		metricRelayName := fmt.Sprintf("%s.%s.pattern_stats.%s.relays.free.value", metricPrefix, global.Hostname, metricMiddle)
		metricRelayValue := float64(cap(line.Relays) - len(line.Relays))
		relayMsg := fmt.Sprintf("%s %.3f %d", metricRelayName, metricRelayValue, currentTS)
		sender.FindMatchedPattern(relayMsg)
	}
}

func handleGlobalStatsMetric() {
	sendMetricName := fmt.Sprintf("%s.%s.global_stats.sendItems.count", metricPrefix, global.Hostname)
	sendOriginMetricName := fmt.Sprintf("%s.%s.global_stats.sendOriginItems.count", metricPrefix, global.Hostname)
	receiveMetricName := fmt.Sprintf("%s.%s.global_stats.receiveItems.count", metricPrefix, global.Hostname)
	receiveOriginMetricName := fmt.Sprintf("%s.%s.global_stats.receiveOriginItems.count", metricPrefix, global.Hostname)
	currentTS := int(time.Now().Unix())

	sendMsg := fmt.Sprintf("%s %.3f %d", sendMetricName, float64(global.GlobalSendItems), currentTS)
	sender.FindMatchedPattern(sendMsg)
	atomic.StoreInt32(&global.GlobalSendItems, 0)

	sendOriginMsg := fmt.Sprintf("%s %.3f %d", sendOriginMetricName, float64(global.GlobalOriginSendItems), currentTS)
	sender.FindMatchedPattern(sendOriginMsg)
	atomic.StoreInt32(&global.GlobalOriginSendItems, 0)

	receiveMsg := fmt.Sprintf("%s %.3f %d", receiveMetricName, float64(global.GlobalReceiveItems), currentTS)
	sender.FindMatchedPattern(receiveMsg)
	atomic.StoreInt32(&global.GlobalReceiveItems, 0)

	receiveOriginMsg := fmt.Sprintf("%s %.3f %d", receiveOriginMetricName, float64(global.GlobalOriginReceiveItems), currentTS)
	sender.FindMatchedPattern(receiveOriginMsg)
	atomic.StoreInt32(&global.GlobalOriginReceiveItems, 0)
}
