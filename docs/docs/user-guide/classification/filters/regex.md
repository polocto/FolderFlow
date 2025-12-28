---
title: regex
sidebar_position: 2
---

# Regex Filter

The regex filter selects files based on regular expression matching against the file name.

It allows fine-grained control over file selection using patterns instead of simple extensions.

---

## Selector name

regex

---

## Configuration

The regex filter is configured using a list of regular expression patterns.

Each pattern is applied to the **file name only** (not the full path).

```yaml
filters:
  - name: "regex"
    config:
      patterns:
        - "^report_.*\\.pdf$"
        - ".*_final\\.docx$"
```

---

## Behavior
- Each pattern is compiled as a regular expression
- Patterns are matched against the file base name
- Matching is case-sensitive unless the pattern specifies otherwise
- A file matches if any pattern matches
- Directory paths are not evaluated

If at least one pattern matches, the filter returns true.

### Example

Given the following files:
```text
source/reports/report_2024.pdf
source/reports/draft.docx
source/reports/summary_final.docx
```

Matched files:
- report_2024.pdf
- summary_final.docx

Ignored files:
- draft.docx