# Subnet Calculator

[![Deploy static content to Pages](https://github.com/anir0y/network/actions/workflows/static.yml/badge.svg)](https://github.com/anir0y/network/actions/workflows/static.yml)

A simple web-based subnet calculator with both web interface and Go CLI tool.

🌐 **[Try it online](https://anir0y.in/network/)**

## Features

*   Calculate subnet details from an IP address and subnet mask.
*   Calculate subnet details from an IP address and number of hosts.
*   Responsive design for different screen sizes.
*   Dark mode for better viewing in low-light environments.

## File Structure

### Web Interface
*   `index.html`: Main HTML structure
*   `styles.css`: Modern CSS styling with dark mode
*   `main.js`: JavaScript logic and validation

### Go CLI Tool
*   `main.go`: Go command-line subnet calculator
*   `go.mod`, `go.sum`: Go module dependencies

## Getting Started

To get started with this project, you can either clone the repository or download the files.

### Cloning the Repository

```bash
git clone https://github.com/anir0y/network-calculator.git
```

### Downloading the Files

You can also download the files directly from the GitHub repository.

## Usage

### Web Interface
1. Visit [https://anir0y.github.io/network/](https://anir0y.github.io/network/)
2. Enter an IP address and either a subnet mask or number of hosts
3. Click "Calculate Subnet" to see the results

### Go CLI Tool
```bash
# Install dependencies
go mod tidy

# Run the CLI tool
go run main.go

# Or build and run
go build -o subnet-calc
./subnet-calc
```

## Contributing

Contributions are welcome! If you have any ideas for improvements or new features, feel free to open an issue or submit a pull request.