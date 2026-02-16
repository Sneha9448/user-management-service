
echo "1. Login with Google (Mock)"
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query": "mutation { loginWithGoogle(idToken: \"mock_token\") { token user { id name email } } }"}')
echo "Login Response: $LOGIN_RESPONSE"

TOKEN=$(echo $LOGIN_RESPONSE | sed -e 's/.*"token":"\([^"]*\)".*/\1/')
echo "Extracted Token: $TOKEN"

echo ""
echo "2. Fetch current user (me) with token"
ME_RESPONSE=$(curl -s -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"query": "query { me { id name email } }"}')
echo "Me Response: $ME_RESPONSE"

echo ""
echo "3. Try protected mutation (createUser) WITHOUT token"
CREATE_FAIL_RESPONSE=$(curl -s -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query": "mutation { createUser(name: \"Auth Test\", email: \"auth@test.com\") { id name } }"}')
echo "Expected Failure Response: $CREATE_FAIL_RESPONSE"

echo ""
echo "4. Try protected mutation (createUser) WITH token"
CREATE_SUCCESS_RESPONSE=$(curl -s -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"query": "mutation { createUser(name: \"Auth Test\", email: \"auth@test.com\") { id name } }"}')
echo "Success Response: $CREATE_SUCCESS_RESPONSE"
