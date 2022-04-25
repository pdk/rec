package main

import (
	"strings"

	"github.com/pdk/rec"
	"github.com/pdk/rec/pipe"
)

const csvData = `
First Name, Last Name, Address,City,State, Zip
John,Doe,120 jefferson st.,Riverside, NJ, 08075
Jack,McGinnis,220 hobo Av.,Phila, PA,09119
"John ""Da Man""",Repici,120 Jefferson St.,Riverside, NJ,08075
Stephen,Tyler,"7452 Terrace ""At the Plaza"" road",SomeTown,SD, 91234
,Blankman,,SomeTown, SD, 00298
"Joan ""the bone"", Anne",Jet,"9th, at Terrace plc",Desert City,CO,00123
`

func main() {

	input := strings.NewReader(csvData)

	pipe.ReadCSV(input).
		Filter(haveName).
		Print()
}

func haveName(r rec.Record) bool {
	return !r.Get("First Name").IsNull()
}
