package exceptions

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDomainError(t *testing.T) {
	t.Run("Should create new domain error", func(t *testing.T) {
		reason := "error creating payment"
		err := NewDomainError(reason)

		assert.NotNil(t, err)
		assert.Equal(t, reason, err.Error())
	})
}
