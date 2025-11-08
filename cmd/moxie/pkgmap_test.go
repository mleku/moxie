package main

import "testing"

func TestPackageMapping(t *testing.T) {
	pm := NewPackageMapping()

	tests := []struct {
		moxie string
		go_   string
	}{
		{"fmt", "fmt"},
		{"os", "os"},
		{"http", "http"},
		{"json", "json"},
		{"sql", "sql"},
		{"testing", "testing"},
	}

	for _, tt := range tests {
		t.Run(tt.moxie, func(t *testing.T) {
			// Test Moxie -> Go
			if got := pm.MoxieToGo(tt.moxie); got != tt.go_ {
				t.Errorf("MoxieToGo(%q) = %q; want %q", tt.moxie, got, tt.go_)
			}

			// Test Go -> Moxie
			if got := pm.GoToMoxie(tt.go_); got != tt.moxie {
				t.Errorf("GoToMoxie(%q) = %q; want %q", tt.go_, got, tt.moxie)
			}

			// Test existence checks
			if !pm.HasMoxieMapping(tt.moxie) {
				t.Errorf("HasMoxieMapping(%q) = false; want true", tt.moxie)
			}

			if !pm.HasGoMapping(tt.go_) {
				t.Errorf("HasGoMapping(%q) = false; want true", tt.go_)
			}
		})
	}
}

func TestPackageMappingUnknown(t *testing.T) {
	pm := NewPackageMapping()

	// Unknown packages should be returned unchanged
	unknown := "mypackage"

	if got := pm.MoxieToGo(unknown); got != unknown {
		t.Errorf("MoxieToGo(%q) = %q; want %q", unknown, got, unknown)
	}

	if got := pm.GoToMoxie(unknown); got != unknown {
		t.Errorf("GoToMoxie(%q) = %q; want %q", unknown, got, unknown)
	}

	if pm.HasMoxieMapping(unknown) {
		t.Errorf("HasMoxieMapping(%q) = true; want false", unknown)
	}

	if pm.HasGoMapping(unknown) {
		t.Errorf("HasGoMapping(%q) = true; want false", unknown)
	}
}
