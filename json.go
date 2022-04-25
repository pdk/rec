package rec

import (
	"encoding/json"
	"fmt"
	"io"
)

func JSONReader(r io.Reader) func() chan Record {
	return func() chan Record {

		ch := make(chan Record)
		go func() {
			defer close(ch)

			d := json.NewDecoder(r)
			if c, ok := r.(io.Closer); ok {
				defer c.Close()
			}

			for {
				next, err := readNextValue(d)
				if err == io.EOF {
					// ????
					return
				}
				if err != nil {
					return
				}

				rec, ok := next.(Record)
				if !ok {
					rec = Record{}.Append(next)
				}

				ch <- rec
			}
		}()

		return ch
	}
}

func readNextValue(d *json.Decoder) (any, error) {

	// Delim, for the four JSON delimiters [ ] { }
	// bool, for JSON booleans
	// float64, for JSON numbers
	// Number, for JSON numbers
	// string, for JSON string literals
	// nil, for JSON null

	t, err := d.Token()
	if err != nil {
		return nil, err
	}

	switch v := t.(type) {
	case json.Delim:
		switch v {
		case '{':
			return readObject(d)
		case '[':
			return readArray(d)
		default:
			return nil, fmt.Errorf("unexpected token %#v reading next value", v)
		}
	default:
		return v, nil
	}
}

func readArray(d *json.Decoder) (any, error) {

	r := Record{}

	for {
		t, err := d.Token()
		if err != nil {
			return r, err
		}

		switch t := t.(type) {
		case json.Delim:
			switch t {
			case ']': // end of the array
				return r, nil
			case '{': // new object
				v, err := readObject(d)
				if err != nil {
					return r, err
				}
				r = r.Append(v)
			case '[': // new array
				v, err := readArray(d)
				if err != nil {
					return r, err
				}
				r = r.Append(v)
			default:
				return r, fmt.Errorf("unexpected token %#v reading next value", t)
			}
		default:
			r = r.Append(t)
		}
	}
}

func readObject(d *json.Decoder) (any, error) {

	r := Record{}

	for {
		t1, err := d.Token()
		if err != nil {
			return r, err
		}

		var k string
		switch t := t1.(type) {
		case json.Delim:
			if t == '}' {
				return r, nil
			}
			return r, fmt.Errorf("unexpected delim %c reading next value", t)
		case string:
			k = t
		default:
			return r, fmt.Errorf("unexpected %t/%#v instead of string key", t1, t1)
		}

		v, err := readNextValue(d)
		if err != nil {
			return r, err
		}

		r = r.Put(k, v)
	}
}
