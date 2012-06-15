//
// Created by Yaz Saito on 06/10/12.
//

package portopt
import "github.com/yasushi-saito/rbtree"

type frontierElement struct {
	mean float64
	stddev float64
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
		if f1.mean < f2.mean {
			return -1
		} else if f1.mean == f2.mean {
			return 0
		}
		return 1
	})
	return f
}

func isConvex(left, middle, right frontierElement) bool {
	return false
}

func (f *frontier) Insert(mean, stddev float64) bool {
	leftIter := f.tree.FindLE(mean)
	rightIter := f.tree.FindGE(mean)
	thisElem := frontierElement{mean: mean, stddev: stddev}

	if !rightIter.Limit() {
		// Exact match found
		if stddev < rightIter.Item().(frontierElement).stddev {
			f.tree.DeleteWithIterator(rightIter)
			f.tree.Insert(thisElem)
			return true
		}
		return false
	}
	if leftIter.NegativeLimit() || rightIter.Limit() {
		// The new mean is beyond the current boundary.
		f.tree.Insert(thisElem)
		return true
	}
	if isConvex(
		leftIter.Item().(frontierElement),
		thisElem,
		rightIter.Item().(frontierElement)) {
		f.tree.Insert(thisElem)
		// Remove elements to the left that make the curve concave
		// after adding thisElem.
		ti := leftIter.Prev()
		for !ti.NegativeLimit() && isConvex(getElement(ti), getElement(leftIter), thisElem) {
			tmp := leftIter
			leftIter = ti
			ti = ti.Prev()
			f.tree.DeleteWithIterator(tmp)
		}

		// Remove elements to the right that make the curve concave
		// after adding thisElem.
		ti = rightIter.Next()
		for !ti.Limit() && isConvex(thisElem, getElement(rightIter), getElement(ti)) {
			tmp := rightIter
			rightIter = ti
			ti = ti.Next()
			f.tree.DeleteWithIterator(tmp)
		}
		return true
	}
	return false
}

