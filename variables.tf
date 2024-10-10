variable "lambda_schedule" {
  type    = string
  default = "cron(0 0 * * ? *)"
}

variable "region" {
  type = string
}
variable "cw_metric_name" {
  type        = string
  description = "CloudWatch metric name"
}
variable "cw_metric_namespace" {
  type        = string
  description = "CloudWatch metric namespace"
}
variable "cw_metric_dimension" {
  type        = string
  description = "CloudWatch metric dimension"
}
variable "name" {
  type = string
}
variable "email_source_address" {
  type = string
}
variable "email_target_address" {
  type = string
}
