---
title: Introduction to Classification in FolderFlow
sidebar_label: Classification
---

# 🗂️ Introduction to Classification

**Classification** is a powerful feature in **FolderFlow** that allows you to **automate the organization of your files** based on customizable rules. Whether you're managing photos, documents, videos, or any other type of file, FolderFlow’s classification system helps you **keep everything tidy, accessible, and well-structured** without manual effort.

---

## 🎯 Why Use Classification?

Classification in FolderFlow is designed to:
- **Automate file organization**: Move files to the right folders automatically based on predefined rules.
- **Save time**: Reduce manual sorting and focus on what matters.
- **Improve efficiency**: Find files faster with a logical and consistent structure.
- **Customize workflows**: Adapt the system to your specific needs with flexible rules and strategies.

---

## 🔧 How Classification Works

FolderFlow’s classification system is built around three core concepts:

### **1. Source Directories**
Define the directories where FolderFlow will look for files to classify. These are the folders you want to **monitor and organize**.

### **2. Filters**
Filters determine **which files** should be moved to a destination directory. You can define filters based on:
- **File extensions** (e.g., `.jpg`, `.pdf`).
- **File name patterns** (e.g., `invoice_*`, `2023_*`).
- **File metadata** (e.g., creation date, size).
- **Custom tags** (e.g., `#projectX`, `#urgent`).

For a file to be selected, **all conditions in a filter must be met**.

### **3. Destination Directories and Strategies**
Define **where** files should be moved and **how** they should be organized. FolderFlow supports:
- **Simple strategies** (e.g., `dirchain` for direct placement).
- **Advanced strategies** (e.g., organizing by date, extension, or custom patterns).

---

## 📌 Key Features

### **Flexible Filtering**
- Use **`extensions`** and **`regex`**.
- Combine multiple filters to create complex selection criteria.

### **Customizable Strategies**
- Choose from predefined strategies like **`dirchain`** or **`date`**.
- Create **custom strategies** to match your unique workflow.

### **Efficient File Handling**
- Use **hard links** or **symbolic links** to avoid file duplication while maintaining accessibility.
- Process files in parallel with **multi-worker support** for faster performance.

### **Easy Configuration**
- Define your rules in a simple **YAML file**.
- Test and refine your configuration with minimal effort.

---

## 📂 Getting Started with Classification

To start using the classification feature in FolderFlow, follow these steps:

1. **Define Source Directories**: Specify the folders you want to monitor.
2. **Set Up Filters**: Create rules to select the files you want to classify.
3. **Configure Destination Directories**: Define where files should be moved and how they should be organized.
4. **Run FolderFlow**: Let FolderFlow handle the rest!

For a step-by-step guide, check out the **[Configuration Page](/user-guide/features/classification/configuration)**.

---

## 📚 Sub-Pages

Explore the following sub-pages to learn more about classification in FolderFlow:

- **[Configuration](/user-guide/features/classification/configuration)**: Learn how to set up your classification rules in the YAML file.
- **[Examples](/user-guide/features/classification/examples)**: See real-world examples of classification rules.
- **[Advanced Use Cases](/user-guide/features/classification/advanced)**: Discover advanced strategies and custom workflows.
- **[Running Classification](./run)**: Once your configuration file is ready, you can run FolderFlow.
- **[Troubleshooting](/user-guide/features/classification/troubleshooting)**: Find solutions to common issues.

---

## 🔄 Example Workflow

Here’s a simple example of how classification works in FolderFlow:

1. **Source Directory**: `/Users/you/Downloads/`
2. **Filters**:
   - Move all `.jpg` and `.png` files to `Photos/`.
   - Move all `.pdf` files to `Documents/`.
   - Move all `.mp4` files to `Videos/`.
3. **Destination Strategy**: Use `dirchain` to place files directly in their respective folders.

For more examples, visit the **[Examples Page](/user-guide/features/classification/examples)**.

---

## 🚀 Best Practices

1. **Start Simple**: Begin with basic rules and gradually add complexity.
2. **Test Your Rules**: Use a small set of files to validate your configuration before applying it to large datasets.
3. **Monitor Logs**: Check the logs to ensure files are being classified correctly.
4. **Backup Files**: Always back up your files before running classification.
5. **Review Regularly**: Update your rules and strategies as your needs evolve.

---

## ❓ Troubleshooting

If you encounter issues with classification, check out the **[Troubleshooting Page](/user-guide/features/classification/troubleshooting)** for common solutions.

---

## 📢 Stay Updated

Follow FolderFlow on **[GitHub](https://github.com/polocto/FolderFlow)** to stay up-to-date with the latest features, improvements, and community discussions.

---

## 🙌 Conclusion

The **Classification** feature in FolderFlow is a powerful tool for **automating file organization** and keeping your digital workspace tidy. By defining flexible rules and strategies, you can create a system that works for your specific needs.

Ready to get started? Head over to the **[Configuration Page](/user-guide/features/classification/configuration)** to set up your first classification rules!
