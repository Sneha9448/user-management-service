
$logFile = "auth_test_results.log"
"Starting Auth Verification" | Out-File $logFile

try {
    # 1. Login with Google (Mock)
    Write-Host "1. Login with Google"
    "1. Login with Google" | Out-File $logFile -Append
    $body = @{ query = 'mutation { loginWithGoogle(idToken: "mock_token") { token user { id name email } } }' } | ConvertTo-Json
    $res = Invoke-RestMethod -Uri http://localhost:8081/graphql -Method Post -Body $body -ContentType "application/json"
    $token = $res.data.loginWithGoogle.token
    "Login Token: $token" | Out-File $logFile -Append
    ($res | ConvertTo-Json -Depth 10) | Out-File $logFile -Append

    # 2. Fetch current user (me) with token
    Write-Host "2. Fetch current user (me)"
    "2. Fetch current user (me)" | Out-File $logFile -Append
    $meBody = @{ query = 'query { me { id name email } }' } | ConvertTo-Json
    $meRes = Invoke-RestMethod -Uri http://localhost:8081/graphql -Method Post -Header @{ Authorization = "Bearer $token" } -Body $meBody -ContentType "application/json"
    ($meRes | ConvertTo-Json -Depth 10) | Out-File $logFile -Append

    # 3. Try protected mutation (createUser) WITHOUT token
    Write-Host "3. Create user WITHOUT token (Should fail)"
    "3. Create user WITHOUT token (Should fail)" | Out-File $logFile -Append
    $createBody = @{ query = 'mutation { createUser(name: "No Auth", email: "noauth@test.com") { id name } }' } | ConvertTo-Json
    try {
        $failRes = Invoke-RestMethod -Uri http://localhost:8081/graphql -Method Post -Body $createBody -ContentType "application/json"
        ($failRes | ConvertTo-Json -Depth 10) | Out-File $logFile -Append
    } catch {
        "Failed as expected: $_" | Out-File $logFile -Append
    }

    # 4. Try protected mutation (createUser) WITH token
    Write-Host "4. Create user WITH token (Should succeed)"
    "4. Create user WITH token (Should succeed)" | Out-File $logFile -Append
    $successBody = @{ query = 'mutation { createUser(name: "Auth Success", email: "authsuccess@test.com") { id name } }' } | ConvertTo-Json
    $successRes = Invoke-RestMethod -Uri http://localhost:8081/graphql -Method Post -Header @{ Authorization = "Bearer $token" } -Body $successBody -ContentType "application/json"
    ($successRes | ConvertTo-Json -Depth 10) | Out-File $logFile -Append

    Write-Host "Verification complete. Results saved to $logFile"
    "Verification complete" | Out-File $logFile -Append
} catch {
    "FATAL ERROR: $_" | Out-File $logFile -Append
    Write-Error $_
}
