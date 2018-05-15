package cartesian_test

import (
	"fmt"
	"testing"

	"github.com/schwarmco/go-cartesian-product"
)

func ExampleIter() {
	a := []interface{}{1, 2, 3}
	b := []interface{}{"a", "b", "c"}

	c := cartesian.Iter(a, b)

	// receive products through channel
	for product := range c {
		fmt.Println(product)
	}

	// Unordered Output:
	// [1 c]
	// [2 c]
	// [3 c]
	// [1 a]
	// [1 b]
	// [2 a]
	// [2 b]
	// [3 a]
	// [3 b]
}

func TestIter(t *testing.T) {

	// the sum on each index should be ( (1+2+3)/3 ) * 3 ^ 4
	// meaning that the mean (2) should occur in every line (which are 81 in total)
	var expected = 162
	var cnt0, cnt1, cnt2, cnt3 int

	a := []interface{}{1, 2, 3}
	c := cartesian.Iter(a, a, a, a)

	for product := range c {
		cnt0 += product[0].(int)
		cnt1 += product[1].(int)
		cnt2 += product[2].(int)
		cnt3 += product[3].(int)
	}

	if cnt0 != expected || cnt1 != expected || cnt2 != expected || cnt3 != expected {
		t.Error("expected counter to be", expected, "got:", cnt0, cnt1, cnt2, cnt3)
	}
}
