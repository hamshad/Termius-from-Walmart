# Troubleshooting PEM Key Issues

## Common Errors and Solutions

### Error: "Load key invalid format"

This means the PEM key format isn't recognized. Here's how to fix it:

#### **Solution 1: Check Your PEM Key Format**

Your PEM key should look EXACTLY like this (with proper line breaks):

```
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAwQz8...
(many lines of base64 text)
...
-----END RSA PRIVATE KEY-----
```

OR for OpenSSH format:

```
-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAA...
(many lines of base64 text)
...
-----END OPENSSH PRIVATE KEY-----
```

#### **Solution 2: Verify Your Key File**

```bash
# Check your original PEM file
cat your-key.pem

# Test if it works manually
ssh -i your-key.pem ubuntu@16.171.3.225

# If it works manually, the key is valid
```

#### **Solution 3: Check Key Permissions**

```bash
# Your key file should have 600 permissions
chmod 600 your-key.pem

# Verify
ls -la your-key.pem
# Should show: -rw------- (600)
```

#### **Solution 4: Convert Key Format**

Some keys need conversion:

```bash
# If you have a .ppk file (PuTTY format)
puttygen your-key.ppk -O private-openssh -o your-key.pem

# If key has passphrase, remove it
ssh-keygen -p -f your-key.pem
# Enter old passphrase, then press Enter twice for no passphrase
```

#### **Solution 5: Re-add the Key Properly**

1. **Get the key as one continuous string:**
```bash
# On Mac/Linux, this shows the key with actual newlines
cat your-key.pem
```

2. **In SSH Manager:**
   - Edit your server (press `e`)
   - Go to PEM field
   - Press `p` to open editor
   - **Copy the ENTIRE key output from cat command**
   - Paste it (Ctrl+Shift+V in most terminals)
   - Press Ctrl+S
   - Submit

### Error: "Permission denied (publickey)"

This means authentication failed. Several possible causes:

#### **Cause 1: Public Key Not on Server**

```bash
# Your PUBLIC key needs to be on the server
# Check if it's there
ssh ubuntu@16.171.3.225 cat ~/.ssh/authorized_keys

# If you need to add it:
# 1. Generate public key from private key
ssh-keygen -y -f your-key.pem > your-key.pub

# 2. Copy it to server
ssh-copy-id -i your-key.pub ubuntu@16.171.3.225
# OR manually add it to ~/.ssh/authorized_keys on server
```

#### **Cause 2: Wrong Username**

Common usernames for different systems:
- **AWS EC2**: `ec2-user` (Amazon Linux), `ubuntu` (Ubuntu), `admin` (Debian)
- **DigitalOcean**: `root`
- **Google Cloud**: Your Google username
- **Azure**: Whatever you set during creation

```bash
# Try different usernames
ssh -i your-key.pem ec2-user@16.171.3.225
ssh -i your-key.pem ubuntu@16.171.3.225
ssh -i your-key.pem admin@16.171.3.225
```

#### **Cause 3: Server-Side Permissions**

The server's `~/.ssh` folder needs correct permissions:

```bash
# Connect to server (using password or another method)
# Then fix permissions:
chmod 700 ~/.ssh
chmod 600 ~/.ssh/authorized_keys
```

#### **Cause 4: Wrong Key**

Make sure you're using the correct private key that matches the public key on the server.

### Error: "Connection timeout"

```bash
# Check if server is reachable
ping 16.171.3.225

# Check if SSH port is open
nc -zv 16.171.3.225 22
# or
telnet 16.171.3.225 22
```

## Testing Your PEM Key

### Step 1: Test Manually First

**Always test your key manually before adding to SSH Manager:**

```bash
# Test the connection
ssh -i your-key.pem -v ubuntu@16.171.3.225

# The -v flag shows verbose output
# Look for lines like:
# - "Offering public key: your-key.pem"
# - "Server accepts key"
# - "Authentication succeeded"
```

### Step 2: Check the Temp Key File

SSH Manager creates temp files in `/tmp/termius-from-walmart-keys/`. Check if they're valid:

```bash
# After trying to connect, check the temp file
ls -la /tmp/termius-from-walmart-keys/

# View the temp key
cat /tmp/termius-from-walmart-keys/key_2.pem

# Compare with your original
diff your-key.pem /tmp/termius-from-walmart-keys/key_2.pem
```

### Step 3: Verify Key Format

```bash
# Check key type
ssh-keygen -l -f your-key.pem

# Should show something like:
# 2048 SHA256:xxx... (RSA)
# or
# 256 SHA256:xxx... (ED25519)
```

## Common PEM Key Issues

### Issue: Key has Windows line endings

```bash
# Convert CRLF to LF
dos2unix your-key.pem

# Or using sed
sed -i 's/\r$//' your-key.pem
```

### Issue: Key is encrypted with passphrase

SSH Manager doesn't support passphrase-protected keys yet.

```bash
# Remove passphrase
ssh-keygen -p -f your-key.pem
# Enter current passphrase
# Press Enter twice for no new passphrase
```

### Issue: Key has extra whitespace

```bash
# Clean the key
# Make sure there are no spaces before -----BEGIN
# Make sure there are no spaces after -----END
# Make sure each line of base64 is exactly as it should be
```

## For Your Specific Error

Based on your error with `ubuntu@16.171.3.225`, try this:

### Quick Fix:

```bash
# 1. Test your key manually first
ssh -i your-key.pem -v ubuntu@16.171.3.225

# 2. If that works, get the exact key content
cat your-key.pem

# 3. Copy the ENTIRE output including -----BEGIN and -----END lines

# 4. In termius-from-walmart:
#    - Edit the server
#    - Clear the PEM field completely
#    - Press 'p' to open editor
#    - Paste the entire key
#    - Press Ctrl+S
#    - Submit

# 5. Try connecting again
```

### Alternative: Use the key file path instead

If pasting doesn't work, you can modify the app to accept file paths, or just use:

```bash
# Create a symlink
ln -s ~/path/to/your-key.pem ~/.ssh/server-key.pem
chmod 600 ~/.ssh/server-key.pem

# Then in SSH Manager, leave PEM field empty
# SSH will automatically try keys in ~/.ssh/
```

## Debug Mode

To see what SSH is actually doing:

```bash
# Check what's in the temp key file
cat /tmp/termius-from-walmart-keys/key_2.pem

# Try connecting with verbose SSH
ssh -vvv -i /tmp/termius-from-walmart-keys/key_2.pem ubuntu@16.171.3.225

# This will show you exactly why authentication is failing
```

## Still Not Working?

1. **Share the output of:**
```bash
ssh -vvv -i your-key.pem ubuntu@16.171.3.225
```

2. **Check if the key format is correct:**
```bash
head -n 1 your-key.pem
# Should be exactly: -----BEGIN RSA PRIVATE KEY-----
# or: -----BEGIN OPENSSH PRIVATE KEY-----
```

3. **Verify server allows public key auth:**
```bash
# On the server, check sshd_config
sudo grep PubkeyAuthentication /etc/ssh/sshd_config
# Should show: PubkeyAuthentication yes
```

---

**Most Common Solution**: The key is correct, but the public key isn't on the server. Add your public key to the server's `~/.ssh/authorized_keys` file.
