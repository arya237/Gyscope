# Gyscope

Gyscope is a lightweight Linux system monitor written in Go.

It provides real-time information about CPU, memory, disk, GPU, and running processes directly in your terminal through a clean and minimal TUI.

---

## Features

- **CPU monitoring**
  - CPU usage
  - Logical CPU cores
  - Load average (1m / 5m / 15m)
  - CPU temperature

- **Memory monitoring**
  - Total memory
  - Used memory
  - Available memory
  - Memory usage percentage

- **Disk monitoring**
  - Total space
  - Used space
  - Free space
  - Disk usage percentage

- **GPU monitoring**
  - Multiple GPU support
  - GPU name
  - GPU utilization
  - GPU memory usage
  - GPU temperature

- **Process monitoring**
  - Top processes by CPU usage
  - Process memory usage
  - Process CPU usage
  - Scrollable process list

- Lightweight and fast
- Native Linux interfaces
- No background daemon required
- Terminal-based UI
- Written in Go

---

## Requirements

Gyscope currently targets **Linux** systems.

Supported architectures:

- `amd64`
- `arm64`

For NVIDIA GPUs, Gyscope uses **NVIDIA NVML** to collect GPU metrics.

---

# Installation

There are two ways to install Gyscope.

## Option 1 — Install with `install.sh`

The easiest way to install Gyscope is using the provided installation script.

Clone the repository:

```bash
git clone https://github.com/arya237/Gyscope.git
cd Gyscope
```

Run the installer:

```bash
./scripts/install.sh
```

The installer will:

- Detect your Linux operating system
- Detect your CPU architecture
- Download the appropriate Gyscope binary
- Install it to `~/.local/bin/gyscope`
- Configure your `PATH` when necessary

After installation, run:

```bash
gyscope
```

If `gyscope` is not immediately available in your current shell, reload your shell configuration:

```bash
source ~/.bashrc
```

Then run:

```bash
gyscope
```

---

## Option 2 — Download the Binary Manually

You can also install Gyscope without cloning the repository.

Prebuilt binaries are available on the GitHub Releases page:

https://github.com/arya237/Gyscope/releases

Available Linux builds:

- `gyscope-linux-amd64`
- `gyscope-linux-arm64`

Download the binary that matches your system architecture.

### Example: AMD64

```bash
wget https://github.com/arya237/Gyscope/releases/latest/download/gyscope-linux-amd64
```

Make the binary executable:

```bash
chmod +x gyscope-linux-amd64
```

Create the local binary directory if it does not exist:

```bash
mkdir -p ~/.local/bin
```

Move the binary into it:

```bash
mv gyscope-linux-amd64 ~/.local/bin/gyscope
```

Make sure `~/.local/bin` is in your `PATH`:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

Now you can run:

```bash
gyscope
```

### Make the PATH change permanent

If `~/.local/bin` is not already in your `PATH`, add it to your shell configuration:

```bash
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
```

Then reload your shell:

```bash
source ~/.bashrc
```

---

# Usage

Simply run:

```bash
gyscope
```

Gyscope will start monitoring your system and refresh the displayed information in real time.

---

## Controls

| Key | Action |
|-----|--------|
| `q` | Quit   |
| `Ctrl+C` | Quit |

---

# Monitoring

## CPU

Gyscope reads CPU statistics directly from Linux system interfaces and calculates CPU utilization from consecutive samples.

Displayed information includes:

- CPU usage
- Logical CPU cores
- Load average
- CPU temperature

---

## Memory

Memory information is collected from Linux `/proc/meminfo`.

Gyscope displays:

- Total memory
- Used memory
- Available memory
- Memory usage percentage

---

## Disk

Disk statistics are collected using Linux filesystem statistics.

Gyscope displays:

- Total capacity
- Used space
- Free space
- Disk usage percentage

---

## GPU

GPU discovery is performed through Linux DRM/sysfs interfaces.

Vendor-specific metrics are collected through the appropriate backend.

For NVIDIA GPUs, Gyscope uses **NVML** to obtain:

- GPU name
- GPU utilization
- GPU memory usage
- GPU temperature

Multiple GPUs are supported.

---

## Processes

Gyscope reads process information from Linux `/proc`.

The process panel displays the processes currently consuming the most CPU.

Displayed information includes:

- PID
- Process name
- CPU usage
- Memory usage

The process list is scrollable when more processes are available than can fit on screen.

---

# Building From Source

If you want to build Gyscope yourself, make sure Go is installed.

Clone the repository:

```bash
git clone https://github.com/arya237/Gyscope.git
cd Gyscope
```

Build the project:

```bash
go build -o gyscope ./cmd/Gyscope
```

Run the application:

```bash
./gyscope
```

---

# Development

Run the complete test suite:

```bash
go test ./...
```

Build the project:

```bash
go build ./...
```

---

# Architecture

Gyscope follows a layered architecture designed to keep Linux-specific implementations isolated from the application and domain logic.

```text
Linux OS
   │
   ▼
Infrastructure
   │
   ▼
Application
   │
   ▼
Domain
   │
   ▼
Delivery
   │
   ▼
TUI
```

Linux-specific data collection is implemented in the infrastructure layer, while application logic and domain models remain independent from Linux implementation details.

---

# Project Structure

```text
Gyscope/
├── cmd/
│   └── Gyscope/
│       └── main.go
│
├── internal/
│   ├── domain/
│   │   ├── cpu/
│   │   ├── memory/
│   │   ├── disk/
│   │   ├── gpu/
│   │   └── process/
│   │
│   ├── application/
│   │   ├── cpu/
│   │   ├── memory/
│   │   ├── disk/
│   │   ├── gpu/
│   │   └── process/
│   │
│   ├── infrastructure/
│   │   ├── cpu/
│   │   ├── memory/
│   │   ├── disk/
│   │   ├── gpu/
│   │   └── process/
│   │
│   └── delivery/
│       └── tui/
│
├── scripts/
│   └── install.sh
│
└── .github/
    └── workflows/
        └── release.yml
```

---

# Releases

Prebuilt Linux binaries are published with each release.

You can find the latest version on the GitHub Releases page:

https://github.com/arya237/Gyscope/releases

Available architectures:

- Linux AMD64
- Linux ARM64

---

# License

See the `LICENSE` file in the repository for license information.

---

# Contributing

Contributions, bug reports, and suggestions are welcome.

If you find a bug or have an idea for improving Gyscope, feel free to open an issue or submit a pull request.

---

## Gyscope

A simple, lightweight system monitor for Linux, built with Go.
