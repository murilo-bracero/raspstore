package parser_test

import (
	"testing"

	"github.com/murilo-bracero/raspstore/file-service/internal/application/parser"
	"github.com/stretchr/testify/assert"
)

var flagtests = []struct {
	in  string
	out int
}{
	{"5M", 5242880},
	{"1G", 1073741824},
}

func TestParseUsage(t *testing.T) {
	for _, tt := range flagtests {
		t.Run(tt.in, func(t *testing.T) {
			res := parser.ParseUsage(tt.in)
			assert.Equal(t, tt.out, res)
		})
	}
}
