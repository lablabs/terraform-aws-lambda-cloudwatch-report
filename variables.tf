variable "lambda_schedule" {
  default = "cron(* 1 * * ? *)"
}
variable "region" {}
variable "cw_metric_name" {}
variable "cw_metric_namespace" {}
variable "cw_metric_dimension" {}
variable "name" {}
variable "email_source_address" {}
variable "email_target_address" {}