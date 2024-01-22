#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$RETRY_INTERVAL" ]] && fail "RETRY_INTERVAL env variable has not been set"
[[ -z "$TIMEOUT_SECONDS" ]] && fail "TIMEOUT_SECONDS env variable has not been set"
[[ -z "$PACKAGES" ]] && fail "PACKAGES env variable has not been set"
[[ -z "$PACKAGE_MANAGER" ]] && fail "PACKAGE_MANAGER env variable has not been set"

install_packages() {
  if [ "$PACKAGES" = "__skip" ]; then
    return 0
  fi

  # Wait for cloud-init to finish so it doesn't race with any of our package installations.
  # Note: Amazon Linux 2 throws Python 2.7 errors when running `cloud-init status` as
  # non-root user (known bug).
  sudo cloud-init status --wait

  echo "Installing Dependencies: $PACKAGES"

  # Use the default package manager of the current Linux distro to install packages
  if [ "$PACKAGE_MANAGER" = "apt" ]; then
    cd /tmp
    sudo apt update
    # Disable this shellcheck rule about double-quoting array expansions; if we use
    # double quotes on ${PACKAGES[@]}, it does not take the packages as separate
    # arguments.
    # shellcheck disable=SC2068,SC2086
    sudo apt install -y ${PACKAGES[@]}
  elif [ "$PACKAGE_MANAGER" = "yum" ]; then
    cd /tmp
    # shellcheck disable=SC2068,SC2086
    sudo yum -y install ${PACKAGES[@]}
  elif [ "$PACKAGE_MANAGER" = "zypper" ]; then
    # shellcheck disable=SC2068,SC2086
    sudo zypper install --no-confirm ${PACKAGES[@]}
  else
    echo "No matching package manager provided."
    exit 1
  fi

  # if [ -f /etc/debian_version ]; then
  #   # Do our best to make sure that we don't race with cloud-init. Wait a reasonable time until we
  #   # see ec2 in the sources list. Very rarely cloud-init will take longer than we wait. In that case
  #   # we'll just install our packages.
  #   grep ec2 /etc/apt/sources.list || true

  #   cd /tmp
  #   sudo apt update
  #   # shellcheck disable=2068
  #   sudo apt install -y ${PACKAGES[@]}
  # else
  #   cd /tmp
  #   # shellcheck disable=2068
  #   sudo yum -y install ${PACKAGES[@]}
  # fi
}

begin_time=$(date +%s)
end_time=$((begin_time + TIMEOUT_SECONDS))
while [ "$(date +%s)" -lt "$end_time" ]; do
  if install_packages; then
    exit 0
  fi

  sleep "$RETRY_INTERVAL"
done

fail "Timed out waiting for packages to install"
