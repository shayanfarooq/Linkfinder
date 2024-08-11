# LiNkFindeR
Welcome to  Link Extractor! 🌐🔍 This Go script is designed to scan JavaScript files and extract URLs and paths embedded within. It's a powerful tool for developers and security analysts looking to identify external resources referenced in JavaScript code.

Features
Comprehensive URL Extraction:
Extracts relative URLs, paths, and absolute URLs.
Identifies URLs from fetch API calls.
Uses advanced regex patterns to catch a wide variety of URL formats.
Color-Coded Output:
Green for URL headers
Cyan for extracted URLs
Usage
From Command-Line Arguments
Provide a URL directly as a command-line argument:
```
go run main.go https://example.com/script.js
```
From Standard Input
If no URL is provided, the script will read from standard input. Pipe or type a URL into the script:

```
echo "https://example.com/script.js" | go run main.go
```
How It Works
Fetch Content:

Retrieves JavaScript content from the specified URL.
Extract URLs:

Uses regex patterns to extract various types of URLs:
Relative URLs and paths (e.g., url('path/to/resource'), href="link")
URLs from fetch API calls
Absolute URLs (e.g., https://example.com)
Additional URL patterns with advanced regex
Output Results:

Prints extracted URLs, with results formatted for clarity.

Example Output
[URL] https://example.com/script.js
  https://cdn.example.com/resource.js
  /assets/images/logo.png
  https://api.example.com/data
  ../styles/main.css

Requirements
Go 1.18 or higher
Contribution
Contributions are welcome! If you have suggestions or find issues, feel free to open issues or submit pull requests.
