variable "lambda_schedule" {
  default = "cron(0 0 * * ? *)"
}
variable "region" {}
variable "cw_metric_name" {}
variable "cw_metric_namespace" {}
variable "cw_metric_dimension" {}
variable "name" {}
variable "email_source_address" {}
variable "email_target_address" {}
