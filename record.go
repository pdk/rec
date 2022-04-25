package rec

import (
	"fmt"
	"strconv"
	"strings"
)

// Record is an ordered list of name/value pairs.
type Record struct {
	fieldMap fieldNamePositionMap
	values   []Field
}

// Put adds a name/value mapping.
func (r Record) Put(name string, value any) Record {
	if r.fieldMap == nil {
		r.fieldMap = fieldNamePositionMap{}
	}
	r.fieldMap.put(name, len(r.values))
	r.values = append(r.values, NewField(value))
	return r
}

// Append adds a field value without any name.
func (r Record) Append(value any) Record {
	r.values = append(r.values, NewField(value))
	return r
}

func (r Record) backMap() map[int]string {

	m := map[int]string{}

	if r.fieldMap != nil {
		for k, v := range r.fieldMap {
			m[v] = k
		}
	}

	return m
}

// String returns a string representation of the rec.
func (r Record) String() string {

	sb := strings.Builder{}

	backmap := r.backMap()

	sb.WriteString("(")
	for i := 0; i < len(r.values); i++ {
		if i > 0 {
			sb.WriteString(", ")
		}
		k, ok := backmap[i]
		if ok {
			sb.WriteString(fmt.Sprintf("%#v", k))
			sb.WriteString(": ")
		}
		sb.WriteString(r.values[i].String())
	}

	sb.WriteString(")")

	return sb.String()
}

func (r Record) Get(fieldName string) Field {
	p, ok := r.fieldMap[fieldName]
	if !ok || p >= len(r.values) {
		return NullField()
	}
	return r.values[p]
}

func (r Record) At(pos int) Field {
	if pos < 0 || pos >= len(r.values) {
		return NullField()
	}
	return r.values[pos]
}

func (r Record) Len() int {
	return len(r.values)
}

// Fields returns the list of field names in the order that the fields appear in the rec.
func (r Record) Fields() []string {

	fields := []string{}

	backmap := r.backMap()

	for i := 0; i < len(r.values); i++ {
		n, ok := backmap[i]
		if ok {
			fields = append(fields, n)
		}
	}

	return fields
}

// EnsureFieldNames makes sure all fields in the record have a name. (Useful when reading arrays.)
func (r Record) EnsureFieldNames() Record {

	if len(r.fieldMap) == len(r.values) {
		return r
	}

	backmap := r.backMap()

	for i := 0; i < len(r.values); i++ {
		n, ok := backmap[i]
		if !ok {
			n = strconv.Itoa(i)
			r.fieldMap[n] = i
		}
	}

	return r
}
