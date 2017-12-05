package transporter_test

import (
	"testing"

	"github.com/msales/transporter"
	"github.com/stretchr/testify/assert"
)

func TestApplication_IsHealthy(t *testing.T) {
	a := transporter.NewApplication()

	assert.Nil(t, a.IsHealthy())
}
