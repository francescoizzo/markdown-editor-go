package ui

import (
	"context"

	"github.com/francescoizzo/markdown-editor-go/internal/config"
	"github.com/francescoizzo/markdown-editor-go/internal/editor"
	"github.com/francescoizzo/markdown-editor-go/internal/ui/theme"
	"github.com/francescoizzo/markdown-editor-go/internal/utils"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// MainWindow represents the main application window
type MainWindow struct {
	ctx       context.Context
	config    *config.Config
	editor    *editor.Editor
	theme     *theme.Theme
	parser    *utils.MarkdownParser
	fileUtils *utils.FileUtils
}

// NewMainWindow creates a new main window instance
func NewMainWindow() *MainWindow {
	return &MainWindow{
		config:    config.DefaultConfig(),
		editor:    editor.NewEditor(),
		theme:     theme.NewTheme(),
		parser:    utils.NewMarkdownParser(),
		fileUtils: &utils.FileUtils{},
	}
}

// OnStartup is called when the app starts
func (w *MainWindow) OnStartup(ctx context.Context) {
	w.ctx = ctx

	// Initialize components
	w.editor.OnStartup(ctx)
	w.theme.Initialize(ctx)

	// Load configuration
	err := w.config.Load()
	if err != nil {
		runtime.LogError(ctx, "Failed to load configuration: "+err.Error())
	}

	// Apply configuration
	w.applyConfiguration()
}

// OnDomReady is called when the DOM is ready
func (w *MainWindow) OnDomReady(ctx context.Context) {
	w.editor.OnDomReady(ctx)

	// Set initial theme
	w.theme.SetTheme(w.getThemeFromConfig())

	// Set window size from config
	runtime.WindowSetSize(ctx, w.config.WindowWidth, w.config.WindowHeight)
}

// OnBeforeClose is called when the app is about to close
func (w *MainWindow) OnBeforeClose(ctx context.Context) bool {
	// Save window size to config
	width, height := runtime.WindowGetSize(ctx)
	w.config.WindowWidth = width
	w.config.WindowHeight = height
	w.config.Save()

	// Ask editor if it's OK to close (e.g., unsaved changes)
	return w.editor.OnBeforeClose(ctx)
}

// OnShutdown is called when the app is shutting down
func (w *MainWindow) OnShutdown(ctx context.Context) {
	w.editor.OnShutdown(ctx)
}

// NewFile creates a new file
func (w *MainWindow) NewFile() {
	w.editor.NewFile()
}

// OpenFile opens a markdown file
func (w *MainWindow) OpenFile() bool {
	success := w.editor.OpenFile()
	if success {
		// Add to recent files in config
		if w.editor.GetCurrentFilePath() != "" {
			w.config.AddRecentFile(w.editor.GetCurrentFilePath())
			w.config.Save()
		}
	}
	return success
}

// SaveFile saves the current file
func (w *MainWindow) SaveFile() bool {
	return w.editor.SaveFile()
}

// SaveFileAs prompts for a filename and saves
func (w *MainWindow) SaveFileAs() bool {
	success := w.editor.SaveFileAs()
	if success {
		// Add to recent files in config
		if w.editor.GetCurrentFilePath() != "" {
			w.config.AddRecentFile(w.editor.GetCurrentFilePath())
			w.config.Save()
		}
	}
	return success
}

// ToggleTheme switches between light and dark mode
func (w *MainWindow) ToggleTheme() {
	w.theme.ToggleTheme()
	w.config.IsDarkMode = w.theme.IsDarkMode()
	w.config.Save()
}

// ToggleAutoSave enables or disables autosave
func (w *MainWindow) ToggleAutoSave() bool {
	enabled := w.editor.ToggleAutoSave()
	w.config.AutoSaveEnabled = enabled
	w.config.Save()
	return enabled
}

// SetAutoSaveDelay updates the autosave delay
func (w *MainWindow) SetAutoSaveDelay(seconds int) {
	w.editor.SetAutoSaveDelay(seconds)
	w.config.AutoSaveDelay = seconds
	w.config.Save()
}

// GetRecentFiles returns the list of recent files
func (w *MainWindow) GetRecentFiles() []string {
	return w.config.RecentFiles
}

// SetContent updates the editor content
func (w *MainWindow) SetContent(content string) {
	w.editor.SetContent(content)
}

// GetContent returns the current content
func (w *MainWindow) GetContent() string {
	return w.editor.GetContent()
}

// GetWordCount returns the word count for the current content
func (w *MainWindow) GetWordCount() int {
	content := w.editor.GetContent()
	return w.parser.WordCount(content)
}

// ExtractTOC generates a table of contents from the markdown
func (w *MainWindow) ExtractTOC() string {
	content := w.editor.GetContent()
	return w.parser.ExtractTOC(content)
}

// Internal helper methods

// applyConfiguration applies the loaded configuration to components
func (w *MainWindow) applyConfiguration() {
	// Apply theme
	w.theme.SetTheme(w.getThemeFromConfig())

	// Apply editor settings
	w.editor.SetAutoSaveEnabled(w.config.AutoSaveEnabled)
	w.editor.SetAutoSaveDelay(w.config.AutoSaveDelay)
}

// getThemeFromConfig gets the theme type from configuration
func (w *MainWindow) getThemeFromConfig() theme.ThemeType {
	if w.config.IsDarkMode {
		return theme.DarkTheme
	}
	return theme.LightTheme
}
