package logger

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"time"
)

/*
	Slog log levels

	File   -> Log levels -3 and above
	Stdout -> Log levels 0 and above

	const (
		LevelDebug Level = -4
		LevelGIN   Level = -3 --> Custom - All GIN middleware SLOG
		LevelFile  Level = -2 --> Custom - Any logs that need to be logged to the file, but not STDOUT
		LevelInfo  Level = 0
		LevelWarn  Level = 4
		LevelError Level = 8
	)
*/

const (
	LevelGIN  = slog.Level(-3)
	LevelFile = slog.Level(-2)
)

func LoadLogger() error {

	//The 'Log' folder will always be in the root directory of Orcha
	err := checkDirectory()
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

	// Removes the time and level information -> STDOUT
	replaceWithoutTimeLevel := func(groups []string, a slog.Attr) slog.Attr {
		if (a.Key == slog.TimeKey && len(groups) == 0) ||
			(a.Key == slog.LevelKey && len(groups) == 0) {
			return slog.Attr{}
		}

		if a.Key == slog.MessageKey && len(groups) == 0 {
			return slog.Attr{Key: "", Value: a.Value} // TODO - Still prints the "" in the stdout
		}

		return a
	}

	multiHandler := slog.NewMultiHandler(
		slog.NewTextHandler(logFile, &slog.HandlerOptions{AddSource: true, Level: LevelGIN}),                              // Logs to the log file, has added source
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: replaceWithoutTimeLevel, Level: slog.LevelInfo}), // Logs to the stdout, does not print the current time and log level
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
	cmd := exec.Command(`/bin/bash`, `-c`, `'find ./Log -type f -mtime +15 -exec rm {} ;'`)
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
