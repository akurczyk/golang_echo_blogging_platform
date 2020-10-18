Golang Echo blogging platform
=============================

Technologies
------------
- Docker
- Docker-compose
- Golang
- Echo framework
- GOORM
- PostgreSQL
- MySQL
- Redis
- MongoDB*
- Cassandra*
- ScyllaDB*
- Prometheus
- Swagger
- Kubernetes*
- Helm*

\* – will be added in the future.

Description
-----------

A group of snippet projects made with Golang, Echo framework and multiple libraries and databases (SQL and NoSQL)
providing an API backend allowing user to register, login and write posts and comments.

Each of this projects is Dockerized with multi-stage build and unprivileged user execution in the final image. In the
future adding of Kubernetes Helm charts is planned.

To run a project, execute `docker-compose up` in the project directory and go to http://127.0.0.1:1323/. Swagger with
API reference will be available at http://127.0.0.1:1323/swagger/index.html. Prometheus metrics will be accessible at
http://127.0.0.1:1323/metrics.

Most of the endpoints requires authorization, for projects from 001 to 003 – with HTTP bearer token with the prefix
`Bearer<space>`, and for projects from 004 to 006 – with JWT token. To obtain the bearer token, issue a POST request to
`/token` endpoint with form-encoded values of user's `name` and `password` gathered by creating the user with POST
request to `/user` endpoint. Example of these transactions has been shown in API reference -
`001_memory/api_reference.md`.

Variants
--------

- **001_memory** – Basic variant with no real database. The data are stored in RAM memory during execution of the
program. Can not be scaled. The API is "docummented" in `api_reference.md` as there is no Swagger. Also there is no
Prometheus included. This is the variant that any other one described below is based on. The authentication is done with
Bearer Tokens.

- **002_postgres_and_redis** – Variant with PostgreSQL database used to store users, posts and comments and Redis used
to store session data. The authentication is done with Bearer Tokens. Swagger and Prometheus are available.

- **003_mysql_and_redis** – Almost the same as described above but with the use od MySQL instead of PostgreSQL.

TODO
----

- Test 002
- Add Redis
- Add 003 with MySQL
- Write tests
- Add and describe NoSQL variants
- Add Helm charts
- Add Swagger screenshot to `README.md`
- Add info on how to get the JWT token
