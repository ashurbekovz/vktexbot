IMAGE_NAME = vktex_proj
CONTAINER_NAME = vktex_container
TESTDATA_DIR = internal/pkg/latex2img/testdata

build:
	docker build -t $(IMAGE_NAME) .

run_container:
	@if docker ps -a --format '{{.Names}}' | grep -q "^$(CONTAINER_NAME)$$"; then \
        echo "Container $(CONTAINER_NAME) is already exists."; \
    else \
        docker run -d --name $(CONTAINER_NAME) vktex_proj tail -f /dev/null; \
    fi

remove_container:
	@if docker ps -a --format '{{.Names}}' | grep -q "^$(CONTAINER_NAME)$$"; then \
		if docker ps --format '{{.Names}}' | grep -q "^$(CONTAINER_NAME)$$"; then \
			docker stop $(CONTAINER_NAME); \
		fi; \
		docker rm $(CONTAINER_NAME); \
		echo "Container $(CONTAINER_NAME) stopped and removed."; \
	else \
		echo "Container $(CONTAINER_NAME) does not exist."; \
	fi

generate_testdata:
	make run_container
	docker cp ./ $(CONTAINER_NAME):./
	docker exec $(CONTAINER_NAME) bash -c \
		"cd internal/pkg/latex2img/testdata_converter/cmd/ && \
		go run main.go -path ./../../testdata"
	docker cp $(CONTAINER_NAME):./internal/pkg/latex2img/testdata ./internal/pkg/latex2img/

test:
	make run_container
	docker cp ./ $(CONTAINER_NAME):/sources
	docker exec $(CONTAINER_NAME) bash -c "cd sources && go test ./..."
	
