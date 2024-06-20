package sign

import "testing"

func TestSign(t *testing.T) {
	type args struct {
		b   []byte
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Good case",
			args: args{
				b:   []byte("body"),
				key: "secret",
			},
			want:    "e6b1b0f6bfd1318bb4825bd1723d88de124d15abad8c3731fb956684f2f063b6",
			wantErr: false,
		},
		{
			name: "Nil signing",
			args: args{
				b:   nil,
				key: "secret",
			},
			want:    "2bb80d537b1da3e38bd30361aa855686bde0eacd7162fef6a25fe97bf527a25b",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sign(tt.args.b, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Sign() = %v, want %v", got, tt.want)
			}
		})
	}
}
