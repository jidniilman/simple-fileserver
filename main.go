package main

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// CSS content as a constant
const cssStyles = `
body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
    line-height: 1.6;
    color: #333;
    max-width: 800px;
    margin: 0 auto;
    padding: 20px;
    background-color: #f5f5f5;
}

.container {
    background: white;
    border-radius: 8px;
    padding: 2rem;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

h1 {
    color: #2c3e50;
    margin-top: 0;
    padding-bottom: 1rem;
    border-bottom: 1px solid #eee;
}

.file-list {
    margin-top: 1.5rem;
}

.file-item {
    padding: 0.75rem 1rem;
    margin: 0.25rem 0;
    border-radius: 4px;
    transition: background-color 0.2s;
}

.file-item:hover {
    background-color: #f8f9fa;
}

.file-link {
    text-decoration: none;
    color: #3498db;
    display: flex;
    align-items: center;
}

.file-link:hover {
    text-decoration: underline;
}

.file-size {
    margin-left: auto;
    color: #7f8c8d;
    font-size: 0.9em;
}

.directory .file-link {
    color: #2ecc71;
    font-weight: 500;
}

/* Responsive design */
@media (max-width: 600px) {
    body {
        padding: 10px;
    }
    
    .container {
        padding: 1rem;
    }
    
    .file-item {
        padding: 0.5rem;
    }
}
`

type FileInfo struct {
	Name  string
	Path  string
	IsDir bool
	Size  int64
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Serve CSS directly
	e.GET("/static/styles.css", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "text/css")
		return c.String(http.StatusOK, cssStyles)
	})

	// Routes
	e.GET("/*", handleFileListing)

	// Start server
	e.Logger.Fatal(e.Start(":4221"))
}

func handleFileListing(c echo.Context) error {
	// Get the requested path
	requestedPath := c.Request().URL.Path
	if requestedPath == "/" {
		requestedPath = "."
	}

	// Get the absolute path
	absPath, err := filepath.Abs(".")
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error getting absolute path")
	}

	// Join with the requested path
	fullPath := filepath.Join(absPath, requestedPath)

	// Check if the path exists
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return c.String(http.StatusNotFound, "File not found")
		}
		return c.String(http.StatusInternalServerError, "Error reading file")
	}

	// If it's a file, serve it
	if !fileInfo.IsDir() {
		return c.File(fullPath)
	}

	// If it's a directory, list its contents
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error reading directory")
	}

	var files []FileInfo
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		files = append(files, FileInfo{
			Name:  entry.Name(),
			Path:  filepath.Join(requestedPath, entry.Name()),
			IsDir: entry.IsDir(),
			Size:  info.Size(),
		})
	}

	// Define the template with proper data structure
	tmpl := template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>File Server</title>
    <style>{{.CSS}}</style>
</head>
<body>
    <div class="container">
        <h1>File Server</h1>
        <div class="file-list">
            {{range .Files}}
                <div class="file-item {{if .IsDir}}directory{{else}}file{{end}}">
                    <a href="{{.Path}}" class="file-link">
                        {{if .IsDir}}📁{{else}}📄{{end}} {{.Name}}
                    </a>
                    {{if not .IsDir}}
                        <span class="file-size">({{.Size}} bytes)</span>
                    {{end}}
                </div>
            {{end}}
        </div>
    </div>
</body>
</html>
	`))

	// Create data structure for template
	data := struct {
		Files []FileInfo
		CSS   template.CSS
	}{
		Files: files,
		CSS:   template.CSS(cssStyles),
	}
	
	// Execute template with the data
	return tmpl.Execute(c.Response(), data)
}
