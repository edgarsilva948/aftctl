/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

package prereqs

import (
	"strings"

	"github.com/edgarsilva948/aftctl/pkg/aws"
	"github.com/edgarsilva948/aftctl/pkg/initialcommit"
	"github.com/spf13/cobra"
)

var args struct {
	// terraform args
	terraformStateBucketName string
	tfVersion                string
	terraformDistribution    string

	// control tower args
	aftManagementAccountID             string
	ctManagementAccountID              string
	logArchiveAccountID                string
	auditAccountID                     string
	ctHomeRegion                       string
	tfBackendSecondaryRegion           string
	aftMetricsReporting                bool
	aftFeatureCloudtrailDataEvents     bool
	aftFeatureEnterpriseSupport        bool
	aftFeatureDeleteDefaultVPCsEnabled bool

	// deployment resources args
	region                     string
	branchName                 string
	gitSourceRepo              string
	codeBuildDockerImage       string
	gitSourceDescription       string
	codePipelineBucketName     string
	codePipelineRoleName       string
	codePipelineRolePolicyName string
	codeBuildRolePolicyName    string
	codeBuildRoleName          string
	projectName                string
	pipelineName               string
}

// Cmd is the exported command for the AFT prerequisites.
var Cmd = &cobra.Command{
	Use:   "prereqs",
	Short: "Setup AFT prerequisites in AFT-Management Account",
	Long:  "Setup AFT prerequisites in AFT-Management Account",
	Example: `# aftctl usage examples"
	  aftctl deploy prereqs -f deployment.yaml
	
	  aftctl deploy prereqs --region="us-east-1"`,
	Run: run,
}

func init() {
	flags := Cmd.Flags()
	flags.SortFlags = false

	flags.StringVar(
		&args.terraformStateBucketName,
		"terraform-state-bucket-name",
		"",
		"Name of the deployment terraform state bucket",
	)

	flags.StringVar(
		&args.region,
		"region",
		"",
		"The region where the aft deployment resources will be created",
	)

	flags.StringVar(
		&args.aftManagementAccountID,
		"aft-account-id",
		"",
		"AFT Management account ID",
	)

	flags.StringVarP(
		&args.branchName,
		"branch",
		"b",
		"main",
		"CodeCommit default branch name",
	)

	flags.StringVarP(
		&args.gitSourceRepo,
		"repository-name",
		"r",
		"aft-deployment",
		"CodeCommit default repository name",
	)

	flags.StringVarP(
		&args.gitSourceDescription,
		"repository-description",
		"",
		"CodeCommit repository to store the AFT "+
			"deployment files",
		"CodeCommit default repository description",
	)

	flags.StringVarP(
		&args.codePipelineBucketName,
		"codepipeline-bucket-name",
		"",
		"aft-deployment-codepipeline-artifact",
		"CodePipeline default artifact bucket",
	)

	flags.StringVarP(
		&args.codeBuildDockerImage,
		"docker-image",
		"",
		"aws/codebuild/amazonlinux2-x86_64-standard:4.0",
		"CodeBuild default Docker Image name",
	)

	flags.StringVarP(
		&args.codePipelineRoleName,
		"code-pipeline-role-name",
		"",
		"aft-deployment-codepipeline-service-role",
		"CodePipeline default role name",
	)

	flags.StringVarP(
		&args.codeBuildRoleName,
		"code-build-role-name",
		"",
		"aft-deployment-codebuild-service-role",
		"CodeBuild default role name",
	)

	flags.StringVarP(
		&args.codePipelineRolePolicyName,
		"code-pipeline-role-policy-name",
		"",
		"aft-deployment-codepipeline-service-role-policy",
		"CodePipeline default role policy name",
	)

	flags.StringVarP(
		&args.codeBuildRolePolicyName,
		"code-build-role-policy-name",
		"",
		"aft-deployment-build-service-role-policy",
		"CodeBuild default role policy name",
	)

	flags.StringVarP(
		&args.projectName,
		"code-build-project-name",
		"",
		"aft-deployment-build",
		"CodeBuild default project to deploy AFT",
	)

	flags.StringVarP(
		&args.pipelineName,
		"codepipeline-pipeline-name",
		"",
		"aft-deployment-pipeline",
		"CodePipeline default pipeline to deploy AFT",
	)

	flags.StringVarP(
		&args.tfVersion,
		"terraform-version",
		"",
		"1.5.6",
		"Terraform version to be used in the deployment and for AFT",
	)

	flags.StringVarP(
		&args.terraformDistribution,
		"terraform-distribution",
		"",
		"oss",
		"Terraform distribution: oss/tfc",
	)

	flags.StringVar(
		&args.ctManagementAccountID,
		"ct-management-account-id",
		"",
		"CT Management account id (aka payer/root/master account)",
	)

	flags.StringVar(
		&args.logArchiveAccountID,
		"ct-log-archive-account-id",
		"",
		"CT Log Archive account id",
	)

	flags.StringVar(
		&args.auditAccountID,
		"ct-audit-account-id",
		"",
		"CT Audit account id",
	)

	flags.StringVar(
		&args.ctHomeRegion,
		"ct-home-region",
		"",
		"CT main region",
	)

	flags.StringVar(
		&args.tfBackendSecondaryRegion,
		"ct-seccondary-region",
		"",
		"CT seccondary region",
	)

	flags.BoolVarP(
		&args.aftMetricsReporting,
		"aft-enable-metrics-reporting",
		"",
		true,
		"Wheter to enable reporting metrics or not",
	)

	flags.BoolVarP(
		&args.aftFeatureCloudtrailDataEvents,
		"aft-enable-cloudtrail-data-events",
		"",
		true,
		"Wheter to enable cloudtrail data events",
	)

	flags.BoolVarP(
		&args.aftFeatureEnterpriseSupport,
		"aft-enable-enterprise-support",
		"",
		true,
		"Wheter to enable enterprise support in created accounts",
	)

	flags.BoolVarP(
		&args.aftFeatureDeleteDefaultVPCsEnabled,
		"aft-delete-default-vpc",
		"",
		true,
		"Wheter to enable enterprise support in created accounts",
	)

}

func run(cmd *cobra.Command, _ []string) {
	awsClient := aws.NewClient()

	// Trim names to remove any leading/trailing invisible characters
	terraformStateBucketName := strings.Trim(args.terraformStateBucketName, " \t")
	aftManagementAccountID := strings.Trim(args.aftManagementAccountID, " \t")

	interpolatedCodeSuiteBucketName := args.aftManagementAccountID + "-" + args.codePipelineBucketName

	interpolatedZIPFileName := args.gitSourceRepo + ".zip"

	interpolatedCloudformationStackName := args.gitSourceRepo + "-cloudformation-stack"

	codebuildTrustRelationshipService := "codebuild.amazonaws.com"
	codePipelineTrustRelationshipService := "codepipeline.amazonaws.com"

	// Ensure the Code Pipeline Service Role is created
	aws.EnsureIamRoleExists(
		awsClient.GetIamClient(),
		args.codePipelineRoleName,
		codePipelineTrustRelationshipService,
		args.codePipelineRolePolicyName,
		args.region,
		aftManagementAccountID,
		args.gitSourceRepo,
		interpolatedCodeSuiteBucketName,
		terraformStateBucketName,
	)

	// Ensure the Code Build Service Role is created
	aws.EnsureIamRoleExists(
		awsClient.GetIamClient(),
		args.codeBuildRoleName,
		codebuildTrustRelationshipService,
		args.codeBuildRolePolicyName,
		args.region,
		aftManagementAccountID,
		args.gitSourceRepo,
		interpolatedCodeSuiteBucketName,
		terraformStateBucketName,
	)

	// Ensure the tfstate bucket is created
	aws.EnsureS3BucketExists(
		awsClient.GetS3Client(),
		terraformStateBucketName,
		aftManagementAccountID,
		"test-kms-key-id",
		args.codeBuildRoleName,
	)

	// Ensure the codepipeline bucket is created
	aws.EnsureS3BucketExists(
		awsClient.GetS3Client(),
		interpolatedCodeSuiteBucketName,
		aftManagementAccountID,
		"test-kms-key-id",
		args.codeBuildRoleName,
	)

	// Ensure the CodeCommit repo is created with initial code
	initialcommit.GenerateCommitFiles(
		args.gitSourceRepo,
		terraformStateBucketName,
		args.region,
		args.tfVersion,
		args.ctManagementAccountID,
		args.logArchiveAccountID,
		args.auditAccountID,
		args.aftManagementAccountID,
		args.ctHomeRegion,
		args.tfBackendSecondaryRegion,
		args.aftMetricsReporting,
		args.aftFeatureCloudtrailDataEvents,
		args.aftFeatureEnterpriseSupport,
		args.aftFeatureDeleteDefaultVPCsEnabled,
		args.terraformDistribution,
	)

	aws.UploadToS3(
		awsClient.GetS3Client(),
		interpolatedCodeSuiteBucketName,
		interpolatedZIPFileName,
		interpolatedZIPFileName,
	)

	// Ensure the repository is created
	// aws.EnsureCodeCommitRepoExists(
	// 	awsClient.GetCodeCommitClient(),
	// 	args.gitSourceRepo,
	// 	args.gitSourceDescription,
	// )

	// Ensure the repository is created
	aws.EnsureCloudformationExists(
		awsClient.CloudformationClient(),
		interpolatedCloudformationStackName,
		args.gitSourceRepo,
		args.gitSourceDescription,
		interpolatedCodeSuiteBucketName,
		interpolatedZIPFileName,
	)

	// Ensure the Code Build Project is created
	aws.EnsureCodeBuildProjectExists(
		awsClient.CodebuildClient(),
		aftManagementAccountID,
		args.codeBuildDockerImage,
		args.projectName,
		args.gitSourceRepo,
		args.branchName,
		args.codeBuildRoleName,
	)

	// Ensure the Code Pipeline Pipe is created
	aws.EnsureCodePipelineExists(
		awsClient.CodePipelineClient(),
		aftManagementAccountID,
		args.codePipelineRoleName,
		args.pipelineName,
		interpolatedCodeSuiteBucketName,
		args.gitSourceRepo,
		args.branchName,
		args.projectName,
	)

}
