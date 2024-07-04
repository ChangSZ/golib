package password

import (
	"reflect"
	"testing"
)

func TestPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		fake     string
		want     bool
	}{
		{
			name:     "test_1",
			password: "hello",
			want:     true,
		},
		{
			name:     "test_2",
			password: "hello",
			fake:     "xxxxxxx",
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoder := NewBcryptPasswordEncoder(10)
			encodedPassword, err := encoder.Encode(tt.password)
			if err != nil {
				t.Error(err)
			}
			if tt.fake != "" {
				encodedPassword = tt.fake
			}
			if match := encoder.Matches(tt.password, encodedPassword); !reflect.DeepEqual(match, tt.want) {
				t.Errorf("Matches() = %v, want %v", match, tt.want)
			}
		})
	}
}
