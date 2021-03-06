package sender

import (
	"net"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	"carbon-go-relay/global"
	"carbon-go-relay/utils"
)

//Define global msg queue from brubeck-end
var (
	BrubeckQueue        = make(chan string, configBrubeckMaxSize)
	BrubeckChanOverflow = make(chan int)
)

//SendToRelay finally push msgs to carbon-c-relay
func sendToRelay() {
	ticker := time.NewTicker(time.Millisecond * time.Duration(global.Config().SendPeriodMilliseconds))
	go func() {
		for {
			select {
			case <-ticker.C:
				handleSendRelay()
			case <-BrubeckChanOverflow:
				handleSendRelay()
			}
		}
	}()
}

//FindMatchedPattern find matched pattern
func FindMatchedPattern(line string) {
	procStatistic(line)
	matchLine(line)
}

func trueHandleSendRelay(line ConnPattern) {
	tmpData := line.Data.PopBackBy(sendBatchSize)
	count := len(tmpData)
	atomic.AddInt32(&global.GlobalOriginSendItems, int32(count))
	if count > 0 {
		var tmpConn net.Conn
		if len(line.Relays) > 0 {
			tmpConn = <-line.Relays
		} else {
			utils.Zlog.Criticalf("pattern %s has no connection to be used", line.AliasName)
			line.Data.PushFrontBatch(tmpData)
			atomic.AddInt32(&global.GlobalOriginSendItems, -1*int32(count))
			return
		}
		pushItems := make([]string, count)
		for i := 0; i < count; i++ {
			pushItems[i] = tmpData[i].(string)
		}

		if global.Config().Debug {
			utils.Zlog.Debug("send relay: ", line.Patterns, pushItems)
		}

		msgs := strings.Join(pushItems, "\n")
		_, err := tmpConn.Write(utils.StringToBytes(msgs + "\n"))
		if err != nil {
			//one retry
			tmpConn.Close()
			tmpConn = initConnection(line.Address)
			if tmpConn == nil {
				line.Data.PushFrontBatch(tmpData)
				atomic.AddInt32(&global.GlobalOriginSendItems, -1*int32(count))
				return
			}
			_, err = tmpConn.Write(utils.StringToBytes(msgs + "\n"))
			if err != nil {
				utils.Zlog.Critical(line.Patterns, " connection write error")
				return
			}
			atomic.AddInt32(&global.GlobalSendItems, int32(count))
			recycleConnection(line, tmpConn)
			return
		}
		atomic.AddInt32(&global.GlobalSendItems, int32(count))
		recycleConnection(line, tmpConn)
		return
	}
}

func handleSendRelay() {
	for _, line := range GlobalConnPatterns {
		go trueHandleSendRelay(line)
	}
}

func procStatistic(msg string) {
	findMatched := false
	globalStatsGroupLength := len(global.Config().StatsGroup)
	for index, line := range global.Config().StatsGroup {
		if strings.HasPrefix(msg, line) {
			atomic.AddInt32(&GlobalStatsMap[index], 1)
			findMatched = true
			break
		}
	}
	if !findMatched {
		atomic.AddInt32(&GlobalStatsMap[globalStatsGroupLength-1], 1)
	}
}

//SendLine push line to matched connections data slice
func matchLine(msg string) {
	for index, line := range GlobalConnPatterns {
		matched := isMatched(line.Patterns, msg)
		if matched {
			if global.Config().Debug {
				utils.Zlog.Debug("send: ", GlobalConnPatterns[index].Patterns, GlobalConnPatterns[index].Data, msg)
			}

			pushSuccess := GlobalConnPatterns[index].Data.PushFront(msg)
			if !pushSuccess {
				utils.Zlog.Errorf("push %s error\n", msg)
				continue
			}
			break
		}
	}
}

func isMatched(patterns []*regexp.Regexp, msg string) bool {
	res := false
	for _, pattern := range patterns {
		if pattern.MatchString(msg) {
			res = true
			break
		} else {
			continue
		}
	}
	return res
}

//CheckBrubeckQueue check brubeck queue is overflow or not every 10ms
func checkBrubeckQueue() {
	ticker := time.NewTicker(time.Millisecond * 10000)
	go func() {
		for _ = range ticker.C {
			if len(BrubeckQueue) >= configBrubeckMaxSize {
				BrubeckChanOverflow <- 1
			}
		}
	}()
}
