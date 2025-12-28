---
title: Troubleshooting Classification in FolderFlow
sidebar_label: Troubleshooting
sidebar_position: 6
---

# 🛠️ Troubleshooting Classification

Encountering issues with classification in FolderFlow? This page provides solutions to common problems and guidance for resolving errors.

---

## 🔍 Common Issues and Solutions

### **1. No Files Are Being Classified**

#### **Symptoms:**
- FolderFlow runs without errors, but no files are moved to the destination directories.

#### **Possible Causes and Solutions:**

| Cause                          | Solution                                                                                     |
|--------------------------------|----------------------------------------------------------------------------------------------|
| Incorrect source directory     | Verify that the `source_dirs` paths in your configuration are correct and accessible.         |
| No matching files              | Ensure that files in the source directory match the defined filters (extensions, patterns, etc.). |
| Incorrect file permissions     | Check that FolderFlow has read access to the source directory and write access to the destination. |
| Misconfigured filters          | Review your filter conditions in the YAML configuration file.                               |

#### **Steps to Debug:**
1. **Check Logs**: Look for warnings or errors in the FolderFlow logs.
2. **Test with Simple Rules**: Temporarily simplify your configuration to isolate the issue.
3. **Verify File Matching**: Manually check if files in the source directory match your filter conditions.

---

### **2. Files Are Moved to the Wrong Destination**

#### **Symptoms:**
- Files are being classified but end up in the wrong destination directory.

#### **Possible Causes and Solutions:**

| Cause                          | Solution                                                                                     |
|--------------------------------|----------------------------------------------------------------------------------------------|
| Overlapping filters            | Ensure that filters are mutually exclusive or use priorities to resolve conflicts.           |
| Incorrect filter conditions    | Double-check the conditions in your filters (e.g., extensions, patterns).                     |
| Misconfigured destination paths | Verify that the `path` field in your destination directories is correct.                      |

#### **Steps to Debug:**
1. **Review Filter Priorities**: If multiple filters match a file, the one with the highest priority is applied.
2. **Test Individual Filters**: Temporarily disable other filters to identify conflicts.
3. **Check Logs**: Look for messages indicating which filter was applied to a file.

---

### **3. Errors Creating Hard Links or Symlinks**

#### **Symptoms:**
- FolderFlow reports errors when creating hard links or symlinks during the regrouping phase.

#### **Possible Causes and Solutions:**

| Cause                          | Solution                                                                                     |
|--------------------------------|----------------------------------------------------------------------------------------------|
| Insufficient permissions       | Ensure FolderFlow has write permissions in the regroup directory.                            |
| Files on different filesystems | Hard links require files to be on the same filesystem. Use symlinks if necessary.             |
| Unsupported filesystem         | Some filesystems (e.g., FAT32) do not support hard links or symlinks.                          |

#### **Steps to Debug:**
1. **Check Filesystem Type**: Use `df -T` (Linux/macOS) or `fsutil fsinfo volumeinfo` (Windows) to verify the filesystem.
2. **Test with Symlinks**: If hard links fail, try using `mode: symlink` in your regroup configuration.
3. **Verify Permissions**: Ensure the user running FolderFlow has the necessary permissions.

---

### **4. Performance Issues**

#### **Symptoms:**
- FolderFlow is slow when processing a large number of files.

#### **Possible Causes and Solutions:**

| Cause                          | Solution                                                                                     |
|--------------------------------|----------------------------------------------------------------------------------------------|
| Too many workers               | Adjust `max_workers` in your configuration (set to `0` for auto-detection).                 |
| Complex filters                | Simplify filters or break them into smaller, more efficient rules.                           |
| Large files                    | Exclude very large files from classification if they are not critical.                       |

#### **Steps to Debug:**
1. **Monitor Resource Usage**: Use tools like `top` (Linux/macOS) or Task Manager (Windows) to check CPU and memory usage.
2. **Adjust Workers**: Start with `max_workers: 0` and gradually increase if needed.
3. **Profile Filters**: Test with a subset of files to identify slow filters.

---

### **5. Files Are Duplicated Instead of Linked**

#### **Symptoms:**
- Files appear to be duplicated in the regroup directory instead of being linked.

#### **Possible Causes and Solutions:**

| Cause                          | Solution                                                                                     |
|--------------------------------|----------------------------------------------------------------------------------------------|
| Incorrect regroup mode         | Ensure `mode: hardlink` or `mode: symlink` is set in your configuration.                       |
| Filesystem limitations         | Hard links require the same filesystem; use symlinks if files are on different drives.       |

#### **Steps to Debug:**
1. **Verify Regroup Configuration**: Check that `mode` is set to `hardlink` or `symlink`.
2. **Test with a Single File**: Manually verify that links are created correctly for a single file.

---

## 📝 Logs and Debugging

### **Where to Find Logs**
FolderFlow logs are typically located in:
- `/var/log/folderflow/` (Linux/macOS)
- `%APPDATA%\FolderFlow\logs\` (Windows)

### **Log Levels**
Adjust the log level in your configuration to get more details:
```yaml
log_level: debug  # Options: error, warn, info, debug
```
#### Common Log Messages
|Message|Meaning|
|:-|:-|
|No files matched filter X|No files in the source directory matched the conditions of filter X.|
|Permission denied: /path/to/file|FolderFlow lacks permissions to access the file or directory.|
|Hard link creation failed|Failed to create a hard link (check filesystem and permissions).|

## 💡 Tips for Effective Troubleshooting

1. Start Small: Test with a small set of files and simple rules before scaling up.
2. Isolate Issues: Temporarily disable parts of your configuration to identify the root cause.
3. Check Permissions: Ensure FolderFlow has the necessary read/write permissions.
4. Review Logs: Logs often contain clues about what went wrong.
5. Update FolderFlow: Ensure you are using the latest version of FolderFlow.


## 📢 Need More Help?
If you’re still encountering issues, consider:

- Opening an issue on [GitHub](https://github.com/polocto/FolderFlow/issues).
