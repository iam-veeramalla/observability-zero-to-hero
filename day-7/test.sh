#!/bin/bash

# Check if both load balancers are provided as input
if [ $# -ne 2 ]; then
    echo "Usage: $0 <LB-1 DNS> <LB-2 DNS>"
    exit 1
fi

# Assign input arguments to variables
LB1=$1
LB2=$2

# Define available routes for LB1 and LB2
LB1_ROUTES=("/call-b" "/hello-a" "/getme-coffee")
LB2_ROUTES=("/call-a" "/hello-b" "/getme-coffee")

# Function to generate random index and request from LB1
request_lb1() {
    RANDOM_INDEX=$((RANDOM % ${#LB1_ROUTES[@]}))
    URL="$LB1${LB1_ROUTES[$RANDOM_INDEX]}"
    echo "Sending request to LB1: $URL"
    curl -s -o /dev/null -w "%{http_code}" $URL
}

# Function to generate random index and request from LB2
request_lb2() {
    RANDOM_INDEX=$((RANDOM % ${#LB2_ROUTES[@]}))
    URL="$LB2${LB2_ROUTES[$RANDOM_INDEX]}"
    echo "Sending request to LB2: $URL"
    curl -s -o /dev/null -w "%{http_code}" $URL
}

# Loop for sending requests to both LBs randomly
while true; do
    # Randomly choose between LB1 and LB2
    if (( RANDOM % 2 == 0 )); then
        request_lb1
    else
        request_lb2
    fi

    # Sleep for 1 second between requests (adjust if needed)
    sleep 1
done
