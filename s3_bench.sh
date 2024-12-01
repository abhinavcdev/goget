#!/bin/bash

# S3 URL to download - use a public bucket for testing
S3_URL="s3://aws-cli-workshop-akumar/sample-file.txt"
S3_REGION="us-east-1"

# Create a directory for test downloads
mkdir -p test_downloads
cd test_downloads

# Clean up any previous downloads
rm -f goget_s3_*.txt wget_s3_download.txt

echo "=== Testing goget with S3 ==="
# Test with different concurrency values
for c in 1 4 8 16; do
  echo "Testing with concurrency: $c"
  time ../goget -c $c -o "goget_s3_c${c}_download.txt" --s3-region "$S3_REGION" "$S3_URL"
  echo ""
done

echo "=== Testing wget with S3 ==="
# For wget, we need to use the AWS CLI to generate a presigned URL
PRESIGNED_URL=$(aws s3 presign "$S3_URL" --region "$S3_REGION")
if [ $? -eq 0 ]; then
  time wget -O wget_s3_download.txt "$PRESIGNED_URL"
else
  echo "Failed to generate presigned URL with AWS CLI. Make sure AWS CLI is installed and configured."
fi
echo ""

# Verify all downloads have the same size
echo "=== File size verification ==="
ls -lh goget_s3_c*_download.txt wget_s3_download.txt 2>/dev/null

# Return to original directory
cd ..
