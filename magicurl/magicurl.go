package magicurl

// Create shortened URL to the database.
func Create(originalURL string) (slug string, err error) {
	//perform input validate the url
	//get id from db, atomic increment value in the db
	//create slug as base64-encoded id
	//create MagicUrlItem, insert in the db
	//get slug from dynamo db
	return "raytran_slug", nil
}

func encodeURLSlug(id uint64) string {
	return ""
}
