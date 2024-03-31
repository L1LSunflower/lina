run:
	go run cmd/main.go
build:
	go build -o lina cmd/main.go
build_migrator:
	go build -tags='no_mysql no_sqlite3 no_ydb' -o goose ./pkg/goose/cmd/goose
local_migrate_up:
	./goose -dir ./migrations postgres "postgresql://user:secret@127.0.0.1:5432/app?sslmode=disable" up
local_migrate_down:
	./goose -dir ./migrations postgres "postgresql://user:secret@127.0.0.1:5432/app?sslmode=disable" reset
