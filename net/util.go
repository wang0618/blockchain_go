package net

import (
	"bytes"
	"encoding/gob"
	"log"
)

func gobDecode(data []byte, e interface{}) {
	var buff bytes.Buffer
	buff.Write(data)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(e)
	if err != nil {
		log.Panic(err)
	}
}

func gobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
