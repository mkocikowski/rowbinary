package rowbinary

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

func fieldIndexes(typ reflect.Type) []int {
	numField := typ.NumField()
	indexes := make([]int, 0, numField)
	for i := 0; i < numField; i++ {
		field := typ.Field(i)
		// https://golang.org/pkg/reflect/#StructField
		if field.PkgPath != "" { // unexported
			continue
		}
		if tag := field.Tag.Get("rowbinary"); tag == "-" {
			continue
		}
		indexes = append(indexes, i)
	}
	return indexes
}

func Columns(structPtr interface{}) (columns []string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("error extracting struct fields from %v: %v", reflect.TypeOf(structPtr), r)
		}
	}()
	typ := reflect.TypeOf(structPtr).Elem()
	for _, i := range fieldIndexes(typ) {
		field := typ.Field(i)
		name := field.Tag.Get("rowbinary")
		if name == "" {
			name = field.Name
		}
		columns = append(columns, name)
	}
	return columns, nil
}

func Marshal(buf io.Writer, structPtr interface{}) error {
	val := reflect.ValueOf(structPtr).Elem()
	typ := val.Type()
	for _, i := range fieldIndexes(typ) {
		if err := marshalValue(buf, val.Field(i)); err != nil {
			return fmt.Errorf("marshal field %q error: %v", typ.Field(i).Name, err)
		}
	}
	return nil
}

// panics if v == nil
func marshalValue(buf io.Writer, v reflect.Value) error {
	switch v.Kind() {
	case reflect.String:
		_, err := MarshalString(buf, v.String())
		return err
	case reflect.Slice:
		switch v.Type().Elem().Kind() {
		case reflect.Uint8:
			_, err := buf.Write(v.Bytes())
			return err
		}
	case reflect.Uint8:
		return binary.Write(buf, binary.LittleEndian, uint8(v.Uint()))
	case reflect.Uint16:
		return binary.Write(buf, binary.LittleEndian, uint16(v.Uint()))
	case reflect.Uint32:
		return binary.Write(buf, binary.LittleEndian, uint32(v.Uint()))
	case reflect.Uint64:
		return binary.Write(buf, binary.LittleEndian, uint64(v.Uint()))
	}
	return fmt.Errorf("value type %s not supported", v.Type())
}

func MarshalString(buf io.Writer, s string) (int, error) {
	var b []byte
	b = appendUleb128(b, uint64(len(s)))
	if n, err := buf.Write(b); err != nil {
		return n, err
	}
	return io.WriteString(buf, s)
}

// https://github.com/golang/go/blob/0ff9df6b53076a9402f691b07707f7d88d352722/src/cmd/internal/dwarf/dwarf.go#L194
// AppendUleb128 appends v to b using DWARF's unsigned LEB128 encoding.
func appendUleb128(b []byte, v uint64) []byte {
	for {
		c := uint8(v & 0x7f)
		v >>= 7
		if v != 0 {
			c |= 0x80
		}
		b = append(b, c)
		if c&0x80 == 0 {
			break
		}
	}
	return b
}

// https://en.wikipedia.org/w/index.php?title=LEB128&section=7#Decode_unsigned_integer
func readUleb128(r io.ByteReader) (uint64, error) {
	var result, shift uint64
	for {
		b, err := r.ReadByte()
		if err != nil {
			return result, err
		}
		result |= uint64(b&0x7f) << shift
		if b&0x80 == 0 {
			break
		}
		shift += 7
	}
	return result, nil
}
