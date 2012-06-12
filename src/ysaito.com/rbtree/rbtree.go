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
	node := iter.node

	for node != nil {
		if node.right != nil {
			node = node.right
			for node.left != nil {
				node = node.left
			}
			iter.node = node
			return
		}
		for node.right == nil {
			node = node.parent
			if node == nil {
				iter.node = nil
				return
			}
		}
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
			} else if (n.right != nil) {
				n = n.right
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
			root.leftRotate(n.parent)
			n = n.left
			continue
		}
		if n.isLeftChild() && n.parent.isRightChild() {
			root.rightRotate(n.parent)
			n = n.right
			continue
		}

		// Case 5: parent is read, uncle is black (2)
		n.parent.color = Black
		grandparent.color = Red
		if n.isLeftChild() {
			root.rightRotate(grandparent)
		} else {
			root.leftRotate(grandparent)
		}
		break
	}
	return true
}

func (root *Root) remove(n *node) {
	if n.parent == nil {
		root.tree = nil
		return
	}
	leftChildIsNonLeaf := (n.left != nil && !n.left.isLeaf())
	rightChildIsNonLeaf := (n.right != nil && !n.right.isLeaf())
	if n.left == nil {
		if !rightChildIsNonLeaf {
			if n.right != nil {
				n.right.parent := n.parent
				n.right.color = n.color
			}
			if n.isLeftChild() {
				n.parent.left = n.right
			} else {
				n.parent.right = n.right
			}
			return
		} else {
			// right child is nonleaf. fallthrough
		}
	}

	if n.right == nil {
		if !leftChildIsNonLeaf {
			if n.left != nil {
				n.left.parent := n.parent
				n.left.color = n.color
			}
			if n.isLeftChild() {
				n.parent.left = n.left
			} else {
				n.parent.right = n.left
			}
			return
		} else {
			// left child is nonleaf. fallthrough
		}
	}
	assert(n.right && n.left)

	var child *node
	if rightChildIsNonLeaf {
		child = n.right
	} else {
		child = n.left
	}

}

/*
    X		     Y
  A   Y	    =>     X   C
     B C 	  A B
*/
func (root *Root) leftRotate(x *node) {
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
func (root *Root) rightRotate(y *node) {
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
