package commands

import (
	"time"

	"github.com/atotto/clipboard" //
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyleViewer = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("205")).
				Padding(0, 1)

	helpStyleViewer = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Padding(1, 1)

	copiedStyleViewer = lipgloss.NewStyle().
				Foreground(lipgloss.Color("46")).
				Padding(1, 1).
				Bold(true)
)

type copiedMessage struct{}

type AITextViewerModel struct {
	viewport      viewport.Model
	rawContent    string
	title         string
	showingCopied bool
}

func NewAITextViewerModel(title, content string) AITextViewerModel {
	const width = 100
	vp := viewport.New(width, 20)

	renderedContent, err := glamour.Render(content, "dark")
	if err != nil {
		renderedContent = content
	}

	vp.SetContent(renderedContent)

	return AITextViewerModel{
		viewport:   vp,
		rawContent: content,
		title:      title,
	}
}

func (m AITextViewerModel) Init() tea.Cmd {
	return nil
}

func (m AITextViewerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case copiedMessage:
		m.showingCopied = false
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "c":
			if !m.showingCopied {
				m.showingCopied = true
				clipboard.WriteAll(m.rawContent)
				cmds = append(cmds, tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
					return copiedMessage{}
				}))
			}
		}

	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - lipgloss.Height(m.headerView()) - lipgloss.Height(m.footerView())
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m AITextViewerModel) View() string {
	if m.showingCopied {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			m.headerView(),
			m.viewport.View(),
			m.copiedView(),
		)
	}
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.headerView(),
		m.viewport.View(),
		m.footerView(),
	)
}

func (m AITextViewerModel) headerView() string {
	return titleStyleViewer.Render(m.title)
}

func (m AITextViewerModel) footerView() string {
	return helpStyleViewer.Render("Scroll: ↑/↓ • Copy: c • Quit: q")
}

func (m AITextViewerModel) copiedView() string {
	return copiedStyleViewer.Render("✓ Copied to clipboard!")
}
