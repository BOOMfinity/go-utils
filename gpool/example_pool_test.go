package gpool

import "bytes"

func ExampleNew_simple() {
	// declare pool of ints
	pool := New[int]()
	// allocate new pointer or reuse if available
	ptr := pool.Get()

	/*
		do things with this pointer
	*/

	// put pointer back to pool
	pool.Put(ptr)
	// pointer can be reused. Do not use pointer after calling Put method!
}

func ExampleNew_buffers() {
	// declare pool of bytes.Buffer
	pool := New[bytes.Buffer](
		// OnInit function is triggered only if new pointer has been allocated
		OnInit[bytes.Buffer](func(b *bytes.Buffer) {
			// increase buffer size by 1024 bytes to limit future allocations
			b.Grow(1024)
		}),
		OnPut[bytes.Buffer](func(b *bytes.Buffer) {
			// reset buffer for future use, so you don't have to worry about it after calling Pool.Get
			b.Reset()
		}),
	)

	// simple example of using bytes.Buffer with pooling

	buff := pool.Get()
	buff.WriteString("Hello World!")
	pool.Put(buff)
}
