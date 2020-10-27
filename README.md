# Automatic-Time-Table-Creation [![Go Report Card](https://goreportcard.com/badge/github.com/yaattc/automatic-time-table-creation)](https://goreportcard.com/report/github.com/yaattc/automatic-time-table-creation) [![godoc](https://godoc.org/github.com/yaattc/automatic-time-table-creation?status.svg)](https://godoc.org/github.com/yaattc/automatic-time-table-creation) ![Go](https://github.com/yaattc/automatic-time-table-creation/workflows/Go/badge.svg) [![codecov](https://codecov.io/gh/yaattc/automatic-time-table-creation/branch/master/graph/badge.svg)](https://codecov.io/gh/yaattc/automatic-time-table-creation)

## Backend REST API

Several notes:
- All timestamps in RFC3339 format, like `2020-06-30T22:01:53+06:00`.
- All durations in RFC3339 format, like `1h30m5s`.
- Clocks should be represented in ISO 8601 format, like `15:04:05`.

### Errors format
Example:
```json
{
	"code"     : 0,
	"details"  : "failed to update event",
	"error"    : "event not found"
}
```

In case of bad client request error might have `null` value.

Supported error codes for client mapping:
```go
const (
	ErrInternal ErrCode = 0 // any internal error
	ErrDecode   ErrCode = 1 // failed to unmarshal incoming request
	ErrBadReq   ErrCode = 2 // request contains incorrect data or doesn't contain data
)
```

### Client methods

- `GET /auth/local/login` - authenticate and get JWT token. The token will be saved in secure cookies. 
  - Body:
    ```json
    {
        "user": "e.duskaliev@innopolis.university",
        "passwd": "verystrongpassword"
    }
    ```
  
  - Response (example, shrinked for the sake of simplicity; next examples will be also shrinked): 
    - Headers:
        ```text
        Set-Cookie: JWT=json.web.token; Path=/; Max-Age=720000; HttpOnly
        ```
    - Body (avatar is not used in the application, it is provided by the auth library):
        ```json
        {
          "name": "e.duskaliev@innopolis.university",
          "id": "local_7f48448389aa065af161c3215237acef139e4ecf",
          "picture": "http://0.0.0.0:8080/avatar/", 
          "email": "e.duskaliev@innopolis.university",
          "attrs": {
            "privileges": [
              "read_users",
              "edit_users",
              "list_users",
              "add_users"
            ]
          }
        }
        ```

- `GET /auth/local/logout` - logout from the app, this will remove the JWT token from cookies.
  - Body: `empty`
  - Response:
    - Headers:
      ```text
      Set-Cookie: JWT=; Path=/; Expires=Thu, 01 Jan 1970 00:00:00 GMT; Max-Age=0
      ```
    - Body: `empty`

