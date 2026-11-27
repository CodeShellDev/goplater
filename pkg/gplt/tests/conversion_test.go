package tests

import (
	"strconv"
	"testing"
)

func createFuncTestTmpl(name, tmpl string) string {
	return `{{ $output := ` + tmpl + ` }} {{ emit ` + strconv.Quote(name) + ` $output }}`
}

func TestToString(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  string
	}{
		{
			name:  "string",
			input: "hello",
			want:  "hello",
		},
		{
			name:  "bytes",
			input: []byte("hello"),
			want:  "hello",
		},
		{
			name:  "integer",
			input: 42,
			want:  "42",
		},
		{
			name:  "float",
			input: 3.14,
			want:  "3.14",
		},
		{
			name:  "boolean",
			input: true,
			want:  "true",
		},
		{
			name:  "nil",
			input: nil,
			want:  "<nil>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl := createFuncTestTmpl(tt.name, "toString .")

			trt := tester.Collect(t, tmpl, tt.input)

			trt.Equal(t, tt.name, tt.want)
		})
	}
}

func TestToBytes(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []byte
	}{
		{
			name:  "plain string",
			input: "hello",
			want:  []byte("hello"),
		},
		{
			name:  "empty string",
			input: "",
			want:  []byte{},
		},
		{
			name:  "unicode string",
			input: "こんにちは",
			want:  []byte("こんにちは"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl := createFuncTestTmpl(tt.name, "toBytes .")

			trt := tester.Collect(t, tmpl, tt.input)

			trt.Equal(t, tt.name, tt.want)
		})
	}
}

func TestToInt(t *testing.T) {
	tests := []struct {
		name	string
		input	string
		want	int
		wantErr bool
	}{
		{
			name:  "positive integer",
			input: "42",
			want:  42,
		},
		{
			name:  "negative integer",
			input: "-42",
			want:  -42,
		},
		{
			name:  "zero",
			input: "0",
			want:  0,
		},
		{
			name:	 "decimal",
			input:	 "3.14",
			wantErr: true,
		},
		{
			name:	 "invalid string",
			input:	 "hello",
			wantErr: true,
		},
		{
			name:	 "empty string",
			input:	 "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl := createFuncTestTmpl(tt.name, "toInt .")

			if tt.wantErr {
				tester.ExpectErr(t, tmpl, nil)
				return
			}

			trt := tester.Collect(t, tmpl, tt.input)

			trt.Equal(t, tt.name, tt.want)
		})
	}
}

func TestToFloat64(t *testing.T) {
	tests := []struct {
		name	string
		input	string
		want	float64
		wantErr bool
	}{
		{
			name:  "positive float",
			input: "3.14",
			want:  3.14,
		},
		{
			name:  "negative float",
			input: "-3.14",
			want:  -3.14,
		},
		{
			name:  "integer",
			input: "42",
			want:  42,
		},
		{
			name:  "scientific notation",
			input: "1.5e2",
			want:  150,
		},
		{
			name:	 "invalid string",
			input:	 "hello",
			wantErr: true,
		},
		{
			name:	 "empty string",
			input:	 "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl := createFuncTestTmpl(tt.name, "toFloat64 .")

			if tt.wantErr {
				tester.ExpectErr(t, tmpl, nil)
				return
			}

			trt := tester.Collect(t, tmpl, tt.input)

			trt.Equal(t, tt.name, tt.want)
		})
	}
}

func TestToFloat32(t *testing.T) {
	tests := []struct {
		name	string
		input	string
		want	float32
		wantErr bool
	}{
		{
			name:  "positive float",
			input: "3.14",
			want:  3.14,
		},
		{
			name:  "negative float",
			input: "-3.14",
			want:  -3.14,
		},
		{
			name:  "integer",
			input: "42",
			want:  42,
		},
		{
			name:  "scientific notation",
			input: "1.5e2",
			want:  150,
		},
		{
			name:	 "invalid string",
			input:	 "hello",
			wantErr: true,
		},
		{
			name:	 "empty string",
			input:	 "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl := createFuncTestTmpl(tt.name, "toFloat32 .")

			if tt.wantErr {
				tester.ExpectErr(t, tmpl, nil)
				return
			}

			trt := tester.Collect(t, tmpl, tt.input)

			trt.Equal(t, tt.name, tt.want)
		})
	}
}

func TestToBool(t *testing.T) {
	tests := []struct {
		name	string
		input	string
		want	bool
		wantErr bool
	}{
		{
			name:  "true",
			input: "true",
			want:  true,
		},
		{
			name:  "false",
			input: "false",
			want:  false,
		},
		{
			name:  "one",
			input: "1",
			want:  true,
		},
		{
			name:  "zero",
			input: "0",
			want:  false,
		},
		{
			name:  "short true",
			input: "T",
			want:  true,
		},
		{
			name:  "short false",
			input: "F",
			want:  false,
		},
		{
			name:	 "invalid string",
			input:	 "yes",
			wantErr: true,
		},
		{
			name:	 "empty string",
			input:	 "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl := createFuncTestTmpl(tt.name, "toBool .")

			if tt.wantErr {
				tester.ExpectErr(t, tmpl, nil)
				return
			}

			trt := tester.Collect(t, tmpl, tt.input)

			trt.Equal(t, tt.name, tt.want)
		})
	}
}