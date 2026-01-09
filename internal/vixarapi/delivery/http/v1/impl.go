package v1

import (
	"context"

	gen "github.com/keenywheels/backend/internal/api/v1"
)

// DeleteUserSearchQuery implements DeleteUserSearchQuery for gen.Handler
func (r *Router) DeleteUserSearchQuery(
	ctx context.Context,
	params gen.DeleteUserSearchQueryParams,
) (gen.DeleteUserSearchQueryRes, error) {
	return r.userController.DeleteUserSearchQuery(ctx, params)
}

// GetUserSearchQueries implements GetUserSearchQueries for gen.Handler
func (r *Router) GetUserSearchQueries(
	ctx context.Context,
	params gen.GetUserSearchQueriesParams,
) (gen.GetUserSearchQueriesRes, error) {
	return r.userController.GetUserSearchQueries(ctx, params)
}

// LogoutUser implements Logout for gen.Handler
func (r *Router) LogoutUser(
	ctx context.Context,
) (gen.LogoutUserRes, error) {
	return r.userController.LogoutUser(ctx)
}

// UserInfo implements UserInfo for gen.Handler
func (r *Router) UserInfo(
	ctx context.Context,
) (gen.UserInfoRes, error) {
	return r.userController.UserInfo(ctx)
}

// SaveUserQuery implements SaveUserQuery for gen.Handler
func (r *Router) SaveUserQuery(
	ctx context.Context,
	req *gen.SaveUserQueryRequest,
) (gen.SaveUserQueryRes, error) {
	return r.userController.SaveUserQuery(ctx, req)
}

// VkAuthCallback implements VkAuthCallback for gen.Handler
func (r *Router) VkAuthCallback(
	ctx context.Context,
	req *gen.VkAuthCallbackRequest,
) (gen.VkAuthCallbackRes, error) {
	return r.userController.VkAuthCallback(ctx, req)
}

// VkAuthRegister implements VkAuthRegister for gen.Handler
func (r *Router) VkAuthRegister(
	ctx context.Context,
	req *gen.VkAuthRegisterRequest,
) (gen.VkAuthRegisterRes, error) {
	return r.userController.VkAuthRegister(ctx, req)
}

// SubscribeUserToToken implements SubscribeUserToToken for gen.Handler
func (r *Router) SubscribeUserToToken(
	ctx context.Context,
	req *gen.SubscribeUserToTokenRequest,
) (gen.SubscribeUserToTokenRes, error) {
	return r.userController.SubscribeUserToToken(ctx, req)
}

// GetUserTokenSubs implements GetUserTokenSubs for gen.Handler
func (r *Router) GetUserTokenSubs(
	ctx context.Context,
	params gen.GetUserTokenSubsParams,
) (gen.GetUserTokenSubsRes, error) {
	return r.userController.GetUserTokenSubs(ctx, params)
}

// DeleteUserTokenSub implements DeleteUserTokenSub for gen.Handler
func (r *Router) DeleteUserTokenSub(
	ctx context.Context,
	params gen.DeleteUserTokenSubParams,
) (gen.DeleteUserTokenSubRes, error) {
	return r.userController.DeleteUserTokenSub(ctx, params)
}

// SearchTokenInfo implements SearchTokenInfo for gen.Handler
func (r *Router) SearchTokenInfo(
	ctx context.Context,
	req *gen.SearchTokenInfoRequest,
) (gen.SearchTokenInfoRes, error) {
	return r.searchController.SearchTokenInfo(ctx, req)
}
