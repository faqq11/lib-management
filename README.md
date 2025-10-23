# library-management

### 1. Register User

**Endpoint:**
```http
POST /api/register
```

**Request Body:**
```json
{
  "username": "string (required)",
  "password": "string (required)", 
  "role": "string (optional)" // user or admin. default = user
}
```

**Success Response (201 Created):**
```json
{
  "message": "User created successfully"
}
```

**Error Responses (400-500):**
```json
{
  "message": "error message"
}
```

### 2. Login

**Endpoint:**
```http
POST /api/login
```

**Request Body:**
```json
{
  "username": "string (required)",
  "password": "string (required)"
}
```

**Success Response (200 OK):**
```json
{
  "access_token": "your token"
}
```

**Error Responses (400-500):**
```json
{
  "message": "error message"
}
```

### 3. Create Category (admin only)

**Endpoint:**
```http
POST /api/create-category
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "name": "string (required)"
}
```

**Success Response (201 created):**
```json
{
  "message": "Category created successfully"
}
```

**Error Responses (400-500):**
```json
{
  "message": "error message"
}
```

### 4. Delete Category (admin only)

**Endpoint:**
```http
DELETE /api/delete-category/{id}
Authorization: Bearer <token>
```

**Success Response (200 OK):**
```json
{
  "message": "Category created successfully"
}
```

**Error Responses (400-500):**
```json
{
  "message": "error message"
}
```

### 5. Create Books (admin only)

**Endpoint:**
```http
POST /api/create-book
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "title": "string (required)",
  "author": "string (required)",
  "stock": "integer (optional)", // default 0
  "category_id": "integer (optional)",
}
```

**Success Response (201 created):**
```json
{
  "message": "Book created successfully"
}
```

**Error Responses (400-500):**
```json
{
  "message": "error message"
}
```

### 6. Get All Books (Need to login)

**Endpoint:**
```http
GET /api/books
Authorization: Bearer <token>
```

**Success Response (200 OK):**
```json
[
  {
      "id": 2,
      "title": "coba2",
      "author": "coba2",
      "category_id": 6,
      "category": "ini judul",
      "stock": 4,
      "created_at": "2025-10-23T20:42:59.300571+07:00"
  },
  {
      "id": 4,
      "title": "halo dunia",
      "author": "123",
      "category_id": 7,
      "category": "ini judul1",
      "stock": 1,
      "created_at": "2025-10-24T01:25:59.808258+07:00"
  }
]
```

**Error Responses (400-500):**
```json
{
  "message": "error message"
}
```

### 7. Search and filter (Need to login)

**Endpoint:**
```http
GET /api/books/search?
Authorization: Bearer <token>
```

**Success Response (200 OK):**
```json
[
  {
      "id": 2,
      "title": "coba2",
      "author": "coba2",
      "category_id": 6,
      "category": "ini judul",
      "stock": 4,
      "created_at": "2025-10-23T20:42:59.300571+07:00"
  },
  {
      "id": 4,
      "title": "halo dunia",
      "author": "123",
      "category_id": 7,
      "category": "ini judul1",
      "stock": 1,
      "created_at": "2025-10-24T01:25:59.808258+07:00"
  }
]
```

**Error Responses (400-500):**
```json
{
  "message": "error message"
}
```

### 8. Get book by id (Need to login)

**Endpoint:**
```http
GET /api/books/{id}
Authorization: Bearer <token>
```

**Success Response (200 OK):**
```json
  {
      "id": 2,
      "title": "coba2",
      "author": "coba2",
      "category_id": 6,
      "category": "ini judul",
      "stock": 4,
      "created_at": "2025-10-23T20:42:59.300571+07:00"
  },
  
```

**Error Responses (400-500):**
```json
{
  "message": "error message"
}
```

### 9. Update book (admin only)

**Endpoint:**
```http
PUT /api/books/{id}
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "title": "string",
  "author": "string", 
  "category_id": "integer",
  "stock": "integer"
}
```

**Success Response (200 OK):**
```json
{
  "message": "Book updated successfully",
}
```

**Error Responses (400-500):**
```json
{
  "message": "error message"
}
```

### 10. Increase book stock(admin only)

**Endpoint:**
```http
PUT /api/books/{id}/increase-stock
Authorization: Bearer <token>
```

**Success Response (200 OK):**
```json
{
  "message": "string",
}
```

**Error Responses (400-500):**
```json
{
  "message": "error message"
}
```

### 11. Decrease book stock(admin only)

**Endpoint:**
```http
PUT /api/books/{id}/decrease-stock
Authorization: Bearer <token>
```

**Success Response (200 OK):**
```json
{
  "message": "string",
}
```

**Error Responses (400-500):**
```json
{
  "message": "error message"
}
```

### 12. Delete book (admin only)

**Endpoint:**
```http
DELETE /api/books/{id}/delete
Authorization: Bearer <token>
```

**Success Response (200 OK):**
```json
{
  "message": "string",
}
```

**Error Responses (400-500):**
```json
{
  "message": "error message"
}
```

### 13. Borrowed books by user (need to login)

**Endpoint:**
```http
GET /api/my-borrowings
Authorization: Bearer <token>
```

**Success Response (200 OK):**
```json
[
  {
    "ID": "integer",
    "BookID": "integer",
    "BookTitle": "string",
    "Author": "string",
    "BorrowedAt": "time",
    "ReturnedAt": "time",
    "Status": "string"
  }
]
```

**Error Responses (400-500):**
```json
{
  "message": "error message"
}
```

### 14. Borrow book (need to login)

**Endpoint:**
```http
POST /api/books/{id}/borrow
Authorization: Bearer <token>
```

**Success Response (200 OK):**
```json
{
  "message": "string"
}
```

**Error Responses (400-500):**
```json
{
  "message": "error message"
}
```

### 15. Return book (need to login)

**Endpoint:**
```http
PUT /api/borrowings/{id}/return
Authorization: Bearer <token>
```

**Success Response (200 OK):**
```json
{
  "message": "string"
}
```

**Error Responses (400-500):**
```json
{
  "message": "error message"
}
```