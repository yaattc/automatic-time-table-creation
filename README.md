# Automatic-Time-Table-Creation [![Go Report Card](https://goreportcard.com/badge/github.com/yaattc/automatic-time-table-creation)](https://goreportcard.com/report/github.com/yaattc/automatic-time-table-creation) [![godoc](https://godoc.org/github.com/yaattc/automatic-time-table-creation?status.svg)](https://godoc.org/github.com/yaattc/automatic-time-table-creation) ![Go](https://github.com/yaattc/automatic-time-table-creation/workflows/Go/badge.svg) [![codecov](https://codecov.io/gh/yaattc/automatic-time-table-creation/branch/master/graph/badge.svg)](https://codecov.io/gh/yaattc/automatic-time-table-creation)

## Build and Deploy

### Environment variables
The application awaits next environment variables provided in .env file in the project folder:

| Environment       | Default  | Description                          | Example                                              |
|-------------------|----------|--------------------------------------|------------------------------------------------------|
| DEBUG             | false    | Turn on debug mode                   | true                                                 |
| POSTGRES_USER     | postgres | Postgres username                    | attc                                                 |
| POSTGRES_PASSWORD |          | Postgres password                    | attcpwd                                              |
| POSTGRES_DB       | postgres | Postgres database name               | attc                                                 |
| DB_CONN_STR       |          | Connection string to database engine | postgres://attc:attcpwd@db:5432/attc?sslmode=disable |
| SERVICE_URL       |          | URL to the backend service           | http://0.0.0.0:8080/                                 |
| SERVICE_PORT      | 8080     | Port of the backend servuce          | 8080                                                 |
| EMAIL             |          | Default admin email                  | e.duskaliev@innopolis.university                     |
| PASSWORD          |          | Default admin password               | test                                                 |

### Run the application
```bash
docker-compose up -d
```

### Testing production build of frontend
To build the frontend prod image:
```bash
docker-compose build frontend
```

To run the frontend prod image on port 8080:
```bash
docker run --rm -e apiUrl=<url to the remote backend> -p 8080:80 semior/attc_frontend:latest
```

## Backend REST API

Several notes:
- All timestamps in RFC3339 format, like `2020-06-30T22:01:53+06:00`.
- All durations in RFC3339 format, like `1h30m5s`.
- Clocks should be represented in ISO 8601 format, like `15:04:05`.

### Errors format

#### Unauthorized
In case if the user requested a route without proper auth, the 401 status code will be returned with the `Unauthorized` body content.

#### General
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
  - `POST /auth/local/login` - authenticate and get JWT token. The token will be saved in secure cookies. 
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
    "user": {
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
    },
    "token": "json.web.token"
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

  - `POST /api/v1/teacher/{id}/preferences` - set teacher preferences
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

  - `POST /api/v1/study_year` - add study year
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
  - `GET /api/v1/study_year` - lists all study years
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
  - `DELETE /api/v1/study_year?id=studyYearID` - remove study year
  - Body: `empty`
  - Response:
  ```json
  {
    "deleted": true
  }
  ```
  - `POST /api/v1/group` - add group
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
  - `GET /api/v1/group` - list groups
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
  - `DELETE /api/v1/group?id=groupID` - remove group
  - Body: `empty`
  - Response:
  ```json
  {
    "deleted": true
  }
  ```

  #### Schedule
  - `GET /api/v1/time_slot` - list time slots
  - Body: `empty`
  - Response:
  ```json
  {
    "time_slots": [
      {
        "id": "965659df-c12a-4fa2-9047-48e8ce5077b7",
        "weekday": 1,
        "start": "09:00:00.000000",
        "duration": "1h30m0s"
      },
      {
        "id": "e7f005f2-4ead-4742-9ad4-a5075baad991",
        "weekday": 2,
        "start": "09:00:00.000000",
        "duration": "1h30m0s"
      }
    ]
  }
  ```

  #### Courses
  - `POST /api/v1/course` - add course
  - Body (**remark**: the assistant lector might be empty - that means that the course does not have tutorials **DO NOT SEND NULL ONLY EMPTY STRING OR UUID**): 
  ```json
  {
    "name": "Operational systems",
    "program": "bachelor",
    "primary_lector": "00000000-0000-0000-0000-200000000001",
    "assistant_lector": "00000000-0000-0000-0000-200000000002",
    "teacher_assistants": [
      "00000000-0000-0000-0000-100000000001",
      "00000000-0000-0000-0000-100000000002",
      "00000000-0000-0000-0000-100000000003",
      "00000000-0000-0000-0000-100000000004",
      "00000000-0000-0000-0000-100000000005"
    ]
  }
  ```
  - Response:
  ```json
  {
    "id": "28b41214-baba-4dbf-a503-ef3e0f33cfae",
    "name": "Operational systems",
    "program": "bachelor",
    "formats": null,
    "groups": null,
    "assistants": [
      {
        "preferences": {
          "time_slots": null,
          "staff": null,
          "rooms": null
        },
        "id": "00000000-0000-0000-0000-100000000001",
        "name": "somename",
        "surname": "somesurname",
        "email": "someemail",
        "degree": "somedegree",
        "about": "something"
      },
      {
        "preferences": {
          "time_slots": null,
          "staff": null,
          "rooms": null
        },
        "id": "00000000-0000-0000-0000-100000000002",
        "name": "somename2",
        "surname": "somesurname2",
        "email": "someemail2",
        "degree": "somedegree2",
        "about": "something2"
      },
      {
        "preferences": {
          "time_slots": null,
          "staff": null,
          "rooms": null
        },
        "id": "00000000-0000-0000-0000-100000000003",
        "name": "somename3",
        "surname": "somesurname3",
        "email": "someemail3",
        "degree": "somedegree3",
        "about": "something3"
      },
      {
        "preferences": {
          "time_slots": null,
          "staff": null,
          "rooms": null
        },
        "id": "00000000-0000-0000-0000-100000000004",
        "name": "somename4",
        "surname": "somesurname4",
        "email": "someemail4",
        "degree": "somedegree4",
        "about": "something4"
      },
      {
        "preferences": {
          "time_slots": null,
          "staff": null,
          "rooms": null
        },
        "id": "00000000-0000-0000-0000-100000000005",
        "name": "somename5",
        "surname": "somesurname5",
        "email": "someemail5",
        "degree": "somedegree5",
        "about": "something5"
      }
    ],
    "primary_lector": {
      "preferences": {
        "time_slots": null,
        "staff": null,
        "rooms": null
      },
      "id": "00000000-0000-0000-0000-200000000001",
      "name": "some primary lector name",
      "surname": "some primary lector surname",
      "email": "some primary lector email",
      "degree": "some primary lector degree",
      "about": "some primary lector about"
    },
    "assistant_lector": {
      "preferences": {
        "time_slots": null,
        "staff": null,
        "rooms": null
      },
      "id": "00000000-0000-0000-0000-200000000002",
      "name": "some assistant lector name",
      "surname": "some assistant lector surname",
      "email": "some assistant lector email",
      "degree": "some assistant lector degree",
      "about": "some assistant lector about"
    },
    "classes": null
  }
  ```
