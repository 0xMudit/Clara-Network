package framing

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	cases := [][]byte{
		{0x01},
		{0x01, 0x02, 0x03},
		bytes.Repeat([]byte{0xAB}, 256),
		bytes.Repeat([]byte{0xFF}, 1000),
	}
	for _, payload := range cases {
		var buf bytes.Buffer
		if err := WriteFrame(&buf, payload); err != nil {
			t.Fatalf("write %d bytes: %v", len(payload), err)
		}
		got, err := ReadFrame(&buf)
		if err != nil {
			t.Fatalf("read %d bytes: %v", len(payload), err)
		}
		if !bytes.Equal(got, payload) {
			t.Fatalf("round-trip %d bytes: got %x, want %x", len(payload), got, payload)
		}
	}
}

func TestMultipleFrames(t *testing.T) {
	var buf bytes.Buffer
	frames := [][]byte{
		[]byte("hello"),
		[]byte("world"),
		{0x00},
		bytes.Repeat([]byte{0x42}, 500),
	}
	for _, f := range frames {
		if err := WriteFrame(&buf, f); err != nil {
			t.Fatal(err)
		}
	}
	for i, want := range frames {
		got, err := ReadFrame(&buf)
		if err != nil {
			t.Fatalf("read frame %d: %v", i, err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("frame %d: got %x, want %x", i, got, want)
		}
	}
}

func TestEmptyFrameRejected(t *testing.T) {
	var buf bytes.Buffer
	binary.BigEndian.PutUint16([]byte{0, 0}, 0)
	buf.Write([]byte{0x00, 0x00})
	_, err := ReadFrame(&buf)
	if err == nil {
		t.Fatal("expected error for zero-length frame")
	}
}

func TestWriteFrameTooLarge(t *testing.T) {
	var buf bytes.Buffer
	payload := make([]byte, 65536)
	err := WriteFrame(&buf, payload)
	if err == nil {
		t.Fatal("expected error for payload > 65535 bytes")
	}
}

func TestReadFrameEOF(t *testing.T) {
	_, err := ReadFrame(bytes.NewReader(nil))
	if err != io.EOF && err != io.ErrUnexpectedEOF {
		t.Fatalf("expected EOF, got %v", err)
	}
}

func TestReadFrameTruncatedPayload(t *testing.T) {
	var buf bytes.Buffer
	binary.BigEndian.PutUint16([]byte{0, 0}, 100)
	buf.Write([]byte{0x00, 0x64})
	buf.Write([]byte{0x01, 0x02, 0x03})
	_, err := ReadFrame(&buf)
	if err == nil {
		t.Fatal("expected error for truncated payload")
	}
}

func TestMaxSizePayload(t *testing.T) {
	payload := bytes.Repeat([]byte{0x55}, 65535)
	var buf bytes.Buffer
	if err := WriteFrame(&buf, payload); err != nil {
		t.Fatal(err)
	}
	got, err := ReadFrame(&buf)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, payload) {
		t.Fatal("max-size payload mismatch")
	}
}

func TestSingleBytePayload(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteFrame(&buf, []byte{0x42}); err != nil {
		t.Fatal(err)
	}
	got, err := ReadFrame(&buf)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != 0x42 {
		t.Fatalf("got %x, want 42", got)
	}
}
