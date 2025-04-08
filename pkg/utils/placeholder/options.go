package placeholder

import "strconv"

type makeDollarsOptions func() string

// WithColumnNum - return placeholder for sql query
//
// example: WithColumnNum(2) -> ($1, $2)
func WithColumnNum(columnNum int) makeDollarsOptions {
	return func() string {
		return makeDollarsWithColumns(columnNum, 1)
	}
}

// WithColumnNumAndRowNum - return placeholder for sql query
//
// example: WithColumnNumAndRowNum(2, 3) -> ($1, $2), ($3, $4), ($5, $6)
func WithColumnNumAndRowNum(columnNum, rowNum int) makeDollarsOptions {
	return func() string {
		return makeDollarsWithColumns(columnNum, rowNum)
	}
}

// makeDollarsWithColumns - table for sql query
//
// example: makeDollarsWithColumns(3, 2) -> ($1, $2, $3), ($4, $5, $6)
func makeDollarsWithColumns(columnNum int, rowNum int) string {
	if columnNum < 1 || rowNum < 1 {
		return ""
	}

	placeholder := makeDollarsWithStartValue(columnNum, 1)
	for i := 1; i < rowNum; i++ {
		placeholder += "), (" + makeDollarsWithStartValue(columnNum, 1+i*columnNum)
	}

	return "(" + placeholder + ")"
}

// makeDollarsWithStartValue - return placeholder for sql query, where start value in changeable
//
// example: makeDollarsWithStartValue(4, 10) -> $10,$11,$12,$13
func makeDollarsWithStartValue(num int, startFrom int) string {
	if num < 1 || startFrom < 1 {
		return ""
	}

	placeholder := "$" + strconv.Itoa(startFrom)
	for i := startFrom + 1; i < startFrom+num; i++ {
		placeholder += ", $" + strconv.Itoa(i)
	}

	return placeholder
}
