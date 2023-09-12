# Deploying AFT with CodeCommit and Terraform OSS

## Deployment

Deploying AFT informing only the required variables:

```sh
aftctl aft deploy \
--region="us-east-1" \ 
--aft-account-id=$AFT_ACCOUNT_ID \ 
--ct-home-region="us-east-1" \ 
--ct-seccondary-region="sa-east-1" \ 
--ct-audit-account-id=$CT_AUDIT_ACCOUNT_ID \ 
--ct-log-archive-account-id=$CT_LOG_ARCHIVE_ACCOUNT_ID \ 
--ct-management-account-id=$CT_MANAGEMENT_ACCOUNT_ID 
```

???+ info
    This documentation is deploying the AFT following the official example found [`here`][AFT Deploy].

[here]: https://github.com/aws-ia/terraform-aws-control_tower_account_factory/blob/main/examples/codecommit%2Btf_oss/main.tf

In cae you want to customize something, this section covers all the available parameters:

Terraform flags:

| flag                             |  type  | use                                                                                        | default value                      |
|----------------------------------|--------|--------------------------------------------------------------------------------------------|------------------------------------|
| --terraform-state-bucket-name    | string | Name of the deployment terraform state bucket (default "aft-deployment-terraform-tfstate") | "aft-deployment-terraform-tfstate" |
| --terraform-version              | string | Terraform version to be used in the deployment and for AFT (default "1.5.6")               | "1.5.6"                            |
| --terraform-distribution         | string | Terraform distribution: oss/tfc                                                            |  oss                               |

Control Tower flags:

| flag                               |  type  | use                                                                 | default value                      |
|------------------------------------|--------|---------------------------------------------------------------------|------------------------------------|
| --ct-management-account-id         | string | Control Tower Management account id (aka payer/root/master account) | ""                                 |
| --ct-log-archive-account-id        | string | Control Tower Log Archive account id                                | ""                                 |
| --ct-audit-account-id              | string | Control Tower Audit account id                                      | ""                                 |
| --ct-home-region                   | string | Control Tower main region                                           | ""                                 |
| --ct-seccondary-region             | string | Control Tower seccondary region                                     | ""                                 |

AFT flags:

| flag                                 |  type  | use                                                                      | default value                      |
|--------------------------------------|--------|--------------------------------------------------------------------------|------------------------------------|
| --aft-account-id                     | string | AFT Management account ID                                                | ""                                 |
| --aft-enable-metrics-reporting       | bool   | Whether to enable reporting metrics or not (default true)                | true                               |
| --aft-enable-cloudtrail-data-events  | bool   | Whether to enable cloudtrail data events (default true)                  | true                               |
| --aft-enable-enterprise-support      | bool   | Whether to enable enterprise support in created accounts (default true)  | true                               |
| --aft-delete-default-vpc             | bool   | Whether to enable enterprise support in created accounts (default true)  | true                               |

Deployment flags:

| flag                              |  type  | use                                                           | default value                                             |
|-----------------------------------|--------|---------------------------------------------------------------|-----------------------------------------------------------|      
| --region                          | string | The region where the aft deployment resources will be created | ""                                                        |
| --branch                          | string | CodeCommit default branch name                                | "main"                                                    |
| --repository-name                 | string | CodeCommit default repository name                            | "aft-deployment"                                          |
| --repository-description          | string | CodeCommit default repository description                     | "CodeCommit repository to store the AFT deployment files" |
| --codepipeline-bucket-name        | string | CodePipeline default artifact bucket                          | "aft-deployment-codepipeline-artifact"                    |
| --docker-image                    | string | CodeBuild default Docker Image name                           | "aws/codebuild/amazonlinux2-x86_64-standard:4.0"          |
| --code-pipeline-role-name         | string | CodePipeline default role name                                | "aft-deployment-codepipeline-service-role"                |
| --code-build-role-name            | string | CodeBuild default role name                                   | "aft-deployment-codebuild-service-role"                   |
| --code-pipeline-role-policy-name  | string | CodePipeline default role policy name                         | "aft-deployment-codepipeline-service-role-policy"         |
| --code-build-role-policy-name     | string | CodeBuild default role policy name                            | "aft-deployment-build-service-role-policy"                |
| --code-build-project-name         | string | CodeBuild default project to deploy AFT                       | "aft-deployment-build"                                    |
| --codepipeline-pipeline-name      | string | CodePipeline default pipeline to deploy AFT                   | "aft-deployment-pipeline"                                 |