all: requests

install-swagger:
	which swagger || go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger.json: install-swagger
	swagger generate spec -o ./swagger.json --scan-models

requests: swagger.json
	go build cmd/main.go 

clean:
	rm -rf swagger.json 

.PHONY: install-swagger clean all