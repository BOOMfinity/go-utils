package nullable

import "github.com/segmentio/encoding/json"

type Nullable[V any] struct {
	data  V
	isSet bool
}

func (v Nullable[V]) MarshalJSON() (data []byte, err error) {
	if !v.isSet {
		return []byte("null"), nil
	}
	raw, err := json.Marshal(v.data)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func (v *Nullable[V]) Set(value V) {
	v.data = value
	v.isSet = true
}

func (v *Nullable[V]) Clear() {
	v.isSet = false
}

func (v *Nullable[V]) Value() *V {
	if !v.isSet {
		return nil
	}
	return &v.data
}

func New[V any]() *Nullable[V] {
	return &Nullable[V]{}
}
