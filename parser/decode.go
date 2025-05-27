package parser

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"os"

	"github.com/jackpal/bencode-go"
)

type InfoDict struct {
	PieceLength int    `bencode:"piece length"`
	Pieces      string `bencode:"pieces"`
	Name        string `bencode:"name"`
	Length      int    `bencode:"length"`
}

type BencodeTorrent struct {
	Announce     string `bencode:"announce"`
	CreatedBy    string `bencode:"created by"`
	CreationDate int64  `bencode:"creation date"`
	Encoding     string `bencode:"encoding"`
	Port         uint16 `bencode:"port"`

	Info        InfoDict `bencode:"info"`
	InfoHash    [20]byte
	PieceHashes [][20]byte
}

func InfoHash(bct *InfoDict) ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *bct)

	if err != nil {
		return [20]byte{}, nil
	}

	return sha1.Sum(buf.Bytes()), nil
}

func PieceHash(bct *InfoDict) ([][20]byte, error) {
	const HashLen = 20
	PiecesSlice := []byte(bct.Pieces)
	SliceLength := len(PiecesSlice)
	HashCount := SliceLength / HashLen
	PiecesList := make([][20]byte, HashCount)

	if SliceLength%HashLen != 0 {
		return nil, fmt.Errorf("malformed piece length received")
	}

	for i := range PiecesList {
		copy(PiecesList[i][:], PiecesSlice[i*HashLen:(i+1)*HashLen])
	}

	return PiecesList, nil

}

func ParseTorrentFile(path string) (*BencodeTorrent, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var tf BencodeTorrent
	err = bencode.Unmarshal(file, &tf)
	if err != nil {
		return nil, err
	}

	return &tf, nil
}

func PrintSummary(tf *BencodeTorrent) {

	fmt.Println("Tracker URL:", tf.Announce)
	fmt.Println("File Name:", tf.Info.Name)
	fmt.Println("Piece Length:", tf.Info.PieceLength)
	fmt.Println("Total Length:", tf.Info.Length)
	fmt.Println("Number of Pieces:", len(tf.Info.Pieces)/20)
}
