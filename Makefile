IMAGE_NAME = vktex_proj
CONTAINER_NAME = vktex_container
TESTDATA_DIR = internal/pkg/latex2img/testdata

build:
	docker build -t $(IMAGE_NAME) .

generate:
	GOOS=linux go build -o gen internal/pkg/latex2img/cmd/generate_testdata_images.go
	docker run -d --name $(CONTAINER_NAME) vktex_proj tail -f /dev/null
	docker cp ./gen $(CONTAINER_NAME):/gen
	docker cp ./internal/pkg/latex2img/testdata $(CONTAINER_NAME):/testdata
	docker exec $(CONTAINER_NAME) ./gen
	docker cp $(CONTAINER_NAME):testdata ./internal/pkg/latex2img/
	docker stop $(CONTAINER_NAME)
	docker rm $(CONTAINER_NAME)
