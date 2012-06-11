//
// Created by Yaz Saito on 06/10/12.
//

package rbtree
import "testing"

func compareInts(i1, i2 interface{}) int {
	return i1.(int) - i2.(int)
}

func TestBasic(t *testing.T) {
	tree := NewTree(compareInts)
	tree.Insert(10, "blah")
}
