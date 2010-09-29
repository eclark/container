// Copyright 2010 Eric Clark. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package set

import "testing"

type EqInt int

func (i EqInt) Equal(j interface{}) bool {
	if k, ok := j.(EqInt); ok {
		return i == k
	}
	return false
}

func TestNew(t *testing.T) {
	s := New()
	if s == nil {
		t.Fatal("New returned nil")
	}

	if s.Len() != 0 {
		t.Fatal("New set has nonzero length")
	}
}

func TestInsertRemove(t *testing.T) {
	s := New()

	if s.Insert(EqInt(3)) == nil {
		t.Error("Failed to insert 3")
	}
	if s.Insert(EqInt(4)) == nil {
		t.Error("Failed to insert 4")
	}
	if s.Insert(EqInt(5)) == nil {
		t.Error("Failed to insert 5")
	}
	if s.Insert(EqInt(5)) != nil {
		t.Error("Insert returned value on duplicate")
	}

	if !s.Remove(EqInt(3)) {
		t.Error("Remove failed")
	}
	if s.Contains(EqInt(3)) {
		t.Error("Set still contains 3")
	}
}

func TestCompare(t *testing.T) {
	x := New()

	x.Insert(EqInt(4))
	x.Insert(EqInt(5))
	x.Insert(EqInt(6))

	y := New()

	y.Insert(EqInt(2))
	y.Insert(EqInt(3))
	y.Insert(EqInt(4))
	y.Insert(EqInt(5))

	u := x.Union(y)
	if u.Len() != 5 {
		t.Error("Union length incorrect")
	}

	i := x.Intersection(y)
	if i.Len() != 2 {
		t.Error("Intersection length incorrect")
	}

	d := x.RelativeComplement(y)
	if d.Len() != 1 {
		t.Error("x diff y Difference length incorrect")
	}

	d = x.Complement(y)
	if d.Len() != 2 {
		t.Error("y diff x Difference length incorrect")
	}
}

func TestSetofSet(t *testing.T) {
	a := New()
	a.Insert(EqInt(1))
	a.Insert(EqInt(2))

	b := New()
	b.Insert(EqInt(1))
	b.Insert(EqInt(2))

	c := New()
	c.Insert(EqInt(2))
	c.Insert(EqInt(3))

	d := New()
	d.Insert(a)
	if d.Insert(b) != nil {
		t.Error("Duplicate set inserted")
	}
	d.Insert(c)

	for v := range d.Iter() {
		for w := range v.(*Set).Iter() {
			t.Log(v, w)
		}
	}
}
