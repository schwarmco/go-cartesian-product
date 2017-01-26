package cartesian_test

import (
    "fmt"
    "github.com/schwarmco/go-cartesian-product"
)

func ExampleIter() {
    a := []interface{}{1,2,3}
    b := []interface{}{"a","b","c"}

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
