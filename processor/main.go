package main

import (
	"github.com/rs-anantmishra/streamsphere/utils/processor/database"
	"github.com/rs-anantmishra/streamsphere/utils/processor/extractor"
	"github.com/rs-anantmishra/streamsphere/utils/processor/requests"
)

func main() {

	//get cli params for scheduled tasks here -- requestId to process as cli param

	//request, scheduler and, process will have repos but no separate service

	// connect to database
	database.ConnectDB()
	defer database.CloseDB()

	//Instantiate extractor service - prerequisites
	extrRepo := extractor.NewRepository(database.DB)
	reqRepo := requests.NewRepository(database.DB)
	extrDownloads := extractor.NewDownload()

	//Instantiate extractor service
	extr := extractor.NewDownloadService(extrRepo, reqRepo, extrDownloads)

	//run-process
	extr.ProcessRequests()
}
