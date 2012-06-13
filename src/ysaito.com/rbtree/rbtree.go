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

func (n *node) isLeftChild() bool {
	return n == n.parent.left
}

func (n *node) isRightChild() bool {
	return n == n.parent.right
}

func (n *node) sibling() *node {
	doAssert(n.parent != nil)
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
	root.doDelete3(iter.node)
}

func doAssert(b bool) {
	if !b {
		panic("rbtree internal assertion failed")
	}
}

func (root *Root) doDelete3(n *node) {
	root.count--
	if n.left != nil && n.right != nil {
		pred := maxPredecessor(n)
		n.item = pred.item
		n = pred
	}
	doAssert(n.left == nil || n.right == nil);
	child := n.right
	if child == nil {
		child = n.left
	}
	if (n.color == Black) {
		n.color = getColor(child)
		root.deleteCase1(n)
	}
	root.replaceNode(n, child)
	if (n.parent == nil && child != nil) {
		child.color = Black
	}
}

func (root* Root) deleteCase1(n *node) {
	if (n.parent == nil) {
		return
	} else {
		root.deleteCase2(n)
	}
}

func (root* Root) deleteCase2(n *node) {
	if (getColor(n.sibling()) == Red) {
		n.parent.color = Red;
		n.sibling().color = Black;
		if (n == n.parent.left) {
			root.rotateLeft(n.parent);
		} else {
			root.rotateRight(n.parent);
		}
	}
	root.deleteCase3(n)
}

func (root* Root) deleteCase3(n *node) {
	if (getColor(n.parent) == Black &&
		getColor(n.sibling()) == Black &&
		getColor(n.sibling().left) == Black &&
		getColor(n.sibling().right) == Black) {
		n.sibling().color = Red;
		root.deleteCase1(n.parent);
	} else {
		root.deleteCase4(n);
	}
}

func (root* Root) deleteCase4(n *node) {
	if (getColor(n.parent) == Red &&
		getColor(n.sibling()) == Black &&
		getColor(n.sibling().left) == Black &&
		getColor(n.sibling().right) == Black) {
		n.sibling().color = Red;
		n.parent.color = Black;
	} else {
		root.deleteCase5(n);
	}
}

func (root* Root) deleteCase5(n *node) {
	if (n == n.parent.left &&
		getColor(n.sibling()) == Black &&
		getColor(n.sibling().left) == Red &&
		getColor(n.sibling().right) == Black) {
		n.sibling().color = Red;
		n.sibling().left.color = Black;
		root.rotateRight(n.sibling());
	} else if (n == n.parent.right &&
		getColor(n.sibling()) == Black &&
		getColor(n.sibling().right) == Red &&
		getColor(n.sibling().left) == Black) {
		n.sibling().color = Red;
		n.sibling().right.color = Black;
		root.rotateLeft(n.sibling());
	}
	root.deleteCase6(n);
}

func (root* Root) deleteCase6(n *node) {
	n.sibling().color = getColor(n.parent);
	n.parent.color = Black;
	if (n == n.parent.left) {
		doAssert(getColor(n.sibling().right) == Red);
		n.sibling().right.color = Black;
		root.rotateLeft(n.parent);
	} else {
		doAssert(getColor(n.sibling().left) == Red);
		n.sibling().left.color = Black;
		root.rotateRight(n.parent);
	}
}

func (root *Root) replaceNode(oldn, newn *node) {
	if (oldn.parent == nil) {
		root.tree = newn;
	} else {
		if (oldn == oldn.parent.left) {
			oldn.parent.left = newn;
		} else {
			oldn.parent.right = newn;
		}
	}
	if (newn != nil) {
		newn.parent = oldn.parent;
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
