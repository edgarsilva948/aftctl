# Running Account Factory for Terraform (AFT) code locally

## Example: `terraform init`

???+ note
    You will need to run the aftctl inside the `terraform` directory, e.g.: aft-global-customizations/terraform.

!!! example "Running a terraform init for aft-global-customizations locally"
    
    ```
    aftctl local -a 111111111111 -c init
    ```

Example output:

```
13/09/2023 21:32:30     INFO    📝 Initializing AWS Client using AFT Account credentials... step (1/4)
13/09/2023 21:32:32     INFO    🔄 Successfully set up AWS profile 000000000000-AWSAFTAdmin Step (2/4)
13/09/2023 21:32:33     INFO    ⚙️ .gitignore successfully generated... (3/4)
13/09/2023 21:32:33     INFO    ⛏️  Executing Terraform command init... (4/4)
Output:

Initializing the backend...

Initializing provider plugins...
- Reusing previous version of hashicorp/aws from the dependency lock file
- Using previously-installed hashicorp/aws v5.16.2

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
```

## Example: `terraform plan`

???+ note
    You will need to run the aftctl inside the `terraform` directory, e.g.: aft-global-customizations/terraform.

!!! example "Running a terraform plan for aft-global-customizations locally"
    
    ```
    aftctl local -a 111111111111 -c plan
    ```

Example output:

```
13/09/2023 21:31:56     INFO    📝 Initializing AWS Client using AFT Account credentials... step (1/4)
13/09/2023 21:31:58     INFO    🔄 Successfully set up AWS profile 000000000000-AWSAFTAdmin Step (2/4)
13/09/2023 21:31:58     INFO    ⚙️ .gitignore successfully generated... (3/4)
13/09/2023 21:31:58     INFO    ⛏️  Executing Terraform command plan... (4/4)
Output:
Acquiring state lock. This may take a few moments...
data.aws_region.current: Reading...
data.aws_caller_identity.current: Reading...
data.aws_region.current: Read complete after 0s [id=us-east-1]
data.aws_caller_identity.current: Read complete after 0s [id=111111111111]

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration
and found no differences, so no changes are needed.
Releasing state lock. This may take a few moments...
```


## Example: `terraform apply`

???+ note
    You will need to run the aftctl inside the `terraform` directory, e.g.: aft-global-customizations/terraform.

!!! example "Running a terraform apply for aft-global-customizations locally"
    
    ```
    aftctl local -a 111111111111 -c apply
    ```

Example output:

```
13/09/2023 21:30:05     INFO    📝 Initializing AWS Client using AFT Account credentials... step (1/4)
13/09/2023 21:30:08     INFO    🔄 Successfully set up AWS profile 000000000000-AWSAFTAdmin Step (2/4)
13/09/2023 21:30:08     INFO    ⚙️ .gitignore successfully generated... (3/4)
13/09/2023 21:30:08     INFO    ⛏️  Executing Terraform command apply... (4/4)
Output:
Acquiring state lock. This may take a few moments...
data.aws_region.current: Reading...
data.aws_caller_identity.current: Reading...
data.aws_region.current: Read complete after 0s [id=us-east-1]
data.aws_caller_identity.current: Read complete after 1s [id=111111111111]

No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration
and found no differences, so no changes are needed.
Releasing state lock. This may take a few moments...

Apply complete! Resources: 0 added, 0 changed, 0 destroyed.
```