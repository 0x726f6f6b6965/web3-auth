# Web3 Auth

Using the crypto wallet to authenticate.

## Table

| #   | access pattern         | target | action      | pk                    | sk                        | done               |
| --- | ---------------------- | ------ | ----------- | --------------------- | ------------------------- | ------------------ |
| 1   | get user information   | table  | get item    | USER#`public_address` | #PROFILE#`public_address` | :white_check_mark: |
| 2   | set user               | table  | put item    | USER#`public_address` | #PROFILE#`public_address` | :white_check_mark: |
| 3   | update user infomation | table  | update item | USER#`public_address` | #PROFILE#`public_address` | :white_check_mark: |

## Flow

```mermaid
sequenceDiagram
    Client->>Service: Send get nonce request.
    Service->>Dynamodb: Get user information by public address.
    alt pk exist
    Dynamodb->>Service: Found user.
    Service->>Dynamodb: Generate nonce and update.
    Service->>Client: Send nonce.
    else pk not exist
    Dynamodb-->>Service: User not found.
    Service-->>Dynamodb: Generate nonce and create user info.
    Service-->>Client: Send nonce.
    end
    Client->>Service: Generate signature with nonce.
    Service->>Dynamodb: Get user information by public address.
    Dynamodb->>Service: Return user.
    Service->>Service: Verify signature and generate JWT token.
    Service->>Client: Send JWT token.
  
```