//
// Created by Yaz Saito on 06/10/12.
//

package portopt
import "github.com/yasushi-saito/rbtree"

type frontierElement struct {
	x float64
	y float64
}

type frontier struct {
	tree *rbtree.Tree
}

func getElement(iter rbtree.Iterator) frontierElement {
	return iter.Item().(frontierElement)
}

func newFrontier() *frontier {
	f := new(frontier)
	f.tree = rbtree.NewTree(func(p1, p2 rbtree.Item) int {
		f1 := p1.(frontierElement)
		f2 := p2.(frontierElement)
		if f1.x < f2.x {
			return -1
		} else if f1.x == f2.x {
			return 0
		}
		return 1
	})
	return f
}

// See if points <left, middle, right> form a convex curve
func isConvex(left, middle, right frontierElement) bool {
	doAssert(left.x < middle.x && middle.x < right.x,
		"Wrong order: ", left, middle, right)
	width := right.x - left.x
	height := right.y - left.y
	expectedY := left.y + height * (middle.x - left.x) / width
	return middle.y > expectedY
}

func (f *frontier) Insert(x, y float64) bool {
	thisElem := frontierElement{x: x, y: y}
	leftIter := f.tree.FindLE(thisElem)
	rightIter := f.tree.FindGE(thisElem)

	if !rightIter.Limit() && x == getElement(rightIter).x {
		// Exact match found
		if y > getElement(rightIter).y {
			f.tree.DeleteWithIterator(rightIter)
			f.tree.Insert(thisElem)
			return true
		}
		return false
	}
	if leftIter.NegativeLimit() || rightIter.Limit() {
		// The new x is beyond the current boundary.
		f.tree.Insert(thisElem)
		return true
	}
	if isConvex(
		getElement(leftIter),
		thisElem,
		getElement(rightIter)) {
		f.tree.Insert(thisElem)
		// Remove elements to the left that make the curve concave
		// after adding thisElem.
		ti := leftIter.Prev()
		for !ti.NegativeLimit() && !isConvex(getElement(ti), getElement(leftIter), thisElem) {
			tmp := leftIter
			leftIter = ti
			ti = ti.Prev()
			f.tree.DeleteWithIterator(tmp)
		}

		// Remove elements to the right that make the curve concave
		// after adding thisElem.
		ti = rightIter.Next()
		for !ti.Limit() && !isConvex(thisElem, getElement(rightIter), getElement(ti)) {
			tmp := rightIter
			rightIter = ti
			ti = ti.Next()
			f.tree.DeleteWithIterator(tmp)
		}
		return true
	}
	return false
}

