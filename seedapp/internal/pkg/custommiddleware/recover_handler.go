package custommiddleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/runsystemid/golog"
)

// gologErrorFunc is a variable that holds the golog.Error function
// This allows us to replace it in tests
var gologErrorFunc = golog.Error

// RecoverHandler handles panics in Fiber applications
func RecoverHandler(c *fiber.Ctx, e interface{}) {
	err, ok := e.(error)
	if !ok {
		err = fmt.Errorf("%v", e)
	}

	gologErrorFunc(c.Context(), err.Error(), err)
}
