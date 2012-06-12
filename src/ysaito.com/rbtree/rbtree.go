//
// Created by Yaz Saito on 06/10/12.
//

package rbtree
const Red = iota
const Black = 1 + iota

type Item interface {}

// CompareFunc returns 0 if a==b, <0 if a<b, >0 if a>b.
type CompareFunc func(a, b Item) int

type node struct {
	item Item
	parent, left, right *node
	color int
}

type Root struct {
	tree *node
	count int
	compare CompareFunc
}

type Iterator struct {
	node *node
}

func newIterator(n *node) Iterator {
	return Iterator{node: n}
}

func (iter Iterator) Done() bool {
	return iter.node == nil
}

func (iter Iterator) Item() interface{} {
	return iter.node.item
}

func (iter Iterator) Next() Iterator {
	n := iter.node

	if n.right != nil {
		return newIterator(minSuccessor(n))
	}

	for n != nil {
		p := n.parent
		if p == nil {
			return newIterator(nil)
		}
		if n.isLeftChild() {
			return newIterator(p)
		}
		n = p
	}
	return newIterator(nil)
}

func (iter *Iterator) Prev() {
	node := iter.node

	for node != nil {
		if node.left != nil {
			node = node.left
			for node.right != nil {
				node = node.right
			}
			iter.node = node
			return
		}
		for node.left == nil {
			node = node.parent
			if node == nil {
				iter.node = nil
				return
			}
		}
	}
}

func (n *node) isLeftChild() bool {
	return n == n.parent.left
}

func (n *node) isRightChild() bool {
	return n == n.parent.right
}

func (n *node) sibling() *node {
	if n.isLeftChild() {
		return n.parent.right
	} else {
		return n.parent.left
	}
	panic("Blah")
}

func NewTree(compare CompareFunc) *Root {
	r := new(Root)
	r.compare = compare
	return r
}

func (root *Root) Len() int {
	return root.count
}

func (root *Root) doInsert(n *node) bool {
	if root.tree == nil {
		n.parent = nil
		root.tree = n
		root.count++
		return true
	}
	parent := root.tree
	for true {
		comp := root.compare(n.item, parent.item)
		if (comp == 0) {
			return false
		} else if (comp < 0) {
			if parent.left == nil {
				n.parent = parent
				parent.left = n
				root.count++
				return true
			} else {
				parent = parent.left
			}
		} else {
			if parent.right == nil {
				n.parent = parent
				parent.right = n
				root.count++
				return true
			} else {
				parent = parent.right
			}
		}
	}
	panic("should not reach here")
}

func (root *Root) Get(key Item) Item {
	iter := root.Find(key)
	if iter.node != nil && root.compare(key, iter.node.item) == 0 {
		return iter.node.item
	}
	return nil
}

func (root *Root) Find(key Item) Iterator {
	n := root.tree
	for true {
		if n == nil {
			return newIterator(nil)
		}
		comp := root.compare(key, n.item)
		if (comp == 0) {
			return newIterator(n)
		} else if (comp < 0) {
			if n.left != nil {
				n = n.left
			} else {
				return newIterator(n)
			}
		} else {
			if n.right != nil {
				n = n.right
			} else {
				return newIterator(n.parent)
			}
		}
	}
	panic("should not reach here")

}

func (root *Root) Insert(item Item) (bool) {
	n := new(node)
	n.item = item
	n.color = Red

	// TODO: delay creating n until it is found to be inserted
	inserted := root.doInsert(n)
	if !inserted { return false }

	n.color = Red

	for true {
		// Case 1: N is at the root
		if n.parent == nil {
			n.color = Black
			break
		}

		// Case 2: The parent is black, so the tree already
		// satisfies the RB properties
		if (n.parent.color == Black) {
			break
		}

		// Case 3: parent and uncle are both red.
		// Then paint both black and make grandparent red.
		grandparent := n.parent.parent
		var uncle *node
		if n.parent.isLeftChild() {
			uncle = grandparent.right
		} else {
			uncle = grandparent.left
		}
		if uncle != nil && uncle.color == Red {
			n.parent.color = Black
			uncle.color = Black
			grandparent.color = Red;
			n = grandparent
			continue
		}

		// Case 4: parent is red, uncle is black (1)
		if n.isRightChild() && n.parent.isLeftChild() {
			root.rotateLeft(n.parent)
			n = n.left
			continue
		}
		if n.isLeftChild() && n.parent.isRightChild() {
			root.rotateRight(n.parent)
			n = n.right
			continue
		}

		// Case 5: parent is read, uncle is black (2)
		n.parent.color = Black
		grandparent.color = Red
		if n.isLeftChild() {
			root.rotateRight(grandparent)
		} else {
			root.rotateLeft(grandparent)
		}
		break
	}
	return true
}

func maxPredecessor(n *node) *node {
	if n.left == nil {
		return n
	}
	m := n.left
	for m.right != nil {
		m = m.right
	}
	return m
}

func minSuccessor(n *node) *node {
	if n.right == nil {
		return n
	}
	m := n.right
	for m.left != nil {
		m = m.left
	}
	return m
}

// Delete an item with the given key. Return true iff the item was
// found.
func (root *Root) DeleteWithKey(key Item) bool {
	iter := root.Find(key)
	if iter.node != nil {
		root.doDelete(iter.node)
		return true
	}
	return false
}

// Delete the current item.
//
// REQUIRES: !iter.Done()
func (root *Root) DeleteWithIterator(iter Iterator) {
	root.doDelete(iter.node)

	// Invalidate the node just to be sure
	iter.node.item = nil
}

func (root *Root) doDelete(toDelete *node) {
	root.count--
	max := maxPredecessor(toDelete)

	n := max
	var child *node
	if n.right == nil {
		child = n.left
	} else {
		child = n.right
	}

	// replace n with child
	child.parent = n.parent
	child.left = n.left
	child.right = n.right

	for true {
		if n.color != Black {
			break
		}
		if child.color == Red {
			child.color = Black
			break
		}
		if n.parent == nil {
			break
		}
		s := n.sibling()
		if s.color == Red {
			n.parent.color = Red
			s.color = Black
			if n.isLeftChild() {
				root.rotateLeft(n.parent)
			} else {
				root.rotateRight(n.parent)
			}
			break
		}
		if (n.parent.color == Black &&
			s.color == Black &&
			s.left.color == Black &&
			s.right.color == Black) {
			s.color = Red
			n = n.parent
			continue
		}
		if (n.parent.color == Red &&
			s.color == Black &&
			s.left.color == Black &&
			s.right.color == Black) {
			s.color = Red
			n.parent.color = Black
			break
		}
		if (s.color == Black) {
			if (n.isLeftChild() &&
				s.right.color == Black &&
				s.left.color == Red) {
				s.color = Red
				s.left.color = Black
				root.rotateLeft(s)
			} else if (n.isRightChild() &&
				s.left.color == Black &&
				s.right.color == Red) {
				s.color = Red
				s.right.color = Black
				root.rotateLeft(s)
			}
		}
		s.color = n.parent.color
		n.parent.color = Black
		if (n.isLeftChild()) {
			s.right.color = Black
			root.rotateLeft(n.parent)
		} else {
			s.left.color = Black
			root.rotateRight(n.parent)
		}
		break
	}

	// replace toDelete with max
	max.parent = toDelete.parent
	max.left = toDelete.left
	max.right = toDelete.right
	max.color = toDelete.color
	if max.parent == nil {
		root.tree = max
	}
}

/*
    X		     Y
  A   Y	    =>     X   C
     B C 	  A B
*/
func (root *Root) rotateLeft(x *node) {
	y := x.right
	x.right = y.left;
	if y.left != nil {
		y.left.parent = x;
	}
	y.parent = x.parent;
	if x.parent == nil {
		root.tree = y;
	} else {
		if x.isLeftChild() {
			x.parent.left = y;
		} else {
			x.parent.right = y;
		}
	}
	y.left = x;
	x.parent = y;
}

/*
     Y           X
   X   C  =>   A   Y
  A B             B C
*/
func (root *Root) rotateRight(y *node) {
	x := y.left

	// Move "B"
	y.left = x.right;
	if x.right != nil {
		x.right.parent = y;
	}

	x.parent = y.parent;
	if y.parent == nil {
		root.tree = x;
	} else {
		if y.isLeftChild() {
			y.parent.left = x;
		} else {
			y.parent.right = x;
		}
	}
	x.right = y;
	y.parent = x;
}
