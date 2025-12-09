/*
 * GoAstra CLI - Buf Config Template
 *
 * Generates Buf configuration files for Protocol Buffer management.
 * Configures code generation for Connect-Go.
 */
package trpc

// BufYAML returns the buf.yaml configuration template.
func BufYAML() string {
	return `version: v1
name: buf.build/goastra/api
breaking:
  use:
    - FILE
lint:
  use:
    - DEFAULT
  except:
    - PACKAGE_VERSION_SUFFIX
`
}

// BufGenYAML returns the buf.gen.yaml code generation template.
func BufGenYAML() string {
	return `version: v1
managed:
  enabled: true
  go_package_prefix:
    default: app/internal/rpc/gen
plugins:
  - plugin: buf.build/protocolbuffers/go
    out: internal/rpc/gen
    opt: paths=source_relative
  - plugin: buf.build/connectrpc/go
    out: internal/rpc/gen
    opt: paths=source_relative
`
}

// BufWorkYAML returns the buf.work.yaml workspace template.
func BufWorkYAML() string {
	return `version: v1
directories:
  - proto
`
}

// ServiceGo returns the RPC service implementations template.
func ServiceGo() string {
	return `package rpc

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
// HEALTH SERVICE
// ============================================================================

/*
 * HealthService implements the health check RPC service.
 */
type HealthService struct {
	protov1connect.UnimplementedHealthServiceHandler
	db *database.Client
}

/*
 * NewHealthService creates a new health service.
 */
func NewHealthService(db *database.Client) *HealthService {
	return &HealthService{db: db}
}

/*
 * Check returns the service health status.
 */
func (s *HealthService) Check(
	ctx context.Context,
	req *connect.Request[pb.HealthCheckRequest],
) (*connect.Response[pb.HealthCheckResponse], error) {
	dbStatus := "connected"
	if s.db != nil {
		if err := s.db.Health(); err != nil {
			dbStatus = "disconnected"
		}
	} else {
		dbStatus = "not configured"
	}

	return connect.NewResponse(&pb.HealthCheckResponse{
		Status:   "healthy",
		Database: dbStatus,
		Version:  "1.0.0",
	}), nil
}

// ============================================================================
// AUTH SERVICE
// ============================================================================

/*
 * AuthService implements the authentication RPC service.
 */
type AuthService struct {
	protov1connect.UnimplementedAuthServiceHandler
	db  *database.Client
	cfg *config.Config
}

/*
 * NewAuthService creates a new auth service.
 */
func NewAuthService(db *database.Client, cfg *config.Config) *AuthService {
	return &AuthService{db: db, cfg: cfg}
}

/*
 * Login authenticates a user.
 */
func (s *AuthService) Login(
	ctx context.Context,
	req *connect.Request[pb.LoginRequest],
) (*connect.Response[pb.AuthResponse], error) {
	// TODO: Implement login
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("not implemented"))
}

/*
 * Register creates a new user account.
 */
func (s *AuthService) Register(
	ctx context.Context,
	req *connect.Request[pb.RegisterRequest],
) (*connect.Response[pb.AuthResponse], error) {
	// TODO: Implement registration
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("not implemented"))
}

/*
 * RefreshToken refreshes the access token.
 */
func (s *AuthService) RefreshToken(
	ctx context.Context,
	req *connect.Request[pb.RefreshTokenRequest],
) (*connect.Response[pb.AuthResponse], error) {
	// TODO: Implement token refresh
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("not implemented"))
}

/*
 * Logout invalidates the user's tokens.
 */
func (s *AuthService) Logout(
	ctx context.Context,
	req *connect.Request[pb.LogoutRequest],
) (*connect.Response[pb.LogoutResponse], error) {
	return connect.NewResponse(&pb.LogoutResponse{Success: true}), nil
}

// ============================================================================
// USER SERVICE
// ============================================================================

/*
 * UserService implements the user management RPC service.
 */
type UserService struct {
	protov1connect.UnimplementedUserServiceHandler
	db  *database.Client
	cfg *config.Config
}

/*
 * NewUserService creates a new user service.
 */
func NewUserService(db *database.Client, cfg *config.Config) *UserService {
	return &UserService{db: db, cfg: cfg}
}

/*
 * GetUser returns a user by ID.
 */
func (s *UserService) GetUser(
	ctx context.Context,
	req *connect.Request[pb.GetUserRequest],
) (*connect.Response[pb.GetUserResponse], error) {
	// TODO: Implement get user
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("not implemented"))
}

/*
 * ListUsers returns a paginated list of users.
 */
func (s *UserService) ListUsers(
	ctx context.Context,
	req *connect.Request[pb.ListUsersRequest],
) (*connect.Response[pb.ListUsersResponse], error) {
	// TODO: Implement list users
	return connect.NewResponse(&pb.ListUsersResponse{
		Users:      []*pb.User{},
		Total:      0,
		Page:       1,
		PageSize:   20,
		TotalPages: 0,
	}), nil
}

/*
 * CreateUser creates a new user.
 */
func (s *UserService) CreateUser(
	ctx context.Context,
	req *connect.Request[pb.CreateUserRequest],
) (*connect.Response[pb.CreateUserResponse], error) {
	// TODO: Implement create user
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("not implemented"))
}

/*
 * UpdateUser updates an existing user.
 */
func (s *UserService) UpdateUser(
	ctx context.Context,
	req *connect.Request[pb.UpdateUserRequest],
) (*connect.Response[pb.UpdateUserResponse], error) {
	// TODO: Implement update user
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("not implemented"))
}

/*
 * DeleteUser deletes a user.
 */
func (s *UserService) DeleteUser(
	ctx context.Context,
	req *connect.Request[pb.DeleteUserRequest],
) (*connect.Response[pb.DeleteUserResponse], error) {
	// TODO: Implement delete user
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("not implemented"))
}
`
}

// InterceptorGo returns the Connect interceptor template.
func InterceptorGo() string {
	return `package rpc

import (
	"context"
	"time"

	"connectrpc.com/connect"

	"app/internal/logger"
)

/*
 * NewLoggingInterceptor returns an interceptor that logs RPC calls.
 */
func NewLoggingInterceptor(log *logger.Logger) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			start := time.Now()
			procedure := req.Spec().Procedure

			resp, err := next(ctx, req)

			latency := time.Since(start)
			if err != nil {
				log.Errorw("RPC failed",
					"procedure", procedure,
					"latency_ms", latency.Milliseconds(),
					"error", err,
				)
			} else {
				log.Infow("RPC completed",
					"procedure", procedure,
					"latency_ms", latency.Milliseconds(),
				)
			}

			return resp, err
		}
	}
}

/*
 * NewAuthInterceptor returns an interceptor that validates JWT tokens.
 */
func NewAuthInterceptor(secret string) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			// Skip auth for public procedures
			procedure := req.Spec().Procedure
			if isPublicProcedure(procedure) {
				return next(ctx, req)
			}

			// Get token from header
			token := req.Header().Get("Authorization")
			if token == "" {
				return nil, connect.NewError(connect.CodeUnauthenticated, nil)
			}

			// TODO: Validate token and add user to context

			return next(ctx, req)
		}
	}
}

func isPublicProcedure(procedure string) bool {
	publicProcedures := map[string]bool{
		"/proto.v1.HealthService/Check": true,
		"/proto.v1.AuthService/Login":    true,
		"/proto.v1.AuthService/Register": true,
	}
	return publicProcedures[procedure]
}
`
}
