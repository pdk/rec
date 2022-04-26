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

			scanBuilder := newScanBuilder(rows)

			for rows.Next() {

				values := scanBuilder()
				if err := rows.Scan(values...); err != nil {
					log.Fatal(err)
				}

				log.Printf("scanned: %#v", values)

				ch <- Record{}
			}
		}()

		return ch
	}
}

// newScanBuilder returns a function that will return a slice of new empty pointers that can be scanned into.
func newScanBuilder(rows *sql.Rows) func() []any {

	makers := []func() any{}

	types, err := rows.ColumnTypes()
	if err != nil {
		log.Fatalf("failed to get column types for query result: %v", err)
	}

	for _, e := range types {
		makers = append(makers, scannable(e.DatabaseTypeName()))
	}

	return func() []any {

		newScannable := []any{}

		for _, e := range makers {
			newScannable = append(newScannable, e())
		}

		return newScannable
	}
}

// scannable returns a function, selected by typeName, that will return a newly allocated, scannable variable.
func scannable(typeName string) func() any {

	switch strings.ToLower(typeName) {
	case "bool", "boolean", "bit":
		return func() any { return &sql.NullBool{} }
	case "byte":
		return func() any { return &sql.NullByte{} }
	case "float", "numeric", "number", "decimal", "double", "double precision",
		"float8", "real", "float4":
		return func() any { return &sql.NullFloat64{} }
	case "int", "int2", "int4", "bigint", "integer", "smallint", "tinyint", "mediumint",
		"serial", "bigserial", "smallserial", "serial2", "serial4":
		return func() any { return &sql.NullInt64{} }
	case "date", "time", "datetime", "timestamp", "timestamp_tz", "timestamp_ntz", "year":
		return func() any { return &sql.NullTime{} }
	}

	// varchar, text, char, etc, etc
	return func() any { return &sql.NullString{} }
}
