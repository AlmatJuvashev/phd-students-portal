#!/bin/bash

# Setup MinIO for local development
# This script creates the bucket and sets CORS policy

echo "⏳ Waiting for MinIO to be ready..."
sleep 5

# Install mc (MinIO Client) if not installed
if ! command -v mc &> /dev/null; then
    echo "📦 Installing MinIO Client..."
    brew install minio/stable/mc || {
        echo "❌ Failed to install mc. Please install manually: https://min.io/docs/minio/linux/reference/minio-mc.html"
        exit 1
    }
fi

# Configure MinIO client
echo "🔧 Configuring MinIO client..."
mc alias set local http://localhost:9000 minioadmin minioadmin

# Create bucket if it doesn't exist
echo "📁 Creating bucket 'phd-portal'..."
mc mb local/phd-portal --ignore-existing

# Set CORS policy for bucket
echo "🌐 Setting CORS policy..."
cat > /tmp/minio-cors.json <<EOF
{
  "CORSRules": [
    {
      "AllowedOrigins": ["http://localhost:5173", "http://localhost:5174"],
      "AllowedMethods": ["GET", "PUT", "POST", "DELETE", "HEAD"],
      "AllowedHeaders": ["*"],
      "ExposeHeaders": ["ETag", "Content-Length", "Content-Type"]
    }
  ]
}
EOF

mc anonymous set-json /tmp/minio-cors.json local/phd-portal

# Set public download policy (for presigned URLs to work)
echo "🔓 Setting bucket policy for presigned URLs..."
mc anonymous set download local/phd-portal

echo "✅ MinIO setup complete!"
echo ""
echo "📊 MinIO Console: http://localhost:9091"
echo "   Username: minioadmin"
echo "   Password: minioadmin"
echo ""
echo "🪣 Bucket 'phd-portal' is ready for file uploads!"
