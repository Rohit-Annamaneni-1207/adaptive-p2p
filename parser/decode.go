package parser

import (
	"os"
	"fmt"
    "github.com/jackpal/bencode-go"
)

type InfoDict struct {
	PieceLength int `bencode:"piece length"`
	Pieces 		string `bencode:"pieces"`
	Name 		string `bencode:"name"`
	Length 		int `bencode:"length"`
}

type TorrentFile struct {
	Announce 	string `bencode:"announce"`
	Info 		InfoDict `bencode:"info"`
}

func ParseTorrentFile(path string) (*TorrentFile, error){
	file, err := os.Open(path)
	if err != nil{
		return nil, err
	}
	defer file.Close()

	var tf TorrentFile
	err = bencode.Unmarshal(file, &tf)
	if err != nil{
		return nil, err
	}

	return &tf, nil
}

func PrintSummary(tf *TorrentFile) (){
	
	fmt.Println("Tracker URL:", tf.Announce)
	fmt.Println("File Name:", tf.Info.Name)
	fmt.Println("Piece Length:", tf.Info.PieceLength)
	fmt.Println("Total Length:", tf.Info.Length)
	fmt.Println("Number of Pieces:", len(tf.Info.Pieces)/20)
}