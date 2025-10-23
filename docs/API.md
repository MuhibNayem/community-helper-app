CommunityConnect Backend API Specification
=========================================

Overview
--------
- **Audience**: Backend (Go) developers, frontend/mobile team, QA engineers, integration partners.  
- **Architecture**: RESTful JSON APIs over HTTPS, stateless except for authenticated sessions (JWT).  
- **Versioning**: Prefix all endpoints with `/v1/`. Breaking changes trigger a new version (`/v2/`).  
- **Authentication**: Firebase-issued ID tokens exchanged for backend JWT session tokens; service-to-service requests use signed HMAC headers.  
- **Idempotency**: All state-changing endpoints accept `Idempotency-Key` header to prevent duplicate processing.  
- **Rate Limiting**: Per-user and per-IP limits via API gateway; `429` responses include `Retry-After`.  

Auth & Session Management
-------------------------

### Request OTP
- `POST /v1/auth/otp/request`  
- **Purpose**: Initiate phone verification via SMS.  
- **Body**:
  ```json
  {
    "phone": "+8801XXXXXXXXX",
    "locale": "bn_BD",
    "channel": "sms" // optional: "call" fallback
  }
  ```
- **Responses**:
  - `200 OK`: `{ "expiresIn": 120, "attemptsRemaining": 3 }`
  - `400 Bad Request`: Invalid phone, throttled.
  - `409 Conflict`: Number blocked/banned.  
- **Notes**: Delegates to Twilio; apply per-number rate limits and reCAPTCHA check.

### Verify OTP
- `POST /v1/auth/otp/verify`
- **Purpose**: Confirm OTP, create/return session JWT.  
- **Body**:
  ```json
  { "phone": "+8801XXXXXXXXX", "otp": "123456", "deviceId": "abc123" }
  ```
- **Responses**:
  - `200 OK`:
    ```json
    {
      "sessionToken": "jwt",
      "refreshToken": "jwt",
      "user": { "id": "uid", "isHelper": true, "onboardingStatus": "COMPLETED" }
    }
    ```
  - `401 Unauthorized`: Invalid/expired code.
  - `423 Locked`: Too many failures; cooldown.  

### Refresh Session
- `POST /v1/auth/refresh`
- Body: `{ "refreshToken": "jwt" }`
- Response: `200 OK` with new `sessionToken`.  
- Errors: `401 Unauthorized` (invalid), `409 Conflict` (revoked).  

### Logout
- `POST /v1/auth/logout`
- Body: `{ "deviceId": "abc123" }`
- Response: `204 No Content`.

Users & Profiles
----------------

### Get Current User
- `GET /v1/users/me`
- Response:
  ```json
  {
    "id": "uid",
    "phone": "+8801...",
    "name": "Rafi",
    "photoUrl": "https://...",
    "language": "bn",
    "isHelper": true,
    "helperStatus": "ACTIVE",
    "notificationPrefs": { "quietHours": { "start": "22:00", "end": "07:00" }, "urgentSms": true },
    "kycStatus": "PENDING"
  }
  ```

### Update Profile
- `PATCH /v1/users/me`
- Body schema supports partial updates:
  ```json
  { "name": "Rafi Rahman", "language": "en", "photoUrl": "https://..." }
  ```
- Response: `200 OK` with updated user.  
- Validations: name length, supported languages, image URL domain.

### Update Helper Toggle
- `POST /v1/helpers/me/toggle`
- Body: `{ "optedIn": true }`
- Response: `200 OK` with helper summary.  
- Side effects: If opting in, check KYC status; if insufficient, respond `409` with required actions.

### Update Skills
- `PUT /v1/helpers/me/skills`
- Body: `{ "skills": ["MEDICAL_FIRST_AID", "MECHANICAL"] }`
- Response: `200 OK`.  
- Constraints: At least one skill; skills validated against admin-defined catalog.

### Manage Availability
- `PUT /v1/helpers/me/availability`
- Body:
  ```json
  {
    "weekly": [
      { "day": "MONDAY", "start": "09:00", "end": "18:00" }
    ],
    "exceptions": [
      { "date": "2025-02-20", "slots": [] }
    ]
  }
  ```
- Response: `200 OK`.  
- Edge cases: Overlapping slots (reject), timezone assumed local (Asia/Dhaka).

### Upload Documents (KYC)
- `POST /v1/helpers/me/kyc`
- Multipart form with `documentType` and file.  
- Response: `202 Accepted` (under review).  
- Security: Virus scanning; limit to allowed formats/PDF/JPEG.

Request Lifecycle APIs
----------------------

### Create Help Request
- `POST /v1/requests`
- Body:
  ```json
  {
    "type": "URGENT",
    "category": "MEDICAL_FIRST_AID",
    "description": "Need help with...",
    "location": { "lat": 23.78, "lng": 90.36, "address": "Dhaka", "placeId": "xyz" },
    "scheduledFor": null,
    "attachments": ["gs://..."]
  }
  ```
- Response: `201 Created` with request object.  
- Validations: category allowed, location present, scheduledFor required for planned.  
- Rate limit: max active urgent request per seeker.

### Get My Requests
- `GET /v1/requests?status=ACTIVE&limit=20&offset=0`
- Response: list with pagination meta. Includes active and history.

### Get Request Detail
- `GET /v1/requests/{requestId}`
- Returns full request info, match status, assigned helper details (limited PII).

### Cancel Request
- `POST /v1/requests/{requestId}/cancel`
- Body: `{ "reason": "HELPER_NOT_NEEDED" }`
- Response: `200 OK` with updated status.  
- Edge cases: If request already `IN_PROGRESS`, require helper consent or admin override.

### Rate Helper
- `POST /v1/requests/{requestId}/rate`
- Body: `{ "rating": 5, "comment": "Great help!" }`
- Response: `201 Created`.  
- Validation: only after completion, one rating per request.

Matching & Invitations
----------------------

### List Invitations (Helper Inbox)
- `GET /v1/matches?status=PENDING`
- Response:
  ```json
  [
    {
      "matchId": "m123",
      "requestId": "r456",
      "type": "URGENT",
      "category": "MEDICAL_FIRST_AID",
      "distanceKm": 1.2,
      "etaMinutes": 5,
      "autoDeclineAt": "2025-02-16T08:15:00Z"
    }
  ]
  ```

### Accept Invitation
- `POST /v1/matches/{matchId}/accept`
- Response: `200 OK` with session detail.  
- Side effects: Payment intent created, request status transitions to `ACCEPTED`.  
- Edge cases: If helper already booked, return `409 Conflict`.

### Decline Invitation
- `POST /v1/matches/{matchId}/decline`
- Body: `{ "reason": "BUSY" }`  
- Response: `200 OK`.

### Update Arrival/Progress
- `POST /v1/matches/{matchId}/status`
- Body:
  ```json
  { "status": "ARRIVED", "location": { "lat": 23.78, "lng": 90.36 } }
  ```
- Allowed statuses: `EN_ROUTE`, `ARRIVED`, `IN_PROGRESS`, `BLOCKED`.  
- Response: `200 OK`; notifies seeker.

### Complete Session Confirmation
- `POST /v1/matches/{matchId}/complete`
- Body: `{ "confirmation": "SUCCESS" }` or `{ "confirmation": "FAILED", "reason": "...", "evidence": ["gs://..."] }`
- Response: `200 OK`; triggers payout or dispute workflow.

Chat & Messaging
----------------

### Fetch Threads
- `GET /v1/threads?status=ACTIVE`
- Returns list of chat threads where user participates.  
- Response fields: `threadId`, `matchId`, `lastMessage`, `unreadCount`.

### Post Message
- `POST /v1/threads/{threadId}/messages`
- Body:
  ```json
  { "type": "TEXT", "body": "On my way!", "clientTimestamp": "2025-02-16T08:05:00Z" }
  ```
- Response: `201 Created`.  
- Attachments handled via pre-signed upload to Storage; message references file URL.

### Mark Thread Read
- `POST /v1/threads/{threadId}/read`
- Body: `{ "lastMessageId": "msg789" }`  
- Response: `204 No Content`.

### Delete Message (User-level Hide)
- `DELETE /v1/threads/{threadId}/messages/{messageId}`
- Soft delete; respect retention policy.  
- Response: `204 No Content`.  
- Restrictions: cannot delete system messages or messages involved in disputes.

Payments & Escrow
-----------------

### Attach Payment Method
- `POST /v1/payments/methods`
- Body: Stripe `paymentMethodId` from client SDK.
- Response: `201 Created` with masked details.  
- Validations: Seekers only; limit to 5 methods.

### List Payment Methods
- `GET /v1/payments/methods`
- Response: array of saved methods.

### Remove Payment Method
- `DELETE /v1/payments/methods/{id}`
- Response: `204 No Content`.  
- Restriction: cannot remove method tied to pending request.

### View Transaction
- `GET /v1/transactions/{transactionId}`
- Response includes escrow status, breakdown, timeline.  

### Issue Dispute
- `POST /v1/transactions/{transactionId}/dispute`
- Body: `{ "reason": "SERVICE_NOT_COMPLETED", "details": "...", "evidence": ["gs://..."] }`
- Response: `202 Accepted` (under review).

### Admin Resolve Dispute
- `POST /v1/transactions/{transactionId}/resolve`
- Body:
  ```json
  {
    "outcome": "REFUND_FULL",
    "notes": "Helper no-show",
    "refundAmount": 1200
  }
  ```
- Authentication: Admin role.  
- Response: `200 OK` with updated transaction and triggered Stripe refund/payout adjustments.

Impact & Gamification
---------------------

### Dashboard Metrics
- `GET /v1/impact/me`
- Response:
  ```json
  {
    "totalPeopleHelped": 24,
    "impactScore": 1280,
    "streakDays": 5,
    "nextAvailability": "2025-02-18T10:00:00Z",
    "badgeSummaries": [
      { "badgeId": "COMMUNITY_HERO", "earnedAt": "2025-01-10T..." }
    ]
  }
  ```

### Badge Catalogue (Public)
- `GET /v1/badges`
- Lists available badges, criteria, assets.

### Referral Links
- `POST /v1/referrals`
- Response: `{ "code": "HELP123", "shareUrl": "https://..." }`

### Claim Referral
- `POST /v1/referrals/{code}/claim`
- Body: optional metadata.  
- Response: `200 OK` with reward status (`PENDING`, `APPLIED`).  
- Edge case: Code expired → `410 Gone`.

Notifications
-------------

### Register Device Token
- `POST /v1/notifications/devices`
- Body: `{ "platform": "ANDROID", "token": "fcmToken", "deviceId": "abc123" }`
- Response: `201 Created`.  
- Handles deduping, toggling quiet hours.

### Update Notification Preferences
- `PATCH /v1/notifications/preferences`
- Body: `{ "urgentSms": false, "quietHours": { "start": "22:00", "end": "06:00" } }`
- Response: `200 OK`.

### Fetch System Messages
- `GET /v1/notifications/system?since=2025-02-01T00:00:00Z`
- Returns announcements, maintenance notices.

Admin & Operations
------------------

### Admin Login
- `POST /v1/admin/login`
- Body: `{ "email": "ops@communityconnect.app", "password": "..." }`
- Response: `200 OK` with admin JWT, scopes.

### User Search
- `GET /v1/admin/users?query=phone:+8801`
- Response: paginated list with status, flags.

### Force Logout
- `POST /v1/admin/users/{userId}/logout`
- Response: `204 No Content`.

### View Active Requests
- `GET /v1/admin/requests?status=MATCHING`
- Includes location snapshots, assigned helpers, escalations.

### Manual Reassignment
- `POST /v1/admin/requests/{requestId}/assign`
- Body: `{ "helperId": "uid" }`
- Response: `200 OK`.

### System Metrics
- `GET /v1/admin/metrics`
- Returns aggregated KPIs (match rate, response time, SMS failures).  
- Protected via admin scope `metrics:read`.

Analytics & Reporting
---------------------

### Event Export Trigger
- `POST /v1/analytics/export`
- Authenticated via scheduler service account.  
- Response: `202 Accepted`; job ID returned.

### Heatmap Data (Aggregated)
- `GET /v1/analytics/heatmap?from=2025-02-01&to=2025-02-15`
- Returns aggregated location density (anonymized) for product team.  
- Requires `analytics:read` scope.

Support & Ticketing
-------------------

### Create Support Ticket
- `POST /v1/support/tickets`
- Body:
  ```json
  {
    "category": "PAYMENT",
    "requestId": "r123",
    "description": "Payment stuck"
  }
  ```
- Response: `201 Created` with ticket ID.

### Update Ticket Status (Ops)
- `POST /v1/support/tickets/{ticketId}/status`
- Body: `{ "status": "RESOLVED", "notes": "Refund processed" }`
- Response: `200 OK`.

System Hooks & Webhooks
-----------------------

### Stripe Webhook
- `POST /v1/hooks/stripe`
- Validates signature; handles `payment_intent.succeeded`, `charge.dispute.created`, etc.  
- Responds `200 OK` on success, `400` if signature invalid (no retries on 2xx).

### Twilio Webhook
- `POST /v1/hooks/twilio`
- Handles delivery receipts; updates notification status; signature validation required.

### Monitoring Hook
- `POST /v1/hooks/healthcheck`
- Auth: internal token.  
- Usage: External uptime monitors call to verify system health.  
- Response: `{ "status": "ok", "uptime": "...", "dependencies": [] }`.

Error Handling
--------------
- **Error Response Format**:
  ```json
  {
    "error": {
      "code": "REQUEST_NOT_FOUND",
      "message": "Help request not found",
      "details": [{ "field": "requestId", "issue": "INVALID" }]
    }
  }
  ```
- **Standard Codes**:
  - `AUTH_FAILED`, `SESSION_EXPIRED`, `PHONE_BLOCKED`
  - `REQUEST_LIMIT_REACHED`, `HELPER_NOT_ELIGIBLE`, `MATCH_ALREADY_CLAIMED`
  - `PAYMENT_REQUIRED`, `PAYMENT_FAILED`, `ESCROW_HELD`
  - `DISPUTE_ALREADY_OPEN`, `BADGE_NOT_FOUND`
  - `VALIDATION_ERROR`, `RATE_LIMITED`, `INTERNAL_SERVER_ERROR`

Security Considerations
-----------------------
- All endpoints require HTTPS with TLS 1.2+.  
- JWTs signed with rotating keys (JWKS endpoint exposed for clients).  
- Access tokens include scopes (`user`, `helper`, `admin`).  
- Audit logging for critical operations (payments, disputes, admin overrides).  
- Input validation and sanitation across all endpoints; use request context with deadlines to prevent hanging.  
- CSRF protection for admin portal (cookie-based).  

Appendix: State Machines
------------------------

### Match Invitation
```
PENDING -> ACCEPTED (helper accept)
PENDING -> DECLINED (helper decline)
PENDING -> TIMEOUT (expires)
ACCEPTED -> IN_PROGRESS (helper en route)
IN_PROGRESS -> COMPLETED (both confirm)
IN_PROGRESS -> FAILED (helper reports failure)
```

### Transaction
```
PENDING -> AUTHORIZED (intent created)
AUTHORIZED -> HELD (service in progress)
HELD -> RELEASED (completion confirmed)
HELD -> REFUNDED (refund issued)
HELD -> DISPUTED (issue raised)
DISPUTED -> RESOLVED (admin outcome)
```

Future API Extensions
---------------------
- WebSocket/SSE for real-time updates (match status, chat).  
- gRPC endpoints for heavy internal services (matching, analytics).  
- Public API keys for partner integrations (e.g., NGO dashboards).  
- Batch APIs for operations (bulk notifications, helper availability import).

