package pgtypes

import (
	"testing"
)

var data = [][]byte{
	[]byte("{}"),
	[]byte("{1,2,3}"),
	[]byte("{4,6,12,4,674,1,4545}"),
	[]byte("{34,6135564}"),
	[]byte("{5}"),
}

func BenchmarkOld(b *testing.B) {
	var s *IntArray
	for i := 0; i < b.N; i++ {
		s.Scan(data[2])
	}
}

/*
func BenchmarkNew(b *testing.B) {
	var s *IntArray
	for i := 0; i < b.N; i++ {
		s.Scan2(data[2])
	}
}
*/

func eq(a IntArray, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

var out = [][]int{
	{},
	{1, 2, 3},
	{4, 6, 12, 4, 674, 1, 4545},
	{34, 6135564},
	{5},
}

func TestInts(t *testing.T) {
	for k, test := range data {
		s := new(IntArray)
		err := s.Scan(test)
		if err != nil {
			t.Error(err)
			continue
		}
		if s == nil {
			t.Error("nil")
			continue
		}
		if !eq(*s, out[k]) {
			t.Errorf("'%s' should be '%v', not '%v'", string(test), out[k], *s)
		}
	}
}

var booldata = [][]byte{
	[]byte("{}"),
	[]byte("{NULL}"),
	[]byte("{f,t,t,t,f,NULL,t,t}"),
	[]byte("{t}"),
	[]byte("{NULL,NULL,NULL,NULL}"),
}

var boolout = [][]bool{
	{},
	{false},
	{false, true, true, true, false, false, true, true},
	{true},
	{false, false, false, false},
}

func booleq(a *BoolArray, b []bool) bool {
	if len(*a) != len(b) {
		return false
	}

	for i, A := range *a {
		if A != b[i] {
			return false
		}
	}
	return true
}

func TestBools(t *testing.T) {
	for k, test := range booldata {
		b := new(BoolArray)
		err := b.Scan(test)
		if err != nil {
			t.Error(err)
			continue
		}
		if !booleq(b, boolout[k]) {
			t.Errorf("'%v' != '%v'", b, boolout[k])
		}
	}
}
