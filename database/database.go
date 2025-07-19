package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bushubdegefu/m-playground/configs"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

var (
	MongoClient *mongo.Client
)

// LoggerFile returns an *os.File for the app's Mongo logs
func LoggerFile(appName string) (*os.File, error) {
	logFileName := fmt.Sprintf("%s_mongo.log", appName)
	logFile, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("error opening log file: %w", err)
	}
	return logFile, nil
}

// ReturnMongoClient initializes and returns a Mongo client with proper logging, tracing, and pooling.
func ReturnMongoClient(appName string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := configs.AppConfig.Get(fmt.Sprintf("%s_MONGODB_URI", strings.ToUpper(appName)))
	if mongoURI == "" {
		return nil, fmt.Errorf("mongo URI is empty for app: %s", appName)
	}

	// Log file setup
	logFile, err := LoggerFile(appName)
	if err != nil {
		return nil, err
	}
	logger := log.New(logFile, "[MongoDB] ", log.LstdFlags|log.Lshortfile)

	// Setup command logging
	otelMonitor := otelmongo.NewMonitor()

	cmdMonitor := &event.CommandMonitor{
		Started: func(ctx context.Context, evt *event.CommandStartedEvent) {
			logger.Printf("Command Started: %s %v", evt.CommandName, evt.Command)
			if otelMonitor.Started != nil {
				otelMonitor.Started(ctx, evt)
			}
		},
		Succeeded: func(ctx context.Context, evt *event.CommandSucceededEvent) {
			logger.Printf("Command Succeeded: %s Duration: %v", evt.CommandName, evt.Duration)
			if otelMonitor.Succeeded != nil {
				otelMonitor.Succeeded(ctx, evt)
			}
		},
		Failed: func(ctx context.Context, evt *event.CommandFailedEvent) {
			logger.Printf("Command Failed: %s Duration: %v Err: %v", evt.CommandName, evt.Duration, evt.Failure)
			if otelMonitor.Failed != nil {
				otelMonitor.Failed(ctx, evt)
			}
		},
	}

	// MongoDB connection options
	clientOpts := options.Client().
		ApplyURI(mongoURI).
		SetConnectTimeout(10 * time.Second).
		SetMaxPoolSize(100).
		SetMinPoolSize(10).
		SetMaxConnIdleTime(30 * time.Minute).
		SetMonitor(cmdMonitor) // log commands
		// SetMonitor(otelmongo.NewMonitor()) // OpenTelemetry tracing

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		logger.Printf("MongoDB connection error: %v", err)
		return nil, fmt.Errorf("mongo connect failed: %w", err)
	}

	MongoClient = client
	return client, nil
}
