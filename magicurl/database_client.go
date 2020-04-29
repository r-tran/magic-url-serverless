package magicurl

//DatabaseService wraps functionality for storing and retrieving the MagicURL from the datastore
type DatabaseService interface {
	Get(slug string) (*MagicURL, error)
	IncrementCounter() (int, error)
	PutMagicURLItem(slug, originalURL string) error
	DeleteMagicURLItem(slug string) (*MagicURL, error)
}
