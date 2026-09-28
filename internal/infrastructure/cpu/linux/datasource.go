package linux

import (
	applicationcpu "Gyscope/internal/application/cpu"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type LinuxCPUDataSource struct{}

const (
	procStatPath    = "/proc/stat"
	procLoadAvgPath = "/proc/loadavg"
	thermalPath = "/sys/class/thermal/thermal_zone*"
)

func NewLinuxCPUDataSource() *LinuxCPUDataSource{
	return &LinuxCPUDataSource{}
}

func (l *LinuxCPUDataSource) Read() (applicationcpu.RawCPUState, error){

	statData, err := os.ReadFile(procStatPath)
	if err != nil {
		return applicationcpu.RawCPUState{}, fmt.Errorf("read %s: %w", procStatPath, err)
	}

	times, err := parseCPUTimes(statData)
	if err != nil {
		return applicationcpu.RawCPUState{}, fmt.Errorf("parse cpu times: %w", err)
	}

	loadAvgData, err := os.ReadFile(procLoadAvgPath)
	if err != nil {
		return applicationcpu.RawCPUState{}, fmt.Errorf("read %s: %w", procLoadAvgPath, err)
	}

	loadAverage, err := parseLoadAverage(loadAvgData)
	if err != nil{
		return applicationcpu.RawCPUState{}, fmt.Errorf("parse load average: %w", err)
	}

	cpuTempreture, err := readCPUTemperature()
	if err != nil{return applicationcpu.RawCPUState{}, fmt.Errorf("read cpu temperature: %w", err)}

	return applicationcpu.RawCPUState{
		Times:        times,
		LogicalCores: runtime.NumCPU(),
		LoadAverage:  loadAverage,
		CPUTemp:      cpuTempreture,
	}, nil
}

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
				return applicationcpu.CPUTimes{}, fmt.Errorf("invalid load average %q: %w", fields[i+1], err)
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
		return applicationcpu.CPUTimes{}, fmt.Errorf("reading load average: %w", err)
	}

	return applicationcpu.CPUTimes{}, fmt.Errorf("cpu stats not found")
}


func parseLoadAverage(data []byte)(applicationcpu.RawLoadAverage, error){
	scanner := bufio.NewScanner(strings.NewReader(string(data)))

	for scanner.Scan(){
		fields := strings.Fields(scanner.Text())

		if len(fields) == 0 || len(fields) < 3{
			return applicationcpu.RawLoadAverage{}, fmt.Errorf("invalid cpu stats: not enough fields")
		}

		values := make([]float64, 3)

		for i := range values{
			value, err := strconv.ParseFloat(fields[i], 64)

			if err != nil{
				return applicationcpu.RawLoadAverage{}, fmt.Errorf("invalid cpu stats value %q: %w", fields[i], err)
			}

			values[i] = value
		}

		return applicationcpu.RawLoadAverage{
			OneMinute: values[0],
			FiveMinutes: values[1],
			FifteenMinutes: values[2],
		}, nil
	}

	if err := scanner.Err(); err != nil {
		return applicationcpu.RawLoadAverage{}, fmt.Errorf("reading cpu stats: %w", err)
	}

	return applicationcpu.RawLoadAverage{}, fmt.Errorf("cpu stats not found")
}

func readCPUTemperature() (float64, error) {
	zones, err := filepath.Glob(thermalPath)
	if err != nil {
		return 0, err
	}

	for _, zone := range zones {
		typeContent, err := os.ReadFile(filepath.Join(zone, "type"))
		if err != nil {
			continue
		}

		zoneType := strings.TrimSpace(string(typeContent))

		if !isCPUThermalZone(zoneType) {
			continue
		}

		tempContent, err := os.ReadFile(filepath.Join(zone, "temp"))
		if err != nil {
			return 0, err
		}

		temp, err := strconv.ParseInt(
			strings.TrimSpace(string(tempContent)),
			10,
			64,
		)
		if err != nil {
			return 0, fmt.Errorf("parse CPU temperature: %w", err)
		}

		return float64(temp) / 1000, nil
	}

	return 0, fmt.Errorf("CPU thermal zone not found")
}

func isCPUThermalZone(zoneType string) bool {
	switch zoneType {
	case "x86_pkg_temp", "cpu-thermal", "cpu_thermal":
		return true
	default:
		return false
	}
}
