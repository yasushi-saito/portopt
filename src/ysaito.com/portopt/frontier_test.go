//
// Created by Yaz Saito on 06/14/12.
//
package portopt
/*
import "testing"
func testAssert(t *testing.T, b bool, messages... interface{}) {
	if !b {
		t.Fatal(messages...)
	}
}

func TestFrontier_SameX(t *testing.T) {
	f := newFrontier()
	testAssert(t, f.Insert(0.0, 1.0), "Insert1")
	testAssert(t, !f.Insert(0.0, 1.5), "Insert2")
	testAssert(t, f.Insert(0.0, 0.5), "Insert3")
	testAssert(t, !f.Insert(0.0, 0.5), "Insert4")
}

func TestFrontier_DifferentX(t *testing.T) {
	f := newFrontier()
	testAssert(t, f.Insert(0.0, 0.0), "Insert1")
	testAssert(t, f.Insert(1.0, 1.0), "Insert2")
	testAssert(t, f.Insert(2.0, 3.0), "Insert3")
	testAssert(t, !f.Insert(0.3, 0.31), "Insert4")
	testAssert(t, f.Insert(0.3, 0.29), "Insert5")
	testAssert(t, !f.Insert(1.1, 1.22), "Insert6")
	testAssert(t, f.Insert(1.1, 1.18), "Insert7")
}

func TestFrontier_RemoveExisting(t *testing.T) {
	f := newFrontier()
	testAssert(t, f.Insert(0.0, 0.0), "Insert1")
	testAssert(t, f.Insert(1.0, 1.0), "Insert2")
	testAssert(t, f.Insert(2.0, 2.0), "Insert3")
	testAssert(t, f.Insert(3.0, 3.0), "Insert3")

	// this should delete entries <1.0,1.0> and <2.0,2.0>
	testAssert(t, f.Insert(1.1, 0.1), "Insert4")

	testAssert(t, !f.Insert(1.0, 1.1), "Insert5")
	testAssert(t, !f.Insert(2.0, 2.1), "Insert6")
}
*/