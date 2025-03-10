package main

import (
	"log"
	"os"
	"sync"
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

	// Read input file path from CLI
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <file_path>")
	}
	filePath := os.Args[1]

	/* ------ stage 1 start------- */
	// Create the URL channel
	urlChan := make(chan string)

	// Initialize reader dependencies
	csvReader := &csv.CSVReader{
		FilePath: filePath,
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
		downloadURLs(downloader, urlChan, contentChan, &downloadWait)
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

	persistWG.Wait()
	/* ------ stage 3 done------- */
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

func downloadURLs(downloader ports.Downloader, urlChan chan string, contentChan chan string, wg *sync.WaitGroup) {
	semaphore := make(chan struct{}, 50)
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
