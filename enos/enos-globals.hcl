# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

globals {
  archs            = ["amd64", "arm64"]
  artifact_sources = ["local", "crt", "artifactory"]
  artifact_types   = ["bundle", "package"]
  backends         = ["consul", "raft"]
  # TO DO: do we need this as a global?
  backend_license_path = var.backend_license_path
  backend_tag_key  = "VaultStorage"

  build_tags = {
    "ce"               = ["ui"]
    "ent"              = ["ui", "enterprise", "ent"]
    "ent.fips1402"     = ["ui", "enterprise", "cgo", "hsm", "fips", "fips_140_2", "ent.fips1402"]
    "ent.hsm"          = ["ui", "enterprise", "cgo", "hsm", "venthsm"]
    "ent.hsm.fips1402" = ["ui", "enterprise", "cgo", "hsm", "fips", "fips_140_2", "ent.hsm.fips1402"]
  }
  consul_editions  = ["ce", "ent"]
  # TO DO: rename consul_version?
  consul_versions = ["1.14.11", "1.15.7", "1.16.3", "1.17.0"]
  distros         = ["amazon_linux", "leap", "rhel", "sles", "ubuntu"]
  distro_packages = {
    amazon_linux = ["nc"]
    leap         = ["netcat"]
    rhel         = ["nc"]
    sles         = ["netcat"]
    ubuntu       = ["netcat"]
  }
  distro_version = {
    "amazon_linux" = var.distro_version_amazon_linux
    "leap"         = var.distro_version_leap
    "rhel"         = var.distro_version_rhel
    "sles"         = var.distro_version_sles
    "ubuntu"       = var.distro_version_ubuntu
  }
  package_manager = {
    "amazon_linux" = "yum"
    "leap"         = "zypper"
    "rhel"         = "yum"
    "sles"         = "zypper"
    "ubuntu"       = "apt"
  }
  editions = ["ce", "ent", "ent.fips1402", "ent.hsm", "ent.hsm.fips1402"]
  packages = ["jq"]
  sample_attributes = {
    aws_region = ["us-east-1", "us-west-2"]
    distro_version_amazon_linux = ["amzn2"]
    distro_version_leap         = ["15.4", "15.5"]
    distro_version_rhel         = ["8.8", "9.1"]
    distro_version_sles         = ["v15_sp4_standard", "v15_sp5_standard"]
    distro_version_ubuntu       = ["18.04", "20.04", "22.04"]
  }
  seals = ["awskms", "pkcs11", "shamir"]
  tags = merge({
    "Project Name" : var.project_name
    "Project" : "Enos",
    "Environment" : "ci"
  }, var.tags)
  // NOTE: when backporting, make sure that our initial versions are less than that
  // release branch's version. Also beware if adding versions below 1.11.x. Some scenarios
  // that use this global might not work as expected with earlier versions. Below 1.8.x is
  // not supported in any way.
  upgrade_initial_versions = ["1.11.12", "1.12.11", "1.13.11", "1.14.7", "1.15.3"]
  # TO DO: make sure all of these vars are transitioned to vault_install_dir
  // vault_install_dir_packages = {
  //   rhel   = "/bin"
  //   ubuntu = "/usr/bin"
  // }
  vault_install_dir = {
    bundle  = "/opt/vault/bin"
    package = "/usr/bin"
  }
  vault_license_path = abspath(var.vault_license_path != null ? var.vault_license_path : joinpath(path.root, "./support/vault.hclic"))
  vault_tag_key      = "Type" // enos_vault_start expects Type as the tag key
}
