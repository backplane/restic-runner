restic-runner:
	go build -o "$@"

dist:
	goreleaser build --auto-snapshot --clean

.PHONY: docker-build
docker-build: dist
	docker build -t backplane/restic-runner .

.PHONY: docker-run
docker-run: docker-build
	docker run --rm -it backplane/restic-runner

.PHONY: clean
clean:
	go clean
	rm -rf dist
