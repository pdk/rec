package rec

import (
	"strings"
)

type RecordBuilder struct {
	fieldMap fieldNamePositionMap
}

func NewRecordBuilder(names ...string) RecordBuilder {

	fMap := fieldNamePositionMap{}

	for p, n := range names {
		fMap[strings.TrimSpace(n)] = p
	}

	return RecordBuilder{
		fieldMap: fMap,
	}
}

func (rb RecordBuilder) RecordFromStrings(values ...string) Record {
	vals := []Field{}

	for _, each := range values {
		each = strings.TrimSpace(each)
		vals = append(vals, NewField(each))
	}

	return Record{
		values:   vals,
		fieldMap: rb.fieldMap,
	}
}

func (rb RecordBuilder) RecordFromSQL(values ...any) Record {
	vals := []Field{}

	for _, each := range values {
		vals = append(vals, fieldFromScanned(each))
	}

	return Record{
		values:   vals,
		fieldMap: rb.fieldMap,
	}
}

func (rb RecordBuilder) Record(values ...any) Record {
	vals := []Field{}

	for _, each := range values {
		vals = append(vals, NewField(each))
	}

	return Record{
		values:   vals,
		fieldMap: rb.fieldMap,
	}
}
