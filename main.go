package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v3"
	"github.com/ngdangkietswe/swe-go-common-shared/config"
	"github.com/ngdangkietswe/swe-integration-service/data/repository"
	"github.com/ngdangkietswe/swe-integration-service/grpc"
	"github.com/ngdangkietswe/swe-integration-service/logger"
	"go.uber.org/fx"
	grpcserver "google.golang.org/grpc"
)

type StravaConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

// stravaConfig is a global variable that holds the Strava configuration.
var stravaConfig StravaConfig

// logger is a global variable that holds the logger instance.
//var logger, _ = commonlogger.NewLogger(
//	"swe-integration-service",
//	"local",
//	"debug",
//	"",
//)

func main() {
	config.Init()

	app := fx.New(
		logger.Module,
		repository.Module,
		grpc.Module,
		fx.Invoke(func(*grpcserver.Server) {}),
	)
	app.Run()

	//
	//stravaConfig = StravaConfig{
	//	ClientID:     config.GetString("STRAVA_CLIENT_ID", ""),
	//	ClientSecret: config.GetString("STRAVA_CLIENT_SECRET", ""),
	//	RedirectURI:  config.GetString("STRAVA_REDIRECT_URI", ""),
	//}
	//
	//app := fiber.New()
	//
	//app.Get("/strava/auth", StravaAuthHandler)
	//app.Get("/strava/callback", StravaCallbackHandler)
	//app.Get("/strava/activities", GetActivitiesHandler)
	//
	//err := app.Listen(":8080")
	//if err != nil {
	//	logger.Error("Failed to start the server", zap.String("error", err.Error()))
	//	return
	//}
}

// StravaAuthHandler - Redirect to Strava for authentication
func StravaAuthHandler(ctx fiber.Ctx) error {
	authURL := fmt.Sprintf(
		"https://www.strava.com/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=read,activity:read_all",
		stravaConfig.ClientID,
		stravaConfig.RedirectURI,
	)
	//logger.Info("Redirecting to Strava for authentication", zap.String("authURL", authURL))

	ctx.Response().Header.Set("Location", authURL)
	ctx.Status(fiber.StatusTemporaryRedirect)
	return nil
}

// StravaCallbackHandler - Handle the callback from Strava and exchange code for access token
func StravaCallbackHandler(c fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).SendString("No code in request")
	}

	client := resty.New()
	response := map[string]interface{}{}

	// Exchange authorization code for access token
	_, err := client.R().
		SetFormData(map[string]string{
			"client_id":     stravaConfig.ClientID,
			"client_secret": stravaConfig.ClientSecret,
			"code":          code,
			"grant_type":    "authorization_code",
		}).
		SetResult(&response).
		Post("https://www.strava.com/oauth/token")

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error exchanging token: " + err.Error())
	}

	// Retrieve and display the access token
	token := response["access_token"].(string)
	//refreshToken := response["refresh_token"].(string)
	//athlete := response["athlete"].(map[string]interface{})
	//
	//logger.Info("Received access token from Strava", zap.String("accessToken", token))
	//logger.Info("Received refresh token from Strava", zap.String("refreshToken", refreshToken))
	//logger.Info("Received athlete info from Strava", zap.Any("athlete", athlete))

	return c.SendString("Authentication successful! Access Token: " + token)
}

// GetActivitiesHandler - Fetch user activities using the Strava API
func GetActivitiesHandler(c fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
	}

	client := resty.New()
	var activities []map[string]interface{}

	// Call Strava API to fetch activities
	_, err := client.R().
		SetAuthToken(token).
		SetResult(&activities).
		Get("https://www.strava.com/api/v3/athlete/activities")

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error fetching activities: " + err.Error())
	}

	// Return activities as JSON
	return c.JSON(activities)
}
