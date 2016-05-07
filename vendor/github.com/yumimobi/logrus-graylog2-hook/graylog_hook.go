package graylog

import (
	"bytes"
	"os"
	"time"

	"github.com/Graylog2/go-gelf/gelf"
	"github.com/Sirupsen/logrus"
)

// Hook send logs to a logging service compatible with the Graylog API and the GELF format.
type Hook struct {
	hostname string
	facility string
	w        *gelf.Writer
	levels   []logrus.Level
}

// can be mocked out for testing
var (
	hostname = os.Hostname
	now      = time.Now
)

// New creates a graylog2 hook
func New(address, facility string, level logrus.Level) (logrus.Hook, error) {
	w, err := gelf.NewWriter(address)
	if err != nil {
		return nil, err
	}

	hostname, err := hostname()
	if err != nil {
		return nil, err
	}

	return &Hook{
		levels:   levelThreshold(level),
		w:        w,
		hostname: hostname,
		facility: facility,
	}, err
}

func levelThreshold(l logrus.Level) []logrus.Level {
	for i := range logrus.AllLevels {
		if logrus.AllLevels[i] == l {
			return logrus.AllLevels[:i+1]
		}
	}
	return logrus.AllLevels
}

// Levels implements logrus.Hook interface
func (h *Hook) Levels() []logrus.Level {
	return h.levels
}

// Fire implements logrus.Hook interface
func (h *Hook) Fire(entry *logrus.Entry) error {
	p := bytes.TrimSpace([]byte(entry.Message))
	short := bytes.NewBuffer(p)
	full := ""
	if i := bytes.IndexRune(p, '\n'); i > 0 {
		full = short.String()
		short.Truncate(i)
	}
	extra := map[string]interface{}{}
	for k, v := range entry.Data {
		extra["_"+k] = v // prefix with _ will be treated as an additional field
	}
	extra["_facility"] = h.facility

	m := &gelf.Message{
		Version:  "1.1",
		Host:     h.hostname,
		Short:    short.String(),
		Full:     full,
		TimeUnix: float64(now().UnixNano()) / 1e9,
		Level:    int32(entry.Level),
		Extra:    extra,
	}
	return h.w.WriteMessage(m)
}
