# 2026-02-13 - Telemetry Exploration

This document defines telemetry scope for Hydragen V2 across three goals:
- Educational outcomes (student learning quality)
- Product reliability (bugs, incidents, alarms)
- Compliance and research ethics (EU GDPR + IRB-ready governance)

## Objectives

- Measure whether students are learning better over time, not just answering more questions.
- Detect regressions, outages, and data quality failures early.
- Collect the minimum data needed to answer product and research questions.
- Keep implementation compatible with GDPR requirements and likely IRB review concerns.

## Telemetry Domains

### 1) Educational Outcomes Telemetry

Core events:
- `question_assigned`
- `question_answered`
- `answer_revealed`
- `hint_requested`
- `attempt_submitted`
- `session_started` / `session_ended`

Core outcome metrics:
- Learning gain (pre/post within-topic delta)
- Time-to-mastery per concept
- Retention proxy (repeat performance after delay)
- Confidence calibration (self-confidence vs correctness)
- Struggle indicators (hint rate, repeated failures, abandon rate)

Minimum event payload (pseudonymized):
- `event_id`, `event_name`, `event_time`
- `learner_id_pseudo` (salted hash / stable pseudonymous key)
- `cohort_id_pseudo` (optional)
- `question_id`, `concept_id`
- `attempt_index`, `is_correct`, `latency_ms`
- `client_version`, `model_version` (if model-generated question)

### 2) Product Reliability Telemetry (Bugs + Alarms)

Signal types:
- Structured logs (server + client error boundaries)
- Metrics (request rate, p95/p99 latency, error rate, DB/query failures)
- Traces (critical API spans, DB spans)
- Alerts (SLO burn, crash loops, ETL/data drift)

MVP alerts:
- API 5xx rate above threshold (rolling 5m + 30m windows)
- Mass spectrum query latency regression
- ETL freshness lag beyond expected interval
- Frontend route crash spike (`/data`, `/data/compounds/*`, `/data/mass_spectra/*`)

## GDPR-by-Design Requirements

Primary legal references for implementation:
- GDPR Regulation (EU) 2016/679 (EUR-Lex consolidated text)
- EDPB Guidelines 4/2019 on Article 25 (Data Protection by Design and by Default)

Engineering controls:
- Data minimization by default: only collect fields required for named metrics.
- Purpose limitation: split educational analytics vs operational monitoring purposes.
- Lawful basis register per event family (Article 6); document any special-category handling (Article 9) as out-of-scope unless explicitly approved.
- Transparent notices and consent/legitimate-interest disclosure where applicable (Articles 12-14).
- Data subject rights workflows for access/rectification/erasure/restriction/portability/objection (Articles 15-21).
- Records of processing activities maintained (Article 30).
- Security controls: encryption in transit/at rest, least-privilege access, key rotation, audit logging (Article 32).
- Breach response runbook with notification workflow and timing constraints (Articles 33-34).
- DPIA gate before large-scale educational profiling or new high-risk analytics (Article 35).
- Cross-border transfer guardrails for non-EU vendors (Chapter V, Article 44+).

### Data Classification (Telemetry)

- `P0` Direct identifiers: email, name, IP-full, free-text with potential PII
- `P1` Pseudonymous learning analytics identifiers
- `P2` Operational metadata (route, version, error code, latency)
- `P3` Aggregated anonymous metrics

Rule:
- Default pipeline accepts `P1+P2`; rejects `P0` unless explicitly approved in DPIA + IRB protocol.

### Retention Policy (MVP)

- Raw event stream (`P1/P2`): 90 days
- Curated analytics tables: 12 months
- Fully aggregated, anonymous trend tables: 24+ months
- Security/audit logs: per policy and legal requirement, separately controlled

## IRB-Readiness (Expected Objections + Mitigations)

Reviewer concerns we should expect:
- Privacy/confidentiality risk from educational behavior tracking
- Re-identification risk across joined datasets
- Potential coercion/undue influence in classroom settings
- Disparate impact or unfair outcomes across groups
- Unclear consent language and participant understanding
- Over-collection not justified by stated research aims

Mitigations to pre-build:
- Protocol packet: purpose, hypothesis, variables, minimization rationale.
- Consent artifacts: plain-language participant notice; opt-out path when required.
- De-identification and access-control plan.
- Risk-benefit statement aligned to Belmont principles (respect, beneficence, justice).
- Data monitoring plan and adverse-event/escalation criteria.
- Pre-registered analysis plan for primary outcomes (reduce p-hacking risk).
- Equity audit slice definitions and periodic fairness review.

## MVP Deliverables

1. Telemetry event taxonomy + schema definitions (versioned).
2. Backend/client instrumentation for core educational and reliability events.
3. Alerting dashboard with baseline SLOs.
4. GDPR controls checklist implemented in pipeline and docs.
5. IRB submission support pack template (protocol summary + data handling appendix).

## Suggested Technologies

### Option A (Recommended): EU-First, Self-Hosted OSS Core

- Instrumentation:
  - OpenTelemetry SDK (Go backend + frontend web instrumentation)
- Collection and routing:
  - OpenTelemetry Collector or Grafana Alloy
- Reliability observability:
  - Grafana + Loki (logs) + Tempo (traces) + Prometheus/Mimir (metrics)
  - Alertmanager (paging and alert routing)
- Educational analytics store:
  - ClickHouse with strict TTL/retention rules for pseudonymous event data
- Product analytics UI (optional for speed):
  - PostHog (self-host) for funnels/retention dashboards if needed
- Error tracking:
  - Sentry (self-host preferred, or EU-region cloud with DPA and transfer review)

Why this is recommended:
- Strong GDPR posture via data residency and tighter control of processing.
- OpenTelemetry avoids vendor lock-in and keeps future backend options open.
- Fits current stack (Go + React + Docker Compose) with incremental rollout.

### Option B: Managed-Hybrid (Faster Setup, More Vendor Governance)

- Keep OpenTelemetry instrumentation and collector.
- Use managed observability/analytics vendors with EU-region hosting.
- Require DPA + transfer impact assessment + documented subprocessors for each vendor.

Tradeoff:
- Faster initial delivery, but higher compliance/vendor management overhead.

## Non-Goals (This Phase)

- Full research-grade causal inference stack
- Automated adaptive interventions based on protected attributes
- Long-term warehousing of raw pseudonymous event-level data

## References

Legal and ethics:
- GDPR (official EU law metadata; CELEX 32016R0679): https://op.europa.eu/en/web/eu-law-in-force/bibliographic-details/-/elif-publication/3e485e15-11bd-11e6-ba9a-01aa75ed71a1
- GDPR ELI identifier (canonical): http://data.europa.eu/eli/reg/2016/679/oj
- EDPB Guidelines 4/2019 (Article 25, Data Protection by Design and by Default): https://www.edpb.europa.eu/our-work-tools/our-documents/guidelines/guidelines-42019-article-25-data-protection-design-and_en
- HHS OHRP 45 CFR 46 (Common Rule): https://www.hhs.gov/ohrp/regulations-and-policy/regulations/45-cfr-46/index.html
- Belmont Report (Respect for persons, beneficence, justice): https://www.hhs.gov/ohrp/regulations-and-policy/belmont-report/index.html

Telemetry and observability technologies:
- OpenTelemetry Collector docs: https://opentelemetry.io/docs/collector/
- OpenTelemetry Go docs: https://opentelemetry.io/docs/languages/go/
- Grafana LGTM (Docker OTel LGTM): https://grafana.com/docs/opentelemetry/docker-lgtm/
- Grafana Alloy docs: https://grafana.com/docs/alloy/latest/
- Prometheus Alertmanager docs: https://prometheus.io/docs/alerting/latest/alertmanager/
- ClickHouse TTL docs: https://clickhouse.com/docs/guides/developer/ttl
- ClickHouse access control docs: https://clickhouse.com/docs/operations/access-rights
- Sentry API + data region note: https://docs.sentry.io/hosted/api/
- Sentry SDK data collection controls: https://docs.sentry.io/platforms/javascript/guides/react/data-management/data-collected
- PostHog trust center (compliance docs): https://trust.posthog.com/
