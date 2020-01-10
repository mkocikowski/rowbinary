package rowbinary

import (
	"bytes"
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

func TestMarshalString(t *testing.T) {
	tests := []struct {
		s string
		b []byte
	}{
		{s: "", b: []byte{0}},
		{s: "foo", b: []byte{0x3, 0x66, 0x6f, 0x6f}},
	}
	for _, tt := range tests {
		buf := new(bytes.Buffer)
		_, err := MarshalString(buf, tt.s)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(buf.Bytes(), tt.b) {
			t.Fatal(buf.Bytes())
		}
	}
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

/*
func BenchmarkMarshal(b *testing.B) {
	r := &record{Foo: "foo", bar: []byte("bar"), baz: 1}
	for i := 0; i < b.N; i++ {
		Marshal(ioutil.Discard, r)
	}
}
*/
