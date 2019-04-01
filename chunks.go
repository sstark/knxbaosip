package main

// makeChunks compresses a list of integers to a list of consecutive chunks
// ({base,count} of consecutive elements)
func makeChunks(l []int) [][]int {
	var cl [][]int
	var count, last int
	for i, n := range l {
		if i == 0 {
			last = n
			continue
		}
		// item is just 1 larger than last item, so add to count instead of
		// appending
		if n-l[i-1] == 1 {
			count += 1
			continue
		}
		cl = append(cl, []int{last, count + 1})
		last = n
		count = 0
	}
	return append(cl, []int{last, count + 1})
}
