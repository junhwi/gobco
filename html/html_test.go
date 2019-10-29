package html

import "testing"

func Test_getNext(t *testing.T) {
	type args struct {
		s   string
		end int
	}
	tests := []struct {
		name      string
		args      args
		wantValue int
		wantNext  int
		wantErr   bool
	}{
		{
			"success",
			args{"123,456", 7},
			456,
			3,
			false,
		},
		{
			"stop",
			args{",123", 4},
			123,
			0,
			false,
		},
		{
			"NaN",
			args{",abc", 4},
			0,
			0,
			true,
		},
		{
			"empty",
			args{"123,", 4},
			0,
			0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotNext, err := getNext(tt.args.s, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("getNext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotValue != tt.wantValue {
				t.Errorf("getNext() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
			if gotNext != tt.wantNext {
				t.Errorf("getNext() gotNext = %v, want %v", gotNext, tt.wantNext)
			}
		})
	}
}
