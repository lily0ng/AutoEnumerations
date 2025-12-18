package logger

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fatih/color"
)

var (
	verbose    bool
	infoLog    *log.Logger
	errorLog   *log.Logger
	warningLog *log.Logger
	successLog *log.Logger
)

func Init(v bool) {
	verbose = v

	infoLog = log.New(os.Stdout, "", 0)
	errorLog = log.New(os.Stderr, "", 0)
	warningLog = log.New(os.Stdout, "", 0)
	successLog = log.New(os.Stdout, "", 0)
}

func Info(format string, args ...interface{}) {
	if infoLog == nil {
		Init(false)
	}

	blue := color.New(color.FgBlue, color.Bold)
	timestamp := time.Now().Format("15:04:05")
	msg := fmt.Sprintf(format, args...)
	infoLog.Printf("%s [%s] %s", blue.Sprint("ℹ"), timestamp, msg)
}

func Error(format string, args ...interface{}) {
	if errorLog == nil {
		Init(false)
	}

	red := color.New(color.FgRed, color.Bold)
	timestamp := time.Now().Format("15:04:05")
	msg := fmt.Sprintf(format, args...)
	errorLog.Printf("%s [%s] %s", red.Sprint("✗"), timestamp, msg)
}

func Warning(format string, args ...interface{}) {
	if warningLog == nil {
		Init(false)
	}

	yellow := color.New(color.FgYellow, color.Bold)
	timestamp := time.Now().Format("15:04:05")
	msg := fmt.Sprintf(format, args...)
	warningLog.Printf("%s [%s] %s", yellow.Sprint("⚠"), timestamp, msg)
}

func Success(format string, args ...interface{}) {
	if successLog == nil {
		Init(false)
	}

	green := color.New(color.FgGreen, color.Bold)
	timestamp := time.Now().Format("15:04:05")
	msg := fmt.Sprintf(format, args...)
	successLog.Printf("%s [%s] %s", green.Sprint("✓"), timestamp, msg)
}

func Debug(format string, args ...interface{}) {
	if !verbose {
		return
	}

	if infoLog == nil {
		Init(false)
	}

	cyan := color.New(color.FgCyan)
	timestamp := time.Now().Format("15:04:05")
	msg := fmt.Sprintf(format, args...)
	infoLog.Printf("%s [%s] %s", cyan.Sprint("⚙"), timestamp, msg)
}

func Progress(current, total int, message string) {
	if infoLog == nil {
		Init(false)
	}

	percentage := float64(current) / float64(total) * 100
	cyan := color.New(color.FgCyan)
	fmt.Printf("\r%s [%.1f%%] %s", cyan.Sprint("▶"), percentage, message)

	if current == total {
		fmt.Println()
	}
}
