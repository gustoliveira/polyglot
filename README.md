# Polyglot

Polyglot is a CLI tool to manage strings in Android resource files and automate translations using the Google Translate API. It helps you:

1. **Check** if translations are normalized (e.g., sorted by key), discover potentially unused string resources, and incosistencies keys between resource.
2. **Normalize** translations by sorting keys automatically.
3. **Remove** string keys across all language files at once.
4. **Translate** new or existing strings to multiple locales using Google Translate.

Below are detailed instructions on how to build, install, configure, and use **Polyglot**.

---

## Table of Contents
1. [Features](#features)
2. [Installation](#installation)
	- [Install from Releases](#install-from-releases)
	- [Build from git repository](#build-from-git-repository)
3. [Configuration](#configuration)
4. [Usage](#usage)
   - [Available Commands](#available-commands)
     - [check](#check)
     - [normalize](#normalize)
     - [remove](#remove)
     - [translate](#translate)
5. [Advanced Topics](#advanced-topics)
   - [Adding New Subcommands](#adding-new-subcommands)
   - [Android Project Detection](#android-project-detection)
6. [License](#license)

---

## Features
- **Sorting**: Ensures all string keys in `strings.xml` are alphabetically sorted.
- **Translation**: Integrates with the Google Translate API to generate localized strings automatically.
- **Cleaning Up**: Removes unused keys quickly across multiple language files if they are not actually referenced in Kotlin code.
- **Interactive Selection**: Provides an interactive UI to select your `res/` directory from multiple Android resource paths in your project.

---

## Installation

### Install from Releases

#### 1. Download
Visit [Polyglot releases page](https://github.com/gustoliveira/polyglot/releases) and download the latest version for your operating system.

#### 2. Extract the downloaded archive

```bash
unzip polyglot_Linux_x86_64.zip
```


#### 3. Move binary

Move the `polyglot` binary to a directory in your `PATH`.

```bash
sudo mv polyglot /usr/local/bin/
```

> [!TIP]
> It is recommended to install it in `/usr/local/bin`, this will ensure that **polyglot** will always be available on your system without interfering with other system programs.
	
#### 4. Verify the installation

```bash
polyglot help
```

### Build from git repository
  #### 1. Clone or Download
  Clone this repository or download the code:

```bash
git clone https://github.com/gustoliveira/polyglot.git
cd polyglot
```

#### 2. Build and Install

Run directly:

```bash
go build -o polyglot main.go
go install
```

#### 3. Create autocompletion

Generate polyglot completion to your specific shell. Is available to `bash`, `fish`, `zsh` and `powershell`

```bash
polyglot completion bash > /tmp/polyglot-completion
source /tmp/polyglot-completion
```

Alternativally you can use Makefile target `build` to install with `bash` autocompletion

```bash
make build  # Builds a binary called 'polyglot', install and source autocompletion
```

---

## Configuration

The CLI needs a Google Translate API key to handle translations. You can supply it in one of two ways:

1. **Environment Variable:** Set `GOOGLE_TRANSLATE_KEY`:
   ```bash
   export GOOGLE_TRANSLATE_KEY="YOUR_API_KEY"
   ```
2. **Command Flag:** Pass `--googleApiKey` to the `translate` command:
   ```bash
   polyglot translate --googleApiKey="YOUR_API_KEY" ...
   ```

---

## Usage

Once installed, run `polyglot <command>` in the root directory of your Android project. Polyglot attempts to detect whether the current directory is an Android project by checking for typical files like `build.gradle`, `settings.gradle`, or an `app/` directory. If these are not found, Polyglot exits with an error.

```bash
polyglot help
```

### Available Commands

#### check
Checks select resource files for:
1. Key sorting: Reports if any file is not sorted.
2. Unused keys: Searches for keys in your `.kt` files. If Polyglot cannot find references like `R.string.<your_key>`, that key is labeled “possibly unused.”
3. Missing translations between files: Report if there're keys that exists in a file and is missing in others.

Run:
```bash
polyglot check
```

#### normalize
Sorts all string keys in `strings.xml` files by alphabetical order across your selected resource directory. If any file is not sorted, Polyglot corrects it in place.

Run:
```bash
polyglot normalize
```

#### remove
Removes a specified key across *all* strings files in a resource directory.

Flags:
- **`--key` or `-k`** *(required)*: The string key to remove.

Usage:
```bash
polyglot remove --key="example_key"
```

#### translate
Translates a single English string (`--value`, `-v`) into every language variant found in your Android `res/` folder (e.g., `values-es`, `values-fr`, etc.). It then appends or substitutes the key in each `strings.xml`.
If the file is sorted, it will be added maintaining the sort property. Otherwise, it will be appended at the end.

Flags:
- **`--key`, `-k`** *(required)*: The key to use for the translated string.
- **`--value`, `-v`** *(required)*: The English text to translate.
- **`--googleApiKey`, `-g`**: Custom Google Translate API key (optional if environment variable is set).
- **`--force`**: Force substitution if a key already exists.
- **`--print-only`**:  Only print the translations instead of adding.

Usage:
```bash
polyglot translate --key="welcome_message" \
                   --value="Welcome to our app!" \
                   --googleApiKey="YOUR_API_KEY" \
                   --force
```

---

## Advanced Topics

### Adding New Subcommands
This CLI is structured using [Cobra](https://cobra.dev/). Each command is defined in a separate file under `cmd/`. For example:
- **`remove.go`** handles the `remove` command.
- **`normalize.go`** handles the `normalize` command.
- etc

To add a new subcommand:
1. Create a new file in `cmd/`.
2. Define a new `*cobra.Command`.
3. Initialize and add it to `rootCmd` in `init()`.

### Android Project Detection
Polyglot checks for any of these in the current directory to confirm you’re in an Android project:
- `build.gradle`
- `settings.gradle`
- `settings.gradle.kts`
- `app/`  
If it doesn’t find at least one, the CLI will exit.

---

## License
Polyglot is [MIT Licensed](./LICENSE). See the [LICENSE](./LICENSE) file for details.
