/*
Copyright © 2022-2024 EVIDEN

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package logs

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var lsugar *zap.SugaredLogger

// caller formatter
func funcCaller(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("\t[" + caller.TrimmedPath() + "]")
}

// init
func init() {

	// Custom configuration
	config := zap.Config{
		Encoding:    "console",                                // Output format (json or console)
		Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel), // Log level
		OutputPaths: []string{"stdout"},                       // Output destinations //"./logs/logfile.log"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:          "ts",                             // Key for the timestamp field
			LevelKey:         "level",                          // Key for the log level field
			NameKey:          "logger",                         // Key for the logger name field
			CallerKey:        "caller",                         // Key for the caller field
			MessageKey:       "msg",                            // Key for the message field
			StacktraceKey:    "stacktrace",                     // Key for the stacktrace field
			LineEnding:       zapcore.DefaultLineEnding,        // Line ending character
			EncodeLevel:      zapcore.CapitalColorLevelEncoder, // CapitalColorLevelEncoder, CapitalLevelEncoder, LowercaseLevelEncoder, // Log level format
			EncodeTime:       zapcore.ISO8601TimeEncoder,       // Timestamp format
			EncodeDuration:   zapcore.StringDurationEncoder,    // Duration format
			EncodeCaller:     funcCaller,                       // Caller format
			ConsoleSeparator: " ",
		},
	}

	// Build the logger with the custom configuration
	logger, _ := config.Build()
	lsugar = logger.Sugar()
	defer logger.Sync() // Flushes buffer, if any

	// Log a message with the custom logger
	lsugar.Info("Logger initialized with custom configuration")
}

// GetLogger Returns global and configured logger
func GetLogger() *zap.SugaredLogger {
	return lsugar
}
