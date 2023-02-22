package dnb

import (
	"bytes"
	"errors"
	"testing"
)

func AssertArrays(t *testing.T,
	want []byte,
	have []byte) {
	if bytes.Compare(want, have) != 0 {
		t.Errorf("arrays are not equal,\nexpected= %v\n  actual= %v", want, have)
	}
}
func AssertValues(t *testing.T,
	want any,
	have any) {
	if want != have {
		t.Errorf("values are not equal,\nexpected= %v\n  actual= %v", want, have)
	}
}

func AssertError(t *testing.T, expected error, actual error) {
	if actual == nil {
		t.Errorf("Actual error is nil")
	}
	if expected == nil {
		t.Errorf("Expected error is nil")
	}
	AssertValues(t, expected.Error(), actual.Error())
}

func TestBinaryWriter_WriteBool(t *testing.T) {
	tests := []struct {
		name     string
		value    bool
		expected []byte
	}{
		{"true value", true, []byte{1}},
		{"false value", false, []byte{0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bw := NewBinaryWriter(bytes.NewBuffer([]byte{}))
			bw.WriteBool(tt.value)
			AssertArrays(t, tt.expected, bw.Bytes())
		})
	}
}

func TestBinaryWriter_WriteInt32(t *testing.T) {
	tests := []struct {
		name     string
		value    int32
		expected []byte
	}{
		{"0 value", 0, []byte{0, 0, 0, 0}},
		{"1000000 value", 1000000, []byte{64, 66, 15, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bw := NewBinaryWriter(bytes.NewBuffer([]byte{}))
			bw.WriteInt32(tt.value)
			AssertArrays(t, tt.expected, bw.Bytes())
		})
	}
}

func TestBinaryWriter_WriteString(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected []byte
	}{
		{"empty string", "", []byte{0}},
		{"test string", "test string", []byte{11, 116, 101, 115, 116, 32, 115, 116, 114, 105, 110, 103}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bw := NewBinaryWriter(bytes.NewBuffer([]byte{}))
			bw.WriteString(tt.value)
			AssertArrays(t, tt.expected, bw.Bytes())
		})
	}
}

func TestBinaryWriter_WriteBytes(t *testing.T) {
	tests := []struct {
		name     string
		value    []byte
		expected []byte
	}{
		{"empty bytes", []byte{}, []byte{0, 0, 0, 0}},
		{"test bytes", []byte{11, 116, 101, 115, 116, 32, 115, 116, 114, 105, 110, 103},
			[]byte{12, 0, 0, 0, 11, 116, 101, 115, 116, 32, 115, 116, 114, 105, 110, 103}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bw := NewBinaryWriter(bytes.NewBuffer([]byte{}))
			bw.WriteBytes(tt.value)
			AssertArrays(t, tt.expected, bw.Bytes())
		})
	}
}

func TestBinaryWriter_WriteByte(t *testing.T) {
	tests := []struct {
		name     string
		value    byte
		expected []byte
	}{
		{"byte 3 value", byte(3), []byte{3}},
		{"byte 0 value", byte(0), []byte{0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bw := NewBinaryWriter(bytes.NewBuffer([]byte{}))
			bw.WriteByte(tt.value)
			AssertArrays(t, tt.expected, bw.Bytes())
		})
	}
}

func TestBinaryWriter(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		expected []byte
	}{
		{"true value", true, []byte{1}},
		{"false value", false, []byte{0}},
		{"byte 3 value", byte(3), []byte{3}},
		{"byte 0 value", byte(0), []byte{0}},
		{"int 0 value", 0, []byte{0, 0, 0, 0}},
		{"int 1000000 value", 1000000, []byte{64, 66, 15, 0}},
		{"empty string", "", []byte{0}},
		{"test string", "test string", []byte{11, 116, 101, 115, 116, 32, 115, 116, 114, 105, 110, 103}},
		{"empty bytes", []byte{}, []byte{0, 0, 0, 0}},
		{"test bytes", []byte{11, 116, 101, 115, 116, 32, 115, 116, 114, 105, 110, 103},
			[]byte{12, 0, 0, 0, 11, 116, 101, 115, 116, 32, 115, 116, 114, 105, 110, 103}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bw := NewBinaryWriter(bytes.NewBuffer([]byte{}))
			bw.Write(tt.value)
			AssertArrays(t, tt.expected, bw.Bytes())
		})
	}
}

func TestBinaryReader_ReadBool(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
		value    []byte
		expErr   error
	}{
		{"true value", true, []byte{1}, nil},
		{"false value", false, []byte{0}, nil},
		{"test err", false, []byte{}, errors.New("not enough bytes, at least 1 byte is needed, array length = 0")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			br := NewBinaryReader(bytes.NewBuffer(tt.value))
			var data bool
			err := br.Read(&data)
			if tt.expErr == nil && err != nil {
				t.Errorf("Error = %v", err)
			}
			if tt.expErr != nil {
				AssertError(t, tt.expErr, err)
			} else {
				AssertValues(t, tt.expected, data)
			}
		})
	}
}

func TestBinaryReader_ReadInt32(t *testing.T) {
	tests := []struct {
		name     string
		expected int32
		value    []byte
		expErr   error
	}{
		{"int 0 value", 0, []byte{0, 0, 0, 0}, nil},
		{"int 1000000 value", 1000000, []byte{64, 66, 15, 0}, nil},
		{"test err", 1000000, []byte{64, 66, 15},
			errors.New("not enough bytes, at least 4 byte is needed, array length = 3")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			br := NewBinaryReader(bytes.NewBuffer(tt.value))
			var data int32
			err := br.Read(&data)
			if tt.expErr == nil && err != nil {
				t.Errorf("Error= %v", err)
			}
			if tt.expErr != nil {
				AssertError(t, tt.expErr, err)
			} else {
				AssertValues(t, tt.expected, data)
			}
		})
	}
}

func TestBinaryReader_ReadString(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		value    []byte
		expErr   error
	}{
		{"empty string", "", []byte{0}, nil},
		{"test err", "", []byte{},
			errors.New("not enough bytes, at least 1 byte is needed, array length = 0")},
		{"test string", "test string", []byte{11, 116, 101, 115, 116, 32, 115, 116, 114, 105, 110, 103}, nil},
		{"test err", "", []byte{11, 116, 101, 115, 116, 32, 115, 116, 114, 105, 110},
			errors.New("not enough bytes, at least 11 byte is needed, array length = 10")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			br := NewBinaryReader(bytes.NewBuffer(tt.value))
			var data string
			err := br.Read(&data)
			if tt.expErr == nil && err != nil {
				t.Errorf("Error= %v", err)
			}
			if tt.expErr != nil {
				AssertError(t, tt.expErr, err)
			} else {
				AssertValues(t, tt.expected, data)
			}
		})
	}
}

func TestBinaryReader_ReadByte(t *testing.T) {
	tests := []struct {
		name     string
		expected byte
		value    []byte
		expErr   error
	}{
		{"byte 3 value", byte(3), []byte{3}, nil},
		{"byte 0 value", byte(0), []byte{0}, nil},
		{"test err", byte(0), []byte{},
			errors.New("not enough bytes, at least 1 byte is needed, array length = 0")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			br := NewBinaryReader(bytes.NewBuffer(tt.value))
			var data byte
			err := br.Read(&data)
			if tt.expErr == nil && err != nil {
				t.Errorf("Error= %v", err)
			}
			if tt.expErr != nil {
				AssertError(t, tt.expErr, err)
			} else {
				AssertValues(t, tt.expected, data)
			}
		})
	}
}

func TestBinaryReader_ReadBytes(t *testing.T) {
	tests := []struct {
		name     string
		expected []byte
		value    []byte
		expErr   error
	}{
		{"empty bytes", []byte{}, []byte{0, 0, 0, 0}, nil},
		{"test err length", []byte{}, []byte{0, 0, 0},
			errors.New("not enough bytes, at least 4 byte is needed, array length = 3")},
		{"test bytes", []byte{11, 116, 101, 115, 116, 32, 115, 116, 114, 105, 110, 103},
			[]byte{12, 0, 0, 0, 11, 116, 101, 115, 116, 32, 115, 116, 114, 105, 110, 103}, nil},
		{"test err body", []byte{},
			[]byte{12, 0, 0, 0, 11, 116, 101, 115, 116, 32, 115, 116, 114, 105, 110},
			errors.New("not enough bytes, at least 12 byte is needed, array length = 11")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			br := NewBinaryReader(bytes.NewBuffer(tt.value))
			var data []byte
			err := br.Read(&data)
			if tt.expErr == nil && err != nil {
				t.Errorf("Error= %v", err)
			}
			if tt.expErr != nil {
				AssertError(t, tt.expErr, err)
			} else {
				AssertArrays(t, tt.expected, data)
			}
		})
	}
}

func TestBinaryReader_Read(t *testing.T) {
	tests := []struct {
		name     string
		expected [][]byte
		value    []byte
		expErr   error
	}{
		{"empty bytes", [][]byte{{}, {11, 116, 101, 115, 116, 32, 115, 116, 114, 105, 110, 103}},
			[]byte{0, 0, 0, 0, 12, 0, 0, 0, 11, 116, 101, 115, 116, 32, 115, 116, 114, 105, 110, 103}, nil},
		{"test bytes", [][]byte{{11, 116, 101, 115, 116, 32, 115, 116, 114, 105, 110, 103}, {}},
			[]byte{12, 0, 0, 0, 11, 116, 101, 115, 116, 32, 115, 116, 114, 105, 110, 103, 0, 0, 0, 0}, nil},
		{"test err", [][]byte{{}, {}},
			[]byte{},
			errors.New("not enough bytes, at least 4 byte is needed, array length = 0")},
		{"test err", [][]byte{{}, {}},
			[]byte{12, 0, 0, 0, 11, 116, 101, 115, 116, 32, 115, 116, 114, 105, 110, 103, 0, 0, 0},
			errors.New("not enough bytes, at least 4 byte is needed, array length = 3")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			br := NewBinaryReader(bytes.NewBuffer(tt.value))
			var data1 []byte
			var data2 []byte
			err := br.Read(&data1)
			if tt.expErr == nil && err != nil {
				t.Errorf("Error= %v", err)
			}
			err = br.Read(&data2)
			if tt.expErr == nil && err != nil {
				t.Errorf("Error= %v", err)
			}
			if tt.expErr != nil {
				AssertError(t, tt.expErr, err)
			} else {
				AssertArrays(t, tt.expected[0], data1)
				AssertArrays(t, tt.expected[1], data2)
			}
		})
	}
}
