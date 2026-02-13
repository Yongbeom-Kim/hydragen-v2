# 2026-02-13 - Mass Spec CRUD Frontend (MVP For Graph Display)

This document defines the MVP plan for a Mass Spectrum CRUD frontend, with graph display as the main product outcome. The goal is to move from a read-only data viewer to a curation-friendly workflow where we can create, edit, review, and delete spectrum records safely.

## Current State (Already Implemented)

- Frontend data viewer pages exist:
  - `/data`: compound list with pagination
  - `/data/compounds/$inchiHash`: compound detail
  - `/data/mass_spectra/$inchiHash`: spectrum metadata + chart(s)
- Backend read endpoints exist:
  - `GET /compounds`
  - `GET /compounds/{inchiKey}`
  - `GET /compounds/{inchiKey}/image`
  - `GET /mass-spectra/{inchiKey}`
- Graph rendering exists with ECharts (`mZ` on X-axis, `peaks` on Y-axis).

## MVP Scope

- CRUD target: `mass_spectra` records.
- Compounds remain read-mostly for MVP (from ETL), with optional linking to existing compounds by `inchiKey`.
- Primary UI outcome: editable spectrum records with immediate graph preview.

### Create
- Create a new spectrum record from a selected compound (`inchiKey`).
- Input fields:
  - Required: `dbNumber`, `source`, `molecularWeight`, `mZ[]`, `peaks[]`
  - Optional: `exactMass`, `precursorMz`, `precursorType`, `ionMode`, `collisionEnergy`, `spectrumType`, `instrument`, `instrumentType`, `splash`, `comments`
- Validate `mZ.length === peaks.length` and reject empty arrays.

### Read
- Continue list/detail/graph flow from existing pages.
- Add filter/sort controls for quick spectrum selection in the compound context.

### Update
- Edit metadata and peak arrays for an existing spectrum record.
- Show before-save graph preview and post-save refreshed graph.

### Delete
- Delete a spectrum record with explicit confirmation modal.
- UI should show affected record identifiers (`id`, `dbNumber`, `source`) before confirming.

## Backend/API Needed for MVP

Existing read endpoints are not enough for CRUD. Add write endpoints:

- `POST /mass-spectra`
- `PATCH /mass-spectra/{id}`
- `DELETE /mass-spectra/{id}`

Recommended behavior:
- Enforce unique constraint semantics (`inchikey`, `db_number`, `source`) with clear 409 responses.
- Return normalized `mZ` and `peaks` arrays in responses.
- Keep current read endpoints unchanged to avoid frontend regression.

## UX and Safety Requirements (MVP)

- Form-level validation with clear inline errors.
- Destructive action guardrail for delete.
- Optimistic update is optional; correctness-first server refresh is acceptable for MVP.
- Basic loading/error/empty states for all CRUD screens.

## Out of Scope (For This MVP)

- Bulk import/edit UX
- Full compound CRUD
- Advanced annotation tooling
- Auth/RBAC (if introduced later, wrap write routes then)

## Delivery Milestones

1. API contracts + backend write handlers for `mass_spectra`.
2. Frontend form screens: create/edit/delete actions.
3. Graph preview integration in create/edit flow.
4. Validation, error states, and regression pass on existing read pages.
