.PHONY: build
	go build C:\Users\NinjaPC\go\src\gophermart\cmd\gophermart\main.go


.PHONY: gen
	gen:
		mockgen -source=C:\Users\NinjaPC\go\src\gophermart\internal\storage\repository.go -destination=C:\Users\NinjaPC\go\src\gophermart\internal\storage\mocks\mock_repository.go


