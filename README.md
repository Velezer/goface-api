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
    - Content-Type: application/json
    - Accept: application/json
- Body:
```json 
{
    "username" : "string",
    "password" : "string"
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

