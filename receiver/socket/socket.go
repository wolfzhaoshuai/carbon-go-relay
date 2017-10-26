package socket

import (
	"net"

	"carbon-go-relay/global"
	"carbon-go-relay/utils"
)

//StartSocket show how to start socket server
func StartSocket() {
	addr := global.Config().Socket.Listen
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		utils.Zlog.Fatalf("net.ResolveTCPAddr fail: %s", err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		utils.Zlog.Fatalf("listen %s fail: %s", addr, err)
	} else {
		utils.Zlog.Info("socket listening", addr)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			utils.Zlog.Warning("listener.Accept occur error:", err)
			continue
		}

		go socketHandle(conn)
	}

}
