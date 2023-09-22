#!/bin/bash
set -eu
set -x

# get to project root
cd ../../../

dnf --setopt install_weak_deps=False install -y \
  subscription-manager \
  podman python3-pip python3-pytest

python3 -m venv venv
. venv/bin/activate
pip install git+https://github.com/ptoscano/pytest-client-tools@main

pytest -v integration-tests
