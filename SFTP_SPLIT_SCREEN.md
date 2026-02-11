# Split-Screen SFTP File Manager

## Overview

Termius from Walmart now includes a **Termius-like split-screen SFTP file manager**. When you open SFTP for a server, you get a dual-pane interface with your local filesystem on the left and the remote VPS filesystem on the right.

## Features

### Split-Screen Interface
- **Left Pane**: Local device filesystem
- **Right Pane**: Remote VPS filesystem
- **Real-time Display**: Both panes show current directory path and file listings
- **Active Pane Indicator**: Shows which pane has focus

### File Navigation
- **Browse Directories**: Press [enter] to navigate into directories
- **Parent Directory**: `../` link to go up one level
- **Tab to Switch**: Use [Tab] to switch between local and remote panes
- **Arrow Keys**: Navigate files within each pane

### File Operations

#### Copy/Paste Files
- **[c]** Copy selected file
- Works in both directions:
  - **Local → Remote**: Copy file from your computer to VPS
  - **Remote → Local**: Copy file from VPS to your computer
- **Progress Indicator**: Shows transfer progress and filename

#### Delete Files
- **[d]** Delete selected file
- Works on whichever pane has focus
- Files are immediately removed
- Status message confirms deletion

### Progress & Feedback
- **Real-time Status Messages**: Shows what operation is in progress
- **Progress Bar**: Displays percentage during transfers
- **Error Messages**: Clear error reporting if something goes wrong
- **Transfer Feedback**: Confirms successful operations

## Usage Guide

### Getting Started

1. **Select a Server** from the main list
2. **Press [s]** to open SFTP for that server
3. **Wait** for connection to establish
4. **Use [Tab]** to switch between panes

### Navigation Example
```
1. Select a file in Local pane (left)
2. Press [Tab] to switch to Remote pane
3. Press [enter] on a directory to navigate into it
4. Press [Tab] back to Local pane
5. Press [c] to copy file to current remote directory
```

### Step-by-Step: Uploading a File

1. Browse to file in **Local pane** (press arrows to select)
2. Press **[Tab]** to focus **Remote pane**
3. Navigate to destination directory (press [enter] on folder names)
4. Press **[Tab]** back to **Local pane**
5. Select your file
6. Press **[c]** to copy
7. Watch the progress indicator
8. File appears in **Remote pane** once complete

### Step-by-Step: Downloading a File

1. Press **[Tab]** to focus **Remote pane**
2. Navigate to file location
3. Press **[Tab]** back to **Local pane**
4. Navigate to destination folder
5. Press **[Tab]** to **Remote pane**
6. Select file to download
7. Press **[c]** to copy
8. Watch progress indicator
9. File appears in **Local pane** once complete

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| **[Tab]** | Switch between Local (left) and Remote (right) panes |
| **[↑/↓]** | Navigate files in active pane |
| **[enter]** | Enter directory or perform action |
| **[c]** | Copy selected file to other side |
| **[d]** | Delete selected file |
| **[q]** / **[esc]** | Close SFTP and return to server list |

## Display Elements

### Header
```
SFTP: username@hostname
[Progress indicator if transferring]
```

### Pane Headers
- Active pane shows: `> PANE_NAME <`
- Inactive pane shows: `PANE_NAME`

### File List
- Directories shown with trailing `/`
- Files shown with normal names
- `../` for parent directory

### Status Line
- Shows operation status
- Displays error messages if needed
- Updates in real-time

## Requirements

### Server Configuration
- SSH/SFTP server running on remote (usually port 22)
- Valid credentials (password or PEM key)
- Read/write permissions on directories

### Local System
- Read/write access to local filesystem
- Sufficient disk space for transfers

## Performance Tips

1. **Large Files**: Use copy for large files - it's designed for streaming
2. **Many Files**: Navigate to specific directory first
3. **Network Issues**: Connection remains open during SFTP session
4. **Permissions**: Check file permissions if copy fails

## Troubleshooting

### Connection Failed
- Verify server credentials
- Check if SFTP is enabled on server (usually enabled with SSH)
- Ensure firewall allows SFTP on port (usually 22)

### Permission Denied
- Check file permissions on remote
- Some directories may require elevated privileges
- Try navigating to different directory

### File Not Showing
- Directory may be empty
- Files may be hidden (dotfiles starting with `.`)
- Refresh by navigating away and back

### Transfer Fails
- Check available disk space
- Verify file permissions
- Ensure connection is still active

## Technical Details

### Authentication
- Reuses SSH credentials from server configuration
- Supports password and PEM key authentication
- Secure connection with SSH tunneling

### File Transfer
- Streaming transfer (efficient memory usage)
- Real-time progress updates
- Automatic reconnection on failure

### Local File Access
- Full read/write to user's home directory
- Can navigate entire filesystem
- Respects OS file permissions

## Limitations & Future Improvements

### Current Limitations
- Cannot create new directories (can add via [mkdir](/) command later)
- Single file operations (batch coming soon)
- No file preview/view mode

### Planned Features
- Batch copy/delete operations
- Directory creation
- File permissions editing
- Drag-and-drop support
- Search functionality
- Rename files
- Archive/compress files
