---
title: Introduction
sidebar_label: Introduction
sidebar_position: 1
---

# 📂 FolderFlow

**FolderFlow** is a tool for **automating file organization** using clear, declarative rules. It is built for people who want their file systems to stay organized over time—without relying on manual sorting or fragile scripts.

Rather than reacting to mess, FolderFlow lets you **define structure once** and apply it consistently.

---

## 🧠 What FolderFlow Is

FolderFlow is a **filesystem automation engine**.

It scans directories, evaluates files against rules, and applies well-defined actions to organize them. Everything it does is driven by configuration, not code, making its behavior easy to understand, review, and reproduce.

FolderFlow focuses on:
- Predictability over magic
- Configuration over scripting
- Explicit behavior over hidden defaults

---

## ❌ What FolderFlow Is Not

FolderFlow is intentionally not:
- A background file watcher or daemon
- A GUI-based file manager
- An AI-driven or heuristic sorter
- A one-off cleanup script

It runs when *you* run it, does exactly what you configured, and then stops.

---

## 🧩 How FolderFlow Thinks About Files

FolderFlow treats file organization as a **classification problem**:

- Files come from defined locations
- Filters decide how they should be grouped
- Strategies define how structure is created
- Optional linking provides alternative access paths

This model scales cleanly from a single downloads folder to large datasets.

---

## 🧾 Configuration-Driven by Design

FolderFlow is fully controlled by a single **YAML configuration file**.

This configuration:
- Describes *what* should happen, not *how*
- Is human-readable and version-control friendly
- Can be reviewed, shared, and audited
- Produces the same results every time it is run

This makes FolderFlow suitable for automation, repeatable workflows, and long-term maintenance.

---

## 🛡️ Predictability & Safety

FolderFlow is designed to be safe by default:

- No silent overwrites
- No implicit behavior
- Deterministic rule evaluation
- Clear logs and error messages

You always know *why* a file was handled a certain way.

---

## 🖥️ Where FolderFlow Fits

FolderFlow works well when:
- Manual file organization no longer scales
- You want consistent structure across machines or datasets
- Scripts have become too complex to maintain
- You care about reproducibility and clarity

It fits equally well in personal setups and technical workflows.

---

## 📚 Where to Go Next

Once you understand the philosophy, the next step is to explore **features**.

Start with **Classification**, which is the primary way FolderFlow organizes files:

👉 **[Classification](./classification)**

From there, you can dive into configuration details, examples, and advanced use cases.

---

## 🧭 Philosophy

FolderFlow follows a simple principle:

> *If you can clearly describe how your files should be organized, FolderFlow should be able to do it — transparently and repeatedly.*

That’s the foundation everything else is built on.
