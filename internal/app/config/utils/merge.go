package utils

func Merge[T any](args ...*T) *T {
	for _, arg := range args {
		if arg != nil {
			return arg
		}
	}
	return nil
}
