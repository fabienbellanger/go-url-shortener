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
        "token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJjcmVhdGVkQXQiOiIyMDIxLTA5LTIzVDE5OjIxOjE4LjIxN1oiLCJleHAiOjE2MzI1MTEyOTQsImZpcnN0bmFtZSI6IlRvdG8iLCJpYXQiOjE2MzI0MjQ4OTQsImlkIjoiMDBkYWVmODMtMGE5ZC00YWY3LWFhMWYtN2ZlZDMwYzlmZmJlIiwibGFzdG5hbWUiOiJUZXN0IiwibmJmIjoxNjMyNDI0ODk0LCJ1c2VybmFtZSI6InRlc3RAZ21haWwuY29tIn0.XT6Cj5WnH1_h8tvagSE4vcXBVu5_5gox0YqbfasyxRKVGu1hvXNOKOyRTXsrYgigokXHR7pGyAJubEriKKjk4w",
        "expires_at": "2021-03-18T21:43:35.641Z"
    }
    ```

- **[POST] `/api/v1/register`**: User creation
    ```bash
    http POST localhost:3000/api/v1/register lastname=Test firstname=Toto username=test@gmail.com password=00000000 "Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJjcmVhdGVkQXQiOiIyMDIxLTA5LTIzVDE5OjIxOjE4LjIxN1oiLCJleHAiOjE2MzI1MTEyOTQsImZpcnN0bmFtZSI6IlRvdG8iLCJpYXQiOjE2MzI0MjQ4OTQsImlkIjoiMDBkYWVmODMtMGE5ZC00YWY3LWFhMWYtN2ZlZDMwYzlmZmJlIiwibGFzdG5hbWUiOiJUZXN0IiwibmJmIjoxNjMyNDI0ODk0LCJ1c2VybmFtZSI6InRlc3RAZ21haWwuY29tIn0.XT6Cj5WnH1_h8tvagSE4vcXBVu5_5gox0YqbfasyxRKVGu1hvXNOKOyRTXsrYgigokXHR7pGyAJubEriKKjk4w"
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
    http GET localhost:3000/api/v1/users "Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJjcmVhdGVkQXQiOiIyMDIxLTA5LTIzVDE5OjIxOjE4LjIxN1oiLCJleHAiOjE2MzI1MTEyOTQsImZpcnN0bmFtZSI6IlRvdG8iLCJpYXQiOjE2MzI0MjQ4OTQsImlkIjoiMDBkYWVmODMtMGE5ZC00YWY3LWFhMWYtN2ZlZDMwYzlmZmJlIiwibGFzdG5hbWUiOiJUZXN0IiwibmJmIjoxNjMyNDI0ODk0LCJ1c2VybmFtZSI6InRlc3RAZ21haWwuY29tIn0.XT6Cj5WnH1_h8tvagSE4vcXBVu5_5gox0YqbfasyxRKVGu1hvXNOKOyRTXsrYgigokXHR7pGyAJubEriKKjk4w"
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
    http GET localhost:3000/api/v1/users/2a40080f-6077-4273-9075-1c5503ac95eb "Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJjcmVhdGVkQXQiOiIyMDIxLTA5LTIzVDE5OjIxOjE4LjIxN1oiLCJleHAiOjE2MzI1MTEyOTQsImZpcnN0bmFtZSI6IlRvdG8iLCJpYXQiOjE2MzI0MjQ4OTQsImlkIjoiMDBkYWVmODMtMGE5ZC00YWY3LWFhMWYtN2ZlZDMwYzlmZmJlIiwibGFzdG5hbWUiOiJUZXN0IiwibmJmIjoxNjMyNDI0ODk0LCJ1c2VybmFtZSI6InRlc3RAZ21haWwuY29tIn0.XT6Cj5WnH1_h8tvagSE4vcXBVu5_5gox0YqbfasyxRKVGu1hvXNOKOyRTXsrYgigokXHR7pGyAJubEriKKjk4w"
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
    http DELETE localhost:3000/api/v1/users/2a40080f-6077-4273-9075-1c5503ac95eb "Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJjcmVhdGVkQXQiOiIyMDIxLTA5LTIzVDE5OjIxOjE4LjIxN1oiLCJleHAiOjE2MzI1MTEyOTQsImZpcnN0bmFtZSI6IlRvdG8iLCJpYXQiOjE2MzI0MjQ4OTQsImlkIjoiMDBkYWVmODMtMGE5ZC00YWY3LWFhMWYtN2ZlZDMwYzlmZmJlIiwibGFzdG5hbWUiOiJUZXN0IiwibmJmIjoxNjMyNDI0ODk0LCJ1c2VybmFtZSI6InRlc3RAZ21haWwuY29tIn0.XT6Cj5WnH1_h8tvagSE4vcXBVu5_5gox0YqbfasyxRKVGu1hvXNOKOyRTXsrYgigokXHR7pGyAJubEriKKjk4w"
    ```
  Response code `204`

- **[PUT] `/api/v1/users/{id}`**: Update user information
    ```bash
    http PUT localhost:3000/api/v1/users/2a40080f-6077-4273-9075-1c5503ac95eb "Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJjcmVhdGVkQXQiOiIyMDIxLTA5LTIzVDE5OjIxOjE4LjIxN1oiLCJleHAiOjE2MzI1MTEyOTQsImZpcnN0bmFtZSI6IlRvdG8iLCJpYXQiOjE2MzI0MjQ4OTQsImlkIjoiMDBkYWVmODMtMGE5ZC00YWY3LWFhMWYtN2ZlZDMwYzlmZmJlIiwibGFzdG5hbWUiOiJUZXN0IiwibmJmIjoxNjMyNDI0ODk0LCJ1c2VybmFtZSI6InRlc3RAZ21haWwuY29tIn0.XT6Cj5WnH1_h8tvagSE4vcXBVu5_5gox0YqbfasyxRKVGu1hvXNOKOyRTXsrYgigokXHR7pGyAJubEriKKjk4w" lastname=Test firstname=Tutu username=test3@gmail.com password=222222222
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

- **[GET] `/api/v1/links`**: Get all links
    ```bash
    http GET "localhost:3000/api/v1/links?page=1&limit=10" "Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJjcmVhdGVkQXQiOiIyMDIxLTA5LTIzVDE5OjIxOjE4LjIxN1oiLCJleHAiOjE2MzI1MTEyOTQsImZpcnN0bmFtZSI6IlRvdG8iLCJpYXQiOjE2MzI0MjQ4OTQsImlkIjoiMDBkYWVmODMtMGE5ZC00YWY3LWFhMWYtN2ZlZDMwYzlmZmJlIiwibGFzdG5hbWUiOiJUZXN0IiwibmJmIjoxNjMyNDI0ODk0LCJ1c2VybmFtZSI6InRlc3RAZ21haWwuY29tIn0.XT6Cj5WnH1_h8tvagSE4vcXBVu5_5gox0YqbfasyxRKVGu1hvXNOKOyRTXsrYgigokXHR7pGyAJubEriKKjk4w"
    ```
    Response:
    ```json
    [
        {
            "expired_at": "2021-12-31T00:00:00Z",
            "id": "63KzMaYN",
            "url": "https://www.apitic.com"
        },
        {
            "expired_at": "2021-12-05T00:00:00Z",
            "id": "FCpwJsD9",
            "url": "http://localhost:8081/shop/pitaya-larochelle/signIn/+33 6 99 05 85 14&12490922"
        },
        {
            "expired_at": "2021-12-31T00:00:00Z",
            "id": "Hr46NkUS",
            "url": "https://google.com"
        }
    ]

- **[POST] `/api/v1/links`**: Create new shortened URL
    ```bash
    http POST localhost:3000/api/v1/links "Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJjcmVhdGVkQXQiOiIyMDIxLTA5LTIzVDE5OjIxOjE4LjIxN1oiLCJleHAiOjE2MzI1MTEyOTQsImZpcnN0bmFtZSI6IlRvdG8iLCJpYXQiOjE2MzI0MjQ4OTQsImlkIjoiMDBkYWVmODMtMGE5ZC00YWY3LWFhMWYtN2ZlZDMwYzlmZmJlIiwibGFzdG5hbWUiOiJUZXN0IiwibmJmIjoxNjMyNDI0ODk0LCJ1c2VybmFtZSI6InRlc3RAZ21haWwuY29tIn0.XT6Cj5WnH1_h8tvagSE4vcXBVu5_5gox0YqbfasyxRKVGu1hvXNOKOyRTXsrYgigokXHR7pGyAJubEriKKjk4w" url=https://google.com expired_at="2021-12-31T00:00:00Z"
    ```
    Response:
    ```json
    {
        "id": "Hr46NkUS",
        "url": "https://google.com",
        "expired_at": "2021-12-31T00:00:00Z"
    }
    ```