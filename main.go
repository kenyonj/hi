package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Content struct {
	Commands  map[string]string `yaml:"commands"`
	ResumeURL string            `yaml:"resume_url"`
}

var content Content

func loadContent() error {
	data, err := os.ReadFile("content.yaml")
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, &content)
}

type CommandResponse struct {
	Command   string
	Output    template.HTML
	IsResume  bool
	ResumeURL string
}

func main() {
	if err := loadContent(); err != nil {
		log.Fatal("Failed to load content.yaml:", err)
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/command", commandHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func commandHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cmd := strings.ToLower(strings.TrimSpace(r.FormValue("command")))
	if cmd == "" {
		return
	}

	if cmd == "clear" {
		w.Header().Set("HX-Reswap", "innerHTML")
		w.Header().Set("HX-Retarget", "#terminal-output")
		return
	}

	output, exists := content.Commands[cmd]
	if !exists {
		output = "Command not found: " + cmd + ". Type \"help\" for available commands."
	}

	// Replace newlines with <br> for HTML
	htmlOutput := template.HTML(strings.ReplaceAll(output, "\n", "<br>"))

	resp := CommandResponse{
		Command:   cmd,
		Output:    htmlOutput,
		IsResume:  cmd == "resume",
		ResumeURL: content.ResumeURL,
	}

	tmpl := template.Must(template.ParseFiles("templates/command.html"))
	tmpl.Execute(w, resp)
}
