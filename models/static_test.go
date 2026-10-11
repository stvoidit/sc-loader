package models

import (
	"encoding/json/v2"
	"testing"
)

func TestResponseStaticGetters(t *testing.T) {
	var rs ResponseStatic
	rs.Static.Features.MMPExternalSourceOrigin = "https://mmp.example.com"
	rs.Static.FeaturesV2.PlayerModuleExternalLoading.MmpVersion = "1.2.3"

	if got := rs.GetMMPOrigin(); got != "https://mmp.example.com" {
		t.Errorf("GetMMPOrigin() = %q, want %q", got, "https://mmp.example.com")
	}
	if got := rs.GetMMPVerstion(); got != "1.2.3" {
		t.Errorf("GetMMPVerstion() = %q, want %q", got, "1.2.3")
	}
	if got, want := rs.GetBaseMMP(), "https://mmp.example.com/1.2.3"; got != want {
		t.Errorf("GetBaseMMP() = %q, want %q", got, want)
	}
}

func TestResponseStaticZeroValue(t *testing.T) {
	var rs ResponseStatic
	if got := rs.GetBaseMMP(); got != "/" {
		t.Errorf("GetBaseMMP() = %q, want %q", got, "/")
	}
}

// TestResponseStaticJSON проверяет json-теги полей, которые используются геттерами.
func TestResponseStaticJSON(t *testing.T) {
	const body = `{
		"static": {
			"env": "production",
			"features": {"MMPExternalSourceOrigin": "https://mmp.cdn"},
			"featuresV2": {"playerModuleExternalLoading": {"mmpVersion": "9.9.9", "previewMMPVersion": "10.0.0"}},
			"unknownField": {"a": 1}
		}
	}`

	var rs ResponseStatic
	if err := json.Unmarshal([]byte(body), &rs); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if got, want := rs.GetBaseMMP(), "https://mmp.cdn/9.9.9"; got != want {
		t.Errorf("GetBaseMMP() = %q, want %q", got, want)
	}
	if got := rs.Static.FeaturesV2.PlayerModuleExternalLoading.PreviewMMPVersion; got != "10.0.0" {
		t.Errorf("PreviewMMPVersion = %q, want %q", got, "10.0.0")
	}
}
