package utils

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/hekmon/transmissionrpc/v3"
	"github.com/sirupsen/logrus"
)

// Check and delete completed downloads
func CheckAndDeleteOldDownloads(client *transmissionrpc.Client, logger *logrus.Logger, dryRun bool, deleteAfter int, targetDownloadDir string, keepFilesDownloadDir string) {
	ctx := context.Background()
	torrents, err := client.TorrentGetAll(ctx)
	if err != nil {
		logger.Errorf("Failed to fetch torrents: %v", err)
		return
	}

	now := time.Now()

	hoursAgo := now.Add(-time.Duration(deleteAfter) * time.Hour)
	logger.Debugf("Delete torrents older than %s", hoursAgo.Format("2006-01-02 15:04:05"))

	for _, torrent := range torrents {
		logger.Debugf("Torrent: %s, Status: %s, Date: %s, Dir: %s\n", *torrent.Name, *torrent.Status, *torrent.DoneDate, *torrent.DownloadDir)
		// Validate torrent status
		if *torrent.Status == transmissionrpc.TorrentStatusStopped || *torrent.Status == transmissionrpc.TorrentStatusSeed {
			// Check if torrent was completed more than 2 hours ago and DoneDate is not 1970-01-01 (for uncompleted downloads)
			if torrent.DoneDate.Before(hoursAgo) && !torrent.DoneDate.Equal(time.Unix(0, 0)) {
				if strings.Contains(*torrent.DownloadDir, targetDownloadDir) {
					deleteTorrent(ctx, client, logger, &torrent, true, dryRun)
				} else if strings.Contains(*torrent.DownloadDir, keepFilesDownloadDir) {
					deleteTorrent(ctx, client, logger, &torrent, false, dryRun)
				}
			}
		}
	}
}

// deleteTorrent removes a torrent by ID and logs the result
func deleteTorrent(ctx context.Context, client *transmissionrpc.Client, logger *logrus.Logger, torrent *transmissionrpc.Torrent, deleteData bool, dryRun bool) {
	if deleteData {
		logger.Infof("The torrent %s is going to be deleted with files\n", *torrent.Name)
	} else {
		logger.Infof("The torrent %s is going to be deleted without files\n", *torrent.Name)
	}

	payload := transmissionrpc.TorrentRemovePayload{
		IDs:             []int64{*torrent.ID},
		DeleteLocalData: deleteData,
	}
	if !dryRun {
		err := client.TorrentRemove(ctx, payload)
		if err != nil {
			logger.Errorf("Failed to delete torrent %s: %v", *torrent.Name, err)
		} else {
			logger.Infof("Torrent %s deleted successfully", *torrent.Name)
		}
	}
}

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
