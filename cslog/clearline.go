package cslog

import (
	log "github.com/sirupsen/logrus"
)

// ClearLineFormatter is a logrus formatter that clears each line before printing
// used in cases where the log line is updated in place (e.g. progress bar)
//
// Make sure the logger is set to print on a terminal before using this or it will
// pollute the log with control characters.
type ClearLineFormatter struct {
    log.TextFormatter
}

func (f *ClearLineFormatter) Format(entry *log.Entry) ([]byte, error) {
	const clearLine = "\r\033[K"

	result, err := f.TextFormatter.Format(entry)
	if err != nil {
		return nil, err
	}
	return append([]byte(clearLine), result...), nil
}
