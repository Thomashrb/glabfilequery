package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Program struct {
	*tea.Program
}

func NewProgram() Program {
	return Program{tea.NewProgram(newModel())}
}

func (p Program) StageMsgSend(msg string) {
	p.Send(stageMsg(msg))
}

func (p Program) QuitMsgSend() {
	p.Send(exitMsg(true))
}
