//
// Created by Yaz Saito on 06/10/12.
//

package rbtree
const Red = iota
const Black = 1 + iota

type node struct {
	key interface{}
	value interface{}
	parent, left, right *node
	color int
}

type Root struct {
	tree *node
	len int
	compare func(k1, k2 interface{}) int
}

type Iterator struct {
	node *node
}

func (iter *Iterator) Done() bool {
	return iter.node == nil
}

func (iter *Iterator) Key() interface{} {
	return iter.node.key
}

func (iter *Iterator) Value() interface{} {
	return iter.node.value
}

func (iter *Iterator) Next() {
	n := iter.node

	if n.right != nil {
		iter.node = minSuccessor(n)
		return
	}

	for n != nil {
		p := n.parent
		if p == nil {
			iter.node = nil
			return
		}
		if n.isLeftChild() {
			iter.node = p
			return
		}
		n = p
	}

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

func NewTree(compare func(k1, k2 interface{}) int) *Root {
	r := new(Root)
	r.compare = compare
	return r
}

func (root *Root) Len() int {
	return root.len
}

func (root *Root) doInsert(n *node) bool {
	if root.tree == nil {
		n.parent = nil
		root.tree = n
		root.len++
		return true
	}
	parent := root.tree
	for true {
		comp := root.compare(n.key, parent.key)
		if (comp == 0) {
			return false
		} else if (comp < 0) {
			if parent.left == nil {
				n.parent = parent
				parent.left = n
				root.len++
				return true
			} else {
				parent = parent.left
			}
		} else {
			if parent.right == nil {
				n.parent = parent
				parent.right = n
				root.len++
				return true
			} else {
				parent = parent.right
			}
		}
	}
	panic("should not reach here")
}

func (root *Root) Find(key interface{}) Iterator {
	n := root.tree
	for true {
		if n == nil {
			return Iterator{node : nil}
		}
		comp := root.compare(key, n.key)
		if (comp == 0) {
			return Iterator{node : n}
		} else if (comp < 0) {
			if n.left != nil {
				n = n.left
			} else {
				return Iterator{node : n}
			}
		} else {
			if n.right != nil {
				n = n.right
			} else {
				return Iterator{node : n.parent}
			}
		}
	}
	panic("should not reach here")

}

func (root *Root) Insert(key interface{}, value interface{}) (bool) {
	n := new(node)
	n.key = key
	n.value = value
	n.color = Red

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

func (root *Root) Remove(n *node) {
	root.remove(n)
}

func (root *Root) remove(toRemove *node) {
	root.len--
	max := maxPredecessor(toRemove)

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
	max.parent = toRemove.parent
	max.left = toRemove.left
	max.right = toRemove.right
	max.color = toRemove.color
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
