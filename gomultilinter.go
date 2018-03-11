package main

import (
	"fmt"

	"github.com/liut0/gomultilinter/api"
	"golang.org/x/net/context"
)

const (
	name            = "GoCyclo"
	category        = "cyclo"
	msg             = "cyclomatic complexity [%v/%v] of function %v is high"
	defaultMaxCyclo = 12
)

type GoCycloLinterFactory struct{}

type GoCycloLinter struct {
	MaxCyclo int `json:"max_cyclo"`
}

var LinterFactory api.LinterFactory = &GoCycloLinterFactory{}

func (f *GoCycloLinterFactory) NewLinterConfig() api.LinterConfig {
	return &GoCycloLinter{
		MaxCyclo: defaultMaxCyclo,
	}
}

func (l *GoCycloLinter) NewLinter() (api.Linter, error) {
	return l, nil
}

func (*GoCycloLinter) Name() string {
	return name
}

func (l *GoCycloLinter) LintFile(ctx context.Context, file *api.File, reporter api.IssueReporter) error {
	stats := buildStats(file.ASTFile, file.FSet, []stat{})
	for _, stat := range stats {
		if stat.Complexity > l.MaxCyclo {
			reporter.Report(&api.Issue{
				Severity: api.SeverityWarning,
				Message:  fmt.Sprintf(msg, stat.Complexity, l.MaxCyclo, stat.FuncName),
				Category: category,
				Position: stat.Pos,
			})
		}
	}
	return nil
}
