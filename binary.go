package dnb

import (
	"bytes"
	"encoding/binary"
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
	case int:
		bw.WriteInt32(int32(v))
	case int16:
		bw.WriteInt32(int32(v))
	case int32:
		bw.WriteInt32(v)
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

func (br *BinaryReader) ReadInt32() int32 {
	var p int32
	_ = binary.Read(br.stream, binary.LittleEndian, &p)
	return p
}

func (br *BinaryReader) ReadString() string {
	length := FromLeb128Bytes(br.stream.Next(1))
	return string(br.stream.Next(length)[:])
}

func (br *BinaryReader) ReadBool() bool {
	i := br.stream.Next(1)
	if i[0] == 0 {
		return false
	} else {
		return true
	}
}

func (br *BinaryReader) ReadBytes() []byte {
	length := br.ReadInt32()
	return br.stream.Next(int(length))
}

func (br *BinaryReader) ReadByte() (byte, error) {
	return br.stream.ReadByte()
}

func (br *BinaryReader) Read(data any) {
	switch data := data.(type) {
	case *bool:
		*data = br.ReadBool()
	case *int:
		*data = int(br.ReadInt32())
	case *int16:
		*data = int16(br.ReadInt32())
	case *int32:
		*data = br.ReadInt32()
	case *string:
		*data = br.ReadString()
	case *byte:
		*data, _ = br.ReadByte()
	case *[]byte:
		*data = br.ReadBytes()
	}
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
