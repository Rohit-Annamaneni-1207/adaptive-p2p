package tracker

import(
	"fmt"
	"crypto/rand"
	"net"
	"net/http"
	"net/url"
	"encoding/binary"
	"strconv"
	"github.com/jackpal/bencode-go"
	"github.com/Rohit-Annamaneni-1207/adaptive-p2p/parser"
)

type Peer struct{
	IPAddress net.IP
	Port uint16
}

type TrackerResp struct{
	Interval int
	Peers string
}

type TrackerReqParams struct{
	url string
	PeerId []byte
}

type p2pInfo struct{
	Peers []Peer
	PeerId []byte
}

func FormatReq (announce string, InfoHash [20]byte, InfoLength int) (*TrackerReqParams, error){
	const PORT = 6881
	base, err := url.Parse(announce)

	if err != nil{
		fmt.Println("Unable to parse url")
		return nil, err
	}

	const PeerIdLen = 20
	PeerIdBuf := make([]byte, PeerIdLen)
	n, err := rand.Read(PeerIdBuf)

	if err != nil{
		return nil, fmt.Errorf("generating peer id failed: %w", err)
	}

	if n != 20{
		return nil, fmt.Errorf("generate peer id: insufficient data")
	}

	params := url.Values{
		"info_hash": 	[]string{string(InfoHash[:])},
		"peer_id":		[]string{string(PeerIdBuf[:])},		
		"port":			[]string{strconv.Itoa(PORT)},
		"uploaded":  	[]string{"0"},
		"downloaded": 	[]string{"0"},
		"compact": 		[]string{"1"},
		"left":			[]string{strconv.Itoa(InfoLength)},
	}

	base.RawQuery = params.Encode()

	return &TrackerReqParams{
		url: base.String(),
		PeerId: PeerIdBuf,
	}, nil
}

func ParsePeers (peers []byte) ([]Peer, error){
	const peerlen = 6
	const portlen = 2

	PeerCount := len(peers)/peerlen

	if len(peers)%peerlen != 0{
		return nil, fmt.Errorf("received corrupted data")
	}

	PeerList := make([]Peer, PeerCount)

	for i:=0; i<PeerCount; i++{

		end := (i+1)*peerlen
		PeerList[i].IPAddress = net.IP(peers[i*peerlen:end-portlen])
		PeerList[i].Port = binary.BigEndian.Uint16([]byte(peers[end-portlen: end]))
	}

	return PeerList, nil
}

func TrackerGetReq (bct *parser.BencodeTorrent) (*p2pInfo, error){
	trp, err := FormatReq(bct.Announce, bct.InfoHash, bct.Info.Length)

	if err != nil{
		return nil, err
	}

	resp, err := http.Get(trp.url)

	if err != nil{
		return nil, err
	}

	defer resp.Body.Close()
	tr := TrackerResp{}

	err = bencode.Unmarshal(resp.Body, &tr)

	if err != nil{
		return nil, err
	}

	PeerList, err := ParsePeers([]byte(tr.Peers))

	if err != nil{
		return nil, err
	}

	return &p2pInfo{
		Peers: PeerList,
		PeerId: trp.PeerId,
	}, nil

}

func StringifyPeer (p Peer) (string){
	return net.JoinHostPort(p.IPAddress.String(), strconv.Itoa(int(p.Port)))
}

