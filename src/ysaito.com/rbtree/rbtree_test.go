//
// Created by Yaz Saito on 06/10/12.
//

package rbtree

import "testing"
import "math/rand"
import "log"
import "sort"

type testItem int

func newTestTree() *Root {
	return NewTree(func(i1, i2 Item) int {
		return int(i1.(testItem)) - int(i2.(testItem))
	})
}

func TestEmpty(t *testing.T) {
	tree := newTestTree()
	if tree.Len() != 0 {
		t.Error("Len!=0")
	}
	iter := tree.FindGE(testItem(10))
	if !iter.Done() {
		t.Error("Not empty")
	}
}

func TestBasic(t *testing.T) {
	tree := newTestTree()
	if !tree.Insert(testItem(10)) {
		t.Error("Insert1")
	}
	if tree.Insert(testItem(10)) {
		t.Error("Insert2")
	}

	if tree.Len() != 1 {
		t.Error("Len!=1")
	}
	iter := tree.FindGE(testItem(10))
	if iter.Done() {
		t.Error()
	}
	if iter.Item().(testItem) != 10 {
		t.Error("Wrong item: ", iter.Item())
	}
	iter = tree.FindGE(testItem(11))
	if !iter.Done() {
		t.Error()
	}

	iter = tree.FindGE(testItem(9))
	if iter.Done() {
		t.Error()
	}
	if iter.Item().(testItem) != 10 {
		t.Error("Wrong item: ", iter.Item())
	}

	item := tree.Get(testItem(10))
	if item == nil || item.(testItem) != 10 {
		t.Error("Wrong Get: ", item)
	}
	log.Print("done")
}

func TestDelete(t *testing.T) {
	tree := newTestTree()
	if tree.DeleteWithKey(testItem(10)) {
		t.Error()
	}
	if tree.Len() != 0 {
		t.Error()
	}

	if !tree.Insert(testItem(10)) {
		t.Error()
	}
	if !tree.DeleteWithKey(testItem(10)) {
		t.Error()
	}
	if tree.Len() != 0 {
		t.Error()
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

func (o *oracle) FindGE(t *testing.T, key testItem) oracleIterator {
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

func (o *oracle) FindLE(t *testing.T, key testItem) oracleIterator {
	iter := o.FindGE(t, key)
	if !iter.Done() && o.data[iter.index] == key {
		return iter
	} else if (iter.index == 0) {
		return oracleIterator{o, len(o.data)}
	}
	return oracleIterator{o, iter.index - 1}
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

func compareContents(t *testing.T, oiter oracleIterator, titer Iterator) {
	for !oiter.Done() {
		if titer.Done() {
			t.Fatal("titer.done")
		}
		// log.Print("Item: ", oiter.Item(), titer.Item())
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
}

const testVerbose = false

func TestRandomized(t *testing.T) {
	const numKeys = 1000

	o := newOracle()
	tree := newTestTree()
	r := rand.New(rand.NewSource(0))
	for i := 0; i < 10000; i++ {
		op := r.Intn(100)
		if op < 50 {
			key := r.Intn(numKeys)
			if testVerbose { log.Print("Insert ", key) }
			o.Insert(testItem(key))
			tree.Insert(testItem(key))
			compareContents(t, o.FindGE(t, testItem(-1)), tree.FindGE(testItem(-1)))
		} else if op < 90 && o.Len() > 0 {
			key := o.RandomExistingKey(r)
			if testVerbose { log.Print("DeleteExisting ", key) }
			o.Delete(key)
			if !tree.DeleteWithKey(key) {
				t.Fatal("DeleteExisting", key)
			}
			compareContents(t, o.FindGE(t, testItem(-1)), tree.FindGE(testItem(-1)))
		} else if (op < 95) {
			key := testItem(r.Intn(numKeys))
			if testVerbose { log.Print("FindGE ", key) }
			compareContents(t, o.FindGE(t, key), tree.FindGE(key))
		} else {
			key := testItem(r.Intn(numKeys))
			if testVerbose { log.Print("FindLE ", key) }
			compareContents(t, o.FindLE(t, key), tree.FindLE(key))
		}
	}
}
