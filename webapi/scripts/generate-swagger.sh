#!/bin/bash

# Set environment variables for swagger generation
export API_HOST=${API_HOST:-"localhost:8001"}
export AUTH_SERVER_URL=${AUTH_SERVER_URL:-"http://localhost:8000"}

# Generate swagger docs
swag init -g main.go

# Replace placeholder values in swagger.json and swagger.yaml
sed -i "s/\${API_HOST}/$API_HOST/g" docs/swagger.json docs/swagger.yaml
sed -i "s/\${AUTH_SERVER_URL}/$AUTH_SERVER_URL/g" docs/swagger.json docs/swagger.yaml

echo "Swagger documentation generated successfully"