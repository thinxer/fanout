package fanout

import "sync"

// A non-blocking fan-out mechanism.
type Fan struct {
	clients []chan<- interface{}
	mu      sync.Mutex
}

// Receive the value with a buffer. If the write function returns an error,
// the receive loop will stop, and the error will be returned to the caller.
func (v *Fan) Receive(buf int, write func(v interface{}) error) error {
	ch := make(chan interface{}, buf)
	v.mu.Lock()
	v.clients = append(v.clients, ch)
	v.mu.Unlock()
	defer func() {
		v.mu.Lock()
		for i := range v.clients {
			if v.clients[i] == ch {
				v.clients = append(v.clients[:i], v.clients[i+1:]...)
				break
			}
		}
		v.mu.Unlock()
	}()
	for v := range ch {
		if err := write(v); err != nil {
			return err
		}
	}
	return nil
}

// Send a value for fanout.
func (f *Fan) Send(value interface{}) {
	f.mu.Lock()
	for _, c := range f.clients {
		select {
		case c <- value:
		default:
		}
	}
	f.mu.Unlock()
}

// Finish closes all clients.
func (f *Fan) Finish() {
	f.mu.Lock()
	for _, ch := range f.clients {
		close(ch)
	}
	f.clients = nil
	f.mu.Unlock()
}
