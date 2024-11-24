#!/bin/bash

# Install wrk if not installed
if ! command -v wrk &> /dev/null; then
  echo "wrk not found, installing..."
  sudo apt update && sudo apt install -y wrk || brew install wrk
fi

# Target URL
URL="http://localhost:8080/api/verve/accept"


# Number of threads and connections
THREADS=8
CONNECTIONS=1000

# Generate Lua script for dynamic IDs
cat <<EOL > dynamic_ids.lua
wrk.method = "GET"
wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"

counter = 0

request = function()
  counter = counter + 1
  local path = "/api/verve/accept?id=" .. counter
  return wrk.format(nil, path)
end
EOL

# Run wrk with the Lua script
wrk -t$THREADS -c$CONNECTIONS -d30s -s dynamic_ids.lua $URL

# Cleanup
rm -f dynamic_ids.lua
