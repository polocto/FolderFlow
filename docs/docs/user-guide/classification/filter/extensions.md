---
title: extensions
sidebar_position: 1
---

# Extensions Filter

The extensions filter matches files based on their file extension.

---

## Selector name

extensions

---

## Configuration

This filter matches JPEG and PNG images.

```yaml
filters:
  - name: "extensions"
    config:
      extensions:
        - ".jpg"
        - ".png"
        - ".jpeg"
```

---

## Behavior

- Only file extensions are inspected
- Directory names are ignored
- Matching is case-insensitive
- The leading dot is required

If a file’s extension matches one of the configured values, the filter returns true.

---

## Notes

- Files without extensions never match
- Extensions must include the leading dot
- The filter does not inspect file content