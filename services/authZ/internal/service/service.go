package authz

import (
	"context"

	"github.com/ashish19912009/zrms/services/authz/internal/domain"
	"github.com/ashish19912009/zrms/services/internal/infrastructure/opa"
)

type Service struct {
	opaClient *opa.Client
	repo      domain.PolicyRepository
}

func NewService(opaClient *opa.Client, repo domain.PolicyRepository) *Service {
	return &Service{
		opaClient: opaClient,
		repo:      repo,
	}
}

func (s *Service) CheckPermission(ctx context.Context, subject, action, resource string) (bool, error) {
	// Get all policies for the subject
	policies, err := s.repo.GetPoliciesForSubject(ctx, subject)
	if err != nil {
		return false, err
	}

	// Prepare input for OPA
	input := map[string]interface{}{
		"subject":  subject,
		"action":   action,
		"resource": resource,
		"policies": policies,
	}

	// Evaluate with OPA
	return s.opaClient.Evaluate(ctx, input)
}

func (s *Service) AddPolicy(ctx context.Context, policy *domain.Policy) error {
	return s.repo.AddPolicy(ctx, policy)
}
