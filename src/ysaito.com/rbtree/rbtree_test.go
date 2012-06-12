//
// Created by Yaz Saito on 06/10/12.
//

package rbtree
import "testing"
import "fmt"
import "math/rand"
import "log"
import "sort"

func compareInts(i1, i2 interface{}) int {
	return i1.(int) - i2.(int)
}

func TestEmpty(t *testing.T) {
	tree := NewTree(compareInts)
	if tree.Len() != 0 { t.Error("Len!=0") }
	iter := tree.Find(10)
	if !iter.Done() { t.Fail() }
}

func TestBasic(t *testing.T) {
	tree := NewTree(compareInts)
	if !tree.Insert(10, "blah") {t.Fail()}
	if tree.Insert(10, "xxx") {t.Fail()}

	if tree.Len() != 1 { t.Error("Len!=1") }
	iter := tree.Find(10)
	if iter.Done() { t.Fail() }
	if iter.Key().(int) != 10 { t.Error("Wrong key: ", iter.Key()) }
	if iter.Value().(string) != "blah" { t.Error("Wrong value: ", iter.Value()) }

	iter = tree.Find(11)
	if !iter.Done() { t.Fail() }

	iter = tree.Find(9)
	if iter.Done() { t.Fail() }
	if iter.Key().(int) != 10 { t.Error("Wrong key: ", iter.Key()) }
	if iter.Value().(string) != "blah" { t.Error("Wrong value: ", iter.Value()) }
	log.Print("done")
}

type oracleElement struct {
	key int
	value string
}

type oracle struct {
	data []oracleElement
}

type oracleIterator struct {
	o *oracle
	index int
}

func newOracle() *oracle {
	return &oracle{data : make([]oracleElement, 100)}
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
	newData := make([]oracleElement, n)
	copy(newData, o.data)
	o.data = newData
	o.data[n - 1].key = key
	o.data[n - 1].value = value
	sort.Sort(o)
	return true
}

func (o *oracle) Find(t *testing.T, key int) oracleIterator {
	prev := oracleElement{key: -1, value: ""}
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

func TestRandomized(t *testing.T) {
	o := newOracle()
	tree := NewTree(compareInts)
	r := rand.New(rand.NewSource(0))
	for i := 0; i < 100; i++ {
		key := r.Intn(1000)
		value := fmt.Sprintf("k%d", key)
		log.Print("Insert ", key)
		o.Insert(key, value)
		tree.Insert(key, value)
	}
}