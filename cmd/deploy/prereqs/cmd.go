/*
Copyright © 2023 Edgar Costa edgarsilva948@gmail.com
*/

package prereqs

import (
	"strings"

	"github.com/edgarsilva948/aftctl/pkg/aws"
	"github.com/spf13/cobra"
)

var args struct {
	terraformStateBucketName string
	aftManagementAccountID   string
	branchName               string
	gitSourceRepo            string
	codeBuildDockerImage     string
	gitSourceDescription     string
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
		&args.codeBuildDockerImage,
		"docker-image",
		"",
		"aws/codebuild/amazonlinux2-x86_64-standard:3.0",
		"CodeBuild default Docker Image name",
	)

}

func run(cmd *cobra.Command, _ []string) {
	awsClient := aws.NewClient()

	// Trim names to remove any leading/trailing invisible characters
	terraformStateBucketName := strings.Trim(args.terraformStateBucketName, " \t")
	aftManagementAccountID := strings.Trim(args.aftManagementAccountID, " \t")

	// Ensure the bucket is created
	aws.EnsureS3BucketExists(awsClient.GetS3Client(), terraformStateBucketName, aftManagementAccountID, "test-kms-key-id")

	// Ensure the repository is created
	aws.EnsureCodeCommitRepoExists(awsClient.GetCodeCommitClient(), args.gitSourceRepo, args.gitSourceDescription)
}
