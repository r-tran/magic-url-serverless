package magicurl

//MagicURL represents the database model
type MagicURL struct {
	Slug, OriginalURL string
}

//CreateMagicURLSlugError indicates unsuccessful slug query
type CreateMagicURLSlugError struct {
	Message string
	Err     error
}

func (m *CreateMagicURLSlugError) Error() string {
	return m.Message + ":" + m.Err.Error()
}

//GetMagicURLSlugError indicates unsuccessful slug query
type GetMagicURLSlugError struct {
	Message string
	Err     error
}

func (m *GetMagicURLSlugError) Error() string {
	return m.Message + ":" + m.Err.Error()
}

//DeleteMagicURLSlugError is error returned from failed delete magic url operation
type DeleteMagicURLSlugError struct {
	Message string
	Err     error
}

func (m *DeleteMagicURLSlugError) Error() string {
	return m.Message + ":" + m.Err.Error()
}
