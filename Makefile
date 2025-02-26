IMAGE_NAME = vktex_proj
CONTAINER_NAME = vktex_container
TESTDATA_DIR = internal/pkg/latex2img/testdata

build:
	docker build -t $(IMAGE_NAME) .

generate:
	GOOS=linux go build -o converter internal/pkg/latex2img/testdata_converter/cmd/main.go
	docker run -d --name $(CONTAINER_NAME) vktex_proj tail -f /dev/null
	docker cp ./converter $(CONTAINER_NAME):/converter
	docker cp ./internal/pkg/latex2img/testdata $(CONTAINER_NAME):/testdata
	docker exec $(CONTAINER_NAME) ./converter
	docker cp $(CONTAINER_NAME):testdata ./internal/pkg/latex2img/
	remove_container

remove_container:
	docker stop $(CONTAINER_NAME)
	docker rm $(CONTAINER_NAME)
