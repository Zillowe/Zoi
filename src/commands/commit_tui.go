package commands

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	errMsgTUI error
)

const (
	idxType = iota
	idxSubject
	idxBody
)

var (
	helpStyleTUI = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	labelStyle   = lipgloss.NewStyle().Bold(true)
)

type CommitTUIModel struct {
	typeInput    textinput.Model
	subjectInput textinput.Model
	bodyInput    textarea.Model

	focusIndex int
	width      int
	height     int
	submitted  bool
	quitting   bool
	err        error

	CommitType string
	Subject    string
	Body       string
}

func NewCommitTUIModel(initialType, initialSubject, initialBody string) CommitTUIModel {
	typeTI := textinput.New()
	typeTI.Placeholder = "e.g. ✨ Feat, Fix, Chore"
	typeTI.Focus()
	typeTI.CharLimit = 50
	typeTI.Width = 50
	typeTI.Prompt = "Type: "
	typeTI.SetValue(initialType)

	subjectTI := textinput.New()
	subjectTI.Placeholder = "Concise explanation of the change"
	subjectTI.CharLimit = 100
	subjectTI.Width = 70
	subjectTI.Prompt = "Subject: "
	subjectTI.SetValue(initialSubject)

	bodyTA := textarea.New()
	bodyTA.Placeholder = "Detailed explanation of your changes. Ctrl+D to submit."
	bodyTA.SetWidth(70)
	bodyTA.SetHeight(5)
	bodyTA.SetValue(initialBody)

	return CommitTUIModel{
		typeInput:    typeTI,
		subjectInput: subjectTI,
		bodyInput:    bodyTA,
		focusIndex:   idxType,
		CommitType:   initialType,
		Subject:      initialSubject,
		Body:         initialBody,
	}
}

func (m CommitTUIModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m CommitTUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:

		if msg.Type == tea.KeyCtrlC {
			m.quitting = true
			return m, tea.Quit
		}

		switch m.focusIndex {
		case idxType:
			m, cmd = m.updateTypeInput(msg)
		case idxSubject:
			m, cmd = m.updateSubjectInput(msg)
		case idxBody:
			m, cmd = m.updateBodyInput(msg)
		}
		cmds = append(cmds, cmd)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		inputWidth := min(70, msg.Width-len(m.typeInput.Prompt)-4)
		if inputWidth < 10 {
			inputWidth = 10
		}
		m.typeInput.Width = min(50, inputWidth)
		m.subjectInput.Width = inputWidth
		m.bodyInput.SetWidth(inputWidth)
		m.bodyInput.SetHeight(min(10, msg.Height-12))

	case errMsgTUI:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(cmds...)
}

func (m CommitTUIModel) updateTypeInput(msg tea.KeyMsg) (CommitTUIModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.Type {
	case tea.KeyEnter, tea.KeyTab:
		m.typeInput.Blur()
		m.subjectInput.Focus()
		m.focusIndex = idxSubject
		return m, textinput.Blink
	case tea.KeyShiftTab:
		m.typeInput.Blur()
		m.bodyInput.Focus()
		m.focusIndex = idxBody

		return m, nil
	case tea.KeyEsc:
		m.quitting = true
		return m, tea.Quit
	}
	m.typeInput, cmd = m.typeInput.Update(msg)
	return m, cmd
}

func (m CommitTUIModel) updateSubjectInput(msg tea.KeyMsg) (CommitTUIModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.Type {
	case tea.KeyEnter, tea.KeyTab:
		m.subjectInput.Blur()
		m.bodyInput.Focus()
		m.focusIndex = idxBody
		return m, nil
	case tea.KeyShiftTab:
		m.subjectInput.Blur()
		m.typeInput.Focus()
		m.focusIndex = idxType
		return m, textinput.Blink
	case tea.KeyEsc:
		m.quitting = true
		return m, tea.Quit
	}
	m.subjectInput, cmd = m.subjectInput.Update(msg)
	return m, cmd
}

func (m CommitTUIModel) updateBodyInput(msg tea.KeyMsg) (CommitTUIModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.Type {
	case tea.KeyCtrlD:
		m.bodyInput.Blur()
		m.submitted = true
		m.quitting = true
		m.CommitType = strings.TrimSpace(m.typeInput.Value())
		m.Subject = strings.TrimSpace(m.subjectInput.Value())
		m.Body = m.bodyInput.Value()
		return m, tea.Quit
	case tea.KeyEsc:
		m.quitting = true
		return m, tea.Quit
	case tea.KeyTab:
		m.bodyInput.Blur()
		m.typeInput.Focus()
		m.focusIndex = idxType
		return m, textinput.Blink
	case tea.KeyShiftTab:
		m.bodyInput.Blur()
		m.subjectInput.Focus()
		m.focusIndex = idxSubject
		return m, textinput.Blink

	}
	m.bodyInput, cmd = m.bodyInput.Update(msg)
	return m, cmd
}

func (m CommitTUIModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\nPress Ctrl+C to quit.", m.err)
	}

	if m.quitting && m.submitted {
		return "Processing commit...\n"
	}
	if m.quitting && !m.submitted {
		return "Commit cancelled.\n"
	}

	var s strings.Builder
	s.WriteString(labelStyle.Render("✨ GCT Commit ✨") + "\n\n")

	s.WriteString(m.typeInput.View())
	s.WriteString("\n\n")

	s.WriteString(m.subjectInput.View())
	s.WriteString("\n\n")

	bodyLabelText := "Body (Ctrl+D to submit when done):"
	if m.focusIndex == idxBody {
		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(bodyLabelText) + "\n")
	} else {
		s.WriteString(labelStyle.Render(bodyLabelText) + "\n")
	}
	s.WriteString(m.bodyInput.View())
	s.WriteString("\n\n")

	help := "Tab/Shift+Tab: Navigate | Enter: Next field (Type/Subject) | Ctrl+C: Quit"
	s.WriteString(helpStyleTUI.Render(help))
	s.WriteString("\n")

	return s.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
