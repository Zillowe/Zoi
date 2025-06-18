package commands

import (
	"fmt"
	"gct/src/ai"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	stateSelectModelPreset = iota
)

type InitModelTUIModel struct {
	InitTUIModel

	presetCursor int
	inPresetView bool
}

func NewInitModelTUIModel() InitModelTUIModel {
	baseModel := NewInitTUIModel()
	baseModel.currentState = stateInitName

	return InitModelTUIModel{
		InitTUIModel: baseModel,
		inPresetView: true,
	}
}

func (m InitModelTUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.inPresetView {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyCtrlC, tea.KeyEsc:
				m.quitting = true
				return m, tea.Quit

			case tea.KeyEnter:
				selectedPreset := ai.ModelPresets[m.presetCursor]

				m.Provider = selectedPreset.Provider
				m.Model = selectedPreset.ModelName

				m.inputs[1].SetValue(selectedPreset.ModelName)

				m.inPresetView = false

				providerID := strings.ToLower(strings.ReplaceAll(m.Provider, " ", ""))
				switch providerID {
				case "googlevertexai", "vertexai":
					m.currentState = stateInitGCPProjectID
					m.inputs[6].Focus()
				case "amazonbedrock", "bedrock":
					m.currentState = stateInitAWSRegion
					m.inputs[8].Focus()
				case "azureopenai":
					m.currentState = stateInitAzureResourceName
					m.inputs[11].Focus()
				case "openaicompatible":
					m.currentState = stateInitEndpoint
					m.inputs[4].Focus()
				default:
					m.currentState = stateInitAPIKey
					m.inputs[2].Focus()
				}
				return m, textinput.Blink

			case tea.KeyUp:
				if m.presetCursor > 0 {
					m.presetCursor--
				}
			case tea.KeyDown:
				if m.presetCursor < len(ai.ModelPresets)-1 {
					m.presetCursor++
				}
			}
		}
		return m, nil
	}

	updatedModel, cmd := m.InitTUIModel.Update(msg)
	m.InitTUIModel = updatedModel.(InitTUIModel)
	return m, cmd
}

func (m InitModelTUIModel) View() string {
	if m.quitting {
		return ""
	}

	if m.inPresetView {
		var s strings.Builder
		s.WriteString(titleStyleInit.Render("Welcome to GCT! Let's set up your project.") + "\n")
		s.WriteString(promptStyleInit.Render("Select a popular model preset to get started (Use ↑/↓):") + "\n\n")

		for i, preset := range ai.ModelPresets {
			cursor := "  "
			if m.presetCursor == i {
				cursor = selectedStyleInit.Render("> ")
			}
			s.WriteString(fmt.Sprintf("%s%s\n", cursor, preset.DisplayName))
		}
		s.WriteString(promptStyleInit.Render("\nPress Esc or Ctrl+C to quit at any time."))
		return s.String()
	}

	return m.InitTUIModel.View()
}
