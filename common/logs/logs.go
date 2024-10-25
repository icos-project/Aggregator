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
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	// zerolog
	// https://blog.logrocket.com/5-structured-logging-packages-for-go/

	log_level := os.Getenv("LOG_LEVEL")

	// UNIX Time is faster and smaller than most timestamps
	//zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// global log level
	if len(log_level) == 0 {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		if strings.ToLower(log_level) == "trace" {
			zerolog.SetGlobalLevel(zerolog.TraceLevel)
		} else if strings.ToLower(log_level) == "debug" {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		} else {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}
	}

	// Console, not JSON
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006-01-02T15:04:05.999Z07:00"})
}

///////////////////////////////////////////////////////////////////////////////

func Println(m string) {
	log.Info().Msg(m)
}

func Printlne(m string, e error) {
	log.Error().Msg(m + " " + e.Error())
}

func Printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func Printfi(m string, s int) {
	log.Printf(m, s)
}

func Trace(args ...interface{}) {
	s := fmt.Sprint(args...)
	log.Trace().Msg(s)
}

func Debug(args ...interface{}) {
	s := fmt.Sprint(args...)
	log.Debug().Msg(s)
}

func Info(m string) {
	log.Info().Msg(m)
}

func Infof(args ...interface{}) {
	s := fmt.Sprint(args...)
	log.Info().Msg(s)
}

func Warn(args ...interface{}) {
	s := fmt.Sprint(args...)
	log.Warn().Msg(s)
}

func Error(args ...interface{}) {
	s := fmt.Sprint(args...)
	log.Error().Msg(s)
}

func Fatal(args ...interface{}) {
	s := fmt.Sprint(args...)
	log.Fatal().Msg(s)
}

func Panic(m string) {
	log.Panic().Msg(m)
}
