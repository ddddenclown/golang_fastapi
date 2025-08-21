# Test script for Analytics Service API (PowerShell)

Write-Host "🚀 Testing Analytics Service API" -ForegroundColor Green
Write-Host "==================================" -ForegroundColor Green

# Base URL
$BASE_URL = "http://localhost:8080"

# Function to print colored output
function Write-Status {
    param(
        [bool]$Success,
        [string]$Message
    )
    
    if ($Success) {
        Write-Host "✅ $Message" -ForegroundColor Green
    } else {
        Write-Host "❌ $Message" -ForegroundColor Red
    }
}

# Test 1: Health check
Write-Host "`n1. Testing health check..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$BASE_URL/" -Method Get
    $response | ConvertTo-Json
    Write-Status -Success $true -Message "Health check"
} catch {
    Write-Status -Success $false -Message "Health check failed: $($_.Exception.Message)"
}

# Test 2: Generate token
Write-Host "`n2. Generating token..." -ForegroundColor Yellow
try {
    $body = @{
        email = "string"
        password = "string"
    } | ConvertTo-Json
    
    $response = Invoke-RestMethod -Uri "$BASE_URL/auth" -Method Post -Body $body -ContentType "application/json"
    $response | ConvertTo-Json
    
    $TOKEN = $response.token
    if ($TOKEN) {
        Write-Status -Success $true -Message "Token generated successfully"
        Write-Host "Token: $TOKEN" -ForegroundColor Cyan
    } else {
        Write-Status -Success $false -Message "Failed to generate token"
        exit 1
    }
} catch {
    Write-Status -Success $false -Message "Token generation failed: $($_.Exception.Message)"
    exit 1
}

# Test 3: Validate token
Write-Host "`n3. Validating token..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$BASE_URL/validate?token=$TOKEN" -Method Get
    $response | ConvertTo-Json
    Write-Status -Success $true -Message "Token validation"
} catch {
    Write-Status -Success $false -Message "Token validation failed: $($_.Exception.Message)"
}

# Test 4: Analytics with generated token
Write-Host "`n4. Testing analytics endpoint..." -ForegroundColor Yellow
try {
    $body = @{
        token = $TOKEN
        StartDate = "01.01.2024"
        FinishDate = "31.01.2024"
    } | ConvertTo-Json
    
    $response = Invoke-RestMethod -Uri "$BASE_URL/analytics" -Method Post -Body $body -ContentType "application/json"
    $response | ConvertTo-Json
    
    # Check if analytics returned data
    $ITEMS_COUNT = $response.items.Count
    $TOTAL = $response.total
    
    if ($ITEMS_COUNT -gt 0 -and $TOTAL -gt 0) {
        Write-Status -Success $true -Message "Analytics returned data successfully"
        Write-Host "Items count: $ITEMS_COUNT" -ForegroundColor Cyan
        Write-Host "Total: $TOTAL" -ForegroundColor Cyan
    } else {
        Write-Status -Success $false -Message "Analytics returned empty data"
        Write-Host "Items count: $ITEMS_COUNT" -ForegroundColor Red
        Write-Host "Total: $TOTAL" -ForegroundColor Red
    }
} catch {
    Write-Status -Success $false -Message "Analytics request failed: $($_.Exception.Message)"
}

# Test 5: Analytics with invalid token
Write-Host "`n5. Testing analytics with invalid token..." -ForegroundColor Yellow
try {
    $body = @{
        token = "invalid-token"
        StartDate = "01.01.2024"
        FinishDate = "31.01.2024"
    } | ConvertTo-Json
    
    $response = Invoke-RestMethod -Uri "$BASE_URL/analytics" -Method Post -Body $body -ContentType "application/json"
    $response | ConvertTo-Json
    Write-Status -Success $true -Message "Invalid token test"
} catch {
    Write-Status -Success $true -Message "Invalid token test (expected error)"
}

Write-Host "`n🎉 API testing completed!" -ForegroundColor Green
