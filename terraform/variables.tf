variable "token" {
  description = "Your Linode API Personal Access Token. (required)"
  default = "95aa2e0d7d05fe4c97238d8b9ac83c230ba009f7a988f7d2ad139db877c4d432"
}

variable "k8s_version" {
  description = "The Kubernetes version to use for this cluster. (required)"
  default = "1.27"
}

variable "label" {
  description = "The unique label to assign to this cluster. (required)"
  default = "big-data-final-project
}

variable "region" {
  description = "The region where your cluster will be located. (required)"
  default = "us-central"
}

variable "tags" {
  description = "Tags to apply to your cluster for organizational purposes. (optional)"
  type = list(string)
  default = ["testing"]
}

variable "pools" {
  description = "The Node Pool specifications for the Kubernetes cluster. (required)"
  type = list(object({
    type = string
    count = number
  }))
  default = [
    {
      type = "g6-standard-4"
      count = 3
    },
    {
      type = "g6-standard-8"
      count = 3
    }
  ]
}