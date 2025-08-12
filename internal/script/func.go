package script

import (
	"strings"

	"github.com/byte4cat/nbx/v2/internal/script/command"
	"github.com/byte4cat/nbx/v2/pkg/clog"
	"github.com/fatih/color"
)

type Func interface {
	// Run executes the ScriptFunc
	Run()
	DryRun() string
}

type FuncImpl struct {
	script ScriptFunc
}

// DryRun implements Func.
func (f *FuncImpl) DryRun() string {
	var sb strings.Builder
	if f.script.StartMsg != "" {
		sb.WriteString(color.BlueString("=====\t%v\t=====\n", f.script.StartMsg))
	}

	for _, cmd := range f.script.Commands {
		for k := range f.script.ReplaceString {
			// Replace placeholders with values
			cmd = strings.ReplaceAll(cmd, k, k)
		}
		sb.WriteString(color.HiBlueString("$ %v\n", cmd))
	}

	if f.script.EndMsg != "" {
		sb.WriteString(color.BlueString("=====\t%v\t=====\n", f.script.EndMsg))
	}

	return sb.String()
}

// Run implements Func.
func (f *FuncImpl) Run() {
	if f.script.StartMsg != "" {
		clog.Line(f.script.StartMsg)
	}

	for _, cmd := range f.script.Commands {
		for k, v := range f.script.ReplaceString {
			cmd = strings.ReplaceAll(cmd, k, v)
		}
		command.Run(cmd, false)
	}

	if f.script.EndMsg != "" {
		clog.Line(f.script.EndMsg)
	}
}

func New(sf ScriptFunc) Func {
	return &FuncImpl{
		script: sf,
	}
}

type ScriptFunc struct {
	StartMsg      string
	Commands      []string
	EndMsg        string
	ReplaceString map[string]string
}
