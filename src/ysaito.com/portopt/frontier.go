//
// Created by Yaz Saito on 06/10/12.
//

package portopt
import "github.com/yasushi-saito/rbtree"

type frontierItem struct {
	x float64
	y float64
	item interface{}
}

type frontier struct {
	tree *rbtree.Tree
}

type frontierIterator struct {
	iter rbtree.Iterator
}

func (iter frontierIterator) Item() interface{} {
	return iter.iter.Item().(frontierItem).item
}

func (iter frontierIterator) Mean() float64 {
	return iter.iter.Item().(frontierItem).x
}

func (iter frontierIterator) Stddev() float64 {
	return iter.iter.Item().(frontierItem).y
}

func (iter frontierIterator) Next() frontierIterator {
	return frontierIterator{iter.iter.Next()}
}

func (iter frontierIterator) Done() bool {
	return iter.iter.Limit()
}

func getItem(iter rbtree.Iterator) frontierItem {
	return iter.Item().(frontierItem)
}

func newFrontier() *frontier {
	f := new(frontier)
	f.tree = rbtree.NewTree(func(p1, p2 rbtree.Item) int {
		f1 := p1.(frontierItem)
		f2 := p2.(frontierItem)
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
func isConvex(left, middle, right frontierItem) bool {
	doAssert(left.x < middle.x && middle.x < right.x,
		"Wrong order: ", left, middle, right)
	width := right.x - left.x
	height := right.y - left.y
	expectedY := left.y + height * (middle.x - left.x) / width
	return middle.y > expectedY
}

func isConcave(left, middle, right frontierItem) bool {
	doAssert(left.x < middle.x && middle.x < right.x,
		"Wrong order: ", left, middle, right)
	width := right.x - left.x
	height := right.y - left.y
	expectedY := left.y + height * (middle.x - left.x) / width
	return middle.y < expectedY
}

func (f *frontier) MaxX() float64 {
	if f.tree.Len() == 0 {
		return -1
	}
	return getItem(f.tree.Max()).x
}

func (f *frontier) Iterate() frontierIterator {
	return frontierIterator{iter: f.tree.Min()}
}

// Returns true if the entry was inserted (i.e., it is part of the
// efficient frontier)
func (f *frontier) Insert(x, y float64, item interface{}) bool {
	maxX := f.MaxX()

	thisElem := frontierItem{x: x, y: y, item: item}
	if f.tree.Len() == 0 {
		f.tree.Insert(thisElem)
		doAssert((x > maxX), " nb=", x, " mean=", x, " maxx=", maxX)
		return true
	}

	leftIter := f.tree.FindLE(thisElem)
	rightIter := f.tree.FindGE(thisElem)

	if !rightIter.Limit() && x == getItem(rightIter).x {
		// Exact match found
		if y < getItem(rightIter).y {
			f.tree.DeleteWithIterator(rightIter)
			f.tree.Insert(thisElem)
			return true
		}
		return false
	}
	if leftIter.NegativeLimit() {
		if y < getItem(f.tree.Min()).y {
			f.tree.Insert(thisElem)
			return true
		}
		return false
	}
	if rightIter.Limit() {
		f.tree.Insert(thisElem)
		doAssert((x > maxX), " nb=", x, " mean=", x, " maxx=", maxX)
		return true
	}

	if isConcave(
		getItem(leftIter),
		thisElem,
		getItem(rightIter)) {
		f.tree.Insert(thisElem)
		// Remove elements to the left that make the curve concave
		// after adding thisElem.
		ti := leftIter.Prev()
		for !ti.NegativeLimit() && !isConcave(getItem(ti), getItem(leftIter), thisElem) {
			tmp := leftIter
			leftIter = ti
			ti = ti.Prev()
			f.tree.DeleteWithIterator(tmp)
		}

		// Remove elements to the right that make the curve concave
		// after adding thisElem.
		ti = rightIter.Next()
		for !ti.Limit() && !isConcave(thisElem, getItem(rightIter), getItem(ti)) {
			tmp := rightIter
			rightIter = ti
			ti = ti.Next()
			f.tree.DeleteWithIterator(tmp)
		}
		return true
	}
	return false
}

