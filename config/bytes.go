package config

import (
	"github.com/dustin/go-humanize"
)

type Bytes struct {
	bytesUint64 uint64
}

func NewBytes(b uint64) Bytes {
	return Bytes{bytesUint64: b}
}

func ParseBytes(bytesStr string) (b Bytes, err error) {
	if err = b.parse(bytesStr); err != nil {
		return Bytes{}, err
	}
	return b, nil
}

func (b Bytes) Int() int {
	return int(b.bytesUint64)
}

func (b Bytes) Int8() int8 {
	return int8(b.bytesUint64)
}

func (b Bytes) Int16() int16 {
	return int16(b.bytesUint64)
}

func (b Bytes) Int32() int32 {
	return int32(b.bytesUint64)
}

func (b Bytes) Int64() int64 {
	return int64(b.bytesUint64)
}

func (b Bytes) Uint8() uint8 {
	return uint8(b.bytesUint64)
}

func (b Bytes) Uint16() uint16 {
	return uint16(b.bytesUint64)
}

func (b Bytes) Uint32() uint32 {
	return uint32(b.bytesUint64)
}

func (b Bytes) Uint64() uint64 {
	return b.bytesUint64
}

func (b *Bytes) UnmarshalYAML(f func(interface{}) error) error {
	var bytesStr string
	err := f(&bytesStr)
	if err != nil {
		return err
	}
	if err = b.parse(bytesStr); err != nil {
		return err
	}
	return nil
}

func (b *Bytes) parse(bytesStr string) (err error) {
	if len(bytesStr) == 0 {
		return nil
	}
	b.bytesUint64, err = humanize.ParseBytes(bytesStr)
	if err != nil {
		return err
	}
	return nil
}
