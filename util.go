package anti_entropy

import (
	"encoding/binary"
	"io"
	"net"
)

func lpWrite(w io.Writer, b []byte) error {

	l := len(b)
	lp := make([]byte, 8)

	binary.BigEndian.PutUint64(lp, uint64(l))

	out := net.Buffers{lp, b}
	_, e := out.WriteTo(w)
	return e

}

func lpRead(r io.Reader) ([]byte, error) {
	lp := make([]byte, 8)
	_, e := io.ReadFull(r, lp)

	if e != nil {
		return lp, e
	}
	l := binary.BigEndian.Uint64(lp)
	out := make([]byte, l)
	_, e = io.ReadFull(r, out)
	return out, e
}
