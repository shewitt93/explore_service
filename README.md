# Explore Service

## Overview

This is my implementation of Muzz' explore microservice

I have sent an invitation to Alex to join a postman workspace I've created so you can easily test my implementation



## API Endpoints

The service implements four gRPC endpoints as defined in the protocol buffer:

1. `ListLikedYou`: Lists all users who liked the specified recipient
2. `ListNewLikedYou`: Lists users who liked the recipient but haven't been liked back
3. `CountLikedYou`: Counts the number of users who liked the recipient
4. `PutDecision`: Records a user's decision (like or pass) about another user

### Example Requests and Responses

#### ListLikedYou

**Request (initial):**
```json
{
  "recipient_user_id": "1",
  "pagination_token": null
}
```

**Response:**
```json
{
  "likers": [
    { "actor_id": "10", "unix_timestamp": 1738754100 },
    { "actor_id": "9", "unix_timestamp": 1738686000 },
    { "actor_id": "7", "unix_timestamp": 1738511100 },
    { "actor_id": "6", "unix_timestamp": 1738404600 },
    { "actor_id": "5", "unix_timestamp": 1738057800 },
    { "actor_id": "3", "unix_timestamp": 1737621000 },
    { "actor_id": "2", "unix_timestamp": 1737390600 }
  ],
  "next_pagination_token": null
}
```

**Request (with pagination):**
```json
{
  "recipient_user_id": "1",
  "pagination_token": "eyJ1cGRhdGVkX2F0IjoiMjAyNS0wMi0wMiAxNTo0NTowMCIsImFjdG9yX2lkIjoiNyJ9"
}
```

**Response:**
```json
{
  "likers": [
    { "actor_id": "6", "unix_timestamp": 1738404600 },
    { "actor_id": "5", "unix_timestamp": 1738057800 },
    { "actor_id": "3", "unix_timestamp": 1737621000 },
    { "actor_id": "2", "unix_timestamp": 1737390600 }
  ],
  "next_pagination_token": "eyJ1cGRhdGVkX2F0IjoiMjAyNS0wMS0yMyAwODozMDowMCIsImFjdG9yX2lkIjoiMyJ9"
}
```

#### CountLikedYou

**Request:**
```json
{
  "recipient_user_id": "1"
}
```

**Response:**
```json
{
  "count": 7
}
```

#### PutDecision

**Request:**
```json
{
  "actor_user_id": "1",
  "recipient_user_id": "2",
  "liked_recipient": true
}
```

**Response:**
```json
{
  "mutual_likes": true
}
```

## Technical Implementation

### Database Schema
This is defined in the `init.sql` file.

### Cursor-Based Pagination

For efficient pagination of large result sets, the service uses cursor-based pagination:

- Each page request can include an optional pagination token
- The token encodes the timestamp and `actor_id` of the last item from the previous page
- This approach is more efficient than offset-based pagination for large datasets

### Performance Considerations

- Optimized database queries with appropriate indexes
- Efficient handling of mutual likes check using transactions
- Cursor-based pagination for consistent performance with large datasets
- Connection pooling for database access
- For the sake of time, I haven't added Redis caching, also partly due to me and Alex previously talking about no caching strategies during our chat

## Deployment

- N/A

### Prerequisites

- Go 1.24
- MySQL 8.0+
- Docker and Docker Compose

### Environment Variables
For the sake of this assessment, I've simply committed this file.

### Running the project

```bash
# Build the service
go build main.go .

# Start the service
docker-compose up --build

# Stop the service
docker-compose down
```

## Testing

To run the test suite:
```bash
# Run tests
go test ./...
```

Sample data for testing can be loaded using:
```bash
# Copy seed file to container
docker cp seed-data.sql mysqldb:/tmp/

# Execute seed file
docker exec -it mysqldb bash -c "mysql -utest -ptest explore_muzz < /tmp/seed-data.sql"
```

## Design Decisions and Assumptions

1. **Composite Primary Key**: Using `actor_id` and `recipient_id` as a composite primary key ensures each user can have only one decision about another user.
2. **TEXT vs VARCHAR**: `VARCHAR` is used for ID columns since they are used in indices and have a known maximum length.
3. **Cursor-Based Pagination**: This approach was chosen for its scalability with large datasets.
4. **Transaction for Mutual Likes**: The mutual like check is performed within a transaction to ensure data consistency.
5. **Error Handling**: Comprehensive error handling and informative error messages are provided.

