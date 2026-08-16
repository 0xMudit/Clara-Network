// Package framing implements the length-prefixed TCP framing used for
// ISO 8583 traffic: a 2-byte big-endian length header followed by the
// message payload.
package framing

import (
	"encoding/binary"
	"fmt"
	"io"
)

// ReadFrame reads one length-prefixed frame from r.
func ReadFrame(r io.Reader) ([]byte, error) {
	var hdr [2]byte
	if _, err := io.ReadFull(r, hdr[:]); err != nil {
		return nil, err
	}
	n := int(binary.BigEndian.Uint16(hdr[:]))
	if n == 0 {
		return nil, fmt.Errorf("framing: empty frame")
	}
	buf := make([]byte, n)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// WriteFrame writes a single length-prefixed frame to w.
func WriteFrame(w io.Writer, payload []byte) error {
	if len(payload) > 65535 {
		return fmt.Errorf("framing: payload too large (%d bytes)", len(payload))
	}
	var hdr [2]byte
	binary.BigEndian.PutUint16(hdr[:], uint16(len(payload)))
	if _, err := w.Write(hdr[:]); err != nil {
		return err
	}
	_, err := w.Write(payload)
	return err
}
