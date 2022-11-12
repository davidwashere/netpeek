package util

import "testing"

func TestConvertOctalStringToUint32(t *testing.T) {
	tt := []struct {
		in   string
		want uint32
	}{
		{"644", 420},
		{"0644", 420},
		{"777", 511},
		{"", 0},
		{"9999999999999999999999999999999999999999999999", 0},
		{"hello", 0},
	}

	for _, test := range tt {
		got := ConvertOctalStringToUint32(test.in)

		if got != test.want {
			t.Errorf("got %v, want %v", got, test.want)
		}
	}
}
