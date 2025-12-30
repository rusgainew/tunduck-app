package service_impl

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rusgainew/tunduck-app/internal/models"
	"github.com/rusgainew/tunduck-app/internal/repository"
	"github.com/rusgainew/tunduck-app/pkg/entity"
	"github.com/rusgainew/tunduck-app/pkg/pagination"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ========== EsfDocumentService Tests ==========
// Comprehensive unit tests для EsfDocumentService

type MockDocumentRepository struct {
	mock.Mock
}

func (m *MockDocumentRepository) GetAllDocuments(ctx context.Context, orgID uuid.UUID) ([]entity.EsfDocument, error) {
	args := m.Called(ctx, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.EsfDocument), args.Error(1)
}

func (m *MockDocumentRepository) GetDocumentByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*entity.EsfDocument, error) {
	args := m.Called(ctx, orgID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EsfDocument), args.Error(1)
}

func (m *MockDocumentRepository) CreateDocument(ctx context.Context, orgID uuid.UUID, doc *entity.EsfDocument) error {
	args := m.Called(ctx, orgID, doc)
	return args.Error(0)
}

func (m *MockDocumentRepository) UpdateDocument(ctx context.Context, orgID uuid.UUID, doc *entity.EsfDocument) error {
	args := m.Called(ctx, orgID, doc)
	return args.Error(0)
}

func (m *MockDocumentRepository) DeleteDocument(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	args := m.Called(ctx, orgID, id)
	return args.Error(0)
}

func (m *MockDocumentRepository) GetAllDocumentsPaginated(ctx context.Context, orgID uuid.UUID, params pagination.PaginationParams, filters pagination.DocumentFilterParams) ([]entity.EsfDocument, int64, error) {
	args := m.Called(ctx, orgID, params, filters)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]entity.EsfDocument), args.Get(1).(int64), args.Error(2)
}

var _ repository.EsfDocumentRepository = (*MockDocumentRepository)(nil)

// ========== GetAllDocuments Tests ==========

func TestEsfDocumentGetAll_Success(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	orgID := uuid.New()

	mockRepo.On("GetAllDocuments", mock.Anything, orgID).Return([]entity.EsfDocument{}, nil)

	service := NewEsfDocumentService(mockRepo, nil, logrus.New())
	result, err := service.GetAllDocuments(context.Background(), orgID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestEsfDocumentGetAll_Error(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	orgID := uuid.New()

	mockRepo.On("GetAllDocuments", mock.Anything, orgID).Return(nil, errors.New("database error"))

	service := NewEsfDocumentService(mockRepo, nil, logrus.New())
	result, err := service.GetAllDocuments(context.Background(), orgID)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

// ========== GetDocumentByID Tests ==========

func TestEsfDocumentGetByID_Success(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	orgID := uuid.New()
	docID := uuid.New()
	doc := &entity.EsfDocument{ID: docID, ForeignName: "Test"}

	mockRepo.On("GetDocumentByID", mock.Anything, orgID, docID).Return(doc, nil)

	service := NewEsfDocumentService(mockRepo, nil, logrus.New())
	result, err := service.GetDocumentByID(context.Background(), orgID, docID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestEsfDocumentGetByID_NotFound(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	orgID := uuid.New()
	docID := uuid.New()

	mockRepo.On("GetDocumentByID", mock.Anything, orgID, docID).Return(nil, errors.New("not found"))

	service := NewEsfDocumentService(mockRepo, nil, logrus.New())
	result, err := service.GetDocumentByID(context.Background(), orgID, docID)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

// ========== CreateDocument Tests ==========

func TestEsfDocumentCreate_Success(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	orgID := uuid.New()

	mockRepo.On("CreateDocument", mock.Anything, orgID, mock.Anything).Return(nil)

	service := NewEsfDocumentService(mockRepo, nil, logrus.New())
	req := &models.EsfCreateDocumentRequest{
		ForeignName: "Test",
	}
	result, err := service.CreateDocument(context.Background(), orgID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestEsfDocumentCreate_Error(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	orgID := uuid.New()

	mockRepo.On("CreateDocument", mock.Anything, orgID, mock.Anything).Return(errors.New("constraint error"))

	service := NewEsfDocumentService(mockRepo, nil, logrus.New())
	req := &models.EsfCreateDocumentRequest{
		ForeignName: "Test",
	}
	result, err := service.CreateDocument(context.Background(), orgID, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

// ========== UpdateDocument Tests ==========

func TestEsfDocumentUpdate_Success(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	orgID := uuid.New()

	mockRepo.On("UpdateDocument", mock.Anything, orgID, mock.Anything).Return(nil)

	service := NewEsfDocumentService(mockRepo, nil, logrus.New())
	req := &models.EsfEditDocumentRequest{
		ID: uuid.New(),
		EsfCreateDocumentRequest: models.EsfCreateDocumentRequest{
			ForeignName: "Updated",
		},
	}
	err := service.UpdateDocument(context.Background(), orgID, req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestEsfDocumentUpdate_Error(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	orgID := uuid.New()

	mockRepo.On("UpdateDocument", mock.Anything, orgID, mock.Anything).Return(errors.New("update failed"))

	service := NewEsfDocumentService(mockRepo, nil, logrus.New())
	req := &models.EsfEditDocumentRequest{
		ID: uuid.New(),
		EsfCreateDocumentRequest: models.EsfCreateDocumentRequest{
			ForeignName: "Updated",
		},
	}
	err := service.UpdateDocument(context.Background(), orgID, req)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// ========== DeleteDocument Tests ==========

func TestEsfDocumentDelete_Success(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	orgID := uuid.New()
	docID := uuid.New()

	mockRepo.On("DeleteDocument", mock.Anything, orgID, docID).Return(nil)

	service := NewEsfDocumentService(mockRepo, nil, logrus.New())
	err := service.DeleteDocument(context.Background(), orgID, docID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestEsfDocumentDelete_Error(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	orgID := uuid.New()
	docID := uuid.New()

	mockRepo.On("DeleteDocument", mock.Anything, orgID, docID).Return(errors.New("deletion failed"))

	service := NewEsfDocumentService(mockRepo, nil, logrus.New())
	err := service.DeleteDocument(context.Background(), orgID, docID)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// ========== GetAllDocumentsPaginated Tests ==========

func TestEsfDocumentGetAllPaginated_Success(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	orgID := uuid.New()
	params := pagination.PaginationParams{Page: 1, PageSize: 10}
	filters := pagination.DocumentFilterParams{}

	mockRepo.On("GetAllDocumentsPaginated", mock.Anything, orgID, params, filters).Return([]entity.EsfDocument{}, int64(0), nil)

	service := NewEsfDocumentService(mockRepo, nil, logrus.New())
	result, total, err := service.GetAllDocumentsPaginated(context.Background(), orgID, params, filters)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(0), total)
	mockRepo.AssertExpectations(t)
}

func TestEsfDocumentGetAllPaginated_Error(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	orgID := uuid.New()
	params := pagination.PaginationParams{Page: 1, PageSize: 10}
	filters := pagination.DocumentFilterParams{}

	mockRepo.On("GetAllDocumentsPaginated", mock.Anything, orgID, params, filters).Return(nil, int64(0), errors.New("database error"))

	service := NewEsfDocumentService(mockRepo, nil, logrus.New())
	result, total, err := service.GetAllDocumentsPaginated(context.Background(), orgID, params, filters)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, int64(0), total)
	mockRepo.AssertExpectations(t)
}

// ========== Context Tests ==========

func TestEsfDocumentWithCanceledContext(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	orgID := uuid.New()

	// Create a canceled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	mockRepo.On("GetAllDocuments", ctx, orgID).Return(nil, errors.New("context canceled"))

	service := NewEsfDocumentService(mockRepo, nil, logrus.New())
	result, err := service.GetAllDocuments(ctx, orgID)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestEsfDocumentWithTimeout(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	orgID := uuid.New()

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	mockRepo.On("GetDocumentByID", ctx, orgID, mock.Anything).Return(nil, errors.New("context deadline exceeded"))

	service := NewEsfDocumentService(mockRepo, nil, logrus.New())
	result, err := service.GetDocumentByID(ctx, orgID, uuid.New())

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

// ========== Edge Cases ==========

func TestEsfDocumentEmptyListResponse(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	orgID := uuid.New()

	mockRepo.On("GetAllDocuments", mock.Anything, orgID).Return([]entity.EsfDocument{}, nil)

	service := NewEsfDocumentService(mockRepo, nil, logrus.New())
	result, err := service.GetAllDocuments(context.Background(), orgID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result))
	mockRepo.AssertExpectations(t)
}

func TestEsfDocumentNilOrgID(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	nilOrgID := uuid.Nil

	mockRepo.On("GetAllDocuments", mock.Anything, nilOrgID).Return(nil, errors.New("invalid organization id"))

	service := NewEsfDocumentService(mockRepo, nil, logrus.New())
	result, err := service.GetAllDocuments(context.Background(), nilOrgID)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestEsfDocumentLargeDocumentList(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	orgID := uuid.New()

	// Create large list of documents
	docs := make([]entity.EsfDocument, 1000)
	for i := 0; i < 1000; i++ {
		docs[i] = entity.EsfDocument{
			ID:          uuid.New(),
			ForeignName: "Doc " + string(rune(i%26+65)),
		}
	}

	mockRepo.On("GetAllDocuments", mock.Anything, orgID).Return(docs, nil)

	service := NewEsfDocumentService(mockRepo, nil, logrus.New())
	result, err := service.GetAllDocuments(context.Background(), orgID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1000, len(result))
	mockRepo.AssertExpectations(t)
}

// ========== Concurrency Tests ==========

func TestEsfDocumentConcurrentReads(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	orgID := uuid.New()
	docID := uuid.New()
	doc := &entity.EsfDocument{ID: docID, ForeignName: "Concurrent Doc"}

	// Allow multiple concurrent calls
	mockRepo.On("GetDocumentByID", mock.Anything, orgID, docID).Return(doc, nil)

	service := NewEsfDocumentService(mockRepo, nil, logrus.New())

	var wg sync.WaitGroup
	numGoroutines := 10
	results := make([]*models.EsfCreateDocumentRequest, numGoroutines)
	errors := make([]error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			result, err := service.GetDocumentByID(context.Background(), orgID, docID)
			results[idx] = result
			errors[idx] = err
		}(i)
	}

	wg.Wait()

	// All concurrent reads should succeed
	for i := 0; i < numGoroutines; i++ {
		assert.NoError(t, errors[i])
		assert.NotNil(t, results[i])
	}

	mockRepo.AssertExpectations(t)
}

func TestEsfDocumentConcurrentOperations(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	orgID := uuid.New()

	// Mock multiple operations
	mockRepo.On("GetAllDocuments", mock.Anything, orgID).Return([]entity.EsfDocument{}, nil)
	mockRepo.On("CreateDocument", mock.Anything, orgID, mock.Anything).Return(nil)
	mockRepo.On("UpdateDocument", mock.Anything, orgID, mock.Anything).Return(nil)
	mockRepo.On("DeleteDocument", mock.Anything, orgID, mock.Anything).Return(nil)

	service := NewEsfDocumentService(mockRepo, nil, logrus.New())

	var wg sync.WaitGroup

	// Concurrent Get
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, _ = service.GetAllDocuments(context.Background(), orgID)
	}()

	// Concurrent Create
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, _ = service.CreateDocument(context.Background(), orgID, &models.EsfCreateDocumentRequest{})
	}()

	// Concurrent Update
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = service.UpdateDocument(context.Background(), orgID, &models.EsfEditDocumentRequest{ID: uuid.New()})
	}()

	// Concurrent Delete
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = service.DeleteDocument(context.Background(), orgID, uuid.New())
	}()

	wg.Wait()

	// All operations should be tracked
	mockRepo.AssertNumberOfCalls(t, "GetAllDocuments", 1)
	mockRepo.AssertNumberOfCalls(t, "CreateDocument", 1)
	mockRepo.AssertNumberOfCalls(t, "UpdateDocument", 1)
	mockRepo.AssertNumberOfCalls(t, "DeleteDocument", 1)
}

// ========== Error Handling Tests ==========

func TestEsfDocumentCreateWithInvalidData(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	orgID := uuid.New()

	mockRepo.On("CreateDocument", mock.Anything, orgID, mock.Anything).Return(errors.New("validation error: invalid foreign name"))

	service := NewEsfDocumentService(mockRepo, nil, logrus.New())
	_, err := service.CreateDocument(context.Background(), orgID, &models.EsfCreateDocumentRequest{
		ForeignName: "", // Invalid empty name
	})

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestEsfDocumentUpdateNonExistentDocument(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	orgID := uuid.New()

	mockRepo.On("UpdateDocument", mock.Anything, orgID, mock.Anything).Return(errors.New("document not found"))

	service := NewEsfDocumentService(mockRepo, nil, logrus.New())
	err := service.UpdateDocument(context.Background(), orgID, &models.EsfEditDocumentRequest{
		ID: uuid.New(),
	})

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestEsfDocumentRepositoryFailure(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	orgID := uuid.New()

	mockRepo.On("GetAllDocuments", mock.Anything, orgID).Return(nil, errors.New("database connection failed"))

	service := NewEsfDocumentService(mockRepo, nil, logrus.New())
	result, err := service.GetAllDocuments(context.Background(), orgID)

	assert.Error(t, err)
	assert.Nil(t, result)
	// Service may wrap or transform the error message
	assert.NotEmpty(t, err.Error())
}

// ========== UserService Tests (Placeholder) ==========
// UserService tests skipped - requires proper UserRepository mock with all methods
// including Create, GetByID, GetByUsername, GetByEmail, Update, Delete

// ========== RoleService Tests (Placeholder) ==========
// RoleService tests skipped - requires RBAC repository interface that doesn't exist yet
// Should be implemented once RBAC layer is fully defined

// ========== Benchmarks ==========

func BenchmarkEsfDocumentGetAll(b *testing.B) {
	mockRepo := new(MockDocumentRepository)
	orgID := uuid.New()

	// Create 100 documents for realistic benchmark
	docs := make([]entity.EsfDocument, 100)
	for i := 0; i < 100; i++ {
		docs[i] = entity.EsfDocument{
			ID:          uuid.New(),
			ForeignName: "Doc " + string(rune(i%26+65)),
		}
	}

	mockRepo.On("GetAllDocuments", mock.Anything, orgID).Return(docs, nil)
	service := NewEsfDocumentService(mockRepo, nil, logrus.New())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.GetAllDocuments(context.Background(), orgID)
	}
}

func BenchmarkEsfDocumentGetByID(b *testing.B) {
	mockRepo := new(MockDocumentRepository)
	orgID := uuid.New()
	docID := uuid.New()
	doc := &entity.EsfDocument{ID: docID, ForeignName: "Test Document"}

	mockRepo.On("GetDocumentByID", mock.Anything, orgID, docID).Return(doc, nil)
	service := NewEsfDocumentService(mockRepo, nil, logrus.New())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.GetDocumentByID(context.Background(), orgID, docID)
	}
}

func BenchmarkEsfDocumentCreate(b *testing.B) {
	mockRepo := new(MockDocumentRepository)
	orgID := uuid.New()

	mockRepo.On("CreateDocument", mock.Anything, orgID, mock.Anything).Return(nil)
	service := NewEsfDocumentService(mockRepo, nil, logrus.New())

	req := &models.EsfCreateDocumentRequest{
		ForeignName: "Benchmark Document",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.CreateDocument(context.Background(), orgID, req)
	}
}

func BenchmarkEsfDocumentUpdate(b *testing.B) {
	mockRepo := new(MockDocumentRepository)
	orgID := uuid.New()

	mockRepo.On("UpdateDocument", mock.Anything, orgID, mock.Anything).Return(nil)
	service := NewEsfDocumentService(mockRepo, nil, logrus.New())

	req := &models.EsfEditDocumentRequest{
		ID: uuid.New(),
		EsfCreateDocumentRequest: models.EsfCreateDocumentRequest{
			ForeignName: "Updated Benchmark Doc",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.UpdateDocument(context.Background(), orgID, req)
	}
}

func BenchmarkEsfDocumentDelete(b *testing.B) {
	mockRepo := new(MockDocumentRepository)
	orgID := uuid.New()
	docID := uuid.New()

	mockRepo.On("DeleteDocument", mock.Anything, orgID, docID).Return(nil)
	service := NewEsfDocumentService(mockRepo, nil, logrus.New())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.DeleteDocument(context.Background(), orgID, docID)
	}
}

func BenchmarkEsfDocumentPaginated(b *testing.B) {
	mockRepo := new(MockDocumentRepository)
	orgID := uuid.New()

	docs := make([]entity.EsfDocument, 50)
	for i := 0; i < 50; i++ {
		docs[i] = entity.EsfDocument{
			ID:          uuid.New(),
			ForeignName: "Doc " + string(rune(i%26+65)),
		}
	}

	params := pagination.PaginationParams{Page: 1, PageSize: 50}
	filters := pagination.DocumentFilterParams{}

	mockRepo.On("GetAllDocumentsPaginated", mock.Anything, orgID, params, filters).Return(docs, int64(500), nil)
	service := NewEsfDocumentService(mockRepo, nil, logrus.New())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.GetAllDocumentsPaginated(context.Background(), orgID, params, filters)
	}
}

func BenchmarkEsfDocumentConcurrentReadOps(b *testing.B) {
	mockRepo := new(MockDocumentRepository)
	orgID := uuid.New()
	docID := uuid.New()
	doc := &entity.EsfDocument{ID: docID, ForeignName: "Benchmark Doc"}

	mockRepo.On("GetDocumentByID", mock.Anything, orgID, docID).Return(doc, nil)
	service := NewEsfDocumentService(mockRepo, nil, logrus.New())

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			service.GetDocumentByID(context.Background(), orgID, docID)
		}
	})
}
