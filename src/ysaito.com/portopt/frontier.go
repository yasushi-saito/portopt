//
// Created by Yaz Saito on 06/10/12.
//

package portopt
import "github.com/stathat/treap"

type frontier struct {
	tree treap.Tree
}

func floatLess(p, q interface{}) bool {
	return p.(float64) < q.(float64)
}

func newFrontier() *frontier {
	f := new(frontier)
	f.tree := NewTree(floatLess)
}

func (f *frontier) Insert(float64 mean, float64 stddev) {
	f.tree.Insert(mean, stddeve)
}

func (f *frontier) Insert(float64 mean, float64 stddev) {
