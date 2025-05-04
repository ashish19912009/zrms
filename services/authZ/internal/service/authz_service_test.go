package service_test

// import (
// 	"context"
// 	"errors"
// 	"testing"

// 	"github.com/ashish19912009/zrms/services/authZ/internal/model"
// 	"github.com/ashish19912009/zrms/services/authZ/internal/service"
// 	"github.com/ashish19912009/zrms/services/authZ/internal/service/mocks"
// 	"github.com/open-policy-agent/opa/rego"
// 	"github.com/stretchr/testify/assert"
// )

// func prepareTestService(t *testing.T) (service.AuthZService, *mocks.AuthZService, rego.PreparedEvalQuery) {
// 	mockRepo := new(mocks.AuthZService)

// 	query, err := rego.New(
// 		rego.Query("data.zrms.authz.allow = true, data.zrms.authz.reason = \"allowed by test\""),
// 	).PrepareForEval(context.Background())
// 	assert.NoError(t, err)

// 	svc, _ := service.NewAuthZService(mockRepo, query)
// 	return svc, mockRepo, query
// }

// func TestIsAuthorized_Success(t *testing.T) {
// 	svc, repo, _ := prepareTestService(t)
// 	ctx := context.Background()

// 	repo.On("GetAccountRole", ctx, "fr1", "acc1").
// 		Return("acc1", "fr1", "role1", nil)

// 	repo.On("GetRolePermissions", ctx, "role1").
// 		Return(map[string][]string{"res1": {"read"}}, nil)

// 	repo.On("GetDirectPermissions", ctx, "acc1").
// 		Return(map[string]bool{"res1:write": true}, nil)

// 	allowed, reason, err := (*svc).IsAuthorized(ctx, "fr1", "acc1", "res1", "read", nil)

// 	assert.NoError(t, err)
// 	assert.True(t, allowed)
// 	assert.Equal(t, "allowed by test", reason)
// }

// func TestIsAuthorized_ContextCancelled(t *testing.T) {
// 	svc, _, _ := prepareTestService(t)

// 	ctx, cancel := context.WithCancel(context.Background())
// 	cancel()

// 	allowed, reason, err := (*svc).IsAuthorized(ctx, "fr1", "acc1", "res1", "read", nil)

// 	assert.Error(t, err)
// 	assert.False(t, allowed)
// 	assert.Empty(t, reason)
// }

// func TestIsAuthorized_AccountFranchiseMismatch(t *testing.T) {
// 	svc, repo, _ := prepareTestService(t)
// 	ctx := context.Background()

// 	repo.On("GetAccountRole", ctx, "fr1", "acc1").
// 		Return("accX", "frX", "role1", nil)

// 	allowed, reason, err := (*svc).IsAuthorized(ctx, "fr1", "acc1", "res1", "read", nil)

// 	assert.Error(t, err)
// 	assert.False(t, allowed)
// 	assert.Empty(t, reason)
// }

// func TestIsAuthorized_RepoError(t *testing.T) {
// 	svc, repo, _ := prepareTestService(t)
// 	ctx := context.Background()

// 	repo.On("GetAccountRole", ctx, "fr1", "acc1").
// 		Return("", "", "", errors.New("db error"))

// 	allowed, reason, err := (*svc).IsAuthorized(ctx, "fr1", "acc1", "res1", "read", nil)

// 	assert.Error(t, err)
// 	assert.False(t, allowed)
// 	assert.Empty(t, reason)
// }

// func TestIsAuthorizedBatch_Success(t *testing.T) {
// 	svc, repo, _ := prepareTestService(t)
// 	ctx := context.Background()

// 	repo.On("GetAccountRole", ctx, "fr1", "acc1").
// 		Return("acc1", "fr1", "role1", nil)

// 	repo.On("GetRolePermissions", ctx, "role1").
// 		Return(map[string][]string{"res1": {"read"}}, nil)

// 	repo.On("GetDirectPermissions", ctx, "acc1").
// 		Return(map[string]bool{"res2:write": true}, nil)

// 	resources := []*model.ResourceAction{
// 		{Resource: "res1", Action: "read"},
// 		{Resource: "res2", Action: "write"},
// 	}

// 	result, err := (*svc).IsAuthorizedBatch(ctx, "fr1", "acc1", resources)
// 	assert.NoError(t, err)
// 	assert.Len(t, result, 2)
// 	assert.True(t, result[0].Allowed)
// 	assert.True(t, result[1].Allowed)
// }

// func TestIsAuthorizedBatch_InvalidMapping(t *testing.T) {
// 	svc, repo, _ := prepareTestService(t)
// 	ctx := context.Background()

// 	repo.On("GetAccountRole", ctx, "fr1", "acc1").
// 		Return("accX", "frX", "role1", nil)

// 	_, err := (*svc).IsAuthorizedBatch(ctx, "fr1", "acc1", []*model.ResourceAction{})
// 	assert.Error(t, err)
// }
