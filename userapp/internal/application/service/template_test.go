package service

import (
	"context"
	"errors"
	"testing"

	"userapp/internal/domain/shared/entity"
	"userapp/internal/domain/shared/identity"
	"userapp/internal/domain/template"
	mock_template "userapp/internal/domain/template/mocks"

	"github.com/runsystemid/golog"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"gopkg.in/guregu/null.v4"
)

type TemplateServiceSuite struct {
	suite.Suite
	mockCtrl    *gomock.Controller
	mockRepo    *mock_template.MockRepository
	templateSvc *Template
	ctx         context.Context
}

func (suite *TemplateServiceSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockRepo = mock_template.NewMockRepository(suite.mockCtrl)
	suite.templateSvc = &Template{
		TemplateRepo: suite.mockRepo,
	}
	suite.ctx = context.Background()
	golog.Load(golog.Config{})
}

func (suite *TemplateServiceSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *TemplateServiceSuite) TestCreate_Success() {
	// Arrange
	request := &template.CreateTemplateRequest{
		Name:      "Test Template",
		Category:  "Test Category",
		Published: true,
	}

	expectedTemplate := &template.Template{
		Entity:    entity.NewEntity(),
		Name:      request.Name,
		Category:  request.Category,
		Published: request.Published,
	}

	suite.mockRepo.EXPECT().
		Create(suite.ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, t *template.Template) (*template.Template, error) {
			// Verify the template was created with correct values
			suite.Equal(request.Name, t.Name)
			suite.Equal(request.Category, t.Category)
			suite.Equal(request.Published, t.Published)
			suite.NotEmpty(t.Entity.ID)
			// Return template with ID set
			expectedTemplate.Entity.ID = t.Entity.ID
			expectedTemplate.Entity.CreatedAt = t.Entity.CreatedAt
			expectedTemplate.Entity.UpdatedAt = t.Entity.UpdatedAt
			return expectedTemplate, nil
		})

	// Act
	result, err := suite.templateSvc.Create(suite.ctx, request)

	// Assert
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(request.Name, result.Name)
	suite.Equal(request.Category, result.Category)
	suite.Equal(request.Published, result.Published)
	suite.NotEmpty(result.Entity.ID)
}

func (suite *TemplateServiceSuite) TestCreate_RepositoryError() {
	// Arrange
	request := &template.CreateTemplateRequest{
		Name:      "Test Template",
		Category:  "Test Category",
		Published: false,
	}

	repoError := errors.New("repository error")
	suite.mockRepo.EXPECT().
		Create(suite.ctx, gomock.Any()).
		Return(nil, repoError)

	// Act
	result, err := suite.templateSvc.Create(suite.ctx, request)

	// Assert
	suite.Error(err)
	suite.Equal(repoError, err)
	suite.Nil(result)
}

func (suite *TemplateServiceSuite) TestCreate_WithAllFields() {
	// Arrange
	request := &template.CreateTemplateRequest{
		Name:      "Complete Template",
		Category:  "Complete Category",
		Published: true,
	}

	expectedTemplate := &template.Template{
		Entity:    entity.NewEntity(),
		Name:      request.Name,
		Category:  request.Category,
		Published: request.Published,
	}

	suite.mockRepo.EXPECT().
		Create(suite.ctx, gomock.Any()).
		Return(expectedTemplate, nil)

	// Act
	result, err := suite.templateSvc.Create(suite.ctx, request)

	// Assert
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal("Complete Template", result.Name)
	suite.Equal("Complete Category", result.Category)
	suite.True(result.Published)
}

func (suite *TemplateServiceSuite) TestUpdate_Success_AllFields() {
	// Arrange
	templateID := identity.NewID()
	existingEntity := entity.NewEntity()
	existingEntity.ID = templateID
	existingTemplate := &template.Template{
		Entity:    existingEntity,
		Name:      "Old Name",
		Category:  "Old Category",
		Published: false,
	}

	request := &template.UpdateTemplateRequest{
		ID:        templateID,
		Name:      "New Name",
		Category:  "New Category",
		Published: null.BoolFrom(true),
	}

	suite.mockRepo.EXPECT().
		GetByID(suite.ctx, templateID).
		Return(existingTemplate, nil)

	suite.mockRepo.EXPECT().
		Update(suite.ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, t *template.Template) error {
			suite.Equal("New Name", t.Name)
			suite.Equal("New Category", t.Category)
			suite.True(t.Published)
			suite.NotZero(t.Entity.UpdatedAt)
			return nil
		})

	// Act
	result, err := suite.templateSvc.Update(suite.ctx, request)

	// Assert
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal("New Name", result.Name)
	suite.Equal("New Category", result.Category)
	suite.True(result.Published)
	suite.NotZero(result.UpdatedAt)
}

func (suite *TemplateServiceSuite) TestUpdate_Success_PartialFields() {
	// Arrange
	templateID := identity.NewID()
	existingEntity := entity.NewEntity()
	existingEntity.ID = templateID
	existingTemplate := &template.Template{
		Entity:    existingEntity,
		Name:      "Original Name",
		Category:  "Original Category",
		Published: true,
	}

	request := &template.UpdateTemplateRequest{
		ID:        templateID,
		Name:      "Updated Name",
		Category:  "",          // Empty, should not update
		Published: null.Bool{}, // Zero value, should not update
	}

	suite.mockRepo.EXPECT().
		GetByID(suite.ctx, templateID).
		Return(existingTemplate, nil)

	suite.mockRepo.EXPECT().
		Update(suite.ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, t *template.Template) error {
			suite.Equal("Updated Name", t.Name)
			suite.Equal("Original Category", t.Category) // Should remain unchanged
			suite.True(t.Published)                      // Should remain unchanged
			return nil
		})

	// Act
	result, err := suite.templateSvc.Update(suite.ctx, request)

	// Assert
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal("Updated Name", result.Name)
	suite.Equal("Original Category", result.Category)
	suite.True(result.Published)
}

func (suite *TemplateServiceSuite) TestUpdate_Success_OnlyName() {
	// Arrange
	templateID := identity.NewID()
	existingEntity := entity.NewEntity()
	existingEntity.ID = templateID
	existingTemplate := &template.Template{
		Entity:    existingEntity,
		Name:      "Original Name",
		Category:  "Original Category",
		Published: false,
	}

	request := &template.UpdateTemplateRequest{
		ID:        templateID,
		Name:      "New Name Only",
		Category:  "",
		Published: null.Bool{},
	}

	suite.mockRepo.EXPECT().
		GetByID(suite.ctx, templateID).
		Return(existingTemplate, nil)

	suite.mockRepo.EXPECT().
		Update(suite.ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, t *template.Template) error {
			suite.Equal("New Name Only", t.Name)
			return nil
		})

	// Act
	result, err := suite.templateSvc.Update(suite.ctx, request)

	// Assert
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal("New Name Only", result.Name)
}

func (suite *TemplateServiceSuite) TestUpdate_Success_OnlyCategory() {
	// Arrange
	templateID := identity.NewID()
	existingEntity := entity.NewEntity()
	existingEntity.ID = templateID
	existingTemplate := &template.Template{
		Entity:    existingEntity,
		Name:      "Original Name",
		Category:  "Original Category",
		Published: true,
	}

	request := &template.UpdateTemplateRequest{
		ID:        templateID,
		Name:      "",
		Category:  "New Category Only",
		Published: null.Bool{},
	}

	suite.mockRepo.EXPECT().
		GetByID(suite.ctx, templateID).
		Return(existingTemplate, nil)

	suite.mockRepo.EXPECT().
		Update(suite.ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, t *template.Template) error {
			suite.Equal("New Category Only", t.Category)
			return nil
		})

	// Act
	result, err := suite.templateSvc.Update(suite.ctx, request)

	// Assert
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal("New Category Only", result.Category)
}

func (suite *TemplateServiceSuite) TestUpdate_Success_OnlyPublished() {
	// Arrange
	templateID := identity.NewID()
	existingEntity := entity.NewEntity()
	existingEntity.ID = templateID
	existingTemplate := &template.Template{
		Entity:    existingEntity,
		Name:      "Original Name",
		Category:  "Original Category",
		Published: false,
	}

	request := &template.UpdateTemplateRequest{
		ID:        templateID,
		Name:      "",
		Category:  "",
		Published: null.BoolFrom(true),
	}

	suite.mockRepo.EXPECT().
		GetByID(suite.ctx, templateID).
		Return(existingTemplate, nil)

	suite.mockRepo.EXPECT().
		Update(suite.ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, t *template.Template) error {
			suite.True(t.Published)
			return nil
		})

	// Act
	result, err := suite.templateSvc.Update(suite.ctx, request)

	// Assert
	suite.NoError(err)
	suite.NotNil(result)
	suite.True(result.Published)
}

func (suite *TemplateServiceSuite) TestUpdate_GetByIDError() {
	// Arrange
	templateID := identity.NewID()
	repoError := errors.New("template not found")

	request := &template.UpdateTemplateRequest{
		ID:        templateID,
		Name:      "New Name",
		Category:  "New Category",
		Published: null.BoolFrom(true),
	}

	suite.mockRepo.EXPECT().
		GetByID(suite.ctx, templateID).
		Return(nil, repoError)

	// Act
	result, err := suite.templateSvc.Update(suite.ctx, request)

	// Assert
	suite.Error(err)
	suite.Equal(repoError, err)
	suite.Nil(result)
}

func (suite *TemplateServiceSuite) TestUpdate_UpdateError() {
	// Arrange
	templateID := identity.NewID()
	existingEntity := entity.NewEntity()
	existingEntity.ID = templateID
	existingTemplate := &template.Template{
		Entity:    existingEntity,
		Name:      "Original Name",
		Category:  "Original Category",
		Published: false,
	}

	request := &template.UpdateTemplateRequest{
		ID:        templateID,
		Name:      "New Name",
		Category:  "New Category",
		Published: null.BoolFrom(true),
	}

	updateError := errors.New("update failed")

	suite.mockRepo.EXPECT().
		GetByID(suite.ctx, templateID).
		Return(existingTemplate, nil)

	suite.mockRepo.EXPECT().
		Update(suite.ctx, gomock.Any()).
		Return(updateError)

	// Act
	result, err := suite.templateSvc.Update(suite.ctx, request)

	// Assert
	suite.Error(err)
	suite.Equal(updateError, err)
	suite.Nil(result)
}

func (suite *TemplateServiceSuite) TestUpdate_PublishedFalse() {
	// Arrange
	templateID := identity.NewID()
	existingEntity := entity.NewEntity()
	existingEntity.ID = templateID
	existingTemplate := &template.Template{
		Entity:    existingEntity,
		Name:      "Original Name",
		Category:  "Original Category",
		Published: true,
	}

	request := &template.UpdateTemplateRequest{
		ID:        templateID,
		Name:      "",
		Category:  "",
		Published: null.BoolFrom(false),
	}

	suite.mockRepo.EXPECT().
		GetByID(suite.ctx, templateID).
		Return(existingTemplate, nil)

	suite.mockRepo.EXPECT().
		Update(suite.ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, t *template.Template) error {
			suite.False(t.Published)
			return nil
		})

	// Act
	result, err := suite.templateSvc.Update(suite.ctx, request)

	// Assert
	suite.NoError(err)
	suite.NotNil(result)
	suite.False(result.Published)
}

func (suite *TemplateServiceSuite) TestFind_Success() {
	// Arrange
	templateID := identity.NewID()
	expectedEntity := entity.NewEntity()
	expectedEntity.ID = templateID
	expectedTemplate := &template.Template{
		Entity:    expectedEntity,
		Name:      "Test Template",
		Category:  "Test Category",
		Published: true,
	}

	suite.mockRepo.EXPECT().
		GetByID(suite.ctx, templateID).
		Return(expectedTemplate, nil)

	// Act
	result, err := suite.templateSvc.Find(suite.ctx, templateID)

	// Assert
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(expectedTemplate.Entity.ID, result.Entity.ID)
	suite.Equal(expectedTemplate.Name, result.Name)
	suite.Equal(expectedTemplate.Category, result.Category)
	suite.Equal(expectedTemplate.Published, result.Published)
}

func (suite *TemplateServiceSuite) TestFind_RepositoryError() {
	// Arrange
	templateID := identity.NewID()
	repoError := errors.New("template not found")

	suite.mockRepo.EXPECT().
		GetByID(suite.ctx, templateID).
		Return(nil, repoError)

	// Act
	result, err := suite.templateSvc.Find(suite.ctx, templateID)

	// Assert
	suite.Error(err)
	suite.Equal(repoError, err)
	suite.Nil(result)
}

func (suite *TemplateServiceSuite) TestFind_WithDifferentID() {
	// Arrange
	templateID1 := identity.NewID()
	templateID2 := identity.NewID()
	expectedEntity := entity.NewEntity()
	expectedEntity.ID = templateID2
	expectedTemplate := &template.Template{
		Entity:    expectedEntity,
		Name:      "Another Template",
		Category:  "Another Category",
		Published: false,
	}

	suite.mockRepo.EXPECT().
		GetByID(suite.ctx, templateID1).
		Return(expectedTemplate, nil)

	// Act
	result, err := suite.templateSvc.Find(suite.ctx, templateID1)

	// Assert
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(templateID2, result.Entity.ID)
}

func TestTemplateServiceSuite(t *testing.T) {
	suite.Run(t, new(TemplateServiceSuite))
}
