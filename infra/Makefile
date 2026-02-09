##@ Utility
help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


##@ Setup
auth: ## Auth into google drive
	gcloud auth application-default login


##@ Tofu
_tofu_run_cmd:
	set -a && . ../.env && set +a && tofu $(CMD)

tofu-init: ## Initialize OpenTofu
	$(MAKE) _tofu_run_cmd CMD="init"

tofu-plan: ## Plan OpenTofu
	$(MAKE) _tofu_run_cmd CMD="plan"

tofu-apply: ## Plan OpenTofu
	$(MAKE) _tofu_run_cmd CMD="apply"

tofu-destroy: ## Destroy OpenTofu
	$(MAKE) _tofu_run_cmd CMD="destroy"