package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ashish19912009/zrms/services/authZ/internal/constants"
	"github.com/ashish19912009/zrms/services/authZ/internal/logger"
	"github.com/ashish19912009/zrms/services/authZ/internal/model"
	"github.com/ashish19912009/zrms/services/authZ/internal/repository"
	"github.com/ashish19912009/zrms/services/authZ/pb"
	"github.com/open-policy-agent/opa/rego"
)

var layer = "service"
var decisionTTL = 24 * time.Hour // update this from config file

type AuthZService interface {
	IsAuthorized(ctx context.Context, franchiseID, accountID, resource, action string, meta map[string]string) (bool, string, int64, int64, string, error)
	IsAuthorizedBatch(ctx context.Context, franchiseID, accountID string, resources []model.ResourceAction) ([]*model.CheckBatchAccessResponse, error)
}

type authZService struct {
	drepo       repository.AuthZRepository
	policyQuery rego.PreparedEvalQuery // precompiled rego query
	cRepo       repository.CacheRepository
}

func NewAuthZService(drepo repository.AuthZRepository, policyPath string, cacheRepo repository.CacheRepository) (AuthZService, error) {
	var method = constants.Methods.NewAuthZService
	ctx := context.Background()
	compiledQuery, err := rego.New(
		rego.Query(`{
					"allow":data.zrms.services.authz.allow,
					"deny_reason": data.zrms.services.authz.deny_reason,
					"policy_version":data.zrms.services.authz.policy_version
		}`),
		rego.Load([]string{policyPath}, nil),
	).PrepareForEval(ctx)

	if err != nil {
		logCtx := logger.BaseLogContext(
			"layer", layer,
			"method", method,
		)
		logger.Error(constants.FailedPreparePolicy, err, logCtx)
		return nil, fmt.Errorf(constants.FailedPreparePolicy, err)
	}

	return &authZService{
		drepo:       drepo,
		policyQuery: compiledQuery,
		cRepo:       cacheRepo,
	}, nil
}

func makeCacheKey(accountID, resource, action string) (string, string) {
	return fmt.Sprintf("account_id:%s:", accountID), fmt.Sprintf("%s:%s:", resource, action)
}

// IsAuthorized checks if an account has permission to perform an action on a resource
func (s *authZService) IsAuthorized(ctx context.Context, franchiseID, accountID, resource, action string, meta map[string]string) (bool, string, int64, int64, string, error) {
	var method = constants.Methods.IsAuthorized
	select {
	case <-ctx.Done():
		return false, "", 0, 0, "", ctx.Err()
	default:
	}
	logCtx := logger.BaseLogContext(
		"layer", layer,
		"method", method,
	)
	accID, frID, roleID, err := s.drepo.GetAccountRole(ctx, franchiseID, accountID)
	if err != nil {
		logger.Error(constants.FailedFetchAccount, err, logCtx)
		return false, "", 0, 0, "", fmt.Errorf(constants.FailedFetchAccount, err)
	}
	if franchiseID == "" {
		logger.Error(constants.AccoutNotAssociated, err, logCtx)
		return false, "", 0, 0, "", fmt.Errorf("Error: %s", constants.AccoutNotAssociated)
	}
	if franchiseID != frID || accountID != accID {
		logger.Error(constants.IdMismatch, err, logCtx)
		return false, "", 0, 0, "", fmt.Errorf("Error: %s", constants.IdMismatch)
	}

	result := &pb.Decision{}
	var tenantPrefix, resourceActionPostfix = makeCacheKey(accountID, resource, action)
	err = s.cRepo.Get(ctx, tenantPrefix, resourceActionPostfix, result)
	if err != nil {
		logger.Error(constants.WrongFetchingData, err, nil)
	}
	if err == nil {
		return result.Allowed, result.Reason, result.IssuedAt, result.ExpiresAt, result.PolicyVersion, nil
	}

	rolePermissions, err := s.drepo.GetRolePermissions(ctx, roleID)
	if err != nil {
		logger.Error(constants.FailedFetchRolePermission, err, logCtx)
		return false, "", 0, 0, "", fmt.Errorf(constants.FailedFetchRolePermission, err)
	}

	directPermissionsFlat, err := s.drepo.GetDirectPermissions(ctx, accountID)
	if err != nil {
		logger.Error(constants.FailedFetchDPermission, err, logCtx)
		return false, "", 0, 0, "", fmt.Errorf(constants.FailedFetchDPermission, err)
	}
	directPermissions := convertDirectPermissions(directPermissionsFlat)

	finalPermissions := mergePermissions(rolePermissions, directPermissions)

	input := map[string]interface{}{
		"resource":    resource,
		"action":      action,
		"permissions": buildOPAInputPermissions(finalPermissions),
	}
	allowed, reason, issued_at, expires_at, policy_version, err := s.evaluatePolicy(ctx, input)
	if err != nil {
		logger.Error(constants.EvaluationErr, err, logCtx)
		return false, "", 0, 0, "", fmt.Errorf(constants.EvaluationErr, err)
	}
	if time.Duration(decisionTTL) > 0 {
		//fmt.Printf("TTL tenantPrefix:%s, resourceActionPostfix: %s", tenantPrefix, resourceActionPostfix)
		s.cRepo.StoreWithTTL(ctx, tenantPrefix, resourceActionPostfix, &pb.Decision{
			Allowed:       allowed,
			Reason:        reason,
			IssuedAt:      issued_at.Unix(),
			ExpiresAt:     expires_at.Unix(),
			PolicyVersion: policy_version,
		}, time.Duration(decisionTTL))
	} else {
		//fmt.Printf("tenantPrefix:%s, resourceActionPostfix: %s", tenantPrefix, resourceActionPostfix)
		s.cRepo.Store(ctx, tenantPrefix, resourceActionPostfix, &pb.Decision{
			Allowed:       allowed,
			Reason:        reason,
			IssuedAt:      issued_at.Unix(),
			ExpiresAt:     expires_at.Unix(),
			PolicyVersion: policy_version,
		})
	}
	return allowed, reason, issued_at.Unix(), expires_at.Unix(), policy_version, nil
}

func (s *authZService) IsAuthorizedBatch(
	ctx context.Context,
	franchiseID, accountID string,
	resources []model.ResourceAction,
) ([]*model.CheckBatchAccessResponse, error) {
	var method = constants.Methods.IsAuthorizedBatch
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	logCtx := logger.BaseLogContext(
		"layer", layer,
		"method", method,
	)

	responses := make([]*model.CheckBatchAccessResponse, len(resources))
	for i := range responses {
		responses[i] = &model.CheckBatchAccessResponse{
			Resource:      resources[i].Resource,
			Action:        resources[i].Action,
			Allowed:       false,
			Reason:        "pending evaluation",
			IssuedAt:      0,
			ExpiresAt:     0,
			PolicyVersion: "",
		}
	}

	// 1. Verify account and get role
	accID, frID, roleID, err := s.drepo.GetAccountRole(ctx, franchiseID, accountID)
	if err != nil {
		logger.Error(constants.FailedFetchAccount, err, logCtx)
		return nil, fmt.Errorf(constants.FailedFetchAccount, err)
	}
	if franchiseID == "" || franchiseID != frID || accountID != accID {
		return nil, fmt.Errorf("Error: %s", constants.InvalidAssociation)
	}

	// 3. First pass - check cache
	cacheMisses := make([]int, 0) // Track indices of cache misses
	for i, rec := range resources {
		result := &pb.Decision{}
		tenantPrefix, resourceActionPostfix := makeCacheKey(accountID, rec.Resource, rec.Action)
		err := s.cRepo.Get(ctx, tenantPrefix, resourceActionPostfix, result)

		if err == nil && result.Allowed {
			// Cache hit with allowed decision
			responses[i] = &model.CheckBatchAccessResponse{
				Resource:      rec.Resource,
				Action:        rec.Action,
				Allowed:       result.Allowed,
				Reason:        result.Reason,
				IssuedAt:      result.IssuedAt,
				ExpiresAt:     result.ExpiresAt,
				PolicyVersion: result.PolicyVersion,
			}
		} else {
			// Cache miss or denied - mark for evaluation
			cacheMisses = append(cacheMisses, i)
		}
	}

	// 4. If all decisions were cached, return early
	if len(cacheMisses) == 0 {
		return responses, nil
	}

	// 5. Fetch permissions for cache misses
	rolePermissions, err := s.drepo.GetRolePermissions(ctx, roleID)
	if err != nil {
		logger.Error(constants.FailedFetchRolePermission, err, logCtx)
		return nil, fmt.Errorf(constants.FailedFetchRolePermission, err)
	}

	directPermissionsFlat, err := s.drepo.GetDirectPermissions(ctx, accountID)
	if err != nil {
		logger.Error(constants.FailedFetchDPermission, err, logCtx)
		return nil, fmt.Errorf(constants.FailedFetchDPermission, err)
	}

	finalPermissions := mergePermissions(rolePermissions, convertDirectPermissions(directPermissionsFlat))
	opaPermissions := buildOPAInputPermissions(finalPermissions)

	// 6. Evaluate policy for cache misses
	for _, idx := range cacheMisses {
		ra := resources[idx]
		input := map[string]interface{}{
			"resource":    ra.Resource,
			"action":      ra.Action,
			"permissions": opaPermissions,
		}

		allowed, reason, issuedAt, expiresAt, policyVersion, err := s.evaluatePolicy(ctx, input)
		if err != nil {
			logger.Error(constants.FailedOPAEval, err, logCtx)
			return nil, fmt.Errorf(constants.FailedOPAEval, ra.Resource, ra.Action, err)
		}

		// Create response
		responses[idx] = &model.CheckBatchAccessResponse{
			Resource:      ra.Resource,
			Action:        ra.Action,
			Allowed:       allowed,
			Reason:        reason,
			IssuedAt:      issuedAt.Unix(),
			ExpiresAt:     expiresAt.Unix(),
			PolicyVersion: policyVersion,
		}

		// Cache the decision
		tenantPrefix, resourceActionPostfix := makeCacheKey(accountID, ra.Resource, ra.Action)
		decision := &pb.Decision{
			Allowed:       allowed,
			Reason:        reason,
			IssuedAt:      issuedAt.Unix(),
			ExpiresAt:     expiresAt.Unix(),
			PolicyVersion: policyVersion,
		}

		if decisionTTL > 0 {
			s.cRepo.StoreWithTTL(ctx, tenantPrefix, resourceActionPostfix, decision, time.Duration(decisionTTL))
		} else {
			s.cRepo.Store(ctx, tenantPrefix, resourceActionPostfix, decision)
		}
	}

	return responses, nil
}

// mergePermissions combines role and direct permissions
func mergePermissions(rolePerms map[string][]string, directPerms map[string]map[string]bool) map[string]map[string]bool {
	final := make(map[string]map[string]bool)

	// Add role permissions first (default allow)
	for res, actions := range rolePerms {
		if final[res] == nil {
			final[res] = make(map[string]bool)
		}
		for _, act := range actions {
			final[res][act] = true
		}
	}

	// Apply direct permissions (override role permissions)
	for res, acts := range directPerms {
		if final[res] == nil {
			final[res] = make(map[string]bool)
		}
		for act, isGranted := range acts {
			final[res][act] = isGranted
		}
	}

	return final
}

// buildOPAInputPermissions prepares the structure needed for OPA input
func buildOPAInputPermissions(perms map[string]map[string]bool) map[string]map[string]interface{} {
	out := make(map[string]map[string]interface{})

	for res, actions := range perms {
		for act, allowed := range actions {
			key := res + ":" + act
			out[key] = map[string]interface{}{
				"resource": res,
				"action":   act,
				"allowed":  allowed,
			}
		}
	}

	return out
}

// convertDirectPermissions transforms flat direct permissions into resource -> action -> allowed structure
func convertDirectPermissions(flat map[string]bool) map[string]map[string]bool {
	result := make(map[string]map[string]bool)

	for key, allowed := range flat {
		var resource, action string
		n, err := fmt.Sscanf(key, "%[^:]:%s", &resource, &action)
		if err != nil || n != 2 {
			continue // skip invalid entries
		}

		if result[resource] == nil {
			result[resource] = make(map[string]bool)
		}
		result[resource][action] = allowed
	}

	return result
}

// evaluatePolicy uses precompiled policy query to check permissions
func (s *authZService) evaluatePolicy(ctx context.Context, input map[string]any) (bool, string, time.Time, time.Time, string, error) {
	select {
	case <-ctx.Done():
		logger.Error("evaluation cancelled", nil, nil)
		return false, "", time.Time{}, time.Time{}, "", ctx.Err()
	default:
	}

	results, err := s.policyQuery.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		logger.Error(constants.RegoEvalFailed, err, nil)
		return false, "", time.Time{}, time.Time{}, "", err
	}

	if len(results) == 0 {
		// Policy didn't return any decision — default deny
		logger.Warn(constants.PolicyDenied, nil)
		return false, "", time.Time{}, time.Time{}, "", nil
	}

	issuedAt := time.Now()
	expiresAt := issuedAt.Add(24 * time.Hour)

	// Default values
	allowed := false
	reason := ""
	policyVersion := ""

	// Process the first result (assuming you only expect one)
	if len(results) > 0 && len(results[0].Expressions) > 0 {
		if resultMap, ok := results[0].Expressions[0].Value.(map[string]interface{}); ok {
			// Get allow
			if allowVal, ok := resultMap["allow"]; ok {
				if allowBool, ok := allowVal.(bool); ok {
					if !allowBool {
						expiresAt = issuedAt.Add(1 * time.Hour)
					}
					allowed = allowBool
				}
			}

			// Get deny_reason (note the correct field name)
			if reasonVal, ok := resultMap["deny_reason"]; ok {
				if reasonStr, ok := reasonVal.(string); ok {
					reason = reasonStr
				}
			}

			// Get policy_version
			if versionVal, ok := resultMap["policy_version"]; ok {
				if versionStr, ok := versionVal.(string); ok {
					policyVersion = versionStr
				}
			}
		}
	}
	return allowed, reason, issuedAt, expiresAt, policyVersion, nil
}
