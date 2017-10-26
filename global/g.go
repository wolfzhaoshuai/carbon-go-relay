package global

import "runtime"

// changelog:
// 0.0.1: init project

//Define some constants
const (
	VERSION = "0.0.1"
)

//define some global varibales
var (
	GlobalOriginReceiveItems int32
	GlobalReceiveItems       int32
	GlobalOriginSendItems    int32
	GlobalSendItems          int32
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
