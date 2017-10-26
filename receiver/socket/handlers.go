package socket

import (
	"bufio"
	"net"
	"strings"
	"sync/atomic"
	"time"

	"carbon-go-relay/global"
	"carbon-go-relay/sender"
	"carbon-go-relay/utils"
)

func socketHandle(conn net.Conn) {
	defer conn.Close()
	buf := bufio.NewReader(conn)
	cfg := global.Config()
	readTimeout := time.Duration(cfg.Socket.Timeout) * time.Second
	for {
		conn.SetReadDeadline(time.Now().Add(readTimeout))
		line, err := buf.ReadString('\n')
		if err != nil {
			utils.Zlog.Warning("There is no data to be read")
			break
		}
		line = strings.Trim(line, "\n")

		if global.Config().Debug {
			utils.Zlog.Debug("receive: ", line)
		}

		if line == "" {
			continue
		}

		atomic.AddInt32(&global.GlobalOriginReceiveItems, 1)
		sender.FindMatchedPattern(line)
		atomic.AddInt32(&global.GlobalReceiveItems, 1)
	}

	return
}
