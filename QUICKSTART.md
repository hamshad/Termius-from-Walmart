# Quick Start Guide - SSH Manager

## Installation in 3 Steps

### 1. Build the Application
```bash
# Make build script executable
chmod +x build.sh

# Run build script
./build.sh

# OR use Make
make build
```

### 2. (Optional) Install Globally
```bash
# Install to /usr/local/bin
sudo mv termius-from-walmart /usr/local/bin/

# OR use Make
make install
```

### 3. Run
```bash
# If installed globally
termius-from-walmart

# OR run locally
./termius-from-walmart
```

## First Time Usage

### Adding Your First Server

#### With Password:
1. Launch the application
2. Press `a` to add a new server
3. Fill in the details:
   - **Name**: My Server (or any friendly name)
   - **Host**: 192.168.1.100 (IP or domain)
   - **Port**: 22 (default SSH port)
   - **User**: root (or your username)
   - **Pass**: your-password
   - **PEM**: (leave empty)
4. Press `Tab` until you reach `[Submit]`
5. Press `Enter`

#### With PEM Key (for AWS, DigitalOcean, etc.):
1. Launch the application
2. Press `a` to add a new server
3. Fill in Name, Host, Port, User
4. Leave Password empty
5. Tab to PEM field, press `p`
6. Paste your entire PEM key (e.g., from `cat my-key.pem`)
7. Press `Ctrl+S` to save the key
8. Navigate to Submit and press Enter

### Connecting to a Server

1. Use arrow keys or `j/k` to select a server
2. Press `Enter` to connect

### Quick Import Example

Want to import multiple servers at once?

1. Create a file at `~/ssh-servers-import.json`:
```bash
cat > ~/ssh-servers-import.json << 'EOF'
[
  {
    "name": "Web Server",
    "host": "192.168.1.100",
    "port": 22,
    "username": "root",
    "password": "mypassword",
    "pem_key": ""
  },
  {
    "name": "AWS Server",
    "host": "ec2-54-123-45-67.compute-1.amazonaws.com",
    "port": 22,
    "username": "ec2-user",
    "password": "",
    "pem_key": "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA...\n-----END RSA PRIVATE KEY-----"
  },
  {
    "name": "DB Server",
    "host": "192.168.1.101",
    "port": 22,
    "username": "admin",
    "password": "",
    "pem_key": ""
  }
]
EOF
```

**Note**: For the `pem_key` field, use `\n` for newlines in JSON format.
To convert a PEM file: `awk '{printf "%s\\n", $0}' your-key.pem`

2. In the application:
   - Press `m` (menu)
   - Select "Import Servers"
   - Press `Enter`

3. Done! Your servers are now in the list.

## Common Tasks

### Export All Servers
1. Press `m`
2. Select "Export Servers"
3. File saved to `~/ssh-servers-export.json`

### Edit a Server
1. Select the server
2. Press `e`
3. Modify fields
4. Submit

### Delete a Server
1. Select the server
2. Press `d`

### Search for a Server
1. Press `/`
2. Type to filter
3. Press `Esc` to clear filter

## Keyboard Shortcuts Cheatsheet

| Key | Action |
|-----|--------|
| `a` | Add new server |
| `e` | Edit selected server |
| `d` | Delete selected server |
| `Enter` | Connect to server |
| `m` | Import/Export menu |
| `/` | Search/filter |
| `↑/↓` or `j/k` | Navigate |
| `Tab` | Next field (in forms) |
| `p` | PEM editor (when on PEM field) |
| `Ctrl+S` | Save (in PEM editor) |
| `Esc` | Cancel/Back |
| `q` or `Ctrl+C` | Quit |

## Tips & Tricks

### 1. Using PEM Keys with Cloud Servers
For AWS, DigitalOcean, Google Cloud, etc.:
```bash
# Copy your PEM key text
cat ~/Downloads/my-aws-key.pem

# In termius-from-walmart:
# - Add server
# - Navigate to PEM field
# - Press 'p' for editor
# - Paste entire key
# - Ctrl+S to save
```

### 2. Use SSH Keys Instead of Passwords
More secure and convenient:
```bash
# Generate SSH key
ssh-keygen -t ed25519

# Copy to server
ssh-copy-id username@server-ip

# In termius-from-walmart, leave password field empty
```

### 3. Backup Your Servers
Regularly export your servers:
```bash
# From the application, press 'm' and export
# Then backup the file
cp ~/ssh-servers-export.json ~/Dropbox/ssh-backup.json
```

### 4. Organize with Names
Use descriptive names:

### 5. Quick Jump Servers
For servers behind a jump host, add both:
1. Add jump host as one entry
2. Add final destination as another
3. Connect to jump host first
4. Then use termius-from-walmart on the jump host

## Troubleshooting

### Can't connect even though details are correct?
```bash
# Try manual SSH first to verify
ssh username@host -p port

# Check if port is blocked
nc -zv host port
```

### Password not working?
Make sure `sshpass` is installed:
```bash
# Ubuntu/Debian
sudo apt-get install sshpass

# macOS
brew install hudochenkov/sshpass/sshpass
```

### Import file not found?
The import file must be at exactly:
```bash
~/ssh-servers-import.json
```
Verify with:
```bash
ls -la ~/ssh-servers-import.json
```

## Next Steps


## Support

Having issues? Common solutions:
1. Check Go version: `go version` (need 1.21+)
2. Rebuild: `make clean && make build`
3. Check logs in terminal output
4. Verify JSON syntax for imports

Happy SSH-ing! 🚀
