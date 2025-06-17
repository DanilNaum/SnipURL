package placeholder

// MakeDollars - return placeholder for sql query
//
// example: MakeDollars(WithColumnNum(2)) -> ($1, $2)
//
// example: MakeDollars(WithColumnNumAndRowNum(2, 3)) -> ($1, $2), ($3, $4), ($5, $6)
//
// options:
//
// [WithColumnNum] - set  number of columns
//
// [WithColumnNumAndRowNum] - set number of columns and rows
func MakeDollars(opt makeDollarsOptions) string {
	return opt()
}
