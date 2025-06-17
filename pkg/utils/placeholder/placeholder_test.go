package placeholder

import (
	"testing"
)

func TestPlaceholder_makeDollarsWithColumns(t *testing.T) {
	tests := []struct {
		name       string
		columnNum  int
		rowNum     int
		wantResult string
	}{
		{
			name:       "single row single column",
			columnNum:  1,
			rowNum:     1,
			wantResult: "($1)",
		},
		{
			name:       "single column multiple rows",
			columnNum:  1,
			rowNum:     3,
			wantResult: "($1), ($2), ($3)",
		},
		{
			name:       "single row multiple columns",
			columnNum:  3,
			rowNum:     1,
			wantResult: "($1, $2, $3)",
		},
		{
			name:       "multiple rows multiple columns",
			columnNum:  2,
			rowNum:     2,
			wantResult: "($1, $2), ($3, $4)",
		},
		{
			name:       "zero rows",
			columnNum:  2,
			rowNum:     0,
			wantResult: "",
		},
		{
			name:       "zero columns",
			columnNum:  0,
			rowNum:     2,
			wantResult: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := makeDollarsWithColumns(tt.columnNum, tt.rowNum)
			if got != tt.wantResult {
				t.Errorf("makeDollarsWithColumns() = %v, want %v", got, tt.wantResult)
			}
		})
	}
}

func TestPlaceholder_makeDollarsWithStartValue(t *testing.T) {
	tests := []struct {
		name       string
		num        int
		startFrom  int
		wantResult string
	}{
		{
			name:       "start from non-zero value",
			num:        3,
			startFrom:  5,
			wantResult: "$5, $6, $7",
		},
		{
			name:       "start from large number",
			num:        2,
			startFrom:  100,
			wantResult: "$100, $101",
		},
		{
			name:       "single placeholder with non-zero start",
			num:        1,
			startFrom:  10,
			wantResult: "$10",
		},
		{
			name:       "zero num with non-zero start",
			num:        0,
			startFrom:  5,
			wantResult: "",
		},
		{
			name:       "negative num with positive start",
			num:        -1,
			startFrom:  5,
			wantResult: "",
		},
		{
			name:       "positive num with negative start",
			num:        2,
			startFrom:  -1,
			wantResult: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := makeDollarsWithStartValue(tt.num, tt.startFrom)
			if got != tt.wantResult {
				t.Errorf("makeDollarsWithStartValue = %v, want %v", got, tt.wantResult)
			}
		})
	}
}
func TestPlaceholder_MakeDollars(t *testing.T) {
	tests := []struct {
		name       string
		opt        makeDollarsOptions
		wantResult string
	}{
		{
			name:       "WithColumnNum",
			opt:        WithColumnNum(2),
			wantResult: "($1, $2)",
		},
		{
			name:       "WithColumnNumAndRowNum",
			opt:        WithColumnNumAndRowNum(2, 3),
			wantResult: "($1, $2), ($3, $4), ($5, $6)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MakeDollars(tt.opt)
			if got != tt.wantResult {
				t.Errorf("MakeDollars() = %v, want %v", got, tt.wantResult)
			}
		})
	}
}
