package config

import (
	"testing"
)

func TestValidDevice(t *testing.T) {
	tests := []struct {
		device string
		want   bool
	}{
		{"note", true},
		{"flagship", true},
		{"Note", false},     // case sensitive
		{"Flagship", false}, // case sensitive
		{"", false},
		{"invalid", false},
		{"board", false},
	}

	for _, tt := range tests {
		got := ValidDevice(tt.device)
		if got != tt.want {
			t.Errorf("ValidDevice(%q) = %v, want %v", tt.device, got, tt.want)
		}
	}
}

func TestConfig_MaskedToken_Empty(t *testing.T) {
	cfg := &Config{Token: ""}
	got := cfg.MaskedToken()
	if got != "(not set)" {
		t.Errorf("MaskedToken() = %q, want \"(not set)\"", got)
	}
}

func TestConfig_MaskedToken_Short(t *testing.T) {
	cfg := &Config{Token: "abc"}
	got := cfg.MaskedToken()
	if got != "****" {
		t.Errorf("MaskedToken() = %q, want \"****\"", got)
	}
}

func TestConfig_MaskedToken_Normal(t *testing.T) {
	cfg := &Config{Token: "abcdefgh1234"}
	got := cfg.MaskedToken()
	if got != "****1234" {
		t.Errorf("MaskedToken() = %q, want \"****1234\"", got)
	}
}

func TestConfig_GetToken_Empty(t *testing.T) {
	cfg := &Config{Token: ""}
	_, err := cfg.GetToken()
	if err == nil {
		t.Error("GetToken() should return error for empty token")
	}
}

func TestConfig_GetToken_Set(t *testing.T) {
	cfg := &Config{Token: "mytoken"}
	token, err := cfg.GetToken()
	if err != nil {
		t.Errorf("GetToken() error = %v, want nil", err)
	}
	if token != "mytoken" {
		t.Errorf("GetToken() = %q, want \"mytoken\"", token)
	}
}

func TestDeviceConstants(t *testing.T) {
	if DeviceNote != "note" {
		t.Errorf("DeviceNote = %q, want \"note\"", DeviceNote)
	}
	if DeviceFlagship != "flagship" {
		t.Errorf("DeviceFlagship = %q, want \"flagship\"", DeviceFlagship)
	}
}
