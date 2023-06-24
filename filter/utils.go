package filter

// FilterByName returns the filter matching the given name.
func FilterByName(fts []Filter, name string) (Filter, bool) {
	for _, f := range fts {
		if f.Name == name {
			return f, true
		}
	}

	return Filter{}, false
}

// RemoveFilterByName removes the filter with the given name and returns the rest of the filters slice.
func RemoveFilterByName(fts []Filter, name string) ([]Filter, bool) {
	var nf []Filter

	for _, f := range fts {
		if f.Name != name {
			nf = append(nf, f)
		}
	}

	return nf, false
}
