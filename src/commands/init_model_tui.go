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
	// We embed the original model. This gives us all its fields and methods.
	InitTUIModel

	// State specific to this "preset selection" view
	presetCursor int
	inPresetView bool // A flag to know if we are in the preset selection screen
}

func NewInitModelTUIModel() InitModelTUIModel {
	// Start with the standard init TUI model
	baseModel := NewInitTUIModel()
	// Override the starting state for our preset flow
	baseModel.currentState = stateInitName // We still start by asking for the project name

	return InitModelTUIModel{
		InitTUIModel: baseModel,
		inPresetView: true, // We start in the preset selection view
	}
}

func (m InitModelTUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// If we are in the preset selection view, handle it here.
	if m.inPresetView {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyCtrlC, tea.KeyEsc:
				m.quitting = true
				return m, tea.Quit

			case tea.KeyEnter:
				// --- THIS IS THE KEY FIX ---
				// 1. A preset is selected.
				selectedPreset := ai.ModelPresets[m.presetCursor]

				// 2. We now populate the fields of the embedded InitTUIModel.
				m.Provider = selectedPreset.Provider
				m.Model = selectedPreset.ModelName // <-- CRITICAL: Set the model field directly.

				// 3. Set the value in the text input as well for display.
				m.inputs[1].SetValue(selectedPreset.ModelName)

				// 4. We are no longer in the preset view.
				m.inPresetView = false

				// 5. Now, we figure out where to go next based on the provider.
				// This reuses the exact same logic from InitTUIModel's Update.
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
					// For simple providers, we go directly to the API key step,
					// because the model name is already set.
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

	// If we are NOT in the preset view, pass all control to the embedded InitTUIModel's Update method.
	// This is the part that was missing.
	updatedModel, cmd := m.InitTUIModel.Update(msg)
	m.InitTUIModel = updatedModel.(InitTUIModel) // Recast the returned model to our base type
	return m, cmd
}

func (m InitModelTUIModel) View() string {
	if m.quitting {
		return ""
	}

	// If we're in the preset selection view, show the list of models.
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

	// Otherwise, just render the standard init TUI view.
	return m.InitTUIModel.View()
}
