package linux

import (
	applicationcpu "Gyscope/internal/application/cpu"
	"bufio"
	"fmt"
	"strconv"
	"strings"
)



func parseCPUTimes(data []byte) (applicationcpu.CPUTimes, error){

	scanner := bufio.NewScanner(strings.NewReader(string(data)))

	for scanner.Scan(){
		fields := strings.Fields(scanner.Text())

		if len (fields) == 0 || fields[0] != "cpu"{
			continue
		}

		if len(fields) < 8 {
			return applicationcpu.CPUTimes{}, fmt.Errorf("invalid cpu stats: not enough fields")
		}

		values := make([]uint64, 8)
		
		for i := range values{
			
			value, err := strconv.ParseUint(fields[i + 1], 10, 64)
			if err != nil{
				return applicationcpu.CPUTimes{}, fmt.Errorf("invalid cpu stats value %q: %w", fields[i+1], err)
			}

			values[i] = value
		}
		
		return applicationcpu.CPUTimes{
			User: values[0] + values[1], // user + nice
			System : values[2],
			Idle: values[3],
			IOWait: values[4],
			IRQ: values[5],
			SoftIRQ: values[6],
			Steal: values[7],
		}, nil
	}

	
	if err := scanner.Err(); err != nil {
		return applicationcpu.CPUTimes{}, fmt.Errorf("reading cpu stats: %w", err)
	}

	return applicationcpu.CPUTimes{}, fmt.Errorf("cpu stats not found")
}
