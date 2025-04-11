package main

import (
	// ... existing imports ...
	"html/template"
	"log"
	"net/http"

	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

var formTemplate = template.Must(template.New("form").Parse(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>Monogram Test Interface</title>
</head>
<body>
  <h1>Monogram Test Interface</h1>
  <form action="/translate" method="post">
    <label for="monogramInput">Monogram Notation:</label><br>
    <textarea name="monogramInput" id="monogramInput" rows="10" cols="80"></textarea><br><br>
    
    <label for="format">Output Format:</label>
    <select name="format" id="format">
       <option value="xml">XML</option>
       <option value="json">JSON</option>
       <option value="yaml">YAML</option>
       <option value="mermaid">Mermaid</option>
       <option value="dot">DOT</option>
    </select>
    <br><br>
    
    <label for="indent">Indent (number of spaces):</label>
    <input type="number" id="indent" name="indent" value="2" min="0"><br><br>
    
    <label for="defaultBreaker">Default Breaker:</label>
    <input type="text" id="defaultBreaker" name="defaultBreaker" value="_"><br><br>
    
    <label for="includeSpans">Include Spans:</label>
    <input type="checkbox" id="includeSpans" name="includeSpans"><br><br>
    
    <input type="submit" value="Translate">
  </form>
  <br>
  <div>
    <h2>Output:</h2>
    <pre>{{.Output}}</pre>
  </div>
</body>
</html>
`))

// startTestServer starts an HTTP listener on the specified port and opens the browser.
func startTestServer(port string) {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/translate", translateHandler)

	if port == "" {
		port = "3000"
	}
	addr := "localhost:" + port
	go openBrowser("http://" + addr)
	log.Printf("Starting test server on %s...", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start test server: %v", err)
	}
}

// indexHandler renders the form page without translation output.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	formTemplate.Execute(w, struct{ Output string }{Output: ""})
}

// translateHandler processes the form and renders the translation output.
func translateHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}
	monogramInput := r.FormValue("monogramInput")
	format := r.FormValue("format")
	indentVal := r.FormValue("indent")
	defaultBreaker := r.FormValue("defaultBreaker")
	includeSpans := r.FormValue("includeSpans") == "on"

	// Convert indent value to integer:
	indent := 2
	if indentParsed, err := strconv.Atoi(indentVal); err == nil {
		indent = indentParsed
	}

	// Set up FormatOptions based on the form values:
	options := FormatOptions{
		Format:       format,
		Input:        "", // Not used in test mode — we’re using form data.
		Output:       "", // Output will be captured in a buffer.
		Indent:       indent,
		Limit:        false,
		UnglueOption: defaultBreaker,
		IncludeSpans: includeSpans,
	}

	// Look up the translator function:
	translator, ok := formatHandlers[format]
	if !ok {
		http.Error(w, "Unknown format: "+format, http.StatusBadRequest)
		return
	}

	// Create reader from input text and a bytes.Buffer for capturing output:
	inputReader := strings.NewReader(monogramInput)
	var outputBuffer bytes.Buffer

	// Perform the translation.
	translator(inputReader, &outputBuffer, &options)

	// Render the same form with the translation output shown:
	formTemplate.Execute(w, struct{ Output string }{Output: outputBuffer.String()})
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Printf("Failed to open browser: %v", err)
	}
}
