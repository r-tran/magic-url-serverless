package magicurl

import (
	"fmt"
	"testing"
)

type MockDatabaseServiceError struct {
	GetFunc                func(string) (*MagicURL, error)
	IncrementCounterFunc   func() (int, error)
	PutMagicURLItemFunc    func(string, string) error
	DeleteMagicURLItemFunc func(string) (*MagicURL, error)
}

func (m *MockDatabaseServiceError) Get(slug string) (*MagicURL, error) {
	return m.GetFunc(slug)
}

func (m *MockDatabaseServiceError) IncrementCounter() (int, error) {
	return m.IncrementCounterFunc()
}

func (m *MockDatabaseServiceError) PutMagicURLItem(slug, originalURL string) error {
	return m.PutMagicURLItemFunc(slug, originalURL)
}

func (m *MockDatabaseServiceError) DeleteMagicURLItem(slug string) (*MagicURL, error) {
	return m.DeleteMagicURLItemFunc(slug)
}

func TestGet(t *testing.T) {
	t.Run("database query failure returns error", func(t *testing.T) {
		slug := "fakeSlug"
		client := &MockDatabaseServiceError{
			GetFunc: func(slug string) (*MagicURL, error) {
				return nil, fmt.Errorf("Failed to database get operation")
			},
		}

		_, err := Get(slug, client)
		if err == nil {
			t.Errorf("Expected non-nil error for Get")
		}
	})
}

func TestCreate(t *testing.T) {
	t.Run("invalid URL returns error", func(t *testing.T) {
		invalidOriginalURL := "thhjkdhsjkfhdkjfhkjdhfkjdsf"

		_, err := Create(invalidOriginalURL, nil)

		if err == nil {
			t.Errorf("Expected error for invalid URL %s\n", invalidOriginalURL)
		}
	})

	t.Run("database increment counter failure returns error", func(t *testing.T) {
		url := "https://google.com"
		client := &MockDatabaseServiceError{
			IncrementCounterFunc: func() (int, error) {
				return -1, fmt.Errorf("Failed database increment counter operation")
			},
		}

		_, err := Create(url, client)

		if err == nil {
			t.Errorf("Expected error for failed increment counter operation\n")
		}
	})

	t.Run("put magicurl database failure returns error", func(t *testing.T) {
		url := "https://google.com"
		client := &MockDatabaseServiceError{
			IncrementCounterFunc: func() (int, error) {
				return 10, nil
			},
			PutMagicURLItemFunc: func(_, _ string) error {
				return fmt.Errorf("Failed database put operation")
			},
		}

		_, err := Create(url, client)

		if err == nil {
			t.Errorf("Expected error for failed put operation")
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("database delete failure returns error", func(t *testing.T) {
		slug := "fakeSlug"
		client := &MockDatabaseServiceError{
			DeleteMagicURLItemFunc: func(_ string) (*MagicURL, error) {
				return nil, fmt.Errorf("Failed database delete operation")
			},
		}

		_, err := Delete(slug, client)
		if err == nil {
			t.Errorf("Expected error for failed delete operation")
		}
	})
}
