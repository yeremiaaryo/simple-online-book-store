# gotu-assignment

## How to Run Services
1. docker-compose up # this will spin up postgres instance in your local
2. install golang migrate: https://github.com/golang-migrate/migrate
3. make migrate-up # this will initiate all the tables needed, as well as the books table that is preloaded with 10 data
4. make run # this will run user service in port 9999
5. Postman collection is included for testing purposes (`Gotu.postman_collection.json`), you can import to your postman apps

## APIs
### Users Service
##### Register
API to register a new users by sending email and password

```
URL: POST /register
Content-Type: application/json
```
##### Request body: (JSON body)
```json
{
    "email": "email@gmail.com",
    "password": "password"
}
```
##### Response:
```json
{
    "result": true,
    "user": {
        "id": 1,
        "email": "email@gmail.com",
        "created_at": 1718290645179,
        "updated_at": 1718290645179
    }
}
```

##### Login
API to log in a user to the system by sending email and password. It will return the respective JWT Token that must be sent on the order API

```
URL: POST /login
Content-Type: application/json
```
##### Request body: (JSON body)
```json
{
    "email": "email@gmail.com",
    "password": "password"
}
```
##### Response:
```json
{
    "result": true,
    "token": "JWT Token"
}
```

### Books Service
##### Book List
API to get book list, this API doesn't need token since an online book store won't need the user to create account just to search books

```
URL: GET /books
Content-Type: application/json
```
```
Parameters:
page_index = int // default will be 1
page_size = int // default will be 10
search = string // can be used to search by title or by author
```
##### Response:
```json
{
    "result": true,
    "books": [
        {
            "id": 1,
            "title": "The Catcher in the Rye",
            "author": "J.D. Salinger",
            "isbn": "9780316769488",
            "published_date": "1951-07-16T00:00:00Z",
            "price": 10.99
        },
        {
            "id": 2,
            "title": "To Kill a Mockingbird",
            "author": "Harper Lee",
            "isbn": "9780061120084",
            "published_date": "1960-07-11T00:00:00Z",
            "price": 7.99
        }
    ]
}
```


### Orders Service
##### Create Order
API to create order, need Bearer token got from the login API to be included in header

```
URL: POST /order
Content-Type: application/json
```
##### Request body: (JSON body)
```json
{
    "total_amount": 35.96,
    "items": [
        {
            "book_id": 10,
            "quantity": 2,
            "price": 9.99
        },
        {
            "book_id": 2,
            "quantity": 2,
            "price": 7.99
        }
    ]
}
```
##### Response:
```json
{
    "result": true,
    "order_id": 3,
    "status": "NEW"
}
```

##### Order History
API to get order history by user, need Bearer token got from the login API to be included in header

```
URL: POST /order
Content-Type: application/json
```
```
Parameters:
page_index = int // default will be 1
page_size = int // default will be 10
search = string // can be used to search by title or by author
```
##### Response:
```json
{
    "result": true,
    "data": [
        {
            "order_id": 2,
            "total_amount": 35.96,
            "status": "NEW",
            "created_at": 1718388109572,
            "updated_at": 1718388109572,
            "items": [
                {
                    "item_id": 2,
                    "book_id": 10,
                    "quantity": 2,
                    "price": 9.99
                },
                {
                    "item_id": 3,
                    "book_id": 2,
                    "quantity": 2,
                    "price": 7.99
                }
            ]
        },
        {
            "order_id": 1,
            "total_amount": 19.98,
            "status": "NEW",
            "created_at": 1718387948631,
            "updated_at": 1718387948631,
            "items": [
                {
                    "item_id": 1,
                    "book_id": 10,
                    "quantity": 2,
                    "price": 9.99
                }
            ]
        }
    ]
}
```