package main

import (
	"fmt"
)

func (m model) WaitView() string {
	// The header
	s := fmt.Sprintf("Selected Date: %s,\nSelected session: %s\n", m.selectedDate.Format("02/01/2006"), m.selectedSession)

	s += "\nwaiting... Please do not close this app."
	
	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}
