# benchmark.ps1 — Benchmark API response time
# Usage: .\benchmark.ps1

$url = "http://localhost:3010/api/v1/events?page=1&limit=10"
$totalRequests = 50

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  API Benchmark - $totalRequests requests" -ForegroundColor Cyan
Write-Host "  URL: $url" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$times = @()

for ($i = 1; $i -le $totalRequests; $i++) {
    $stopwatch = [System.Diagnostics.Stopwatch]::StartNew()
    
    try {
        $response = Invoke-RestMethod -Uri $url -Method GET -ErrorAction Stop
        $stopwatch.Stop()
        $ms = $stopwatch.Elapsed.TotalMilliseconds
        $times += $ms
        
        # Progress bar
        $bar = "#" * [math]::Floor($ms / 2)
        $color = if ($ms -lt 5) { "Green" } elseif ($ms -lt 20) { "Yellow" } else { "Red" }
        Write-Host ("[{0,3}] {1,8:F2} ms  {2}" -f $i, $ms, $bar) -ForegroundColor $color
    }
    catch {
        Write-Host ("[{0,3}] FAILED - {1}" -f $i, $_.Exception.Message) -ForegroundColor Red
    }
}

# Statistik
Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  RESULTS" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan

$avg = ($times | Measure-Object -Average).Average
$min = ($times | Measure-Object -Minimum).Minimum
$max = ($times | Measure-Object -Maximum).Maximum
$total = ($times | Measure-Object -Sum).Sum

Write-Host ("  Total Requests : {0}" -f $times.Count)
Write-Host ("  Total Time     : {0:F2} ms" -f $total)
Write-Host ("  Average        : {0:F2} ms" -f $avg) -ForegroundColor Yellow
Write-Host ("  Min            : {0:F2} ms" -f $min) -ForegroundColor Green
Write-Host ("  Max            : {0:F2} ms" -f $max) -ForegroundColor Red
Write-Host ("  Req/sec        : {0:F1}" -f (1000 / $avg))
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
