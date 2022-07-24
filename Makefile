gen:
	mockgen -source=internal/storage/repository.go -destination=internal/storage/mock/mock_repository.go
	mockgen -source=internal/service/service.go -destination=internal/service/mock/mock_service.go