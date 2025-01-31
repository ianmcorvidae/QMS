all: QMS

install-swagger:
	which swagger || go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger.json: install-swagger
	swagger generate spec -o ./swagger.json --scan-models

QMS: swagger.json
	go build .

clean:
	rm -rf QMS swagger.json 

.PHONY: install-swagger clean all
