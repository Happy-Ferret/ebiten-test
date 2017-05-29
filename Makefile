install:
	cd src/sample && glide install -v

update:
	cd src/sample && glide up

run:
	cd src/sample && :; go run ./main.go

.PHONY: install
.PHONY: update
.PHONY: run
