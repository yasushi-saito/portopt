//
// Created by Yaz Saito on 06/10/12.
//

package rbtree

import "testing"
import "math/rand"
import "log"
import "sort"

type testItem int

func compareItems(i1, i2 Item) int {
	return int(i1.(testItem)) - int(i2.(testItem))
}

func TestEmpty(t *testing.T) {
	tree := NewTree(compareItems)
	if tree.Len() != 0 {
		t.Error("Len!=0")
	}
	iter := tree.Find(testItem(10))
	if !iter.Done() {
		t.Fail()
	}
}

func TestBasic(t *testing.T) {
	tree := NewTree(compareItems)
	if !tree.Insert(testItem(10)) {
		t.Fail()
	}
	if tree.Insert(testItem(10)) {
		t.Fail()
	}

	if tree.Len() != 1 {
		t.Error("Len!=1")
	}
	iter := tree.Find(testItem(10))
	if iter.Done() {
		t.Fail()
	}
	if iter.Item().(testItem) != 10 {
		t.Error("Wrong item: ", iter.Item())
	}
	iter = tree.Find(testItem(11))
	if !iter.Done() {
		t.Fail()
	}

	iter = tree.Find(testItem(9))
	if iter.Done() {
		t.Fail()
	}

	if iter.Item().(testItem) != 10 {
		t.Error("Wrong item: ", iter.Item())
	}
	log.Print("done")
}

func TestDelete(t *testing.T) {
	tree := NewTree(compareItems)
	if tree.DeleteWithKey(testItem(10)) {
		t.Fail()
	}
	if tree.Len() != 0 {
		t.Fail()
	}

	if !tree.Insert(testItem(10)) {
		t.Fail()
	}
	if !tree.DeleteWithKey(testItem(10)) {
		t.Fail()
	}
	if tree.Len() != 0 {
		t.Fail()
	}
}

type oracle struct {
	data []testItem
}

func newOracle() *oracle {
	return &oracle{data: make([]testItem, 0)}
}

func (o *oracle) Len() int {
	return len(o.data)
}

func (o *oracle) Less(i, j int) bool {
	return o.data[i] < o.data[j]
}

func (o *oracle) Swap(i, j int) {
	e := o.data[j]
	o.data[j] = o.data[i]
	o.data[i] = e
}

func (o *oracle) Insert(key testItem) bool {
	for _, e := range o.data {
		if e == key {
			return false
		}
	}

	n := len(o.data) + 1
	newData := make([]testItem, n)
	copy(newData, o.data)
	o.data = newData
	o.data[n-1] = key
	sort.Sort(o)
	return true
}

func (o *oracle) RandomExistingKey(rand *rand.Rand) testItem {
	index := rand.Intn(len(o.data))
	return o.data[index]
}

func (o *oracle) Find(t *testing.T, key testItem) oracleIterator {
	prev := testItem(-1)
	for i, e := range o.data {
		if e <= prev {
			t.Fatal("Nonsorted oracle ", e, prev)
		}
		if e >= key {
			return oracleIterator{o: o, index: i}
		}
	}
	return oracleIterator{o: o, index: len(o.data)}
}

func (o *oracle) Delete(key testItem) bool {
	for i, e := range o.data {
		if e == key {
			newData := make([]testItem, len(o.data) - 1)
			copy(newData, o.data[0:i])
			copy(newData[i:], o.data[i+1:])
			o.data = newData
			return true
		}
	}
	return false
}


//
// Test iterator
//
type oracleIterator struct {
	o     *oracle
	index int
}

func (oiter oracleIterator) Done() bool {
	return oiter.index >= len(oiter.o.data)
}

func (oiter oracleIterator) Item() testItem {
	return oiter.o.data[oiter.index]
}

func (oiter oracleIterator) Next() oracleIterator {
	return oracleIterator{oiter.o, oiter.index + 1}
}

func compareContents(t *testing.T, o *oracle, tree *Root) {
	log.Print("Start compare")
	oiter := o.Find(t, testItem(-1))
	titer := tree.Find(testItem(-1))
	for !oiter.Done() {
		if titer.Done() {
			t.Fatal("titer.done")
		}
		log.Print("Item: ", oiter.Item(), titer.Item())
		if titer.Item().(testItem) != oiter.Item() {
			t.Fatal(titer.Item(), oiter.Item())
		}
		oiter = oiter.Next()
		titer = titer.Next()
	}
	if !titer.Done() {
		log.Print("Excess item: ", titer.Item())
		t.Fatal("!titer.done")
	}
	log.Print("End compare")
}

func TestRandomized(t *testing.T) {
	o := newOracle()
	tree := NewTree(compareItems)
	r := rand.New(rand.NewSource(0))
	for i := 0; i < 100; i++ {
		op := r.Intn(100)
		if op < 50 {
			key := r.Intn(1000)
			log.Print("Insert ", key)
			o.Insert(testItem(key))
			tree.Insert(testItem(key))
			compareContents(t, o, tree)
		} else if (op < 75 && o.Len() > 0) {
			key := o.RandomExistingKey(r)
			log.Print("DeleteExisting ", key)
			o.Delete(key)
			if !tree.DeleteWithKey(key) {
				t.Fatal("DeleteExisting", key)
			}
			compareContents(t, o, tree)
		}

	}
}
