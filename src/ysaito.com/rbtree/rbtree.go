//
// Created by Yaz Saito on 06/10/12.
//

package rbtree

const Red = iota
const Black = 1 + iota

type Item interface{}

// CompareFunc returns 0 if a==b, <0 if a<b, >0 if a>b.
type CompareFunc func(a, b Item) int

type node struct {
	item                Item
	parent, left, right *node
	color               int
}

type Root struct {
	tree    *node
	count   int
	compare CompareFunc
}

type Iterator struct {
	root *Root
	node *node
}

func getColor(n *node) int {
	if n == nil {
		return Black
	}
	return n.color
}

func newIterator(r *Root, n *node) Iterator {
	return Iterator{r, n}
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
		return newIterator(iter.root, minSuccessor(n))
	}

	for n != nil {
		p := n.parent
		if p == nil {
			return newIterator(iter.root, nil)
		}
		if n.isLeftChild() {
			return newIterator(iter.root, p)
		}
		n = p
	}
	return newIterator(iter.root, nil)
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

func (n *node) isLeaf() bool {
	return n.left == nil && n.right == nil
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
		if comp == 0 {
			return false
		} else if comp < 0 {
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
			return newIterator(root, nil)
		}
		comp := root.compare(key, n.item)
		if comp == 0 {
			return newIterator(root, n)
		} else if comp < 0 {
			if n.left != nil {
				n = n.left
			} else {
				return newIterator(root, n)
			}
		} else {
			if n.right != nil {
				n = n.right
			} else {
				return newIterator(root, n.parent)
			}
		}
	}
	panic("should not reach here")

}

func (root *Root) Insert(item Item) bool {
	n := new(node)
	n.item = item
	n.color = Red

	// TODO: delay creating n until it is found to be inserted
	inserted := root.doInsert(n)
	if !inserted {
		return false
	}

	n.color = Red

	for true {
		// Case 1: N is at the root
		if n.parent == nil {
			n.color = Black
			break
		}

		// Case 2: The parent is black, so the tree already
		// satisfies the RB properties
		if n.parent.color == Black {
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
			grandparent.color = Red
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
		root.DeleteWithIterator(iter)
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

	// If toDelete is not the minimum node, find its predecessor
	// (call it pred). We will eventually replace toDelete with
	// pred. Below, T=toDelete, P=pred
 	//
	//     T            P
	//   P   R  =>    C   R
	//     C
	pred := toDelete
	if pred.left != nil {
		pred := maxPredecessor(toDelete)
		toDelete.item = pred.item
		// TODO: this will invalidate the iterator.
		// fix.
	}

	// pred should have at most one child. Replace pred's contents with
	// the child's.
	n := pred
	var child *node
	if n.right != nil {
		child := n.right
		n.item = child.item
		n.left = child.left
		n.right = child.right
	} else if n.left != nil {
		child = n.left
		n.item = child.item
		n.left = child.left
		n.right = child.right
	} else {
		// n is a leaf
		if n.parent == nil {
			root.tree = nil
			return;
		} else if n.isLeftChild() {
			n.parent.left = nil
			child = n.parent.right
		} else {
			n.parent.right = nil
			child = n.parent.left
		}
	}

	// Fix the color of the child
	for true {
		if n.color != Black {
			break
		}
		if child != nil && child.color == Red {
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
		if n.parent.color == Black &&
			s.color == Black &&
			getColor(s.left) == Black &&
			getColor(s.right) == Black {
			s.color = Red
			n = n.parent
			continue
		}
		if n.parent.color == Red &&
			s.color == Black &&
			getColor(s.left) == Black &&
			getColor(s.right) == Black {
			s.color = Red
			n.parent.color = Black
			break
		}
		if s.color == Black {
			if n.isLeftChild() &&
				getColor(s.right) == Black &&
				getColor(s.left) == Red {
				s.color = Red
				s.left.color = Black
				root.rotateLeft(s)
			} else if n.isRightChild() &&
				getColor(s.left) == Black &&
				getColor(s.right) == Red {
				s.color = Red
				s.right.color = Black
				root.rotateLeft(s)
			}
		}
		s.color = n.parent.color
		n.parent.color = Black
		if n.isLeftChild() {
			s.right.color = Black
			root.rotateLeft(n.parent)
		} else {
			s.left.color = Black
			root.rotateRight(n.parent)
		}
		break
	}

}

/*
    X		     Y
  A   Y	    =>     X   C
     B C 	  A B
*/
func (root *Root) rotateLeft(x *node) {
	y := x.right
	x.right = y.left
	if y.left != nil {
		y.left.parent = x
	}
	y.parent = x.parent
	if x.parent == nil {
		root.tree = y
	} else {
		if x.isLeftChild() {
			x.parent.left = y
		} else {
			x.parent.right = y
		}
	}
	y.left = x
	x.parent = y
}

/*
     Y           X
   X   C  =>   A   Y
  A B             B C
*/
func (root *Root) rotateRight(y *node) {
	x := y.left

	// Move "B"
	y.left = x.right
	if x.right != nil {
		x.right.parent = y
	}

	x.parent = y.parent
	if y.parent == nil {
		root.tree = x
	} else {
		if y.isLeftChild() {
			y.parent.left = x
		} else {
			y.parent.right = x
		}
	}
	x.right = y
	y.parent = x
}
