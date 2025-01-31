migrationsDir = cmd/server/database/migrations
dbstring = 'host=localhost port=5432 user=postgres password=kali dbname=metrics sslmode=disable'
driver = postgres

migrate-up:
	goose -dir ${migrationsDir} ${driver} ${dbstring} up

migrate-down:
	goose -dir ${migrationsDir} ${driver} ${dbstring} down-to 0

heap:
	curl -v http://localhost:8888/debug/pprof/heap > heap.out
	go tool pprof -http=":9090" -seconds=30 heap.out

goroutine:
	curl -v http://localhost:8888/debug/pprof/goroutine > goroutine.out
	go tool pprof -http=":9090" -seconds=30 goroutine.out

trace:
	curl -v http://localhost:8888/debug/pprof/trace > trace.out
	go tool pprof -http=":9090" -seconds=30 trace.out

profile:
	curl -v http://localhost:8888/debug/pprof/profile > profile.out
	go tool pprof -http=":9090" -seconds=30 profile.out