package linux

import (
	"fmt"
	"strconv"
	"strings"

	applicationNetwork "Gyscope/internal/application/network"
)

const (
	receiveBytesField    = 0
	transmitBytesField   = 8
	networkStatsFieldNum = 16
)

func parseProcNetDev(data []byte) (applicationNetwork.RawNetworkState, error) {
	interfaces := make([]applicationNetwork.RawNetworkInterface, 0)
	seen := make(map[string]struct{})

	for lineIndex, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		name, counters, found := strings.Cut(line, ":")
		if !found {
			continue
		}

		name = strings.TrimSpace(name)
		if name == "" {
			return applicationNetwork.RawNetworkState{}, fmt.Errorf(
				"line %d: missing interface name",
				lineIndex+1,
			)
		}

		if name == "lo" {
			continue
		}

		fields := strings.Fields(counters)
		if len(fields) < networkStatsFieldNum {
			return applicationNetwork.RawNetworkState{}, fmt.Errorf(
				"line %d (%s): expected at least %d counters, got %d",
				lineIndex+1,
				name,
				networkStatsFieldNum,
				len(fields),
			)
		}

		if _, exists := seen[name]; exists {
			return applicationNetwork.RawNetworkState{}, fmt.Errorf(
				"line %d: duplicate interface %q",
				lineIndex+1,
				name,
			)
		}

		receivedBytes, err := strconv.ParseUint(fields[receiveBytesField], 10, 64)
		if err != nil {
			return applicationNetwork.RawNetworkState{}, fmt.Errorf(
				"line %d (%s): parse received bytes %q: %w",
				lineIndex+1,
				name,
				fields[receiveBytesField],
				err,
			)
		}

		transmittedBytes, err := strconv.ParseUint(fields[transmitBytesField], 10, 64)
		if err != nil {
			return applicationNetwork.RawNetworkState{}, fmt.Errorf(
				"line %d (%s): parse transmitted bytes %q: %w",
				lineIndex+1,
				name,
				fields[transmitBytesField],
				err,
			)
		}

		seen[name] = struct{}{}
		interfaces = append(interfaces, applicationNetwork.RawNetworkInterface{
			Name:             name,
			ReceivedBytes:    receivedBytes,
			TransmittedBytes: transmittedBytes,
		})
	}

	return applicationNetwork.RawNetworkState{
		Interfaces: interfaces,
	}, nil
}

