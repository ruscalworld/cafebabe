package cafebabe

import (
	"encoding/binary"
	"io"
)

var order = binary.BigEndian

type Primitive interface {
	io.ReaderFrom
}

type U1 uint8

func (u *U1) ReadFrom(r io.Reader) (int64, error) {
	buf := [1]byte{}
	n, err := r.Read(buf[:])
	if err != nil {
		return int64(n), err
	}

	*u = U1(buf[0])
	return int64(n), nil
}

type U2 uint16

func (u *U2) ReadFrom(r io.Reader) (int64, error) {
	buf := [2]byte{}
	n, err := r.Read(buf[:])
	if err != nil {
		return int64(n), err
	}

	*u = U2(order.Uint16(buf[:]))
	return int64(n), nil
}

type U4 uint32

func (u *U4) ReadFrom(r io.Reader) (int64, error) {
	buf := [4]byte{}
	n, err := r.Read(buf[:])
	if err != nil {
		return int64(n), err
	}

	*u = U4(order.Uint32(buf[:]))
	return int64(n), nil
}
