package conf

import "github.com/kataras/golog"

var (
	// Logger is a standard stdOut logging for printing require logs and can be use with any io.Writer
	Logger *golog.Logger
)

func init() {
	// Initialize the Logger with golog.new()
	Logger = golog.New()
	// Note: below functions are Logger customization and can be customize for your needs
	// Logger.SetOutput()
	// Logger.SetPrefix()
	// Logger.SetLevel()
	// Logger.SetTimeFormat()
}
