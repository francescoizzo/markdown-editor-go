package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"github.com/francescoizzo/markdown-editor-go/internal/ui"
)

var assets embed.FS

func main() {
	// Create a new instance of the MainWindow
	mainWindow := ui.NewMainWindow()

	// Create application with options
	err := wails.Run(&options.App{
		Title:             "Markdown Editor",
		Width:             1024,
		Height:            768,
		MinWidth:          800,
		MinHeight:         600,
		DisableResize:     false,
		Fullscreen:        false,
		StartHidden:       false,
		HideWindowOnClose: false,
		BackgroundColour:  &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Menu:             nil,
		Logger:           nil,
		LogLevel:         0,
		OnStartup:        mainWindow.OnStartup,
		OnDomReady:       mainWindow.OnDomReady,
		OnBeforeClose:    mainWindow.OnBeforeClose,
		OnShutdown:       mainWindow.OnShutdown,
		WindowStartState: options.Normal,
		Bind: []interface{}{
			mainWindow,
		},
		// Windows specific configuration
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
		},
		// MacOS specific configuration
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: true,
				HideTitle:                  false,
				HideTitleBar:               false,
				FullSizeContent:            false,
				UseToolbar:                 false,
				HideToolbarSeparator:       true,
			},
			WebviewIsTransparent: true,
			WindowIsTranslucent:  false,
			About: &mac.AboutInfo{
				Title:   "Markdown Editor",
				Message: "A modern Markdown editor built with Go and Wails",
			},
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
