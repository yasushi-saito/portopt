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
	compare func(k1, k2 interface{}) int
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

func (r *Root) doInsert(n *node) bool {
	if r.tree == nil {
		n.parent = nil
		r.tree = n
		return true
	}
	parent := r.tree
	for true {
		comp := r.compare(n.key, parent.key)
		if (comp == 0) {
			return false
		} else if (comp < 0) {
			if parent.left == nil {
				n.parent = parent.left
				parent.left = n
				return true
			} else {
				parent = parent.left
			}
		} else {
			if parent.right == nil {
				n.parent = parent.right
				parent.right = n
				return true
			} else {
				parent = parent.right
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
	for n != root.tree && n.parent.color == Red {
		if n.parent == n.parent.parent.left {
			/* If x's parent is a left, y is x's right 'uncle' */
			y := n.parent.parent.right;
			if y.color == Red {
				/* case 1 - change the colors */
				n.parent.color = Black;
				y.color = Black;
				n.parent.parent.color = Red;
				/* Move x up the tree */
				n = n.parent.parent;
			} else {
				if n == n.parent.right {
					/* and x is to the right */
					/* case 2 - move x up and rotate */
					n = n.parent;
					root.leftRotate(n);
				}
				/* case 3 */
				n.parent.color = Black;
				n.parent.parent.color = Red;
				root.rightRotate(n.parent.parent);
			}
		} else {
			/* repeat the "if" part with right and left
			 exchanged */
		}
	}
	/* Color the root black */
	root.tree.color = Black
	return true
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
