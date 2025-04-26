package editor

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// FileUtils provides file operation utilities within the editor
type FileUtils struct{}

// SaveToFile writes content to a file
func (f *FileUtils) SaveToFile(path string, content string) error {
	return ioutil.WriteFile(path, []byte(content), 0644)
}

// ReadFromFile reads content from a file
func (f *FileUtils) ReadFromFile(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// GetFilenameFromPath extracts the filename from a path
func (f *FileUtils) GetFilenameFromPath(path string) string {
	return filepath.Base(path)
}

// Editor represents the main editor functionality
type Editor struct {
	ctx             context.Context
	currentFilePath string
	content         string
	lastSaveTime    time.Time
	isDarkMode      bool
	autoSaveEnabled bool
	autoSaveDelay   time.Duration
	lastEditTime    time.Time
	autoSaveTimer   *time.Timer
	fileUtils       *FileUtils
}

// NewEditor creates a new instance of the Markdown editor
func NewEditor() *Editor {
	return &Editor{
		isDarkMode:      false,
		autoSaveEnabled: true,
		autoSaveDelay:   5 * time.Second, // 5 second autosave delay by default
		fileUtils:       &FileUtils{},
		lastSaveTime:    time.Now(),
		lastEditTime:    time.Now(),
	}
}

// OnStartup is called when the app starts
func (e *Editor) OnStartup(ctx context.Context) {
	e.ctx = ctx
}

// OnDomReady is called when the DOM is ready
func (e *Editor) OnDomReady(ctx context.Context) {
	// Initialize UI components when DOM is ready
	runtime.EventsEmit(e.ctx, "theme:update", e.isDarkMode)
}

// OnBeforeClose is called when the app is about to close
func (e *Editor) OnBeforeClose(ctx context.Context) bool {
	// Check for unsaved changes and prompt the user
	return true
}

// OnShutdown is called when the app is shutting down
func (e *Editor) OnShutdown(ctx context.Context) {
	// Perform cleanup
	if e.autoSaveEnabled && e.hasUnsavedChanges() {
		e.SaveFile()
	}
}

// SetContent updates the editor content
func (e *Editor) SetContent(content string) {
	e.content = content
	e.lastEditTime = time.Now()

	// Schedule autosave
	if e.autoSaveEnabled {
		if e.autoSaveTimer != nil {
			e.autoSaveTimer.Stop()
		}
		e.autoSaveTimer = time.AfterFunc(e.autoSaveDelay, func() {
			if e.hasUnsavedChanges() {
				e.AutoSave()
			}
		})
	}

	// Emit event to update the preview
	runtime.EventsEmit(e.ctx, "content:update", e.RenderHTML())
}

// GetContent returns the current content
func (e *Editor) GetContent() string {
	return e.content
}

// RenderHTML converts markdown to HTML
func (e *Editor) RenderHTML() string {
	if e.content == "" {
		return ""
	}

	md := []byte(e.content)
	html := markdown.ToHTML(md, nil, nil)
	return string(html)
}

// SaveFile saves the content to the current file
func (e *Editor) SaveFile() bool {
	// If no file is currently open, prompt for a location
	if e.currentFilePath == "" {
		return e.SaveFileAs()
	}

	// Save to the existing file
	err := e.saveToFile(e.currentFilePath)
	if err != nil {
		runtime.EventsEmit(e.ctx, "error", "Failed to save file: "+err.Error())
		return false
	}

	e.lastSaveTime = time.Now()
	runtime.EventsEmit(e.ctx, "status:update", "File saved")
	return true
}

// SaveFileAs prompts for a filename and saves the content
func (e *Editor) SaveFileAs() bool {
	// Show file dialog
	filePath, err := runtime.SaveFileDialog(e.ctx, runtime.SaveDialogOptions{
		DefaultDirectory: "",
		DefaultFilename:  "untitled.md",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Markdown Files (*.md, *.markdown)",
				Pattern:     "*.md;*.markdown",
			},
			{
				DisplayName: "All Files (*.*)",
				Pattern:     "*.*",
			},
		},
	})

	if err != nil || filePath == "" {
		// User cancelled
		return false
	}

	err = e.saveToFile(filePath)
	if err != nil {
		runtime.EventsEmit(e.ctx, "error", "Failed to save file: "+err.Error())
		return false
	}

	e.currentFilePath = filePath
	e.lastSaveTime = time.Now()
	runtime.EventsEmit(e.ctx, "status:update", "File saved")

	// Update window title with filename
	runtime.WindowSetTitle(e.ctx, "Markdown Editor - "+e.getFilenameFromPath())

	return true
}

// OpenFile opens a markdown file
func (e *Editor) OpenFile() bool {
	// Show file dialog
	filePath, err := runtime.OpenFileDialog(e.ctx, runtime.OpenDialogOptions{
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Markdown Files (*.md, *.markdown)",
				Pattern:     "*.md;*.markdown",
			},
			{
				DisplayName: "All Files (*.*)",
				Pattern:     "*.*",
			},
		},
	})

	if err != nil || filePath == "" {
		// User cancelled
		return false
	}

	// Read file content
	content, err := e.readFromFile(filePath)
	if err != nil {
		runtime.EventsEmit(e.ctx, "error", "Failed to open file: "+err.Error())
		return false
	}

	e.currentFilePath = filePath
	e.SetContent(content)
	e.lastSaveTime = time.Now()

	// Update window title with filename
	runtime.WindowSetTitle(e.ctx, "Markdown Editor - "+e.getFilenameFromPath())

	runtime.EventsEmit(e.ctx, "status:update", "File opened")
	return true
}

// AutoSave automatically saves the file if changes exist
func (e *Editor) AutoSave() {
	if e.currentFilePath != "" && e.hasUnsavedChanges() {
		e.SaveFile()
		runtime.EventsEmit(e.ctx, "status:update", "Auto-saved")
	}
}

// NewFile creates a new file
func (e *Editor) NewFile() {
	// Check for unsaved changes
	if e.hasUnsavedChanges() {
		// Prompt user to save changes
		// This will be handled via a frontend dialog
	}

	e.currentFilePath = ""
	e.SetContent("")
	runtime.WindowSetTitle(e.ctx, "Markdown Editor - Untitled")
}

// ToggleDarkMode switches between light and dark mode
func (e *Editor) ToggleDarkMode() {
	e.isDarkMode = !e.isDarkMode
	runtime.EventsEmit(e.ctx, "theme:update", e.isDarkMode)
}

// ToggleAutoSave enables or disables autosave functionality
func (e *Editor) ToggleAutoSave() bool {
	e.autoSaveEnabled = !e.autoSaveEnabled
	return e.autoSaveEnabled
}

// SetAutoSaveDelay sets the autosave delay in seconds
func (e *Editor) SetAutoSaveDelay(seconds int) {
	e.autoSaveDelay = time.Duration(seconds) * time.Second
}

// Helper methods
func (e *Editor) hasUnsavedChanges() bool {
	return e.lastEditTime.After(e.lastSaveTime)
}

func (e *Editor) getFilenameFromPath() string {
	if e.currentFilePath == "" {
		return "Untitled"
	}
	return e.fileUtils.GetFilenameFromPath(e.currentFilePath)
}

func (e *Editor) saveToFile(path string) error {
	return e.fileUtils.SaveToFile(path, e.content)
}

func (e *Editor) readFromFile(path string) (string, error) {
	return e.fileUtils.ReadFromFile(path)
}

// GetCurrentFilePath returns the path of the currently open file
func (e *Editor) GetCurrentFilePath() string {
	return e.currentFilePath
}

// SetAutoSaveEnabled enables or disables autosave
func (e *Editor) SetAutoSaveEnabled(enabled bool) {
	e.autoSaveEnabled = enabled
}

// GetAutoSaveEnabled returns whether autosave is enabled
func (e *Editor) GetAutoSaveEnabled() bool {
	return e.autoSaveEnabled
}
