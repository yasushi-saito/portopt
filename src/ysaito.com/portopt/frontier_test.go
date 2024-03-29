//
// Created by Yaz Saito on 06/14/12.
//
package portopt

import "testing"
import "math/rand"

func testAssert(t *testing.T, b bool, messages... interface{}) {
	if !b {
		t.Fatal(messages...)
	}
}

func TestFrontier_SameX(t *testing.T) {
	f := newFrontier()
	testAssert(t, f.Insert(0.0, 1.0, nil), "Insert1")
	testAssert(t, !f.Insert(0.0, 1.5, nil), "Insert2")
	testAssert(t, f.Insert(0.0, 0.5, nil), "Insert3")
	testAssert(t, !f.Insert(0.0, 0.5, nil), "Insert4")
}

func TestFrontier_DifferentX(t *testing.T) {
	f := newFrontier()
	testAssert(t, f.Insert(0.0, 0.0, nil), "Insert1")
	testAssert(t, f.Insert(1.0, 1.0, nil), "Insert2")
	testAssert(t, f.Insert(2.0, 3.0, nil), "Insert3")
	testAssert(t, !f.Insert(0.3, 0.31, nil), "Insert4")
	testAssert(t, f.Insert(0.3, 0.29, nil), "Insert5")
	testAssert(t, !f.Insert(1.1, 1.22, nil), "Insert6")
	testAssert(t, f.Insert(1.1, 1.18, nil), "Insert7")
}

func TestFrontier_RemoveExisting(t *testing.T) {
	f := newFrontier()
	testAssert(t, f.Insert(0.0, 0.0, nil), "Insert1")
	testAssert(t, f.Insert(1.0, 1.0, nil), "Insert2")
	testAssert(t, f.Insert(2.0, 2.0, nil), "Insert3")
	testAssert(t, f.Insert(3.0, 3.0, nil), "Insert3")

	// this should delete entries <1.0,1.0> and <2.0,2.0>
	testAssert(t, f.Insert(1.1, 0.1, nil), "Insert4")

	testAssert(t, !f.Insert(1.0, 1.1, nil), "Insert5")
	testAssert(t, !f.Insert(2.0, 2.1, nil), "Insert6")
}

func TestFrontier_Randomized(t *testing.T) {
	r := rand.New(rand.NewSource(0))
	f := newFrontier()
	for i := 0; i < 100; i++ {
		x := float64(r.Intn(10000)) / 10000.0
		y := float64(r.Intn(10000)) / 10000.0
		f.Insert(x, y, nil)

		lastX := -1.0
		lastY := -1.0
		for iter := f.Iterate(); !iter.Done(); iter = iter.Next() {
			testAssert(t, lastX <= iter.Mean(), f.String())
			testAssert(t, lastY <= iter.Stddev(), f.String())
			lastX, lastY = iter.Mean(), iter.Stddev()
		}
	}
}
