package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	textInput  textinput.Model
	err        error
	config     []map[string]string
	createMode bool
	fileName   string
}

func (m *model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "c":
			m.createMode = true
			m.textInput.Focus()
		case "s":
			err := saveConfig(m.config, m.fileName)
			if err != nil {
				log.Fatal(err)
			}
		case "o":
			err := openConfig(&m.config, m.fileName)
			if err != nil {
				log.Fatal(err)
			}
		case "a":
			addKeyValuePair(&m.config)
		case "l":
			listKeys(m.config)
		case "r":
			removeKey(&m.config)
		case "q":
			return nil, tea.Quit
		}
	// Remove the unnecessary import statement

	// ...

	func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "c":
				m.createMode = true
				m.textInput.Focus()
			case "s":
				err := saveConfig(m.config, m.fileName)
				if err != nil {
					log.Fatal(err)
				}
			case "o":
				err := openConfig(&m.config, m.fileName)
				if err != nil {
					log.Fatal(err)
				}
			case "a":
				addKeyValuePair(&m.config)
			case "l":
				listKeys(m.config)
			case "r":
				removeKey(&m.config)
			case "q":
				return nil, tea.Quit
			}
		case tea.KeyRunes:
			if m.createMode {
				m.fileName = string(msg)
				m.createMode = false
			}
		}
		return m.(tea.Model), nil

	}


	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func initialModel() model {
	ti := textinput.NewModel()
	ti.Placeholder = "Pikachu"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		textInput:  ti,
		err:        nil,
		config:     []map[string]string{},
		createMode: false,
		fileName:   "",
	}
}

func (m *model) View() string {
	s := "My Config App\n\n" // Add the header here
	s += "Config:\n"
	for _, section := range m.config {
		s += "  "
		for key, value := range section {
			s += key + ": " + value + "\n"
		}
	}
	s += "\nCommands:\n"
	s += "  c: Create new config\n"
	s += "  s: Save config\n"
	s += "  o: Open config\n"
	s += "  a: Add new key-value pair\n"
	s += "  l: List all config keys\n"
	s += "  r: Remove config key\n"
	s += "  q: Quit\n"

	if m.createMode {
		s += "\nPlease enter the name of your config file: " + m.textInput.Value()
	}

	return s
}

func addKeyValuePair(config *[]map[string]string) {
	fmt.Println("Enter a new key-value pair:")
	var key, value string
	fmt.Print("Key: ")
	fmt.Scan(&key)
	fmt.Print("Value: ")
	fmt.Scan(&value)
	(*config)[len(*config)-1][key] = value
}

func removeKey(config *[]map[string]string) {
	fmt.Println("Enter the key to remove:")
	var key string
	fmt.Print("Key: ")
	fmt.Scan(&key)
	for i, section := range *config {
		if _, ok := section[key]; ok {
			delete((*config)[i], key)
			break
		}
	}
}

func saveConfig(config []map[string]string, fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, section := range config {
		for key, value := range section {
			fmt.Fprintf(f, "%s: %s\n", key, value)
		}
	}
	return nil
}

func openConfig(config *[]map[string]string, fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			return fmt.Errorf("invalid config file")
		}
		(*config) = append((*config), map[string]string{
			strings.TrimSpace(parts[0]): strings.TrimSpace(parts[1]),
		})
	}
	return nil
}

func listKeys(config []map[string]string) {
	for _, section := range config {
		for key, value := range section {
			fmt.Println(key, ":", value)
		}
	}
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
