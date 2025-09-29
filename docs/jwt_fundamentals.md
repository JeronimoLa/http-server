# 🔐 JWT Fundamentals and Security

## 1. **What a JWT is**
- A JWT is a **compact, URL-safe string** used to represent claims between two parties.  
- Structure:  
  ```
  <header>.<payload>.<signature>
  ```
  - **Header** → describes metadata (algorithm, type).  
  - **Payload** → the claims (user ID, roles, expiration, etc.).  
  - **Signature** → cryptographic proof the token hasn’t been tampered with.

- Example:  
  ```
  eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.
  eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkFsaWNlIn0.
  TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ
  ```

## 2. **Encoding vs Encryption**
- **Header** and **Payload** are only **Base64URL encoded** → anyone can decode them.  
- They are *not secret*.  
- **Signature** ensures **integrity** (data wasn’t modified), not **confidentiality**.  
- If you need confidentiality → use JWE (JSON Web Encryption) or store minimal info in the JWT.

## 3. **How JWT is secured**
- The server signs the token using:
  - Symmetric key (HMAC, e.g. HS256) → same secret for signing and verifying.  
  - Asymmetric key (RSA/ECDSA, e.g. RS256/ES256) → private key signs, public key verifies.  

- When verifying, the server recomputes the signature. If it doesn’t match, token is invalid.  

✅ This prevents tampering.  
❌ This does **not** prevent reading claims.

## 4. **Why JWTs are useful in security**
- **Stateless authentication**: server doesn’t need to keep sessions in memory. The token itself contains user info + expiry.  
- **Scalable**: multiple services can trust the same token without central state.  
- **Standardized**: widely supported in APIs, OAuth2, OpenID Connect, etc.

## 5. **Security considerations**
1. **Confidentiality**  
   - Never put sensitive data (like passwords, SSNs) in the payload.  
   - Always use HTTPS to protect JWTs in transit.  

2. **Integrity**  
   - Always verify the signature.  
   - Never accept tokens with `alg":"none"`.  

3. **Expiration (`exp`)**  
   - Tokens should have short lifetimes.  
   - Use refresh tokens if you need longer sessions.  

4. **Revocation**  
   - JWTs are stateless → you can’t “log out” a token easily.  
   - Common mitigations: short TTLs, token blacklists, rotating secrets.  

5. **Storage**  
   - Store JWTs securely (e.g., HTTP-only cookies to mitigate XSS).  
   - Avoid localStorage for sensitive JWTs (vulnerable to XSS).  

6. **Audience / Issuer validation**  
   - Check `aud` (audience) and `iss` (issuer) claims to prevent token reuse across apps.  

## 6. **JWT in practice**
- **Login** → server verifies user, issues JWT with `sub=userID`, `exp=15m`.  
- **Client** → sends token in `Authorization: Bearer <token>` header.  
- **API** → validates signature, checks `exp`, extracts claims, grants access.  

---

## 🔑 Key takeaway
- JWT ≠ encrypted secret.  
- JWT = signed proof of claims (integrity, not confidentiality).  
- Security depends on:
  - Strong keys,  
  - Short expiration,  
  - Proper signature verification,  
  - Safe transport/storage.  
