Community Helper Project Pricing
===============================

Cost Model Assumptions
----------------------
- Developer rate: BDT 2,500/hr (part-time 16 h/week baseline).
- Designer rate: BDT 2,000/hr; Product/QA/Compliance support: BDT 1,800/hr.
- Contingency buffer applied per phase to cover integration surprises and schedule slack.
- Exchange rate reference: 1 USD ≈ BDT 110 (for context only).

Phase-by-Phase Estimate
-----------------------
| Phase | Duration (weeks) | Core Deliverables | Est. Hours (Dev/Design/PM) | Cost (BDT) | Rationale |
|-------|------------------|-------------------|----------------------------|------------|-----------|
| Discovery & UX Foundations | 4 | Requirements lock, pricing strategy workshop, low-fi UX, localization plan | 160 / 80 / 40 | 651,000 | Early validation reduces downstream rework; includes workshops with stakeholders. |
| Core Build Sprint 1 | 8 | Firebase env setup, onboarding, helper profiles, availability scheduler, request submission | 256 / 80 / 80 | 1,058,000 | Covers initial app architecture, CI skeleton, and QA passes on foundational flows. |
| Core Build Sprint 2 | 8 | Matching engine, radar UI, chat MVP, push/SMS orchestration | 256 / 60 / 100 | 1,052,000 | Real-time features require extra QA/time for Firestore rules and messaging edge cases. |
| Payments & Compliance | 6 | Stripe escrow, payouts, disputes, KYC workflow, legal review | 192 / 40 / 120 | 892,000 | Higher risk work; compliance consultations and sandbox certifications included. |
| Gamification & Ops Tooling | 5 | Impact dashboard, badges, admin console, monitoring alerts | 160 / 60 / 100 | 784,000 | Mix of UX polish and backend tooling to support operations and engagement metrics. |
| QA & Launch Readiness | 5 | Automated test coverage, performance/security testing, release runbook | 120 / 0 / 200 | 726,000 | Final hardening with extended QA cycles, staging distribution, and incident drills. |

Totals
------
- Core delivery subtotal: **BDT 5,163,000**
- Recommended contingency reserve (8%): **BDT 413,000**
- **Grand total budget: BDT 5,576,000**

Notes
-----
- Rates reflect Dhaka market for senior contractors; adjust for in-house salaries if different.
- External service fees (Twilio sender IDs, Stripe onboarding, device lab access) are not included and should be budgeted separately.
- Accelerating delivery would require additional developer/designer capacity, increasing the hourly burn accordingly.
