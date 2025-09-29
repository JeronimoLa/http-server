# 🌐 Modern Authentication & Authorization Technologies

Authentication and authorization have evolved significantly over the past decades. Here’s a comprehensive overview of the technologies in use today, their purposes, and trends.

---

## 1. **OAuth 2.0**
- **Purpose:** Authorization framework for delegated access.
- **Use Cases:** APIs, mobile apps, web apps.
- **Description:** Allows a user to grant a third-party app access to resources without sharing their password.
- **Modern Best Practice:** OAuth 2.1 (simplified, more secure version of OAuth2).
- **Examples:** Google APIs, GitHub API, Microsoft Graph.

---

## 2. **OpenID Connect (OIDC)**
- **Purpose:** Authentication layer built on top of OAuth2.
- **Use Cases:** Web login, Single Sign-On (SSO).
- **Description:** Provides ID Tokens (JWTs) containing user identity claims. Can be used in conjunction with OAuth2 for API access.
- **Examples:** Login with Google, Apple, GitHub; enterprise logins via Okta, Azure AD.

---

## 3. **Single Sign-On (SSO)**
- **Purpose:** User experience that allows one login across multiple applications.
- **Use Cases:** Enterprise applications, SaaS platforms, consumer apps.
- **Description:** Once authenticated with an Identity Provider (IdP), users can access all trusted applications without re-entering credentials.
- **Examples:** Okta, Ping Identity, Google Workspace SSO.

---

## 4. **SAML (Security Assertion Markup Language)**
- **Purpose:** Authentication standard (legacy, XML-based).
- **Use Cases:** Enterprise SSO, legacy systems.
- **Description:** Uses XML-based assertions to convey authentication info. Often replaced by OIDC in modern setups.
- **Examples:** Older Oracle, SAP, and enterprise web apps.

---

## 5. **FIDO2 / WebAuthn**
- **Purpose:** Passwordless authentication.
- **Use Cases:** Modern web and mobile apps.
- **Description:** Enables strong authentication using devices like security keys, biometrics, or platform authenticators. Works alongside OIDC/OAuth2.
- **Examples:** Passkeys, YubiKey, Touch ID, Windows Hello.

---

## 6. **Kerberos / LDAP**
- **Purpose:** Internal enterprise authentication.
- **Use Cases:** Corporate networks, Windows Active Directory environments.
- **Description:** Kerberos uses tickets to authenticate users on a network. LDAP provides directory services. Often bridged to SAML/OIDC for cloud access.

---

## 🔹 Modern Usage Summary

- **For web/mobile APIs →** OAuth2 + OIDC (standard today).
- **For enterprise SSO →** SAML (legacy, but still common) + OIDC (modern).
- **For the future of login →** WebAuthn / Passkeys (passwordless).
- **For internal corporate networks →** Kerberos, LDAP (still used).

---

### Key Takeaways
- **OAuth2**: Delegated access, not authentication by itself.
- **OIDC**: Adds authentication to OAuth2, enables login and SSO.
- **SSO**: User convenience across multiple apps.
- **SAML**: Older enterprise SSO standard, XML-based.
- **WebAuthn / Passkeys**: Modern passwordless authentication.
- **Kerberos / LDAP**: Traditional internal network authentication.

This ecosystem ensures secure, scalable, and flexible authentication and authorization across web, mobile, and enterprise environments.