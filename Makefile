# phonies
.PHONY: build clean cleandb cleanall

clean:
	@docker compose -f test/docker-compose.yaml down
	@docker rmi layeredarch/example:latest

cleandb:
	@rm -fr test/servers/mysql/data
	@mkdir -p test/servers/mysql/data

cleanall: cleandb clean

build:
	@docker compose -f test/docker-compose.yaml up -d
