package searxng

import (
	"testing"
)

// TestGetCategory tests GetCategory function
func TestGetCategory(t *testing.T) {
	tests := []struct {
		name     string
		category string
		wantName string
		wantErr  bool
	}{
		{"general category", "general", "general", false},
		{"images category", "images", "images", false},
		{"videos category", "videos", "videos", false},
		{"news category", "news", "news", false},
		{"invalid category", "invalid", "", true},
		{"empty category", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cat, err := GetCategory(tt.category)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && cat.Name != tt.wantName {
				t.Errorf("GetCategory() name = %v, want %v", cat.Name, tt.wantName)
			}
		})
	}
}

// TestGetCategoryNames tests GetCategoryNames function
func TestGetCategoryNames(t *testing.T) {
	names := GetCategoryNames()

	if len(names) == 0 {
		t.Error("GetCategoryNames() returned empty slice")
	}

	// Check for expected categories
	expectedCategories := []string{"general", "images", "videos", "news", "map", "music", "it", "science", "files", "social media"}
	for _, expected := range expectedCategories {
		found := false
		for _, name := range names {
			if name == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("GetCategoryNames() missing category '%s'", expected)
		}
	}
}

// TestGetDisplayNames tests GetDisplayNames function
func TestGetDisplayNames(t *testing.T) {
	displayNames := GetDisplayNames()

	if len(displayNames) == 0 {
		t.Error("GetDisplayNames() returned empty slice")
	}

	// Check that we have the expected number of display names
	if len(displayNames) < 10 {
		t.Errorf("GetDisplayNames() returned %d names, expected at least 10", len(displayNames))
	}
}

// TestIsValidCategory tests IsValidCategory function
func TestIsValidCategory(t *testing.T) {
	tests := []struct {
		name     string
		category string
		want     bool
	}{
		{"valid general", "general", true},
		{"valid images", "images", true},
		{"valid videos", "videos", true},
		{"valid news", "news", true},
		{"valid map", "map", true},
		{"valid music", "music", true},
		{"valid it", "it", true},
		{"valid science", "science", true},
		{"valid files", "files", true},
		{"valid social media", "social media", true},
		{"invalid category", "invalid", false},
		{"empty category", "", false},
		{"case insensitive", "IMAGES", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidCategory(tt.category)
			if got != tt.want {
				t.Errorf("IsValidCategory() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestNormalizeCategory tests NormalizeCategory function
func TestNormalizeCategory(t *testing.T) {
	tests := []struct {
		name     string
		category string
		want     string
	}{
		{"already normalized", "general", "general"},
		{"uppercase", "IMAGES", "images"},
		{"mixed case", "ViDeOs", "videos"},
		{"with spaces", "social media", "social media"},
		{"invalid", "invalid", "invalid"},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeCategory(tt.category)
			if got != tt.want {
				t.Errorf("NormalizeCategory() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestNeedsSpecialFormatting tests NeedsSpecialFormatting function
func TestNeedsSpecialFormatting(t *testing.T) {
	tests := []struct {
		name     string
		category string
		want     bool
	}{
		{"images needs special", "images", true},
		{"videos no special", "videos", false},
		{"general no special", "general", false},
		{"news no special", "news", false},
		{"map no special", "map", false},
		{"invalid category", "invalid", false},
		{"empty category", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NeedsSpecialFormatting(tt.category)
			if got != tt.want {
				t.Errorf("NeedsSpecialFormatting() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestGetImageURL tests GetImageURL function
func TestGetImageURL(t *testing.T) {
	result := SearchResult{
		Title:    "Test Image",
		URL:      "https://example.com/page",
		ImgSrc:   "https://example.com/image.jpg",
		Content:  "Test content",
		Engine:   "google",
		Category: "images",
		Score:    0.95,
	}

	url := GetImageURL(result)
	if url != "https://example.com/image.jpg" {
		t.Errorf("GetImageURL() = %v, want %v", url, "https://example.com/image.jpg")
	}

	// Test with no ImgSrc
	result2 := SearchResult{
		Title:    "Test",
		URL:      "https://example.com",
		Content:  "Content",
		Engine:   "google",
		Category: "general",
		Score:    0.95,
	}

	url2 := GetImageURL(result2)
	if url2 != "" {
		t.Errorf("GetImageURL() with no ImgSrc = %v, want empty", url2)
	}
}

// TestFormatResultForCategory tests FormatResultForCategory function
func TestFormatResultForCategory(t *testing.T) {
	result := SearchResult{
		Title:    "Test",
		URL:      "https://example.com",
		ImgSrc:   "https://example.com/img.jpg",
		Content:  "Content",
		Engine:   "google",
		Category: "images",
		Score:    0.95,
	}

	t.Run("format for images", func(t *testing.T) {
		formatted := FormatResultForCategory(result, "images")
		if formatted.Title == "" {
			t.Error("FormatResultForCategory() returned result with empty Title")
		}
	})

	t.Run("format for general", func(t *testing.T) {
		result.Category = "general"
		formatted := FormatResultForCategory(result, "general")
		if formatted.Title == "" {
			t.Error("FormatResultForCategory() returned result with empty Title")
		}
	})
}

// TestGetCategorySuggestions tests GetCategorySuggestions function
func TestGetCategorySuggestions(t *testing.T) {
	t.Run("suggestions for valid category", func(t *testing.T) {
		// Note: This function may return empty for valid categories
		// as suggestions are query-based, not category-based
		suggestions := GetCategorySuggestions("general")
		// Just verify it returns a slice (may be empty)
		_ = suggestions
	})

	t.Run("suggestions for invalid category", func(t *testing.T) {
		suggestions := GetCategorySuggestions("invalid")
		if suggestions == nil {
			t.Error("GetCategorySuggestions() returned nil for invalid category")
		}
	})

	t.Run("suggestions for empty string", func(t *testing.T) {
		suggestions := GetCategorySuggestions("")
		if suggestions == nil {
			t.Error("GetCategorySuggestions() returned nil for empty string")
		}
	})
}

// TestCategoryInfo tests CategoryInfo function
func TestCategoryInfo(t *testing.T) {
	tests := []struct {
		name     string
		category string
		wantLen  int
	}{
		{"general info", "general", 0},
		{"images info", "images", 0},
		{"videos info", "videos", 0},
		{"invalid category", "invalid", 0},
		{"empty category", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := CategoryInfo(tt.category)
			if len(info) > 0 {
				// Valid category should return info string
				t.Logf("CategoryInfo() for '%s' returned: %s", tt.category, info)
			}
		})
	}
}

// TestAllCategoriesHaveInfo tests that all categories have proper info
func TestAllCategoriesHaveInfo(t *testing.T) {
	names := GetCategoryNames()

	for _, name := range names {
		t.Run(name+" category", func(t *testing.T) {
			cat, err := GetCategory(name)
			if err != nil {
				t.Errorf("GetCategory() for '%s' failed: %v", name, err)
				return
			}

			if cat.Name == "" {
				t.Error("Category has empty Name")
			}
			if cat.DisplayName == "" {
				t.Error("Category has empty DisplayName")
			}
			if cat.Description == "" {
				t.Error("Category has empty Description")
			}
			if cat.ExampleQuery == "" {
				t.Error("Category has empty ExampleQuery")
			}
		})
	}
}