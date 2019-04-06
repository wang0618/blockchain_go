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

func init() {
	CLI = log.New(os.Stdout, "CLI: ", log.Lshortfile|log.Ltime)
	Net = log.New(os.Stdout, "Net: ", log.Lshortfile|log.Ltime)
	Miner = log.New(os.Stdout, "Miner: ", log.Lshortfile|log.Ltime)
	MISC = log.New(os.Stdout, "MISC: ", log.Lshortfile|log.Ltime)
}
