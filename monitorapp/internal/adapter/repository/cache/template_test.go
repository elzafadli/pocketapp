package cache

import (
	"context"
	"testing"

	"monitorapp/internal/domain/shared/entity"
	"monitorapp/internal/domain/shared/identity"
	"monitorapp/internal/domain/template"
	mock_template "monitorapp/internal/domain/template/mocks"

	"github.com/runsystemid/gocache"
	mock_gocache "github.com/runsystemid/gocache/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type TemplateRepositorySuite struct {
	suite.Suite
	mockCtrl  *gomock.Controller
	mockCache *mock_gocache.MockService
	mockRepo  *mock_template.MockRepository
	repo      *TemplateRepository
}

func (suite *TemplateRepositorySuite) SetupTest() {
	// Set up any necessary test fixtures
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockCache = mock_gocache.NewMockService(suite.mockCtrl)
	suite.mockRepo = mock_template.NewMockRepository(suite.mockCtrl)
	suite.repo = &TemplateRepository{
		Cache: suite.mockCache,
		Next:  suite.mockRepo,
	}
}

func (suite *TemplateRepositorySuite) TearDownTest() {
	// Tear down any test fixtures
	suite.mockCtrl.Finish()
}

func (suite *TemplateRepositorySuite) TestCreate() {
	ctx := context.TODO()
	data := &template.Template{}

	suite.mockRepo.EXPECT().Create(ctx, data).Return(data, nil)

	result, err := suite.repo.Create(ctx, data)

	suite.NoError(err)
	suite.NotNil(result)
}

func (suite *TemplateRepositorySuite) TestUpdate() {
	ctx := context.TODO()
	data := &template.Template{
		Entity: entity.NewEntity(),
	}

	suite.mockCache.EXPECT().Delete(ctx, gomock.Any())
	suite.mockRepo.EXPECT().Update(ctx, data).Return(nil)

	err := suite.repo.Update(ctx, data)

	suite.NoError(err)
}

func (suite *TemplateRepositorySuite) TestGetByID_CacheAvailable() {
	ctx := context.TODO()
	id := identity.ID{}

	data := &template.Template{}
	cached := &template.Template{
		Entity:    entity.Entity{ID: id},
		Name:      "test",
		Category:  "test",
		Published: true,
	}

	suite.mockCache.EXPECT().Get(ctx, gomock.Any(), &data).SetArg(2, cached).Return(nil)

	result, err := suite.repo.GetByID(ctx, id)

	suite.NoError(err)
	suite.Equal(cached, result)
}

func (suite *TemplateRepositorySuite) TestGetByID_ErrCache() {
	ctx := context.TODO()
	id := identity.ID{}

	suite.mockCache.EXPECT().Get(ctx, gomock.Any(), gomock.Any()).Return(assert.AnError)

	result, err := suite.repo.GetByID(ctx, id)

	suite.Error(err)
	suite.Nil(result)
}

func (suite *TemplateRepositorySuite) TestGetByID_ErrRepo() {
	ctx := context.TODO()
	id := identity.ID{}
	// data := &template.Template{
	// 	Entity:    entity.Entity{ID: id},
	// 	Name:      "test",
	// 	Category:  "test",
	// 	Published: true,
	// }

	suite.mockCache.EXPECT().Get(ctx, gomock.Any(), gomock.Any()).Return(gocache.ErrNil)
	suite.mockRepo.EXPECT().GetByID(ctx, id).Return(nil, assert.AnError)

	result, err := suite.repo.GetByID(ctx, id)

	suite.Error(err)
	suite.Nil(result)
}

func (suite *TemplateRepositorySuite) TestGetByID_NoErr() {
	ctx := context.TODO()
	id := identity.ID{}

	suite.mockCache.EXPECT().Get(ctx, gomock.Any(), gomock.Any()).Return(gocache.ErrNil)
	suite.mockRepo.EXPECT().GetByID(ctx, id).Return(&template.Template{}, nil)
	suite.mockCache.EXPECT().Put(ctx, gomock.Any(), gomock.Any(), gomock.Any())

	result, err := suite.repo.GetByID(ctx, id)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
}

func TestTemplateRepositorySuite(t *testing.T) {
	suite.Run(t, new(TemplateRepositorySuite))
}
