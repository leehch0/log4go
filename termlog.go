// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package log4go

import (
	"fmt"
	"io"
	"os"
)

var stdout io.Writer = os.Stdout

var consoleLogQuit = make(chan struct{})

// This is the standard writer that prints to standard output.
type ConsoleLogWriter chan *LogRecord

// This creates a new ConsoleLogWriter
func NewConsoleLogWriter() ConsoleLogWriter {
	records := make(ConsoleLogWriter, LogBufferLength)
	go records.run(stdout)
	return records
}

func (w ConsoleLogWriter) run(out io.Writer) {
	var timestr string
	var timestrAt int64

	for rec := range w {
		if at := rec.Created.UnixNano() / 1e9; at != timestrAt {
			timestr, timestrAt = rec.Created.Format("01/02/06 15:04:05"), at
		}
		fmt.Fprint(out, "[", timestr, "] [", levelStrings[rec.Level], "] (", rec.Source, ") ", rec.Message, "\n")
	}

	// inflight logs flushed, safe to quit
	close(consoleLogQuit)
}

// This is the ConsoleLogWriter's output method.  This will block if the output
// buffer is full.
func (w ConsoleLogWriter) LogWrite(rec *LogRecord) {
	w <- rec
}

// Close stops the logger from sending messages to standard output.  Attempts to
// send log messages to this logger after a Close have undefined behavior.
func (w ConsoleLogWriter) Close() {
	close(w)

	// wait for inflight logs flush
	<-consoleLogQuit
}
