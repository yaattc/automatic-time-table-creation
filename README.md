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
	ErrInternal   ErrCode = 0 // any internal error
	ErrDecode     ErrCode = 1 // failed to unmarshal incoming request
	ErrBadRequest ErrCode = 2 // request contains incorrect data or doesn't contain data
)
```

### Client methods

#### Auth methods
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

#### Teachers
- `POST /api/v1/teacher` - add teacher.
  - Body (weekday counts from 0 - Sunday, 6 - Saturday):
    ```json
    {
      "id": "uuid",
      "name": "Ivan",
      "surname": "Konyukhov",
      "email": "i.konyukhov@innopolis.ru",
      "degree": "Dr.",
      "about": "A good professor :)"
    }
    ```
  - Response - updated or added teacher, in case of adding `preferences` might be shrinked: 
    ```json
    {
      "id": "uuid",
      "name": "Ivan",
      "surname": "Konyukhov",
      "email": "i.konyukhov@innopolis.ru",
      "degree": "Dr.",
      "about": "A good professor :)",
      "preferences": {
        "time_slots": [
            {
              "weekday": 1,
              "start": "19:24:00.000000",
              "duration": "1h30m0s",
              "location": "room #108"
            }
        ],
        "staff": [
          {
            "id": "uuid",
            "name": "Nursultan",
            "surname": "Askarbekuly",
            "email": "n.askarbekuly@innopolis.ru",
            "degree": "Mr.",
            "about": "A good TA"
          }
        ],
        "locations": [
          "room #108",
          "room #109",
          "room #231"
        ]
      }
    }
    ```

- `DELETE /api/v1/teacher?id=teacherID` - delete teacher
  - Body: `empty`
  - Response:
    ```json
    {
      "deleted": true
    }
    ```

- `GET /api/v1/teacher?id=teacherID` - list teachers or get a teacher
  - Body: `empty`
  - Response (in case if `teacherID` is not provided, the preferences will be shrinked):
    ```json
    {
      "teachers": [
        {
          "id": "uuid",
          "name": "Ivan",
          "surname": "Konyukhov",
          "email": "i.konyukhov@innopolis.ru",
          "degree": "Dr.",
          "about": "A good professor :)",
          "preferences": {
            "time_slots": [
                {
                  "weekday": 1,
                  "start": "19:24:00.000000",
                  "duration": "1h30m0s",
                  "location": "room #108"
                }
            ],
            "staff": [
              {
                "id": "uuid",
                "name": "Nursultan",
                "surname": "Askarbekuly",
                "email": "n.askarbekuly@innopolis.ru",
                "degree": "Mr.",
                "about": "A good TA"
              }
            ],
            "locations": [
              "room #108",
              "room #109",
              "room #231"
            ]
          }
        }
      ]
    }
    ```

- `POST /teacher/{id}/preferences` - set teacher preferences
  - Body:
    ```json
    {
      "time_slots": [
          {
            "weekday": 1,
            "start": "19:24:00.000000",
            "duration": "1h30m0s",
            "location": "room #108"
          }
      ],
      "staff": [
        {
          "id": "uuid",
          "name": "Nursultan",
          "surname": "Askarbekuly",
          "email": "n.askarbekuly@innopolis.ru",
          "degree": "Mr.",
          "about": "A good TA"
        }
      ],
      "locations": [
        "room #108",
        "room #109",
        "room #231"
      ]
    }
    ```
  - Response:
    ```json
    {
      "id": "uuid",
      "name": "Ivan",
      "surname": "Konyukhov",
      "email": "i.konyukhov@innopolis.ru",
      "degree": "Dr.",
      "about": "A good professor :)",
      "preferences": {
        "time_slots": [
            {
              "weekday": 1,
              "start": "19:24:00.000000",
              "duration": "1h30m0s",
              "location": "room #108"
            }
        ],
        "staff": [
          {
            "id": "uuid",
            "name": "Nursultan",
            "surname": "Askarbekuly",
            "email": "n.askarbekuly@innopolis.ru",
            "degree": "Mr.",
            "about": "A good TA"
          }
        ],
        "locations": [
          "room #108",
          "room #109",
          "room #231"
        ]
      }
    }
    ```

#### Groups and study years

- `POST /study_year` - add study year
  - Body:
    ```json
    {
        "name": "BS - Year 1 (Computer Science)"
    }
    ```
  - Response:
    ```json
    {
        "id": "a2e5feae-1467-4fab-9729-cfefcd71a1c0",
        "name": "BS - Year 1 (Computer Science)"
    }
    ```
- `GET /study_year` - lists all study years
  - Body: `empty`
  - Response:
    ```json
    {
        "study_years": [
            {
                "id": "a2e5feae-1467-4fab-9729-cfefcd71a1c0",
                "name": "BS - Year 1 (Computer Science)"
            },
            {
                "id": "bb42aaf8-b1b2-45d7-b886-e910f356792b",
                "name": "BS - Year 3 (Computer Science)"
            }
        ]
    }
    ```
- `DELETE /study_year?id=studyYearID` - remove study year
  - Body: `empty`
  - Response:
    ```json
    {
        "deleted": true
    }
    ```
- `POST /group` - add group
  - Body:
    ```json
    {
      "name": "B20-02",
      "study_year_id": "bb42aaf8-b1b2-45d7-b886-e910f356792b"
    }
    ```
  - Response:
    ```json
    {
        "id": "c0b709bd-5ad7-4809-832c-446875551ca1",
        "name": "B20-02",
        "study_year": {
            "id": "bb42aaf8-b1b2-45d7-b886-e910f356792b",
            "name": "BS - Year 3 (Computer Science)"
        }
    }
    ```
- `GET /group` - list groups
  - Body: `empty`
  - Response:
    ```json
    {
      "groups": [
          {
              "id": "c0b709bd-5ad7-4809-832c-446875551ca1",
              "name": "B20-02",
              "study_year": {
                  "id": "bb42aaf8-b1b2-45d7-b886-e910f356792b",
                  "name": "BS - Year 3 (Computer Science)"
              }
          },
          {
              "id": "dcd4c252-254b-4b20-a2ff-65c376609eb8",
              "name": "B20-03",
              "study_year": {
                  "id": "bb42aaf8-b1b2-45d7-b886-e910f356792b",
                  "name": "BS - Year 3 (Computer Science)"
              }
          },
          {
              "id": "576c79ae-5611-41e2-8245-a3a7842b9f3e",
              "name": "B20-04",
              "study_year": {
                  "id": "bb42aaf8-b1b2-45d7-b886-e910f356792b",
                  "name": "BS - Year 3 (Computer Science)"
              }
          },
          {
              "id": "ac6fb317-25b2-46e9-9c23-2d1ca526bc64",
              "name": "B20-01",
              "study_year": {
                  "id": "bb42aaf8-b1b2-45d7-b886-e910f356792b",
                  "name": "BS - Year 3 (Computer Science)"
              }
          }
      ]
    }
    ```
- `DELETE /group?id=groupID` - remove group
  - Body: `empty`
  - Response:
    ```json
    {
        "deleted": true
    }
    ```
