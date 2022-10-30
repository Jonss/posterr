### Author: João Marcos Santana Santos Júnior

# Posterr

### Requisites:

- Docker and docker-compose

## How to run locally?

There are three options to run locally:

- Run the database container and app container together through the command `make run`.
- Run the database container and app container individually using the commands `make env-up` and `make run-docker`.
- Run the database container individually and the main file itself, using the commands `make env-up` and `make run-local`.

The makefile might help if needed.

### Migrations:

I used [migrate](https://github.com/golang-migrate/migrate) as migration library.
It's possible to run the migration manually if the key `SHOULD_MIGRATE` is set to false.
You can use the command below to execute manually after install `migrate`.

```bash
migrate -path db/migration -database 'DATABASE_URL_HERE.EX:postgres://user:password@host:port/db_name' -verbose up|down
```

The key `SHOULD_MIGRATE` is set to `true` to run on application startup. I decided to let both options because in a eventual "production mode", the migrations should run when decided by the team.

### Libraries used in this project:

- [gorilla/mux](https://github.com/gorilla/mux) - I used gorilla/mux as http router. I usually use this lib in personal projects, but it takes some time to configure.
- [viper](https://github.com/spf13/viper) - Viper is a lib to set configurations. I executed the project by calling the main file. The env variables come from the file `.env`. Viper allows overriding the env-var from the file when the value comes from the env-var set in the machine or docker.
  You can check the host and port on `DATABASE_URL` is different in the file and on docker-compose.
- [sqlc](https://github.com/kyleconroy/sqlc) - sqlc It's a lib to generate type-safe SQL code. It helped me to generate the most uncomplicated queries, but to fetch the posts conditionally, I had to extend the original Querier (an interface that handles the queries)
- [mockgen](https://github.com/golang/mock) - Mockgen generates the mocks used in most of the tests. In addition, it gave me the possibility to generate mocks simply. I usually create my mocks by myself, but this project was a great help.
- [testcontainers-go-wrapper](http://github.com/jonss/testcontainers-go-wrapper)/[testcontainers-go](https://github.com/testcontainers/testcontainers-go) - testcontainers-go makes tests cleanly using a container. I created a wrapper for Postgres and used it in the DB layer to test the actual queries.

### How to use the REST API?

Fetch posts: `GET /api/posts?start_date=2022-05-02&end_date=2022-05-20&page=0&size=10&only_mine=true`
Response:

```json
{
  "content": [
    {
      "id": 9,
      "message": "This is my first post by api",
      "type": "ORIGINAL",
      "username": "fferdinand",
      "originalPost": null
    },
    {
      "id": 8,
      "message": null,
      "type": "REPOSTING",
      "username": "fferdinand",
      "originalPost": {
        "id": 7,
        "message": "maoe 3"
      }
    },
    {
      "id": 7,
      "message": "maoe 3",
      "type": "QUOTE_POST",
      "username": "sissi",
      "originalPost": {
        "id": 2,
        "message": null
      }
    }
  ],
  "hasNext": true,
  "hasPrev": false
}
```

Create post: `POST /api/posts`
Request:

```json
{
  "user_id": 1,
  "message": "a message",
  "originalPostId": 2
}
```

Fetch user info: `GET /api/users/{username}`
Response

```json
{
  "username": "example",
  "dateJoined": "March 25, 2021",
  "postsCount": 42
}
```

Homepage: requires 10 itens in the first request. It can be set using the query parameters `page` and `size`. The other parameters can be used when required.
User profile can use the fetch user info endpoint and fetch posts endpoint. I created the API independent of resources, so `/api/posts` knows post details only and `users/{username}` knows user details only.

## Critique

I should have used some library to help me on assertion, such as testify or is. Unfortunately, doing the assertions using `if` took more time than expected, and the code is not so clean.

I created a simple database with only posts and user tables.
The posts table contains the post itself and a reference for another post, which is used to check if the post is original, repost, and quote-post.

- A original post has a null original_post_id.
- A repost post has an `original_post_id` and a null message.
- A quote-post has an original_post_id and a message.
  I should have created an enum and saved it in the database to make it more straightforward. Instead, I handled every type in the service layer. Thinking about this part, in an eventual refactoring, this part might be confusing , so I'd improve it before.

- If this project were to grow and have many users and posts, which parts do you think would fail first?

- In a real-life situation, what steps would you take to scale this product? What other types of technology and infrastructure might you need to use?
  I'd create a spike to apply cache the posts and users.
  The database could be a bottleneck, so a distributed database such as `Apache Cassandra` or `cochroachDB` could improve the scalability.
  Also, creating posts can be something other than synchronous. For example, an event-based architecture could be applied, using an event broken to consume posts creation.
