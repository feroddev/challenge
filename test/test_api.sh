#!/bin/bash

echo "Testando POST do giroscópio..."
curl -X POST http://localhost:8080/telemetry/gyroscope \
  -H "Content-Type: application/json" \
  -d '{
    "x": 10.5,
    "y": -5.2,
    "z": 3.7,
    "timestamp": "2025-06-28T20:30:00Z",
    "device_id": "abc123"
  }'

echo -e "\n\nTestando POST do GPS..."
curl -X POST http://localhost:8080/telemetry/gps \
  -H "Content-Type: application/json" \
  -d '{
    "latitude": -23.5505,
    "longitude": -46.6333,
    "timestamp": "2025-06-28T20:30:00Z",
    "device_id": "abc123"
  }'

echo -e "\n\nTestando POST da foto..."
curl -X POST http://localhost:8080/telemetry/photo \
  -H "Content-Type: application/json" \
  -d '{
    "photo": "base64_encoded_string_here",
    "timestamp": "2025-06-28T20:30:00Z",
    "device_id": "abc123"
  }'

echo -e "\n\nTestando GET do giroscópio..."
curl -X GET http://localhost:8080/telemetry/gyroscope

echo -e "\n\nTestando GET do GPS..."
curl -X GET http://localhost:8080/telemetry/gps

echo -e "\n\nTestando GET da foto..."
curl -X GET http://localhost:8080/telemetry/photo
