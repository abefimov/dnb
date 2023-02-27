package dnb

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type BinaryWriter struct {
	stream *bytes.Buffer
}

func NewBinaryWriter(stream *bytes.Buffer) *BinaryWriter {
	return &BinaryWriter{stream: stream}
}

func (bw *BinaryWriter) WriteInt32(value int32) {
	_ = binary.Write(bw.stream, binary.LittleEndian, value)
}

func (bw *BinaryWriter) WriteInt64(value int64) {
	_ = binary.Write(bw.stream, binary.LittleEndian, value)
}

func (bw *BinaryWriter) WriteString(value string) {
	leb := GetLeb128Bytes(len(value))
	_, _ = bw.stream.Write(leb)
	_, _ = bw.stream.Write([]byte(value))
}

// WriteBool writes a bool value to the stream.
func (bw *BinaryWriter) WriteBool(value bool) {
	_ = binary.Write(bw.stream, binary.LittleEndian, value)
}

func (bw *BinaryWriter) WriteBytes(value []byte) {
	bw.WriteInt32(int32(len(value)))
	bw.stream.Write(value)
}

func (bw *BinaryWriter) WriteByte(value byte) {
	bw.stream.Write([]byte{value})
}

func (bw *BinaryWriter) Write(value any) {
	switch v := value.(type) {
	case *bool:
		bw.WriteBool(*v)
	case bool:
		bw.WriteBool(v)
	case *int:
		bw.WriteInt32(int32(*v))
	case *int16:
		bw.WriteInt32(int32(*v))
	case *int32:
		bw.WriteInt32(*v)
	case *int64:
		bw.WriteInt64(*v)
	case int:
		bw.WriteInt32(int32(v))
	case int16:
		bw.WriteInt32(int32(v))
	case int32:
		bw.WriteInt32(v)
	case int64:
		bw.WriteInt64(v)
	case *string:
		bw.WriteString(*v)
	case string:
		bw.WriteString(v)
	case *byte:
		bw.WriteByte(*v)
	case byte:
		bw.WriteByte(v)
	case *[]byte:
		bw.WriteBytes(*v)
	case []byte:
		bw.WriteBytes(v)
	}
}

func (bw *BinaryWriter) Bytes() []byte {
	return bw.stream.Bytes()
}

type BinaryReader struct {
	stream *bytes.Buffer
}

func FromBytes(b []byte) *BinaryReader {
	return NewBinaryReader(bytes.NewBuffer(b))
}
func NewBinaryReader(stream *bytes.Buffer) *BinaryReader {
	return &BinaryReader{stream: stream}
}

func (br *BinaryReader) ReadInt32() (int32, error) {
	if br.stream.Len() < 4 {
		return 0, fmt.Errorf("not enough bytes, at least %v byte is needed, array length = %v", 4, br.stream.Len())
	}
	var p int32
	_ = binary.Read(br.stream, binary.LittleEndian, &p)
	return p, nil
}

func (br *BinaryReader) ReadString() (string, error) {
	if br.stream.Len() < 1 {
		return "", fmt.Errorf("not enough bytes, at least %v byte is needed, array length = %v", 1, br.stream.Len())
	}
	length := FromLeb128Bytes(br.stream.Next(1))
	if br.stream.Len() < length {
		return "", fmt.Errorf("not enough bytes, at least %v byte is needed, array length = %v", length, br.stream.Len())
	}
	return string(br.stream.Next(length)[:]), nil
}

func (br *BinaryReader) ReadBool() (bool, error) {
	if br.stream.Len() < 1 {
		return false, fmt.Errorf("not enough bytes, at least %v byte is needed, array length = %v", 1, br.stream.Len())
	}
	i, err := br.stream.ReadByte()
	if err != nil {
		return false, err
	}
	if i == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (br *BinaryReader) ReadBytes() ([]byte, error) {
	length, err := br.ReadInt32()
	if err != nil {
		return nil, err
	}
	l := int(length)
	if br.stream.Len() < l {
		return nil, fmt.Errorf("not enough bytes, at least %v byte is needed, array length = %v", l, br.stream.Len())
	}
	return br.stream.Next(l), nil
}

func (br *BinaryReader) ReadOneByte() (byte, error) {
	if br.stream.Len() < 1 {
		return 0, fmt.Errorf("not enough bytes, at least %v byte is needed, array length = %v", 1, br.stream.Len())
	}
	b, err := br.stream.ReadByte()
	if err != nil {
		return 0, err
	}
	return b, nil
}

func (br *BinaryReader) Read(data any) error {
	switch data := data.(type) {
	case *bool:
		v, err := br.ReadBool()
		if err != nil {
			return err
		}
		*data = v
	case *int:
		v, err := br.ReadInt32()
		if err != nil {
			return err
		}
		*data = int(v)
	case *int16:
		v, err := br.ReadInt32()
		if err != nil {
			return err
		}
		*data = int16(v)
	case *int32:
		v, err := br.ReadInt32()
		if err != nil {
			return err
		}
		*data = v
	case *string:
		v, err := br.ReadString()
		if err != nil {
			return err
		}
		*data = v
	case *byte:
		v, err := br.ReadOneByte()
		if err != nil {
			return err
		}
		*data = v
	case *[]byte:
		v, err := br.ReadBytes()
		if err != nil {
			return err
		}
		*data = v
	}
	return nil
}

func (br *BinaryReader) Len() int {
	return br.stream.Len()
}

func (br *BinaryReader) Cap() int {
	return br.stream.Cap()
}

func FromLeb128Bytes(l []byte) int {
	var n uint
	for i := 0; i < len(l); i++ {
		b := uint(0x7F & l[i])
		n |= b << (i * 7)
		if b := l[i]; b&0x80 == 0 && b&0x40 != 0 {
			return int(n) | (^0 << ((i + 1) * 7))
		}
	}
	return int(n)
}

func GetLeb128Bytes(n int) []byte {
	leb := make([]byte, 0)
	for {
		var (
			b    = byte(n & 0x7F)
			sign = byte(n & 0x40)
		)
		if n >>= 7; sign == 0 && n != 0 || n != -1 && (n != 0 || sign != 0) {
			b |= 0x80
		}
		leb = append(leb, b)
		if b&0x80 == 0 {
			break
		}
	}
	return leb
}
