//
// Created by Yaz Saito on 06/10/12.
//

package rbtree
import "testing"
import "fmt"
import "math/rand"
import "log"
import "sort"

type testItem struct {
	key int
	value string
}

func testKey(k int) testItem {
	return testItem{k, ""}
}

func compareInts(i1, i2 Item) int {
	return i1.(testItem).key - i2.(testItem).key
}

func TestEmpty(t *testing.T) {
	tree := NewTree(compareInts)
	if tree.Len() != 0 { t.Error("Len!=0") }
	iter := tree.Find(testKey(10))
	if !iter.Done() { t.Fail() }
}

func TestBasic(t *testing.T) {
	tree := NewTree(compareInts)
	if !tree.Insert(testItem{10, "blah"}) {t.Fail()}
	if tree.Insert(testItem{10, "xxx"}) {t.Fail()}

	if tree.Len() != 1 { t.Error("Len!=1") }
	iter := tree.Find(testKey(10))
	if iter.Done() { t.Fail() }
	if iter.Item().(testItem).key != 10 || iter.Item().(testItem).value != "blah" {
		t.Error("Wrong item: ", iter.Item())
	}
	iter = tree.Find(testKey(11))
	if !iter.Done() { t.Fail() }

	iter = tree.Find(testKey(9))
	if iter.Done() { t.Fail() }

	if iter.Item().(testItem).key != 10 || iter.Item().(testItem).value != "blah" {
		t.Error("Wrong item: ", iter.Item())
	}
	log.Print("done")
}

func TestRemove(t *testing.T) {

}

type oracle struct {
	data []testItem
}

func newOracle() *oracle {
	return &oracle{data : make([]testItem, 0)}
}

func (o *oracle) Len() int {
	return len(o.data)
}

func (o *oracle) Less(i, j int) bool {
	return o.data[i].key < o.data[j].key
}

func (o *oracle) Swap(i, j int) {
	e := o.data[j]
	o.data[j] = o.data[i]
	o.data[i] = e
}

func (o *oracle) Insert(key int, value string) bool {
	for _, e := range o.data {
		if e.key == key { return false }
	}

	n := len(o.data) + 1
	newData := make([]testItem, n)
	copy(newData, o.data)
	o.data = newData
	o.data[n - 1].key = key
	o.data[n - 1].value = value
	sort.Sort(o)
	return true
}

func (o *oracle) Find(t *testing.T, key int) oracleIterator {
	prev := testItem{key: -1, value: ""}
	for i, e := range o.data {
		if e.key <= prev.key {
			t.Fatal("Nonsorted oracle ", e, prev)
		}
		if e.key >= key {
			return oracleIterator{o: o, index: i}
		}
	}
	return oracleIterator{o: o, index: len(o.data)}
}

type oracleIterator struct {
	o *oracle
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
	oiter := o.Find(t, -1)
	titer := tree.Find(testKey(-1))
	for !oiter.Done() {
		if titer.Done() {
			t.Fatal("titer.done")
		}
		if titer.Item().(testItem).key != oiter.Item().key {
			t.Fatal(titer.Item(), oiter.Item())
		}
		if titer.Item().(testItem).value != oiter.Item().value {
			t.Fatal(titer.Item(), oiter.Item())
		}
		oiter = oiter.Next()
		titer = titer.Next()
	}
	if !titer.Done() {
		t.Fatal("!titer.done")
	}
}

func TestRandomized(t *testing.T) {
	o := newOracle()
	tree := NewTree(compareInts)
	r := rand.New(rand.NewSource(0))
	for i := 0; i < 100; i++ {
		key := r.Intn(1000)
		value := fmt.Sprintf("k%d", key)
		log.Print("Insert ", key)
		o.Insert(key, value)
		tree.Insert(testItem{key, value})

		compareContents(t, o, tree)
	}
}