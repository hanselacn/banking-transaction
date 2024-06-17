# Banking Transaction - Service

## Quick Start
```
# Build migration 
Execute SQL on ./database/migration files to migrate all table

# Run App
go mod tidy
go mod vendor
go mod download
go run server.go

You have to make a super admin account to grant access to other users.

there is two ways :
1. Without Seeder
uncomment/enable this route :
/banking-transaction/users/create/supadmin
and use that for one time purpose only (create account without authorization)
after that directly change the role on that user in table users to 'super_admin'

2. With Seeder :
use env.example as .env to prevent error while using existing AES encrypt decryption (seeder data use the example AES KEY)
excecute SQL on ./database/seeder



After those steps, you can disable that route and begin to use the features.
```

## Project Structure

```
database
    migrations  # Contains required database migrations
    seeder      # Contains seeder data for testing

internal             # Contain all internal dependecies
    handler          # Contain handler and routing management.
    business         # Contain the business logic.
    consts           # Contains reusable consts.
    entity           # Contains reusable entities.
    middleware       # Contains middleware for http requests.
    pkg              # Contains reusable packages.
    repo             # Contain all repository dependecies
    
```

## Environment
```
Server will load dependencies based on .env file, you can use the env.example or custom it by yourself.
If you want to mantain a secure connection, you can set the value "API_TLS" on .env to "true"

Custom Default Interest by changing "DEFAULT_INTEREST_RATE"
If needed, you can custom interest rate per-User by API Request

AES key are used for Encrypt and Decrypt Balance and Interest Rate. 
You can custom "AES_KEY" with 256 bits long key

"PAYOUT_INTERVAL" will set interval each time the interest will be paid.
"PAYOUT_TIME_UNIT" sets the desired unit (ex: YEAR, MONTH, DAYS, etc).
Payout schedule will be calculate based on those values
example :
PAYOUT_INTERVAL = 3
PAYOUT_TIME_UNIT = DAYS
means interest payout occured every 3 days
```

## Health Check
- `/ping` basic end point for health check. Returns the value below
```
Pong!
```

## Features
```
This Application have multiple features :
- Create Users
- Update Users Role
- Deposit
- Withdrawal
- Manual Interest Payout
- Update Custom Interest
- Scheduled Interest Payout

```

## Security Measurements

```
All REST API on the application applied with security procedures
That includes :
- Authentication Middleware to prevent unauthorized access
- RBAC implementation to grant permissions which the access rights to perform certain operations on resources.
- Data Encryption to database (Balance and Interest Rate).
- TLS enabled for secure connection
- Payload and Input validation to prevent XSS attack
- Placeholder ($) usage on queries to prevent SQL Injection
```