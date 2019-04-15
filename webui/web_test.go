package webui

import (
	"blockchain_go/log"
	"testing"
	"time"
)

func TestStartWebUIServer(t *testing.T) {

	go func() {
		i := 0
		for {
			log.MISC.Println(i)
			time.Sleep(time.Millisecond * 200)
		}
	}()

	StartWebUIServer("./html", "")
}
