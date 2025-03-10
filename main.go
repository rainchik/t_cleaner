package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/hekmon/transmissionrpc/v3"
	"github.com/rainchik/t_cleaner/utils"
	"github.com/sirupsen/logrus"
)

var (
	targetDownloadDir    = utils.GetEnv("TARGET_DOWNLOAD_DIR", "/downloads")
	keepFilesDownloadDir = utils.GetEnv("TARGET_DOWNLOAD_DIR", "/media")
	username             = os.Getenv("TRANSMISSION_USERNAME")
	password             = os.Getenv("TRANSMISSION_PASSWORD")
	transmissionUrl      = os.Getenv("TRANSMISSION_URL")
	deleteAfter          = utils.GetEnv("DELETE_AFTER", "2")
	logLevel             = utils.GetEnv("LOG_LEVEL", "INFO")
	dryRun               = utils.GetEnv("DRY_RUN", "false")
	checkInterval        = utils.GetEnv("CHECK_INTERVAL", "1")
)

func main() {
	// Configure logger
	logger := logrus.New()
	level, _ := logrus.ParseLevel(logLevel)
	logger.SetLevel(level)

	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	dryRun, _ := strconv.ParseBool(dryRun)
	deleteAfter, err := strconv.Atoi(deleteAfter)
	if err != nil {
		logger.Errorf("Can't convert deleteAfter to int: %v", err)
		return
	}
	// Validate environment variables
	if username == "" || password == "" || transmissionUrl == "" {
		logger.Fatal("Missing required environment variables (TRANSMISSION_URL, TRANSMISSION_USERNAME, TRANSMISSION_PASSWORD)")
	}

	// Connect to Transmission
	urlString := fmt.Sprintf("https://%s:%s@%s", username, password, transmissionUrl)
	endpoint, err := url.Parse(urlString)
	if err != nil {
		panic(err)
	}

	client, err := transmissionrpc.New(endpoint, nil)
	if err != nil {
		logger.Fatalf("Failed to connect to Transmission: %v", err)
	}

	logger.Info("Transmission cleaner service started")
	if dryRun {
		logger.Info("DryRun mode is on")
	}

	checkInterval, err := strconv.Atoi(checkInterval)
	if err != nil {
		logger.Errorf("Can't convert checkInterval to int: %v", err)
		return
	}
	// Periodic check loop
	for {
		logger.Debug("Reconcillation loop is starting")
		utils.CheckAndDeleteOldDownloads(client, logger, dryRun, deleteAfter, targetDownloadDir, keepFilesDownloadDir)
		time.Sleep(time.Duration(checkInterval) * time.Minute) // Check every some minutes
	}
}
