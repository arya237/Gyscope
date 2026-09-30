package linux

import (
	"strings"
	"testing"
	"time"
)

func TestNormalizePCIAddress(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "full PCI address",
			input:    "0000:00:02.0",
			expected: "0:0:02.0",
		},
		{
			name:     "short PCI address",
			input:    "00:02.0",
			expected: "0:0:02.0",
		},
		{
			name:     "another full address",
			input:    "0000:03:00.0",
			expected: "0:3:00.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizePCIAddress(tt.input)

			if got != tt.expected {
				t.Fatalf(
					"normalizePCIAddress(%q) = %q, expected %q",
					tt.input,
					got,
					tt.expected,
				)
			}
		})
	}
}

func TestParseIntelFDInfo(t *testing.T) {
	data := `
drm-driver: i915
drm-client-id: 42
drm-pdev: 0000:00:02.0
drm-engine-render: 120000000
drm-engine-copy: 30000000
drm-engine-video: 10000000
drm-engine-capacity-video: 2
`

	info, err := parseIntelFDInfo(
		strings.NewReader(data),
		"00:02.0",
	)
	if err != nil {
		t.Fatalf(
			"parse Intel fdinfo: %v",
			err,
		)
	}

	if info.driver != "i915" {
		t.Fatalf(
			"expected driver i915, got %q",
			info.driver,
		)
	}

	if info.pdev != "0000:00:02.0" {
		t.Fatalf(
			"expected pdev 0000:00:02.0, got %q",
			info.pdev,
		)
	}

	render := info.engines["render"]

	if render.busyNs != 120000000 {
		t.Errorf(
			"expected render busy %d, got %d",
			uint64(120000000),
			render.busyNs,
		)
	}

	if render.capacity != 1 {
		t.Errorf(
			"expected render capacity 1, got %d",
			render.capacity,
		)
	}

	video := info.engines["video"]

	if video.busyNs != 10000000 {
		t.Errorf(
			"expected video busy %d, got %d",
			uint64(10000000),
			video.busyNs,
		)
	}

	if video.capacity != 2 {
		t.Errorf(
			"expected video capacity 2, got %d",
			video.capacity,
		)
	}
}

func TestParseIntelFDInfo_Xe(t *testing.T) {
	data := `
drm-driver: xe
drm-client-id: 10
drm-pdev: 0000:03:00.0
drm-engine-render: 50000000
`

	info, err := parseIntelFDInfo(
		strings.NewReader(data),
		"0000:03:00.0",
	)
	if err != nil {
		t.Fatalf(
			"parse xe fdinfo: %v",
			err,
		)
	}

	if info.driver != "xe" {
		t.Fatalf(
			"expected driver xe, got %q",
			info.driver,
		)
	}

	if info.engines["render"].busyNs != 50000000 {
		t.Fatalf(
			"unexpected render busy: %d",
			info.engines["render"].busyNs,
		)
	}
}

func TestParseIntelFDInfo_RejectsNonIntelDriver(t *testing.T) {
	data := `
drm-driver: amdgpu
drm-pdev: 0000:03:00.0
drm-engine-render: 100
`

	_, err := parseIntelFDInfo(
		strings.NewReader(data),
		"0000:03:00.0",
	)

	if err == nil {
		t.Fatal("expected error for non-Intel DRM driver")
	}
}

func TestCalculateIntelUsage(t *testing.T) {
	previous := intelUsageSample{
		timestamp: time.Unix(0, 0),
		renderNs:  100,
	}

	current := intelUsageSample{
		timestamp: time.Unix(1, 0),
		renderNs:  600_000_000,
	}

	usage := calculateIntelUsage(
		previous,
		current,
	)

	if usage < 59.9 || usage > 60.1 {
		t.Fatalf(
			"expected usage around 60%%, got %.2f%%",
			usage,
		)
	}
}

func TestCalculateIntelUsage_CounterReset(t *testing.T) {
	previous := intelUsageSample{
		timestamp: time.Unix(0, 0),
		renderNs:  1000,
	}

	current := intelUsageSample{
		timestamp: time.Unix(1, 0),
		renderNs:  100,
	}

	usage := calculateIntelUsage(
		previous,
		current,
	)

	if usage != 0 {
		t.Fatalf(
			"expected 0 usage after counter reset, got %.2f",
			usage,
		)
	}
}

func TestCalculateIntelUsage_NoElapsedTime(t *testing.T) {
	previous := intelUsageSample{
		timestamp: time.Unix(1, 0),
		renderNs:  100,
	}

	current := intelUsageSample{
		timestamp: time.Unix(1, 0),
		renderNs:  200,
	}

	usage := calculateIntelUsage(
		previous,
		current,
	)

	if usage != 0 {
		t.Fatalf(
			"expected 0 usage when elapsed time is zero, got %.2f",
			usage,
		)
	}
}

func TestCalculateIntelUsage_CappedAt100(t *testing.T) {
	previous := intelUsageSample{
		timestamp: time.Unix(0, 0),
		renderNs:  0,
	}

	current := intelUsageSample{
		timestamp: time.Unix(1, 0),
		renderNs:  2_000_000_000,
	}

	usage := calculateIntelUsage(
		previous,
		current,
	)

	if usage != 100 {
		t.Fatalf(
			"expected usage capped at 100%%, got %.2f%%",
			usage,
		)
	}
}
