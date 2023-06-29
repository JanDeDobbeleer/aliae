# aliae 🌱

Cross shell and platform alias management.

## Installation

### Go

```bash
go install github.com/jandedobbeleer/aliae/src
```

### GitHub Release

Download the binary for your platform/architecture from the latest [release][latest].

## Usage

### zsh

```bash
eval "$(aliae init zsh)"
```

### bash

Add the following to `~/.bashrc` (could be `~/.profile` or `~/.bash_profile` depending on your environment):

```bash
eval "$(aliae init bash)"
```

### pwsh

```powershell
aliae init pwsh | Invoke-Expression
```

### cmd

There's no out-of-the-box support for Windows CMD when it comes to custom prompts.
There is however a way to do it using [Clink][clink], which at the same time supercharges
your cmd experience. Follow the installation instructions and make sure you select autostart.

Integrating aliae with Clink is easy: create a new file called aliae.lua in your Clink
scripts directory (run `clink info` inside cmd to find that file's location).

```lua title="aliae.lua"
load(io.popen('aliae init cmd'):read("*a"))()
```

### fish

Add the following line to `~/.config/fish/config.fish:`:

```fish
aliae init fish | source
```

### nu

> **Warning**
> aliae requires Nushell v0.78.0 or higher.

Add the following line to the Nushell env file (`$nu.env-path`):

```bash
aliae init nu
```

This saves the initialization script to `~/.aliae.nu`.
Now, edit the Nushell config file (`$nu.config-path`) and add the following line at the bottom:

```bash
source ~/.aliae.nu
```

### tcsh

Add the following at the end of `~/.tcshrc`:

```bash
eval `aliae init tcsh`
```

### xonsh

Add the following line at the end of `~/.xonshrc`:

```bash
execx($(aliae init xonsh))
```

## Configuration

By default, aliae expects a `YAML` file called `.aliae.yaml` in the user `HOME` folder.
You can however also specify a custom location using the `--config` flag on initialization.

```bash
eval "$(aliae init zsh --config '/Users/aliae/configs/aliae.yaml')"
```

The custom location can be a local file, or a URL pointing to a remote config.

### Config layout

All shared aliae are defined using the `aliae` property:

```yaml
aliae:
  - alias: a
    value: aliae
  - alias: hello-world
    value: echo "hello world"
    type: function
```

An alias can be a command (default), or a function definition.

[latest]: https://github.com/JanDeDobbeleer/aliae/releases/latest
[clink]: https://chrisant996.github.io/clink/
