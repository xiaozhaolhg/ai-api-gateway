# auth-service API Contracts

> Identity and access control service gRPC API

## Service Definition

```protobuf
service AuthService {
  // Authentication
  rpc ValidateAPIKey(ValidateAPIKeyRequest) returns (UserIdentity);
  
  // Model Authorization
  rpc CheckModelAuthorization(CheckModelAuthorizationRequest) returns (AuthorizationResult);
  
  // User Management
  rpc GetUser(GetUserRequest) returns (User);
  rpc CreateUser(CreateUserRequest) returns (User);
  rpc UpdateUser(UpdateUserRequest) returns (User);
  rpc DeleteUser(DeleteUserRequest) returns (Empty);
  rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
  
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
  string role = 4;
  string status = 5;    // "active" | "disabled"
  int64 created_at = 6;
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
  string parent_group_id = 3;
  int64 created_at = 4;
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