package types

// Query endpoints supported by the check querier
const (
	QueryIssue = "issue"
)

/*
Below you will be able how to set your own queries:


// QueryResList Queries Result Payload for a query
type QueryResList []string

// implement fmt.Stringer
func (n QueryResList) String() string {
	return strings.Join(n[:], "\n")
}

*/
