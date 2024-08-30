ifeq ($(env), test)
dbname=""
else
dbname=""
env="prod"
endif

dbnetwork="master-db_net"

$(info    ENV is $(env))
$(info    DB is $(dbname))

migrate_make:
	docker run -i --rm -v "$(PWD)/migrations:/migrations" --network ${dbnetwork} migrate/migrate create -ext sql -dir /migrations -seq $(name)
migrate_up:
	docker run -i --rm -v "$(PWD)/migrations:/migrations" --network ${dbnetwork} migrate/migrate -path=/migrations -database "postgres://root:localhost@pgsql:5432/${dbname}?sslmode=disable" -verbose up
migrate_down:
	docker run -i --rm -v "$(PWD)/migrations:/migrations" --network ${dbnetwork} migrate/migrate -path=/migrations -database "postgres://root:localhost@pgsql:5432/${dbname}?sslmode=disable" -verbose down
migrate_drop:
	docker run -i --rm -v "$(PWD)/migrations:/migrations" --network ${dbnetwork} migrate/migrate -path=/migrations -database "postgres://root:localhost@pgsql:5432/${dbname}?sslmode=disable" -verbose drop -f
migrate_seed:
	go run seeds/default_seed.go -env=$(env)
migrate_reset: migrate_drop migrate_up migrate_seed