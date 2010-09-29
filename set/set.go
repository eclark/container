// Copyright 2010 Eric Clark. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Implements an unordered list of unique elements.
package set

import "container/list"

/*
	A type that satisfies set.ElementValue can be inserted into a Set.  The Equal method should test both the type and equality.

	Here is a simple example of a type that conforms to ElementValue:

		type EqInt int
		func (i EqInt) Equal(j interface{}) bool {
			if k, ok := j.(EqInt); ok {
				return i == k
			}
			return false
		}
*/
type ElementValue interface {
	Equal(v interface{}) bool
}

// Set represents an unordered list of unique elements.  The zero value for Set is not usable until Init() is called. 
type Set struct {
	l *list.List
}

// New returns an initialized Set.
func New() *Set { return &Set{list.New()} }

// Init initializes or clears a Set.
func (s *Set) Init() *Set { s.l.Init(); return s }

// Len returns the number of elements in the Set.
func (s *Set) Len() int { return s.l.Len() }

func (s *Set) wrapIter(c chan<- ElementValue) {
	for v := range s.l.Iter() {
		c <- v.(ElementValue)
	}
	close(c)
}

// Iter returns a channel of values in the Set.
func (s *Set) Iter() <-chan ElementValue {
	c := make(chan ElementValue)
	go s.wrapIter(c)
	return c
}

// Contains returns a boolean indicating if a value is part of the Set.
func (s *Set) Contains(value ElementValue) bool {
	for e := s.l.Front(); e != nil; e = e.Next() {
		if value.Equal(e.Value) {
			return true
		}
	}
	return false
}

// Insert adds a value in the Set.
func (s *Set) Insert(value ElementValue) *list.Element {
	if s.Contains(value) {
		return nil
	}
	return s.l.PushFront(value)
}

// Remove removes a value from the Set.
func (s *Set) Remove(value ElementValue) bool {
	for e := s.l.Front(); e != nil; e = e.Next() {
		if value.Equal(e.Value) {
			s.l.Remove(e)
			return true
		}
	}

	return false;
}

// Subset returns true if all values in the second set are also in the first.
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

// Superset returns true if all values in the first set are also in the second.
func (s *Set) Superset(s2 *Set) bool {
	return s2.Subset(s)
}

// Equal returns if the all values in the first set are also in the second, and the second contains no additional values.
func (s *Set) Equal(s2 *Set) bool {
	if s.Len() != s2.Len() {
		return false
	}

	return s.Subset(s2);
}

// Union returns a new set which contains the distinct values from both sets.
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

// Intersection returns a new set which contains values which only exist in both sets.
func (s *Set) Intersection(s2 *Set) *Set {
	in := New()

	for e := s.l.Front(); e != nil; e = e.Next() {
		if s2.Contains(e.Value.(ElementValue)) {
			in.l.PushFront(e.Value)
		}
	}
	return in
}

// Difference returns a new set which contains values from the first set which do not exist in the second.
func (s *Set) Difference(s2 *Set) *Set {
	df := New()

	for e := s.l.Front(); e != nil; e = e.Next() {
		if !s2.Contains(e.Value.(ElementValue)) {
			df.l.PushFront(e.Value)
		}
	}
	return df
}
