all: setup build run

setup: 
		docker-compose up -d
		#sleep 10 # wait for mysql to start
		mysql -uroot -prushingpass -h 127.0.0.1 < ${CURDIR}/db.sql

build:
		cd ${CURDIR}/cmd/app; go get -d; go clean -r; go build;

run:
		cd ${CURDIR}/cmd/app/; ./app

clean:
		docker-compose down

test:
		mysql -uroot -prushingpass -h 127.0.0.1 < ${CURDIR}/db.sql
		go test -v
