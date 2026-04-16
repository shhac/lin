package linear

// FetchAll paginates through all pages of a connection query, collecting all nodes.
// The fetch function takes (first int, after *string) and returns (nodes, hasNextPage, endCursor, error).
func FetchAll[T any](
	fetch func(first int, after *string) ([]T, bool, *string, error),
) ([]T, error) {
	const pageSize = 50
	var all []T
	var cursor *string
	for {
		nodes, hasNext, endCursor, err := fetch(pageSize, cursor)
		if err != nil {
			return nil, err
		}
		all = append(all, nodes...)
		if !hasNext {
			break
		}
		cursor = endCursor
	}
	return all, nil
}
