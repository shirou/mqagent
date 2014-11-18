package main

import (
	"bytes"
	"fmt"
	"sort"
	"time"

	"github.com/Sirupsen/logrus"
)

const (
	nocolor = 0
	red     = 31
	green   = 32
	yellow  = 33
	blue    = 34
)

var (
	isTerminal bool
)

type LTSVFormatter struct {
	ForceColors   bool
	DisableColors bool
}

func (f *LTSVFormatter) Format(entry *logrus.Entry) ([]byte, error) {

	var keys []string
	for k := range entry.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	b := [][]byte{
		f.appendKeyValue("time", entry.Time.Format(time.RFC3339)),
		f.appendKeyValue("level", entry.Level.String()),
		f.appendKeyValue("msg", entry.Message),
	}
	for _, key := range keys {
		b = append(b, f.appendKeyValue(key, entry.Data[key]))
	}

	out := bytes.Join(b, []byte("\t"))
	out = append(out, '\n')

	return out, nil
}

func (f *LTSVFormatter) appendKeyValue(key, value interface{}) []byte {
	b := &bytes.Buffer{}
	switch value.(type) {
	case string, error:
		fmt.Fprintf(b, "%s:%s", key, value)
	default:
		fmt.Fprintf(b, "%s:%v", key, value)
	}
	return b.Bytes()
}
