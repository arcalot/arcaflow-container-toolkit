package dto

import (
	"os"

	"go.arcalot.io/log"
)

func LookupEnvVar(registries string, key string, logger log.Logger) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		logger.Errorf("%s not set for %s", key, registries)
	} else if len(val) == 0 {
		logger.Errorf("%s empty for %s", key, registries)
	}
	return val
}
