package sender

import (
	"net"
	"regexp"
	"sync"

	"carbon-go-relay/global"
	"carbon-go-relay/utils"

	nlist "github.com/toolkits/container/list"
)

//ConnPattern give connection and pattern
type ConnPattern struct {
	Conn      net.Conn
	Patterns  []*regexp.Regexp
	AliasName string
	Data      *nlist.SafeListLimited
}

//ApmStatsMap give metric map
type ApmStatsMap struct {
	sync.RWMutex
	Data map[string]int
}

//GlobalConnPatterns includes all relay connection and relavent patterns
var (
	GlobalConnPatterns []ConnPattern
	//GlobalStatsMap     ApmStatsMap
	GlobalStatsMap []int32
)

//GetConnPatterns get relay patterns
func getConnPatterns() {
	cfg := global.Config()
	for _, line := range cfg.RelayCluster.RelayClusterList {
		tmpConn, err := net.Dial("tcp", line.Address)
		if err != nil {
			utils.Zlog.Warningf("address %s connected error", line.Address)
			continue
		}
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
			Conn:      tmpConn,
			Patterns:  tmpPatternList,
			AliasName: line.AliasName,
			Data:      nlist.NewSafeListLimited(configBrubeckMaxSize)})
	}
}

func getStatsGroup() {
	cfg := global.Config()
	GlobalStatsMap = make([]int32, len(cfg.StatsGroup))
}
