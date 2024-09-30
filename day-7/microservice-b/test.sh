#!/bin/bash

# Set the base URL of your Node.js application
BASE_URL="http://localhost:8081"

echo $BASE_URL

# Define an array of endpoints
ENDPOINTS=(
  "/hello-b"
  "/call-a"
  "/getme-coffee"
)

# Function to make a random request to one of the endpoints
make_random_request() {
  local endpoint=${ENDPOINTS[$RANDOM % ${#ENDPOINTS[@]}]}
  curl -s -o /dev/null -w "%{http_code}" "$BASE_URL$endpoint"
}

# Make 1000 random requests
for ((i=1; i<=1000; i++)); do
  make_random_request
  echo "Request $i completed"
  sleep 0.1  # Optional: Sleep for a short duration between requests to simulate real traffic
done

echo "Completed 1000 requests"
