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
    # Need to refresh repositories first, otherwise searching for them can sometimes fail
    sudo zypper --gpg-auto-import-keys ref
    # shellcheck disable=SC2068,SC2086
    sudo zypper --non-interactive install ${PACKAGES[@]}
  else
    echo "No matching package manager provided."
    exit 1
  fi
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
