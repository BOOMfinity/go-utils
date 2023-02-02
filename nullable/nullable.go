package nullable

import (
	"bytes"
	"github.com/segmentio/encoding/json"
)

type Nullable[V any] struct {
	Value V
	isSet bool
}

func (v *Nullable[V]) Null() bool {
	return !v.isSet
}

func (v *Nullable[V]) MarshalJSON() (data []byte, err error) {
	if v.Null() {
		return []byte("null"), nil
	}
	raw, err := json.Marshal(v.Value)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func (v *Nullable[V]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) || len(data) == 0 {
		v.isSet = false
		return nil
	}
	err := json.Unmarshal(data, &v.Value)
	if err != nil {
		return err
	}
	v.isSet = true
	return nil
}

func (v *Nullable[V]) Set(value V) {
	v.Value = value
	v.isSet = true
}

func (v *Nullable[V]) Clear() {
	v.isSet = false
}

func New[V any]() *Nullable[V] {
	return &Nullable[V]{}
}
