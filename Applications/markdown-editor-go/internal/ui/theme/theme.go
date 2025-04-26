package theme

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// ThemeType represents the available themes
type ThemeType string

const (
	// LightTheme is the light color scheme
	LightTheme ThemeType = "light"

	// DarkTheme is the dark color scheme
	DarkTheme ThemeType = "dark"
)

// ThemeColors represents the color palette for a theme
type ThemeColors struct {
	Background          string `json:"background"`
	BackgroundSecondary string `json:"backgroundSecondary"`
	Text                string `json:"text"`
	TextSecondary       string `json:"textSecondary"`
	Border              string `json:"border"`
	Accent              string `json:"accent"`
	AccentHover         string `json:"accentHover"`
	EditorBackground    string `json:"editorBackground"`
	PreviewBackground   string `json:"previewBackground"`
	Toolbar             string `json:"toolbar"`
	StatusBar           string `json:"statusBar"`
	Highlight           string `json:"highlight"`
}

// Theme manages application theming
type Theme struct {
	ctx          context.Context
	currentTheme ThemeType
	lightColors  ThemeColors
	darkColors   ThemeColors
}

// NewTheme creates a new Theme manager
func NewTheme() *Theme {
	return &Theme{
		currentTheme: LightTheme,
		lightColors: ThemeColors{
			Background:          "#f9f7f7",
			BackgroundSecondary: "#f0f0f0",
			Text:                "#2d3436",
			TextSecondary:       "#636e72",
			Border:              "#dfe6e9",
			Accent:              "#74b9ff",
			AccentHover:         "#0984e3",
			EditorBackground:    "#ffffff",
			PreviewBackground:   "#f9f7f7",
			Toolbar:             "#f5f5f5",
			StatusBar:           "#f0f0f0",
			Highlight:           "rgba(116, 185, 255, 0.2)",
		},
		darkColors: ThemeColors{
			Background:          "#2d3436",
			BackgroundSecondary: "#222626",
			Text:                "#dfe6e9",
			TextSecondary:       "#b2bec3",
			Border:              "#636e72",
			Accent:              "#6c5ce7",
			AccentHover:         "#a29bfe",
			EditorBackground:    "#232323",
			PreviewBackground:   "#2d3436",
			Toolbar:             "#222626",
			StatusBar:           "#1e2022",
			Highlight:           "rgba(108, 92, 231, 0.2)",
		},
	}
}

// Initialize sets up the theme with application context
func (t *Theme) Initialize(ctx context.Context) {
	t.ctx = ctx
}

// ToggleTheme switches between light and dark themes
func (t *Theme) ToggleTheme() {
	if t.currentTheme == LightTheme {
		t.SetTheme(DarkTheme)
	} else {
		t.SetTheme(LightTheme)
	}
}

// SetTheme changes the active theme
func (t *Theme) SetTheme(themeType ThemeType) {
	t.currentTheme = themeType
	isDark := themeType == DarkTheme

	// Emit event to frontend to update theme
	if t.ctx != nil {
		runtime.EventsEmit(t.ctx, "theme:update", isDark)
	}
}

// GetCurrentTheme returns the current theme type
func (t *Theme) GetCurrentTheme() ThemeType {
	return t.currentTheme
}

// IsDarkMode returns true if dark mode is active
func (t *Theme) IsDarkMode() bool {
	return t.currentTheme == DarkTheme
}

// GetCurrentColors returns the color palette for the current theme
func (t *Theme) GetCurrentColors() ThemeColors {
	if t.currentTheme == DarkTheme {
		return t.darkColors
	}
	return t.lightColors
}

// GetColors returns the color palette for a specific theme
func (t *Theme) GetColors(themeType ThemeType) ThemeColors {
	if themeType == DarkTheme {
		return t.darkColors
	}
	return t.lightColors
}

// SetCustomColors updates the color palette for a theme
func (t *Theme) SetCustomColors(themeType ThemeType, colors ThemeColors) {
	if themeType == DarkTheme {
		t.darkColors = colors
	} else {
		t.lightColors = colors
	}

	// If updating the current theme, emit event to refresh UI
	if themeType == t.currentTheme && t.ctx != nil {
		runtime.EventsEmit(t.ctx, "theme:colors-update", colors)
	}
}
