package rowbinary

import (
	"encoding/binary"
	"io"
)

/*
func Unmarshal(r io.Reader, structPtr interface{}) error {
	val := reflect.ValueOf(structPtr).Elem()
	typ := val.Type()
	for _, i := range fieldIndexes(typ) {
		if err := unmarshalValue(buf, val.Field(i)); err != nil {
			return fmt.Errorf("unmarshal field %q error: %v", typ.Field(i).Name, err)
		}
	}
	return nil
}
*/

/*
// panics if v == nil
func marshalValue(r io.Reader, v reflect.Value) error {
	switch v.Kind() {
	case reflect.String:
		_, err := UnmarshalString(r, v)
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
*/

func UnmarshalString(r io.Reader) (string, error) {
	n, err := readUleb128(r.(io.ByteReader))
	if err != nil {
		return "", err
	}
	b := make([]byte, int(n))
	if _, err := io.ReadFull(r, b); err != nil {
		return "", err
	}
	return string(b), nil
}

func UnmarshalUint8(r io.Reader) (uint8, error) {
	var i uint8
	err := binary.Read(r, binary.LittleEndian, &i)
	return i, err
}

func UnmarshalUint16(r io.Reader) (uint16, error) {
	var i uint16
	err := binary.Read(r, binary.LittleEndian, &i)
	return i, err
}

func UnmarshalUint32(r io.Reader) (uint32, error) {
	var i uint32
	err := binary.Read(r, binary.LittleEndian, &i)
	return i, err
}

func UnmarshalUint64(r io.Reader) (uint64, error) {
	var i uint64
	err := binary.Read(r, binary.LittleEndian, &i)
	return i, err
}

func UnmarshalFloat32(r io.Reader) (float32, error) {
	var i float32
	err := binary.Read(r, binary.LittleEndian, &i)
	return i, err
}

func UnmarshalFloat64(r io.Reader) (float64, error) {
	var i float64
	err := binary.Read(r, binary.LittleEndian, &i)
	return i, err
}
