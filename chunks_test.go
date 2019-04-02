package knxbaosip

import (
	"reflect"
	"testing"
)

var in = [][]int{
	{3, 5, 6, 7, 11, 12, 13, 18, 22, 23, 24},
	{0, 1, 2},
	{5, 7, 9, 0, 32},
	{255, 256, 243, 477, 478, 479, 480, 4},
}

var out = [][][]int{
	{{3, 1}, {5, 3}, {11, 3}, {18, 1}, {22, 3}},
	{{0, 3}},
	{{5, 1}, {7, 1}, {9, 1}, {0, 1}, {32, 1}},
	{{255, 2}, {243, 1}, {477, 4}, {4, 1}},
}

func TestChunks(t *testing.T) {
	var got, wanted [][]int
	for i, l := range in {
		got = makeChunks(l)
		wanted = out[i]
		if !reflect.DeepEqual(got, wanted) {
			t.Errorf("got %v, wanted %v", got, wanted)
		}
	}
}
