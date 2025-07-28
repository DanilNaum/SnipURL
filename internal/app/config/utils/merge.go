package utils

// Merge returns the first non-nil pointer from the provided arguments.
// If all arguments are nil, it returns nil.
func Merge[T any](args ...*T) *T {
	for _, arg := range args {
		if arg != nil {
			return arg
		}
	}
	return nil
}
