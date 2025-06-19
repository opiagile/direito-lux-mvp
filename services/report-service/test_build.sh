#!/bin/bash
cd /Users/franc/Opiagile/SAAS/direito-lux/services/report-service
echo "Building Report Service..."
go mod tidy
go build -o report-service ./cmd/server
echo "Build completed successfully!"