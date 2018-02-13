APP?=mdal
RELEASE?=0.0.4
COMMIT?=$(shell git rev-parse --short HEAD)
PROJECT?=github.com/gtfierro/mdal
PERSISTDIR?=/etc/mdal

clean:
	rm -f ${APP}

build: clean
	go build \
		-ldflags "-s -w -X ${PROJECT}/version.Release=${RELEASE} \
						-X ${PROJECT}/version.Commit=${COMMIT}" \
						-o ${APP}
install:
	go install \
		-ldflags "-s -w -X ${PROJECT}/version.Release=${RELEASE} \
						-X ${PROJECT}/version.Commit=${COMMIT}"

run: build
		${APP}

container: build
	docker build -t gtfierro/$(APP):$(RELEASE) .
	docker build -t gtfierro/$(APP):latest .

push: container
	docker push gtfierro/$(APP):$(RELEASE)
	docker push gtfierro/$(APP):latest

containerRun: container
	docker stop $(APP):$(RELEASE) || true && docker rm $(APP):$(RELEASE) || true
	docker run --name $(APP) \
			   --mount type=bind,source=$(shell pwd)/$(PERSISTDIR),target=/etc/mdal \
			   -it \
			   -e BW2_AGENT=$(BW2_AGENT) -e BW2_DEFAULT_ENTITY=$(BW2_DEFAULT_ENTITY) \
			   --rm \
			   gtfierro/$(APP):$(RELEASE)

