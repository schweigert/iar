
func main() {
	var s []int
	printSlice(s)

	// append works on nil slices.
	s = append(s, 0)
	printSlice(s)

	// The slice grows as needed.
	s = append(s, 1)
	printSlice(s)

	// We can add more than one element at a time.
	s = append(s, 2, 3, 4)
	printSlice(s)
}

func printSlice(s []int) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}

// package main
//
// import (
// 	"fmt"
// )
//
// func main() {
// 	a := []int{1,2,3}
// 	var b int
//
// 	b, a = a[0], a[1:]
//
// 	fmt.Println(a, b)
//
// 	a = append(a, b)
//
// 	fmt.Println(a, b)
// [2 3] 1
// [2 3 1] 1
// }
