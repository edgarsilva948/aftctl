## aftctl completion

Generates completion scripts

### Synopsis

To load completions:

Bash:

  $ source <(aftctl completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ aftctl completion bash > /etc/bash_completion.d/aftctl
  # macOS:
  $ aftctl completion bash > /usr/local/etc/bash_completion.d/aftctl

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ aftctl completion zsh > "${fpath[1]}/_aftctl"

  # You will need to start a new shell for this setup to take effect.

fish:

  $ aftctl completion fish | source

  # To load completions for each session, execute once:
  $ aftctl completion fish > ~/.config/fish/completions/aftctl.fish

PowerShell:

  PS> aftctl completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> aftctl completion powershell > aftctl.ps1
  # and source this file from your PowerShell profile.


```
aftctl completion
```

### Options

```
  -h, --help   help for completion
```

### Options inherited from parent commands

```
      --color string   Surround certain characters with escape sequences to display them in color on the terminal. Allowed options are [auto never always] (default "auto")
```

### SEE ALSO

* [aftctl](aftctl.md)	 - Command line tool for Amazon Account Factory for Terraform (AFT).

