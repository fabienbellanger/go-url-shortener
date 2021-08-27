# go-url-shortener
A simple URL shortener written in Go with [Fiber](https://github.com/gofiber/fiber)


## Routes

### Web

#### Links

- **[GET] `/:id`**: Redirect to original URL from link ID

### API

#### Users

- **[POST] `/api/v1/login`**: Authentication
    ```bash
    http POST localhost:3000/api/v1/login username=test@gmail.com password=00000000
    ```
    Response:
    ```json
    {
        "id": "2a40080f-6077-4273-9075-1c5503ac95eb",
        "username": "test@gmail.com",
        "lastname": "Test",
        "firstname": "Toto",
        "created_at": "2021-03-08T20:43:28.345Z",
        "updated_at": "2021-03-08T20:43:28.345Z",
        "token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJjcmVhdGVkQXQiOiIyMDIxLTA4LTI3VDE1OjIzOjMyLjIyNloiLCJleHAiOjE2MzA5NDE4MTQsImZpcnN0bmFtZSI6IlRvdG8iLCJpYXQiOjE2MzAwNzc4MTQsImlkIjoiYjc4ODg5M2UtNjFlNi00YTg0LThkZWItMWY4ZjVhYTU2MzFlIiwibGFzdG5hbWUiOiJUZXN0IiwibmJmIjoxNjMwMDc3ODE0LCJ1c2VybmFtZSI6InRlc3RAZ21haWwuY29tIn0.vgccmOprDnbY3m-WurT8CaqKtN7VXT8u4WXtmiIN1w65iD80uGmkDO2iSpkVdG1pgsA3crNvCfgOMwSpHrf3FQ",
        "expires_at": "2021-03-18T21:43:35.641Z"
    }
    ```

- **[POST] `/api/v1/register`**: User creation
    ```bash
    http POST localhost:3000/api/v1/register lastname=Test firstname=Toto username=test@gmail.com password=00000000
    ```
    Response:
    ```json
    {
        "id": "cb13cc29-13bb-4b84-bf30-17da00ec7400",
        "username": "test@gmail.com",
        "lastname": "Test",
        "firstname": "Toto",
        "created_at": "2021-03-09T21:05:35.564747+01:00",
        "updated_at": "2021-03-09T21:05:35.564747+01:00"
    }
    ```

- **[GET] `/api/v1/users`**: Users list
    ```bash
    http GET localhost:3000/api/v1/users "Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJjcmVhdGVkQXQiOiIyMDIxLTA4LTI3VDE1OjIzOjMyLjIyNloiLCJleHAiOjE2MzA5NDE4MTQsImZpcnN0bmFtZSI6IlRvdG8iLCJpYXQiOjE2MzAwNzc4MTQsImlkIjoiYjc4ODg5M2UtNjFlNi00YTg0LThkZWItMWY4ZjVhYTU2MzFlIiwibGFzdG5hbWUiOiJUZXN0IiwibmJmIjoxNjMwMDc3ODE0LCJ1c2VybmFtZSI6InRlc3RAZ21haWwuY29tIn0.vgccmOprDnbY3m-WurT8CaqKtN7VXT8u4WXtmiIN1w65iD80uGmkDO2iSpkVdG1pgsA3crNvCfgOMwSpHrf3FQ"
    ```
    Response:
    ```json
    [
        {
            "id": "2a40080f-6077-4273-9075-1c5503ac95ed",
            "username": "test@gmail.com",
            "lastname": "Test",
            "firstname": "Toto",
            "created_at": "2021-03-08T20:43:28.345Z",
            "updated_at": "2021-03-08T20:43:28.345Z"
        },
        {
            "id": "2a40080f-6077-4273-9075-1c5503ac95eb",
            "username": "test1@gmail.com",
            "lastname": "Test",
            "firstname": "Toto",
            "created_at": "2021-03-08T20:45:51.16Z",
            "updated_at": "2021-03-08T20:45:51.16Z"
        }
    ]
    ```

- **[GET] `/api/v1/users/{id}`**: Get user information
    ```bash
    http GET localhost:3000/api/v1/users/2a40080f-6077-4273-9075-1c5503ac95eb "Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJjcmVhdGVkQXQiOiIyMDIxLTA4LTI3VDE1OjIzOjMyLjIyNloiLCJleHAiOjE2MzA5NDE4MTQsImZpcnN0bmFtZSI6IlRvdG8iLCJpYXQiOjE2MzAwNzc4MTQsImlkIjoiYjc4ODg5M2UtNjFlNi00YTg0LThkZWItMWY4ZjVhYTU2MzFlIiwibGFzdG5hbWUiOiJUZXN0IiwibmJmIjoxNjMwMDc3ODE0LCJ1c2VybmFtZSI6InRlc3RAZ21haWwuY29tIn0.vgccmOprDnbY3m-WurT8CaqKtN7VXT8u4WXtmiIN1w65iD80uGmkDO2iSpkVdG1pgsA3crNvCfgOMwSpHrf3FQ"
    ```
    Response:
    ```json
    {
        "id": "2a40080f-6077-4273-9075-1c5503ac95eb",
        "username": "test@gmail.com",
        "lastname": "Test",
        "firstname": "Toto",
        "created_at": "2021-03-08T20:43:28.345Z",
        "updated_at": "2021-03-08T20:43:28.345Z"
    }
    ```

- **[DELETE] `/api/v1/users/{id}`**: Delete user
    ```bash
    http DELETE localhost:3000/api/v1/users/2a40080f-6077-4273-9075-1c5503ac95eb "Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJjcmVhdGVkQXQiOiIyMDIxLTA4LTI3VDE1OjIzOjMyLjIyNloiLCJleHAiOjE2MzA5NDE4MTQsImZpcnN0bmFtZSI6IlRvdG8iLCJpYXQiOjE2MzAwNzc4MTQsImlkIjoiYjc4ODg5M2UtNjFlNi00YTg0LThkZWItMWY4ZjVhYTU2MzFlIiwibGFzdG5hbWUiOiJUZXN0IiwibmJmIjoxNjMwMDc3ODE0LCJ1c2VybmFtZSI6InRlc3RAZ21haWwuY29tIn0.vgccmOprDnbY3m-WurT8CaqKtN7VXT8u4WXtmiIN1w65iD80uGmkDO2iSpkVdG1pgsA3crNvCfgOMwSpHrf3FQ"
    ```
  Response code `204`

- **[PUT] `/api/v1/users/{id}`**: Update user information
    ```bash
    http PUT localhost:3000/api/v1/users/2a40080f-6077-4273-9075-1c5503ac95eb "Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJjcmVhdGVkQXQiOiIyMDIxLTA4LTI3VDE1OjIzOjMyLjIyNloiLCJleHAiOjE2MzA5NDE4MTQsImZpcnN0bmFtZSI6IlRvdG8iLCJpYXQiOjE2MzAwNzc4MTQsImlkIjoiYjc4ODg5M2UtNjFlNi00YTg0LThkZWItMWY4ZjVhYTU2MzFlIiwibGFzdG5hbWUiOiJUZXN0IiwibmJmIjoxNjMwMDc3ODE0LCJ1c2VybmFtZSI6InRlc3RAZ21haWwuY29tIn0.vgccmOprDnbY3m-WurT8CaqKtN7VXT8u4WXtmiIN1w65iD80uGmkDO2iSpkVdG1pgsA3crNvCfgOMwSpHrf3FQ" lastname=Test firstname=Tutu username=test3@gmail.com password=222222222
    ```
  Response:
    ```json
    {
        "id": "2a40080f-6077-4273-9075-1c5503ac95eb",
        "username": "test3@gmail.com",
        "lastname": "Test",
        "firstname": "Tutu",
        "created_at": "2021-03-08T20:43:28.345Z",
        "updated_at": "2021-03-12T20:43:28.345Z"
    }
    ```
#### Links

- **[POST] `/api/v1/links`**: Create new shortened URL
    ```bash
    http POST localhost:3000/api/v1/links "Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJjcmVhdGVkQXQiOiIyMDIxLTA4LTI3VDE1OjIzOjMyLjIyNloiLCJleHAiOjE2MzA5NDE4MTQsImZpcnN0bmFtZSI6IlRvdG8iLCJpYXQiOjE2MzAwNzc4MTQsImlkIjoiYjc4ODg5M2UtNjFlNi00YTg0LThkZWItMWY4ZjVhYTU2MzFlIiwibGFzdG5hbWUiOiJUZXN0IiwibmJmIjoxNjMwMDc3ODE0LCJ1c2VybmFtZSI6InRlc3RAZ21haWwuY29tIn0.vgccmOprDnbY3m-WurT8CaqKtN7VXT8u4WXtmiIN1w65iD80uGmkDO2iSpkVdG1pgsA3crNvCfgOMwSpHrf3FQ" url=https://google.com expired_at="2021-12-31T00:00:00Z"
    ```
    Response:
    ```json
    {
        "id": "Hr46NkUS",
        "url": "https://google.com",
        "expired_at": "2021-12-31T00:00:00Z"
    }
    ```