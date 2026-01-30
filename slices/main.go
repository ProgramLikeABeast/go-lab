package main

import (
	"fmt"
	"slices"
)

// =============================================================================
// ARRAYS vs SLICES: Arrays are fixed, slices are dynamic views
// =============================================================================

func demoArraysVsSlices() {
	fmt.Println("\n=== ARRAYS vs SLICES ===")

	// Array: fixed size, value type (copying copies all elements)
	arr := [5]int{1, 2, 3, 4, 5}
	fmt.Printf("Array: %v (type: %T)\n", arr, arr)

	// Slice: dynamic, reference type (header: ptr, len, cap)
	slice := []int{1, 2, 3, 4, 5}
	fmt.Printf("Slice: %v (type: %T)\n", slice, slice)

	// Arrays are values - assignment copies
	arr2 := arr
	arr2[0] = 999
	fmt.Printf("After arr2[0]=999: arr=%v, arr2=%v\n", arr, arr2)

	// Slices share underlying array
	slice2 := slice
	slice2[0] = 999
	fmt.Printf("After slice2[0]=999: slice=%v, slice2=%v\n", slice, slice2)
}

// =============================================================================
// SLICE INTERNALS: Length, Capacity, and the underlying array
// =============================================================================

func demoSliceInternals() {
	fmt.Println("\n=== SLICE INTERNALS ===")

	// make([]T, length, capacity)
	s := make([]int, 3, 5)
	fmt.Printf("make([]int, 3, 5): %v, len=%d, cap=%d\n", s, len(s), cap(s))

	// Slice header structure (conceptually):
	// type slice struct {
	//     ptr *T   // pointer to underlying array
	//     len int  // number of elements
	//     cap int  // capacity (can grow up to this without reallocation)
	// }

	// Capacity grows when exceeded
	s = append(s, 1, 2) // len=5, cap=5
	fmt.Printf("After append 2 elements: %v, len=%d, cap=%d\n", s, len(s), cap(s))

	s = append(s, 3) // len=6, cap grows (typically doubles)
	fmt.Printf("After exceeding cap: %v, len=%d, cap=%d\n", s, len(s), cap(s))
}

// =============================================================================
// SLICING: Creating sub-slices (they share memory!)
// =============================================================================

func demoSlicing() {
	fmt.Println("\n=== SLICING SYNTAX ===")

	s := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	fmt.Println("Original:", s)
	fmt.Println("s[2:5] =", s[2:5])   // [2 3 4] - indices 2,3,4
	fmt.Println("s[:3] =", s[:3])     // [0 1 2] - first 3
	fmt.Println("s[7:] =", s[7:])     // [7 8 9] - from index 7 to end
	fmt.Println("s[:] =", s[:])       // [0 1 2 ... 9] - copy of slice header

	// GOTCHA: Sub-slices share the underlying array!
	fmt.Println("\n=== SHARED MEMORY GOTCHA ===")
	original := []int{1, 2, 3, 4, 5}
	sub := original[1:3] // [2, 3]

	fmt.Printf("original: %v\n", original)
	fmt.Printf("sub: %v (cap=%d)\n", sub, cap(sub))

	sub[0] = 999 // Modifies original too!
	fmt.Printf("After sub[0]=999: original=%v, sub=%v\n", original, sub)

	// Appending to sub can overwrite original's elements
	sub = append(sub, 888) // This overwrites original[3]!
	fmt.Printf("After append to sub: original=%v, sub=%v\n", original, sub)
}

// =============================================================================
// FULL SLICE EXPRESSION: Limit capacity to prevent overwrites
// =============================================================================

func demoFullSliceExpression() {
	fmt.Println("\n=== FULL SLICE EXPRESSION s[low:high:max] ===")

	original := []int{1, 2, 3, 4, 5}
	
	// s[low:high:max] - max limits capacity
	sub := original[1:3:3] // len=2, cap=2 (not 4!)

	fmt.Printf("original: %v\n", original)
	fmt.Printf("sub with limited cap: %v (cap=%d)\n", sub, cap(sub))

	// Now append will allocate new array instead of overwriting
	sub = append(sub, 888)
	fmt.Printf("After append: original=%v, sub=%v\n", original, sub)
}

// =============================================================================
// COPYING SLICES: How to truly copy
// =============================================================================

func demoCopying() {
	fmt.Println("\n=== COPYING SLICES ===")

	src := []int{1, 2, 3, 4, 5}

	// Method 1: copy() built-in
	dst1 := make([]int, len(src))
	copy(dst1, src)

	// Method 2: append to nil slice
	dst2 := append([]int(nil), src...)

	// Method 3: slices.Clone (Go 1.21+)
	dst3 := slices.Clone(src)

	src[0] = 999
	fmt.Println("src:", src)
	fmt.Println("dst1 (copy):", dst1)
	fmt.Println("dst2 (append):", dst2)
	fmt.Println("dst3 (Clone):", dst3)
}

// =============================================================================
// SLICE OPERATIONS
// =============================================================================

func demoOperations() {
	fmt.Println("\n=== COMMON SLICE OPERATIONS ===")

	s := []int{1, 2, 3, 4, 5}

	// Delete element at index 2
	i := 2
	deleted := append(s[:i], s[i+1:]...)
	fmt.Println("After deleting index 2:", deleted)

	// Insert element at index 2
	s = []int{1, 2, 4, 5}
	i = 2
	s = append(s[:i], append([]int{3}, s[i:]...)...)
	fmt.Println("After inserting 3 at index 2:", s)

	// Using slices package (Go 1.21+)
	fmt.Println("\n=== SLICES PACKAGE (Go 1.21+) ===")

	nums := []int{3, 1, 4, 1, 5, 9, 2, 6}
	fmt.Println("Original:", nums)

	sorted := slices.Clone(nums)
	slices.Sort(sorted)
	fmt.Println("Sorted:", sorted)

	fmt.Println("Contains 5?", slices.Contains(nums, 5))
	fmt.Println("Index of 5:", slices.Index(nums, 5))
	fmt.Println("Max:", slices.Max(nums))
	fmt.Println("Min:", slices.Min(nums))
}

// =============================================================================
// NIL vs EMPTY SLICE
// =============================================================================

func demoNilVsEmpty() {
	fmt.Println("\n=== NIL vs EMPTY SLICE ===")

	var nilSlice []int          // nil slice
	emptySlice := []int{}       // empty slice (allocated)
	makeEmpty := make([]int, 0) // also empty slice

	fmt.Printf("nilSlice: %v, len=%d, cap=%d, nil=%t\n",
		nilSlice, len(nilSlice), cap(nilSlice), nilSlice == nil)
	fmt.Printf("emptySlice: %v, len=%d, cap=%d, nil=%t\n",
		emptySlice, len(emptySlice), cap(emptySlice), emptySlice == nil)
	fmt.Printf("makeEmpty: %v, len=%d, cap=%d, nil=%t\n",
		makeEmpty, len(makeEmpty), cap(makeEmpty), makeEmpty == nil)

	// Both work with append, len, range, etc.
	nilSlice = append(nilSlice, 1)
	emptySlice = append(emptySlice, 1)
	fmt.Println("\nAfter append:")
	fmt.Println("nilSlice:", nilSlice)
	fmt.Println("emptySlice:", emptySlice)

	// Prefer nil slice as zero value; use empty only when marshaling to JSON []
}

// =============================================================================
// MAPS: Key-value storage
// =============================================================================

func demoMaps() {
	fmt.Println("\n=== MAPS ===")

	// Creation
	m := map[string]int{
		"alice": 30,
		"bob":   25,
	}
	fmt.Println("Initial map:", m)

	// Add/Update
	m["charlie"] = 35
	m["alice"] = 31
	fmt.Println("After add/update:", m)

	// Access with ok pattern (distinguish zero value from missing)
	age, ok := m["david"]
	if !ok {
		fmt.Println("david not found, age would be:", age)
	}

	// Delete
	delete(m, "bob")
	fmt.Println("After delete bob:", m)

	// Iteration (order is random!)
	fmt.Println("Iterating:")
	for name, age := range m {
		fmt.Printf("  %s: %d\n", name, age)
	}

	// Nil map gotcha
	fmt.Println("\n=== NIL MAP GOTCHA ===")
	var nilMap map[string]int
	fmt.Println("Read from nil map:", nilMap["key"]) // OK, returns zero
	// nilMap["key"] = 1 // PANIC! Can't write to nil map
	fmt.Println("Writing to nil map would panic!")
}

func main() {
	fmt.Println("========================================")
	fmt.Println("GO SLICES AND MAPS")
	fmt.Println("========================================")

	demoArraysVsSlices()
	demoSliceInternals()
	demoSlicing()
	demoFullSliceExpression()
	demoCopying()
	demoOperations()
	demoNilVsEmpty()
	demoMaps()
}
