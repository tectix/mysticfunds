package mana

import (
	"context"

	"github.com/stretchr/testify/mock"
	wizardpb "github.com/tectix/mysticfunds/proto/wizard"
	"google.golang.org/grpc"
)

// MockWizardServiceClient is a mock implementation of WizardServiceClient
type MockWizardServiceClient struct {
	mock.Mock
}

func (m *MockWizardServiceClient) GetManaBalance(ctx context.Context, req *wizardpb.GetManaBalanceRequest, opts ...grpc.CallOption) (*wizardpb.GetManaBalanceResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*wizardpb.GetManaBalanceResponse), args.Error(1)
}

func (m *MockWizardServiceClient) UpdateManaBalance(ctx context.Context, req *wizardpb.UpdateManaBalanceRequest, opts ...grpc.CallOption) (*wizardpb.UpdateManaBalanceResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*wizardpb.UpdateManaBalanceResponse), args.Error(1)
}

func (m *MockWizardServiceClient) TransferMana(ctx context.Context, req *wizardpb.TransferManaRequest, opts ...grpc.CallOption) (*wizardpb.TransferManaResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*wizardpb.TransferManaResponse), args.Error(1)
}

// Add all other required methods to satisfy the interface (with minimal implementations)
func (m *MockWizardServiceClient) CreateWizard(ctx context.Context, req *wizardpb.CreateWizardRequest, opts ...grpc.CallOption) (*wizardpb.Wizard, error) {
	return nil, nil
}

func (m *MockWizardServiceClient) GetWizard(ctx context.Context, req *wizardpb.GetWizardRequest, opts ...grpc.CallOption) (*wizardpb.Wizard, error) {
	return nil, nil
}

func (m *MockWizardServiceClient) UpdateWizard(ctx context.Context, req *wizardpb.UpdateWizardRequest, opts ...grpc.CallOption) (*wizardpb.Wizard, error) {
	return nil, nil
}

func (m *MockWizardServiceClient) ListWizards(ctx context.Context, req *wizardpb.ListWizardsRequest, opts ...grpc.CallOption) (*wizardpb.ListWizardsResponse, error) {
	return nil, nil
}

func (m *MockWizardServiceClient) DeleteWizard(ctx context.Context, req *wizardpb.DeleteWizardRequest, opts ...grpc.CallOption) (*wizardpb.DeleteWizardResponse, error) {
	return nil, nil
}

func (m *MockWizardServiceClient) JoinGuild(ctx context.Context, req *wizardpb.JoinGuildRequest, opts ...grpc.CallOption) (*wizardpb.Wizard, error) {
	return nil, nil
}

func (m *MockWizardServiceClient) LeaveGuild(ctx context.Context, req *wizardpb.LeaveGuildRequest, opts ...grpc.CallOption) (*wizardpb.Wizard, error) {
	return nil, nil
}

func (m *MockWizardServiceClient) CreateJob(ctx context.Context, req *wizardpb.CreateJobRequest, opts ...grpc.CallOption) (*wizardpb.Job, error) {
	return nil, nil
}

func (m *MockWizardServiceClient) GetJob(ctx context.Context, req *wizardpb.GetJobRequest, opts ...grpc.CallOption) (*wizardpb.Job, error) {
	return nil, nil
}

func (m *MockWizardServiceClient) ListJobs(ctx context.Context, req *wizardpb.ListJobsRequest, opts ...grpc.CallOption) (*wizardpb.ListJobsResponse, error) {
	return nil, nil
}

func (m *MockWizardServiceClient) UpdateJob(ctx context.Context, req *wizardpb.UpdateJobRequest, opts ...grpc.CallOption) (*wizardpb.Job, error) {
	return nil, nil
}

func (m *MockWizardServiceClient) DeleteJob(ctx context.Context, req *wizardpb.DeleteJobRequest, opts ...grpc.CallOption) (*wizardpb.DeleteJobResponse, error) {
	return nil, nil
}

func (m *MockWizardServiceClient) AssignWizardToJob(ctx context.Context, req *wizardpb.AssignWizardToJobRequest, opts ...grpc.CallOption) (*wizardpb.JobAssignment, error) {
	return nil, nil
}

func (m *MockWizardServiceClient) GetJobAssignments(ctx context.Context, req *wizardpb.GetJobAssignmentsRequest, opts ...grpc.CallOption) (*wizardpb.GetJobAssignmentsResponse, error) {
	return nil, nil
}

func (m *MockWizardServiceClient) CompleteJobAssignment(ctx context.Context, req *wizardpb.CompleteJobAssignmentRequest, opts ...grpc.CallOption) (*wizardpb.JobAssignment, error) {
	return nil, nil
}

func (m *MockWizardServiceClient) CancelJobAssignment(ctx context.Context, req *wizardpb.CancelJobAssignmentRequest, opts ...grpc.CallOption) (*wizardpb.JobAssignment, error) {
	return nil, nil
}

func (m *MockWizardServiceClient) UpdateJobProgress(ctx context.Context, req *wizardpb.UpdateJobProgressRequest, opts ...grpc.CallOption) (*wizardpb.JobProgress, error) {
	return nil, nil
}

func (m *MockWizardServiceClient) GetJobProgress(ctx context.Context, req *wizardpb.GetJobProgressRequest, opts ...grpc.CallOption) (*wizardpb.JobProgress, error) {
	return nil, nil
}

func (m *MockWizardServiceClient) GetActivities(ctx context.Context, req *wizardpb.GetActivitiesRequest, opts ...grpc.CallOption) (*wizardpb.GetActivitiesResponse, error) {
	return nil, nil
}

func (m *MockWizardServiceClient) GetRealms(ctx context.Context, req *wizardpb.GetRealmsRequest, opts ...grpc.CallOption) (*wizardpb.GetRealmsResponse, error) {
	return nil, nil
}
