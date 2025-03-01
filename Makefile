ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))

IMAGE_NAME = vktex_proj
CONTAINER_NAME = vktex_container
TESTDATA_DIR = internal/pkg/latex2img/testdata

build:
	docker build -t $(IMAGE_NAME) $(ROOT_DIR)

run_container:
	@if docker ps -a --format '{{.Names}}' | grep -q "^$(CONTAINER_NAME)$$"; then \
        echo "Container $(CONTAINER_NAME) is already exists."; \
    else \
        docker run -d --name $(CONTAINER_NAME) vktex_proj tail -f /dev/null; \
    fi

remove_container:
	@if docker ps -a --format '{{.Names}}' | grep -q "^$(CONTAINER_NAME)$$"; then \
		if docker ps --format '{{.Names}}' | grep -q "^$(CONTAINER_NAME)$$"; then \
			docker kill $(CONTAINER_NAME); \
		fi; \
		docker rm $(CONTAINER_NAME); \
		echo "Container $(CONTAINER_NAME) killed and removed."; \
	else \
		echo "Container $(CONTAINER_NAME) does not exist."; \
	fi

update_app:
	docker exec $(CONTAINER_NAME) bash -c "rm -rf /app/*"
	docker cp $(ROOT_DIR)/. $(CONTAINER_NAME):/app/

generate_testdata:
	make run_container
	make update_app
	\
	docker exec $(CONTAINER_NAME) bash -c \
		"cd /app/internal/pkg/latex2img/testdata_converter/cmd/ && \
		go run main.go -path ./../../testdata"
	docker cp $(CONTAINER_NAME):/app/internal/pkg/latex2img/testdata $(ROOT_DIR)/internal/pkg/latex2img/
	\
	docker exec $(CONTAINER_NAME) bash -c \
		"cd /app/internal/pkg/template2img/testdata_converter/cmd/ && \
		go run main.go -path ./../../testdata"
	docker cp $(CONTAINER_NAME):/app/internal/pkg/template2img/testdata $(ROOT_DIR)/internal/pkg/template2img/

test:
	make run_container
	make update_app
	docker exec $(CONTAINER_NAME) bash -c "cd app && go test ./..."
	
