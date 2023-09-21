package cartesian_test

import (
	"fmt"
	"math/rand"
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

func TestIterStress(t *testing.T) {
	const max = 6
	rng := rand.New(rand.NewSource(2))
	for n := 1; n < max; n++ {
		for m := 1; m < max; m++ {
			for x := 1; x < max; x++ {
				for y := 1; y < max; y++ {
					for z := 1; z < max; z++ {
						p0 := randomSlice(rng, n)
						p1 := randomSlice(rng, m)
						p2 := randomSlice(rng, x)
						p3 := randomSlice(rng, y)
						p4 := randomSlice(rng, z)

						wantHashes := make(map[uint32]struct{})
						for in := 0; in < n; in++ {
							n0 := p0[in].(uint32)
							for im := 0; im < m; im++ {
								n1 := p1[im].(uint32)
								for ix := 0; ix < x; ix++ {
									n2 := p2[ix].(uint32)
									for iy := 0; iy < y; iy++ {
										n3 := p3[iy].(uint32)
										for iz := 0; iz < z; iz++ {
											n4 := p4[iz].(uint32)
											sum := n0 + n1 + n2 + n3 + n4
											wantHashes[sum] = struct{}{}
										}
									}
								}
							}
						}

						// We then check that we got that hash for every product we received.
						for product := range cartesian.Iter(p0, p1, p2, p3, p4) {
							sum := uint32(0)
							for i := range product {
								sum += product[i].(uint32)
							}
							if _, ok := wantHashes[sum]; !ok {
								t.Error("unexpected product:", product)
							}
						}
					}
				}
			}
		}
	}
}

func randomSlice(rng *rand.Rand, n int) []interface{} {
	s := make([]interface{}, n)
	for i := 0; i < n; i++ {
		s[i] = rng.Uint32()
	}
	return s
}
