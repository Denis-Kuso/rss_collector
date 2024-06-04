# RSS collector

------------------------------------------------------------------------------

## Description

### What

A client-server model that implements an RSS feeder (collects posts from
various sites). The client is implemented as a command-line tool and the server
as a RESTful API. You can host the server somewhere or use the already
hosted server. It can be used with any tool that sends requests (e.g. curl/wget)
but the command-line tool `clipr` is made to query this api/server.

## "Why?"

The benefits of a RSS collector are probably well documented. I wanted a
self-hosted solution, no browser interaction, no reader, a set-and-forget
collector, which I can use with a shell when I want to, query what might be
interesting and only then open a browser to interact with the post.

## 🚀 Quick Start

### Choose suitable binary file

In the `client/releases` folder there are platform-specific binaries. Choose one
that matches your machine, for example, on "apple silicon" machines choose 
`releases/darwin/arm64/clipr`.

### OR install using the Go toolchain

```bash
go install github.com/Denis-Kuso/rss_collector@latest
```

Using this service/api can be done using a tool such as curl/wget, however
the command-line tool `clipr` was built to interact with the api.

## 📖 Usage of CLI tool clipr

`clipr` will look for a configuration file called `.env` inside the `client`
directory. The file must contain a `SERVER_URL=<address_of_server>` key, value
pair that represents the address of your/demo server.

If deployed on your local machine (e.g. port 3000) add the following line:
`SERVER_URL=http://localhost:3000/v1`

If using the server that is already running:
`SERVER_URL=http://142.93.236.158/api/`

The configuration file can also contain a `CRED_LOC` key followed by a valid
pathname (relative to `./client/` or absolute pathname) which represents where
the apikey will be stored. If empty or missing, apikey stored in `.credentials`.

## Examples

> Create a user

```bash
clipr create Frodo
```

> Overwrite user

```bash
clipr create -o Smeagol
```

> Get user info

```bash
clipr info
```

> Get available feeds to follow

```bash
clipr show
```

> Add a feed

```bash
clipr add XKCD https://xkcd.com/rss.xml
```

> Follow an existing feed

```bash
clipr follow c607531a-832a-4b44-9b13-3acd9839d165
```

> Stop following a feed

```bash
clipr rm c607531a-832a-4b44-9b13-3acd9839d165
```

> Get posts from feeds you are following

```bash
clipr fetch
```

All commands can refer to a different config file (other than default .env) and
thus load data for a different user or from a different server.

```bash
clipr info --config someEnvFile
```

```bash
clipr fetch --config .privateUserEnv
```

------------------------------------------------------------------------------

## 🤝 Contributing

### Clone the repo

```bash
git clone https://github.com/Denis-Kuso/rss_collector@latest
cd rss_collector/client
```

### Build the project

```bash
go build -o clipr
```

### Submit a pull request

If you'd like to contribute, please fork the repository and open a pull request
to the `main` branch.
