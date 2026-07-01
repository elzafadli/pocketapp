package service

import (
	"context"
	"testing"

	"userapp/internal/domain/pocket"
	mock_pocket "userapp/internal/domain/pocket/mocks"
	"userapp/internal/domain/shared/identity"

	"github.com/runsystemid/golog"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type PocketServiceSuite struct {
	suite.Suite
	mockCtrl  *gomock.Controller
	mockRepo  *mock_pocket.MockRepository
	pocketSvc *Pocket
	ctx       context.Context
}

func (suite *PocketServiceSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockRepo = mock_pocket.NewMockRepository(suite.mockCtrl)
	suite.pocketSvc = &Pocket{
		PocketRepo: suite.mockRepo,
	}
	suite.ctx = context.Background()
	golog.Load(golog.Config{})
}

func (suite *PocketServiceSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *PocketServiceSuite) TestCreate_Success() {
	request := &pocket.CreatePocketRequest{
		Title:       "React Performance Guide",
		URL:         "https://example.com/react-performance",
		Description: "A guide about React rendering optimization",
		ContentType: "article",
		Tags:        []string{"frontend", "react"},
	}

	suite.mockRepo.EXPECT().
		Create(suite.ctx, "public", gomock.Any()).
		DoAndReturn(func(ctx context.Context, schema string, item *pocket.PocketItem) (*pocket.PocketItem, error) {
			suite.Equal(request.Title, item.Title)
			suite.Equal(request.URL, *item.URL)
			suite.Equal(request.Description, *item.Description)
			suite.Equal(request.ContentType, item.ContentType)
			suite.Equal("unread", item.Status)
			suite.False(item.IsFavorite)
			return item, nil
		})

	result, err := suite.pocketSvc.Create(suite.ctx, "public", request)
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(request.Title, result.Title)
}

func (suite *PocketServiceSuite) TestUpdate_Success() {
	id := identity.NewID()
	existing := &pocket.PocketItem{
		ID:          id,
		Title:       "Old Title",
		ContentType: "article",
	}

	request := &pocket.UpdatePocketRequest{
		ID:          id,
		Title:       "Updated Title",
		ContentType: "article",
		URL:         "https://example.com/updated",
	}

	suite.mockRepo.EXPECT().GetByID(suite.ctx, "public", id).Return(existing, nil)
	suite.mockRepo.EXPECT().Update(suite.ctx, "public", gomock.Any()).Return(nil)

	result, err := suite.pocketSvc.Update(suite.ctx, "public", request)
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal("Updated Title", result.Title)
}

func (suite *PocketServiceSuite) TestUpdate_NotFound() {
	id := identity.NewID()
	request := &pocket.UpdatePocketRequest{
		ID:          id,
		Title:       "Updated Title",
		ContentType: "article",
	}

	suite.mockRepo.EXPECT().GetByID(suite.ctx, "public", id).Return(nil, pocket.ErrPocketNotFound)

	result, err := suite.pocketSvc.Update(suite.ctx, "public", request)
	suite.ErrorIs(err, pocket.ErrPocketNotFound)
	suite.Nil(result)
}

func (suite *PocketServiceSuite) TestFind_Success() {
	id := identity.NewID()
	expected := &pocket.PocketItem{ID: id, Title: "Some Title"}

	suite.mockRepo.EXPECT().GetByID(suite.ctx, "public", id).Return(expected, nil)

	result, err := suite.pocketSvc.Find(suite.ctx, "public", id)
	suite.NoError(err)
	suite.Equal(expected, result)
}

func (suite *PocketServiceSuite) TestDelete_Success() {
	id := identity.NewID()

	suite.mockRepo.EXPECT().Delete(suite.ctx, "public", id).Return(nil)

	err := suite.pocketSvc.Delete(suite.ctx, "public", id)
	suite.NoError(err)
}

func (suite *PocketServiceSuite) TestList_Success() {
	query := &pocket.PocketListQuery{
		Search: "test",
		Page:   1,
		Limit:  10,
		Sort:   "createdAt:desc",
	}
	expectedFilter := query.ToFilter()
	expectedList := []*pocket.PocketItem{
		{Title: "Result 1"},
	}

	suite.mockRepo.EXPECT().List(suite.ctx, "public", expectedFilter).Return(expectedList, nil)
	suite.mockRepo.EXPECT().Count(suite.ctx, "public", expectedFilter).Return(uint64(1), nil)

	list, count, err := suite.pocketSvc.List(suite.ctx, "public", query)
	suite.NoError(err)
	suite.Equal(expectedList, list)
	suite.Equal(uint64(1), count)
}

func (suite *PocketServiceSuite) TestUpdateStatus_Success() {
	id := identity.NewID()
	existing := &pocket.PocketItem{ID: id, Status: "unread"}

	suite.mockRepo.EXPECT().GetByID(suite.ctx, "public", id).Return(existing, nil)
	suite.mockRepo.EXPECT().Update(suite.ctx, "public", gomock.Any()).Return(nil)

	result, err := suite.pocketSvc.UpdateStatus(suite.ctx, "public", id, "read")
	suite.NoError(err)
	suite.Equal("read", result.Status)
}

func (suite *PocketServiceSuite) TestToggleFavorite_Success() {
	id := identity.NewID()
	existing := &pocket.PocketItem{ID: id, IsFavorite: false}

	suite.mockRepo.EXPECT().GetByID(suite.ctx, "public", id).Return(existing, nil)
	suite.mockRepo.EXPECT().Update(suite.ctx, "public", gomock.Any()).Return(nil)

	result, err := suite.pocketSvc.ToggleFavorite(suite.ctx, "public", id, true)
	suite.NoError(err)
	suite.True(result.IsFavorite)
}

func (suite *PocketServiceSuite) TestGetSummary_Success() {
	expectedSummary := &pocket.PocketSummary{
		TotalItems:    150,
		UnreadItems:   50,
		ReadingItems:  20,
		ReadItems:     70,
		ArchivedItems: 10,
		FavoriteItems: 35,
	}

	suite.mockRepo.EXPECT().GetSummary(suite.ctx, "public").Return(expectedSummary, nil)

	result, err := suite.pocketSvc.GetSummary(suite.ctx, "public")
	suite.NoError(err)
	suite.Equal(expectedSummary, result)
}

func TestPocketService(t *testing.T) {
	suite.Run(t, new(PocketServiceSuite))
}
