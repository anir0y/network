document.addEventListener('DOMContentLoaded', () => {
    const themeToggle = document.getElementById('theme-toggle');
    const subnetForm = document.getElementById('subnetForm');
    const outputDiv = document.getElementById('output');
    const resetBtn = document.getElementById('reset-btn');
    const typeSelect = document.getElementById('type');
    const valueLabel = document.getElementById('value-label');
    const valueInput = document.getElementById('value');
    const valueHelp = document.getElementById('value-help');
    const valueIcon = document.getElementById('value-icon');

    // Theme switcher
    themeToggle.addEventListener('change', () => {
        document.documentElement.setAttribute('data-theme', themeToggle.checked ? 'dark' : 'light');
    });

    // Update form labels based on calculation type
    typeSelect.addEventListener('change', updateFormLabels);
    
    function updateFormLabels() {
        const type = typeSelect.value;
        if (type === 'mask') {
            valueLabel.textContent = 'Subnet Mask (CIDR)';
            valueInput.placeholder = '24';
            valueHelp.textContent = 'Enter a value between 0-32 for CIDR notation';
            valueIcon.textContent = '📊';
        } else {
            valueLabel.textContent = 'Number of Hosts';
            valueInput.placeholder = '254';
            valueHelp.textContent = 'Enter the number of hosts needed';
            valueIcon.textContent = '👥';
        }
        clearValidationErrors();
    }

    // Form validation
    function validateForm() {
        const ip = document.getElementById('ip').value.trim();
        const value = valueInput.value.trim();
        let isValid = true;

        clearValidationErrors();

        if (!ip) {
            showFieldError('ip', 'IP address is required');
            isValid = false;
        } else if (!isValidIP(ip)) {
            showFieldError('ip', 'Please enter a valid IP address');
            isValid = false;
        }

        if (!value) {
            showFieldError('value', 'This field is required');
            isValid = false;
        } else if (typeSelect.value === 'mask') {
            const maskBits = parseInt(value, 10);
            if (isNaN(maskBits) || maskBits < 0 || maskBits > 32) {
                showFieldError('value', 'Subnet mask must be between 0 and 32');
                isValid = false;
            }
        } else {
            const numHosts = parseInt(value, 10);
            if (isNaN(numHosts) || numHosts <= 0) {
                showFieldError('value', 'Number of hosts must be a positive integer');
                isValid = false;
            }
        }

        return isValid;
    }

    function isValidIP(ip) {
        const parts = ip.split('.');
        return parts.length === 4 && parts.every(part => {
            const num = parseInt(part, 10);
            return !isNaN(num) && num >= 0 && num <= 255;
        });
    }

    function showFieldError(fieldId, message) {
        const errorElement = document.getElementById(`${fieldId}-error`);
        if (errorElement) {
            errorElement.textContent = message;
        }
    }

    function clearValidationErrors() {
        document.querySelectorAll('.error-message').forEach(el => el.textContent = '');
    }

    // Form submission
    subnetForm.addEventListener('submit', (e) => {
        e.preventDefault();
        if (validateForm()) {
            calculateSubnet();
        }
    });

    // Reset button
    resetBtn.addEventListener('click', () => {
        outputDiv.innerHTML = '';
        clearValidationErrors();
    });

    // Initialize form labels
    updateFormLabels();

    function calculateSubnet() {
        const ip = document.getElementById('ip').value.trim();
        const type = document.getElementById('type').value;
        const value = valueInput.value.trim();

        outputDiv.innerHTML = `
            <div class="loading-card fade-in">
                <div class="loading-spinner"></div>
                <p class="loading-text">Calculating subnet details...</p>
            </div>
        `;

        setTimeout(() => {
            try {
                let result;
                if (type === 'mask') {
                    const maskBits = parseInt(value, 10);
                    result = calculateFromMask(ip, maskBits);
                } else {
                    const numHosts = parseInt(value, 10);
                    result = calculateFromHosts(ip, numHosts);
                }

                displayResult(result);
            } catch (error) {
                displayError(error.message);
            }
        }, 800);
    }

    function displayResult(result) {
        outputDiv.innerHTML = `
            <div class="result-card fade-in">
                <div class="result-header">
                    <h2 class="result-title">
                        <span>✅</span>
                        <span>Subnet Calculation Results</span>
                    </h2>
                </div>
                <div class="result-grid">
                    ${createResultItem('IP Address', result.ip)}
                    ${createResultItem('Subnet Mask', `/${result.subnetMask} (${result.subnetMaskDotted})`)}
                    ${createResultItem('Network Address', result.networkAddress)}
                    ${createResultItem('Broadcast Address', result.broadcastAddress)}
                    ${createResultItem('First Usable IP', result.firstUsableIP)}
                    ${createResultItem('Last Usable IP', result.lastUsableIP)}
                    ${createResultItem('Total Hosts', result.numHosts.toLocaleString())}
                    ${createResultItem('Usable Range', result.usableHostsRange)}
                </div>
            </div>
        `;
        addCopyFunctionality();
    }

    function createResultItem(label, value) {
        return `
            <div class="result-item">
                <div class="result-label">${label}</div>
                <div class="result-value">
                    <span>${value}</span>
                    <button class="copy-btn" data-value="${value}" aria-label="Copy ${label}">
                        📋
                        <span class="tooltip">Copy to clipboard</span>
                    </button>
                </div>
            </div>
        `;
    }

    function displayError(message) {
        outputDiv.innerHTML = `
            <div class="error-card fade-in">
                <div class="error-content">
                    <span class="error-icon">❌</span>
                    <span>${message}</span>
                </div>
            </div>
        `;
    }

    function addCopyFunctionality() {
        document.querySelectorAll('.copy-btn').forEach(button => {
            button.addEventListener('click', async () => {
                const valueToCopy = button.getAttribute('data-value');
                try {
                    await navigator.clipboard.writeText(valueToCopy);
                    button.classList.add('copied');
                    button.innerHTML = '✅<span class="tooltip">Copied!</span>';
                    setTimeout(() => {
                        button.classList.remove('copied');
                        button.innerHTML = '📋<span class="tooltip">Copy to clipboard</span>';
                    }, 2000);
                } catch (err) {
                    // Fallback for older browsers
                    const textArea = document.createElement('textarea');
                    textArea.value = valueToCopy;
                    document.body.appendChild(textArea);
                    textArea.select();
                    document.execCommand('copy');
                    document.body.removeChild(textArea);
                    
                    button.classList.add('copied');
                    button.innerHTML = '✅<span class="tooltip">Copied!</span>';
                    setTimeout(() => {
                        button.classList.remove('copied');
                        button.innerHTML = '📋<span class="tooltip">Copy to clipboard</span>';
                    }, 2000);
                }
            });
        });
    }

    function calculateFromMask(ip, maskBits) {
        const ipParts = ip.split('.').map(Number);
        if (ipParts.length !== 4 || ipParts.some(part => isNaN(part) || part < 0 || part > 255)) {
            throw new Error('Invalid IP address.');
        }

        const subnetMask = Array(4).fill(0).map((_, i) => {
            if (maskBits >= (i + 1) * 8) {
                return 255;
            } else if (maskBits <= i * 8) {
                return 0;
            } else {
                return (255 << (8 - (maskBits % 8))) & 255;
            }
        });

        const networkAddress = ipParts.map((part, i) => part & subnetMask[i]);
        const broadcastAddress = ipParts.map((part, i) => part | (~subnetMask[i] & 0xFF));

        const numHosts = Math.pow(2, 32 - maskBits) - 2;
        let firstUsableIP, lastUsableIP, usableHostsRange;

        if (numHosts > 0) {
            firstUsableIP = [...networkAddress];
            incrementIP(firstUsableIP);

            lastUsableIP = [...broadcastAddress];
            decrementIP(lastUsableIP);

            usableHostsRange = `${firstUsableIP.join('.')} - ${lastUsableIP.join('.')}`;
        } else {
            firstUsableIP = 'N/A';
            lastUsableIP = 'N/A';
            usableHostsRange = 'No usable hosts';
        }

        return {
            ip,
            subnetMask: maskBits,
            subnetMaskDotted: subnetMask.join('.'),
            networkAddress: networkAddress.join('.'),
            broadcastAddress: broadcastAddress.join('.'),
            firstUsableIP: firstUsableIP === 'N/A' ? firstUsableIP : firstUsableIP.join('.'),
            lastUsableIP: lastUsableIP === 'N/A' ? lastUsableIP : lastUsableIP.join('.'),
            numHosts: numHosts > 0 ? numHosts : 0,
            usableHostsRange,
        };
    }

    function calculateFromHosts(ip, numHosts) {
        const ipParts = ip.split('.').map(Number);
        if (ipParts.length !== 4 || ipParts.some(part => isNaN(part) || part < 0 || part > 255)) {
            throw new Error('Invalid IP address.');
        }

        const maskBits = 32 - Math.ceil(Math.log2(numHosts + 2));
        if (maskBits < 0 || maskBits > 32) {
            throw new Error('Number of hosts exceeds the maximum possible for an IPv4 subnet.');
        }

        return calculateFromMask(ip, maskBits);
    }

    function incrementIP(ipParts) {
        for (let i = 3; i >= 0; i--) {
            if (++ipParts[i] <= 255) break;
            ipParts[i] = 0;
        }
    }

    function decrementIP(ipParts) {
        for (let i = 3; i >= 0; i--) {
            if (--ipParts[i] >= 0) break;
            ipParts[i] = 255;
        }
    }
});