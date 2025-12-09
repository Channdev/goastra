/*
 * GoAstra CLI - tRPC Generator
 *
 * Generates Protocol Buffer definitions and Connect-Go service
 * implementations for new resources.
 */
package generator

import (
	"fmt"
	"os"
	"path/filepath"
)

/*
 * TRPCGenerator handles proto and service generation for tRPC.
 */
type TRPCGenerator struct {
	name       string
	pascalName string
	camelName  string
	snakeName  string
}

/*
 * NewTRPCGenerator creates a new tRPC generator instance.
 */
func NewTRPCGenerator(name string) *TRPCGenerator {
	return &TRPCGenerator{
		name:       name,
		pascalName: toPascalCase(name),
		camelName:  toCamelCase(name),
		snakeName:  toSnakeCase(name),
	}
}

/*
 * GenerateProto creates a Protocol Buffer definition for the resource.
 */
func (g *TRPCGenerator) GenerateProto() error {
	content := fmt.Sprintf(`syntax = "proto3";

package proto.v1;

option go_package = "app/internal/rpc/gen/proto/v1;protov1";

// ============================================================================
// %s SERVICE
// ============================================================================

/*
 * %sService provides CRUD operations for %s resources.
 */
service %sService {
  rpc Get%s(Get%sRequest) returns (Get%sResponse);
  rpc List%ss(List%ssRequest) returns (List%ssResponse);
  rpc Create%s(Create%sRequest) returns (Create%sResponse);
  rpc Update%s(Update%sRequest) returns (Update%sResponse);
  rpc Delete%s(Delete%sRequest) returns (Delete%sResponse);
}

// ============================================================================
// %s MESSAGES
// ============================================================================

/*
 * %s represents a %s entity.
 */
message %s {
  uint64 id = 1;
  // TODO: Add your fields here
  // string name = 2;
  // string description = 3;
  string created_at = 10;
  string updated_at = 11;
}

// ============================================================================
// REQUEST/RESPONSE MESSAGES
// ============================================================================

message Get%sRequest {
  uint64 id = 1;
}

message Get%sResponse {
  %s %s = 1;
}

message List%ssRequest {
  int32 page = 1;
  int32 page_size = 2;
}

message List%ssResponse {
  repeated %s %ss = 1;
  int32 total = 2;
  int32 page = 3;
  int32 page_size = 4;
  int32 total_pages = 5;
}

message Create%sRequest {
  // TODO: Add your create fields here
  // string name = 1;
  // string description = 2;
}

message Create%sResponse {
  %s %s = 1;
}

message Update%sRequest {
  uint64 id = 1;
  // TODO: Add your update fields here
  // optional string name = 2;
  // optional string description = 3;
}

message Update%sResponse {
  %s %s = 1;
}

message Delete%sRequest {
  uint64 id = 1;
}

message Delete%sResponse {
  bool success = 1;
}
`,
		g.pascalName,
		g.pascalName, g.name,
		g.pascalName,
		g.pascalName, g.pascalName, g.pascalName,
		g.pascalName, g.pascalName, g.pascalName,
		g.pascalName, g.pascalName, g.pascalName,
		g.pascalName, g.pascalName, g.pascalName,
		g.pascalName, g.pascalName, g.pascalName,
		g.pascalName,
		g.pascalName, g.name,
		g.pascalName,
		g.pascalName,
		g.pascalName,
		g.pascalName, g.snakeName,
		g.pascalName,
		g.pascalName, g.pascalName, g.snakeName,
		g.pascalName,
		g.pascalName,
		g.pascalName, g.snakeName,
		g.pascalName,
		g.pascalName,
		g.pascalName, g.snakeName,
		g.pascalName,
		g.pascalName,
	)

	path := filepath.Join("app/proto/v1", g.snakeName+".proto")
	return os.WriteFile(path, []byte(content), 0644)
}

/*
 * GenerateService creates a Connect-Go service implementation.
 */
func (g *TRPCGenerator) GenerateService() error {
	content := fmt.Sprintf(`package rpc

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	"app/internal/config"
	"app/internal/database"
	pb "app/internal/rpc/gen/proto/v1"
	"app/internal/rpc/gen/proto/v1/protov1connect"
)

// ============================================================================
// %s SERVICE
// ============================================================================

/*
 * %sService implements the %s RPC service.
 */
type %sService struct {
	protov1connect.Unimplemented%sServiceHandler
	db  *database.Client
	cfg *config.Config
}

/*
 * New%sService creates a new %s service.
 */
func New%sService(db *database.Client, cfg *config.Config) *%sService {
	return &%sService{db: db, cfg: cfg}
}

/*
 * Get%s returns a %s by ID.
 */
func (s *%sService) Get%s(
	ctx context.Context,
	req *connect.Request[pb.Get%sRequest],
) (*connect.Response[pb.Get%sResponse], error) {
	// TODO: Implement get %s
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("not implemented"))
}

/*
 * List%ss returns a paginated list of %s.
 */
func (s *%sService) List%ss(
	ctx context.Context,
	req *connect.Request[pb.List%ssRequest],
) (*connect.Response[pb.List%ssResponse], error) {
	// TODO: Implement list %ss
	return connect.NewResponse(&pb.List%ssResponse{
		%ss:        []*pb.%s{},
		Total:      0,
		Page:       1,
		PageSize:   20,
		TotalPages: 0,
	}), nil
}

/*
 * Create%s creates a new %s.
 */
func (s *%sService) Create%s(
	ctx context.Context,
	req *connect.Request[pb.Create%sRequest],
) (*connect.Response[pb.Create%sResponse], error) {
	// TODO: Implement create %s
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("not implemented"))
}

/*
 * Update%s updates an existing %s.
 */
func (s *%sService) Update%s(
	ctx context.Context,
	req *connect.Request[pb.Update%sRequest],
) (*connect.Response[pb.Update%sResponse], error) {
	// TODO: Implement update %s
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("not implemented"))
}

/*
 * Delete%s deletes a %s.
 */
func (s *%sService) Delete%s(
	ctx context.Context,
	req *connect.Request[pb.Delete%sRequest],
) (*connect.Response[pb.Delete%sResponse], error) {
	// TODO: Implement delete %s
	return connect.NewResponse(&pb.Delete%sResponse{Success: true}), nil
}
`,
		g.pascalName,
		g.pascalName, g.name,
		g.pascalName,
		g.pascalName,
		g.pascalName, g.name,
		g.pascalName, g.pascalName,
		g.pascalName,
		g.pascalName, g.name,
		g.pascalName, g.pascalName,
		g.pascalName, g.pascalName,
		g.name,
		g.pascalName, g.name,
		g.pascalName, g.pascalName,
		g.pascalName, g.pascalName,
		g.name,
		g.pascalName,
		g.pascalName, g.pascalName,
		g.pascalName, g.name,
		g.pascalName, g.pascalName,
		g.pascalName, g.pascalName,
		g.name,
		g.pascalName, g.name,
		g.pascalName, g.pascalName,
		g.pascalName, g.pascalName,
		g.name,
		g.pascalName, g.name,
		g.pascalName, g.pascalName,
		g.pascalName, g.pascalName,
		g.name,
		g.pascalName,
	)

	path := filepath.Join("app/internal/rpc", g.snakeName+"_service.go")
	return os.WriteFile(path, []byte(content), 0644)
}

/*
 * GenerateAll generates both proto and service.
 */
func (g *TRPCGenerator) GenerateAll() error {
	if err := g.GenerateProto(); err != nil {
		return err
	}
	return g.GenerateService()
}
