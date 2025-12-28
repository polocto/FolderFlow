---
title: Installation
sidebar_position: 2
---

FolderFlow is a command-line application written in Go (Golang).  
It is distributed as source code and can be built into a single standalone binary.

---
# Download Prebuilt Binary (Recommended)

FolderFlow provides prebuilt binaries for common platforms via **GitHub Releases**.

No additional dependencies are required to run the binary.

&rarr; Latest release:  
https://github.com/polocto/FolderFlow/releases/latest

&rarr; All releases:    
https://github.com/polocto/FolderFlow/releases

---

## Download

Download the archive corresponding to your operating system and architecture.

Common artifacts include:

- `folderflow-linux-amd64`
- `folderflow-linux-arm64`
- `folderflow-darwin-amd64`
- `folderflow-darwin-arm64`
- `folderflow-windows-amd64.exe`

Extract the downloaded file.

---

## Make Executable (Linux / macOS)

```bash
chmod +x folderflow
```

Optional [Intall Globaly](./installation.md#optional-install-globally)

---

# Install by Building from Source
## Requirements

- Go **1.25 or newer**
- Windows, macOS, or Linux
- Git (recommended)

## Clone the Repository

```bash
git clone https://github.com/polocto/FolderFlow.git
cd FolderFlow
```
## Build the Binary

From project root
```bash
go build ./cmd/folderflow
```
This will generate an executable named `folderflow` (or `folderflow.exe` on Windows).

# (Optional) Install Globally

Move the binary into a directory included in your system PATH.

## Linux/macOS

```bash
sudo mv folderflow /usr/local/bin/

```

## Windows (PowerShell)
Move `folderflow.exe` to a folder such as `C:\Windows\System32` or any directory already in your `PATH`.

# Verify Installation

Run:
```bash
folderflow --help
```
If installed correctly, the command usage information will be displayed.

# Cross-Compilation

You can build FolderFlow for other platforms using Go environment variables.

Example: build for Linux from macOS or Windows
```bash
GOOS=linux GOARCH=amd64 go build ./cmd/folderflow
```

# Uninstall 
To uninstall FolderFlow, simply remove the binary from your system:
```bash
rm /usr/local/bin/folderflow
```
(or delete `folderflow.exe` on Windows)