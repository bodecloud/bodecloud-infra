.PHONY: all install reset ping prereq

ANSIBLE_CMD = ansible-playbook -i inventory/hosts.yml

# Default target
all: install

# Check prerequisites
prereq:
	ansible-galaxy collection install -r collections/requirements.yml
	ansible --version

# Ping all nodes
ping:
	ansible all -i inventory/hosts.yml -m ping

# Install k3s cluster
install: prereq
	$(ANSIBLE_CMD) playbooks/site.yml

# Reset k3s cluster 
reset:
	$(ANSIBLE_CMD) playbooks/reset.yml

# Apply a specific playbook
apply:
	@if [ -z "$(playbook)" ]; then \
		echo "Error: specify a playbook with 'make apply playbook=<playbook-name>'"; \
		echo "Available playbooks:"; \
		ls -1 playbooks/*.yml; \
		exit 1; \
	fi
	$(ANSIBLE_CMD) playbooks/$(playbook) 