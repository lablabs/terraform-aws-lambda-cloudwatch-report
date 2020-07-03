# AWS Lambda CloudWatch report Terraform module

[![Labyrinth Labs logo](ll-logo.png)](https://www.lablabs.io)

We help companies build, run, deploy and scale software and infrastructure by embracing the right technologies and principles. Check out our website at https://lablabs.io/

---

![Terraform validation](https://github.com/lablabs/terraform-aws-lambda-cloudwatch-report/workflows/Terraform%20validation/badge.svg?branch=master)

## Description

A terraform module to deploy a CloudWatch report lambda function and to provide AWS CloudWatch resources.

<!-- BEGINNING OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
## Requirements

No requirements.

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| cw\_metric\_dimension | n/a | `any` | n/a | yes |
| cw\_metric\_name | n/a | `any` | n/a | yes |
| cw\_metric\_namespace | n/a | `any` | n/a | yes |
| email\_source\_address | n/a | `any` | n/a | yes |
| email\_target\_address | n/a | `any` | n/a | yes |
| name | n/a | `any` | n/a | yes |
| region | n/a | `any` | n/a | yes |
| lambda\_schedule | n/a | `string` | `"cron(0 0 * * ? *)"` | no |

## Outputs

No output.

<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->

## Contributing and reporting issues

Feel free to create an issue in this repository if you have questions, suggestions or feature requests.

## License

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

See [LICENSE](LICENSE) for full details.

    Licensed to the Apache Software Foundation (ASF) under one
    or more contributor license agreements.  See the NOTICE file
    distributed with this work for additional information
    regarding copyright ownership.  The ASF licenses this file
    to you under the Apache License, Version 2.0 (the
    "License"); you may not use this file except in compliance
    with the License.  You may obtain a copy of the License at

      https://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing,
    software distributed under the License is distributed on an
    "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
    KIND, either express or implied.  See the License for the
    specific language governing permissions and limitations
    under the License.
