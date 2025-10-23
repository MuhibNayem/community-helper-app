Community Helper Technical Documentation
========================================

Version History
---------------
- v0.1 (2025-02-16): Initial draft compiled from proposal and SRS requirements.

1. System Overview
------------------
Community Helper is an Android-first platform that connects individuals needing assistance with nearby helpers. The solution combines:
- Flutter client (FlutterFlow + custom code) for user interfaces.
- Firebase backend providing authentication, data storage, messaging, and serverless logic.
- Third-party integrations for SMS (Twilio), payments (Stripe), and geospatial services (Google Maps/Places).

Primary goals:
- Support urgent and scheduled help requests.
- Provide trusted matching, secure payments, and transparent impact tracking.
- Maintain high performance and reliability on budget Android devices.

2. Architecture Breakdown
-------------------------

2.1 High-Level Components
- **Mobile Client (Flutter)**  
  Handles onboarding, request creation, matching visualizations, chat, dashboard, and payments UI. Generates and consumes app state via Firestore streams and Cloud Function endpoints.

- **Firebase Services**  
  - Authentication: Phone number/OTP login.  
  - Firestore: NoSQL database for users, requests, sessions, chat, transactions.  
  - Storage: Media assets (profile photos, attachments).  
  - Cloud Functions: Business logic (matching, notifications, payments, analytics jobs).  
  - Cloud Messaging (FCM): Push notifications.  
  - Remote Config: Feature flags, experiment toggles.  
  - App Check: Protects backend from abuse.  

- **Third-Party Integrations**  
  - Twilio Programmable SMS: Urgent alerts and OTP fallback.  
  - Stripe Connect: Escrow processing, payouts to helpers.  
  - Google Maps SDK, Directions API, Places API: Location display, routing, address autocompletion.  
  - Lottie: Celebration animations.  
  - BigQuery/Data Studio (optional): Analytics reporting.

2.2 Environment Topology
```
Android Client (Dev)  ->  Firebase Dev Project  ->  External Sandboxes (Stripe test, Twilio test)
Android Client (Stage)->  Firebase Stage Project->  External Stage Accounts
Android Client (Prod) ->  Firebase Prod Project ->  Stripe Live, Twilio Live
```
- Separate Google Cloud projects per environment.  
- Configuration managed via `.env` files and Flutter build flavors (`dev`, `stage`, `prod`).  
- Remote Config/Feature flags allow toggling features without redeploy.  
- Secret keys stored in Firebase Secret Manager; Cloud Functions load via environment variables.

2.3 Deployment Pipeline
- Source hosted on GitHub.  
- GitHub Actions pipeline stages: lint → unit tests → widget/integration tests → Flutter build → deploy Cloud Functions → distribute APK/AAB to Firebase App Distribution (QA) or Play Console (prod).  
- Automated semantic version bumping for release branches.  
- Release candidate gating: QA sign-off + stakeholders.  
- Post-deploy verification includes smoke tests and monitoring dashboards.

3. Module Design
----------------

3.1 Client Modules
- `core/`  
  Shared utilities (logging, analytics events, app configuration, localization, theme tokens).
- `auth/`  
  OTP entry, helper opt-in modal, onboarding guard, account recovery.
- `profile/`  
  Profile editing, skill selection, availability scheduler, verification status.
- `requests/`  
  Creation flows (urgent/planned), request detail screens, cancellation.
- `matching/`  
  Radar visualization, helper list, matching status states, fallback messaging.
- `chat/`  
  Conversation UI, attachments, message state management, retention enforcement.
- `payments/`  
  Payment method capture (Stripe SDK), escrow status, receipts, dispute submission.
- `dashboard/`  
  Impact metrics, badges, mind-map visualization, referral flows.
- `notifications/`  
  Push handler, local notifications scheduling, quiet hours.
- `admin/` (future)  
  lightweight operations tooling if built into mobile client for support staff.

3.2 Cloud Functions Modules
- `auth`  
  Post-auth triggers, KYC initialization, abuse detection.
- `requests`  
  Matching engine, SLA timers, escalation to SMS, state transitions.
- `payments`  
  Stripe webhook handler, payout scheduler, ledger reconciliation.
- `chat`  
  Retention cleanup (scheduled), content moderation hooks.
- `analytics`  
  Metric aggregation, badge issuance, BigQuery export.
- `admin`  
  Slack/email alert dispatch, audit logging, manual override endpoints.

3.3 Integration Edge Services
- Twilio: Hosted functions for fallback messaging if Cloud Functions exceed throughput.  
- Stripe: Connected accounts management, manual review dashboards for compliance.  
- Monitoring: Stackdriver logs + third-party alerting (PagerDuty/Slack) configured via Cloud Monitoring.

4. Data Model Specification
---------------------------

### 4.1 Firestore Collections

`users`  
- `id`: string (UID)  
- `phone`: string  
- `name`: string  
- `roles`: { `isHelper`: bool, `isAdmin`: bool }  
- `helperStatus`: { `optedIn`: bool, `skills`: [string], `availability`: schedule object, `rating`: float, `ratingCount`: int, `badges`: [string], `kycStatus`: enum }  
- `preferences`: { `language`: enum, `notifications`: { `quietHours`: range, `urgentSms`: bool } }  
- `location`: { `lastKnown`: geopoint, `accuracy`: meters, `updatedAt`: timestamp }  
- `createdAt`, `updatedAt`

`helperAvailability` (optional if schedule large)  
- `userId`  
- `weeklySlots`: array of { day, startTime, endTime }  
- `exceptions`: { date: [time ranges] }  

`helpRequests`  
- `id`  
- `requesterId`  
- `type`: enum (`URGENT`, `PLANNED`)  
- `status`: enum (`DRAFT`, `SUBMITTED`, `MATCHING`, `ACCEPTED`, `IN_PROGRESS`, `COMPLETED`, `CANCELLED`, `DISPUTED`)  
- `category`: string (from taxonomy)  
- `description`: string  
- `attachments`: [storage path]  
- `location`: { `geopoint`, `address`, `placeId`, `accuracy` }  
- `scheduledFor`: timestamp (for planned)  
- `createdAt`, `updatedAt`, `expiresAt`  
- `sla`: { `matchDeadline`, `completionDeadline` }  
- `pricing`: { `estimatedAmount`, `currency`, `platformFee` }  
- `cancellation`: { `reason`, `initiator`, `timestamp`, `penaltyApplied` }  

`matchSessions`  
- `id`  
- `requestId`  
- `helperId`  
- `status`: enum (`INVITED`, `ACCEPTED`, `DECLINED`, `TIMED_OUT`, `CANCELLED`, `COMPLETED`)  
- `invitedAt`, `respondedAt`, `acceptedAt`, `arrivedAt`, `completedAt`  
- `eta`: seconds  
- `smsSent`: bool  
- `pushSent`: bool  
- `reason`: string (for decline/cancel)  
- `metrics`: { `distance`: km, `travelTime`: seconds }  

`chatThreads`  
- `id`  
- `sessionId`  
- `participants`: [userId]  
- `createdAt`, `expiresAt`

`chatMessages` (subcollection under `chatThreads/{id}/messages`)  
- `senderId`  
- `body`: string  
- `type`: enum (`TEXT`, `IMAGE`, `SYSTEM`)  
- `attachment`: storage path (if type=IMAGE)  
- `sentAt`  
- `readBy`: [userId]  
- `deletedAt`

`transactions`  
- `id`  
- `requestId`  
- `payerId`, `payeeId`  
- `stripePaymentIntentId`  
- `status`: enum (`PENDING`, `AUTHORIZED`, `HELD`, `RELEASED`, `REFUNDED`, `DISPUTED`, `FAILED`)  
- `amount`: { `total`, `platformFee`, `helperShare`, `currency` }  
- `events`: array of { `status`, `timestamp`, `notes` }  
- `dispute`: { `initiatedBy`, `reason`, `evidence`, `resolvedAt`, `outcome` }  

`reviews`  
- `id`  
- `requestId`  
- `fromUserId`, `toUserId`  
- `rating`: int (1-5)  
- `comment`: string  
- `createdAt`, `flagged`: bool

`badges`  
- `id`  
- `name`, `description`, `criteria` metadata  

`badgeEvents`  
- `userId`  
- `badgeId`  
- `earnedAt`  
- `source`: enum (`COUNT`, `STREAK`, `CAMPAIGN`)  

`referrals`  
- `code`  
- `referrerId`  
- `status`: enum (`CREATED`, `CLAIMED`, `REWARDED`)  
- `reward`: object  

`auditLogs`  
- `actorId`  
- `action`  
- `targetType`, `targetId`  
- `metadata` (JSON)  
- `createdAt`

### 4.2 Storage Buckets
- `user-avatars/{uid}/{filename}`  
- `request-attachments/{requestId}/{filename}`  
- `kyc-documents/{uid}/{documentId}`  
- Retention policies enforce cleanup for expired content (chat attachments, disputes).

5. Backend Workflows
--------------------

5.1 Phone Authentication
1. Client requests Firebase OTP (verifies via Play Integrity/Recaptcha).  
2. Firebase sends SMS; user submits code; client obtains `AuthCredential`.  
3. On success, Cloud Function `onCreate` trigger seeds `users` document with default roles and helper opt-in flag.  
4. If OTP fails repeatedly, enforce cooldown and optionally mark suspicious numbers.

5.2 Help Request Submission
1. Client validates form locally (required fields, location).  
2. Document written to `helpRequests` with status `SUBMITTED`.  
3. Cloud Function trigger executes matching pipeline:  
   - Query eligible helpers (geohash neighbors, skill match, availability).  
   - Score helpers by distance, rating, previous acceptance rate.  
   - For urgent: send push + SMS to top five simultaneously.  
   - For planned: send push sequentially with configurable delay.  
   - Create `matchSessions` for each invite; update request status to `MATCHING`.  
4. Helper acceptance transitions request to `ACCEPTED`, logs acceptance time, triggers payment intent creation (pending capture).  
5. If no helper responds within SLA, escalate via additional SMS, fallback helpers, or advise seeker to call helpline (configurable).

5.3 Payment & Escrow Flow
1. Upon acceptance, Cloud Function creates Stripe Payment Intent with amount + platform fee, captures customer payment method.  
2. Funds authorized and held until completion.  
3. After both parties confirm completion:  
   - Cloud Function captures payment.  
   - Stripe transfers helper share to connected account (instant or scheduled).  
   - Transaction status changes to `RELEASED`.  
4. If confirmation mismatched or dispute raised:  
   - Status becomes `DISPUTED`.  
   - Operations team reviews evidence (chat logs, photos).  
   - Possible outcomes: refund to seeker, partial refund, release to helper.  
5. Failed payments trigger cancellation or fallback payment options; helper is notified to pause service until resolved.

5.4 Chat Lifecycle
1. Chat thread created when match session accepted.  
2. Messages stored under `chatThreads/{id}/messages`.  
3. Cloud Function scheduled job deletes messages older than 30 days, redacting attachments.  
4. Offensive content detection (optional integration) flags and notifies moderators.  
5. On dispute, chat transcripts exported to secure channel for review.

5.5 Badge & Impact Calculations
1. Cloud Function listens to `matchSessions` completions.  
2. Aggregates helper stats (total requests served, hours contributed).  
3. Evaluates badge criteria; on unlock, writes `badgeEvent`, triggers celebration notification.  
4. Updates helper dashboard metrics via cached documents to lower read cost.

6. Client Application Flows
--------------------------

### 6.1 Navigation Structure
- `Splash` → `Welcome` → `Onboarding` (phone entry, OTP, skills) → `Dashboard`  
- `Dashboard` tabs: Home, Requests, Messages, Impact, Profile  
- `Home`: welcome message, helper toggle, “Get Help Now” button.  
- `Requests`: active requests list, history, filters.  
- `Messages`: chat threads.  
- `Impact`: badges, stats, mind-map.  
- `Profile`: personal info, availability, preferences, documents.

### 6.2 State Management
- Use `Riverpod` or `Bloc` (to be determined) for reactive state.  
- Firestore streams supply real-time updates; caching via local database (`hive` or `sqflite`) for offline resilience.  
- Critical data (helper opt-in, availability) cached with expiration timestamps; revalidated on app resume.

### 6.3 Edge Case Handling
- **No GPS**: Prompt user to enable location; allow manual address entry; degrade radar to list view.  
- **No Network**: Display cached data, allow draft requests; background sync when online.  
- **OTP SMS delayed**: Provide call verification fallback (if Twilio/Telephony supported) or new attempt with countdown.  
- **Multiple Requests**: Limit seekers to one active urgent request; queue planned requests to avoid double booking.  
- **Helper Double Booking**: Check `matchSessions` before allowing acceptance; enforce cooldown after declines.  
- **Payment Failures**: Notify seeker to update payment; helper sees hold state; auto-cancel if unresolved within SLA.  
- **Disputed Service**: Freeze payouts, disable ratings until resolved, route to operations.  
- **Chat Abuse**: Allow users to block/report; restrict helper until review completed.  
- **SMS Delivery Failure**: Retry with alternate sender ID; fallback to voice call (optional).  
- **App Updates**: Remote Config prompts mandatory update when incompatible API changes shipped.  
- **Data Retention**: Auto-delete expired chats and attachments; anonymize user data on deletion requests.  
- **Helper Inactivity**: Send re-engagement notifications; auto-toggle to inactive after prolonged inactivity (configurable).  
- **Security Breach**: Incident response plan—revoke tokens, invalidate sessions, notify affected users.

7. API Contracts
----------------
Although Firestore uses client SDKs, define expected document shapes to maintain integrity.

### 7.1 Cloud Function HTTP Endpoints

| Endpoint | Method | Auth | Description |
|----------|--------|------|-------------|
| `/v1/match/retry` | POST | Service role | Force re-run matching for a request (admin/manual override). |
| `/v1/payment/hooks` | POST | Stripe signature | Stripe webhook (payment_succeeded, payment_failed, charge.dispute.created). |
| `/v1/sms/hooks` | POST | Twilio signature | Delivery receipts for urgent SMS. |
| `/v1/admin/override` | POST | Admin | Apply manual status change (cancel request, reassign helper). |
| `/v1/analytics/export` | POST | Scheduler | Nightly job to export snapshots to BigQuery. |
| `/v1/audit/log` | POST | Internal | Append log entry from admin tools. |

All endpoints validate App Check or server-to-server tokens, log correlation IDs, and enforce rate limits.

### 7.2 Firestore Security Considerations
- Seekers can read/write only their `helpRequests`; read limited fields of matched helper profiles.  
- Helpers can view requests they are invited to; cannot access other seekers’ personal data.  
- Chats restricted to participants + admins.  
- Transactions accessible by payer/payee; financial details hidden from other users.  
- Admin roles flagged in custom claims; Cloud Functions enforce privilege-based operations.  
- Security rules include data validation (e.g., rating between 1–5, status transitions allowed only from permitted states).  

8. Error Handling & Monitoring
-------------------------------
- Centralized client error reporter pushing logs to Crashlytics with user consent.  
- Cloud Functions use structured logging with correlation IDs and severity levels.  
- Alert thresholds:  
  - OTP failure rate >10% over 15 min → Slack alert.  
  - Matching queue backlog > 20 pending → incident page.  
  - Payment failures >5 per hour → escalate to finance.  
- Twilio/Stripe webhook retries handled idempotently (use event IDs).  
- Circuit breakers for third-party outages (disable urgent SMS after repeated failures; show banner).  
- Incident response runbook stored in knowledge base with RACI roles defined.

9. Testing Strategy
-------------------
- **Unit Tests**:  
  - Matching scoring function.  
  - Payment state machine.  
  - Badge calculation logic.  
- **Widget Tests**:  
  - Onboarding forms across locales.  
  - Radar animation state transitions.  
  - Chat message rendering with attachments.  
- **Integration Tests**:  
  - OTP flow using Firebase emulator.  
  - Request creation → helper acceptance → completion.  
  - Payment flow hitting Stripe test sandbox.  
- **Emulators**: Firebase Emulator Suite used for CI to avoid mutating live services.  
- **Manual QA**: Device matrix covering low-end and mid-range Android phones, offline scenarios, different locales.  
- **Performance Tests**: Gatling/JMeter hitting Cloud Functions, Firestore load simulations.  
- **Security Tests**: Automated Firebase rules tests; third-party pentest before launch.  
- **Release Testing**: Smoke tests executed post-deploy via automated scripts (e.g., Maestro).

10. Deployment & Release Management
-----------------------------------
- Feature branches → pull requests with required reviews.  
- Main branch protected; merges trigger staging build.  
- Stage environment tested by QA and stakeholders weekly.  
- Release branch created bi-weekly (or per milestone), triggers production build upon approval.  
- Play Store: staged rollout (10% → 50% → 100%) with monitoring between steps.  
- Rollback plan: revert to previous build via Play Console, disable new features via Remote Config.  
- Post-release review: metrics (crash-free sessions, match rate), collect feedback for next iteration.

11. Security & Compliance
-------------------------
- Enforce App Check and device integrity signals.  
- Secure storage of PII: minimal retention, encryption at rest (Firestore, Storage).  
- GDPR-like rights: data export, deletion, consent logs.  
- Stripe compliance (PCI DSS) handled via hosted elements; store only last 4 digits & brand.  
- KYC documents stored with limited access, short retention after verification.  
- Access reviews conducted quarterly; admin actions audited.  
- Disaster recovery: automated backups (Firestore export + Storage snapshots) with restore procedures tested quarterly.

12. Operations & Support
------------------------
- Customer support tools:  
  - Lookup by phone number; view recent requests, chats (if not expired), transactions.  
  - Trigger resend OTP, cancel requests, issue refunds (with permissions).  
- Escalation matrix: L1 support → L2 operations → L3 engineering.  
- Knowledge base articles for common issues (OTP failure, payment pending, helper no-show).  
- Incident postmortem template with blameless analysis.  
- Integration with ticketing system (Zendesk/Freshdesk) through webhooks.

13. Future Considerations
-------------------------
- Web portal for seekers lacking Android phones.  
- iOS client using same Flutter codebase (requires Apple-specific adjustments).  
- Machine-learning driven helper recommendations; dynamic pricing.  
- Integration with government/NGO databases for verified helpers.  
- Offline-first capabilities using local job queues.  
- Media-rich chat (voice notes, video) with increased storage policies.  
- Advanced trust model (background checks, ID verification services).  
- Gamification expansions: seasonal challenges, team challenges, corporate sponsorship.

Appendices
----------

A. State Transition Table (Help Request)
```
SUBMITTED -> MATCHING (Cloud Function)
MATCHING -> ACCEPTED (Helper acceptance)
MATCHING -> CANCELLED (Seeker cancel; timeout without helper)
ACCEPTED -> IN_PROGRESS (Helper indicates en route/arrived)
IN_PROGRESS -> COMPLETED (Both confirm success)
ACCEPTED/IN_PROGRESS -> CANCELLED (Either party cancel; reason logged)
COMPLETED -> DISPUTED (Seeker files dispute within allowable window)
DISPUTED -> RESOLVED (Admin outcome; status remains COMPLETED or REFUNDED)
```

B. Notification Matrix
```
Event                     Seeker                    Helper                   Admin
-------------------------------------------------------------------------------------------------
Request Submitted         Push confirmation         n/a                      n/a
Urgent invite sent        n/a                       Push + SMS               n/a
Helper accepted           Push + in-app banner      Push confirmation        Optional Slack
Helper declined           Push with retry option    n/a                      n/a
Payment pending           In-app prompt             In-app status            Finance alert (if > SLA)
Completion confirmed      Push + rating prompt      Push + payout pending    n/a
Dispute opened            Email + in-app message    Email + in-app message   Support ticket
Badge earned              Push + celebration        Push + celebration       n/a
System outage             Status banner             Status banner            PagerDuty alert
```

C. Glossary
- **Helper**: User offering assistance; default state upon signup unless opted out.  
- **Seeker**: User requesting help (may also act as helper later).  
- **Match Session**: Pairing between a help request and a helper candidate.  
- **ESCROW**: Payment funds held until completion confirmation.  
- **SLA**: Service Level Agreement for matching and completion times.  
- **KYC**: Know Your Customer verification process.

End of Document

