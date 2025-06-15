#!/bin/bash

echo "Starting HTTP server on http://localhost:8000"
echo "Open your browser and navigate to the URL above to play Island Merge"
echo "Press Ctrl+C to stop the server"

cd web && python3 -m http.server 8000