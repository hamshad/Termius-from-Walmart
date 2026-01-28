# SSH Manager - Project Overview

## 🎯 What is This?

A beautiful, terminal-based SSH connection manager built in Go. Think Termius, but for the command line, lightweight, and open source.

## ✨ Key Features

### Core Functionality
- ✅ **Save & Manage Servers**: Store unlimited SSH connection details
- ✅ **PEM Key Authentication**: Store PEM private keys as text (RSA, ED25519, ECDSA, OpenSSH)
- ✅ **One-Click Connect**: Press Enter to connect instantly
- ✅ **Password Storage**: Optional password saving (with security warnings)
- ✅ **Import/Export**: Bulk import/export in JSON format
- ✅ **Search & Filter**: Quickly find servers by name
- ✅ **Beautiful TUI**: Built with Bubble Tea framework

### User Experience
- 🎨 **Intuitive Interface**: Clean, modern terminal UI
- ⌨️ **Keyboard Driven**: All actions via keyboard shortcuts
- 🔍 **Fuzzy Search**: Built-in filtering for large server lists
- 📝 **Easy Forms**: Simple add/edit forms with validation
- 💾 **Auto-Save**: Changes saved immediately

### Technical Highlights
- 🔐 **Secure by Default**: Config files with 0600 permissions
- 🚀 **Fast & Lightweight**: Single binary, minimal dependencies
- 🔧 **Cross-Platform**: Works on Linux, macOS, and more
- 📦 **No Database**: Simple JSON file storage
- 🎯 **Zero Configuration**: Works out of the box

## 📁 Project Structure

```
termius-from-walmart/
├── main.go                          # Main application code
├── go.mod                           # Go dependencies
├── Makefile                         # Build automation
├── build.sh                         # Build script
├── README.md                        # Full documentation
├── QUICKSTART.md                    # Quick start guide
├── example-import-template.json     # Import template
└── .gitignore                       # Git ignore rules
```

## 🏗️ Architecture

### Data Model
```go
type Server struct {
    ID       int    // Auto-incremented
    Name     string // Display name
    Host     string // IP or domain
    Port     int    // SSH port (default 22)
    Username string // SSH username
    Password string // Optional password
    PemKey   string // Optional PEM private key as text
}

type Config struct {
    Servers []Server // All saved servers
    NextID  int      // Next available ID
}
```

### Storage
- Location: `~/.termius-from-walmart/config.json`
- Format: JSON
- Permissions: 0600 (owner read/write only)

### UI Framework
Built with Charm's TUI stack:
- **Bubble Tea**: TUI framework (event-driven architecture)
- **Bubbles**: Pre-built components (list, text input)
- **Lipgloss**: Styling and layout

## 📥 Import/Export Format

### Template Structure
```json
[
  {
    "name": "Server Name",
    "host": "IP or domain",
    "port": 22,
    "username": "username",
    "password": "optional-password",
    "pem_key": "optional-pem-key-text"
  }
]
```

**Note**: Use either `password` OR `pem_key`, not both.
For `pem_key`, include full PEM text with `\n` for newlines in JSON.

### Import Location
Place file at: `~/ssh-servers-import.json`

### Export Location
File saved to: `~/ssh-servers-export.json`

## 🔐 Security Considerations

### Current Implementation
- ⚠️ Passwords stored in **plain text**
- ⚠️ PEM keys stored in **plain text**
- ✅ Config file has restricted permissions (0600)
- ✅ Config directory is user-only (~/.termius-from-walmart)
- ✅ Temporary PEM key files created with 0600 permissions

### Recommendations
1. **Use PEM Keys or SSH Agent**: More secure than passwords
2. **Personal Use Only**: Not recommended for production environments
3. **Backup Carefully**: Export files contain plain text credentials
4. **Consider Alternatives**: For production, use encrypted vaults

### Future Security Enhancements
- Encrypt config file with master password
- Integration with system keychain
- Support for SSH agent
- Hardware key support (YubiKey)
- Passphrase-protected PEM keys

## 🎮 User Interface

### Main List View
```
┌─────────────────────────────────────────┐
│ SSH Connection Manager                  │
├─────────────────────────────────────────┤
│ > Production Web Server                 │
│   root@192.168.1.100:22                 │
│                                          │
│   Database Server                       │
│   admin@db.example.com:22               │
│                                          │
│   Development Environment               │
│   dev@dev.mycompany.com:2222            │
└─────────────────────────────────────────┘
Keys: [a]dd • [e]dit • [d]elete • [enter] connect
```

### Add/Edit Form
```
Add New Server

Name: Production Server█
Host: 192.168.1.100
Port: 22
User: root
Pass: ••••••••

> [Submit] <

Navigate: [tab] • Submit: [enter] • Cancel: [esc]
```

### Import/Export Menu
```
Import/Export Menu

> Import Servers
  Export Servers
  Back to List

Use ↑/↓ to navigate, [enter] to select
```

## 🛠️ Development

### Prerequisites
- Go 1.21 or higher
- Terminal with 256 color support
- (Optional) sshpass for password authentication

### Build
```bash
# Download dependencies
go mod download

# Build binary
go build -o termius-from-walmart main.go

# Or use Make
make build

# Or use build script
./build.sh
```

### Run in Development
```bash
go run main.go
```

### Testing
```bash
# Run tests
make test

# Or directly
go test -v ./...
```

## 📊 Technical Specifications

### Dependencies
| Package | Version | Purpose |
|---------|---------|---------|
| bubbletea | 0.25.0 | TUI framework |
| bubbles | 0.18.0 | UI components |
| lipgloss | 0.9.1 | Styling |

### Performance
- **Startup Time**: < 100ms
- **Memory Usage**: ~5-10MB
- **Binary Size**: ~8-10MB (statically linked)
- **Max Servers**: Limited only by available memory

### Compatibility
- ✅ Linux (all distributions)
- ✅ macOS (Intel & Apple Silicon)
- ✅ BSD variants
- ✅ Windows (via WSL)

## 🎯 Use Cases

### Personal Development
- Manage dev, staging, production servers
- Quick access to multiple projects
- Backup server configurations

### System Administration
- Organize servers by client or project
- Share configurations across team (via export)
- Quick server switching

### Learning & Education
- Practice SSH management
- Learn Go and TUI development
- Study clean code architecture

## 🚀 Getting Started

### 1. Installation
```bash
# Clone or download the project
cd termius-from-walmart

# Build
./build.sh

# Install (optional)
sudo mv termius-from-walmart /usr/local/bin/
```

### 2. Add First Server
```bash
# Run application
termius-from-walmart

# Press 'a' to add server
# Fill in details and submit
```

### 3. Import Existing Servers
```bash
# Create import file
cat > ~/ssh-servers-import.json << 'EOF'
[
  {
    "name": "My Server",
    "host": "192.168.1.100",
    "port": 22,
    "username": "root",
    "password": ""
  }
]
EOF

# In application, press 'm' and import
```

## 🔮 Future Enhancements

### Planned Features
- [ ] SSH key management interface
- [ ] Connection history and logs
- [ ] Port forwarding configuration
- [ ] Jump host support
- [ ] Server groups/folders
- [ ] SFTP file browser
- [ ] Config file encryption
- [ ] Cloud sync (optional)
- [ ] Terminal session recording

### Advanced Features
- [ ] Multi-session support
- [ ] Automated backups
- [ ] Custom SSH options
- [ ] Proxy configuration
- [ ] Integration with cloud providers

## 📝 Code Quality

### Best Practices
- ✅ Clean, readable code
- ✅ Proper error handling
- ✅ Input validation
- ✅ Secure file permissions
- ✅ Modular design

### Code Organization
- Single-file design for simplicity
- Clear separation of concerns
- Reusable functions
- Type-safe operations

## 🤝 Contributing

While this is a single-file project for simplicity, contributions are welcome:
1. Bug fixes
2. Feature requests
3. Documentation improvements
4. Security enhancements

## 📄 License

MIT License - Free to use, modify, and distribute.

## 🙏 Credits

Built with:
- [Charm.sh](https://charm.sh) - Amazing TUI tools
- Go community - Best language community
- Termius - Inspiration for features

## 📞 Support

Need help?
1. Check QUICKSTART.md for common tasks
2. Read README.md for detailed docs
3. Open an issue on GitHub
4. Review example files

---

**Made with ❤️ and Go**

*A simple, secure, and efficient way to manage your SSH connections.*
