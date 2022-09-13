package sets

import (
	"reflect"

	"github.com/segmentio/encoding/json"
)

type getKeyFunc[K comparable, V any] func(a V) K

func NewCustomSet[K comparable, V any](handler getKeyFunc[K, V]) Set[K, V] {
	return &Safe[K, V]{
		unsafe: Unsafe[K, V]{disableTypeCheck: true, getKey: handler},
	}
}

func NewLimitedCustomSet[K comparable, V any](handler getKeyFunc[K, V], limit int) Set[K, V] {
	return &Safe[K, V]{
		unsafe: Unsafe[K, V]{disableTypeCheck: true, getKey: handler, limit: limit},
	}
}

type Set[K comparable, V any] interface {
	PushStart(x V)
	PushEnd(x V)
	Size() int
	RemoveFirst() (x *V)
	RemoveLast() (x *V)
	Remove(x K)
	Exists(x K) bool
	Get(x K) *V
	First() *V
	Last() *V
	Each(cb func(a V) bool)
	Reset()
	UnmarshalJSON(data []byte) error
	MarshalJSON() ([]byte, error)
}

func checkType(x interface{}) (ok bool) {
	switch reflect.TypeOf(x).Kind() {
	case reflect.Struct:
		ok = false
	case reflect.Map:
		ok = false
	case reflect.Slice:
		ok = false
	case reflect.Interface:
		ok = false
	case reflect.Array:
		ok = false
	case reflect.Func:
		ok = false
	default:
		ok = true
	}
	return
}

type Unsafe[K comparable, V any] struct {
	keys             map[K]int
	getKey           getKeyFunc[K, V]
	data             []V
	limit            int
	disableTypeCheck bool
}

func (u *Unsafe[K, V]) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &u.data)
	if err != nil {
		return err
	}
	for i := range u.data {
		key := u.getKey(u.data[i])
		u.keys[key] = i
	}
	return nil
}

func (u *Unsafe[K, V]) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.data)
}

func (u *Unsafe[K, V]) Get(x K) *V {
	index, ok := u.keys[x]
	if !ok {
		return nil
	}
	val := u.data[index]
	return &val
}

func (u *Unsafe[K, V]) First() *V {
	if len(u.data) == 0 {
		return nil
	}
	val := u.data[0]
	return &val
}

func (u *Unsafe[K, V]) Last() *V {
	if len(u.data) == 0 {
		return nil
	}
	val := u.data[len(u.data)-1]
	return &val
}

func (u *Unsafe[K, V]) Each(cb func(a V) bool) {
	for i := range u.data {
		if !cb(u.data[i]) {
			break
		}
	}
}

func (u *Unsafe[K, V]) Exists(x K) bool {
	u.initKeysMap()
	_, ok := u.keys[x]
	return ok
}

func (u *Unsafe[K, V]) initKeysMap() {
	if u.keys == nil {
		u.keys = map[K]int{}
	}
	if u.getKey == nil {
		u.getKey = func(a V) K {
			x, ok := (interface{})(a).(K)
			if ok {
				return x
			}
			panic("types doesnt match")
		}
	}
}

func (u *Unsafe[K, V]) PushStart(x V) {
	if !checkType(x) && !u.disableTypeCheck {
		panic("Only primitive data types are supported!")
	}
	u.initKeysMap()
	key := u.getKey(x)
	if u.Exists(key) {
		return
	}
	if u.limit != 0 {
		if u.Size() == u.limit {
			u.RemoveLast()
		}
	}
	u.data = append([]V{x}, u.data...)
	for i := range u.keys {
		u.keys[i]++
	}
	u.keys[key] = 0
}

func (u *Unsafe[K, V]) PushEnd(x V) {
	if !checkType(x) && !u.disableTypeCheck {
		panic("Only primitive data types are supported!")
	}
	u.initKeysMap()
	key := u.getKey(x)
	if u.Exists(key) {
		return
	}
	if u.limit != 0 {
		if u.Size() == u.limit {
			u.RemoveLast()
		}
	}
	u.keys[key] = len(u.data)
	u.data = append(u.data, x)
}

func (u *Unsafe[K, V]) Size() int {
	return len(u.data)
}

func (u *Unsafe[K, V]) RemoveFirst() (x *V) {
	u.initKeysMap()
	x = u.First()
	if x == nil {
		return
	}
	key := u.getKey(*x)
	u.Remove(key)
	return
}

func (u *Unsafe[K, V]) RemoveLast() (x *V) {
	u.initKeysMap()
	x = u.Last()
	if x == nil {
		return
	}
	key := u.getKey(*x)
	u.Remove(key)
	return
}

func (u *Unsafe[K, V]) Remove(x K) {
	u.initKeysMap()
	key, ok := u.keys[x]
	if !ok {
		return
	}
	u.data = append(u.data[:key], u.data[key+1:]...)
	delete(u.keys, x)
	for i := range u.keys {
		if u.keys[i] > key {
			u.keys[i]--
		}
	}
}

func (u *Unsafe[K, V]) Reset() {
	u.initKeysMap()
	u.data = []V{}
	u.keys = map[K]int{}
}
