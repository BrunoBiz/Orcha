package logger

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Legacy diagnostic levels remain available to callers; both destinations now
// use LOG_LEVEL (info by default), rather than always writing verbose file logs.
const (
	LevelGIN  = slog.Level(-3)
	LevelFile = slog.Level(-2)
)

func LoadLogger() error {
	level, err := logLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		return err
	}
	//The 'Log' folder will always be in the root directory of Orcha
	err = checkDirectory()
	if err != nil {
		return err
	}

	// Will delete any file older than 15 days
	err = logRotation()
	if err != nil {
		return err
	}

	// Every log file will be named accordingly to the current date
	fileName := time.Now().Format(time.DateOnly)
	logFile, err := os.OpenFile(fmt.Sprintf("Log/%s-orcha.log", fileName), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		slog.Error("Can't create log file - " + err.Error())
		return err
	}

	multiHandler := slog.NewMultiHandler(
		slog.NewTextHandler(logFile, &slog.HandlerOptions{AddSource: level <= slog.LevelDebug, Level: level}),
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}),
	)
	slog.SetDefault(slog.New(multiHandler))

	return nil
}

func checkDirectory() error {
	cmd := exec.Command(`/bin/bash`, `-c`, `if [ ! -d Log ]; then mkdir 'Log' && echo 'Created'; else echo 'Already exists'; fi`) // TODO - There might be a better way to do this, but it works
	checkDir, err := cmd.CombinedOutput()

	slog.Debug("checkDirectory()")
	slog.Debug(cmd.String())
	slog.Debug(string(checkDir))

	// The return here (checkDir) does not matter, either it returns `Created` or  `Already exists`, both scenarios are acceptable, all that matters is that there was no errors
	if err != nil {
		slog.Error("Can't create log directory - " + err.Error())
		return err
	}

	return nil
}

func logRotation() error {
	// Will delete any log file older than 15 days
	cmd := exec.Command(
		"find",
		"./Log",
		"-type",
		"f",
		"-mtime",
		"+15",
		"-delete",
	)
	deleteOldLogs, err := cmd.CombinedOutput()

	slog.Debug("logRotation()")
	slog.Debug(cmd.String())
	slog.Debug(string(deleteOldLogs))

	if err != nil {
		slog.Error("Can't delete old logs - " + err.Error())
		return err
	}

	return nil
}

// LOG_LEVEL controls both destinations; verbose internal steps are opt-in.
func logLevel(value string) (slog.Level, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "info":
		return slog.LevelInfo, nil
	case "debug":
		return slog.LevelDebug, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, fmt.Errorf("invalid LOG_LEVEL %q: use debug, info, warn or error", value)
	}
}
