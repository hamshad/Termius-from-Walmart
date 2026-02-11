# Split-Screen SFTP - Quick Visual Guide

## The Interface Layout

```
SFTP: user@192.168.1.100
[100%] Uploading myfile.txt

> LOCAL <              REMOTE
/home/user             /var/www

file1.txt              index.html
file2.txt              config.php
docs/                  uploads/
downloads/             backup/

Keys: [Tab] switch pane • [c]opy • [d]elete • [enter] navigate • [q]uit
```

## Copy Local → Remote (Upload)

```
STEP 1: Select file in LOCAL       STEP 2: Navigate in REMOTE
┌──────────────────────┐          ┌──────────────────────┐
│ > LOCAL < │ REMOTE   │          │ LOCAL │ > REMOTE <  │
│ /home/user │ /root   │          │ /home │ /var/www    │
│ ▶ file.txt │ config/ │          │ file  │ ▶ public/   │
│   doc/     │ script/ │          │ code/ │   private/  │
└──────────────────────┘          └──────────────────────┘
        [Tab to switch]                  [Enter to navigate]

STEP 3: Copy!
┌──────────────────────┐
│ LOCAL │ > REMOTE <   │
│ file.txt │ /var/www │
│ doc/     │           │
│ code/    │           │
│                      │
│ [50%] Uploading...   │
└──────────────────────┘
        [c] key pressed
```

## Copy Remote → Local (Download)

```
STEP 1: Navigate REMOTE      STEP 2: Switch to LOCAL
┌────────────────────┐      ┌────────────────────┐
│ LOCAL │ > REMOTE   │      │ > LOCAL │ REMOTE   │
│ /home │ /var/www   │      │ /home   │ /var/www │
│ doc/  │ ▶ file.txt │      │ ▶ down/ │ file.txt │
│ code/ │   image.jp │      │ code/   │ image.jp │
└────────────────────┘      └────────────────────┘
    [Tab to switch]          [Enter to navigate]

STEP 3: Select & Copy!
┌────────────────────┐
│ > LOCAL │ REMOTE   │
│ /down   │ /var/www │
│         │          │
│ [50%] Downloading  │
└────────────────────┘
```

## Keyboard Commands Reference

```
┌─────────────────────────────────────────────────┐
│ SFTP COMMAND REFERENCE                          │
├─────────────────────────────────────────────────┤
│ [Tab]         → Switch between Local/Remote     │
│ [↑] [↓]       → Navigate files                  │
│ [enter]       → Enter directory                 │
│ [c]           → Copy file to other side         │
│ [d]           → Delete selected file            │
│ [q] / [esc]   → Exit SFTP, return to list       │
└─────────────────────────────────────────────────┘
```

## File Symbols

```
file.txt          Regular file
directory/        Directory (with / suffix)
../               Parent directory
```

## Status Messages

```
Uploading myfile.txt
[████████████░░░░░] 75%

Uploaded config.php

Error downloading: Permission denied

Deleted old_backup.tar
```

## Active Pane Indicator

```
> LOCAL < │ REMOTE        (Local pane has focus)

LOCAL │ > REMOTE <        (Remote pane has focus)
```

## Flow Diagrams

### Upload Workflow
```
Select file in LOCAL
       ↓
[Tab] to REMOTE
       ↓
Navigate to destination
       ↓
[Tab] back to LOCAL
       ↓
[c] to copy
       ↓
Watch progress [■■■░░] 60%
       ↓
File appears in REMOTE ✓
```

### Download Workflow
```
[Tab] to REMOTE
       ↓
Find file to download
       ↓
[Tab] back to LOCAL
       ↓
Navigate to save location
       ↓
[Tab] to REMOTE
       ↓
Select file, [c] to copy
       ↓
Watch progress [■■■■■░] 90%
       ↓
File appears in LOCAL ✓
```

### Delete Workflow
```
Select file (Local or Remote)
       ↓
Press [d]
       ↓
File deleted ✓
       ↓
List refreshes automatically
```

## Tips & Tricks

```
Pro Tip #1: Navigate to your destination FIRST
  Before copying, make sure you're in the right
  directory on the destination pane!

Pro Tip #2: Use [Tab] freely
  Switch between panes as many times as you need.
  The connection stays active.

Pro Tip #3: Check paths carefully
  Notice the path display at top of each pane
  to avoid copying to wrong location.

Pro Tip #4: Watch the progress bar
  Don't close terminal while file is transferring.
  Wait for completion message.
```

## Error Recovery

```
If connection drops:
  → Error message appears
  → Navigate back to list with [q]
  → Try opening SFTP again

If copy fails:
  → Error message shows reason
  → Fix issue (permissions, space, etc.)
  → Try again with [c]

If file stuck:
  → Wait a moment
  → Check connection
  → Try [q] to exit and reconnect
```
