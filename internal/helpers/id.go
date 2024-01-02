package helpers

import (
	"fmt"
	"strings"
)

func JoinIDs(ids []uint) string {
	var id strings.Builder
	for _, v := range ids {
		id.WriteString(fmt.Sprintf("%v,", v))
	}

	return id.String()[:len(id.String())-1]
}
