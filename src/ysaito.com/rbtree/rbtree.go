//
// Created by Yaz Saito on 06/10/12.
//

package rbtree
type color {
	Red = iota
	Black = 1 + iota
}

type node struct {
	key interface{}
	value interface{}
	parent, left, right *tree
	color color
}

type Root {
	tree *node
	func compare(k1, k2 interface{}) bool
}

func (n *node) isLeftChild() bool {
	return n == n.parent.left
}

func (n *node) isRightChild() bool {
	return n == n.parent.right
}

func NewTree(func compare(k1, k2 interface{})) *Root {
	r := new(Root)
	r.compare := compare
	return r
}

func (r *Root) doInsert(n *tree) {
	if r.parent == nil {
		n.Position(r.parent)
		r.parent = &n
		return true
	}
	parent := r.parent
	var next tree
	comp := r.compare(key, parent.key)
	if (comp == 0) {
		return false
	} else if (comp < 0) {
		if !parent.left {
			n.Position(parent.left)
			parent.left = &n
			return true
		} else {
			parent = parent.left
		}
	} else {
		if !parent.right {
			n.Position(parent.right)
			parent.right = &n
			return true
		} else {
			parent = parent.right
		}
	}
}

func (r *Root) Insert(key interface{}, value interface{}) (bool) {
	n := new(tree)
	n.key = key
	n.value = value
	n.color = Red

	if r.parent == nil {
		n.color = black
		n.Position(r.parent)
		r.parent = n
		return true
	}
	inserted := r.doInsert(&n)
	if !inserted { return false }

	n.color = Red
	for n != r.tree && n.parent.color == Red {
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
					left_rotate(T, n);
				}
				/* case 3 */
				n.parent.color = black;
				n.parent.parent.color = red;
				root.rightRotate(n.parent.parent);
			}
		} else {
			/* repeat the "if" part with right and left
			 exchanged */
		}
	}
	/* Color the root black */
	root.tree.color = black
}

/*
    X		     Y
  A   Y	    =>     X   C
     B C 	  A B
*/
func (root *Root) leftRotate(x node) {
	y := x.right
	x.right = y.left;
	if y.left != NULL {
		y.left.parent = x;
	}
	y.parent = x.parent;
	if x.parent == root.tree {
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
func (root *Root) rightRotate(y node) {
	x := y.left

	// Move "B"
	y.left = x.right;
	if x.right != NULL {
		x.right.parent = y;
	}

	x.parent = y.parent;
	if y.parent == root.tree {
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
