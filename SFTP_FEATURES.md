# SFTP Features Added

## Overview
Added comprehensive SFTP (SSH File Transfer Protocol) support to Termius from Walmart. Users can now manage remote files directly from the application.

## New Features

### 1. SFTP Server Configuration
- Added `SFTPPort` field to the Server struct
- SFTP port defaults to SSH port if not specified
- Imported/exported server configurations include SFTP port settings

### 2. SFTP Connection
- Press **[s]** from the server list to open SFTP browser for the selected server
- Automatic SSH+SFTP connection using existing server credentials
- Support for both PEM key and password authentication
- Secure handling of temporary keys

### 3. File Operations

#### Copy Files (Remote to Remote)
- **Key**: [c]
- Copy files from one location to another on the remote server
- Two-step process: Select source file, then select destination directory

#### Delete Files
- **Key**: [d]
- Delete files from the remote server
- Confirmation via status message

#### Rename Files
- **Key**: [r]
- Rename files on the remote server
- Two-step process: Select file to rename, enter new name

#### Navigate Directories
- **Key**: [enter]
- Navigate into subdirectories
- Parent directory (..) link for going up

### 4. UI Components

#### SFTP Browser View
- Shows current remote path
- Lists all files and directories
- Displays operation messages
- Help text for available commands

#### File Picker for Operations
- Modal interface for selecting source/destination files
- Handles multi-step operations (copy source → destination)
- Clear operation mode indicators

## File Structure

### New Files
- **sftp.go**: SFTP client implementation with connection management and file operations

### Modified Files
- **main.go**: Added SFTP views, state management, and key handlers
- **storage.go**: No changes (data model handles SFTP port automatically)
- **go.mod**: Added dependencies:
  - `github.com/pkg/sftp v1.13.6` - SFTP client library
  - `golang.org/x/crypto v0.21.0` - SSH/cryptography support

## Implementation Details

### Authentication
- Uses existing server credentials (PEM key or password)
- Reuses SSH connection logic from main application
- Insecure host key checking (suitable for private networks, can be enhanced)

### File Operations
The SFTP manager provides:
- `ConnectSFTP()` - Establish SFTP connection
- `ListFiles()` - List directory contents
- `CopyFile()` - Copy remote files
- `DeleteFile()` - Delete remote files
- `RenameFile()` - Rename remote files
- `UploadFile()` - Upload from local (extensible)
- `DownloadFile()` - Download to local (extensible)

### UI States
New view states:
- `sftpView` - Main SFTP browser
- `sftpFilePickerView` - File/directory selection for operations

## Usage Example

1. **Add/Edit a server** with SFTP port (defaults to SSH port)
2. **Select server** from list and press **[s]**
3. **Navigate directories** using arrow keys and [enter]
4. **Copy files**: Press [c], select source, navigate to destination directory, confirm
5. **Delete files**: Press [d], confirm
6. **Rename files**: Press [r], enter new name

## Future Enhancements

Possible improvements:
- Upload files from local filesystem
- Download files to local filesystem
- Create directories
- Batch file operations
- File preview/viewing
- Drag-and-drop file management
- Search within remote directories
- File permissions management
