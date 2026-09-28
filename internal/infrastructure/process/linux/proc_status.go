package linux

import (
	"fmt"
	"strconv"
	"strings"
)

func parseProcessStatus(content string) (uint64, error) {
	for _, line := range strings.Split(content, "\n") {
		fields := strings.Fields(line)

		if len(fields) < 2 {
			continue
		}

		if fields[0] != "VmRSS:" {
			continue
		}

		value, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("parse VmRSS: %w", err)
		}

		// VmRSS is reported in kB.
		return value * 1024, nil
	}

	return 0, fmt.Errorf("VmRSS not found")
}