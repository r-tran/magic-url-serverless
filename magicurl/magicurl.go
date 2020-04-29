package magicurl

import (
	"fmt"
	"net/url"
	"strconv"
)

var magicURLTable = "magicUrl"

//MagicURL represents the database model
type MagicURL struct {
	Slug, OriginalURL string
}

// IsEmpty method used to determine uninitialized MagicUrl
func (m *MagicURL) IsEmpty() bool {
	return len(m.Slug) == 0 && len(m.OriginalURL) == 0
}

//Get returns the Magic URL given the shortened URL slug.
func Get(urlSlug string, client DatabaseService) (*MagicURL, error) {
	magicURLItem, err := client.Get(urlSlug)

	if err != nil {
		return nil, fmt.Errorf("Failed query to find MagicURL with slug %s : %w", urlSlug, err)
	}

	return magicURLItem, err
}

// Create returns a MagicURL containing the shortened URL slug for the originalURL.
func Create(originalURL string, client DatabaseService) (*MagicURL, error) {
	sanitizedURL, err := validateURL(originalURL)
	if err != nil {
		return nil, fmt.Errorf("Create MagicURL failed with invalid URL format for %s : %w", originalURL, err)
	}

	count, err := client.IncrementCounter()
	_, err = client.IncrementCounter()
	if err != nil {
		return nil, fmt.Errorf("Create MagicURL failed to update base counter for %s : %w", originalURL, err)
	}

	countStr := strconv.Itoa(count)
	slug, err := EncodeToBase62(countStr)
	if err != nil {
		return nil, fmt.Errorf("Create MagicURL failed base62 encoding for decimal %s : %w", countStr, err)
	}

	err = client.PutMagicURLItem(slug, sanitizedURL)
	if err != nil {
		return nil, fmt.Errorf("Create MagicURL failed creating item : %w", err)
	}

	return Get(slug, client)
}

//Delete removes MagicURl slug from datastore, returns the deleted slug.
func Delete(urlSlug string, client DatabaseService) (*MagicURL, error) {
	magicURLItem, err := client.DeleteMagicURLItem(urlSlug)
	if err != nil {
		return nil, fmt.Errorf("Failed to delete MagicURL with slug %s : %w", urlSlug, err)
	}

	return magicURLItem, nil
}

func validateURL(urlTarget string) (string, error) {
	parsedURL, err := url.ParseRequestURI(urlTarget)
	if err != nil {
		return "", err
	}
	return parsedURL.String(), err
}
