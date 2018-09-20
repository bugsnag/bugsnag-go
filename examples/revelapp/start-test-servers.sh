set -e

echo "Starting test server..."
json-server -p 1234 payloads.json
echo "Shutting down test server..."
