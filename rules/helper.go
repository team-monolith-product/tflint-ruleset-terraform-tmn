package rules

import (
	"github.com/hashicorp/hcl/v2"
	"strings"
)

func ref(hr hcl.Range) *hcl.Range {
	return &hr
}

// RemoveSpaceAndLine remove space, "\t" and "\n" from the given string
func RemoveSpaceAndLine(str string) string {
	newStr := strings.ReplaceAll(str, " ", "")
	newStr = strings.ReplaceAll(newStr, "\t", "")
	newStr = strings.ReplaceAll(newStr, "\n", "")
	return newStr
}
