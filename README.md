Golang Echo blogging platform
=============================

Technologies
------------

<img src="https://img.shields.io/static/v1?label=&message=Docker&color=99E7FF&style=flat" height="30"/> <img src="https://img.shields.io/static/v1?label=&message=Docker-compose&color=E6E6E6&style=flat" height="30"/>
<img src="https://img.shields.io/static/v1?label=&message=Golang&color=99FFF5&style=flat" height="30"/> <img src="https://img.shields.io/static/v1?label=&message=Echo%20framework&color=99EEFF&style=flat" height="30"/> <img src="https://img.shields.io/static/v1?label=&message=GORM&color=CCCCCC&style=flat" height="30"/>
<img src="https://img.shields.io/static/v1?label=&message=PostgreSQL&color=99D6FF&style=flat" height="30"/> <img src="https://img.shields.io/static/v1?label=&message=MySQL&color=99D3FF&style=flat" height="30"/> <img src="https://img.shields.io/static/v1?label=&message=Redis&color=FFA099&style=flat" height="30"/> <img src="https://img.shields.io/static/v1?label=&message=MongoDB*&color=A2FF99&style=flat" height="30"/> <img src="https://img.shields.io/static/v1?label=&message=Cassandra*&color=99DFFF&style=flat" height="30"/> <img src="https://img.shields.io/static/v1?label=&message=ScyllaDB*&color=99F1FF&style=flat" height="30"/>
<img src="https://img.shields.io/static/v1?label=&message=Prometheus&color=FFAD99&style=flat" height="30"/> <img src="https://img.shields.io/static/v1?label=&message=Swagger&color=9EFF99&style=flat" height="30"/>
<img src="https://img.shields.io/static/v1?label=&message=Kubernetes*&color=99BEFF&style=flat" height="30"/> <img src="https://img.shields.io/static/v1?label=&message=Helm*&color=99BEFF&style=flat" height="30"/>

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

Authorization
-------------

Most of the endpoints requires authorization, for projects from 001 to 003 – with HTTP bearer token with the prefix
`Bearer<space>`, and for projects from 004 to 006 – with JWT token.

To obtain the bearer token, issue a POST request to `/token` endpoint with form-encoded values of user's `name` and
`password` gathered by creating the user with POST request to `/user` endpoint. Example of these transactions has been
shown in the Swagger API reference.

Variants
--------

- **001_postgres_and_redis** – Variant with PostgreSQL database used to store users, posts and comments and Redis used
to store session data. The authentication is done with Bearer Tokens. Swagger and Prometheus are available.

- **002_mysql_and_redis** – Almost the same as described above but with the use od MySQL instead of PostgreSQL.

- **999_memory** – Basic variant with no real database. The data are stored in RAM memory during execution of the
program. Can not be scaled. The API is "docummented" in `api_reference.md` as there is no Swagger. Also there is no
Prometheus included. This is the variant that any other one described above is based on. The authentication is done with
Bearer Tokens.

TODO
----

- Write tests
- Add and describe NoSQL variants
- Add Helm charts
- Add Swagger screenshot to `README.md`
- Add info on how to get the JWT token
