package parser

import (
	"bytes"
	"os"
	"testing"

	"github.com/jackpal/bencode-go"
)

func TestParseTorrentFile(t *testing.T) {
	// Create a sample BencodeTorrent struct for testing
	info := InfoDict{
		PieceLength: 16384,
		Pieces:      string(bytes.Repeat([]byte{'a'}, 20)),
		Name:        "testfile",
		Length:      12345,
	}
	torrent := BencodeTorrent{
		Announce: "http://tracker",
		Info:     info,
	}

	var buf bytes.Buffer
	err := bencode.Marshal(&buf, torrent)
	if err != nil {
		t.Fatalf("Failed to marshal torrent: %v", err)
	}

	// Print the bencoded data for inspection
	t.Logf("Bencoded data: %q", buf.Bytes())

	tmpfile, err := os.CreateTemp("", "test.torrent")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	if _, err := tmpfile.Write(buf.Bytes()); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpfile.Close()

	tf, err := ParseTorrentFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("ParseTorrentFile failed: %v", err)
	}

	// Print all parsed info for inspection
	t.Logf("Parsed Announce: %s", tf.Announce)
	t.Logf("Parsed Info.Name: %s", tf.Info.Name)
	t.Logf("Parsed Info.Length: %d", tf.Info.Length)
	t.Logf("Parsed Info.PieceLength: %d", tf.Info.PieceLength)
	t.Logf("Parsed Info.Pieces: %q", tf.Info.Pieces)

	if tf.Info.Name != "testfile" {
		t.Errorf("Expected file name 'testfile', got '%s'", tf.Info.Name)
	}
	if tf.Info.Length != 12345 {
		t.Errorf("Expected length 12345, got %d", tf.Info.Length)
	}
	if tf.Announce != "http://tracker" {
		t.Errorf("Expected announce 'http://tracker', got '%s'", tf.Announce)
	}
}
