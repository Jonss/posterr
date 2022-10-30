### Author: João Marcos Santana Santos Júnior

# Posterr

## How to run locally?

### Requisites:

- Docker and docker-compose
  There's three options to run locally:

- Run the database container and app container together, through the command `make run`.
- Run the database container and app container individually, using the commands `make env-up` and `make run-docker`.
- Run the database container individually and the main file itself, using the commands `make env-up` and `make run-local`.

The makefile might help if needed.

### Migrations:

I used [migrate](https://github.com/golang-migrate/migrate) as migration library.
It's possible to run the migration manually if the key `SHOULD_MIGRATE` is set to false.
You can use the command below to execute manually.

```bash
migrate -path db/migration -database 'DATABASE_URL_HERE' -verbose up|down
```

The key `SHOULD_MIGRATE` is set to `true` to run on application startup. I decided to let both options because in a eventual "production" mode, the migrations should run when decided by the team.

### Libraries used in this project:

- [gorilla/mux](https://github.com/gorilla/mux) - I used gorilla/mux as http router. This is the library that usually I use in personal projects, but it take some time configuring.
- [viper](https://github.com/spf13/viper) - Viper is a lib to set configurations. I executed the project calling the main file, the env variables comes from the file `.env`. Viper allows to override the env-var from file when the value comes from the env-var set in the machine or docker.
  You can check the host and port on `DATABASE_URL` is different in the file and on docker-compose.
- [sqlc](https://github.com/kyleconroy/sqlc) - sqlc It's a lib to generate type safe sql code. It helped me to generate the simplest queries, but to fetch the posts conditionally I had to extend the original Querier (interface that handles the queries)
- [mockgen](https://github.com/golang/mock) - Mockgen generates the mocks used in most of the tests. It gave me the possibility to generate mocks in a simple way. I usually create my mocks by myself but for this project it went a great help.
- [testcontainers-go-wrapper](http://github.com/jonss/testcontainers-go-wrapper)/[testcontainers-go](https://github.com/testcontainers/testcontainers-go) - testcontainers-go makes test using a container in a clean way. I created a wrapper for postgres and used in the db layer to test the actual queries.

### How to used.

Fetch posts: `GET /api/posts?start_date=2022-05-02&end_date=2022-05-20&page=0&size=10&only_mine=true`
Create post: `POST /api/posts`
Request:

```{
        "user_id": 1,
        "message": "a message",
        "originalPostId": 2
    }
```

Homepage requires 10 itens in the first request. It can be set using the query parameters `page` and `size`.

User profile [WIP]

## Critique [WIP]

I regret to not add logs when I created the base of the project.
