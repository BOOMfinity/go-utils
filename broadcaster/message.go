package broadcaster

type Message[V any] struct {
	sender *Member[V]
	data   V
	ack    chan bool
}

// Data stores value sent by sender
func (v *Message[V]) Data() V {
	return v.data
}

// ACK will tell sender that message is free to be reused if needed
//
// Use only for SYNC communication. No action if message is ASYNC
//
// IMPORTANT: DON'T FORGET TO CALL THIS FUNCTION IF USING Group.Send. THIS WILL BLOCK THREAD FOREVER.
func (v *Message[V]) ACK() {
	if v.ack == nil {
		return
	}
	v.ack <- true
}
