package ch_encoding

import (
	"math"
	"reflect"
	"testing"
)

func TestCustomEncoder_Decode(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    Data
		wantErr bool
	}{
		{
			name: "decode golden path",
			args: args{
				s: "\u0002\u0000\u0000\u0000\u0005hello\u0001\u0000\u0000\u0000\u0007",
			},
			want: Data{
				"hello",
				int32(7),
			},
			wantErr: false,
		},
		{
			name: "decode int data",
			args: args{
				s: "\u0001\u0000\u0000\u0000\u0008",
			},
			want: Data{
				int32(8),
			},
			wantErr: false,
		},
		{
			name: "test nested",
			args: args{
				s: "\u0002\u0000\u0000\u0000\u0005hello\u0001\u0000\u0000\u0000\u0005\u0003\u0000\u0000\u0000\u0012\u0002\u0000\u0000\u0000\bgood bye\u0001\u0000\u0000\u0000\t",
			},
			want: Data{
				"hello",
				int32(5),
				Data{
					"good bye",
					int32(9),
				},
			},
			wantErr: false,
		},
		{
			name: "test deeply nested",
			args: args{
				s: "\u0002\u0000\u0000\u0000\u0003one\u0001\u0000\u0000\u0000\u0001\u0003\u0000\u0000\u0000I\u0002\u0000\u0000\u0000\u0003two\u0001\u0000\u0000\u0000\u0002\u0003\u0000\u0000\u00007\u0002\u0000\u0000\u0000\u0005three\u0001\u0000\u0000\u0000\u0003\u0003\u0000\u0000\u0000#\u0002\u0000\u0000\u0000\u0004four\u0001\u0000\u0000\u0000\u0004\u0003\u0000\u0000\u0000\u0010\u0002\u0000\u0000\u0000\u0006five e\u0001\u0000\u0000\u0000\u0005",
			},
			want: Data{
				"one",
				int32(1),
				Data{
					"two",
					int32(2),
					Data{
						"three",
						int32(3),
						Data{
							"four",
							int32(4),
							Data{
								"five e",
								int32(5),
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "test empty",
			args: args{
				s: "",
			},
			want:    Data{},
			wantErr: false,
		},
		{
			name: "test nested empty",
			args: args{
				s: "\u0003\u0000\u0000\u0000\u0000",
			},
			want:    Data{Data{}},
			wantErr: false,
		},
		{
			name: "test several un-nested empty",
			args: args{
				s: "\u0003\u0000\u0000\u0000\b\u0002\u0000\u0000\u0000\u0003one\u0003\u0000\u0000\u0000\b\u0002\u0000\u0000\u0000\u0003two\u0003\u0000\u0000\u0000\n\u0002\u0000\u0000\u0000\u0005three",
			},
			want: Data{
				Data{"one"},
				Data{"two"},
				Data{"three"},
			},
			wantErr: false,
		},
		{
			name: "invalid code point",
			args: args{
				s: "\u0009",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &CustomEncoder{}
			got, err := e.Decode(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCustomEncoder_Encode(t *testing.T) {
	type args struct {
		data Data
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test string",
			args: args{
				data: Data{
					"hello",
				},
			},
			want:    "\u0002\u0000\u0000\u0000\u0005hello",
			wantErr: false,
		},
		{
			name: "test int",
			args: args{
				data: Data{
					int32(7),
				},
			},
			want:    "\u0001\u0000\u0000\u0000\u0007",
			wantErr: false,
		},
		{
			name: "test nested",
			args: args{
				data: Data{
					"hello",
					int32(5),
					Data{
						"good bye",
						int32(9),
					},
				},
			},
			want:    "\u0002\u0000\u0000\u0000\u0005hello\u0001\u0000\u0000\u0000\u0005\u0003\u0000\u0000\u0000\u0012\u0002\u0000\u0000\u0000\bgood bye\u0001\u0000\u0000\u0000\t",
			wantErr: false,
		},
		{
			name: "test deep nested",
			args: args{
				data: Data{
					"one",
					int32(1),
					Data{
						"two",
						int32(2),
						Data{
							"three",
							int32(3),
							Data{
								"four",
								int32(4),
								Data{
									"five e",
									int32(5),
								},
							},
						},
					},
				},
			},
			want:    "\u0002\u0000\u0000\u0000\u0003one\u0001\u0000\u0000\u0000\u0001\u0003\u0000\u0000\u0000I\u0002\u0000\u0000\u0000\u0003two\u0001\u0000\u0000\u0000\u0002\u0003\u0000\u0000\u00007\u0002\u0000\u0000\u0000\u0005three\u0001\u0000\u0000\u0000\u0003\u0003\u0000\u0000\u0000#\u0002\u0000\u0000\u0000\u0004four\u0001\u0000\u0000\u0000\u0004\u0003\u0000\u0000\u0000\u0010\u0002\u0000\u0000\u0000\u0006five e\u0001\u0000\u0000\u0000\u0005",
			wantErr: false,
		},
		{
			name: "test empty",
			args: args{
				data: Data{},
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "test nested empty",
			args: args{
				data: Data{
					Data{},
				}},
			want:    "\u0003\u0000\u0000\u0000\u0000",
			wantErr: false,
		},
		{
			name: "test several un-nested array",
			args: args{
				data: Data{
					Data{"one"},
					Data{"two"},
					Data{"three"},
				},
			},
			want:    "\u0003\u0000\u0000\u0000\b\u0002\u0000\u0000\u0000\u0003one\u0003\u0000\u0000\u0000\b\u0002\u0000\u0000\u0000\u0003two\u0003\u0000\u0000\u0000\n\u0002\u0000\u0000\u0000\u0005three",
			wantErr: false,
		},
		{
			name: "unsupported type",
			args: args{
				data: Data{uint16(666)},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "test large number",
			args: args{
				data: Data{int32(math.MaxInt32)},
			},
			// easier to write in hex
			want:    string([]byte{IntCodepoint, 0x7F, 0xFF, 0xFF, 0xFF}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &CustomEncoder{}
			got, err := e.Encode(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Encode() got = %v, want %v", got, tt.want)
			}
		})
	}
}
