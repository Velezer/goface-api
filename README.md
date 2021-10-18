# API SPEC

## Auth
Some API has auth with jwt-token 

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
    - Content-Type: application/json
    - Accept: application/json

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