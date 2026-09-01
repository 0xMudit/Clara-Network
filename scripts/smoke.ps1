# smoke.ps1 - one-click end-to-end smoke test for Clara Network (Windows twin
# of scripts/smoke.sh). PowerShell 5.1 compatible. Prefer running smoke.sh from
# Git Bash for full parity; this is a native-PowerShell fallback.
#
# Boots the local docker-compose stack, seeds demo data via the one-shot sims,
# then runs three check suites (DB, backend API, frontend build + BFF guard)
# and prints a PASS/FAIL table. Exit code non-zero if any check fails.
#
# Usage:
#   .\scripts\smoke.ps1                # full run
#   .\scripts\smoke.ps1 -NoSeed        # reuse an already-seeded running stack
#   .\scripts\smoke.ps1 -Keep          # leave stack running
#   .\scripts\smoke.ps1 -SkipFrontend  # skip the slow next build
#   .\scripts\smoke.ps1 -Help

param(
    [switch]$Help,
    [switch]$NoSeed,
    [switch]$Keep,
    [switch]$SkipFrontend
)

$ErrorActionPreference = "Stop"
$scriptDir = $PSScriptRoot
if (-not $scriptDir) { $scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition }
$repoRoot = [System.IO.Path]::GetFullPath((Join-Path $scriptDir ".."))
$compose = Join-Path $repoRoot "deploy\docker-compose.yml"
$adminPort = if ($env:CLARA_ADMINAPI_PORT) { $env:CLARA_ADMINAPI_PORT } else { 18083 }
$adminBase = "http://localhost:$adminPort"
$webDir = Join-Path $repoRoot "web"

function Green { param($s) Write-Host $s -ForegroundColor Green }
function Red   { param($s) Write-Host $s -ForegroundColor Red }
function Yellow{ param($s) Write-Host $s -ForegroundColor Yellow }
function Step  { param($s) Write-Host ""; Write-Host "== $s ==" -ForegroundColor Cyan }

$script:pass = 0
$script:fail = 0
$script:results = New-Object System.Collections.Generic.List[string]

function Check([string]$name, [bool]$ok, [string]$detail = "") {
    if ($ok) {
        $script:pass++
        $script:results.Add($name + "|PASS")
        Write-Host ("  PASS  " + $name) -ForegroundColor Green
    } else {
        $script:fail++
        $script:results.Add($name + "|FAIL|" + $detail)
        $msg = "  FAIL  $name"
        if ($detail) { $msg += "  ($detail)" }
        Write-Host $msg -ForegroundColor Red
    }
}

function Invoke-JsonOk([string]$name, [string]$url) {
    try {
        $resp = Invoke-WebRequest -Uri $url -UseBasicParsing -TimeoutSec 20 -ErrorAction Stop
        $body = $resp.Content
        $json = $body | ConvertFrom-Json -ErrorAction Stop
        if ($null -eq $json) {
            Check $name $false "empty JSON body"
        } else {
            Check $name $true
        }
    } catch {
        if ($_.Exception.Response) {
            $code = [int]$_.Exception.Response.StatusCode
            Check $name $false "HTTP $code"
        } else {
            Check $name $false $_.Exception.Message
        }
    }
}

function Wait-Health([string]$url, [string]$label, [int]$tries = 60) {
    Write-Host "  waiting for $label ($url) ..."
    for ($i = 0; $i -lt $tries; $i++) {
        try {
            $r = Invoke-WebRequest -Uri $url -UseBasicParsing -TimeoutSec 5 -ErrorAction Stop
            if ($r.StatusCode -eq 200) { return $true }
        } catch { }
        Start-Sleep -Seconds 2
    }
    return $false
}

function Status-Of([string]$url) {
    try {
        $r = Invoke-WebRequest -Uri $url -UseBasicParsing -TimeoutSec 15 -ErrorAction Stop -MaximumRedirection 0
        return [int]$r.StatusCode
    } catch {
        if ($_.Exception.Response) { return [int]$_.Exception.Response.StatusCode }
        return 0
    }
}

function PsqlQuery([string]$q) {
    $out = & docker compose -f $compose exec -T postgres psql -U clara -d clara -tA -c $q 2>&1
    return ($out | Out-String).Trim()
}

if ($Help) {
    Get-Content $MyInvocation.MyCommand.Path | Where-Object { $_ -match '^#' } | ForEach-Object { $_ -replace '^#\s?','' }
    exit 0
}

if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
    Red "ERROR: docker not found on PATH"
    exit 2
}

Yellow "Clara Network smoke test"
Yellow "  adminapi: $adminBase   compose: $compose"

Step "Boot docker-compose stack"
& docker compose -f $compose up -d --build
if ($LASTEXITCODE -ne 0) {
    Red "ERROR: docker compose up failed"
    & docker compose -f $compose logs --tail=50
    exit 2
}

if (-not (Wait-Health "$adminBase/health" "adminapi" 60)) {
    Red "ERROR: adminapi did not become healthy"
    & docker compose -f $compose logs adminapi --tail=80
    exit 2
}
Green "  adminapi is healthy"

if (-not $NoSeed) {
    Step "Seed demo data (one-shot sims)"
    foreach ($svc in @("clearing-sim","ledger-sim","card-sim","acquiring-sim","disputes-sim","acquirer-sim")) {
        Write-Host "  running $svc ..."
        & docker compose -f $compose run --rm --no-deps $svc | Out-Null
    }
}

Step "DB suite (schema + seed data)"
foreach ($table in @("switch_transactions","clearing_records","net_positions","settlement_instructions","prefund_accounts","default_fund","ledger_accounts","ledger_entries","bin_ranges","cards","tokens","merchants","funding_lines","screening_lists","disputes","dispute_transactions")) {
    $n = PsqlQuery "SELECT count(*) FROM information_schema.tables WHERE table_schema='public' AND table_name='$table';"
    Check ("db: table " + $table + " exists") ($n -eq "1") $n
}
foreach ($t in @("switch_transactions","clearing_records","net_positions","settlement_instructions","ledger_accounts","ledger_entries","bin_ranges","cards","tokens","merchants","funding_lines","screening_lists","disputes")) {
    $n = PsqlQuery "SELECT count(*) FROM $t;"
    $ok = $false
    if ($n -match '^\d+$') { [int]$cnt = $n; if ($cnt -gt 0) { $ok = $true } }
    Check ("db: " + $t + " has data") $ok ("row count=$n")
}

Step "Backend API suite (adminapi endpoints)"
Check "api: /health" $true "via boot wait"
Invoke-JsonOk "api: /dashboard"                 "$adminBase/api/v1/dashboard"
Invoke-JsonOk "api: /transactions"              "$adminBase/api/v1/transactions?limit=5"
Invoke-JsonOk "api: clearing/cycles"            "$adminBase/api/v1/clearing/cycles"

# clearing/records, clearing/positions and settlement/instructions require a
# ?cycle= param (matching how the frontend drives them), so resolve a real one
# from /clearing/cycles first.
$cycle = ""
try {
    $cycles = (Invoke-WebRequest -Uri "$adminBase/api/v1/clearing/cycles" -UseBasicParsing -TimeoutSec 20 -ErrorAction Stop).Content | ConvertFrom-Json
    $cycle = @($cycles.items)[0]
} catch { }
if (-not $cycle) {
    Check "api: resolve clearing cycle" $false "no cycles returned by /clearing/cycles"
} else {
    Check "api: resolve clearing cycle" $true "cycle=$cycle"
    Invoke-JsonOk "api: clearing/records?cycle"           "$adminBase/api/v1/clearing/records?cycle=$cycle"
    Invoke-JsonOk "api: clearing/positions?cycle"         "$adminBase/api/v1/clearing/positions?cycle=$cycle"
    Invoke-JsonOk "api: settlement/instructions?cycle"    "$adminBase/api/v1/settlement/instructions?cycle=$cycle"
}
Invoke-JsonOk "api: settlement/prefunds"        "$adminBase/api/v1/settlement/prefunds"
Invoke-JsonOk "api: settlement/default-fund"    "$adminBase/api/v1/settlement/default-fund"
Invoke-JsonOk "api: ledger/accounts"            "$adminBase/api/v1/ledger/accounts"
Invoke-JsonOk "api: ledger/entries"             "$adminBase/api/v1/ledger/entries?limit=5"
Invoke-JsonOk "api: cards"                      "$adminBase/api/v1/cards?limit=5"
Invoke-JsonOk "api: bin-ranges"                 "$adminBase/api/v1/bin-ranges"
Invoke-JsonOk "api: tokens"                     "$adminBase/api/v1/tokens?limit=5"
Invoke-JsonOk "api: merchants"                  "$adminBase/api/v1/merchants?limit=5"
Invoke-JsonOk "api: disputes"                   "$adminBase/api/v1/disputes?limit=5"
Invoke-JsonOk "api: disputes/overdue"           "$adminBase/api/v1/disputes/overdue"

if (-not $SkipFrontend) {
    Step "Frontend suite (next build + runtime auth/BFF guard)"
    if (Test-Path (Join-Path $webDir "node_modules")) {
        Push-Location $webDir
        try {
            $env:NEXT_PUBLIC_SUPABASE_URL = if ($env:NEXT_PUBLIC_SUPABASE_URL) { $env:NEXT_PUBLIC_SUPABASE_URL } else { "http://localhost:54321" }
            $env:NEXT_PUBLIC_SUPABASE_ANON_KEY = if ($env:NEXT_PUBLIC_SUPABASE_ANON_KEY) { $env:NEXT_PUBLIC_SUPABASE_ANON_KEY } else { "smoke-placeholder-anon" }
            & npm.cmd run build *> (Join-Path ([System.IO.Path]::GetTempPath()) ("clara_nextbuild_" + $PID + ".log"))
            Check "fe: next build" ($LASTEXITCODE -eq 0) "see $([System.IO.Path]::GetTempPath())"
        } finally {
            Pop-Location
        }

        # Boot the production server with smoke env and probe route/auth behaviour.
        Write-Host "  next start (port 3112) ..."
        $fePort = 3112
        $feLog = Join-Path $env:TEMP ("clara_nextstart_" + $PID + ".log")
        $proc = Start-Process -FilePath "C:\Program Files\nodejs\npx.cmd" -ArgumentList "next","start","-p","$fePort" -WorkingDirectory $webDir -RedirectStandardOutput $feLog -RedirectStandardError $feLog -PassThru -WindowStyle Hidden
        $env:SUPABASE_URL = $env:NEXT_PUBLIC_SUPABASE_URL
        $env:SUPABASE_ANON_KEY = $env:NEXT_PUBLIC_SUPABASE_ANON_KEY
        $env:CLARA_API_URL = if ($env:CLARA_API_URL) { $env:CLARA_API_URL } else { $adminBase }

        $feReady = $false
        for ($i = 0; $i -lt 30; $i++) {
            if ((Status-Of "http://localhost:$fePort/login") -gt 0) { $feReady = $true; break }
            Start-Sleep -Seconds 1
        }
        if (-not $feReady) {
            Check "fe: server startup" $false "next start did not come up; see $feLog"
        } else {
            Check "fe: server startup" $true
            Check "fe: /login returns 200" ((Status-Of "http://localhost:$fePort/login") -eq 200)
            Check "fe: / redirects to login (307)" ((Status-Of "http://localhost:$fePort/") -eq 307)
            Check "fe: protected page /ops redirects (307)" ((Status-Of "http://localhost:$fePort/ops") -eq 307)
            # BFF proxy must return 401 JSON for unauthenticated API calls, NOT
            # a 307 redirect (middleware matcher bug fixed in web/src/middleware.ts).
            $bff = Status-Of "http://localhost:$fePort/api/data/dashboard"
            Check "fe: BFF /api/data/dashboard unauthenticated -> 401" ($bff -eq 401) "HTTP $bff (307 = middleware hijack bug)"
            $bff2 = Status-Of "http://localhost:$fePort/api/data/clearing/records"
            Check "fe: BFF /api/data/clearing/records -> 401" ($bff2 -eq 401) "HTTP $bff2"
        }
        if ($proc -and -not $proc.HasExited) { Stop-Process -Id $proc.Id -Force -ErrorAction SilentlyContinue }
    } else {
        Check "fe: next build" $false "web/node_modules missing; run npm install in web/"
    }
}

Step "Results"
Yellow ("  summary: " + $script:pass + " passed, " + $script:fail + " failed")
foreach ($r in $script:results) {
    if ($r -match '^([^|]+)\|(PASS|FAIL)') {
        if ($Matches[2] -eq "PASS") { Green ("  PASS  " + $Matches[1]) } else { Red ("  FAIL  " + $Matches[1]) }
    }
}

if (-not $Keep) {
    Write-Host "  tearing down stack ..."
    & docker compose -f $compose down | Out-Null
} else {
    Yellow "  stack left running (use -Keep). Tear down: docker compose -f deploy\docker-compose.yml down"
}

if ($script:fail -gt 0) { Red "  SOME CHECKS FAILED"; exit 1 }
Green "  ALL CHECKS PASSED"
exit 0
