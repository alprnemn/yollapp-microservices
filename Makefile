run-gway:
	@go run services/apigateway/cmd/main.go
run-user:
	@go run services/user/cmd/main.go
run-auth:
	@go run services/auth/cmd/main.go