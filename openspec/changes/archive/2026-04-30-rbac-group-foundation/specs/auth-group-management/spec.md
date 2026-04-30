## ADDED Requirements

### Requirement: Group CRUD operations
The auth-service SHALL provide CreateGroup, UpdateGroup, DeleteGroup, and ListGroup gRPC handlers that persist Group entities with name, description, parent_group_id, model_patterns, token_limit, and rate_limit fields.

#### Scenario: Create a new group
- **WHEN** CreateGroup is called with name "developers" and description "Developer team"
- **THEN** a Group entity is persisted with a generated UUID, the provided fields, and empty model_patterns/token_limit/rate_limit defaults

#### Scenario: Create group with model patterns and limits
- **WHEN** CreateGroup is called with name "power-users", model_patterns=["gpt-4","claude-*"], token_limit={prompt_tokens:100000,completion_tokens:100000,period:"daily"}, rate_limit={requests_per_minute:60,requests_per_day:10000}
- **THEN** the Group entity is persisted with all provided configuration

#### Scenario: Update a group
- **WHEN** UpdateGroup is called with an existing group ID and new name "senior-devs"
- **THEN** the Group entity is updated with the new name and updated_at timestamp

#### Scenario: Delete a group
- **WHEN** DeleteGroup is called with an existing group ID
- **THEN** the Group entity and all associated UserGroupMembership records are removed

#### Scenario: List groups with pagination
- **WHEN** ListGroups is called with page=1, page_size=10
- **THEN** up to 10 groups are returned with total count

### Requirement: User-Group membership management
The auth-service SHALL provide AddUserToGroup and RemoveUserToGroup gRPC handlers that manage UserGroupMembership records linking users to groups.

#### Scenario: Add user to group
- **WHEN** AddUserToGroup is called with user_id and group_id
- **THEN** a UserGroupMembership record is created with a generated UUID and added_at timestamp

#### Scenario: Add user to group that they are already in
- **WHEN** AddUserToGroup is called with a user_id and group_id that already has a membership
- **THEN** the operation SHALL return an error indicating duplicate membership

#### Scenario: Remove user from group
- **WHEN** RemoveUserFromGroup is called with user_id and group_id
- **THEN** the UserGroupMembership record is deleted

#### Scenario: Remove user from group they are not in
- **WHEN** RemoveUserFromGroup is called with a user_id and group_id that has no membership
- **THEN** the operation SHALL return success (idempotent)

### Requirement: Group-scoped model patterns
A Group entity SHALL carry a model_patterns field (list of glob patterns) that defines which models members of the group are authorized to access. This field is stored but not enforced in this sprint.

#### Scenario: Group with model patterns
- **WHEN** a Group is created with model_patterns=["gpt-4","claude-*"]
- **THEN** the patterns are persisted and retrievable via GetByID/List

### Requirement: Group-scoped token limits
A Group entity SHALL carry an optional token_limit field with prompt_tokens, completion_tokens, and period. This field is stored but not enforced in this sprint.

#### Scenario: Group with token limit
- **WHEN** a Group is created with token_limit={prompt_tokens:50000,completion_tokens:50000,period:"daily"}
- **THEN** the limit is persisted and retrievable via GetByID/List

### Requirement: Group-scoped rate limits
A Group entity SHALL carry an optional rate_limit field with requests_per_minute and requests_per_day. This field is stored but not enforced in this sprint.

#### Scenario: Group with rate limit
- **WHEN** a Group is created with rate_limit={requests_per_minute:30,requests_per_day:5000}
- **THEN** the limit is persisted and retrievable via GetByID/List
