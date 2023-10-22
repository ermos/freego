package term

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"time"
)

type toggleDomain struct {
	domains   []string
	nbDomains int
	index     int
	status    string
	spinner   spinner.Model
	done      bool
}

var (
	checkMark = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
)

func newToggleDomain(domains []string, isDisable bool) toggleDomain {
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))

	status := "Added"
	if isDisable {
		status = "Removed"
	}

	return toggleDomain{
		domains:   domains,
		nbDomains: len(domains),
		status:    status,
		spinner:   s,
	}
}

func (m toggleDomain) Init() tea.Cmd {
	return tea.Batch(downloadAndInstall(m.domains[m.index]), m.spinner.Tick)
}

func (m toggleDomain) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case installedPkgMsg:
		if m.index >= len(m.domains)-1 {
			// Everything's been installed. We're done!
			m.done = true
			return m, tea.Quit
		}

		t := tea.Batch(
			tea.Printf("%s %s [%s]", checkMark, m.domains[m.index], m.status),
			downloadAndInstall(m.domains[m.index]),
		)

		m.index++
		return m, t
	}
	return m, nil
}

func (m toggleDomain) View() string {
	if m.done {
		return fmt.Sprintf("%s %s [%s]\n", checkMark, m.domains[m.index], m.status) +
			fmt.Sprintf("Done! %s %d domain(s).\n", m.status, m.nbDomains)
	}

	state := "Adding"
	if m.status == "Removed" {
		state = "Removing"
	}

	return fmt.Sprintf("%s %s %s (%d/%d)", m.spinner.View(), state, m.domains[m.index], m.index, m.nbDomains)
}

type installedPkgMsg string

func downloadAndInstall(pkg string) tea.Cmd {
	d := time.Millisecond * time.Duration(500)
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return installedPkgMsg(pkg)
	})
}

func ToggleDomain(domains []string, isDisable bool) error {
	_, err := tea.NewProgram(newToggleDomain(domains, isDisable)).Run()
	return err
}
