#!/bin/bash

# Deploys the myprecious chart to a custom, non-ElfHosted Kubernetes cluster.
# This script uses Ansible to install FluxCD and point it to a custom
# configuration path within your own fork of the repository.

# --- Configuration ---

# Your GitHub username and the name of your forked repository
# Flux will be configured to monitor this repository.
GITHUB_USER="th3w1zard1"
GITHUB_REPO="my-media-stack"

# The branch in your fork that Flux should monitor.
GIT_BRANCH="main"

# The path within your repository that contains your custom HelmRelease configuration.
# This should match the directory you created in Step 4.
CONFIG_PATH="./k8s/elfhosted/src/infra/tenants/custom-cluster"

# The path to your custom Ansible inventory file created in Step 2.
INVENTORY_FILE="k8s/elfhosted/src/infra/custom-inventory.ini"

# The path to your custom Ansible playbook created in Step 3.
PLAYBOOK="k8s/elfhosted/src/infra/deploy-custom.yml"

# --- Pre-flight Checks ---

if [ "$GITHUB_USER" == "th3w1zard1" ]; then
    echo "ERROR: Please edit this script and set GITHUB_USER to your GitHub username."
    exit 1
fi

if [ ! -f "$INVENTORY_FILE" ]; then
    echo "ERROR: Inventory file not found at $INVENTORY_FILE"
    echo "Please ensure you have created the inventory file as described in the documentation."
    exit 1
fi

if [ ! -f "$PLAYBOOK" ]; then
    echo "ERROR: Playbook not found at $PLAYBOOK"
    exit 1
fi


# --- Deployment ---

echo "Starting deployment of myprecious to your custom cluster..."

# The `flux` role in the ElfHosted repository uses Ansible variables to configure
# the bootstrap command. We will pass these as extra variables on the command line
# to avoid modifying any source files.
# This is derived from how roles like `k8s/elfhosted/src/infra/roles/flux` are designed to be used.
ANSIBLE_EXTRA_VARS="flux_github_repo=$GITHUB_REPO flux_github_user=$GITHUB_USER flux_path=$CONFIG_PATH flux_branch=$GIT_BRANCH"

echo "Running Ansible playbook..."
echo "This will install/update Flux on your cluster and configure it to sync with your repository."

ansible-playbook -i "$INVENTORY_FILE" "$PLAYBOOK" --extra-vars "$ANSIBLE_EXTRA_VARS"

if [ $? -eq 0 ]; then
    echo "Ansible playbook completed successfully."
    echo ""
    echo "Flux has been installed and is configured to watch the path '$CONFIG_PATH' in your repo '$GITHUB_USER/$GITHUB_REPO'."
    echo "You can monitor the deployment status with the following commands:"
    echo "  flux get sources git"
    echo "  flux get kustomizations"
    echo "  flux get helmreleases -A"
    echo ""
    echo "Your applications should start deploying shortly."
else
    echo "ERROR: Ansible playbook failed. Please check the output for errors."
fi 