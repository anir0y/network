/*
* Author: Animesh Roy 
* Author: Amreen Zariwala
* Code for Subnet Calculator
*/

package main

import (
	"fmt"
	"image/png"
	"io"
	"math"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/layout"
)

// Logo settings
const logoURL = "https://raw.githubusercontent.com/anir0y/cdn/refs/heads/main/anir0y.in%20logo%20(Curve)_4.%20Blue%20(Without%20Text).png"
const localLogoPath = "logo.png"

// Download logo if not exists
func downloadLogo() error {
    // Check if logo already exists and is valid
    if _, err := os.Stat(localLogoPath); err == nil {
        if verifyLogo(localLogoPath) {
            return nil
        }
        os.Remove(localLogoPath)
    }

    // Create a client with timeout
    client := &http.Client{
        Timeout: 10 * time.Second,
    }

    // Download the logo
    resp, err := client.Get(logoURL)
    if err != nil {
        return fmt.Errorf("failed to download logo: %v", err)
    }
    defer resp.Body.Close()

    // Check response status
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("bad status: %s", resp.Status)
    }

    // Create the file
    out, err := os.Create(localLogoPath)
    if err != nil {
        return fmt.Errorf("failed to create file: %v", err)
    }
    defer out.Close()

    // Copy the response body to the file
    _, err = io.Copy(out, resp.Body)
    if err != nil {
        os.Remove(localLogoPath)
        return fmt.Errorf("failed to save file: %v", err)
    }

    // Verify the downloaded file
    if !verifyLogo(localLogoPath) {
        os.Remove(localLogoPath)
        return fmt.Errorf("downloaded file is corrupted")
    }

    return nil
}

// Verify if logo file is valid
func verifyLogo(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	_, err = png.Decode(file)
	return err == nil
}

// Load logo from local file
func loadLocalLogo(filePath string) fyne.CanvasObject {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening local file:", err)
		return widget.NewLabel("Logo not available")
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		fmt.Println("Error decoding local image:", err)
		os.Remove(filePath)
		return widget.NewLabel("Logo not available")
	}

	imageWidget := canvas.NewImageFromImage(img)
	imageWidget.FillMode = canvas.ImageFillContain
	imageWidget.SetMinSize(fyne.NewSize(180, 180))

	return imageWidget
}

// Subnet Calculation Logic
type SubnetResult struct {
	IP               string
	SubnetMask       int
	SubnetMaskDotted string
	NetworkAddress   string
	BroadcastAddress string
	FirstUsableIP    string
	LastUsableIP     string
	NumHosts         int
	UsableHostsRange string
}

func parseIP(ip string) (net.IP, error) {
	parsed := net.ParseIP(ip).To4()
	if parsed == nil {
		return nil, fmt.Errorf("invalid IP address")
	}
	return parsed, nil
}

func calculateFromMask(ip string, maskBits int) (SubnetResult, error) {
	parsedIP, err := parseIP(ip)
	if err != nil {
		return SubnetResult{}, err
	}

	mask := net.CIDRMask(maskBits, 32)
	network := parsedIP.Mask(mask)
	broadcast := net.IP(make([]byte, 4))

	for i := 0; i < 4; i++ {
		broadcast[i] = network[i] | ^mask[i]
	}

	numHosts := int(math.Pow(2, float64(32-maskBits))) - 2
	if maskBits == 31 {
		numHosts = 2
	} else if maskBits == 32 {
		numHosts = 0
	}

	firstUsable, lastUsable := "N/A", "N/A"
	if numHosts > 0 {
		firstUsable = incrementIP(network)
		lastUsable = decrementIP(broadcast)
	}

	return SubnetResult{
		IP:               ip,
		SubnetMask:       maskBits,
		SubnetMaskDotted: net.IP(mask).String(),
		NetworkAddress:   network.String(),
		BroadcastAddress: broadcast.String(),
		FirstUsableIP:    firstUsable,
		LastUsableIP:     lastUsable,
		NumHosts:         numHosts,
		UsableHostsRange: fmt.Sprintf("%s - %s", firstUsable, lastUsable),
	}, nil
}

// Increment an IP address
func incrementIP(ip net.IP) string {
	ipCopy := append([]byte(nil), ip...)
	for i := 3; i >= 0; i-- {
		ipCopy[i]++
		if ipCopy[i] != 0 {
			break
		}
	}
	return net.IP(ipCopy).String()
}

// Decrement an IP address
func decrementIP(ip net.IP) string {
	ipCopy := append([]byte(nil), ip...)
	for i := 3; i >= 0; i-- {
		if ipCopy[i]--; ipCopy[i] != 255 {
			break
		}
	}
	return net.IP(ipCopy).String()
}

// Show subnet results
func showResults(result SubnetResult, output *widget.Label) {
	output.SetText(fmt.Sprintf(`IP Address: %s
Subnet Mask: /%d (%s)
Network Address: %s
Broadcast Address: %s
First Usable IP: %s
Last Usable IP: %s
Number of Hosts: %d
Usable Hosts Range: %s`,
		result.IP, result.SubnetMask, result.SubnetMaskDotted, result.NetworkAddress,
		result.BroadcastAddress, result.FirstUsableIP, result.LastUsableIP, result.NumHosts, result.UsableHostsRange))
}

func main() {
    // Download logo and handle any errors before launching the app
    err := downloadLogo()
    if err != nil {
        fmt.Printf("Failed to download logo: %v\nLaunching without logo...\n", err)
    }

    a := app.New()
    w := a.NewWindow("Subnet Calculator")
	w.Resize(fyne.NewSize(699, 600))
	w.SetFixedSize(true)

	// Load Logo
	logo := loadLocalLogo(localLogoPath)

	// Input Fields
	ipEntry := widget.NewEntry()
	ipEntry.SetPlaceHolder("e.g., 192.168.1.1")

	typeSelect := widget.NewSelect([]string{"Subnet Mask", "Number of Hosts"}, nil)
	typeSelect.SetSelected("Subnet Mask")

	valueEntry := widget.NewEntry()
	valueEntry.SetPlaceHolder("e.g., 24 or 254")

	output := widget.NewLabel("Results will appear here.")
	output.Wrapping = fyne.TextWrapWord

	// Calculate Button
	calculateButton := widget.NewButton("Calculate", func() {
		ip := ipEntry.Text
		maskBits, err := strconv.Atoi(valueEntry.Text)
		if err != nil || maskBits < 0 || maskBits > 32 {
			dialog.ShowError(fmt.Errorf("invalid subnet mask"), w)
			return
		}
		result, err := calculateFromMask(ip, maskBits)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		showResults(result, output)
	})

	// Custom Menu Button positioned at top left
    menuButton := widget.NewButton("☰ About", func() {
        dialog.ShowInformation("About", "App Developer:\nAnimesh Roy\nAmreen Zariwala\n\nORG: Arishti Consolidated Private Limited\nEmail: admin@arishtisecurity.com", w)
    })
    menuButton.Resize(fyne.NewSize(80, 30))
    
    // Create a top container for the menu button
    topContainer := container.NewHBox(menuButton, layout.NewSpacer())

// UI Layout
    // UI Layout
    form := container.NewVBox(ipEntry, typeSelect, valueEntry, calculateButton)
    
    // Right padding container with both horizontal and vertical padding
    paddingContainer := container.NewVBox(
        widget.NewLabel(" "),  // Top padding
        container.NewHBox(form, widget.NewLabel("   ")),  // Horizontal padding
        widget.NewLabel(" "),  // Bottom padding
    )

    mainContent := container.NewBorder(topContainer, output, logo, paddingContainer)

w.SetContent(mainContent)
w.ShowAndRun()
}
