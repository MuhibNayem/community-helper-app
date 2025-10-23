Community Helper Software Requirements Specification
====================================================

1. Product Overview
-------------------
- **Purpose**: Formalize requirements for the Community Helper Android application that connects help seekers with nearby helpers for urgent and scheduled assistance.
- **Scope**: Android-first MVP delivered with FlutterFlow-generated UI plus custom Flutter code, Firebase backend services, Twilio SMS, and Stripe-based escrow payments. Future web/iOS channels sit outside this SRS except where constraints prepare for them.
- **Objectives**: Enable secure onboarding, reliable matching, protected payments, and motivational progress tracking while maintaining trust, speed, and affordability for Bangladeshi communities.
- **Success Metrics**: Cold start < 2 s; ≥95% OTP success; ≥90% urgent requests matched within 5 min; ≥4.5 helper rating average; <0.5% disputed payouts; daily active helpers / requests growth ≥10% month-over-month post-launch.

2. Stakeholders and Personas
----------------------------
- **Help Seeker**: Needs immediate or planned help; prioritizes fast response, clear status, transparent pricing, and personal safety.
- **Helper**: Offers skills/availability; values dependable requests, timely payouts, reputation tracking, and schedule control.
- **Operations/Admin**: Monitors platform health, resolves disputes, executes manual overrides; requires dashboards, audit trails, and compliance tools.
- **Customer Support**: Assists users via chat or phone; needs visibility into request history, chat transcripts, payment status, and SMS logs.
- **Business Stakeholders**: Track KPIs, growth, and engagement to inform marketing, partnerships, and roadmap prioritization.

3. Assumptions, Dependencies, Constraints
-----------------------------------------
- Android devices running API level 24+ with Google Play services available.
- Users possess active phone numbers capable of receiving OTPs and SMS alerts.
- FlutterFlow subscription active during build; custom Flutter code maintained in GitHub.
- Firebase (Auth, Firestore, Storage, Cloud Functions, FCM, Analytics) provides serverless backend. Separate dev/stage/prod projects maintained.
- Third-party services: Twilio Programmable SMS (with local sender IDs), Stripe Connect (or compliant local payment alternative), Google Maps/Directions, Lottie animations, Figma for design assets.
- Regulatory compliance: Bangladeshi KYC/AML standards for payouts, data privacy obligations similar to GDPR/PDPA.
- Network connectivity intermittent; app must degrade gracefully and cache critical state.
- Project delivered by a part-time developer (16 h/week) with cross-functional stakeholder inputs available on demand.

4. Functional Requirements
--------------------------
4.1 Onboarding and Identity  
- Collect name, phone number, consent acknowledgements, and helper opt-in toggle (default ON).  
- Validate phone via Firebase OTP; implement retry limits, cooldowns, and fallback verification.  
- Capture user location permissions with rationale prompts; support manual address entry.  
- Persist helper settings and resend onboarding reminders for inactive helpers.

4.2 Helper Availability and Profile Management  
- Define skill taxonomy (e.g., medical, mechanical, tutoring) manageable by admins.  
- Allow helpers to select skills, set recurring availability windows, and toggle temporary unavailability.  
- Display earned badges, impact metrics, ratings, and verification status.  
- Enable helpers to update profile photo, bio, and communication preferences.  
- Record KYC progress, verification documents, and manual review notes.

4.3 Help Request Lifecycle  
- Provide “Urgent” and “Plan Ahead” flows with contextual copy explaining expectations and pricing.  
- Auto-detect location via GPS; permit manual adjustments and saved addresses.  
- Collect request details: category, description, optional attachment, scheduled time (if planned).  
- Persist requests in Firestore with lifecycle states (Draft, Submitted, Matching, Accepted, In Progress, Completed, Cancelled, Disputed).  
- Allow seekers to cancel with reason codes; enforce penalties where applicable.  
- Present request history with filters, status badges, receipts, and ability to repeat a request.

4.4 Matching and Location Services  
- Compute helper eligibility based on proximity, skill match, availability, and rating threshold.  
- Implement radar-inspired animation showing nearby helpers, updating in real time as responses arrive.  
- For urgent requests, broadcast to top N helpers via push plus SMS; throttle to respect per-minute limits.  
- For planned requests, notify helpers sequentially with configurable wait periods.  
- Display ETA and map route once a helper accepts; update both parties with location deltas every 10 seconds (tunable).  
- Handle declines, timeouts, and fallback matching logic.

4.5 Communication and Notifications  
- Provide in-app chat between matched parties with text, emoji, and image attachments; redact after 30 days.  
- Show typing indicators, read receipts, message timestamps, and offline caching.  
- Support push notifications for request updates, chat messages, reminders, and streak encouragement.  
- Deliver SMS notifications for urgent escalation and critical system alerts; log delivery status.  
- Respect user preferences for notification quiet hours and opt-out scenarios.  
- Capture system messages (e.g., badge unlocks, policy updates) with deep links.

4.6 Payments and Escrow  
- Define pricing model (base fee + platform fee + helper payout) configurable via Remote Config.  
- Initiate Stripe payment intent upon helper acceptance; place funds in escrow until completion.  
- Provide secure payment UI with saved payment methods, receipts, and refunds.  
- Require dual confirmation (“Help completed successfully”) to release funds; escalate to dispute workflow if mismatched.  
- Support payouts to helpers’ connected accounts with compliance checks (KYC, bank verification).  
- Maintain transaction ledger with statuses (Pending, Held, Released, Refunded, Disputed).  
- Enable admins to trigger manual adjustments, refunds, or freeze accounts.

4.7 Impact Dashboard and Gamification  
- Calculate helper metrics: total people helped, cumulative hours, impact score, streaks.  
- Render mind-map visualization of connected help network with zoom/pan controls.  
- Award badges based on milestones (e.g., Community Hero, Lifesaver) and limited-time campaigns.  
- Trigger celebration animations (confetti, haptics) upon achievements.  
- Provide “Invite Friends” flow with referral tracking and incentives.  
- Surface suggestions for next available slots to encourage repeat contributions.

4.8 Admin and Operations Tooling  
- Provide secure web portal (role-based) for monitoring active requests, user status, and system health.  
- Allow admins to reassign helpers, resolve disputes, manage content (welcome messages, skill taxonomy), and view analytics.  
- Expose audit logs for critical actions (payout overrides, verification decisions).  
- Integrate alerting for failed SMS, payment issues, or unusual activity patterns.  
- Support exporting operational data (CSV/BigQuery) for reporting.

5. Non-Functional Requirements
------------------------------
- **Performance**: App cold start < 2 s on mid-tier Android; screen transitions < 500 ms; chat message latency < 1 s; matching responses processed < 3 s.  
- **Scalability**: Support 10,000 concurrent users; Firestore indexes tuned to prevent throttling; Cloud Functions auto-scale with queue-based throttling.  
- **Reliability**: 99.5% monthly uptime target; retry strategies with exponential backoff for network calls; graceful degradation for external service outages.  
- **Security**: TLS for all traffic; Firebase security rules using least privilege; encryption of sensitive fields (KYC data, financial tokens); periodic security audits.  
- **Privacy**: Compliance with consent and data minimization; auto-deletion of chats after 30 days; user controls for data download/delete; privacy policy accessible in-app.  
- **Usability & Accessibility**: WCAG AA color contrast; large tap targets; simple copy in Bangla and English (localization-ready); support for font scaling and screen readers.  
- **Maintainability**: Modular codebase, documented APIs, automated tests, lint/format enforcement, CI/CD pipeline with gated releases.

6. System Architecture
----------------------
- **Client**: FlutterFlow-generated components extended with custom Flutter modules for radar animation, chat, payment UI, and platform-specific plugins.  
- **Backend**: Firebase Auth (phone OTP), Firestore (primary data store), Storage (media), Cloud Functions (business logic, cron jobs), FCM (push), Remote Config (feature flags).  
- **Integrations**: Twilio SMS (urgent alerts, OTP fallback), Stripe Connect (escrow, payouts), Google Maps & Places (location services), Lottie (animations), BigQuery/Data Studio (analytics).  
- **Environments**: Dedicated dev/stage/prod Firebase projects; configuration handled via environment files and remote config; feature flags for staged rollout.  
- **Monitoring**: Firebase Crashlytics, Performance Monitoring, Stripe/Twilio webhooks, uptime checks with alert routing to Slack/email.

7. Data Model Overview
-----------------------
- **User**: ID, phone, name, role flags, helper opt-in, location consent, notification preferences, KYC status, created/updated timestamps.  
- **HelperProfile**: User ref, skills, availability schedule, rating average/count, badge progress, last active location, verification metadata.  
- **HelpRequest**: Requester ref, type (Urgent/Planned), description, location geohash/address, scheduled time, attachments, lifecycle status, timestamps, cancellation reason.  
- **MatchSession**: Request ref, helper ref, acceptance timestamps, ETA, status transitions, rejection reasons, analytics metadata.  
- **ChatMessage**: Session ref, sender, message body/attachment path, read status, created timestamp, deletion timestamp.  
- **PaymentTransaction**: Request ref, amount breakdown, Stripe intent IDs, escrow status, payout status, refund/dispute references.  
- **BadgeEvent**: User ref, badge ID, earned timestamp, source event.  
- **AuditLog**: Actor, action, target entity, timestamp, metadata (IP, device ID).

8. Key User Journeys
--------------------
1. **Registration**: Install → permissions prompts → enter phone → OTP verification → helper toggle modal → select skills → dashboard intro.  
2. **Urgent Request**: Tap “Get Help Now” → choose urgent → confirm location → submit details → radar & notifications → helper accepts → chat & map guidance → completion confirmation → review & rating.  
3. **Planned Request**: Choose planned option → pick date/time → helper invitation sequence → acceptance → reminders → service completion → payout release.  
4. **Helper Response**: Receive push/SMS → view request summary → accept/decline → navigate to location → perform help → confirm completion → receive rating and payout.  
5. **Dispute Resolution**: Either party reports issue → auto-hold funds → admin reviews evidence → resolves with refund/payout → notify parties.  
6. **Impact Tracking**: Helper opens dashboard → views badges/network → checks next availability → shares invite link → receives streak reminder.

9. Risks and Mitigation Strategies
----------------------------------
- **SMS Delivery Failures**: Implement provider status callbacks, retry via alternative sender ID, supplement with push notifications and in-app alerts.  
- **Payment Compliance Hurdles**: Validate Stripe availability; prepare local gateway integration fallback; maintain KYC documentation workflow.  
- **Location Inaccuracy**: Provide manual location editing, accuracy indicators, caching; allow help requests to include landmarks/directions.  
- **Matching Scalability**: Pre-compute geohashes, use Firestore composite indexes, queue matching operations with Cloud Functions to avoid rate limits.  
- **Trust & Safety Incidents**: Introduce verification badges, rating thresholds, block/report features, incident response SOP, and background checks where feasible.  
- **Part-Time Velocity Risk**: Maintain groomed backlog, prioritize high-risk items early, enforce weekly demos/check-ins to catch issues quickly.

10. Quality Assurance Plan
--------------------------
- **Automated Testing**: Unit tests for matching, payment, and badge logic; widget tests for key UI components; integration tests covering OTP, chat, and payment flows.  
- **Manual Testing**: Scenario-based QA for urgent/planned journeys, error states, offline mode, localization, accessibility; regression suite per milestone.  
- **Performance Testing**: Load tests simulating concurrent requests/chat; monitor Firestore throughput; evaluate cold start on representative devices.  
- **Security Testing**: Regular review of Firebase rules, dependency vulnerability scans, penetration testing of client endpoints.  
- **Release Process**: CI pipeline running tests and linters; Firebase App Distribution for QA; staged rollout on Play Store with rollback plan; post-release analytics review within 48 hours.

11. Open Questions
------------------
- Finalize pricing tiers, helper compensation percentages, and platform fee structure.  
- Determine legal requirements for operating payment escrow within Bangladesh.  
- Confirm acceptable identity verification process (NID upload, manual review turnaround).  
- Decide on default language strategy (Bangla-first vs bilingual).  
- Clarify post-launch support SLAs and bug-fix window.  
- Validate long-term roadmap (web portal for seekers, telemedicine tie-ins, B2B partnerships).

