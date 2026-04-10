run-gway:
	@go run services/apigateway/cmd/main.go
run-user:
	@go run services/user/cmd/main.go
run-auth:
	@go run services/auth/cmd/main.go


consul:
	@docker run -d -p 8500:8500 -p 8600:8600/udp --name=alp-consul \
	hashicorp/consul agent -server -ui \
	-node=server-1 -bootstrap-expect=1 -client='0.0.0.0'

compose-up:
	@cd services/user && docker compose up -d

run-all:
	@$(MAKE) run-gway & \

	$(MAKE) run-user & \

	$(MAKE) run-auth & \

	wait


