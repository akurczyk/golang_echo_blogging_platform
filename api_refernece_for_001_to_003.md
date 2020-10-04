## Users and authentication

** List: ```> http GET http://127.0.0.1:1323/users```**
```
HTTP/1.1 200 OK
Content-Length: 3
Content-Type: application/json; charset=UTF-8
Date: Tue, 01 Sep 2020 21:08:47 GMT

{}
```

** Create: ```> http -f POST http://127.0.0.1:1323/users name='aaa' password='zaq1@WSX' email='aaa@bbb.ccc'```**
```
HTTP/1.1 201 Created
Content-Length: 37
Content-Type: application/json; charset=UTF-8
Date: Tue, 01 Sep 2020 21:08:52 GMT

{
    "email": "aaa@bbb.ccc",
    "name": "aaa"
}
```

**```> http GET http://127.0.0.1:1323/users```**
```
HTTP/1.1 200 OK
Content-Length: 45
Content-Type: application/json; charset=UTF-8
Date: Tue, 01 Sep 2020 21:08:58 GMT

{
    "aaa": {
        "email": "aaa@bbb.ccc",
        "name": "aaa"
    }
}
```

**```> http -f POST http://127.0.0.1:1323/token name='aaa' password='zaq1@WSX'```**
```
HTTP/1.1 200 OK
Content-Length: 32
Content-Type: text/plain; charset=UTF-8
Date: Tue, 01 Sep 2020 21:09:04 GMT

KMNFQOtSny8TOn32pTUuIfUjSVdjdVGk
```

**```> http -f PUT http://127.0.0.1:1323/users password='QWERTY123' email='xxx@yyy.zzz' Authorization:'Bearer KMNFQOtSny8TOn32pTUuIfUjSVdjdVGk'```**
```
HTTP/1.1 200 OK
Content-Length: 37
Content-Type: application/json; charset=UTF-8
Date: Tue, 01 Sep 2020 21:10:00 GMT

{
    "email": "xxx@yyy.zzz",
    "name": "aaa"
}
```

**```> http GET http://127.0.0.1:1323/users```**
```
HTTP/1.1 200 OK
Content-Length: 45
Content-Type: application/json; charset=UTF-8
Date: Tue, 01 Sep 2020 21:10:06 GMT

{
    "aaa": {
        "email": "xxx@yyy.zzz",
        "name": "aaa"
    }
}
```

**```> http -f DELETE http://127.0.0.1:1323/users Authorization:'Bearer KMNFQOtSny8TOn32pTUuIfUjSVdjdVGk'```**
```
HTTP/1.1 204 No Content
Date: Tue, 01 Sep 2020 21:10:14 GMT

<no body>
```

**```> http GET http://127.0.0.1:1323/users```**
```
HTTP/1.1 200 OK
Content-Length: 3
Content-Type: application/json; charset=UTF-8
Date: Tue, 01 Sep 2020 21:10:22 GMT

{}
```

**```> http -f POST http://127.0.0.1:1323/users name='aaa' password='zaq1@WSX' email='aaa@bbb.ccc'```**
```
...
```

**```> http -f POST http://127.0.0.1:1323/token name='aaa' password='zaq1@WSX'```**
```
... AZwm5kg3TC0Saa0HwEat8BWfEBd7Bc0E
```

**```> http DELETE http://127.0.0.1:1323/token Authorization:'Bearer AZwm5kg3TC0Saa0HwEat8BWfEBd7Bc0E'```**
```
HTTP/1.1 204 No Content
Date: Tue, 01 Sep 2020 22:23:51 GMT

<no body>
```

## Posts

**```> http -f POST http://127.0.0.1:1323/token name='aaa' password='zaq1@WSX'```**
```
... LBW9OjFrEwJYilVnpOkCIYyIDkoGqV7B
```

**```> http GET http://127.0.0.1:1323/posts```**
```
HTTP/1.1 200 OK
Content-Length: 3
Content-Type: application/json; charset=UTF-8
Date: Tue, 01 Sep 2020 22:28:07 GMT

{}
```

**```> http -f POST http://127.0.0.1:1323/posts title='aaa' content='bbb' Authorization:'Bearer LBW9OjFrEwJYilVnpOkCIYyIDkoGqV7B'```**
```
HTTP/1.1 201 Created
Content-Length: 149
Content-Type: application/json; charset=UTF-8
Date: Tue, 01 Sep 2020 22:28:12 GMT

{
    "author_name": "aaa",
    "content": "bbb",
    "create_date": "2020-09-01T22:28:12.6177187Z",
    "id": 1,
    "modify_date": "2020-09-01T22:28:12.6177438Z",
    "title": "aaa"
}
```

**```> http GET http://127.0.0.1:1323/posts```**
```
HTTP/1.1 200 OK
Content-Length: 155
Content-Type: application/json; charset=UTF-8
Date: Tue, 01 Sep 2020 22:28:16 GMT

{
    "1": {
        "author_name": "aaa",
        "content": "bbb",
        "create_date": "2020-09-01T22:28:12.6177187Z",
        "id": 1,
        "modify_date": "2020-09-01T22:28:12.6177438Z",
        "title": "aaa"
    }
}
```

**```> http GET http://127.0.0.1:1323/posts\?author_name\=aaa```**
```
HTTP/1.1 200 OK
Content-Length: 155
Content-Type: application/json; charset=UTF-8
Date: Tue, 01 Sep 2020 22:28:21 GMT

{
    "1": {
        "author_name": "aaa",
        "content": "bbb",
        "create_date": "2020-09-01T22:28:12.6177187Z",
        "id": 1,
        "modify_date": "2020-09-01T22:28:12.6177438Z",
        "title": "aaa"
    }
}
```

**```> http -f PUT http://127.0.0.1:1323/posts/1 title='xxx' content='yyy' Authorization:'Bearer LBW9OjFrEwJYilVnpOkCIYyIDkoGqV7B'```**
```
HTTP/1.1 200 OK
Content-Length: 149
Content-Type: application/json; charset=UTF-8
Date: Tue, 01 Sep 2020 22:28:26 GMT

{
    "author_name": "aaa",
    "content": "yyy",
    "create_date": "2020-09-01T22:28:12.6177187Z",
    "id": 1,
    "modify_date": "2020-09-01T22:28:26.6296401Z",
    "title": "xxx"
}
```

**```> http GET http://127.0.0.1:1323/posts```**
```
HTTP/1.1 200 OK
Content-Length: 155
Content-Type: application/json; charset=UTF-8
Date: Tue, 01 Sep 2020 22:28:31 GMT

{
    "1": {
        "author_name": "aaa",
        "content": "yyy",
        "create_date": "2020-09-01T22:28:12.6177187Z",
        "id": 1,
        "modify_date": "2020-09-01T22:28:26.6296401Z",
        "title": "xxx"
    }
}
```

**```> http -f DELETE http://127.0.0.1:1323/posts/1 Authorization:'Bearer LBW9OjFrEwJYilVnpOkCIYyIDkoGqV7B'```**
```
HTTP/1.1 204 No Content
Date: Tue, 01 Sep 2020 22:28:41 GMT

<no body>
```

**```> http GET http://127.0.0.1:1323/posts```**
```
HTTP/1.1 200 OK
Content-Length: 3
Content-Type: application/json; charset=UTF-8
Date: Tue, 01 Sep 2020 22:28:45 GMT

{}
```

## Comments

**```> http GET http://127.0.0.1:1323/comments```**
```
HTTP/1.1 200 OK
Content-Length: 3
Content-Type: application/json; charset=UTF-8
Date: Tue, 01 Sep 2020 23:22:39 GMT

{}
```

**```> http -f POST http://127.0.0.1:1323/comments postid='1' content='111' Authorization:'Bearer LBW9OjFrEwJYilVnpOkCIYyIDkoGqV7B'```**
```
HTTP/1.1 201 Created
Content-Length: 146
Content-Type: application/json; charset=UTF-8
Date: Tue, 01 Sep 2020 23:22:43 GMT

{
    "author_name": "aaa",
    "content": "111",
    "create_date": "2020-09-01T23:22:43.278341Z",
    "id": 1,
    "modify_date": "2020-09-01T23:22:43.2783673Z",
    "post_id": 1
}
```

**```> http GET http://127.0.0.1:1323/comments```**
```
HTTP/1.1 200 OK
Content-Length: 152
Content-Type: application/json; charset=UTF-8
Date: Tue, 01 Sep 2020 23:22:47 GMT

{
    "1": {
        "author_name": "aaa",
        "content": "111",
        "create_date": "2020-09-01T23:22:43.278341Z",
        "id": 1,
        "modify_date": "2020-09-01T23:22:43.2783673Z",
        "post_id": 1
    }
}
```

**```> http GET http://127.0.0.1:1323/comments\?author_name\=aaa\&post_id\=1```**
```
HTTP/1.1 200 OK
Content-Length: 152
Content-Type: application/json; charset=UTF-8
Date: Tue, 01 Sep 2020 23:22:52 GMT

{
    "1": {
        "author_name": "aaa",
        "content": "111",
        "create_date": "2020-09-01T23:22:43.278341Z",
        "id": 1,
        "modify_date": "2020-09-01T23:22:43.2783673Z",
        "post_id": 1
    }
}
```

**```> http -f PUT http://127.0.0.1:1323/comments/1 content='222' Authorization:'Bearer LBW9OjFrEwJYilVnpOkCIYyIDkoGqV7B'```**
```
HTTP/1.1 200 OK
Content-Length: 146
Content-Type: application/json; charset=UTF-8
Date: Tue, 01 Sep 2020 23:23:53 GMT

{
    "author_name": "aaa",
    "content": "222",
    "create_date": "2020-09-01T23:22:43.278341Z",
    "id": 1,
    "modify_date": "2020-09-01T23:23:53.1694098Z",
    "post_id": 1
}
```

**```> http GET http://127.0.0.1:1323/comments```**
```
HTTP/1.1 200 OK
Content-Length: 152
Content-Type: application/json; charset=UTF-8
Date: Tue, 01 Sep 2020 23:23:58 GMT

{
    "1": {
        "author_name": "aaa",
        "content": "222",
        "create_date": "2020-09-01T23:22:43.278341Z",
        "id": 1,
        "modify_date": "2020-09-01T23:23:53.1694098Z",
        "post_id": 1
    }
}
```

**```> http -f DELETE http://127.0.0.1:1323/comments/1 Authorization:'Bearer LBW9OjFrEwJYilVnpOkCIYyIDkoGqV7B'```**
```
HTTP/1.1 204 No Content
Date: Tue, 01 Sep 2020 23:24:03 GMT

<no body>
```

**```> http GET http://127.0.0.1:1323/comments```**
```
HTTP/1.1 200 OK
Content-Length: 3
Content-Type: application/json; charset=UTF-8
Date: Tue, 01 Sep 2020 23:24:08 GMT

{}
```
