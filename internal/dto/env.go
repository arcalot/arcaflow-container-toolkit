package dto

import (
	"fmt"
	"go.arcalot.io/log"
	"os"
)

func LookupEnvVar(key string, logger log.Logger) verbose {
	val, ok := os.LookupEnv(key)
	var msg string
	if !ok {
		msg = fmt.Sprintf("%s not set", key)
	} else if len(val) == 0 {
		msg = fmt.Sprintf("%s is empty", key)
	}
	logger.Infof(msg)
	return verbose{return_value: val, msg: msg}
}

type verbose struct {
	msg          string
	return_value string
}
