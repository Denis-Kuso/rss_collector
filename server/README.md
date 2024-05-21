# API description

| METHOD | URL                | DESCRIPTION                                                      |
|------- | ------------------ | ---------------------------------------------------------------- |
| POST   | "/users"           | Creates a new user                                               |
| GET    | "/users"           | Auth: return user's data                                         |
| POST   | "/feeds"           | Auth: Creates a new feed for user to follow                      |
| GET    | "/feeds"           | Auth: returns available feeds to follow                          |
| POST   | "/feed_follows"    | Auth: follow an existing feed                                    |
| DELETE | "/feed_follows/id" | Auth: unfollows a currently followed feed                        |
| GET    | "/feed_follows"    | Auth: retrieve all followed feeds                                |
| GET    | "/posts{?limit}"   | Auth: retrieve default num of posts from followed feeds OR limit |


## users
<details>
 <summary><code>POST</code> <code><b>/users</b></code> <code>(create a new user)</code></summary>

> | name              |  type     | data type      | description                         |
> |-------------------|-----------|----------------|-------------------------------------|
> | `name  `          |  required | string         | desired username                    |


##### Response

- HTTP CODE: `200`

- Content:
```json
{
    "name": "Frodo",
    "apiKey": "bXkgcHJlY2lvdXM-aXRzLW1pbmU-bXkgZGVhciBnYW5kYWxm"
}
```

- HTTP CODE: `400`

- Content:

```json
{
    "error": "invalid username"
}
```

##### example cURL

 ```bash
 curl -X POST http://localhost:8080/v1/users -d '{"name": "Frodo"}'
 ```
</details>

<details>
 <summary><code>GET</code> <code><b>/users</b></code> <code>(retrieve user's data)</code></summary>

##### Parameters:

> | name              |  type     | data type      | description                         |
> |-------------------|-----------|----------------|-------------------------------------|
> | `apiKey`          |  required | string         | apiKey used for authentication      |


##### Response

- HTTP CODE: `200`

- Content:

```json
{
    "name": "Smeagol",
    "apiKey": "bXkgcHJlY2lvdXM-aXRzLW1pbmU-bXkgZGVhciBnYW5kYWxm",
    "followedFeeds": [
        {
            "name": "rings blog",
            "url": "https://www.precious.com"
        },
        {
            "name": "fish blog",
            "url": "https://www.juicy-sweet.com"
        }
    ]
}
```

- HTTP CODE: `404`

- Content:

```json
{
    "error": "no such user"
}
```

##### example cURL

```bash
curl 'http://localhost:8080/v1/users' -H 'Authorization: ApiKey bXkgcHJlY2lvdXM-aXRzLW1pbmU-bXkgZGVhciBnYW5kYWxm'
```
</details>

------------------------------------------------------------------------------------------

## feeds

<details>
 <summary><code>POST</code> <code><b>/feeds</b></code> <code>(create and follow new feed)</code></summary>

##### Parameters

> | name              |  type     | data type      | description           |
> |-------------------|-----------|----------------|-----------------------|
> | `feed name`       |  required | int ($int64)   | Desired name for feed |
> | `feed url `       |  required | int ($int64)   | URL of the feed       |
> | `apiKey`          |  required | string         | apiKey used for authentication      |

```json
{
    "name": "AI blog",
    "url": "https://www.aiblog.com"
}
```

##### Response

- HTTP CODE: `200`

- Content:

```json
{
    "name": "AI blog",
    "url": "https://www.aiblog.com",
    "id":"297d4b48-d12f-45f6-bbe5-c5fc673066f4"
}
```

- HTTP CODE: `400`

- Content:

```json
{
    "error":"invalid url format"
}
```

##### example cURL

```bash
curl -X POST http://localhost:8080/v1/feeds -H "Authorization: ApiKey 6711c5359a5bb4a60bfd37113689bc003e128764d2599a7974fbc77e1580c27c" -d '{"name": "AI blog", "url": "www.aiblog.com/xml"}'

```
</details>

<details>
 <summary><code>GET</code> <code><b>/feeds</b></code> <code>(get already available feeds)</code></summary>

##### Parameters

`None`

##### Response

- HTTP CODE:` 200`

- Content:

```json
[
    {
        "name": "Smashing the stack",
        "url": "https://www.set-rsp.com",
        "id": "e62054ae-a46c-40cc-9c1b-1dac7c37581a"
    },
    {
        "name": "Data brokers",
        "url": "https://www.sell-your-data-today.com",
        "id": "cf444a75-8ac7-41cc-a7a7-7154736dbdde" 
    }
]
```

- HTTP CODE: `404`

- Content:

```json
{
    "error":"no feeds found:
}
```

##### example cURL

```bash
curl http://localhost:8080/v1/feeds
```
</details>
------------------------------------------------------------------------------------------

## feed_follows

<details>
 <summary><code>POST</code> <code><b>/feed_follows/{feedID}</b></code> <code>(follow a feed)</code></summary>

##### Parameters

> | name     |  type     | data type      | description                         |
> |----------|-----------|----------------|-------------------------------------|
> | `uuid`   |  required | string         | id of the desired feed to follow    |
> | `apiKey` |  required | string         | apiKey used for authentication      |

##### Response

- HTTP CODE: `200`

- Content:

```json
{
    "name": "AI blog",
    "url": "https://www.aiblog.com",
    "id": "cf444a75-8ac7-41cc-a7a7-7154736dbdde"
}
```

- HTTP CODE: `404`

- Content:

```json
{
    "error":"cannot follow feed"
}
```

##### example cURL

```bash
curl -X POST http://localhost:8080/v1/feed_follows/b410779d-f9b2-436f-a3ef-7e7c31ccf2f5 -H "Authorization: ApiKey 50352ca57f321f95c016a4782751cb155b1340aa8a97fc59aa9c9d5edd96c3d4"
```
</details>

<details>
 <summary><code>GET</code> <code><b>/feed_follows/</b></code> <code>(retrieve all feeds you are following)</code></summary>

##### Parameters

> | name              |  type     | data type      | description                         |
> |-------------------|-----------|----------------|-------------------------------------|
> | `apiKey`          |  required | string         | apiKey used for authentication      |


##### Response

- HTTP CODE: `200`

- Content:

```json
[
    {
        "name": "AI blog",
        "url": "https://www.ai-blog.com",
        "id": "df484a75-8ac7-41cc-a7a7-7154736dbdde"
    },
    {
        "name": "compression blog",
        "url": "https://www.xzblog.com",
        "id": "cf444a75-8ac7-41cc-a7a7-7154736dbdde"
    }
]
```

##### example cURL
```bash
curl http://localhost:8080/v1/feed_follows/ -H "Authorization: ApiKey 50352ca57f321f95c016a4782751cb155b1340aa8a97fc59aa9c9d5edd96c3d4"
```

</details>

<details>
 <summary><code>DELETE</code> <code><b>/feed_follows/{feedID}</b></code> <code>(unfollow feed corresponding to the provided feedID)</code></summary>

##### Parameters
> | name              |  type     | data type      | description                         |
> |-------------------|-----------|----------------|-------------------------------------|
> | `uuid`            |  required | string         | id of the feed                      |
> | `apiKey`          |  required | string         | apiKey used for authentication      |


##### Response

- HTTP CODE:`200`

- Content:

```json
{
    "unfollowedFeed": "cf444a75-8ac7-41cc-a7a7-7154736dbdde"
}
```

- HTTP CODE: `400`

- Content:

```json
{
    "error":"Cannot parse feed id"
}
```

##### example cURL

```bash
curl -X DELETE http://localhost:8080/v1/feed_follows/b410779d-f9b2-436f-a3ef-7e7c31ccf2f5 -H "Authorization: ApiKey 50352ca57f321f95c016a4782751cb155b1340aa8a97fc59aa9c9d5edd96c3d4"
```
</details>

------------------------------------------------------------------------------------------

## posts

<details>
 <summary><code>GET</code> <code><b>/posts/{limit}</b></code> <code>(retrieve posts from followed feeds)</code></summary>

##### Parameters
> | name              |  type     | data type      | description                         |
> |-------------------|-----------|----------------|-------------------------------------|
> | `limit`           |  optional | int ($int64)   | Number of posts to show             |
> | `apiKey`          |  required | string         | apiKey used for authentication      |


##### Response

- HTTP CODE: `200`

- Content:

```json
[
    {
        "feedName": "xkcd",
        "title": "Sphere Tastiness",
        "url": "https://xkcd.com/2893/"
    },
    {
        "feedName": "ai blog",
        "title": "LLM models are not...",
        "url": "https://ai-blog.com/2892/"
    }
]
```

- HTTP CODE: `400`

- Content:

```json
{
    "error":"Provided limit value not supported"
}
```

##### example cURL

```bash
curl http://localhost:8080/v1/posts -H "Authorization: ApiKey 6711c5359a5bb4a60bfd37113689bc003e128764d2599a7974fbc77e1580c27c"
```

</details>
