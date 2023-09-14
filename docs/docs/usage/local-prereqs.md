### Credentials

1. **Set the AWS Region for AFT**:  

   Export the AWS region that corresponds to your AFT environment by executing the following command in your terminal:

   ```bash
    export AWS_REGION="us-east-1"
   ```

    Make sure you have valid AWS credentials for accessing the AFT account. These credentials can either be:

    - Stored as environment variables (AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY)
    - Configured in an AWS credentials file (commonly located at ~/.aws/credentials)

2. **Clone the Necessary Repository**:

Obtain a local copy of either the `aft-account-customizations` or `aft-global-customizations` repository by running:

```bash
git clone <REPOSITORY-URL>
```

3. **Install Terraform:**

Confirm that the Terraform CLI is installed and that its binary is accessible from your system's `PATH`. You can verify this by running:

```bash
terraform --version
```
If Terraform is not yet installed, you can follow the [official installation][Terraform Installation guide] to set it up.

[Terraform Installation guide]: https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli

By satisfying these prerequisites, you will be well-prepared to utilize the local command effectively.
