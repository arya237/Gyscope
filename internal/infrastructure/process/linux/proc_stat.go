package linux

import (
	"fmt"
	"strconv"
	"strings"

	applicationProcess "Gyscope/internal/application/process"
)

func parseProcessStat(content string) (applicationProcess.RawProcessStat, error) {
	closeParen := strings.LastIndex(content, ")")
	if closeParen == -1 {
		return applicationProcess.RawProcessState{}, fmt.Errorf("invalid process stat")
	}

	prefix := content[:closeParen+1]
	fields := strings.Fields(content[closeParen+1:])

	openParen := strings.Index(prefix, "(")
	if openParen == -1 {
		return applicationProcess.RawProcessState{}, fmt.Errorf("invalid process stat")
	}

	pid, err := strconv.Atoi(strings.TrimSpace(content[:openParen]))
	if err != nil {
		return applicationProcess.RawProcessState{}, fmt.Errorf("parse pid: %w", err)
	}

	name := prefix[openParen+1 : closeParen]

	if len(fields) < 13 {
		return applicationProcess.RawProcessState{}, fmt.Errorf("invalid process stat fields")
	}

	utime, err := strconv.ParseUint(fields[11], 10, 64)
	if err != nil {
		return applicationProcess.RawProcessState{}, fmt.Errorf("parse utime: %w", err)
	}

	stime, err := strconv.ParseUint(fields[12], 10, 64)
	if err != nil {
		return applicationProcess.RawProcessState{}, fmt.Errorf("parse stime: %w", err)
	}

	return applicationProcess.RawProcessState{
		PID:     pid,
		Name:    name,
		CPUTime: utime + stime,
	}, nil
}
