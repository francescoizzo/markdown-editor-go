// Import Wails runtime
// No explicit import needed; Wails runtime is injected automatically

// Global variables
let editor; // Monaco editor instance
let editorValue = "";
let isDarkMode = false;
let wordCount = 0;
let autoSaveEnabled = true;
let hasUnsavedChanges = false;
let editorChangeTimeout;

// Initialize application when DOM content is loaded
document.addEventListener("DOMContentLoaded", () => {
  // In Wails v2, the runtime is already initialized and available as window.runtime
  console.log("DOM content loaded, initializing app");
  initializeApp();
});

// Main application initialization
async function initializeApp() {
  // Initialize Monaco Editor
  initMonacoEditor();

  // Initialize event listeners
  setupEventListeners();

  // Initialize theme based on system preference or saved preference
  const prefersDarkMode =
    window.matchMedia &&
    window.matchMedia("(prefers-color-scheme: dark)").matches;
  setTheme(prefersDarkMode);

  // Set up auto save indicator
  updateAutoSaveIndicator(autoSaveEnabled);

  // Register Wails event listeners
  registerWailsEvents();
}

// Initialize Monaco Editor
function initMonacoEditor() {
  require(["vs/editor/editor.main"], () => {
    // Define Markdown language configuration
    monaco.languages.register({ id: "markdown" });

    // Create and configure Monaco editor
    editor = monaco.editor.create(document.getElementById("editor-pane"), {
      value: "",
      language: "markdown",
      theme: isDarkMode ? "vs-dark" : "vs",
      lineNumbers: "on",
      wordWrap: "on",
      wrappingIndent: "same",
      minimap: { enabled: true },
      scrollBeyondLastLine: false,
      automaticLayout: true,
      fontSize: 14,
      fontFamily: "'Roboto Mono', monospace",
      tabSize: 4,
      renderWhitespace: "boundary",
      contextmenu: true,
      folding: true,
      lineDecorationsWidth: 10,
      renderLineHighlight: "all",
    });

    // Set up editor change event
    editor.onDidChangeModelContent(() => {
      editorValue = editor.getValue();
      hasUnsavedChanges = true;

      // Update word count
      updateWordCount(editorValue);

      // Schedule content update
      clearTimeout(editorChangeTimeout);
      editorChangeTimeout = setTimeout(() => {
        window.go.main.Editor.SetContent(editorValue);
      }, 300);
    });

    // Initial focus on editor
    editor.focus();
  });
}

// Set up event listeners for UI elements
function setupEventListeners() {
  // File operations
  document.getElementById("btn-new").addEventListener("click", newFile);
  document.getElementById("btn-open").addEventListener("click", openFile);
  document.getElementById("btn-save").addEventListener("click", saveFile);
  document.getElementById("btn-save-as").addEventListener("click", saveFileAs);

  // Theme toggle
  document
    .getElementById("btn-theme-toggle")
    .addEventListener("click", toggleTheme);

  // Keyboard shortcuts
  document.addEventListener("keydown", handleKeyboardShortcuts);
}

// Register Wails event listeners
function registerWailsEvents() {
  // Handle content updates
  window.runtime.EventsOn("content:update", (html) => {
    document.getElementById(
      "preview-pane"
    ).innerHTML = `<div class="markdown-preview">${html}</div>`;
    highlightCodeBlocks();
  });

  // Handle status updates
  window.runtime.EventsOn("status:update", (message) => {
    document.getElementById("status-message").textContent = message;
    setTimeout(() => {
      document.getElementById("status-message").textContent = "Ready";
    }, 3000);
  });

  // Handle errors
  window.runtime.EventsOn("error", (message) => {
    console.error("Error:", message);
    document.getElementById("status-message").textContent = `Error: ${message}`;
  });

  // Handle theme updates
  window.runtime.EventsOn("theme:update", (darkMode) => {
    setTheme(darkMode);
  });
}

// File operations
function newFile() {
  if (hasUnsavedChanges) {
    if (!confirm("You have unsaved changes. Do you want to continue?")) {
      return;
    }
  }

  window.go.main.MainWindow.NewFile();
  editor.setValue("");
  hasUnsavedChanges = false;
  updateWordCount("");
}

function openFile() {
  if (hasUnsavedChanges) {
    if (!confirm("You have unsaved changes. Do you want to continue?")) {
      return;
    }
  }

  window.go.main.MainWindow.OpenFile().then((success) => {
    if (success) {
      window.go.main.MainWindow.GetContent().then((content) => {
        editor.setValue(content);
        hasUnsavedChanges = false;
        updateWordCount(content);
      });
    }
  });
}

function saveFile() {
  window.go.main.MainWindow.SaveFile().then((success) => {
    if (success) {
      hasUnsavedChanges = false;
    }
  });
}

function saveFileAs() {
  window.go.main.MainWindow.SaveFileAs().then((success) => {
    if (success) {
      hasUnsavedChanges = false;
    }
  });
}

// Theme operations
function toggleTheme() {
  isDarkMode = !isDarkMode;
  setTheme(isDarkMode);
  window.go.main.MainWindow.ToggleTheme();
}

function setTheme(darkMode) {
  isDarkMode = darkMode;

  // Update body class
  if (darkMode) {
    document.body.classList.add("dark-theme");
  } else {
    document.body.classList.remove("dark-theme");
  }

  // Update editor theme
  if (editor) {
    monaco.editor.setTheme(darkMode ? "vs-dark" : "vs");
  }

  // Update theme icons
  document.getElementById("theme-icon-light").style.display = darkMode
    ? "none"
    : "inline";
  document.getElementById("theme-icon-dark").style.display = darkMode
    ? "inline"
    : "none";
  document.getElementById("theme-text").textContent = darkMode
    ? "Light Mode"
    : "Dark Mode";

  // If highlight.js is loaded, update code highlighting theme
  const highlightTheme = document.querySelector('link[href*="highlight.js"]');
  if (highlightTheme) {
    highlightTheme.href = darkMode
      ? "https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.7.0/styles/atom-one-dark.min.css"
      : "https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.7.0/styles/github.min.css";
  }
}

// Autosave operations
function toggleAutoSave() {
  window.go.main.MainWindow.ToggleAutoSave().then((enabled) => {
    autoSaveEnabled = enabled;
    updateAutoSaveIndicator(enabled);
  });
}

function updateAutoSaveIndicator(enabled) {
  const indicator = document.getElementById("autosave-indicator");
  indicator.querySelector(".indicator-text").textContent = `Autosave: ${
    enabled ? "On" : "Off"
  }`;
}

// Word count
function updateWordCount(text) {
  if (!text || text.trim() === "") {
    wordCount = 0;
  } else {
    // Count words by splitting on whitespace
    const words = text.trim().split(/\s+/);
    wordCount = words.length;
  }

  document
    .getElementById("wordcount-indicator")
    .querySelector(".indicator-text").textContent = `${wordCount} ${
    wordCount === 1 ? "word" : "words"
  }`;
}

// Keyboard shortcuts
function handleKeyboardShortcuts(event) {
  // Ctrl+S or Command+S: Save
  if ((event.ctrlKey || event.metaKey) && event.key === "s") {
    event.preventDefault();
    if (event.shiftKey) {
      saveFileAs();
    } else {
      saveFile();
    }
  }

  // Ctrl+N or Command+N: New File
  if ((event.ctrlKey || event.metaKey) && event.key === "n") {
    event.preventDefault();
    newFile();
  }

  // Ctrl+O or Command+O: Open File
  if ((event.ctrlKey || event.metaKey) && event.key === "o") {
    event.preventDefault();
    openFile();
  }
}

// Syntax highlighting for code blocks in preview
function highlightCodeBlocks() {
  document.querySelectorAll(".markdown-preview pre code").forEach((block) => {
    hljs.highlightBlock(block);
  });
}
