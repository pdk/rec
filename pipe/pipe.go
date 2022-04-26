package pipe

import (
	"fmt"
	"io"

	"github.com/pdk/rec"
)

type Producer func() chan rec.Record
type Pipe func(chan rec.Record) chan rec.Record
type Consumer func(chan rec.Record)

type Pipeline struct {
	first Producer
	then  []Pipe
	fini  Consumer
}

func New(p Producer) Pipeline {
	return Pipeline{
		first: p,
	}
}

func (p Pipeline) Then(t Pipe) Pipeline {
	p.then = append(p.then, t)
	return p
}

func (p Pipeline) Fini(f Consumer) {
	p.fini = f

	c := p.first()
	for _, t := range p.then {
		c = t(c)
	}

	p.fini(c)
}

func (p Pipeline) Print() {
	p.Fini(Print)
}

func ReadJSON(r io.Reader) Pipeline {
	return New(rec.JSONReader(r))
}

func ReadCSV(r io.Reader) Pipeline {
	return New(rec.CSVReader(r))
}

func (p Pipeline) Filter(f func(r rec.Record) bool) Pipeline {

	filter := func(in chan rec.Record) chan rec.Record {

		out := make(chan rec.Record)
		go func() {
			defer close(out)

			for r := range in {
				if f(r) {
					out <- r
				}
			}
		}()

		return out
	}

	return p.Then(filter)
}

func Print(ch chan rec.Record) {
	for rec := range ch {
		fmt.Println(rec)
	}
}

func (p Pipeline) Limit(limit int) Pipeline {

	return p.Then(func(in chan rec.Record) chan rec.Record {

		out := make(chan rec.Record)
		go func() {
			defer close(out)

			recCount := 0
			for r := range in {
				recCount++
				if recCount <= limit {
					out <- r
				}
			}
		}()

		return out
	})
}
