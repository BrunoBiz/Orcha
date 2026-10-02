# Build for Linux
$env:GOOS = "linux"
$env:GOARCH = "amd64"

$deployIP = "192.168.18.162"  # Orcha Container

go build -o Orcha

if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed."
    exit 1
}

# Deploy to server - app
scp .\Orcha root@[$deployIP]:/home/api/Orcha

if ($LASTEXITCODE -ne 0) {
    Write-Host "Copy failed - .\Orcha"
    exit 1
}

# Deploy to server - env
scp .\pm.env root@[$deployIP]:/home/api/pm.env

if ($LASTEXITCODE -ne 0) {
    Write-Host "Copy failed - \pm.env"
    exit 1
}

# CHMOD API
ssh "api@$deployIP" "sudo chmod +x Orcha;"

if ($LASTEXITCODE -ne 0) {
    Write-Host "SSH Failed - CHMOD."
    exit 1
}

Write-Host "Deployment successful!"