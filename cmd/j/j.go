package main

import (
	"strings"

	"github.com/pdk/rec"
	"github.com/pdk/rec/pipe"
)

const jsonStream = `
{"Name": "Ed", "Text": "Knock knock.", "count": 10, "ready": false}
{"Name": "Sam", "Text": "Who's there?"}
1234
true 
[]
{}
"blue" 
["aa", 12, 23.93]
{"A": {"B": true}}
{"null": null}
{"null": true}
[true, false, "a", 12, 12.3]
`

const jsonStreamx = `
{"A": {"B": true}}
`

func main() {

	input := strings.NewReader(jsonStream)

	pipe.ReadJSON(input).
		Filter(isBlue).
		Filter(hasFields).
		Print()
}

func isBlue(r rec.Record) bool {
	return r.At(0).StringValue() != "blue"
}

func hasFields(r rec.Record) bool {
	return r.Len() > 0
}
