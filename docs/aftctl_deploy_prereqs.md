## aftctl deploy prereqs

Setup AFT prerequisites in AFT-Management Account

### Synopsis

Setup AFT prerequisites in AFT-Management Account

```
aftctl deploy prereqs [flags]
```

### Examples

```
# aftctl usage examples"
	  aftctl deploy prereqs -f deployment.yaml
	
	  aftctl deploy prereqs --region="us-east-1"
```

### Options

```
      --terraform-state-bucket-name string      Name of the deployment terraform state bucket
      --aft-account-id string                   AFT Management account ID
  -b, --branch string                           CodeCommit default branch name (default "main")
  -r, --repository-name string                  CodeCommit default repository name (default "aft-deployment")
      --repository-description string           CodeCommit default repository description (default "CodeCommit repository to store the AFT deployment files")
      --codepipeline-bucket-name string         CodePipeline default artifact bucket (default "aft-deployment-codepipeline-artifact")
      --docker-image string                     CodeBuild default Docker Image name (default "aws/codebuild/amazonlinux2-x86_64-standard:4.0")
      --code-pipeline-role-name string          CodePipeline default role name (default "aft-deployment-codepipeline-service-role")
      --code-build-role-name string             CodeBuild default role name (default "aft-deployment-codebuild-service-role")
      --code-pipeline-role-policy-name string   CodePipeline default role policy name (default "aft-deployment-codepipeline-service-role-policy")
      --code-build-role-policy-name string      CodeBuild default role policy name (default "aft-deployment-build-service-role-policy")
      --code-build-project-name string          CodeBuild default project to deploy AFT (default "aft-deployment-build")
      --codepipeline-pipeline-name string       CodePipeline default pipeline to deploy AFT (default "aft-deployment-pipeline")
      --terraform-version string                Terraform version to be used in the deployment and for AFT (default "1.5.6")
      --terraform-distribution string           Terraform distribution: oss/tfc (default "oss")
      --ct-management-account-id string         CT Management account id (aka payer/root/master account)
      --ct-log-archive-account-id string        CT Log Archive account id
      --ct-audit-account-id string              CT Audit account id
      --ct-home-region string                   CT main region
      --ct-seccondary-region string             CT seccondary region
      --aft-enable-metrics-reporting            Wheter to enable reporting metrics or not (default true)
      --aft-enable-cloudtrail-data-events       Wheter to enable cloudtrail data events (default true)
      --aft-enable-enterprise-support           Wheter to enable enterprise support in created accounts (default true)
      --aft-delete-default-vpc                  Wheter to enable enterprise support in created accounts (default true)
  -h, --help                                    help for prereqs
```

### Options inherited from parent commands

```
      --color string   Surround certain characters with escape sequences to display them in color on the terminal. Allowed options are [auto never always] (default "auto")
```

### SEE ALSO

* [aftctl deploy](aftctl_deploy.md)	 - Deploy AFT from from stdin

