package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	userStorage := &mockUserStorage{}
	refreshTokenStorage := &mockRefreshTokenStorage{}
	got := New(userStorage, refreshTokenStorage)
	assert.NotNil(t, got)
}
