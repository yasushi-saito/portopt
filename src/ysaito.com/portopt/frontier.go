//
// Created by Yaz Saito on 06/10/12.
//

package portopt
import "bytes"
import "fmt"
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
			f.maybeRemoveLeftElements(thisElem, leftIter.Prev())
			f.maybeRemoveRightElements(thisElem, rightIter.Next())
			f.tree.DeleteWithIterator(rightIter)
			f.tree.Insert(thisElem)
			return true
		}
		return false
	}
	if leftIter.NegativeLimit() {
		if y < getItem(f.tree.Min()).y {
			f.tree.Insert(thisElem)
			f.maybeRemoveRightElements(thisElem, rightIter)
			return true
		}
		return false
	}
	if rightIter.Limit() {
		f.tree.Insert(thisElem)
		f.maybeRemoveLeftElements(thisElem, leftIter)
		return true
	}

	if isConcave(
		getItem(leftIter),
		thisElem,
		getItem(rightIter)) {
		f.tree.Insert(thisElem)

		f.maybeRemoveLeftElements(thisElem, leftIter)
		f.maybeRemoveRightElements(thisElem, rightIter)
		return true
	}
	return false
}

// Remove elements <= leftIter that make the curve concave after
// adding thisElem.
func (f *frontier) maybeRemoveLeftElements(
	thisElem frontierItem, leftIter rbtree.Iterator) {
	if leftIter.NegativeLimit() {
		return
	}

	if getItem(leftIter).y >= thisElem.y {
		tmp := leftIter
		leftIter = leftIter.Prev()
		f.tree.DeleteWithIterator(tmp)
	}

	if leftIter.NegativeLimit() {
		return
	}
	ti := leftIter.Prev()

	for !ti.NegativeLimit() && !isConcave(getItem(ti), getItem(leftIter), thisElem) {
		tmp := leftIter
		leftIter = ti
		ti = ti.Prev()
		f.tree.DeleteWithIterator(tmp)
	}
}


// Remove elements >= rightIter that make the curve concave after
// adding thisElem.
func (f *frontier) maybeRemoveRightElements(
	thisElem frontierItem, rightIter rbtree.Iterator) {
	if rightIter.Limit() {
		return
	}

	ti := rightIter.Next()
	for !ti.Limit() && !isConcave(thisElem, getItem(rightIter), getItem(ti)) {
		tmp := rightIter
		rightIter = ti
		ti = ti.Next()
		f.tree.DeleteWithIterator(tmp)
	}
}

func (f *frontier) String() string {
	buf := bytes.NewBufferString("")
	for iter := f.Iterate(); !iter.Done(); iter = iter.Next() {
		fmt.Fprint(buf, "P: mean=", iter.Mean(), " stddev=", iter.Stddev(),
			" port=", iter.Item(), "\n")
	}
	return buf.String()
}

