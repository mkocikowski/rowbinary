package rowbinary

import (
	"bytes"
	"math"
	"reflect"
	"testing"
)

func TestColumns(t *testing.T) {
	record := struct {
		Foo    uint8
		bar    uint8 `rowbinary:"BAR"`
		Baz    uint8 `rowbinary:"monkey"`
		Banana uint8 `rowbinary:"-"`
	}{}
	c, err := Columns(&record)
	if !reflect.DeepEqual(c, []string{"Foo", "monkey"}) {
		t.Fatal(c)
	}
	if err != nil {
		t.Fatal(err)
	}
}

func TestColumnsError(t *testing.T) {
	for _, tt := range []interface{}{nil, 1, struct{}{}} {
		_, err := Columns(tt)
		if err == nil {
			t.Fatal(tt, err)
		}
		t.Log(err)
	}
}

func TestMarshalValue(t *testing.T) {
	tests := []struct {
		v    interface{}
		want []byte
	}{
		{v: "", want: []byte{0}},
		{v: "foo", want: []byte{0x3, 0x66, 0x6f, 0x6f}},
		{v: []byte("foo"), want: []byte{102, 111, 111}},
		{v: uint16(1), want: []byte{1, 0}},
	}
	for _, tt := range tests {
		buf := new(bytes.Buffer)
		if err := marshalValue(buf, reflect.ValueOf(tt.v)); err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(buf.Bytes(), tt.want) {
			t.Fatal(buf.Bytes())
		}
	}
}

func TestMarshalValueError(t *testing.T) {
	i := uint32(1)
	b := []byte{1}
	tests := []struct {
		v       interface{}
		wantErr bool
	}{
		{v: i, wantErr: false},
		{v: &i, wantErr: true},
		{v: b, wantErr: false},
		{v: &b, wantErr: true},
		{v: []int16{1}, wantErr: true},
		{v: float64(1), wantErr: true},
	}
	for _, tt := range tests {
		buf := new(bytes.Buffer)
		err := marshalValue(buf, reflect.ValueOf(tt.v))
		if (err != nil) != tt.wantErr {
			t.Fatal(err)
		}
	}
}

func TestMarshalValueNil(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	buf := new(bytes.Buffer)
	marshalValue(buf, reflect.ValueOf(nil))
}

func TestMarshal(t *testing.T) {
	record := struct{ Foo uint16 }{Foo: 0xff}
	buf := new(bytes.Buffer)
	if err := Marshal(buf, &record); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf.Bytes(), []byte{0xff, 0x0}) {
		t.Fatal(buf.Bytes())
	}
}

func TestMarshalErrors(t *testing.T) {
	// TODO
}

func TestUleb128(t *testing.T) {
	for _, i := range []uint64{0, 1, math.MaxUint32, math.MaxUint32 + 1, math.MaxUint64} {
		b := appendUleb128(nil, i)
		j, _ := readUleb128(bytes.NewReader(b))
		if j != i {
			t.Fatal(i, j)
		}
	}
}

/*
func BenchmarkMarshal(b *testing.B) {
	r := &record{Foo: "foo", bar: []byte("bar"), baz: 1}
	for i := 0; i < b.N; i++ {
		Marshal(ioutil.Discard, r)
	}
}
*/
