# SFTP Split-Screen File Manager - Quick Start

## The Experience

When you open SFTP in Termius from Walmart, you get a **professional split-screen interface** just like the real Termius app:

```
SFTP: user@192.168.1.100

> LOCAL <              REMOTE
/home/user             /var/www/html

config.yaml            index.php
backup.tar.gz          assets/
uploads/               config/
documents/             logs/
../                    ../

Keys: [Tab] switch pane • [c]opy • [d]elete • [enter] navigate • [q]uit
```

## Main Commands

| What You Want | How To Do It |
|--------------|------------|
| **Upload a file** | Select file (left) → [Tab] → Navigate folder (right) → [Tab] → [c] |
| **Download a file** | [Tab] → Find file (right) → [Tab] → Navigate folder (left) → [Tab] → Select → [c] |
| **Navigate folder** | Press [enter] on directory name |
| **Go back up** | Press [enter] on `../` |
| **Delete file** | Select file → [d] |
| **Switch panes** | Press [Tab] anytime |

## Step-by-Step: Upload a Config File

```
1. Server List appears
2. Select your server
3. Press [s] for SFTP
   → App connects...
   → Shows LOCAL: /home/user
   → Shows REMOTE: /var/www

4. In LOCAL pane (left), select "myconfig.yaml"
5. Press [Tab] to focus REMOTE pane (right)
6. Press [enter] on "config/" folder
7. Press [Tab] back to LOCAL pane
8. Press [c] to copy
   → Progress bar appears: [████░░░░] 50%
   → Success: "Uploaded myconfig.yaml"

9. Now both panes show the file
10. Done! Press [q] to exit SFTP
```

## What Makes This Great

✨ **Like Termius**
- Professional split-screen interface
- Intuitive navigation
- Real-time feedback

⚡ **Fast & Efficient**
- Upload and download in seconds
- No file size limits
- Works with large files

🔒 **Secure**
- Uses your SSH credentials
- Encrypted connection
- Proper permission handling

📊 **Smart Progress**
- See transfer percentage
- Current operation displayed
- Error messages if needed

## Key Features

### Left Pane (LOCAL)
- Browse your computer
- Browse from home directory (~)
- Can go to any accessible folder
- Delete local files

### Right Pane (REMOTE)
- Browse remote VPS
- Browse from home directory
- Navigate any accessible folder
- Delete remote files

### Both Panes
- File and folder browsing
- Up/down arrow navigation
- Enter to navigate folders
- `../` to go back

### Transfer Operations
- **[c]** Copy from active pane to other
- **[d]** Delete in active pane
- Progress updates in real-time
- Status messages confirm operations

## What to Try

```
BEGINNER:
1. Open SFTP on a server
2. See files in both panes
3. Press [Tab] to switch panes
4. Press arrows to browse files
5. Press [q] to exit

INTERMEDIATE:
1. Navigate to a folder in both panes
2. Copy a file from local to remote ([c])
3. Watch the progress bar
4. See the file appear in remote
5. Copy it back with [c] again

ADVANCED:
1. Open multiple SFTP sessions (one at a time)
2. Copy files between servers
3. Organize files by navigating
4. Delete old files you don't need
5. Use as your daily file manager
```

## Common Scenarios

### Upload Your Web App
```
1. Open SFTP on web server
2. Navigate LOCAL: project/ folder
3. [Tab] to REMOTE, go to /var/www/
4. [Tab] back, select app.js file
5. [c] to upload
6. Done!
```

### Backup Important Files
```
1. Open SFTP on backup server
2. [Tab] to REMOTE, go to /backups/
3. [Tab] back to LOCAL
4. [c] to copy important files
5. Watch progress bar
6. Files are now backed up!
```

### Clean Up Old Files
```
1. Open SFTP on server
2. Navigate to old files
3. Select with arrows
4. Press [d] to delete
5. Repeat as needed
6. Done!
```

## Help & Troubleshooting

### Can't Connect
- Check server credentials
- Make sure SFTP is enabled (usually is)
- Verify SSH port is accessible

### Upload/Download Fails
- Check available disk space
- Verify file permissions
- Try smaller file first
- Check internet connection

### Can't See Files
- Wrong directory? Press [enter] on folders
- Files might be hidden (starts with .)
- Check read permissions

### Questions?
- See [SFTP_SPLIT_SCREEN.md](SFTP_SPLIT_SCREEN.md) for detailed guide
- See [SFTP_VISUAL_GUIDE.md](SFTP_VISUAL_GUIDE.md) for visual examples
- Check main README.md for general help

---

**Happy file managing! 🚀**
