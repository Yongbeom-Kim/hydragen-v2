# 2026-02-13 - Auth & Permission Model

This document defines the Authentication and Authorization model for Hydragen V2.
Goal: strict syllabus adherence with explicit, enforceable access rules.

## Core Policy

> [!IMPORTANT]
> Syllabus adherence is a hard product requirement.
> Instructors define session scope first; personalization happens only inside that scope.

## Authentication

- Authentication is handled by [`authentik`](https://goauthentik.io/) via [`OIDC`](https://openid.net/connect/).
- Global role claims come from the IdP token.

## Authorization Model

Authorization is deny-by-default and uses intersection checks:

1. Global role claim (Student, Instructor, Admin) grants capability class.
2. Session membership in the application database grants access to a specific session.
3. Request is allowed only if all required checks pass.

Rules:
- Session membership stores membership only (no per-session role field).
- Every session-protected endpoint must re-check membership in the database per request.
- If no allowed skills are selected in a session, return no spectra.

## Roles

Global roles:

1. Student
2. Instructor
3. Admin

Defaults and powers:
- Everyone starts as Student, except internal Hydragen team bootstrap admins.
- Admins can promote/demote roles.
- Instructors can manage sessions they belong to.
- Students can participate in sessions they belong to.

Admin safety controls:
- The last remaining Admin cannot be demoted or removed.
- Self-demotion is allowed only if at least one other Admin remains.
- Long-term: require two-admin approval for Instructor -> Admin promotion.

## Session Model

A session is a course container (for example, semester-long).
When an Instructor or Admin creates a session, they define an allowed skill set.

Skill model:
- Skills are the only curriculum scope control.
- Session-visible spectra are derived from `allowedSkills`.
- We maintain an index: `Skill -> Molecule/Spectrum`.
- Empty `allowedSkills` means restricted mode (no spectra).

Membership and management model:
- Session creator is automatically added as a session member.
- Any other user (including Instructors) must be invited/added to join.
- Session mutation (membership and allowed skills) requires:
  - global role in `{Instructor, Admin}`
  - and active membership in that session.

## Invite and Membership Lifecycle

Hybrid invite model:
- Existing registered users can be directly added as members.
- New or unregistered users receive a tokenized invite.

Tokenized invite policy:
- Single-use token.
- Expires in 7 days.
- Revocable by Instructor/Admin session managers before acceptance.

Removal policy:
- On member removal, access is revoked immediately for new API requests.
- Existing tokens are not force-revoked; protected endpoints rely on per-request membership checks.

## Endpoint Policy

Session endpoints:
- `POST /sessions` -> Global `Instructor` or `Admin`.
- `POST /sessions/{id}/members` -> Global `Instructor` or `Admin`, and session member.
- `DELETE /sessions/{id}/members/{userId}` -> Global `Instructor` or `Admin`, and session member.
- `PUT /sessions/{id}/allowed-skills` -> Global `Instructor` or `Admin`, and session member.
- `GET /sessions/{id}/allowed-skills` -> Session member or global `Admin`.
- `GET /sessions/{id}/mass-spectra` -> Session member or global `Admin`; response filtered by `allowedSkills`.

Invite endpoints:
- `POST /sessions/{id}/invites` -> Global `Instructor` or `Admin`, and session member.
- `POST /sessions/{id}/invites/{inviteId}/revoke` -> Global `Instructor` or `Admin`, and session member.
- `POST /invites/{token}/accept` -> Authenticated invited user; token must be valid, unexpired, and unused.

Admin role endpoints:
- `POST /admin/roles/instructors/{userId}` -> Admin only.
- `DELETE /admin/roles/instructors/{userId}` -> Admin only.
- `POST /admin/roles/admins/{userId}` -> Admin only.
- `DELETE /admin/roles/admins/{userId}` -> Admin only, but cannot remove last Admin.

## Audit Logging Requirements

Must be audit logged:
- Role changes (all admin/instructor promote/demote events).
- Session create/delete.
- Session membership add/remove.
- Session `allowedSkills` changes.

Each log record should include:
- Acting user ID.
- Target resource/user ID.
- Action type.
- Timestamp.
- Result (success/denied).

## Assumptions and Future Work

Current assumption:
- Single-tenant deployment for now.
- Session boundaries are the primary access boundary at current scale.

Future work:
- Add explicit organization/tenant boundaries for multi-school deployments.
- Add stronger privileged-role workflows (for example, two-admin approval).
