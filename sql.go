package rec

import (
	"database/sql"
	"log"
	"strings"
)

func SQLReader(db *sql.DB, query string, args ...any) func() chan Record {
	return func() chan Record {

		ch := make(chan Record)
		go func() {
			defer close(ch)

			rows, err := db.Query(query, args...)
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()

			names, scanBuilder := newScanBuilder(rows)
			rowBuilder := NewRecordBuilder(names...)

			for rows.Next() {

				values := scanBuilder()
				if err := rows.Scan(values...); err != nil {
					log.Fatal(err)
				}
				// log.Printf("scanned: %#v", values)

				ch <- rowBuilder.RecordFromSQL(values...)
			}
		}()

		return ch
	}
}

// newScanBuilder returns a function that will return a slice of new empty pointers that can be scanned into.
func newScanBuilder(rows *sql.Rows) ([]string, func() []any) {

	names := []string{}
	makers := []func() any{}

	types, err := rows.ColumnTypes()
	if err != nil {
		log.Fatalf("failed to get column types for query result: %v", err)
	}

	for _, e := range types {
		names = append(names, e.Name())
		makers = append(makers, scannable(e.DatabaseTypeName()))
	}

	return names, func() []any {

		newScannable := []any{}

		for _, e := range makers {
			newScannable = append(newScannable, e())
		}

		return newScannable
	}
}

// scannable returns a function, selected by typeName, that will return a newly allocated, scannable variable.
func scannable(typeName string) func() any {

	// notes:
	// 1. handling mysql "datetime" type requires connect param parseTime=true
	// 2. mysql type "bit" just handled as string

	switch strings.ToLower(typeName) {
	case "bool", "boolean":
		return func() any { return &sql.NullBool{} }
	case "byte":
		return func() any { return &sql.NullByte{} }
	case "float", "numeric", "number", "decimal", "double", "double precision",
		"float8", "real", "float4":
		return func() any { return &sql.NullFloat64{} }
	case "int", "int2", "int4", "bigint", "integer", "smallint", "tinyint", "mediumint",
		"serial", "bigserial", "smallserial", "serial2", "serial4":
		return func() any { return &sql.NullInt64{} }
	case "date", "time", "timestamp", "timestamp_tz", "timestamp_ntz", "year", "datetime":
		return func() any { return &sql.NullTime{} }
	}

	// default for every type, including varchar, text, char, etc, etc
	return func() any { return &sql.NullString{} }
}

func fieldFromScanned(val any) Field {

	switch v := val.(type) {
	case *sql.NullBool:
		if v.Valid {
			return NewField(v.Bool)
		}
	case *sql.NullByte:
		if v.Valid {
			return NewField(v.Byte)
		}
	case *sql.NullFloat64:
		if v.Valid {
			return NewField(v.Float64)
		}
	case *sql.NullInt64:
		if v.Valid {
			return NewField(v.Int64)
		}
	case *sql.NullString:
		if v.Valid {
			return NewField(v.String)
		}
	case *sql.NullTime:
		if v.Valid {
			return NewField(v.Time)
		}
	default:
		log.Fatalf("tried to convert unhandled SQL result type %t", val)
	}

	return NullField()
}
