package utils

import "github.com/op/go-logging"

func createLogger() *logging.Logger {
	//set logger format globally
	var format = logging.MustStringFormatter(
		`[%{level:.4s}] [%{shortpkg} %{shortfile} %{shortfunc}] %{message}`,
	)
	logging.SetFormatter(format)
	return logging.MustGetLogger("")
}

//Zlog is a global logging instance
var Zlog = createLogger()
