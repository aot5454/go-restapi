package book

import "testing"

func TestNewFunc(t *testing.T) {
	t.Run("Should return new book handler", func(t *testing.T) {
		bookStorage := &bookStorageMockSuccess{}
		bookHandler := New(bookStorage)

		if bookHandler == nil {
			t.Errorf("Book handler should not be nil")
		}
	})
}