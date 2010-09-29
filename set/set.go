package set

import "container/list"

type ElementValue interface {
	Equal(v interface{}) bool
}

type Set struct {
	l *list.List
}

func New() *Set { return &Set{list.New()} }
func (s *Set) Len() int { return s.l.Len() }

func (s *Set) wrapIter(c chan<- ElementValue) {
	for v := range s.l.Iter() {
		c <- v.(ElementValue)
	}
	close(c)
}

func (s *Set) Iter() <-chan ElementValue {
	c := make(chan ElementValue)
	go s.wrapIter(c)
	return c
}

func (s *Set) Contains(value ElementValue) bool {
	for e := s.l.Front(); e != nil; e = e.Next() {
		if value.Equal(e.Value) {
			return true
		}
	}
	return false
}

func (s *Set) Insert(value ElementValue) *list.Element {
	if s.Contains(value) {
		return nil
	}
	return s.l.PushFront(value)
}

func (s *Set) Remove(value ElementValue) bool {
	for e := s.l.Front(); e != nil; e = e.Next() {
		if value.Equal(e.Value) {
			s.l.Remove(e)
			return true
		}
	}

	return false;
}

func (s *Set) Subset(s2 *Set) bool {
	if s2.Len() > s.Len() {
		return false
	}

	for e := s2.l.Front(); e != nil; e = e.Next() {
		if !s.Contains(e.Value.(ElementValue)) {
			return false
		}
	}

	return true
}

func (s *Set) Superset(s2 *Set) bool {
	return s2.Subset(s)
}

func (s *Set) Equal(s2 *Set) bool {
	if s.Len() != s2.Len() {
		return false
	}

	return s.Subset(s2);
}

func (s *Set) Union(s2 *Set) *Set {
	u := New()

	for e := s.l.Front(); e != nil; e = e.Next() {
		u.l.PushFront(e.Value)
	}

	for e := s2.l.Front(); e != nil; e = e.Next() {
		u.Insert(e.Value.(ElementValue))
	}

	return u
}

func (s *Set) Intersection(s2 *Set) *Set {
	in := New()

	for e := s.l.Front(); e != nil; e = e.Next() {
		if s2.Contains(e.Value.(ElementValue)) {
			in.l.PushFront(e.Value)
		}
	}
	return in
}

func (s *Set) Difference(s2 *Set) *Set {
	df := New()

	for e := s.l.Front(); e != nil; e = e.Next() {
		if !s2.Contains(e.Value.(ElementValue)) {
			df.l.PushFront(e.Value)
		}
	}
	return df
}
