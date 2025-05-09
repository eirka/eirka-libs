# Secret Management Migration Guide

This document provides guidance on migrating from the legacy approach to the new config-based JWT secret management.

## Background

The previous implementation used a global `Secret` variable that was set from the main application. This approach had several limitations:

1. No validation of secret strength
2. No ability to rotate secrets without downtime
3. Limited error handling
4. No atomic operations when working with the secret

The new implementation addresses these issues by:
1. Moving JWT secrets into the config file
2. Supporting both "old" and "new" secrets for smooth rotation
3. Always using the "new" secret for signing tokens while validating against both

## Configuration Changes

The session secret configuration has moved from individual microservice config files to the central `/etc/pram/pram.conf` file:

```json
{
  "Session": {
    "OldSecret": "",
    "NewSecret": "your-jwt-signing-secret"
  }
}
```

### Key Features

- The `NewSecret` is used for signing all new tokens and as the primary validation secret
- The `OldSecret` is used as a fallback for validation during secret rotation
- Secret validation now enforces minimum length requirements (16 chars in production)
- Thread-safe operations for all secret access

## Migration Steps

### Step 1: Update Config File

Add the `Session` section to your `/etc/pram/pram.conf` file:

```json
{
  "Session": {
    "OldSecret": "",
    "NewSecret": "your-current-jwt-secret"
  }
}
```

### Step 2: Remove Session Struct from Microservices

Remove the `Session` struct from your microservice configs, as the JWT secrets will now be managed by the `eirka-libs/config` package.

**Old approach in microservice config:**
```go
// Session holds secret for JWT key
type Session struct {
    Secret string
}
```

This section should be removed as it's now handled by `eirka-libs`.

### Step 3: Remove Secret Setting in Microservices

Remove any code that manually sets the JWT secret:

**Old approach:**
```go
// Don't do this anymore
user.Secret = config.Settings.Session.Secret
```

This is no longer needed as the JWT secret is loaded directly from the central config file.

## Secret Rotation Workflow

When you need to rotate to a new secret:

1. **Start Rotation:** Update the config file to move the current secret to `OldSecret` and set a new value for `NewSecret`:
   ```json
   {
     "Session": {
       "OldSecret": "your-previous-jwt-secret",
       "NewSecret": "your-new-jwt-secret"
     }
   }
   ```

2. **Deploy Changes:** Deploy the updated configuration to all services.
   - All new tokens will be signed with `NewSecret`
   - Existing tokens signed with `OldSecret` will continue to work

3. **Complete Rotation:** After all tokens using the old secret have expired (typically 90 days), you can clear the `OldSecret`:
   ```json
   {
     "Session": {
       "OldSecret": "",
       "NewSecret": "your-new-jwt-secret"
     }
   }
   ```

This process allows for zero-downtime rotation of JWT secrets.

## Security Recommendations

1. Use a strong random secret of at least 32 characters
2. Store secrets securely in the config file with appropriate permissions
3. Rotate secrets periodically (every 90-180 days)
4. Ensure the JWT expiration period is shorter than your rotation cycle
5. Follow a systematic process for secret rotation across all services