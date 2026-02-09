# Hydragen V2 Master Plan

- GitHub: https://github.com/Yongbeom-Kim/hydragen-v2
- (To Be) Hosted on: https://hydragen.senpailearn.com/v2/

## Tech Stack

- Frontend: Tanstack Start + React
	- React Query | Jotai | MUI Joy | Tailwind
- Backend: Go Backend
- DB: Postgresql
- ETL: Python
- Infrastructure: Hetzner (server) + GCP (DNS)
- Deploy: Docker Compose
- SSL Termination / Reverse Proxy: Caddy


## Technical Documents & TODOs

### Phase 1
- [2026-02-08 - Mass Spec Database from Massbank (~100k molecules)](./features/2026-02-08%20-%20Mass%20Spec%20Database%20from%20Massbank%20(~100k%20molecules))
- CI/CD + Infrastructure + Deployment
- Mass Spec CRUD Frontend (MVP For Graph Display)
- Telemetry Exploration
	- What data can/should we collect? What can we interpret from the data?
- MVP: Question generation + ELO system (same as Hydragen V1)
	- Though, I’m not so sure the MCQ is the right system
- Allow instructors to scope students’ questions to be roughly aligned with curricula
	- Optional: as students master curricula, they slowly expand beyond it.
	- The initial approach in Hydragen V1 is a **mistake**, we should have had a **very small** number of molecules (which can get very accurate difficulty rating), and slowly increase the number of molecules. Quality > Quantity
- Explore & Expand Question Types
- Student-centric UX Exploration
- Instructor-centric UX Exploration
- Explore LLMs & Instructional Scaffolding

### Phase 2
- Model LLM as students to figure out common pitfalls & mistakes
- Reinforcment learning of LLM to understand & solve mass spec
- Expand approach to > M/S to other spectra

