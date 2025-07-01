# 🧠 AI Agent Instructions: Implement Feature and Unit Tests in Go Web Application

## 📝 Objective

You will receive a **feature definition** via prompt. Your task is to:

1. Implement the feature in a Go web application.
2. Write comprehensive unit tests for it.

---

## 📁 Project Structure Assumptions

Assume a standard Go web application layout:
```
├── main.go
├── cmd
│   ├── domain 
│   │   ├── charge 
│   │   │   ├── dao.go # Here should be the chgarge dao inteface to be implemented by infra
│   │   │   ├── builder.go # Here should be the charge builder if its necessary
│   │   │   └── entity.go # Here must have the charge struct and its rules
│   │   ├── payment 
│   │   │   ├── dao.go # Here should be the payment dao inteface to be implemented by infra
│   │   │   ├── builder.go # Here should be the payment builder if its necessary
│   │   │   └── entity.go # Here must have the payment struct and its rules
│   │   └── order
│   │       ├── dao.go # Here should be the order dao inteface to be implemented by infra
│   │       ├── builder.go # Here should be the order builder if its necessary
│   │       └── entity.go # Here must have the order struct and its rules
│   ├── infra
│   │   ├── conf 
│   │   │   ├── routes.go # This is an implementation of the routes
│   │   │   └── runtime.go # This is must have all instantiations of the needed dependencies
│   │   ├── dao 
│   │   │   ├── payment_dao.go # This is an implementation of dao.PaymentDao interface
│   │   │   ├── charge_dao.go # This is an implementation of dao.ChargeDao interface
│   │   │   └── order_dao.go # This is an implementation of dao.OrderDao interface
│   │   ├── db 
│   │   │   ├── mysql
│   │   │   │   └── client.go # This is an implementation of db.Client interface in MySQL
│   │   │   └── client.go # This is an interface for the database client
│   │   ├── handler # Here should be all handler to the features
│   │   └── configuration.go
│   └── usecases # Here should be all use cases to the features
└── go.mod
```

---

## 🚧 Implementation Guidelines

### 1. Parse the Feature Definition

Carefully analyze the feature definition provided in the prompt. Identify:

- Endpoints and HTTP methods
- Data models (input/output)
- Business rules and constraints
- Error handling expectations

### 2. Create/Update Files

Depending on the feature, implement or update:

- `cmd/infra/handler` – Implement the handler for the feature.
- `cmd/usecases` – Implement the use case for the feature.
- `cmd/infra/dao` – Implement the data access object (DAO) for the feature.
- `cmd/infra/conf/routes.go` Register any new routes/endpoints.
- `cmd/domain` – Create or update the domain model for the feature.


### 3. Create/Update Tests

- For each feature, write tests in the same directory as the implementation;
- The file name should be `<feature_name>_test.go`;
- The package name should be `<feature_name>_test`;
- The test should cover the implementation of the feature;
- The test should cover the business rules and constraints;
- The test should cover the error handling expectations;
- The test should cover the edge cases and invalid inputs;
- The test should cover the successful execution paths.
 
### 4. Run the tests to assert the implementation
- After the implementation, run the tests to ensure the implementation is correct.
- If the tests fail, fix the implementation to make the tests pass.

---

### Recommended Tools

- Standard `testing` package
- `net/http/httptest` for handler testing
- `github.com/stretchr/testify/assert` (optional, if project allows)
- `github.com/stretchr/testify/mock`

## 🧠 Best Practices

- Use the others files in the same directory as template for the implementation
- Keep business logic out of handlers
- Use clear and consistent error messages
- Structure code for readability and maintainability
- Favor composition over inheritance

---


## When complete, provide:

- Updated or new Go source files
- A list of files modified or created
- The complete set of unit tests

Do not include build artifacts or third-party binaries.

---

## ⏳ Waiting for Feature Definition...

Respond with the full implementation once the prompt provides the specific feature details.

