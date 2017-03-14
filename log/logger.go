package log

import (
	"fmt"
	
    log15 "gopkg.in/inconshreveable/log15.v2"
)

func NewLogger() (log15.Logger, error) {
    logger := log15.New()
    SetFilterHandler("warn", logger, log15.StdoutHandler)
    return logger, nil
}

func SetFilterHandler(level string, logger log15.Logger, handler log15.Handler) error {
	if level == "none" {
		logger.SetHandler(log15.DiscardHandler())
		return nil
	}

	lvl, err := log15.LvlFromString(level)
	if err != nil {
		return fmt.Errorf("Bad log level: %v", err)
	}
	logger.SetHandler(log15.LvlFilterHandler(lvl, handler))

	return nil
}