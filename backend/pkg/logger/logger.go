package logger

import (
	"log"

	"github.com/fatih/color"
)

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

func (lg *Logger) FatalError(errorText string) {
	log.Fatal(color.RedString("%s", errorText))
}

func (lg *Logger) SubmitSuccess() {
	log.Print(color.WhiteString("SUBMISSION: OK"))
}

func (lg *Logger) SubmitError(errText string) {
	log.Print(color.YellowString("SUMBISSION: ERROR: %s", errText))
}
