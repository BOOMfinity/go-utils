package gpool

type options[V any] struct {
	onInit  func(*V)
	onPut   func(*V)
	discard func(*V) bool
}

type Option[V any] func(opts *options[V])

func OnInit[V any](fn func(*V)) Option[V] {
	return func(opts *options[V]) {
		opts.onInit = fn
	}
}

func OnPut[V any](fn func(*V)) Option[V] {
	return func(opts *options[V]) {
		opts.onPut = fn
	}
}

func Discard[V any](fn func(*V) bool) Option[V] {
	return func(opts *options[V]) {
		opts.discard = fn
	}
}

func parseOpts[V any](opts []Option[V]) *options[V] {
	o := new(options[V])
	for i := range opts {
		opts[i](o)
	}
	return o
}
