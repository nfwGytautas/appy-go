package utility

import "fmt"

type PagingSettings struct {
	Offset uint64
	Count  uint64
}

func (ps PagingSettings) String() string {
	return fmt.Sprintf("{Offset: %v, Count: %v}", ps.Offset, ps.Count)
}
