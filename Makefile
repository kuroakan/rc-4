db:
	docker rm -f db
	docker run -v rest:/var/lib/postgresql/data/ --name db -p "5432:5432" --restart=always -e POSTGRES_PASSWORD=dev -e POSTGRES_USER=kuro -e POSTGRES_DB=userdb -d postgres:16.2
migrate-up:
	 goose -dir migrations postgres "user=kr host=localhost port=5432 password=dev dbname=userdb sslmode=disable"  up
migrate-down:
	 goose -dir migrations postgres "user=kr host=localhost port=5432 password=dev dbname=userdb sslmode=disable"  down

 migrate-reset:
	 goose -dir migrations postgres "user=kr host=localhost port=5432 password=dev dbname=userdb sslmode=disable"  reset && \
 	 goose -dir migrations postgres "user=kr host=localhost port=5432 password=dev dbname=userdb sslmode=disable"  up