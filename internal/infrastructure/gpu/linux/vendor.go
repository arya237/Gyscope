package linux

const (
	nvidiaVendorID = "0x10de"
	amdVendorID    = "0x1002"
	intelVendorID  = "0x8086"
)

func gpuVendorName(vendorID string) string {
	switch vendorID {
	case nvidiaVendorID:
		return "NVIDIA"

	case amdVendorID:
		return "AMD"

	case intelVendorID:
		return "Intel"

	default:
		return "Unknown"
	}
}