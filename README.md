# amartha-test
Amarha Backend Engineer Test Using GO


# Shorty Service
## _Build with_

[![N|Solid](https://upload.wikimedia.org/wikipedia/commons/thumb/0/05/Go_Logo_Blue.svg/429px-Go_Logo_Blue.svg.png)](https://golang.org/)


## Endpoint List

- Make Shorten URL
- Go to Shorten URL
- Get Shorten Stats

## Installation

Shorty requires [Docker Compose](https://docs.docker.com/compose/). 

After clone from repository dan already installed docker and docker-compose just follow this command

```sh
docker-compose up --build
```
It will build new docker and will run the service


If already have image just use this

```sh
docker-compose up 
```
Service will run in [localhost:8008][PlDb]

## Detail Endpoint

### POST /shorten

| Attribute | Description |
| ------ | ------ |
| url | url to shorten |
| shortcode | preferential shortcode |

Curl Example :

```sh
curl --location --request POST 'http://localhost:8008/shorten' \
--header 'Content-Type: application/json' \
--data-raw '{
    "url":"http://foo",
    "shortcode":"bar123"
}'
```

### GET /:shortcode/{shortcode}

| Attribute | Description |
| ------ | ------ |
| shortcode | url encoded shortcode |

Just hit

```sh
"http://localhost:8008/shortcode/bar123"
```


### GET /:shortcode/stats/{shortcode}

| Attribute | Description |
| ------ | ------ |
| shortcode | url encoded shortcode |

Curl Example :

```sh
curl --location --request GET 'http://localhost:8008/shortcode/stats/bar123'
```

