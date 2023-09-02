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
  -f, --file string                                 This file contains the deploymentinstructions to deploy AFT
      --dry-run                                     Simulate deploying AFT
  -w, --watch                                       Watch cluster installation logs.
      --controltower-management-account-id string   The Management Account ID that will be used during the deployment process.
  -h, --help                                        help for aft
```

### Options inherited from parent commands

```
      --color string     Surround certain characters with escape sequences to display them in color on the terminal. Allowed options are [auto never always] (default "auto")
      --debug            Enable debug mode.
      --profile string   Use a specific AWS profile from your credential file.
      --region string    Use a specific AWS region, overriding the AWS_REGION environment variable.
```

### SEE ALSO

* [aftctl deploy](aftctl_deploy.md)	 - Deploy AFT from from stdin

