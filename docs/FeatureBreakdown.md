Community Helper Feature Breakdown
==================================

1. Onboarding & Identity
------------------------
- **F1.1 Requirements Analysis**: Draft onboarding copy, consent language, helper opt-in defaults, and localization checklist.
- **F1.2 UX & Visuals**: Design screens for phone entry, OTP, profile setup, and permissions prompts.
- **F1.3 Phone Authentication**: Integrate Firebase phone auth, handle retries, throttling, and reCAPTCHA.
- **F1.4 Profile Completion**: Implement name, avatar, skill selection UI; store helper preferences.
- **F1.5 Error/Offline Handling**: Surface friendly errors, cached retries, and blocked-number flows.
- **F1.6 Analytics & Events**: Track onboarding step success/failure, helper opt-in rates, drop-off points.

2. Helper Availability & Reputation
-----------------------------------
- **F2.1 Skill Taxonomy**: Define categories, icons, admin CRUD; sync with request types.
- **F2.2 Availability Scheduler**: Build recurring slots, quick pause, and immediate toggle.
- **F2.3 Helper Dashboard**: Display upcoming commitments, earnings summary, impact snippet.
- **F2.4 Ratings & Reviews**: Collect post-service feedback, moderation queue, reporting tools.
- **F2.5 Verification/KYC**: Capture documents, verification status, reminders, manual review notes.
- **F2.6 Helper Analytics**: Track streaks, badge progress, utilization metrics.

3. Help Request Lifecycle
-------------------------
- **F3.1 Request UX**: Separate urgent/planned forms, contextual helper expectations, pricing hints.
- **F3.2 Location Services**: GPS auto-detect, manual map pin, saved addresses, accuracy indicator.
- **F3.3 Submission Pipeline**: Firestore schema, Cloud Functions triggers, SLA timers.
- **F3.4 Cancellation & Penalties**: Define policy, implement cancellation reasons, automated notices.
- **F3.5 Request History**: List past/current requests, receipts, repeat request shortcut.
- **F3.6 Accessibility & Localization**: Ensure form readability, adaptive UI, bilingual support.

4. Matching & Radar Visualization
---------------------------------
- **F4.1 Eligibility Engine**: Filter helpers by distance, skill, availability, rating thresholds.
- **F4.2 Matching Logic**: Prioritize urgent broadcasts vs planned sequential invites.
- **F4.3 Radar Component**: Build animation, real-time helper dots, status transitions.
- **F4.4 Notifications Orchestration**: Coordinate push, SMS, and in-app alerts with rate limits.
- **F4.5 Acceptance Handling**: Manage accept/decline, timeouts, fallback to next helper.
- **F4.6 ETA & Navigation**: Fetch directions, update ETA, display route progress.

5. Communication & Notifications
--------------------------------
- **F5.1 Chat UI**: Message list, input bar, attachments, typing indicators, read receipts.
- **F5.2 Chat Backend**: Firestore schema, retention job (30-day delete), encryption at rest.
- **F5.3 Urgent SMS Escalation**: Template management, localization, delivery reporting.
- **F5.4 Push Notifications**: Configure FCM topics, deep links, background handling.
- **F5.5 System Messaging**: Broadcast updates, policy changes, badge notifications.
- **F5.6 Quiet Hours & Preferences**: Allow per-user notification schedules, overrides for emergencies.

6. Payments & Escrow
--------------------
- **F6.1 Pricing Model Definition**: Document fee structure, discounts, and subsidy options.
- **F6.2 Stripe Integration**: Payment intent creation, customer management, webhook handling.
- **F6.3 Escrow Flow**: Hold funds, dual confirmation, automated release, timeout rules.
- **F6.4 Payouts**: Helper onboarding, bank verification, payout scheduling, receipts.
- **F6.5 Dispute Workflow**: Evidence submission, escalation path, resolution outcomes.
- **F6.6 Financial Reporting**: Ledger views, CSV exports, reconciliation scripts, audit logs.

7. Impact Dashboard & Gamification
----------------------------------
- **F7.1 Metrics Engine**: Calculate impact score, total assists, hours contributed.
- **F7.2 Visualization**: Mind-map network, zoom/pan, highlight recent activity.
- **F7.3 Badges & Levels**: Define tiers, unlock rules, limited-time campaigns.
- **F7.4 Celebrations**: Integrate Lottie animations, haptics, shareable moments.
- **F7.5 Referral Flow**: Invite links, referral tracking, reward distribution.
- **F7.6 Helper Nudges**: Remind helpers of next availability, highlight demand hotspots.

8. Admin & Operations
---------------------
- **F8.1 Admin Portal Auth**: Role-based access control, SSO considerations, audit logging.
- **F8.2 Request Command Center**: Live board for urgent/planned requests, manual reassignment.
- **F8.3 Trust & Safety Tools**: Block/report management, incident escalation, user suspension.
- **F8.4 Content Management**: Edit welcome messages, announcements, skill taxonomy updates.
- **F8.5 Analytics Dashboard**: KPIs (match rate, churn, response time), BigQuery/Data Studio integrations.
- **F8.6 System Monitoring**: Health checks, alert routing for SMS/payment errors, cron job status.

9. Quality, Security, Compliance
--------------------------------
- **F9.1 Test Automation**: Unit/integration test suites, coverage tracking, CI enforcement.
- **F9.2 Performance Testing**: Load scripts, scaling thresholds, caching strategy evaluation.
- **F9.3 Security Hardening**: Review Firebase rules, dependency scanning, penetration tests.
- **F9.4 Accessibility Audit**: WCAG compliance review, assistive tech testing, localization QA.
- **F9.5 Release Management**: Versioning policy, staged rollout, rollback procedure.
- **F9.6 Support Playbook**: Post-launch support SLAs, incident response, knowledge base assets.

Task Grouping by Phase
----------------------
- **Discovery & Planning**: F1.1, F2.1, F3.1, F6.1, F7.1, F9.6.  
- **Design & Prototyping**: F1.2, F2.2, F3.1, F4.3, F5.1, F7.2.  
- **Core Implementation**: F1.3–F1.5, F2.3–F2.6, F3.2–F3.5, F4.1–F4.6, F5.2–F5.6, F6.2–F6.6.  
- **Enhancements & Ops**: F7.3–F7.6, F8.1–F8.6, F9.1–F9.5.  
- **Launch Readiness**: Consolidated testing, compliance checks, store submission, go-live runbooks.

