package linux

import (
	"fmt"
	"strconv"
	"strings"
)

func parseTotalCPUTime(content string) (uint64, error) {
	for _, line := range strings.Split(content, "\n") {
		fields := strings.Fields(line)

		if len(fields) < 8 {
			continue
		}

		if fields[0] != "cpu" {
			continue
		}

		var total uint64

		for _, field := range fields[1:] {
			value, err := strconv.ParseUint(field, 10, 64)
			if err != nil {
				return 0, fmt.Errorf("parse cpu time: %w", err)
			}

			total += value
		}

		return total, nil
	}

	return 0, fmt.Errorf("cpu line not found")
}