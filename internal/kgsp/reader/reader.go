package reader

import (
	"fmt"
	"gontainer/internal/kgsp/spvalue"
	"strconv"
)

func (spr *SPReader) readLine() (line []byte, n int, err error) {
	for {
		b, err := spr.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

func (spr *SPReader) readInteger() (x int, n int, err error) {
	line, n, err := spr.readLine()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}
	return int(i64), n, nil
}

func (spr *SPReader) Read() (spvalue.Value, error) {
	_type, err := spr.reader.ReadByte()

	if err != nil {
		return spvalue.Value{}, err
	}

	switch _type {
	case spvalue.ARRAY:
		return spr.readArray()
	case spvalue.BULK:
		return spr.readBulk()
	default:
		fmt.Printf("Unknown type: %v", string(_type))
		return spvalue.Value{}, nil
	}
}

func (spr *SPReader) readArray() (spvalue.Value, error) {
	v := spvalue.Value{}
	v.Typ = "array"

	// read length of array
	length, _, err := spr.readInteger()
	if err != nil {
		return v, err
	}

	// foreach line, parse and read the value
	v.Array = make([]spvalue.Value, length)
	for i := 0; i < length; i++ {
		val, err := spr.Read()
		if err != nil {
			return v, err
		}

		v.Array[i] = val
	}

	return v, nil
}

func (spr *SPReader) readBulk() (spvalue.Value, error) {
	v := spvalue.Value{}

	v.Typ = "bulk"

	l, _, err := spr.readInteger()
	if err != nil {
		return v, err
	}

	bulk := make([]byte, l)

	_, err = spr.reader.Read(bulk)
	if err != nil {
		return v, err
	}

	v.Bulk = string(bulk)

	// Read the trailing CRLF
	_, _, err = spr.readLine()
	if err != nil {
		return v, err
	}

	return v, nil
}
