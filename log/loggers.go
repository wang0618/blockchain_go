package log

import (
	"log"
	"os"
)

var (
	CLI   *log.Logger
	Net   *log.Logger
	Miner *log.Logger
	MISC  *log.Logger
)

var LogChan = make(chan string, 10)

type logWriter struct{}

func (_ logWriter) Write(p []byte) (n int, err error) {
	LogChan <- string(p)

	return len(p), nil
}

func init() {
	CLI = log.New(os.Stdout, "CLI: ", log.Lshortfile|log.Ltime)
	Net = log.New(os.Stdout, "Net: ", log.Lshortfile|log.Ltime)
	Miner = log.New(os.Stdout, "Miner: ", log.Lshortfile|log.Ltime)
	MISC = log.New(os.Stdout, "MISC: ", log.Lshortfile|log.Ltime)
}

func LogToChan() {
	CLI = log.New(logWriter{}, "CLI: ", log.Lshortfile|log.Ltime)
	Net = log.New(logWriter{}, "Net: ", log.Lshortfile|log.Ltime)
	Miner = log.New(logWriter{}, "Miner: ", log.Lshortfile|log.Ltime)
	MISC = log.New(logWriter{}, "MISC: ", log.Lshortfile|log.Ltime)
}
