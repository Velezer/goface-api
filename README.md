# Deploying

This API deployed on heroku
click the [link](https://goface-api-echo.herokuapp.com/)

Attention!
Due to [go-face](https://github.com/Kagami/go-face) dependency, it's only support jpg image

# API SPEC

## Auth

Some endpoints has auth with jwt-token

Request:

- Header:
  - Authorization: "Bearer your_jwt_token"

## Admin Account

### JWT Register

Request:

- Method: POST
- Endpoint: `/jwt/register`
- Header:
  - Content-Type: application/json
  - Accept: application/json
- Body:

```json
{
    "username" : "string",
    "password" : "string"
}
```

Response :

```json
{
    "detail" : "string"
}
```

### JWT Login

Request:

- Method: POST
- Endpoint: `/jwt/login`
  - Content-Type: application/json
  - Accept: application/json
- Body:

```json
{
    "username" : "string",
    "password" : "string"
}
```

Response :

```json
{
    "token" : "string"
}
```

## Face Data

### Get Face Data

Request:

- Method: GET
- Endpoint: `api/face`
- Header:
  - Accept: application/json
  - Authorization: "Bearer your_jwt_token"

Response:

```json
{
    "data": [
        {
            "id" : "string",
            "name" : "string",
            "descriptors" : "int"
        }
    ]
    
}
```

### Get Face Data By Id

Request:

- Method: GET
- Endpoint: `api/face/:id`
- Header:
  - Accept: application/json
  - Authorization: "Bearer your_jwt_token"

Response:

```json
{
    "data": {
            "id" : "string",
            "name" : "string",
            "descriptors" : "int"
        }
    
}
```

### Delete Face Data

Request:

- Method: DELETE
- Endpoint: `api/face/:id`
- Header:
  - Accept: application/json
  - Authorization: "Bearer your_jwt_token"

Response:

```json
{
    "detail" : "string"
}
```

## Face Action

### Find Face

Request:

- Method: POST
- Endpoint: `api/face/find`
- Header:
  - Content-Type: multipart/form-data
  - Accept: application/json
  - Authorization: "Bearer your_jwt_token"

- Body:

```form-data

"file" : "image/jpeg"

```

Response:

```json
{
    "detail" : "string",
    "response_time" : "string",
    "data" : [
      {
        "id" : "string",
        "name" : "string",
        "distance" : "float64"
      }
    ]
}
```
### Register Face

Request:

- Method: POST
- Endpoint: `api/face/register`
- Header:
  - Content-Type: multipart/form-data
  - Accept: application/json
  - Authorization: "Bearer your_jwt_token"

- Body:

```form-data
"id" : "string",
"name" : "string",
"file" : "image/jpeg"

```

Response:

```json
{
    "detail" : "string",
    "data" : 
      {
        "id" : "string",
        "name" : "string",
        "descriptors" : [
          "[128]float32"
        ]
      }
    
}
### Add Face

Request:

- Method: PUT
- Endpoint: `api/face/register`
- Header:
  - Content-Type: multipart/form-data
  - Accept: application/json
  - Authorization: "Bearer your_jwt_token"

- Body:

```form-data
"id" : "string",
"name" : "string",
"file" : "image/jpeg"

```

Response:

```json
{
    "detail" : "string",
    "data" : 
      {
        "id" : "string",
        "name" : "string",
        "descriptors" : [
          "[128]float32"
        ]
      }
    
}
```
