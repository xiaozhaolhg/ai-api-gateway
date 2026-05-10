# auth-service API Contracts

> Identity and access control service gRPC API

## Service Definition

```protobuf
service AuthService {
  // Authentication
  rpc ValidateAPIKey(ValidateAPIKeyRequest) returns (UserIdentity);
  rpc Login(LoginRequest) returns (LoginResponse);  // Supports email or username
  rpc Register(RegisterRequest) returns (User);     // Creates user with mandatory username

  // Model Authorization
  rpc CheckModelAuthorization(CheckModelAuthorizationRequest) returns (AuthorizationResult);

  // User Management
  rpc GetUser(GetUserRequest) returns (User);
  rpc CreateUser(CreateUserRequest) returns (User);
  rpc UpdateUser(UpdateUserRequest) returns (User);
  rpc DeleteUser(DeleteUserRequest) returns (Empty);
  rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
  rpc CheckUsernameAvailability(CheckUsernameAvailabilityRequest) returns (CheckUsernameAvailabilityResponse);
  
  // API Key Management
  rpc CreateAPIKey(CreateAPIKeyRequest) returns (CreateAPIKeyResponse);
  rpc DeleteAPIKey(DeleteAPIKeyRequest) returns (Empty);
  rpc ListAPIKeys(ListAPIKeysRequest) returns (ListAPIKeysResponse);
  
  // Group Management (Phase 2+)
  rpc CreateGroup(CreateGroupRequest) returns (Group);
  rpc UpdateGroup(UpdateGroupRequest) returns (Group);
  rpc DeleteGroup(DeleteGroupRequest) returns (Empty);
  rpc ListGroups(ListGroupsRequest) returns (ListGroupsResponse);
  rpc AddUserToGroup(AddUserToGroupRequest) returns (Empty);
  rpc RemoveUserFromGroup(RemoveUserFromGroupRequest) returns (Empty);
  
  // Permission Management (Phase 2+)
  rpc GrantPermission(GrantPermissionRequest) returns (Permission);
  rpc RevokePermission(RevokePermissionRequest) returns (Empty);
  rpc ListPermissions(ListPermissionsRequest) returns (ListPermissionsResponse);
  rpc CheckPermission(CheckPermissionRequest) returns (CheckPermissionResponse);
}
```

## Request/Response Messages

### Authentication

```protobuf
message ValidateAPIKeyRequest {
  string api_key = 1;
}

message UserIdentity {
  string user_id = 1;
  string role = 2;              // "admin" | "user"
  repeated string group_ids = 3;
  repeated string scopes = 4;
}
```

### Model Authorization

```protobuf
message CheckModelAuthorizationRequest {
  string user_id = 1;
  repeated string group_ids = 2;
  string model = 3;
}

message AuthorizationResult {
  bool allowed = 1;
  string reason = 2;
  repeated string authorized_models = 3;
}
```

### User

```protobuf
message User {
  string id = 1;
  string name = 2;
  string email = 3;
  string username = 4;  // Mandatory for new users, immutable after creation
  string role = 5;
  string status = 6;    // "active" | "disabled"
  int64 created_at = 7;
  repeated string group_ids = 8;  // User's group memberships
}
```

### Login (New)

```protobuf
message LoginRequest {
  string email = 1;      // Optional: used if username not provided
  string username = 2;  // Optional: used if email not provided
  string password = 3;
}

message LoginResponse {
  string token = 1;     // JWT token
  User user = 2;
}
```

### Register (New)

```protobuf
message RegisterRequest {
  string username = 1;  // Mandatory, unique, immutable after creation
  string email = 2;
  string name = 3;
  string password = 4;
  string role = 5;      // Optional: defaults to "user"
}
```

### CheckUsernameAvailability (New)

```protobuf
message CheckUsernameAvailabilityRequest {
  string username = 1;
}

message CheckUsernameAvailabilityResponse {
  bool available = 1;
}
```

### API Key

```protobuf
message CreateAPIKeyRequest {
  string user_id = 1;
  string name = 2;
}

message CreateAPIKeyResponse {
  string api_key_id = 1;
  string api_key = 2;    // returned once only
}
```

### Group (Phase 2+)

```protobuf
message Group {
  string id = 1;
  string name = 2;
  string description = 3;    // Group description
  string parent_group_id = 4;
  int64 created_at = 5;
  int32 member_count = 6;    // Number of members in the group
}
```

### CreateGroupRequest

```protobuf
message CreateGroupRequest {
  string name = 1;
  string description = 2;   // Optional group description
  string parent_group_id = 3;
}
```

### UpdateGroupRequest

```protobuf
message UpdateGroupRequest {
  string id = 1;
  string name = 2;
  string description = 3;   // Optional group description
}
```

### Permission (Phase 2+)

```protobuf
message Permission {
  string id = 1;
  string group_id = 2;
  string resource_type = 3;    // "model" | "provider" | "admin_feature"
  string resource_id = 4;
  string action = 5;           // "access" | "manage" | "view"
}
```