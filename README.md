# SnapURL

**SnapURL** is a Go-based command-line application that downloads content from URLs listed in a CSV file and saves them into `.txt` files. It supports concurrent downloading and allows you to specify various parameters such as the input file path and maximum download concurrency.

---

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/anilvpatel21/snapurl.git
   ```

2. Navigate into the project directory:
   ```bash
   cd snapurl
   cd cmd
   ```

---

## Usage

### Running the Application

To run the application locally, use the following command:

```bash
go run main.go --filePath="../../external/example.csv"
```

In this command:
- `--filePath` specifies the relative or absolute path to the CSV file containing the URLs.

### Build the Application

Alternatively, you can build the executable and then run it:

```bash
go build -o snapurl
```

After building, run the compiled executable using:

```bash
./snapurl --filePath=<relative_path_or_absolute_path_to_csv>
```

### Optional Flag

- `--maxDownloadConcurrency` (optional): Specifies the maximum number of concurrent downloads. Default is `50`. Use this flag to limit the number of concurrent downloads, for example:

```bash
./snapurl --filePath="../../external/example.csv" --maxDownloadConcurrency=10
```

This will run the application with a maximum of 10 concurrent downloads.

---

## Command-Line Flags

- `--filePath`: (Required) Path to the input CSV file containing the URLs to download.
- `--maxDownloadConcurrency`: (Optional) The maximum number of concurrent downloads. Default is `50`.

Example:

```bash
go run main.go --filePath="../../external/example.csv" --maxDownloadConcurrency=10
```

---

## Example CSV File

The input CSV file (`example.csv`) should contain URLs in the following format:

```
urls
https://anilpatel.online
https://anilpatel.offline
```

---
