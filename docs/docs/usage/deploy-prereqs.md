### Credentials

You will need to have AWS API credentials from your Management Account configured.

???+ info
    As described [`here`][AFT Deploy] you will need to have AdministratorAccess to allow AFT Account to launch products from AWS Control Tower Account Factory Portfolio.

[AFT Deploy]: https://docs.aws.amazon.com/controltower/latest/userguide/aft-getting-started.html

You can use [`~/.aws/credentials` file][awsconfig]
or [environment variables][awsenv]. For more information read [AWS documentation](https://docs.aws.amazon.com/cli/latest/userguide/cli-environment.html).

[awsenv]: https://docs.aws.amazon.com/cli/latest/userguide/cli-environment.html
[awsconfig]: https://docs.aws.amazon.com/cli/latest/userguide/cli-config-files.html

### Foundation

An AWS Control Tower landing zone. For more information, see [Plan your AWS Control Tower landing zone.][plan]

A home Region for your AWS Control Tower landing zone. For more information, see [How AWS Regions work with AWS Control Tower][region].

A Terraform version and distribution. For more information, see [Terraform][tfversion] and [AFT versions][aftversion].

A VCS provider for tracking and managing changes to code and other files.

[tfversion]: https://releases.hashicorp.com/terraform/
[aftversion]: https://docs.aws.amazon.com/controltower/latest/userguide/version-supported.html
[region]: https://docs.aws.amazon.com/controltower/latest/userguide/region-how.html
[plan]: https://docs.aws.amazon.com/controltower/latest/userguide/planning-your-deployment.html

### Organizations

Create a new organizational unit for AFT (Optional)

Provision the AFT management account