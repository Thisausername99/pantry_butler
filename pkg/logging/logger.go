package logging

import (
	"sync"

	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	once   sync.Once
)

// GetLogger returns a singleton zap.Logger instance.
func GetLogger() *zap.Logger {
	once.Do(func() {
		var err error
		logger, err = zap.NewDevelopment() // Use zap.NewProduction() in production
		if err != nil {
			panic(err)
		}
	})
	return logger
}
