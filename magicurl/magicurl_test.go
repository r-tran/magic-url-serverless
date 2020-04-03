package magicurl

import "testing"

func TestMagicUrl(t *testing.T) { 
	res := "hi"
	actual := encodeURLSlug(100)
	if res != actual {
		t.Errorf("Expected: %v Actual: %v", res, actual)
	}
}
