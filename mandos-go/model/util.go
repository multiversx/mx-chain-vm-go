package mandosjsonmodel

// JSONBytesFromTreeValues extracts values from a slice of JSONBytesFromTree into a list
func JSONBytesFromTreeValues(jbs []JSONBytesFromTree) [][]byte {
	result := make([][]byte, len(jbs))
	for i, jb := range jbs {
		result[i] = jb.Value
	}
	return result
}

// ToInt returns the int representation of the current TraceGasStatus
func (tgs TraceGasStatus) ToInt() int {
	return int(tgs)
}
