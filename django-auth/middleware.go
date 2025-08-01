package django_auth

import (
	"fmt"
	"github.com/bushubdegefu/m-playground/database"
	"github.com/bushubdegefu/m-playground/observe"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"strings"
)

func otelechospanstarter(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		routeName := ctx.Path() + "_" + strings.ToLower(ctx.Request().Method)
		tracer, span := observe.EchoAppSpanner(ctx, fmt.Sprintf("%v-root", routeName))
		ctx.Set("tracer", &observe.RouteTracer{Tracer: tracer, Span: span})

		// Process request
		err := next(ctx)
		if err != nil {
			return err
		}

		span.SetAttributes(attribute.String("response", fmt.Sprintf("%v", ctx.Response().Status)))
		span.End()
		return nil
	}
}

func dbsessioninjection(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		client, err := database.ReturnMongoClient("django_auth")
		if err != nil {
			return err
		}

		ctx.Set("db", client)

		nerr := next(ctx)
		if nerr != nil {
			return nerr
		}
		return nil
	}
}

// Custom Middlewares can be added here specfic to the app
