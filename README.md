# order-service

**Microservice for Managing Order Operations**

## Getting Started

### 1. Set Up Environment Variables

Create a `.env` file with the following configurations:

```bash
PRODUCT_SERVICE_HOST=http://localhost:8080
GET_PRODUCT_URI=/v1/products/
INCREASE_BOOKED_QUOTA=/v1/products/increase-booked-quota
DECREASE_BOOKED_QUOTA=/v1/products/decrease-booked-quota

AUTH_SERVICE_HOST=http://localhost:8081
AUTHENTICATE_URI=/v1/authenticate

DB_HOST=localhost
DB_PORT=3306
DB_USER={{YOUR-DB-USER}}
DB_PASSWORD={{YOUR-DB-PASSWORD}}
DB_NAME=order_svc
```

### 2. Create a new database:

```sql
CREATE DATABASE order_svc;
```

### 3. Run Database Migrations

Execute the SQL commands in `db/migrations` to set up the database schema.

### 4. Run the Application

Launch the application using the following command:

```bash
go run main.go
```

### 5. Access the Server

The server will be accessible at [http://localhost:8082](http://localhost:8082).

## Contributing

We welcome contributions! Feel free to submit issues, feature requests, or pull requests.

## License

This project is licensed under the [MIT License](LICENSE).

