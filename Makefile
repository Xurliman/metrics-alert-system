migrationsDir = cmd/server/database/migrations
dbstring = 'host=localhost port=5432 user=postgres password=kali dbname=metrics sslmode=disable'
driver = postgres

migrate-up:
	goose -dir ${migrationsDir} ${driver} ${dbstring} up

migrate-down:
	goose -dir ${migrationsDir} ${driver} ${dbstring} down-to 0

