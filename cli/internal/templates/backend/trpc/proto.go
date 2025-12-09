/*
 * GoAstra CLI - tRPC Proto Template
 *
 * Generates Protocol Buffer definitions for tRPC services.
 * Provides User, Auth, and Health service definitions.
 */
package trpc

// ServiceProto returns the service.proto template.
func ServiceProto() string {
	return `syntax = "proto3";

package proto.v1;

option go_package = "app/internal/rpc/gen/proto/v1;protov1";

// ============================================================================
// HEALTH SERVICE
// ============================================================================

service HealthService {
  rpc Check(HealthCheckRequest) returns (HealthCheckResponse);
}

message HealthCheckRequest {}

message HealthCheckResponse {
  string status = 1;
  string database = 2;
  string version = 3;
}

// ============================================================================
// AUTH SERVICE
// ============================================================================

service AuthService {
  rpc Login(LoginRequest) returns (AuthResponse);
  rpc Register(RegisterRequest) returns (AuthResponse);
  rpc RefreshToken(RefreshTokenRequest) returns (AuthResponse);
  rpc Logout(LogoutRequest) returns (LogoutResponse);
}

message LoginRequest {
  string email = 1;
  string password = 2;
}

message RegisterRequest {
  string email = 1;
  string password = 2;
  string name = 3;
}

message RefreshTokenRequest {
  string refresh_token = 1;
}

message LogoutRequest {}

message LogoutResponse {
  bool success = 1;
}

message AuthResponse {
  string token = 1;
  string refresh_token = 2;
  int64 expires_at = 3;
  User user = 4;
}

// ============================================================================
// USER SERVICE
// ============================================================================

service UserService {
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
  rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
  rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse);
  rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse);
}

message User {
  uint64 id = 1;
  string email = 2;
  string name = 3;
  string role = 4;
  bool active = 5;
  string created_at = 6;
  string updated_at = 7;
}

message GetUserRequest {
  uint64 id = 1;
}

message GetUserResponse {
  User user = 1;
}

message ListUsersRequest {
  int32 page = 1;
  int32 page_size = 2;
}

message ListUsersResponse {
  repeated User users = 1;
  int32 total = 2;
  int32 page = 3;
  int32 page_size = 4;
  int32 total_pages = 5;
}

message CreateUserRequest {
  string email = 1;
  string password = 2;
  string name = 3;
  string role = 4;
}

message CreateUserResponse {
  User user = 1;
}

message UpdateUserRequest {
  uint64 id = 1;
  optional string email = 2;
  optional string name = 3;
  optional string role = 4;
  optional bool active = 5;
}

message UpdateUserResponse {
  User user = 1;
}

message DeleteUserRequest {
  uint64 id = 1;
}

message DeleteUserResponse {
  bool success = 1;
}
`
}
