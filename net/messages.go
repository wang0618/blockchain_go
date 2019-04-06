package net

type version struct {
	Version    int
	BestHeight int
}

type verack struct{}

type addr struct {
	AddrList []string
}

type block struct {
	Block []byte
}

type getblocks struct {
}

type getdata struct {
	Type string
	ID   []byte
}

type inv struct {
	Type  string
	Items [][]byte
}

type tx struct {
	Transaction []byte
}
