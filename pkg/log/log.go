package log

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger = nil

func GetLogger() (*zap.Logger, error) {

	if logger != nil {
		return logger, nil
	}

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, err := config.Build()
	if err != nil {
		fmt.Printf("Error creating logger: %v", err)
		return nil, err
	}

	return logger, nil

}
