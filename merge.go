package anti_entropy

import "sync"

type Serde interface {
	Ser() []byte
	De([]byte) error
}

// Similar to state based crdts or a monoid states can be merged and there is an identity
type Merger interface {
	Merge(a Serde, b Serde) Serde
	Zero() Serde
}

// Similar to delta crdts  or how git works
// There's smaller metadata that can be diffed against the state that can be diffed against
type DeltaState interface {
	Merger

	// Gets Metadata to Diff against
	GetDelta(s Serde) Serde
	Diff(s Serde, d Serde) Serde
	ZeroDelta() Serde
}

type State struct {
	lock sync.RWMutex
	data Serde
}

func (s *State) Write(f func(d Serde) Serde) {
	defer s.lock.Unlock()
	s.lock.Lock()
	s.data = f(s.data)
}
func (s *State) Read(f func(d Serde)) {
	defer s.lock.RUnlock()
	s.lock.RLock()
	f(s.data)
}

func NewState(d Serde) State {
	var mu sync.RWMutex
	return State{mu, d}

}
