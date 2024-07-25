package bencode

import "io"

type bInfo struct {
	Name         string `bencode:"name"`
	Pieces       string `bencode:"pieces"`
	PiecesLength int    `bencode:"pieces length"`
	Length       int    `bencode:"length"`
}

type bTorrent struct {
	Annouce string `bencode:"announce"`
	Info    bInfo  `bencode:"info"`
}

// Open and parses .torrent file
func Open(r io.Reader) (*bTorrent, error) {
	bto := bTorrent{}
	// Unmarshal algorithm

	return &bto, nil
}
