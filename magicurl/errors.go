package magicurl

//MagicURL represents the database model
type MagicURL struct {
	Slug, OriginalURL string
}

//CreateMagicURLSlugError indicates unsuccessful slug query
type CreateMagicURLSlugError struct {
	Message string
}

func (m *CreateMagicURLSlugError) Error() string {
	return m.Message
}

//GetMagicURLSlugError indicates unsuccessful slug query
type GetMagicURLSlugError struct {
	Message string
}

func (m *GetMagicURLSlugError) Error() string {
	return m.Message
}

//DeleteMagicURLSlugError is error returned from failed delete magic url operation
type DeleteMagicURLSlugError struct {
	Message string
}

func (m *DeleteMagicURLSlugError) Error() string {
	return m.Message
}
