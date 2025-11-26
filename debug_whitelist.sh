#!/bin/bash

echo "===== Debugging Whitelist Endpoint ====="
echo ""

# Register user
echo "1️⃣  Registering test user..."
RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email":"debug'$(date +%s)'@uii.ac.id",
    "password":"testpass123",
    "full_name":"Debug User",
    "phone_number":"081234567890"
  }')

SUCCESS=$(echo "$RESPONSE" | jq -r '.success')
if [ "$SUCCESS" != "true" ]; then
  echo "❌ Registration failed!"
  echo "$RESPONSE" | jq '.'
  exit 1
fi

TOKEN=$(echo "$RESPONSE" | jq -r '.data.token')
USER_ID=$(echo "$RESPONSE" | jq -r '.data.user.id')
echo "✅ User registered"
echo "   User ID: $USER_ID"
echo "   Token: ${TOKEN:0:30}..."
echo ""

# Create valid PDF
echo "2️⃣  Creating test PDF file..."
cat > /tmp/test_whitelist.pdf << 'EOF'
%PDF-1.4
1 0 obj
<<
/Type /Catalog
/Pages 2 0 R
>>
endobj
2 0 obj
<<
/Type /Pages
/Kids [3 0 R]
/Count 1
>>
endobj
3 0 obj
<<
/Type /Page
/Parent 2 0 R
/MediaBox [0 0 612 792]
/Contents 4 0 R
>>
endobj
4 0 obj
<<
/Length 44
>>
stream
BT
/F1 12 Tf
100 700 Td
(Test Document) Tj
ET
endstream
endobj
xref
0 5
0000000000 65535 f 
0000000009 00000 n 
0000000058 00000 n 
0000000115 00000 n 
0000000214 00000 n 
trailer
<<
/Size 5
/Root 1 0 R
>>
startxref
306
%%EOF
EOF

ls -lh /tmp/test_whitelist.pdf
echo "✅ PDF created"
echo ""

# Check storage permissions
echo "3️⃣  Checking storage directory..."
if [ -d "storage/documents" ]; then
  ls -ld storage/documents
  echo "✅ Storage directory exists"
else
  echo "⚠️  Creating storage/documents..."
  mkdir -p storage/documents
fi
echo ""

# Submit whitelist request with verbose output
echo "4️⃣  Submitting whitelist request..."
echo "   Organization: Test Organization"
echo "   Document: /tmp/test_whitelist.pdf"
echo ""

# Use -v for verbose and capture everything
FULL_OUTPUT=$(curl -v -X POST http://localhost:8080/api/v1/whitelist/request \
  -H "Authorization: Bearer $TOKEN" \
  -F "organization_name=Test Organization" \
  -F "document=@/tmp/test_whitelist.pdf" 2>&1)

# Extract HTTP status
HTTP_CODE=$(echo "$FULL_OUTPUT" | grep "< HTTP" | awk '{print $3}')

echo "HTTP Status Code: $HTTP_CODE"
echo ""

# Show response body
RESPONSE_BODY=$(echo "$FULL_OUTPUT" | sed -n '/^{/,/^}/p')

if [ -n "$RESPONSE_BODY" ]; then
  echo "Response Body:"
  echo "$RESPONSE_BODY" | jq '.' 2>/dev/null || echo "$RESPONSE_BODY"
else
  echo "⚠️  No response body (empty)"
fi

echo ""
echo "===== Full Curl Output ====="
echo "$FULL_OUTPUT"
echo ""

# Check result
if [ "$HTTP_CODE" == "201" ]; then
  echo "✅ SUCCESS! Whitelist request submitted"
elif [ "$HTTP_CODE" == "500" ]; then
  echo "❌ FAILED with 500 Internal Server Error"
  echo ""
  echo "🔍 Possible causes:"
  echo "   1. Check server terminal for error stack trace"
  echo "   2. Database connection issue"
  echo "   3. File upload permissions"
  echo "   4. Missing columns in whitelist_requests table"
  echo ""
  echo "💡 Next steps:"
  echo "   - Look at server logs in the terminal running 'go run cmd/api/main.go'"
  echo "   - Check storage/documents/ permissions"
  echo "   - Verify database schema"
else
  echo "❌ FAILED with HTTP $HTTP_CODE"
fi

# Cleanup
rm -f /tmp/test_whitelist.pdf
