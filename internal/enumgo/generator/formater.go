package generator

import (
	"fmt"
	"go/format"
)

func formatCode(src []byte) ([]byte, error) {
	formatted, err := format.Source(src)
	if err != nil {
		return nil, fmt.Errorf("failed to format code: %v", err)
	}
	return formatted, nil
}
