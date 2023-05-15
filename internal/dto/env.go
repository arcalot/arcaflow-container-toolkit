package dto

import (
	"fmt"
	"os"

	"go.arcalot.io/log"
)

func LookupEnvVar(registries string, key string, logger log.Logger) (string, error) {
	val, ok := os.LookupEnv(key)
	if !ok {
		err := fmt.Errorf("%s environment variable not set to push to %s", key, registries)
		return "", err
	} else if len(val) == 0 {
		err := fmt.Errorf("%s environment variable empty to push to %s", key, registries)
		return "", err
	}
	return val, nil
}
