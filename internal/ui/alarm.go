package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/umtdemr/spor-istanbul-cli/internal/alarm"
	"github.com/umtdemr/spor-istanbul-cli/internal/service"
	"github.com/umtdemr/spor-istanbul-cli/internal/session"
	"strings"
	"time"
)

type AlarmModel struct {
	api                    *service.Service
	selectedSession        *session.SelectedSession
	selectedSubscriptionId string
	checkCount             int
	sub                    chan bool
	found                  bool // if the spot is found
	spinner                spinner.Model
	err                    error
}

type responseMsg bool

func initialAlarmModel(api *service.Service) AlarmModel {
	return AlarmModel{
		api:     api,
		sub:     make(chan bool),
		spinner: spinner.New(),
	}
}

func (m AlarmModel) Init() tea.Cmd {
	return nil
}

func (m AlarmModel) listenForActivity() tea.Cmd {
	return func() tea.Msg {
		for {
			val := m.api.CheckSessionApplicable(m.selectedSubscriptionId, m.selectedSession.Id)
			if val {
				m.sub <- val
				return nil
			}
			m.sub <- val
			time.Sleep(5 * time.Second)
		}
	}
}

func (m AlarmModel) waitForActivity() tea.Cmd {
	return func() tea.Msg {
		return responseMsg(<-m.sub)
	}
}

func (m AlarmModel) alarmCmd() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.listenForActivity(),
		m.waitForActivity(),
	)
}

func (m AlarmModel) InitAlarm() tea.Cmd {
	return m.alarmCmd()
}

func (m AlarmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case responseMsg:
		m.checkCount++
		if msg {
			close(m.sub)
			m.found = true
			go alarm.PlayAlarm()
			return m, nil
		}

		return m, m.waitForActivity()
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m AlarmModel) View() string {
	doc := strings.Builder{}

	if m.found {
		successBox := lipgloss.
			NewStyle().
			Width(20).
			Border(lipgloss.RoundedBorder()).
			Background(lipgloss.Color("#53A653")).
			Foreground(lipgloss.Color("#FFF")).
			AlignHorizontal(lipgloss.Center).
			Padding(2)
		doc.WriteString(successBox.Render("✓ A spot found!!"))
		return doc.String()
	}

	doc.WriteString(fmt.Sprintf("%s Checking an empty spot", m.spinner.View()))
	doc.WriteString("\n")
	doc.WriteString(fmt.Sprintf("Session Date: %s %s", m.selectedSession.Date, m.selectedSession.Day))
	doc.WriteString("\n")
	doc.WriteString(fmt.Sprintf("Session time: %s", m.selectedSession.Time))
	doc.WriteString("\n")
	doc.WriteString("\n")
	doc.WriteString(fmt.Sprintf("%v times checked so far", m.checkCount))
	doc.WriteString("\n")
	doc.WriteString("\n")

	return doc.String()
}
