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

func DeltaMerge(c *net.Conn, s State, ds DeltaState) (Serde, error) {
	zd := ds.ZeroDelta()
	out := make([]byte)

	rf := func(d Serde) {
		out = d.Ser()
	}

	s.Read(rf)

	meta, e := Xchg(c, out, zd)

	var diff Serde

	rf1 := func(d Serde) {
		diff = ds.Diff(d, meta)
	}

	if e != nil {
		return meta, e
	}
	s.Read(rf1)

	remote, e1 := Xchg(c, diff.Ser(), ds.Zero())

	complete := func(d Serde) Serde {
		return ds.Merge(d, remote)
	}

	if e1 != nil {
		return
	}

	result := s.Write(complete)

}
