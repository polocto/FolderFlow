---
title: Advanced Classification in FolderFlow
sidebar_label: Advanced Use Cases
sidebar_position: 4
---

# 🚀 Advanced Classification

Take your file organization to the next level with **advanced classification techniques** in FolderFlow. This page covers complex filtering, custom strategies, and integration with other tools.

---

## 🎯 Advanced Filtering Techniques

### **1. Combining Multiple Filters**

You can combine multiple filters to create **complex selection criteria**. For example, select files based on **both extension and name pattern**:

```yaml
dest_dirs:
  - name: "project_documents"
    path: "./documents/projects"
    filters:
      - name: "extensions"
        config:
          extensions: [".pdf", ".docx"]
      - name: "regex"
        config:
          pattern: "project_.*"
    strategy:
      name: "dirchain"
```

In this example, **only `.pdf` or `.docx` files with names starting with `project_`** will be moved to `./documents/projects`.

### **2. Using Metadata Filters**
FolderFlow supports filtering based on **file metadata** such as creation date, modification date, and file size.

**Example: Filter by Creation Date**
```yaml
dest_dirs:
  - name: "old_files"
    path: "./archive/old"
    filters:
      - name: "creation_date"
        config:
          before: "2023-01-01"  # Files created before January 1, 2023
    strategy:
      name: "dirchain"
```

**Example: Filter by File Size**
```yaml
dest_dirs:
  - name: "large_files"
    path: "./large_files"
    filters:
      - name: "size"
        config:
          min: "10MB"  # Files larger than 10MB
    strategy:
      name: "dirchain"
```

### **3. Custom Tags**

Use **custom tags** to classify files based on user-defined labels. Tags can be added manually or automatically using scripts.

**Example: Filter by Tag**

```yaml
dest_dirs:
  - name: "urgent_documents"
    path: "./priority/urgent"
    filters:
      - name: "tag"
        config:
          tag: "#urgent"
    strategy:
      name: "dirchain"
```

## 📁 Advanced Destination Strategies


### **1. Dynamic Paths with Variables**
Use variables in destination paths to create dynamic folder structures. For example, organize files by year and month based on their creation date:

```yaml
dest_dirs:
  - name: "monthly_documents"
    path: "./documents/{year}/{month}"
    filters:
      - name: "extensions"
        config:
          extensions: [".pdf", ".docx"]
    strategy:
      name: "by_date"
```

In this example, `{year}` and `{month}` are dynamically replaced by the file's creation date.

### **2. Custom Strategies**
Create **custom strategies** to implement complex organization logic. For example, organize files by **project name** extracted from the filename:


```yaml
dest_dirs:
  - name: "project_files"
    path: "./projects/{project_name}"
    filters:
      - name: "pattern"
        config:
          pattern: "project_(?P<project_name>.*)_.*"
    strategy:
      name: "custom"
      config:
        path_template: "{project_name}"
```

In this example, files named like `project_alpha_report.pdf` will be moved to `./projects/alpha/`.


### **3. Multi-Destination Classification**

Classify files into **multiple destination directories** based on different criteria. For example, a `.pdf` file could be moved to both `Documents/` and `Projects/` if it matches multiple filters.

```yaml
dest_dirs:
  - name: "documents"
    path: "./documents"
    filters:
      - name: "extensions"
        config:
          extensions: [".pdf"]
    strategy:
      name: "dirchain"

  - name: "projects"
    path: "./projects"
    filters:
      - name: "pattern"
        config:
          pattern: "project_*"
    strategy:
      name: "dirchain"
```

## 🔄 Advanced Regrouping Techniques

