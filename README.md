# ghokin-json

Fork of [ghokin](https://github.com/antham/ghokin) that treats docstrings as JSON by default. You'll need `jq`.

I've removed some tests so that this builds, added one, and added some Nix stuff.

The rest of this README is the original.

Ghokin format and apply transformation on gherkin files.

---

- [Install](#install)
- [Usage](#usage)
- [Documentation](#documentation)
- [Contribute](#contribute)

---

## Usage

```
Clean and/or apply transformation on gherkin files

Usage:
  ghokin [command]

Available Commands:
  check       Check a file/folder is well formatted
  fmt         Format a feature file/folder
  help        Help about any command

Flags:
      --config string   config file
  -h, --help            help for ghokin

Use "ghokin [command] --help" for more information about a command.
```

⚠️ Ghokin works only on `UTF-8` encoded files, it will detect and convert automatically files that are not encoded in this charset.

### fmt stdout

Dump stdin or a feature file formatted on stdout

```
ghokin fmt stdout features/test.feature
```

or

```
cat features/test.feature|ghokin fmt stdout
```

### fmt replace

Format and replace a file or all files in a directory

```
ghokin fmt replace features/test.feature
```

or

```
ghokin fmt replace features/
```

### check

Ensure a file or all files in a directory are well formatted, exit with an error code otherwise

```
ghokin check features/test.feature
```

or

```
ghokin check features/
```

## Documentation

### Shell commands

You can run shell commands from within your feature file to transform some datas with annotations, to do so you need first to define in the config an alias and afterwards you can simply "comment" the line before the line you want to transform with that alias.
For instance let say `@json` calls behind the curtain `jq`, we could validate and format some json in our feature like so :

```
Feature: A Feature
  Description

  Scenario: A scenario to test
    Given a thing
    # @json
    """
    {
      "test": "test"
    }
    """
```

### Config

Defaut config is to use 2 spaces for indentation.

It's possible to override configuration by defining a `.ghokin.yml` file in the home directory or in the current directory where we are running the binary from :

```
indent: 2
aliases:
  json: "jq ."
```

Aliases key defined [shell commands](#shell-commands) callable in comments as we discussed earlier.

It's possible to use environments variables instead of a static config file :

```
export GHOKIN_INDENT=2
export GHOKIN_ALIASES='{"json":"jq ."}'
```

## Contribute

If you want to add a new feature to ghokin project, the best way is to open a ticket first to know exactly how to implement your changes in code.

### Setup

After cloning the repository you need to install vendors with `go mod vendor`
To test your changes locally you can run go tests with : `make test-all`
