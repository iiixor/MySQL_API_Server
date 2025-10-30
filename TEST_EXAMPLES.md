# Server Testing Examples

## Health Check
```bash
curl http://localhost:8080/health
```

## Basic Query
```bash
curl -X POST http://localhost:8080/api/v1/execute \
  -H "Content-Type: application/json" \
  -d '{"query":"CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(100)); INSERT INTO users VALUES (1, \"John\"); SELECT * FROM users;"}'
```

## Multiple Tables with JOIN
```bash
curl -X POST http://localhost:8080/api/v1/execute \
  -H "Content-Type: application/json" \
  -d '{"query":"CREATE TABLE users (id INT PRIMARY KEY, name VARCHAR(50)); CREATE TABLE orders (id INT PRIMARY KEY, user_id INT, product VARCHAR(50)); INSERT INTO users VALUES (1, \"Alice\"), (2, \"Bob\"); INSERT INTO orders VALUES (1, 1, \"Laptop\"), (2, 1, \"Mouse\"), (3, 2, \"Keyboard\"); SELECT u.name, o.product FROM users u JOIN orders o ON u.id = o.user_id;"}'
```

## Aggregate Functions
```bash
curl -X POST http://localhost:8080/api/v1/execute \
  -H "Content-Type: application/json" \
  -d '{"query":"CREATE TABLE sales (id INT, amount DECIMAL(10,2)); INSERT INTO sales VALUES (1, 100.50), (2, 200.75), (3, 150.25); SELECT COUNT(*) as total, SUM(amount) as total_amount, AVG(amount) as average FROM sales;"}'
```

## Security Tests

### DROP DATABASE (should be blocked)
```bash
curl -X POST http://localhost:8080/api/v1/execute \
  -H "Content-Type: application/json" \
  -d '{"query":"DROP DATABASE test;"}'
```

### LOAD_FILE (should be blocked)
```bash
curl -X POST http://localhost:8080/api/v1/execute \
  -H "Content-Type: application/json" \
  -d '{"query":"SELECT LOAD_FILE(\"/etc/passwd\");"}'
```

### CREATE USER (should be blocked)
```bash
curl -X POST http://localhost:8080/api/v1/execute \
  -H "Content-Type: application/json" \
  -d '{"query":"CREATE USER \"hacker\"@\"localhost\" IDENTIFIED BY \"password\";"}'
```

## Error Handling

### Syntax Error
```bash
curl -X POST http://localhost:8080/api/v1/execute \
  -H "Content-Type: application/json" \
  -d '{"query":"SELEC * FROM users;"}'
```

### Empty Query
```bash
curl -X POST http://localhost:8080/api/v1/execute \
  -H "Content-Type: application/json" \
  -d '{"query":""}'
```
