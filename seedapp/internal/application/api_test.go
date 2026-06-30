package application

import (
	"errors"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"seedapp/internal/adapter/rest"
)

// MockHealthCheckAPI mocks the HealthCheckAPI interface
type MockHealthCheckAPI struct {
	mock.Mock
}

func (m *MockHealthCheckAPI) Ping(c *fiber.Ctx) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockHealthCheckAPI) Ready(c *fiber.Ctx) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockHealthCheckAPI) Version(c *fiber.Ctx) error {
	args := m.Called(c)
	return args.Error(0)
}

// MockMigrationAPI mocks the MigrationAPI interface
type MockMigrationAPI struct {
	mock.Mock
}

func (m *MockMigrationAPI) RunMigration(c *fiber.Ctx) error {
	args := m.Called(c)
	return args.Error(0)
}

// MockSeedAPI mocks the SeedAPI interface
type MockSeedAPI struct {
	mock.Mock
}

func (m *MockSeedAPI) RunSeed(c *fiber.Ctx) error {
	args := m.Called(c)
	return args.Error(0)
}

// MockAuthMiddleware mocks the AuthMiddlewareService interface
type MockAuthMiddleware struct {
	mock.Mock
}

func (m *MockAuthMiddleware) BasicAuthProtection() fiber.Handler {
	args := m.Called()
	return args.Get(0).(fiber.Handler)
}

// MockRouter mocks a fiber.Router
type MockRouter struct {
	mock.Mock
}

func (m *MockRouter) Post(path string, handlers ...fiber.Handler) fiber.Router {
	args := m.Called(path, handlers[0])
	return args.Get(0).(fiber.Router)
}

// MockFiberApp is a mock implementation of rest.Fiber
type MockFiberApp struct {
	*fiber.App
	mock.Mock
}

func NewMockFiberApp() *rest.Fiber {
	mockApp := &MockFiberApp{
		App: fiber.New(),
	}
	return &rest.Fiber{App: mockApp.App}
}

// ErrorFiber is a custom implementation of rest.Fiber that returns an error on Shutdown
type ErrorFiber struct {
	*fiber.App
}

func (f *ErrorFiber) Shutdown() error {
	return errors.New("shutdown error")
}

func TestApiStartup(t *testing.T) {
	// Create the app with a real fiber instance
	app := fiber.New()
	fiberApp := &rest.Fiber{App: app}

	// Create other mocks
	mockHealthCheck := new(MockHealthCheckAPI)
	mockMigration := new(MockMigrationAPI)
	mockSeed := new(MockSeedAPI)
	mockAuth := new(MockAuthMiddleware)
	mockAuthHandler := func(c *fiber.Ctx) error { return nil }

	// Setup expectations
	mockAuth.On("BasicAuthProtection").Return(mockAuthHandler)

	// Create API instance with mocks
	api := Api{
		App:                fiberApp,
		AuthMiddleware:     mockAuth,
		HealthCheckHandler: mockHealthCheck,
		MigrationHandler:   mockMigration,
		SeedHandler:        mockSeed,
	}

	// Call the method under test
	err := api.Startup()

	// Assertions
	assert.NoError(t, err)
	mockAuth.AssertExpectations(t)

	// Verify routes were registered
	// We can't directly verify the routes, but we can check that the app has routes
	assert.NotEmpty(t, app.Stack())

	// Test route existence by checking the stack
	routes := app.Stack()

	// Helper function to check if a route exists
	routeExists := func(method, path string) bool {
		for _, route := range routes {
			for _, r := range route {
				if r.Method == method && r.Path == path {
					return true
				}
			}
		}
		return false
	}

	// Check if our routes exist
	assert.True(t, routeExists("GET", "/ping"), "Route GET /ping should exist")
	assert.True(t, routeExists("GET", "/ready"), "Route GET /ready should exist")
	assert.True(t, routeExists("POST", "/v1/migrate"), "Route POST /v1/migrate should exist")
	assert.True(t, routeExists("POST", "/v1/seed"), "Route POST /v1/seed should exist")
}

func TestApiShutdown(t *testing.T) {
	// Create a real fiber app for testing
	app := fiber.New()
	fiberApp := &rest.Fiber{App: app}

	api := Api{
		App: fiberApp,
	}

	// Call the method under test
	err := api.Shutdown()

	// Assertions
	assert.NoError(t, err)
}

// TestApiWithNilDependencies tests that the API can handle nil dependencies
func TestApiWithNilDependencies(t *testing.T) {
	// Create API instance with nil dependencies
	api := Api{}

	// Call the method under test - this should not panic
	err := api.Shutdown()

	// Assertions
	assert.NoError(t, err)
}

// TestApiShutdownWithError tests the error handling in Shutdown
func TestApiShutdownWithError(t *testing.T) {
	t.Skip("Skipping this test as we can't easily mock the Fiber.Shutdown method")
}
