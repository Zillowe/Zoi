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
	stateInitEndpoint
	stateInitGCPProjectID
	stateInitGCPRegion
	stateInitAWSRegion
	stateInitAWSAccessKeyID
	stateInitAWSSecretAccessKey
	stateInitAzureResourceName
	stateInitModel
	stateInitAPIKey
	stateInitCommitGuides
	stateInitChangelogGuides
	stateSubmit
)

var (
	titleStyleInit     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	promptStyleInit    = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	selectedStyleInit  = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	checkmarkStyleInit = lipgloss.NewStyle().Foreground(lipgloss.Color("46"))
)

type InitTUIModel struct {
	currentState   int
	inputs         []textinput.Model
	providerCursor int
	width          int

	Name            string
	Provider        string
	Endpoint        string
	Model           string
	APIKey          string
	CommitGuides    string
	ChangelogGuides string

	GCPProjectID string
	GCPRegion    string

	AWSRegion          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string

	AzureResourceName string

	quitting  bool
	submitted bool
}

func NewInitTUIModel() InitTUIModel {
	inputs := make([]textinput.Model, 12)

	inputs[0] = textinput.New()
	inputs[0].Placeholder = "My Awesome Project"
	inputs[0].Focus()
	inputs[0].CharLimit = 50
	inputs[0].Width = 50

	inputs[1] = textinput.New()
	inputs[1].Placeholder = "e.g. gpt-4o, claude-3-haiku-20240307"
	inputs[1].CharLimit = 100
	inputs[1].Width = 50

	inputs[2] = textinput.New()
	inputs[2].Placeholder = "sk-..."
	inputs[2].EchoMode = textinput.EchoPassword
	inputs[2].EchoCharacter = '•'
	inputs[2].CharLimit = 128
	inputs[2].Width = 50

	inputs[3] = textinput.New()
	inputs[3].Placeholder = "./path/to/commits.md (optional)"
	inputs[3].CharLimit = 256
	inputs[3].Width = 50

	inputs[4] = textinput.New()
	inputs[4].Placeholder = "https://api.example.com/v1"
	inputs[4].CharLimit = 256
	inputs[4].Width = 50

	inputs[5] = textinput.New()
	inputs[5].Placeholder = "./path/to/changelogs.md (optional)"
	inputs[5].CharLimit = 256
	inputs[5].Width = 50

	inputs[6] = textinput.New()
	inputs[6].Placeholder = "gct-project-12345"
	inputs[6].CharLimit = 100
	inputs[6].Width = 50

	inputs[7] = textinput.New()
	inputs[7].Placeholder = "us-central1"
	inputs[7].CharLimit = 50
	inputs[7].Width = 50

	inputs[8] = textinput.New()
	inputs[8].Placeholder = "us-east-1"
	inputs[8].CharLimit = 50
	inputs[8].Width = 50

	inputs[9] = textinput.New()
	inputs[9].Placeholder = "AKIA..."
	inputs[9].EchoMode = textinput.EchoPassword
	inputs[9].EchoCharacter = '•'
	inputs[9].CharLimit = 128
	inputs[9].Width = 50

	inputs[10] = textinput.New()
	inputs[10].Placeholder = "Your very long secret key"
	inputs[10].EchoMode = textinput.EchoPassword
	inputs[10].EchoCharacter = '•'
	inputs[10].CharLimit = 128
	inputs[10].Width = 50

	inputs[11] = textinput.New()
	inputs[11].Placeholder = "my-openai-resource"
	inputs[11].CharLimit = 100
	inputs[11].Width = 50

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
	var cmds []tea.Cmd

	providerID := strings.ToLower(strings.ReplaceAll(m.Provider, " ", ""))

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
				providerID = strings.ToLower(strings.ReplaceAll(m.Provider, " ", ""))

				switch providerID {
				case "openaicompatible":
					m.currentState = stateInitEndpoint
					m.inputs[4].Focus()
				case "googlevertexai", "vertexai", "vertex":
					m.currentState = stateInitGCPProjectID
					m.inputs[6].Focus()
				case "amazonbedrock", "bedrock", "amazon", "aws":
					m.currentState = stateInitAWSRegion
					m.inputs[8].Focus()
				case "azureopenai":
					m.currentState = stateInitAzureResourceName
					m.inputs[11].Focus()
				default:
					m.currentState = stateInitModel
					m.inputs[1].Focus()
				}
				return m, textinput.Blink

			case stateInitEndpoint:
				m.Endpoint = m.inputs[4].Value()
				m.currentState = stateInitModel
				m.inputs[4].Blur()
				m.inputs[1].Focus()
				return m, textinput.Blink

			case stateInitGCPProjectID:
				m.GCPProjectID = m.inputs[6].Value()
				m.currentState = stateInitGCPRegion
				m.inputs[6].Blur()
				m.inputs[7].Focus()
				return m, textinput.Blink

			case stateInitGCPRegion:
				m.GCPRegion = m.inputs[7].Value()
				m.currentState = stateInitModel
				m.inputs[7].Blur()
				m.inputs[1].Focus()
				return m, textinput.Blink

			case stateInitAWSRegion:
				m.AWSRegion = m.inputs[8].Value()
				m.currentState = stateInitAWSAccessKeyID
				m.inputs[8].Blur()
				m.inputs[9].Focus()
				return m, textinput.Blink

			case stateInitAWSAccessKeyID:
				m.AWSAccessKeyID = m.inputs[9].Value()
				m.currentState = stateInitAWSSecretAccessKey
				m.inputs[9].Blur()
				m.inputs[10].Focus()
				return m, textinput.Blink

			case stateInitAWSSecretAccessKey:
				m.AWSSecretAccessKey = m.inputs[10].Value()
				m.currentState = stateInitModel
				m.inputs[10].Blur()
				m.inputs[1].Focus()
				return m, textinput.Blink

			case stateInitAzureResourceName:
				m.AzureResourceName = m.inputs[11].Value()
				m.currentState = stateInitModel
				m.inputs[11].Blur()
				m.inputs[1].Focus()
				return m, textinput.Blink

			case stateInitModel:
				m.Model = m.inputs[1].Value()
				if providerID == "amazonbedrock" || providerID == "bedrock" || providerID == "amazon" || providerID == "aws" {
					m.currentState = stateInitCommitGuides
					m.inputs[1].Blur()
					m.inputs[3].Focus()
				} else {
					m.currentState = stateInitAPIKey
					m.inputs[1].Blur()
					m.inputs[2].Focus()
				}
				return m, textinput.Blink

			case stateInitAPIKey:
				m.APIKey = m.inputs[2].Value()
				m.currentState = stateInitCommitGuides
				m.inputs[2].Blur()
				m.inputs[3].Focus()
				return m, textinput.Blink

			case stateInitCommitGuides:
				m.CommitGuides = m.inputs[3].Value()
				m.currentState = stateInitChangelogGuides
				m.inputs[3].Blur()
				m.inputs[5].Focus()
				return m, textinput.Blink

			case stateInitChangelogGuides:
				m.ChangelogGuides = m.inputs[5].Value()
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
	case stateInitCommitGuides:
		m.inputs[3], cmd = m.inputs[3].Update(msg)
	case stateInitEndpoint:
		m.inputs[4], cmd = m.inputs[4].Update(msg)
	case stateInitChangelogGuides:
		m.inputs[5], cmd = m.inputs[5].Update(msg)
	case stateInitGCPProjectID:
		m.inputs[6], cmd = m.inputs[6].Update(msg)
	case stateInitGCPRegion:
		m.inputs[7], cmd = m.inputs[7].Update(msg)
	case stateInitAWSRegion:
		m.inputs[8], cmd = m.inputs[8].Update(msg)
	case stateInitAWSAccessKeyID:
		m.inputs[9], cmd = m.inputs[9].Update(msg)
	case stateInitAWSSecretAccessKey:
		m.inputs[10], cmd = m.inputs[10].Update(msg)
	case stateInitAzureResourceName:
		m.inputs[11], cmd = m.inputs[11].Update(msg)
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m InitTUIModel) View() string {
	if m.quitting {
		return ""
	}
	var s strings.Builder

	s.WriteString(titleStyleInit.Render("Welcome to GCT! Let's set up your project.") + "\n")
	s.WriteString(promptStyleInit.Render("Fill in the details below. Press Enter to confirm each step.") + "\n\n")

	if m.currentState > stateInitName {
		s.WriteString(fmt.Sprintf("%s Project Name: %s\n", checkmarkStyleInit.Render("✓"), m.Name))
	} else {
		s.WriteString("Project Name:\n" + m.inputs[0].View() + "\n")
	}

	if m.currentState > stateInitProvider {
		s.WriteString(fmt.Sprintf("%s AI Provider: %s\n", checkmarkStyleInit.Render("✓"), m.Provider))
	} else if m.currentState == stateInitProvider {
		s.WriteString("Select AI Provider (Use ↑/↓):\n")
		for i, choice := range ai.SupportedProviders {
			cursor := "  "
			if m.providerCursor == i {
				cursor = selectedStyleInit.Render("> ")
			}
			s.WriteString(fmt.Sprintf("%s%s\n", cursor, choice))
		}
	}

	if m.currentState > stateInitEndpoint {
		if m.Endpoint != "" {
			s.WriteString(fmt.Sprintf("%s Endpoint URL: %s\n", checkmarkStyleInit.Render("✓"), m.Endpoint))
		}
	} else if m.currentState == stateInitEndpoint {
		s.WriteString("\nEndpoint Base URL:\n" + m.inputs[4].View() + "\n")
	}

	if m.currentState > stateInitGCPProjectID {
		if m.GCPProjectID != "" {
			s.WriteString(fmt.Sprintf("%s GCP Project ID: %s\n", checkmarkStyleInit.Render("✓"), m.GCPProjectID))
		}
	} else if m.currentState == stateInitGCPProjectID {
		s.WriteString("\nGCP Project ID:\n" + m.inputs[6].View() + "\n")
	}
	if m.currentState > stateInitGCPRegion {
		if m.GCPRegion != "" {
			s.WriteString(fmt.Sprintf("%s GCP Region: %s\n", checkmarkStyleInit.Render("✓"), m.GCPRegion))
		}
	} else if m.currentState == stateInitGCPRegion {
		s.WriteString("\nGCP Region:\n" + m.inputs[7].View() + "\n")
	}

	if m.currentState > stateInitAWSRegion {
		if m.AWSRegion != "" {
			s.WriteString(fmt.Sprintf("%s AWS Region: %s\n", checkmarkStyleInit.Render("✓"), m.AWSRegion))
		}
	} else if m.currentState == stateInitAWSRegion {
		s.WriteString("\nAWS Region:\n" + m.inputs[8].View() + "\n")
	}
	if m.currentState > stateInitAWSAccessKeyID {
		if m.AWSAccessKeyID != "" {
			s.WriteString(fmt.Sprintf("%s AWS Access Key ID: %s\n", checkmarkStyleInit.Render("✓"), "[hidden]"))
		}
	} else if m.currentState == stateInitAWSAccessKeyID {
		s.WriteString("\nAWS Access Key ID:\n" + m.inputs[9].View() + "\n")
	}
	if m.currentState > stateInitAWSSecretAccessKey {
		if m.AWSSecretAccessKey != "" {
			s.WriteString(fmt.Sprintf("%s AWS Secret Access Key: %s\n", checkmarkStyleInit.Render("✓"), "[hidden]"))
		}
	} else if m.currentState == stateInitAWSSecretAccessKey {
		s.WriteString("\nAWS Secret Access Key:\n" + m.inputs[10].View() + "\n")
	}

	if m.currentState > stateInitAzureResourceName {
		if m.AzureResourceName != "" {
			s.WriteString(fmt.Sprintf("%s Azure Resource Name: %s\n", checkmarkStyleInit.Render("✓"), m.AzureResourceName))
		}
	} else if m.currentState == stateInitAzureResourceName {
		s.WriteString("\nAzure Resource Name:\n" + m.inputs[11].View() + "\n")
	}

	if m.currentState > stateInitModel {
		s.WriteString(fmt.Sprintf("%s Model/Deployment: %s\n", checkmarkStyleInit.Render("✓"), m.Model))
	} else if m.currentState == stateInitModel {
		modelPrompt := "\nAI Model Name:\n"
		if m.Provider == "Azure OpenAI" {
			modelPrompt = "\nAzure Deployment Name (this is your 'model'):\n"
		}
		s.WriteString(modelPrompt + m.inputs[1].View() + "\n")
	}

	if m.currentState > stateInitAPIKey {
		if m.APIKey != "" {
			s.WriteString(fmt.Sprintf("%s API Key: %s\n", checkmarkStyleInit.Render("✓"), "[hidden]"))
		}
	} else if m.currentState == stateInitAPIKey {
		s.WriteString("\nAPI Key:\n" + m.inputs[2].View() + "\n")
	}

	if m.currentState > stateInitCommitGuides {
		guidesDisplay := m.CommitGuides
		if guidesDisplay == "" {
			guidesDisplay = "None"
		}
		s.WriteString(fmt.Sprintf("%s Commit Guides: %s\n", checkmarkStyleInit.Render("✓"), guidesDisplay))
	} else if m.currentState == stateInitCommitGuides {
		s.WriteString("\nCommit Guides (optional, separate paths with a space):\n" + m.inputs[3].View() + "\n")
	}

	if m.currentState > stateInitChangelogGuides {
		guidesDisplay := m.ChangelogGuides
		if guidesDisplay == "" {
			guidesDisplay = "None"
		}
		s.WriteString(fmt.Sprintf("%s Changelog Guides: %s\n", checkmarkStyleInit.Render("✓"), guidesDisplay))
	} else if m.currentState == stateInitChangelogGuides {
		s.WriteString("\nChangelog Guides (optional, separate paths with a space):\n" + m.inputs[5].View() + "\n")
	}

	s.WriteString(promptStyleInit.Render("\nPress Esc or Ctrl+C to quit at any time."))
	return s.String()
}
