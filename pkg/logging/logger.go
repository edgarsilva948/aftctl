/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

package logging

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

// getColorCode returns ANSI color code for a given color name
func getColorCode(colorName string) string {
	colorMap := map[string]string{
		"red":     "31",
		"green":   "32",
		"yellow":  "33",
		"blue":    "34",
		"magenta": "35",
		"cyan":    "36",
		"white":   "37",
	}

	if code, exists := colorMap[colorName]; exists {
		return code
	}

	return "37"
}

// CustomLog centralizes the custom logging function
func CustomLog(emoji string, colorName string, msg string) {

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.DisableCaller = true
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("02/01/2006 15:04:05"))
	})

	var err error
	logger, err = config.Build()
	if err != nil {
		fmt.Printf("Error initializing logger: %v", err)
	}

	colorCode := getColorCode(colorName)
	coloredMsg := fmt.Sprintf("\x1b[%sm%s %s\x1b[0m", colorCode, emoji, msg)
	logger.Info(coloredMsg)
}
