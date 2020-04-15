package magicurl

import "testing"


func TestGet(t *testing.T) {
	t.Run("found magic url contains original URL", func (t *testing.T){
		t.Errorf("Not implemented\n")
	})

	t.Run("found magic url contains short URL", func (t *testing.T){
		t.Errorf("Not implemented\n")
	})
	t.Run("slug not found returns GetMagicUrlSlugError", func (t *testing.T){
		t.Errorf("Not implemented\n")
	})
}

func TestCreate(t *testing.T) {
	t.Run("invalid URL returns CreateMagicURLSlugError", func (t *testing.T){
		t.Errorf("Not implemented\n")
	})
	t.Run("failure returns CreateMagicURLSlugError", func (t *testing.T){
		t.Errorf("Not implemented\n")
	})
	t.Run("created magic url contains original URL", func (t *testing.T){
		t.Errorf("Not implemented\n")
	})
	t.Run("created magic url contains short URL", func (t *testing.T){
		t.Errorf("Not implemented\n")
	})
}

func TestDelete(t *testing.T) {
	t.Run("target slug not found returns DeleteMagicURLSlugError", func (t *testing.T){
		t.Errorf("Not implemented\n")
	})
	t.Run("failure returns DeleteMagicURLSlugError", func (t *testing.T){
		t.Errorf("Not implemented\n")
	})
	t.Run("deleted magic url contains original URL", func (t *testing.T){
		t.Errorf("Not implemented\n")
	})
	t.Run("deleted magic url contains short URL", func (t *testing.T){
		t.Errorf("Not implemented\n")
	})
}