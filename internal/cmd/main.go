package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/anilvpatel21/snapurl/internal/ports"
	"github.com/anilvpatel21/snapurl/pkg/csv"
	"github.com/anilvpatel21/snapurl/pkg/downloader"
	"github.com/anilvpatel21/snapurl/pkg/persister"
)

var totalURLs, successFetch, errorFetch float64
var totalDuration float64

func main() {
	// Set up logging
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// Define flags
	filePath := flag.String("filePath", "", "Path to the input file")
	maxDownloadConcurrency := flag.Int("maxDownloadConcurrency", 50, "Maximum number of concurrent downloads")

	// Parse the flags
	flag.Parse()

	// Validate flags
	if *filePath == "" {
		log.Fatal("Argument Missing: --filePath=<file_path> should be added with command")
	}

	if *maxDownloadConcurrency < 0 {
		log.Fatal("Invalid Argument: --maxDownloadConcurrency should be greater than zero. Default 50.")
	}

	log.Println("File Path", *filePath)
	log.Println("Max Download Concurrency", *maxDownloadConcurrency)

	/* ------ stage 1 start------- */
	// Create the URL channel
	urlChan := make(chan string)

	// Initialize reader dependencies
	csvReader := &csv.CSVReader{
		FilePath: *filePath,
	}

	// Start goroutine to reading file
	var readingWait sync.WaitGroup
	readingWait.Add(1)
	go func() {
		defer func() {
			readingWait.Done()
		}()
		startReading(csvReader, urlChan)
	}()

	go func() {
		readingWait.Wait()
		close(urlChan)
	}()

	/* ------ stage 1 done ------- */

	/* ------ stage 2 start ------- */
	// Channels for passing data between stages
	contentChan := make(chan string)

	// Initialize downloader dependencies
	downloader := &downloader.HttpDownloader{
		Timeout: 10,
	}

	// Start goroutine to download URLs
	var downloadWait sync.WaitGroup
	downloadWait.Add(1)
	go func() {
		defer func() {
			downloadWait.Done()
		}()
		downloadURLs(*maxDownloadConcurrency, downloader, urlChan, contentChan, &downloadWait)
	}()

	go func() {
		downloadWait.Wait()
		close(contentChan)
	}()

	/* ------ stage 2 done ------- */
	/* ------ stage 3 start------- */

	// Initialize persister dependencies
	persister := &persister.FilePersister{
		BaseDir: "../../external/downloads",
	}

	// Start the goroutine to persist content
	var persistWG sync.WaitGroup
	persistWG.Add(1)
	go func() {
		defer func() {
			persistWG.Done()
		}()
		persistContent(persister, contentChan)
	}()

	var ctx context.Context
	var cancel context.CancelFunc
	doneCh := make(chan struct{})
	go func() {
		persistWG.Wait()
		ctx, cancel = context.WithTimeout(context.Background(), 1*time.Microsecond)
		close(doneCh)
	}()
	/* ------ stage 3 done------- */

	// Create a channel to listen for interrupt signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-signalChan:
		log.Println("Recieved an interrupt signal.")
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	case <-doneCh:
		log.Println("All goroutine task is completed.")
	}

	defer cancel()

	// Timeout or interrupt while waiting for completion
	select {
	case <-ctx.Done():
		log.Println("Graceful shutdown application.")
	}

	log.Printf("total URL processed from file: %.0f\n", totalURLs)
	log.Printf("success percentage: %.2f \n", (successFetch/totalURLs)*100)
	log.Printf("failure percentage: %.2f\n", (errorFetch/totalURLs)*100)
	log.Printf("average download duration: %.2f(ms)\n", totalDuration/totalURLs)
}

// Read URLs from CSV file line by line
func startReading(r ports.Reader, readChan chan<- string) {
	if err := r.ReadLines(readChan); err != nil {
		log.Fatalf("Error reading: %v", err)
	}
}

func downloadURLs(maxDownloadConcurrency int, downloader ports.Downloader, urlChan chan string, contentChan chan string, wg *sync.WaitGroup) {
	semaphore := make(chan struct{}, maxDownloadConcurrency)
	var mu sync.Mutex
	for url := range urlChan {
		totalURLs++
		wg.Add(1)

		go func(url string) {
			defer func() {
				wg.Done()
				<-semaphore
			}()
			semaphore <- struct{}{}
			downloadStartTime := time.Now()
			content, err := downloader.Download(url)
			if err != nil {
				log.Printf("Error downloading %s: %v", url, err)
				mu.Lock()
				errorFetch++
				mu.Unlock()
				return
			}
			mu.Lock()
			totalDuration += float64(time.Since(downloadStartTime).Milliseconds())
			mu.Unlock()

			contentChan <- content
		}(url)
	}
}

func persistContent(persister ports.Persister, contentChan chan string) {
	for content := range contentChan {
		successFetch++
		_, err := persister.Persist(content)
		if err != nil {
			log.Printf("Error persisting content: %v", err)
		}
	}
}
