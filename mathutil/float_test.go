package mathutil

import "testing"

func TestRoundHalfEven(t *testing.T) {
	type args struct {
		value  float64
		places int
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "test_1",
			args: args{
				value:  1.23456,
				places: 1,
			},
			want: 1.2,
		},
		{
			name: "test_2",
			args: args{
				value:  1.23456,
				places: 2,
			},
			want: 1.23,
		},
		{
			name: "test_3",
			args: args{
				value:  1.23456,
				places: 3,
			},
			want: 1.235,
		},
		{
			name: "test_4",
			args: args{
				value:  1.23456,
				places: 4,
			},
			want: 1.2346,
		},
		{
			name: "test_5",
			args: args{
				value:  1.2,
				places: 2,
			},
			want: 1.2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RoundHalfEven(tt.args.value, tt.args.places); got != tt.want {
				t.Errorf("RoundHalfEven() = %v, want %v", got, tt.want)
			}
		})
	}
}
