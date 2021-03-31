package multidim_test

import (
	"fmt"

	"github.com/dolmen-go/multidim"
)

func ExampleInit_noDimension() {
	var n int
	multidim.Init(&n, nil)
	fmt.Println(n)

	multidim.Init(&n, -1)
	fmt.Println(n)

	var s string
	multidim.Init(&s, "a")
	fmt.Println(s)

	multidim.Init(&s, func() string {
		return "x"
	})
	fmt.Println(s)

	multidim.Init(&s, func(ps *string) {
		*ps = "y"
	})
	fmt.Println(s)

	// Output:
	// 0
	// -1
	// a
	// x
	// y
}

func ExampleInit_square() {
	var square [][]int

	multidim.Init(&square, nil, 2, 2)
	fmt.Println(square)

	square = nil
	multidim.Init(&square, 4, 2, 2)
	fmt.Println(square)

	// Output:
	// [[0 0] [0 0]]
	// [[4 4] [4 4]]
}

func ExampleInit_cube() {
	var cube [][][]int
	multidim.Init(&cube, 8, 2, 2, 2)

	fmt.Println(cube)
	// Output:
	// [[[8 8] [8 8]] [[8 8] [8 8]]]
}
