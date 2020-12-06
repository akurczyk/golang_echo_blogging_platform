Golang Echo blogging platform
=============================

Technologies
------------

<img src="https://img.shields.io/static/v1?label=&message=Docker&color=99E7FF&style=flat" height="30"/> <img src="https://img.shields.io/static/v1?label=&message=Docker-compose&color=E6E6E6&style=flat" height="30"/>

<img src="https://img.shields.io/static/v1?label=&message=Golang&color=99FFF5&style=flat" height="30"/> <img src="https://img.shields.io/static/v1?label=&message=Echo%20framework&color=99EEFF&style=flat" height="30"/> <img src="https://img.shields.io/static/v1?label=&message=GORM&color=CCCCCC&style=flat" height="30"/>

<img src="https://img.shields.io/static/v1?label=&message=PostgreSQL&color=99D6FF&style=flat" height="30"/> <img src="https://img.shields.io/static/v1?label=&message=MySQL&color=99D3FF&style=flat" height="30"/> <img src="https://img.shields.io/static/v1?label=&message=Redis&color=FFA099&style=flat" height="30"/> <img src="https://img.shields.io/static/v1?label=&message=MongoDB*&color=A2FF99&style=flat" height="30"/>

<img src="https://img.shields.io/static/v1?label=&message=Prometheus&color=FFAD99&style=flat" height="30"/> <img src="https://img.shields.io/static/v1?label=&message=Swagger&color=9EFF99&style=flat" height="30"/>

<img src="https://img.shields.io/static/v1?label=&message=Kubernetes&color=99BEFF&style=flat" height="30"/> <img src="https://img.shields.io/static/v1?label=&message=Helm&color=99BEFF&style=flat" height="30"/>

Description
-----------

A group of snippet projects made with Golang, Echo framework, and multiple libraries and databases (SQL and NoSQL)
providing an API backend allowing users to register, login, and write posts and comments.

Each of these projects is Dockerized with the multi-stage build and unprivileged user execution in the final image. In
the future addition of Kubernetes Helm charts is planned and another flavors with Cassandra and Scylla is planned.

To run a single project, execute `docker-compose up` in the project directory and go to http://127.0.0.1:1323/. Swagger
with API reference will be available at http://127.0.0.1:1323/swagger/index.html. Prometheus metrics will be accessible
at http://127.0.0.1:1323/metrics.

Authorization
-------------

Most of the endpoints require an authorization with HTTP bearer token with the prefix `Bearer<space>`. To obtain the
token, issue a POST request to `/token` endpoint with form-encoded values of user's `name` and `password` gathered by
creating the user with POST request to `/user` endpoint. An example of these transactions has been shown in the API
reference in 999 directory.

Versions
--------

- **001_postgres_and_redis** – Version with PostgreSQL database used to store users, posts, and comments and Redis used
to store session data. The authentication is done with Bearer Tokens. Swagger and Prometheus are available.

- **002_mysql_and_redis** – Almost the same as described above but with the use of MySQL instead of PostgreSQL.

- **003_MONGO_and_redis** – The same as above but with MongoDB used to store content.

- **999_memory** – Basic version with no real database. The data are stored in RAM during the execution of the program.
This can not be scaled. The API is "documented" in `api_reference.md` as there is no Swagger. Also, there is no
Prometheus included. This is the version that any other one described above is based on. The authentication is done with
Bearer Tokens.

TODO
----

- Write tests
- Add Helm charts
- Add Cassandra and Scylla flavours
