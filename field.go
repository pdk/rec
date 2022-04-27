package rec

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// FieldType indicates what type of value is in the field.
type FieldType int

const (
	TypeNull FieldType = iota
	TypeString
	TypeInteger
	TypeFloat
	TypeBoolean
	TypeRecord
	TypeDateTime
)

// String returns a string for the particular FieldType.
func (t FieldType) String() string {
	switch t {
	case TypeNull:
		return "Null"
	case TypeString:
		return "String"
	case TypeInteger:
		return "Integer"
	case TypeFloat:
		return "Float"
	case TypeBoolean:
		return "Boolean"
	case TypeRecord:
		return "Record"
	case TypeDateTime:
		return "DateTime"
	}

	return "Invalid"
}

// Field holds one datum, of a particular type.
type Field struct {
	fieldType     FieldType
	stringValue   string
	intValue      int64
	floatValue    float64
	boolValue     bool
	recordValue   Record
	datetimeValue time.Time
}

func NewField(value any) Field {

	if value == nil {
		return NullField()
	}

	switch val := value.(type) {
	case int64:
		return IntegerField(val)
	case float64:
		return FloatField(val)
	case string:
		if strings.TrimSpace(val) == "" {
			return NullField()
		}
		return StringField(val)
	case bool:
		return BooleanField(val)
	case time.Time:
		return DateTimeField(val)
	case Record:
		return RecordField(val)
	}

	s := fmt.Sprintf("%v", value)
	return StringField(s)
}

func (f Field) StringValue() string {

	switch f.fieldType {
	case TypeNull:
		return ""
	case TypeString:
		return f.stringValue
	case TypeInteger:
		return strconv.FormatInt(f.intValue, 10)
	case TypeFloat:
		return strconv.FormatFloat(f.floatValue, 'g', -1, 64)
	case TypeBoolean:
		return strconv.FormatBool(f.boolValue)
	case TypeRecord:
		return f.recordValue.String()
	case TypeDateTime:
		return f.datetimeValue.Format(time.RFC3339)
	}

	return fmt.Sprintf("unhandled string type %d", f.fieldType)
}

func (f Field) AsString() (string, bool) {

	switch f.fieldType {
	case TypeNull:
		return "", true
	case TypeString:
		return f.stringValue, true
	case TypeInteger:
		return strconv.FormatInt(f.intValue, 10), true
	case TypeFloat:
		return strconv.FormatFloat(f.floatValue, 'g', -1, 64), true
	case TypeBoolean:
		return strconv.FormatBool(f.boolValue), true
	case TypeRecord:
		return f.recordValue.String(), true
	case TypeDateTime:
		return f.datetimeValue.Format(time.RFC3339), true
	}

	return fmt.Sprintf("unhandled string type %d", f.fieldType), false
}

func (f Field) String() string {

	switch f.fieldType {
	case TypeNull:
		return "null"
	case TypeString:
		return fmt.Sprintf("%#v", f.stringValue)
	case TypeInteger:
		return strconv.FormatInt(f.intValue, 10)
	case TypeFloat:
		return strconv.FormatFloat(f.floatValue, 'g', -1, 64)
	case TypeBoolean:
		return strconv.FormatBool(f.boolValue)
	case TypeRecord:
		return f.recordValue.String()
	case TypeDateTime:
		return f.datetimeValue.Format(time.RFC3339)
	}

	return fmt.Sprintf("unhandled string type %d", f.fieldType)
}

func NullField() Field {
	return Field{
		fieldType: TypeNull,
	}
}

func StringField(value string) Field {
	return Field{
		fieldType:   TypeString,
		stringValue: value,
	}
}

func IntegerField(value int64) Field {
	return Field{
		fieldType: TypeInteger,
		intValue:  value,
	}
}

func FloatField(value float64) Field {
	return Field{
		fieldType:  TypeFloat,
		floatValue: value,
	}
}

func BooleanField(value bool) Field {
	return Field{
		fieldType: TypeBoolean,
		boolValue: value,
	}
}

func DateTimeField(value time.Time) Field {
	return Field{
		fieldType:     TypeDateTime,
		datetimeValue: value,
	}
}

func RecordField(value Record) Field {
	return Field{
		fieldType:   TypeRecord,
		recordValue: value,
	}
}

func (f Field) IsNull() bool {
	return f.fieldType == TypeNull
}
