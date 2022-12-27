compose_file = test/docker-compose.yml
compose_project = paperless-cli

clean_targets += local-uninstall

.PHONY: local-install
local-install: ## Install paperless-ngx in docker-compose
	docker-compose -f $(compose_file) -p $(compose_project) up -d

.PHONY: local-uninstall
local-uninstall: ## Uninstall paperless-ngx in docker-compose
	docker-compose -f $(compose_file) -p $(compose_project) rm --force --stop -v
