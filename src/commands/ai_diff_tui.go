package commands

import (
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
)

type DiffViewerModel struct {
	viewport viewport.Model
	content  string
}

func NewDiffViewerModel(content string) DiffViewerModel {
	const width = 100

	vp := viewport.New(width, 20)

	renderedContent, err := glamour.Render(content, "dark")
	if err != nil {
		renderedContent = content
	}

	vp.SetContent(renderedContent)

	return DiffViewerModel{
		viewport: vp,
		content:  content,
	}
}

func (m DiffViewerModel) Init() tea.Cmd {
	return nil
}

func (m DiffViewerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - lipgloss.Height(m.headerView()) - lipgloss.Height(m.footerView())
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m DiffViewerModel) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.headerView(),
		m.viewport.View(),
		m.footerView(),
	)
}

func (m DiffViewerModel) headerView() string {
	title := titleStyleViewer.Render("🤖 AI Explanation of Changes")
	return title
}

func (m DiffViewerModel) footerView() string {
	help := helpStyleViewer.Render("Scroll: ↑/↓/pgup/pgdn • Quit: q")
	return help
}
