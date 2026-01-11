# Test your plugin

For each plugin you add. You must add a test to see if it works as excpected.

Under the corresponding folder for your plugin (`fileters` or `strategy`), create a folder that you'll name after your plugin's. 

The architecture of your test foler must follow this schem:

```text
plugin_name/
├── expected/
└── config.yaml

```

## `config.yaml` requirements

> 1. Paths in source_dirs and dest_dirs must be relative to the test folder, `testdata/`.
> 2. For the test to work, at least one file in source_dirs must match the filters specified in the configuration.
> 3. Dummy files can be created automatically by the test if necessary.
> 4. The test will compare the contents of expected/ with the output generated in the temporary result directory.
> 5. To make your plugin test work correctly, the config.yaml file must include at least the following sections:

### Source Directories

```yaml
source_dirs:
  - "testdata/input/<dir>"   # relative path to the folder containing files to working directory
```

- You can specify multiple source directories.
- Each folder will be copied into a temporary directory during the test.

### Destination Directories

```yaml
dest_dirs:
  - name: "documents"
    path: "destination/documents"
    filters:
      - name: "extensions"   # Name of the filter to test or 'extensions' as default
        config:
          extensions: [".txt", ".pdf", ".md"]   # Plugin configuration if needed
    strategy:
      name: "dirchain"       # Name of the strategy to test or 'dirchain' as default
```

- name: logical name for the destination folder.
- path: path that will be create in a tmp directory that should match the path under expected
- filters: a list of filters applied to this destination. Each filter must have a name and config.
- strategy: the classification strategy applied to this folder.

You can define multiple dest_dirs to test different types of filters or strategies.

### Regrouping (Optional)

```yaml
regroup:
  path: "regrouped"  # Folder where all processed files will also be grouped
  mode: hardlink     # Mode: "copy", "move", or "hardlink"
```

- regroup is optional. If present, all files processed by your plugin will also be copied or hardlinked into this folder.
- Useful for verifying that files are accessible after processing.


## Expected Folder

The `expected/` folder should contain the folder structure and files that you expect your plugin to produce.

- The Go test will generate files in a temporary directory using the source directories defined in `config.yaml`.
- After classification, the test compares the generated output with `expected/` using file paths and `SHA256` hashes.

For example, if your configuration defines `documents` and `regrouped`, your `expected/` folder could look like:

```text
expected/
├── documents/
│   ├── example.txt
│   └── notes.md
└── regrouped/
    ├── example.txt
    └── notes.md
```