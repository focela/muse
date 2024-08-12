// Copyright (c) 2024 Focela Technologies. All rights reserved.
// Use of this source code is governed by an MIT style
// license that can be found in the LICENSE file.

package hash

// RS implements the classic RS hash algorithm for 32-bit integers.
// This algorithm, developed by Robert Sedgewick, is efficient for hashing strings.
func RS(str []byte) uint32 {
	var (
		b    uint32 = 378551
		a    uint32 = 63689
		hash uint32 = 0
	)
	for _, ch := range str {
		hash = hash*a + uint32(ch)
		a *= b
	}
	return hash
}

// RS64 implements the classic RS hash algorithm for 64-bit integers.
// This version extends the RS algorithm to operate on 64-bit values.
func RS64(str []byte) uint64 {
	var (
		b    uint64 = 378551
		a    uint64 = 63689
		hash uint64 = 0
	)
	for _, ch := range str {
		hash = hash*a + uint64(ch)
		a *= b
	}
	return hash
}
