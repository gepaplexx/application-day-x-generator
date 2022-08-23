package util

import (
	"sync"
)

var lock = &sync.Mutex{}

type configuration struct {
	debugSealedSecrets bool
}

func (s *configuration) GetDebugSealedSecrets() bool {
	return s.debugSealedSecrets
}

func (s *configuration) SetDebugSealedSecrets(debugSealedSecrets bool) {
	s.debugSealedSecrets = debugSealedSecrets
}

var config *configuration

func GetConfig() *configuration {
	if config == nil {
		lock.Lock()
		defer lock.Unlock()
		if config == nil {
			config = &configuration{
				debugSealedSecrets: false,
			}
		}
	}

	return config
}
