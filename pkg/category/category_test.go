package category

import "testing"

func TestValidateEnglishCategory(t *testing.T) {
	if err := ValidateEnglishCategory("business"); err != nil {
		t.Fatalf("expected business to validate: %v", err)
	}

	if err := ValidateEnglishCategory("Business"); err == nil {
		t.Fatalf("expected Business to fail")
	}

	if err := ValidateEnglishCategory("unknown"); err == nil {
		t.Fatalf("expected unknown to fail")
	}
}
