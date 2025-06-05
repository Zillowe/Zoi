package commands

import (
	"fmt"
	"gct/src/ai"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	stateInitName = iota
	stateInitProvider
	stateInitModel
	stateInitAPIKey
	stateInitGuides
	stateSubmit
)

var (
	titleStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	promptStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	selectedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	checkmarkStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("46"))
)

type InitTUIModel struct {
	currentState   int
	inputs         []textinput.Model
	providerCursor int
	width          int

	Name     string
	Provider string
	Model    string
	APIKey   string
	Guides   string

	quitting  bool
	submitted bool
}

func NewInitTUIModel() InitTUIModel {
	inputs := make([]textinput.Model, 4)

	inputs[0] = textinput.New()
	inputs[0].Placeholder = "My Awesome Project"
	inputs[0].Focus()
	inputs[0].CharLimit = 50
	inputs[0].Width = 50

	inputs[1] = textinput.New()
	inputs[1].Placeholder = "e.g., gpt-4o, claude-3-haiku-20240307"
	inputs[1].CharLimit = 100
	inputs[1].Width = 50

	inputs[2] = textinput.New()
	inputs[2].Placeholder = "sk-..."
	inputs[2].EchoMode = textinput.EchoPassword
	inputs[2].EchoCharacter = '•'
	inputs[2].CharLimit = 128
	inputs[2].Width = 50

	inputs[3] = textinput.New()
	inputs[3].Placeholder = "./path/to/guide.md (optional)"
	inputs[3].CharLimit = 256
	inputs[3].Width = 50

	return InitTUIModel{
		currentState: stateInitName,
		inputs:       inputs,
	}
}

func (m InitTUIModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m InitTUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			return m, tea.Quit

		case tea.KeyEnter:
			switch m.currentState {
			case stateInitName:
				m.Name = m.inputs[0].Value()
				m.currentState = stateInitProvider
				m.inputs[0].Blur()
				return m, nil
			case stateInitProvider:
				m.Provider = ai.SupportedProviders[m.providerCursor]
				m.currentState = stateInitModel
				m.inputs[1].Focus()
				return m, textinput.Blink
			case stateInitModel:
				m.Model = m.inputs[1].Value()
				m.currentState = stateInitAPIKey
				m.inputs[1].Blur()
				m.inputs[2].Focus()
				return m, textinput.Blink
			case stateInitAPIKey:
				m.APIKey = m.inputs[2].Value()
				m.currentState = stateInitGuides
				m.inputs[2].Blur()
				m.inputs[3].Focus()
				return m, textinput.Blink
			case stateInitGuides:
				m.Guides = m.inputs[3].Value()
				m.currentState = stateSubmit
				m.submitted = true
				m.quitting = true
				return m, tea.Quit
			}

		case tea.KeyUp, tea.KeyShiftTab:
			if m.currentState == stateInitProvider {
				if m.providerCursor > 0 {
					m.providerCursor--
				}
			}
		case tea.KeyDown, tea.KeyTab:
			if m.currentState == stateInitProvider {
				if m.providerCursor < len(ai.SupportedProviders)-1 {
					m.providerCursor++
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
	}

	switch m.currentState {
	case stateInitName:
		m.inputs[0], cmd = m.inputs[0].Update(msg)
	case stateInitModel:
		m.inputs[1], cmd = m.inputs[1].Update(msg)
	case stateInitAPIKey:
		m.inputs[2], cmd = m.inputs[2].Update(msg)
	case stateInitGuides:
		m.inputs[3], cmd = m.inputs[3].Update(msg)
	}

	return m, cmd
}

func (m InitTUIModel) View() string {
	if m.quitting {
		return ""
	}
	var s strings.Builder

	s.WriteString(titleStyle.Render("Welcome to GCT! Let's set up your project.") + "\n")
	s.WriteString(promptStyle.Render("Fill in the details below. Press Enter to confirm each step.") + "\n\n")

	if m.currentState > stateInitName {
		s.WriteString(fmt.Sprintf("%s Project Name: %s\n", checkmarkStyle.Render("✓"), m.Name))
	} else {
		s.WriteString("Project Name:\n" + m.inputs[0].View() + "\n")
	}

	if m.currentState > stateInitProvider {
		s.WriteString(fmt.Sprintf("%s AI Provider: %s\n", checkmarkStyle.Render("✓"), m.Provider))
	} else if m.currentState == stateInitProvider {
		s.WriteString("Select AI Provider (Use ↑/↓):\n")
		for i, choice := range ai.SupportedProviders {
			cursor := "  "
			if m.providerCursor == i {
				cursor = selectedStyle.Render("> ")
			}
			s.WriteString(fmt.Sprintf("%s%s\n", cursor, choice))
		}
	}

	if m.currentState > stateInitModel {
		s.WriteString(fmt.Sprintf("%s AI Model: %s\n", checkmarkStyle.Render("✓"), m.Model))
	} else if m.currentState == stateInitModel {
		s.WriteString("\nAI Model Name:\n" + m.inputs[1].View() + "\n")
	}

	if m.currentState > stateInitAPIKey {
		s.WriteString(fmt.Sprintf("%s API Key: %s\n", checkmarkStyle.Render("✓"), "[hidden]"))
	} else if m.currentState == stateInitAPIKey {
		s.WriteString("\nAPI Key:\n" + m.inputs[2].View() + "\n")
	}

	if m.currentState > stateInitGuides {
		guidesDisplay := m.Guides
		if guidesDisplay == "" {
			guidesDisplay = "None"
		}
		s.WriteString(fmt.Sprintf("%s Commit Guides: %s\n", checkmarkStyle.Render("✓"), guidesDisplay))
	} else if m.currentState == stateInitGuides {
		s.WriteString("\nCommit Guides (optional, separate multiple paths with a space):\n" + m.inputs[3].View() + "\n")
	}

	s.WriteString(promptStyle.Render("\nPress Esc or Ctrl+C to quit at any time."))
	return s.String()
}
