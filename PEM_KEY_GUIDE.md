# PEM Key Authentication Guide

## What is PEM Key Authentication?

PEM (Privacy Enhanced Mail) keys are a secure way to authenticate SSH connections without using passwords. They're commonly used with cloud providers like AWS, DigitalOcean, and Google Cloud.

## Supported PEM Key Formats

The SSH Manager supports various PEM key formats:

### RSA Keys
```
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA...
(base64 encoded key data)
...
-----END RSA PRIVATE KEY-----
```

### OpenSSH Keys
```
-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAA...
(base64 encoded key data)
...
-----END OPENSSH PRIVATE KEY-----
```

### EC Keys
```
-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIIGHjP...
(base64 encoded key data)
...
-----END EC PRIVATE KEY-----
```

### ED25519 Keys (Modern, Recommended)
```
-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAA...
(base64 encoded key data)
...
-----END OPENSSH PRIVATE KEY-----
```

## How to Add PEM Keys

### Method 1: Using the TUI (Recommended for Small Keys)

1. Press `a` to add a new server or `e` to edit existing
2. Fill in Name, Host, Port, Username
3. **Leave Password field empty**
4. Navigate to the PEM field using Tab
5. Press `p` to open the multiline editor
6. Paste your entire PEM key (Ctrl+Shift+V in most terminals)
7. Press `Ctrl+S` to save the PEM key
8. Submit the form

### Method 2: Using Import (Recommended for Multiple Servers)

1. Create `~/ssh-servers-import.json`:
```json
[
  {
    "name": "AWS Server",
    "host": "ec2-54-123-45-67.compute-1.amazonaws.com",
    "port": 22,
    "username": "ec2-user",
    "password": "",
    "pem_key": "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAA...\n-----END RSA PRIVATE KEY-----"
  }
]
```

2. In the app, press `m` → Select "Import Servers"

**Important**: In JSON, newlines must be represented as `\n`

## Converting PEM Files to Text

If you have a `.pem` file, you need to convert it to text:

### On Linux/Mac:
```bash
# Read the entire file
cat your-key.pem

# Or copy to clipboard
cat your-key.pem | pbcopy  # Mac
cat your-key.pem | xclip -selection clipboard  # Linux
```

### For JSON Import:
```bash
# Convert PEM file to JSON-safe format (with \n for newlines)
awk '{printf "%s\\n", $0}' your-key.pem
```

Then paste the output into the `pem_key` field.

## Example: Full Workflow

### Scenario: Add AWS EC2 Instance

1. **Get your PEM key:**
   ```bash
   cat ~/Downloads/my-aws-key.pem
   ```

2. **Copy the output**, which looks like:
   ```
   -----BEGIN RSA PRIVATE KEY-----
   MIIEpAIBAAKCAQEAwmKwN8zqEXAMPLE...
   (many lines of base64 text)
   ...
   -----END RSA PRIVATE KEY-----
   ```

3. **In SSH Manager:**
   - Press `a` to add server
   - Name: "AWS Production Server"
   - Host: "ec2-54-123-45-67.compute-1.amazonaws.com"
   - Port: "22"
   - User: "ec2-user"
   - Pass: (leave empty)
   - Tab to PEM field, press `p`
   - Paste entire key
   - Press `Ctrl+S`
   - Submit form

4. **Connect:**
   - Select the server
   - Press Enter
   - You'll connect using the PEM key!

## Security Best Practices

### ✅ DO:
- Use PEM keys instead of passwords when possible
- Keep your PEM keys secure (chmod 600)
- Use different keys for different environments
- Regularly rotate your keys
- Use ED25519 keys for new setups (most secure)

### ❌ DON'T:
- Share your private keys
- Commit keys to version control
- Use the same key everywhere
- Use both password AND PEM key (app will reject this)

## Troubleshooting

### "Invalid PEM key format" Error
**Cause**: The key doesn't contain "BEGIN" and "PRIVATE KEY" markers

**Solution**: Make sure you copied the entire key including:
```
-----BEGIN ... PRIVATE KEY-----
... (key data) ...
-----END ... PRIVATE KEY-----
```

### "Permission denied (publickey)" when connecting
**Possible causes**:
1. Wrong username (use `ec2-user` for AWS, `ubuntu` for Ubuntu, etc.)
2. Public key not added to server's `~/.ssh/authorized_keys`
3. Wrong private key

**Solution**:
```bash
# Test the key manually
ssh -i /path/to/key.pem username@host

# Check what user you should use
# AWS: ec2-user, ubuntu, admin
# DigitalOcean: root
# Google Cloud: your-username
```

### Key works manually but not in SSH Manager
**Cause**: Key might have passphrase

**Solution**: SSH Manager currently doesn't support passphrase-protected keys. Remove passphrase:
```bash
# Remove passphrase from key
ssh-keygen -p -f your-key.pem
# Press Enter for empty passphrase
```

### Connection hangs
**Solution**: Add `-o StrictHostKeyChecking=no` option (we'll add this in future update)

## Converting Different Key Formats

### From .ppk (PuTTY) to PEM:
```bash
# Install puttygen
sudo apt-get install putty-tools

# Convert
puttygen your-key.ppk -O private-openssh -o your-key.pem
```

### Generate New ED25519 Key (Recommended):
```bash
# Generate new key
ssh-keygen -t ed25519 -C "your_email@example.com" -f ~/.ssh/my-server-key

# Copy public key to server
ssh-copy-id -i ~/.ssh/my-server-key.pub user@host

# Then add private key to SSH Manager
cat ~/.ssh/my-server-key
```

## Import/Export with PEM Keys

### Export Format:
When you export servers with PEM keys, they're saved in JSON:
```json
[
  {
    "name": "Server",
    "host": "1.2.3.4",
    "port": 22,
    "username": "user",
    "password": "",
    "pem_key": "-----BEGIN RSA PRIVATE KEY-----\nMIIE..."
  }
]
```

### Import Format:
Same as above, but remember:
- Use `\n` for newlines in JSON
- Escape quotes if needed: `\"`
- Leave `password` empty when using `pem_key`

## Quick Reference

| Scenario | Password | PEM Key |
|----------|----------|---------|
| Basic server with password | ✅ | ❌ |
| AWS EC2 | ❌ | ✅ |
| DigitalOcean Droplet | Either | Either |
| Google Cloud | ❌ | ✅ |
| Personal VPS | Either | ✅ (Recommended) |

## Advanced: Using the Same Key for Multiple Servers

You can reuse the same PEM key across servers:

1. Generate one key:
```bash
ssh-keygen -t ed25519 -f ~/.ssh/mykey
```

2. Copy to all servers:
```bash
ssh-copy-id -i ~/.ssh/mykey.pub user@server1
ssh-copy-id -i ~/.ssh/mykey.pub user@server2
```

3. In SSH Manager, add the same PEM key to both server entries

## Need Help?

- Test your key manually first: `ssh -i key.pem user@host`
- Check server logs: `sudo tail -f /var/log/auth.log`
- Verify key permissions: `ls -la your-key.pem` (should be 600)
- Check if public key is on server: `cat ~/.ssh/authorized_keys`

---

**Remember**: PEM keys are like passwords - keep them secure! 🔐
