package rec

// fieldNamePositionMap provides mapping from field/column names to column position, i.e. offset in values Field slice.
type fieldNamePositionMap map[string]int

// put adds a name/position mapping in the field map.
func (f fieldNamePositionMap) put(name string, p int) fieldNamePositionMap {
	f[name] = p
	return f
}
