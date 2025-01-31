package main

import (
    "encoding/json"
    "fmt"
    "log"
    "math"
    "net/http"
    "strconv"
    "strings"
)

// SubnetResult represents the result of the subnet calculation
type SubnetResult struct {
    IP               string `json:"ip"`
    SubnetMask       int    `json:"subnetMask"`
    SubnetMaskDotted string `json:"subnetMaskDotted"`
    NetworkAddress   string `json:"networkAddress"`
    BroadcastAddress string `json:"broadcastAddress"`
    FirstUsableIP    string `json:"firstUsableIP"`
    LastUsableIP     string `json:"lastUsableIP"`
    NumHosts         int    `json:"numHosts"`
    UsableHostsRange string `json:"usableHostsRange"`
}

func calculateFromMask(ip string, maskBits int) (SubnetResult, error) {
    ipParts := parseIP(ip)
    if ipParts == nil {
        return SubnetResult{}, fmt.Errorf("invalid IP address")
    }

    subnetMask := make([]int, 4)
    for i := 0; i < 4; i++ {
        if maskBits >= (i+1)*8 {
            subnetMask[i] = 255
        } else if maskBits <= i*8 {
            subnetMask[i] = 0
        } else {
            subnetMask[i] = (255 << (8 - (maskBits % 8))) & 255
        }
    }

    networkAddress := make([]int, 4)
    broadcastAddress := make([]int, 4)
    for i := 0; i < 4; i++ {
        networkAddress[i] = ipParts[i] & subnetMask[i]
        broadcastAddress[i] = ipParts[i] | (^subnetMask[i] & 0xFF)
    }

    numHosts := int(math.Pow(2, float64(32-maskBits))) - 2
    var firstUsableIP, lastUsableIP []int
    var usableHostsRange string

    if numHosts > 0 {
        firstUsableIP = incrementIP(networkAddress)
        lastUsableIP = decrementIP(broadcastAddress)
        usableHostsRange = fmt.Sprintf("%d.%d.%d.%d - %d.%d.%d.%d",
            firstUsableIP[0], firstUsableIP[1], firstUsableIP[2], firstUsableIP[3],
            lastUsableIP[0], lastUsableIP[1], lastUsableIP[2], lastUsableIP[3])
    } else {
        firstUsableIP = []int{-1, -1, -1, -1}
        lastUsableIP = []int{-1, -1, -1, -1}
        usableHostsRange = "No usable hosts"
    }

    return SubnetResult{
        IP:               ip,
        SubnetMask:       maskBits,
        SubnetMaskDotted: fmt.Sprintf("%d.%d.%d.%d", subnetMask[0], subnetMask[1], subnetMask[2], subnetMask[3]),
        NetworkAddress:   fmt.Sprintf("%d.%d.%d.%d", networkAddress[0], networkAddress[1], networkAddress[2], networkAddress[3]),
        BroadcastAddress: fmt.Sprintf("%d.%d.%d.%d", broadcastAddress[0], broadcastAddress[1], broadcastAddress[2], broadcastAddress[3]),
        FirstUsableIP:    fmt.Sprintf("%d.%d.%d.%d", firstUsableIP[0], firstUsableIP[1], firstUsableIP[2], firstUsableIP[3]),
        LastUsableIP:     fmt.Sprintf("%d.%d.%d.%d", lastUsableIP[0], lastUsableIP[1], lastUsableIP[2], lastUsableIP[3]),
        NumHosts:         numHosts,
        UsableHostsRange: usableHostsRange,
    }, nil
}

func calculateFromHosts(ip string, numHosts int) (SubnetResult, error) {
    maskBits := 32 - int(math.Ceil(math.Log2(float64(numHosts)+2)))
    if maskBits < 0 || maskBits > 32 {
        return SubnetResult{}, fmt.Errorf("number of hosts exceeds the maximum possible for an IPv4 subnet")
    }
    return calculateFromMask(ip, maskBits)
}

func parseIP(ip string) []int {
    parts := strings.Split(ip, ".")
    if len(parts) != 4 {
        return nil
    }
    ipParts := make([]int, 4)
    for i, part := range parts {
        val, err := strconv.Atoi(part)
        if err != nil || val < 0 || val > 255 {
            return nil
        }
        ipParts[i] = val
    }
    return ipParts
}

func incrementIP(ip []int) []int {
    result := make([]int, 4)
    copy(result, ip)
    for i := 3; i >= 0; i-- {
        if result[i]++; result[i] <= 255 {
            break
        }
        result[i] = 0
    }
    return result
}

func decrementIP(ip []int) []int {
    result := make([]int, 4)
    copy(result, ip)
    for i := 3; i >= 0; i-- {
        if result[i]--; result[i] >= 0 {
            break
        }
        result[i] = 255
    }
    return result
}

func calculateSubnetHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    ip := r.FormValue("ip")
    typeInput := r.FormValue("type")
    value := r.FormValue("value")

    var result SubnetResult
    var err error

    switch typeInput {
    case "mask":
        maskBits, err := strconv.Atoi(value)
        if err != nil || maskBits < 0 || maskBits > 32 {
            http.Error(w, "Invalid subnet mask", http.StatusBadRequest)
            return
        }
        result, err = calculateFromMask(ip, maskBits)
    case "hosts":
        numHosts, err := strconv.Atoi(value)
        if err != nil || numHosts <= 0 {
            http.Error(w, "Invalid number of hosts", http.StatusBadRequest)
            return
        }
        result, err = calculateFromHosts(ip, numHosts)
    default:
        http.Error(w, "Invalid input type", http.StatusBadRequest)
        return
    }

    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}

func uiHandler(w http.ResponseWriter, r *http.Request) {
    html := `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Subnet Calculator</title>
  <style>
    body {
      font-family: Arial, sans-serif;
      background-color: #f4f4f9;
      margin: 0;
      padding: 20px;
      display: flex;
      justify-content: center;
      align-items: center;
      min-height: 100vh;
    }
    .container {
      max-width: 600px;
      width: 100%;
      background: #fff;
      padding: 20px;
      border-radius: 8px;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    }
    h1 {
      text-align: center;
      color: #333;
      margin-bottom: 20px;
    }
    label {
      display: block;
      margin-top: 15px;
      font-weight: bold;
      color: #555;
    }
    input, select {
      width: 100%;
      padding: 10px;
      margin-top: 5px;
      border: 1px solid #ccc;
      border-radius: 4px;
      font-size: 16px;
    }
    button {
      display: block;
      width: 100%;
      padding: 10px;
      margin-top: 20px;
      background-color: #007bff;
      color: white;
      border: none;
      border-radius: 4px;
      cursor: pointer;
      font-size: 16px;
    }
    button:hover {
      background-color: #0056b3;
    }
    .result {
      margin-top: 20px;
      padding: 10px;
      background-color: #e9f7ef;
      border-radius: 4px;
      border: 1px solid #c3e6cb;
      color: #155724;
    }
    .error {
      margin-top: 20px;
      padding: 10px;
      background-color: #f8d7da;
      border-radius: 4px;
      border: 1px solid #f5c6cb;
      color: #721c24;
    }
    .explanation {
      font-size: 14px;
      color: #555;
      margin-top: 5px;
    }

    @media (max-width: 600px) {
      .container {
        padding: 15px;
      }
      h1 {
        font-size: 20px;
      }
      input, select, button {
        font-size: 14px;
      }
    }
  </style>
</head>
<body>
  <div class="container">
    <h1>Subnet Calculator</h1>
    <form id="subnetForm">
      <label for="ip">IP Address:</label>
      <input type="text" id="ip" name="ip" placeholder="e.g., 192.168.1.1" required>

      <label for="type">Input Type:</label>
      <select id="type" name="type" required>
        <option value="mask">Subnet Mask (e.g., 24)</option>
        <option value="hosts">Number of Hosts</option>
      </select>

      <label for="value">Value:</label>
      <input type="text" id="value" name="value" placeholder="e.g., 24 or 254" required>

      <button type="submit">Calculate</button>
    </form>

    <div id="output"></div>
  </div>

  <script>
    document.getElementById('subnetForm').addEventListener('submit', function(event) {
      event.preventDefault();
      const formData = new FormData(this);
      fetch('/calculate', {
        method: 'POST',
        body: new URLSearchParams(formData),
      })
      .then(response => response.json())
      .then(result => {
        let outputHTML = '<div class="result">';
        outputHTML += '<strong>IP Address:</strong> ' + result.ip + '<br>';
        outputHTML += '<div class="explanation">The IP address you entered.</div>';
        
        outputHTML += '<strong>Subnet Mask:</strong> /' + result.subnetMask + ' (' + result.subnetMaskDotted + ')<br>';
        outputHTML += '<div class="explanation">The subnet mask defines the network portion of the IP address.</div>';
        
        outputHTML += '<strong>Network Address:</strong> ' + result.networkAddress + '<br>';
        outputHTML += '<div class="explanation">The first address in the subnet, used to identify the network.</div>';
        
        outputHTML += '<strong>Broadcast Address:</strong> ' + result.broadcastAddress + '<br>';
        outputHTML += '<div class="explanation">The last address in the subnet, used for broadcasting messages to all devices in the network.</div>';
        
        outputHTML += '<strong>First Usable IP:</strong> ' + result.firstUsableIP + '<br>';
        outputHTML += '<div class="explanation">The first usable IP address for assigning to devices in the network.</div>';
        
        outputHTML += '<strong>Last Usable IP:</strong> ' + result.lastUsableIP + '<br>';
        outputHTML += '<div class="explanation">The last usable IP address for assigning to devices in the network.</div>';
        
        outputHTML += '<strong>Number of Hosts:</strong> ' + result.numHosts + '<br>';
        outputHTML += '<div class="explanation">The total number of usable IP addresses available in the subnet.</div>';
        
        outputHTML += '<strong>Usable Hosts Range:</strong> ' + result.usableHostsRange + '<br>';
        outputHTML += '<div class="explanation">The range of IP addresses that can be assigned to devices in the network.</div>';
        
        outputHTML += '</div>';
        document.getElementById('output').innerHTML = outputHTML;
      })
      .catch(error => {
        document.getElementById('output').innerHTML = '<div class="error">An error occurred. Please try again.</div>';
      });
    });
  </script>
</body>
</html>
`
    w.Header().Set("Content-Type", "text/html")
    w.Write([]byte(html))
}

func main() {
    http.HandleFunc("/", uiHandler)
    http.HandleFunc("/calculate", calculateSubnetHandler)

    fmt.Println("Starting server on :8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
