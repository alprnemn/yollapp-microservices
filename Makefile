run-gway:
	@go run services/apigateway/cmd/main.go
run-user:
	@go run services/user/cmd/main.go
run-auth:
	@go run services/auth/cmd/main.go


run-all:
	@$(MAKE) run-gway & \
	$(MAKE) run-user & \
	$(MAKE) run-auth & \
	wait

compose-up:
	@cd services/user && docker compose up -d
