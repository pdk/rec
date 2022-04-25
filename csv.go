package rec

import (
	"encoding/csv"
	"io"
)

func CSVReader(r io.Reader) func() chan Record {
	return func() chan Record {

		ch := make(chan Record)
		go func() {
			defer close(ch)

			c := csv.NewReader(r)
			if c, ok := r.(io.Closer); ok {
				defer c.Close()
			}

			header, err := c.Read()
			if err != nil {
				return
			}

			b := NewRecordBuilder(header...)

			for {
				rec, err := c.Read()
				if err == io.EOF {
					return
				}
				if err != nil {
					return
				}

				ch <- b.RecordFromStrings(rec...)
			}

		}()

		return ch
	}
}

// func ReadCSV(r io.Reader) error {

// 	c := csv.NewReader(r)
// 	if c, ok := r.(io.Closer); ok {
// 		defer c.Close()
// 	}

// 	header, err := c.Read()
// 	if err != nil {
// 		return err
// 	}

// 	b := NewRecordBuilder(header...)

// 	for {
// 		record, err := c.Read()
// 		if err == io.EOF {
// 			return nil
// 		}
// 		if err != nil {
// 			return err
// 		}

// 		rec := b.RecordFromStrings(rec...)

// 		log.Printf("got: %s", rec)
// 	}
// }
