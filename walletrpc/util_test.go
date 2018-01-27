package walletrpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXMRToDecimal(t *testing.T) {
	assert.Equal(t, "0.034000200000", XMRToDecimal(34000200000))
	assert.Equal(t, "15.000000000000", XMRToDecimal(15e12))
}
