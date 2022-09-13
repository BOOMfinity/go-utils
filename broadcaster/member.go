package broadcaster

type Member[V any] struct {
	group  *Group[V]
	Out    chan Message[V]
	in     chan Message[V]
	filter func(msg Message[V]) bool
	close  chan bool
	once   bool
}

// Close will leave current Group and close communication
func (v *Member[V]) Close() {
	v.group.Leave(v)
}

func (v *Member[V]) loop() {
	for {
		select {
		case msg := <-v.in:
			if msg.sender == v {
				break
			}
			if v.filter != nil {
				if v.filter(msg) {
					v.Out <- msg
				} else {
					msg.ACK()
				}
			} else {
				v.Out <- msg
			}
			if v.once {
				v.Close()
			}
		case <-v.close:
			close(v.Out)
			return
		}
	}
}

// Recv reads one message from channel
//
// The "more" variable is false if channel has been closed and attempts for reading more messages will cause panic.
// If "more" is false just exit your read loop
func (v *Member[V]) Recv() (msg Message[V], more bool) {
	msg, more = <-v.Out
	return
}

// WithFilter is built-in message filtering feature.
// Messages which doesn't meet filter function will be automatically "acked"
//
// # Only valid messages can be read
//
// Default filter is set to nil which means ALL messages are forwarded for reading
func (v *Member[V]) WithFilter(f func(msg Message[V]) bool) {
	v.filter = f
}

// SetOnce will automatically call Member.Close after reading exactly ONE VALID message
func (v *Member[V]) SetOnce() {
	v.once = true
}
