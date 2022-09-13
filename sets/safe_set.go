package sets

import "sync"

type Safe[K comparable, V any] struct {
	unsafe Unsafe[K, V]
	mutex  sync.RWMutex
}

func (s *Safe[K, V]) UnmarshalJSON(data []byte) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.unsafe.UnmarshalJSON(data)
}

func (s *Safe[K, V]) MarshalJSON() ([]byte, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.unsafe.MarshalJSON()
}

func (s *Safe[K, V]) Get(x K) *V {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.unsafe.Get(x)
}

func (s *Safe[K, V]) First() *V {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.unsafe.First()
}

func (s *Safe[K, V]) Last() *V {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.unsafe.Last()
}

func (s *Safe[K, V]) Each(cb func(a V) bool) {
	s.mutex.RLock()
	s.unsafe.Each(cb)
	s.mutex.RUnlock()
}

func (s *Safe[K, V]) PushStart(x V) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.unsafe.PushStart(x)
}

func (s *Safe[K, V]) PushEnd(x V) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.unsafe.PushEnd(x)
}

func (s *Safe[K, V]) Size() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.unsafe.Size()
}

func (s *Safe[K, V]) RemoveFirst() (x *V) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.unsafe.RemoveFirst()
}

func (s *Safe[K, V]) RemoveLast() (x *V) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.unsafe.RemoveLast()
}

func (s *Safe[K, V]) Remove(x K) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.unsafe.Remove(x)
}

func (s *Safe[K, V]) Exists(x K) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.unsafe.Exists(x)
}

func (s *Safe[K, V]) Reset() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.unsafe.Reset()
}
