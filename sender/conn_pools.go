package sender

import (
	"net"
	"regexp"

	"carbon-go-relay/global"
	"carbon-go-relay/utils"

	nlist "github.com/toolkits/container/list"
)

//ConnPattern give connection and pattern
type ConnPattern struct {
	Address   string
	Patterns  []*regexp.Regexp
	AliasName string
	Data      *nlist.SafeListLimited
	Relays    chan net.Conn
}

//GlobalConnPatterns includes all relay connection and relavent patterns
var (
	GlobalConnPatterns []ConnPattern
	GlobalStatsMap     []int32
)

func initConnection(address string) net.Conn {
	tmpConn, err := net.Dial("tcp", address)
	if err != nil {
		utils.Zlog.Warningf("address %s connected error", address)
		return nil
	}
	return tmpConn
}

func recycleConnection(line ConnPattern, connection net.Conn) {
	if len(line.Relays) < cap(line.Relays) {
		line.Relays <- connection
	} else {
		dropConn := <-line.Relays
		dropConn.Close()
		line.Relays <- connection
	}
}

//GetConnPatterns get relay patterns
func getConnPatterns() {
	cfg := global.Config()
	for _, line := range cfg.RelayCluster.RelayClusterList {
		tmpPatternList := make([]*regexp.Regexp, 0)
		for _, patternExpression := range line.Patterns {
			tmpPattern, err := regexp.Compile(patternExpression)
			if err != nil {
				utils.Zlog.Criticalf("%s can not compile", patternExpression)
				continue
			}
			tmpPatternList = append(tmpPatternList, tmpPattern)
		}

		GlobalConnPatterns = append(GlobalConnPatterns, ConnPattern{
			Address:   line.Address,
			Patterns:  tmpPatternList,
			AliasName: line.AliasName,
			Data:      nlist.NewSafeListLimited(configBrubeckMaxSize),
			Relays:    make(chan net.Conn, line.MaxWorkerNumber)})
	}
}

func getStatsGroup() {
	cfg := global.Config()
	GlobalStatsMap = make([]int32, len(cfg.StatsGroup))
}
