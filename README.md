# Postty

A terminal-based REST API client inspired by lazygit.

## Features

- **Beautiful TUI** - Clean interface using Bubbletea + Lipgloss
- **Keyboard-Driven** - Navigate with numbers 1-5, vim-style j/k, and arrow keys
- **Fast & Lightweight** - Built with Go, instant startup
- **Full HTTP Support** - GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS
- **Multiple Content Types** - JSON, XML, plain text, form-urlencoded, multipart
- **Auto-Formatting** - Automatic JSON pretty-printing
- **Status Indicators** - Color-coded HTTP status codes

## Quick Start

```bash
# Build and run
go build -o postty
./postty
```

## Usage

### Keybindings

| Key | Action |
|-----|--------|
| `Tab` | Cycle to next pane |
| `Shift+Tab` | Cycle to previous pane |
| `1` | Jump to URL pane (from Method/Header/Response) |
| `2` | Jump to Method pane (from Method/Header/Response) |
| `3` | Jump to Body pane (from Method/Header/Response) |
| `4` | Jump to Content-Type pane (from Method/Header/Response) |
| `5` | Jump to Response pane (from Method/Header/Response) |
| `↑/↓` or `j/k` | Navigate lists (in Method/Content-Type) / Scroll (in Result) |
| `PgUp/PgDown` | Scroll half page (in Result pane) |
| `Home/End` or `g/G` | Jump to top/bottom (in Result pane) |
| `Enter` | Send request (or new line in Body pane) |
| `Alt+Enter` | Send request from Body pane |
| `Esc` | Quit (from any pane) |
| `q` | Quit (from Method/Header/Response only) |
| `Ctrl+C` | Quit (from any pane) |

**Notes:**
- When typing in URL or Body panes, all characters (including q and 1-5) are typed into the input. Use `Tab` to navigate between panes while in text input fields.
- When viewing large responses in the Result pane, use arrow keys or j/k to scroll through the content
- In the Body pane, press `Enter` for new lines and `Alt+Enter` to send the request

### Example: Making a GET Request

1. Press `1` or `Tab` to focus URL pane (default pane on startup)
2. Type: `https://jsonplaceholder.typicode.com/posts/1`
3. Press `Enter`
4. View formatted response in Result pane

### Example: Making a POST Request

1. Press `Tab` → Navigate to Method pane
2. Select POST with `↓`
3. Press `Tab` or `1` → Focus URL pane
4. Enter URL: `https://jsonplaceholder.typicode.com/posts`
5. Press `Tab` twice → Navigate to Body pane
6. Enter body:
   ```json
   {
     "title": "Test",
     "body": "Hello",
     "userId": 1
   }
   ```
7. Press `Tab` → Navigate to Content-Type
8. Ensure `application/json` is selected
9. Press `Enter` to send

## Testing

```bash
# Run tests
go test -v

# Run with coverage
go test -v -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run integration tests
go test -tags=integration -v

```
## Building

```bash
# Build
go build -o postty

# Install to $GOPATH/bin
go install

# Install to custom location
go build -o /usr/local/bin/postty
```

## Inspiration

- [lazygit](https://github.com/jesseduffield/lazygit) - UI/UX design
- [Postman](https://www.postman.com/) - Feature set
