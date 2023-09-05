# `aftctl` - A CLI for Amazon Account Factory for Terraform (AFT)

[![Go Report Card](https://goreportcard.com/badge/github.com/edgarsilva948/aftctl)](https://goreportcard.com/report/github.com/edgarsilva948/aftctl) [![codecov](https://codecov.io/gh/edgarsilva948/aftctl/graph/badge.svg?token=PIGXFII1NG)](https://codecov.io/gh/edgarsilva948/aftctl) [![CI](https://github.com/edgarsilva948/aftctl/actions/workflows/main.yml/badge.svg)](https://github.com/edgarsilva948/aftctl/actions/workflows/main.yml)


`aftctl` is a straightforward command-line interface (CLI) tool to perform the AFT deploy. 

## Deployment Prerequisites

Before you configure and launch your AFT environment, ensure you have the following prerequisites:

- **An AWS Control Tower landing zone**: For more information, see [Plan your AWS Control Tower landing zone](#).
  
- **A home Region for your AWS Control Tower landing zone**: For more information, see [How AWS Regions work with AWS Control Tower](#).

- **Terraform version and distribution**: For more information, see [Terraform and AFT versions](#).

- **VCS Provider**: A Version Control System (VCS) provider for tracking and managing changes to code and other files.

  > **Note**: By default, AFT uses AWS CodeCommit. For more details, see [What is AWS CodeCommit?](#) in the AWS CodeCommit User Guide.
  >
  > If you'd like to choose a different VCS provider, see [Alternatives for version control of source code in AFT](#).

- **Runtime Environment**: A suitable environment where you can run the Terraform module that installs AFT.

- **AFT Feature Options**: For more information, see [Enable feature options](#).

## Installation

`aftctl` is available for installation as outlined below. 

### For Unix

To download the latest release, run:

```sh
TBD
```

### For Windows

```sh
TBD
```