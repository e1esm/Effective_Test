package envParser

import (
	"fmt"
	"log"
	"os"
)

func ParseEnvVariable(keys ...string) (map[string]string, error) {

	envVars := make(map[string]string)

	for i := 0; i < len(keys); i++ {
		if val := os.Getenv(keys[i]); val != "" {
			envVars[keys[i]] = val
		}
	}

	if len(envVars) != len(keys) {
		log.Println(envVars)
		return nil, fmt.Errorf("some variables are missing: %v", keys)
	}

	return envVars, nil
}
