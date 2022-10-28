package log

import (
	"errors"
	"fmt"
	"os"

	"github.com/jinzhu/copier"
)

var defaultLogger *logger

func Init(options ...Option) error {
	if defaultLogger != nil {
		return errors.New("logger has been initialized")
	}

	soup := defaultOption()

	for _, salt := range options {
		err := copier.CopyWithOption(soup, salt, copier.Option{IgnoreEmpty: true})
		if err != nil {
			return err
		}
	}

	// If output dir does not exist, we make it
	_, err := os.Stat(soup.OutputDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		err = os.Mkdir(soup.OutputDir, 0755)
		if err != nil {
			return err
		}
	}

	defaultLogger, err = newLogger(soup)

	return err
}

func Error(err error, message string) {
	defaultLogger.Error(err, message)
}

func Errorf(err error, format string, a ...any) {
	defaultLogger.Error(err, fmt.Sprintf(format, a...))
}

func Warn(message string) {
	defaultLogger.Warn(message)
}

func Warnf(format string, a ...any) {
	defaultLogger.Warn(fmt.Sprintf(format, a...))
}

func Info(message string) {
	defaultLogger.Info(message)
}

func Infof(format string, a ...any) {
	defaultLogger.Info(fmt.Sprintf(format, a...))
}

func Debug(message string) {
	defaultLogger.Debug(message)
}

func Debugf(format string, a ...any) {
	defaultLogger.Debug(fmt.Sprintf(format, a...))
}
