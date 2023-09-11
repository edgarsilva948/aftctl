# aftctl - Facilitates the AFT deployment process.

`aftctl` is a CLI designed to simplify the deployment process of Amazon Account Factory for Terraform (AFT). It follows best practices for Terraform state file isolation, making it easier for you to manage state files securely and efficiently. Additionally, `aftctl` facilitates seamless upgrades and modifications to your AFT configurations in the future. It is written in Go and it welcomes contributions from the community.

!!! example "Deploy aft in the management account"
    
    ```
    aftctl aft deploy \
    --region="us-east-1" \ 
    --aft-account-id=$AFT_ACCOUNT_ID \ 
    --ct-home-region="us-east-1" \ 
    --ct-seccondary-region="sa-east-1" \ 
    --ct-audit-account-id=$CT_AUDIT_ACCOUNT_ID \ 
    --ct-log-archive-account-id=$CT_LOG_ARCHIVE_ACCOUNT_ID \ 
    --ct-management-account-id=$CT_MANAGEMENT_ACCOUNT_ID 
    ```

    ![aftctl deploy](./static/logo.png){ align=right width=25% }

    The deployment will happen in the management account with:

    - the resources will have the default name
    - the resources will reside in the `us-east-1` region
    - the deployment will happen inside the AWS CodePipeline
    - the state will be stored in the S3 Bucket
    - the terraform files are in the AWS CodeCommit

Example output:

```
09/09/2023 23:29:13     INFO    🔒 IAM Role aft-deployment-codepipeline-service-role doesn't exists... creating
09/09/2023 23:29:13     INFO    🔒 IAM Role aft-deployment-codepipeline-service-role successfully created
09/09/2023 23:29:13     INFO    🔒 IAM Role aft-deployment-codebuild-service-role doesn't exists... creating
09/09/2023 23:29:14     INFO    🔒 IAM Role aft-deployment-codebuild-service-role successfully created
09/09/2023 23:29:14     INFO    🪣 S3 bucket ************-aft-deployment-terraform-tfstate doesn't exists... creating
09/09/2023 23:29:15     INFO    🪣 Waiting for bucket "************-aft-deployment-terraform-tfstate" to be created...
09/09/2023 23:29:26     INFO    🪣 S3 Bucket ************-aft-deployment-terraform-tfstate successfully created
09/09/2023 23:29:27     INFO    🪣 S3 bucket ************-aft-deployment-codepipeline-artifact doesn't exists... creating
09/09/2023 23:29:28     INFO    🪣 Waiting for bucket "************-aft-deployment-codepipeline-artifact" to be created...
09/09/2023 23:29:39     INFO    🪣 S3 Bucket ************-aft-deployment-codepipeline-artifact successfully created
09/09/2023 23:29:39     INFO    📁 Directory aft-deployment successfully created
09/09/2023 23:29:39     INFO    📄 File ./aft-deployment/backend.tf successfully created
09/09/2023 23:29:39     INFO    📄 File ./aft-deployment/buildspec.yaml successfully created
09/09/2023 23:29:39     INFO    📄 File ./aft-deployment/main.tf successfully created
09/09/2023 23:29:39     INFO    📦 File ./aft-deployment.zip successfully created
09/09/2023 23:29:39     INFO    ⬆️ zip file aft-deployment.zip successfully uploaded
09/09/2023 23:29:40     INFO    📚 Cloudformation stack aft-deployment-cloudformation-stack doesn't exists... creating
09/09/2023 23:29:40     INFO    🔒 Cloudformation stack aft-deployment-cloudformation-stack successfuly created
09/09/2023 23:29:41     INFO    🛠️ CodeBuild project aft-deployment-build doesn't exists... creating
09/09/2023 23:29:41     INFO    🛠️ CodeBuild Project aft-deployment-build successfully created
09/09/2023 23:29:42     INFO    👷 CodePipeline pipeline aft-deployment-pipeline doesn't exists... creating
09/09/2023 23:29:43     INFO    👷 CodePipeline Pipeline aft-deployment-pipeline successfully created
```