# Termius from Walmart (SSH Connection Manager)

A terminal-based SSH connection manager written in Go, similar to Termius. Manage your SSH servers with an intuitive TUI interface.

## Features

- 📝 **Save SSH Connections**: Store server details including host, port, username, and password
- 🔑 **PEM Key Support**: Add PEM private keys as text (RSA, ED25519, ECDSA, OpenSSH formats)
- 🔐 **Secure Storage**: Configurations saved in `~/.termius-from-walmart/config.json` with 0600 permissions
- 📤 **Export Servers**: Export all server configurations to JSON
- 📥 **Import Servers**: Import server configurations from JSON file
- 🎨 **Beautiful TUI**: Built with Bubble Tea for a smooth terminal experience
- 🔍 **Search & Filter**: Quickly find servers by name
- ⌨️ **Keyboard Navigation**: Efficient keyboard shortcuts

## Installation

### Prerequisites

- Go 1.21 or higher
- (Optional) `sshpass` for password-based authentication

Install sshpass (optional):
```bash
# Ubuntu/Debian
sudo apt-get install sshpass

# macOS
brew install hudochenkov/sshpass/sshpass

# Fedora
sudo dnf install sshpass
```

### Build from Source

```bash
# Clone or create the project directory
cd termius-from-walmart

# Download dependencies
go mod download

# Build the application
go build -o termius-from-walmart main.go

# (Optional) Install to your PATH
sudo mv termius-from-walmart /usr/local/bin/
```

## Usage

### Starting the Application

```bash
./termius-from-walmart
# or if installed globally
termius-from-walmart
```

### Keyboard Shortcuts

#### Main List View
- `a` - Add new server
- `e` - Edit selected server
- `d` - Delete selected server
- `Enter` - Connect to selected server
- `m` - Open Import/Export menu
- `/` - Filter/search servers
- `↑/↓` or `j/k` - Navigate list
- `q` or `Ctrl+C` - Quit

#### Add/Edit Form
- `Tab` / `Shift+Tab` - Move between fields
- `↑/↓` - Move between fields
- `p` - Open PEM key multiline editor (when on PEM field)
- `Enter` - Submit (when on submit button)
- `Esc` - Cancel and return to list

#### PEM Key Editor
- `Ctrl+S` - Save PEM key and return to form
- `Esc` - Cancel and return to form without saving
- Type normally to enter key data

#### Import/Export Menu
- `↑/↓` or `j/k` - Navigate options
- `Enter` - Select option
- `Esc` - Back to list

## Import/Export Format

### Export
When you export servers (press `m` then select "Export Servers"), all servers are saved to:
```
~/ssh-servers-export.json
```

### Import
To import servers, create a file at:
```
~/ssh-servers-import.json
```

### JSON Template

```json
[
  {
    "name": "Production Server",
    "host": "192.168.1.100",
    "port": 22,
    "username": "root",
    "password": "your-password",
    "pem_key": ""
  },
  {
    "name": "AWS EC2 Instance",
    "host": "ec2-54-123-45-67.compute-1.amazonaws.com",
    "port": 22,
    "username": "ec2-user",
    "password": "",
    "pem_key": "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA...\n-----END RSA PRIVATE KEY-----"
  },
  {
    "name": "Development Server",
    "host": "dev.example.com",
    "port": 2222,
    "username": "admin",
    "password": "dev-password",
    "pem_key": ""
  },
  {
    "name": "Staging Server",
    "host": "staging.example.com",
    "port": 22,
    "username": "deploy",
    "password": "",
    "pem_key": ""
  }
]
```

**Note**: You can use either `password` OR `pem_key`, not both. Leave one empty.
For `pem_key`, include the entire PEM private key as text. In JSON, use `\n` for newlines.

### Import Steps

1. Create `~/ssh-servers-import.json` with your servers
2. Run the application
3. Press `m` to open menu
4. Select "Import Servers"
5. Your servers will be added to the list

## Configuration Storage

All configurations are stored in:
```
~/.termius-from-walmart/config.json
```

This file contains:
- All saved servers
- Auto-incrementing IDs for each server
- Passwords stored in plain text (use with caution)
- PEM private keys stored as text (use with caution)

**Security Note**: Both passwords and PEM keys are stored unencrypted. For production environments, consider using SSH agent or system keychain. The config file has 0600 permissions (owner read/write only), but the data itself is not encrypted.

## Connecting to Servers

### With PEM Key (recommended for cloud servers)
If you've saved a PEM private key, the application will automatically use it for authentication.

**Supported PEM formats:**
- RSA: `-----BEGIN RSA PRIVATE KEY-----`
- OpenSSH: `-----BEGIN OPENSSH PRIVATE KEY-----`
- EC: `-----BEGIN EC PRIVATE KEY-----`
- ED25519 (recommended)

**How to add PEM key:**
1. When adding/editing a server, navigate to PEM field
2. Press `p` to open the multiline editor
3. Paste your entire PEM key
4. Press `Ctrl+S` to save
5. Leave password field empty

**Get your PEM key text:**
```bash
cat your-key.pem
```

See [PEM_KEY_GUIDE.md](PEM_KEY_GUIDE.md) for detailed instructions.

### With Password (requires sshpass)
If you've saved a password, the application will use `sshpass` to automatically authenticate.

### With SSH Keys (recommended)
Leave both password and PEM fields empty and ensure your SSH keys are properly configured in `~/.ssh/`.

## Example Workflow

1. **Add your first server:**
   - Press `a`
   - Fill in server details
   - Navigate to Submit and press Enter

2. **Connect to a server:**
   - Select server from list
   - Press Enter

3. **Export for backup:**
   - Press `m`
   - Select "Export Servers"
   - File saved to `~/ssh-servers-export.json`

4. **Import on another machine:**
   - Copy `ssh-servers-export.json` to `~/ssh-servers-import.json`
   - Press `m`
   - Select "Import Servers"

## Troubleshooting

### "sshpass: command not found"
Install sshpass or use SSH keys instead of passwords.

### Cannot connect to server
- Verify the host and port are correct
- Check that the server is accessible
- Ensure your credentials are valid
- Try connecting manually with: `ssh username@host -p port`

### Import fails
- Verify the JSON file exists at `~/ssh-servers-import.json`
- Check JSON syntax is valid
- Ensure file has read permissions

## Development

### Project Structure
```
.
├── main.go      # Main application code
├── go.mod       # Go module dependencies
└── README.md    # This file
```

### Dependencies
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Style definitions

## Security Considerations

⚠️ **Important Security Notes:**

1. **Password Storage**: Passwords are stored in plain text in the config file. This is suitable for personal/development use but NOT recommended for production environments.

2. **File Permissions**: The config file is created with 0600 permissions (owner read/write only), but still contains sensitive data.

3. **Best Practice**: Use SSH keys instead of passwords when possible:
   ```bash
   ssh-keygen -t ed25519 -C "your_email@example.com"
   ssh-copy-id username@hostname
   ```

4. **Alternative**: Consider encrypting the config file or using a system keychain for password storage in production.

## License

MIT License - Feel free to use and modify as needed.

## Contributing

Contributions are welcome! Feel free to submit issues or pull requests.

## Future Enhancements

Potential features for future versions:
- SSH key management
- Connection history
- Port forwarding management
- SFTP file browser
- Terminal session recording
- Config encryption
- Group/folder organization
- Jump host support
