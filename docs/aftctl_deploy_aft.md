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
  -w, --watch                                Watch logs during deployment.
      --dry-run                              Simulate deploying AFT
  -f, --file string                          This file contains the deployment instructions to deploy AFT
  -n, --name string                          A metadata Name for the deployment (default "aft-deploy-configuration")
      --create-terraform-state-bucket        Whether to create a terraform state bucket (default true)
      --terraform-state-bucket-name string   Name of the deployment terraform state bucket
      --terraform-state-bucket-path string   Path to save the state file inside the terraform state bucket
  -h, --help                                 help for aft
```

### Options inherited from parent commands

```
      --color string   Surround certain characters with escape sequences to display them in color on the terminal. Allowed options are [auto never always] (default "auto")
```

### SEE ALSO

* [aftctl deploy](aftctl_deploy.md)	 - Deploy AFT from from stdin

