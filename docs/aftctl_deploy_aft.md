## aftctl deploy aft

Setup AFT in AFT-Management Account

### Synopsis

Setup AFT in AFT-Management Account

```
aftctl deploy aft [flags]
```

### Examples

```
# aftctl usage examples"
	  aftctl deploy aft -f deployment.yaml
	
	  aftctl deploy aft --region="us-east-1"
```

### Options

```
  -w, --watch     Watch logs during deployment.
      --dry-run   Simulate deploying AFT
  -h, --help      help for aft
```

### Options inherited from parent commands

```
      --color string   Surround certain characters with escape sequences to display them in color on the terminal. Allowed options are [auto never always] (default "auto")
```

### SEE ALSO

* [aftctl deploy](aftctl_deploy.md)	 - Deploy AFT from from stdin

