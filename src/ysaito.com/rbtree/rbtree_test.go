//
// Created by Yaz Saito on 06/10/12.
//

package rbtree
import "testing"
import "log"

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
