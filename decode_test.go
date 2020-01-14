package rowbinary

import (
	"bytes"
	"encoding/base64"
	"testing"
)

func TestUnmarshalUint64(t *testing.T) {
	// clickhouse-client -q"select toUInt64(1) format RowBinary" | base64 -w0
	b, err := base64.StdEncoding.DecodeString("AQAAAAAAAAA=")
	if err != nil {
		t.Fatal(err)
	}
	r := bytes.NewReader(b)
	i, err := UnmarshalUint64(r)
	if err != nil {
		t.Fatal(err)
	}
	if i != uint64(1) {
		t.Fatal(i)
	}
	if _, err := UnmarshalUint64(r); err == nil {
		t.Fatal("expected error")
	}
}

func TestUnmarshalString(t *testing.T) {
	// clickhouse-client -q"select 'foo' format RowBinary" | base64 -w0
	b, err := base64.StdEncoding.DecodeString("A2Zvbw==")
	if err != nil {
		t.Fatal(err)
	}
	r := bytes.NewReader(b)
	s, err := UnmarshalString(r)
	if err != nil {
		t.Fatal(err)
	}
	if s != "foo" {
		t.Fatal(s)
	}
	if _, err := UnmarshalString(r); err == nil {
		t.Fatal("expected error")
	}
}
