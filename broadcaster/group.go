package broadcaster

import (
	"golang.org/x/exp/slices"
	"sync"
)

var ackPool = sync.Pool{
	New: func() interface{} {
		return make(chan bool)
	},
}

type Group[V any] struct {
	in         chan Message[V]
	close      chan bool
	m          *sync.Mutex
	members    []*Member[V]
	ackEnabled bool
}

func (v *Group[V]) loop() {
	for {
		select {
		case msg := <-v.in:
			go func() {
				v.m.Lock()
				for _, member := range v.members {
					member.in <- msg
				}
				v.m.Unlock()
			}()
		case <-v.close:
			v.m.Lock()
			members := v.members[:]
			v.m.Unlock()
			for _, member := range members {
				v.Leave(member)
			}
			return
		}
	}
}

// Close will close ALL channels and sending new messages will be no longer available
func (v *Group[V]) Close() {
	if v.close == nil {
		return
	}
	v.close <- true
	close(v.close)
	v.close = nil
}

// SendAsync sends value to all goroutines like Group.Send but won't block current thread
func (v *Group[V]) SendAsync(value V) {
	v.in <- Message[V]{
		data: value,
		ack:  nil,
	}
}

// Send sends value to the ALL goroutines listening this group and WAITING for acks from them
func (v *Group[V]) Send(value V) {
	v.m.Lock()
	wantedAcks := len(v.members)
	if wantedAcks == 0 {
		v.m.Unlock()
		return
	}
	v.m.Unlock()
	msg := Message[V]{
		data: value,
		ack:  ackPool.Get().(chan bool),
	}
	v.in <- msg
	acks := 0
	for {
		select {
		case <-msg.ack:
			acks++
			if acks >= wantedAcks {
				ackPool.Put(msg.ack)
				msg.ack = nil
				return
			}
		}
	}
}

// Join adds new goroutine to the group and listen for messages
func (v *Group[V]) Join() *Member[V] {
	m := &Member[V]{
		group:  v,
		Out:    make(chan Message[V]),
		in:     make(chan Message[V]),
		filter: nil,
		once:   false,
		close:  make(chan bool, 1),
	}
	go m.loop()
	v.m.Lock()
	v.members = append(v.members, m)
	defer v.m.Unlock()
	return m
}

// Leave removes specific group member and stop listening
func (v *Group[V]) Leave(m *Member[V]) {
	v.m.Lock()
	index := slices.Index[*Member[V]](v.members, m)
	if index == -1 {
		return
	}
	member := v.members[index]
	v.members[index] = nil
	v.members = append(v.members[:index], v.members[index+1:]...)
	v.m.Unlock()
	member.close <- true
}

func NewGroup[V any]() *Group[V] {
	g := &Group[V]{
		in:    make(chan Message[V]),
		close: make(chan bool),
		m:     new(sync.Mutex),
	}
	go g.loop()
	return g
}
