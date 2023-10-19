package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	storage := &mockUserStorage{}
	handler := New(storage)
	assert.NotNil(t, handler)
}
