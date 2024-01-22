# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
  }
}

locals {
  arch = {
    "amd64" = "x86_64"
    "arm64" = "aarch64"
  }
  package_manager = {
    "amzn"  = "yum"
    "opensuse-leap" = "zypper"
    "rhel"          = "yum"
    "sles"     = "zypper"
    "ubuntu"        = "apt"
  }
  # "packages"        = join(" ", var.packages)
  # "package_manager" = var.package_manager
}

variable "packages" {
  type    = list(string)
  default = []
}

variable "hosts" {
  type = map(object({
    private_ip = string
    public_ip  = string
  }))
  description = "The hosts to install packages on"
}

variable "timeout" {
  type        = number
  description = "The max number of seconds to wait before timing out"
  default     = 120
}

variable "retry_interval" {
  type        = number
  description = "How many seconds to wait between each retry"
  default     = 2
}

resource "enos_host_info" "hosts" {
  for_each = var.hosts

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

resource "enos_remote_exec" "install_packages" {
  for_each = var.hosts

  environment = {
    DISTRO = enos_host_info.hosts[each.key].distro
    PACKAGE_MANAGER = local.package_manager[enos_host_info.hosts[each.key].distro]
    # PACKAGE_MANAGER = "zypper"
    PACKAGES        = length(var.packages) >= 1 ? join(" ", var.packages) : "__skip"
    RETRY_INTERVAL  = var.retry_interval
    TIMEOUT_SECONDS = var.timeout
  }

  scripts = [abspath("${path.module}/scripts/install-packages.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
