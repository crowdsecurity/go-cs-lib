package slicetools

import (
	"context"
	"slices"
)

// a simple wrapper around slices.Chunk
// if you want context cancelation and don't need parallelism
//
// also: doesn't panic for size = 0

// Batch applies fn to successive chunks of at most "size" elements.
// A size of 0 (or negative) processes all the elements in one chunk.
// Stops at the first error and returns it.
func Batch[T any](ctx context.Context, elems []T, size int, fn func(context.Context, []T) error) error {
	n := len(elems)

	if n == 0 {
		return nil
	}

	if size <= 0 || size > n {
		size = n
	}

	// delegate to stdlib

	for part := range slices.Chunk(elems, size) {
		if err := ctx.Err(); err != nil {
			return err
		}

		if err := fn(ctx, part); err != nil {
			return err
		}
	}

	// we have stdlib at home

	//	for start := 0; start < n; start += size {
	//		if ctx.Err() != nil {
	//			return ctx.Err()
	//		}
	//
	//		end := min(start+size, n)
	//		if err := fn(ctx, elems[start:end]); err != nil {
	//			return err
	//		}
	//	}

	return nil
}
