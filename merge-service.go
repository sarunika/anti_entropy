package anti_entropy

import (
	"net"
)

// The core handshake for everything its just a push pull exchange for data
func Xchg(c net.Conn, out []byte, in Serde) (Serde, error) {
	e := lpWrite(c, out)
	b, e := lpRead(c)

	if e != nil {
		return in, e
	}

	e = in.De(b)

	return in, e
}

func Merge(c net.Conn, s State, m Merger) (Serde, error) {

	var out []byte

	rf := func(d Serde) {
		out = d.Ser()
	}

	s.Read(rf)

	data, e := Xchg(c, out, m.Zero())

	var res Serde
	if e != nil {
		return data, e
	}

	wf := func(d Serde) Serde {
		res = m.Merge(d, data)
		return res
	}

	s.Write(wf)

	return res, nil
}
