# 2026-02-13 - Question Generation Recommendation (Track A)

This document defines how Hydragen V2 should generate and recommend questions per student for Track A only.
Goal: personalized practice that adapts to mastery and uncertainty, without relying on a single global difficulty score.

## Related Document

- [2026-02-13 - Mass Spectrometry Skills](./2026-02-13%20-%20Mass%20Spectrometry%20Skills.md)
- [Hydragen Plan](../PLAN.md)
- [2026-02-13 - Telemetry Exploration](./2026-02-13%20-%20Telemetry%20Exploration.md)

## Problem Statement

A fixed question set or one-size-fits-all queue is not acceptable.
Students differ in:
- Skill strengths and weaknesses.
- Confidence calibration.
- Pace and frustration tolerance.

The recommendation system must therefore optimize for:
- Learning gain per student.
- Productive difficulty (challenging but not discouraging).
- Coverage over required session skills.

Track A scope constraint:
- Question generation is limited to inputs with molecule identity + single mass spectrum.
- No chromatogram-dependent, MS/MS-dependent, quantitative workflow, or proteomics workflow questions in this document.

## Non-Goals

- Pure Elo-only ranking as the primary model.
- Opaque black-box recommendation with no auditability.
- Personalization outside instructor-selected session scope.

## Core Policy Constraint

> [!IMPORTANT]
> Recommendation happens only inside session `allowedSkills`.
> No question should be generated or served outside the authorized skill scope.
> This document applies only to Track A data availability.

## Data Availability Scope

Track A (in scope):
- Molecule identity + single mass spectrum.
- Supports spectrum-only question generation and recommendation.

Track B (out of scope in this document):
- Requires additional datasets (chromatograms, MS/MS, quant workflows, proteomics context).
- Should be specified in a separate Track B recommendation document.

## Recommendation Model (Hybrid, Not Elo-Only)

Maintain a per-student state vector, updated after every attempt:
- `mastery[skill]`: estimated competence per skill.
- `uncertainty[skill]`: confidence interval for each mastery estimate.
- `engagement_state`: struggle, flow, or boredom indicators.
- `recent_history`: last N attempts (correctness, latency, hints, retries, confidence).

Candidate questions are scored with a multi-objective policy:
- `learning_gain_score`: expected improvement on weak or uncertain skills.
- `difficulty_fit_score`: target zone relative to current mastery.
- `coverage_score`: reinforces session-required skill distribution.
- `novelty_score`: avoids near-duplicate repetition.
- `fatigue_penalty`: down-rank when repeated failure risk is high.

Final ranking:
- Weighted sum with configurable policy weights per course/session.
- Top-K sampled with controlled randomness to prevent overfitting to one path.

## Engagement Measurement Spec

`engagement_state` is treated as an operational learning-risk signal, not a psychological diagnosis.
It must be inferred only from observable product behavior.

States:
- `flow`: challenge is appropriate and progress is stable.
- `productive_struggle`: errors occur, but persistence and recovery are healthy.
- `frustration_risk`: repeated failure with worsening support dependence.
- `disengagement_risk`: avoidance patterns (skip/abandon/inactivity) are rising.

### Observable Inputs (Per Attempt)

- `is_correct`
- `latency_ms`
- `hint_count`
- `attempt_index` (within a question)
- `abandoned` (started but not submitted)
- `skipped`
- `confidence_self_report` (optional)
- `question_expected_time_ms` (from template/archetype baseline)

### Rolling Features (MVP Window: Last 10 Attempts)

- `latency_ratio`: median(`latency_ms / question_expected_time_ms`)
- `hint_rate`: total hints / attempts
- `retry_depth`: mean `attempt_index` for completed questions
- `abandon_rate`: abandoned / started
- `skip_rate`: skipped / assigned
- `accuracy_trend`: slope of correctness over window
- `recovery_after_error`: probability(next attempt correct | previous attempt incorrect)
- `confidence_gap`: mean(`confidence_self_report - is_correct`) when confidence exists
- `return_gap_ratio`: time since last session / student's median inter-session gap

### Normalization

- Primary normalization is student-relative z-score over trailing 30-day baseline.
- Cold start (<20 attempts): blend student signal with cohort priors.
- Cap feature z-scores to `[-3, +3]` to limit outlier effects.

### Composite Scores

Intermediate signals:
- `persistence_signal = z(retry_depth) - z(abandon_rate) - z(skip_rate)`
- `focus_signal = -abs(z(latency_ratio))` (very slow or very fast both reduce focus)
- `recovery_signal = z(recovery_after_error) - z(hint_rate)`
- `frustration_signal = z(abandon_rate) + z(skip_rate) + z(hint_rate) - z(accuracy_trend)`

Overall score:
- `engagement_score = 0.30*persistence_signal + 0.25*focus_signal + 0.25*recovery_signal - 0.20*frustration_signal`
- Smooth with EWMA: `E_t = 0.6*E_(t-1) + 0.4*engagement_score`

### State Classification (MVP Thresholds)

- `flow`: `E_t >= +0.50` and `frustration_signal < +0.75`
- `productive_struggle`: `-0.25 <= E_t < +0.50` and `persistence_signal >= 0`
- `frustration_risk`: `E_t < -0.25` and `frustration_signal >= +0.75`
- `disengagement_risk`: `skip_rate >= 0.30` or `abandon_rate >= 0.25` or `return_gap_ratio >= 2.5`

If multiple conditions match:
- `disengagement_risk` overrides all.
- Then `frustration_risk`.
- Then `flow` vs `productive_struggle`.

### Recommendation Actions by State

- `flow`: slightly increase ambiguity or reasoning depth; keep pace.
- `productive_struggle`: maintain difficulty; add targeted feedback and contrastive examples.
- `frustration_risk`: reduce ambiguity/noise, shorten task length, increase scaffolded hints.
- `disengagement_risk`: serve short high-success re-entry tasks and reduce consecutive hard items.

### Fairness and Governance Checks

- Never use protected attributes in state estimation.
- Monitor state-label prevalence by cohort slices for disparate impact.
- Log feature inputs and decision outputs per recommendation for audit.
- Keep this telemetry aligned with GDPR/IRB controls in `Telemetry Exploration`.

## Skill Sufficiency Check (Algorithmic)

Goal:
- Verify that a generated question is solvable using only its declared required skills.
- Reject questions with hidden skill dependencies or unnecessary required-skill tags.

### Required Inputs

- `required_skills[]` from template/question contract.
- `evidence_operations[]` used by the canonical solution path.
- `operation_to_skill_map` (each operation mapped to exactly one skill).
- `operation_dependency_dag` for the question archetype.
- Bound spectrum data and canonical answer criteria.

### Procedure

1. Build canonical solution graph:
- Instantiate the archetype DAG with bound data.
- Keep only operations that are actually used for a valid answer path.

2. Skill projection:
- Map each used operation to skill via `operation_to_skill_map`.
- Compute `derived_required_skills = union(mapped_skills)`.

3. Sufficiency check:
- Restrict graph to operations whose skills are in declared `required_skills`.
- Pass only if at least one full path reaches a valid terminal answer, no dependency node is missing on that path, and no operation outside declared `required_skills` is required.

4. Necessity ablation (minimality):
- For each skill `s` in declared `required_skills`, remove all operations mapped to `s`.
- Re-run reachability.
- If a valid path still exists, mark `s` as non-necessary.

5. Classification:
- `sufficient = true` if step 3 passes.
- `minimal = true` if all declared skills are necessary from step 4.
- `missing_skills = derived_required_skills - declared_required_skills`
- `extraneous_skills = declared_required_skills - minimal_required_skills`

### Output Contract

- `sufficient` (bool)
- `minimal` (bool)
- `minimal_required_skills[]`
- `missing_skills[]`
- `extraneous_skills[]`
- `validation_confidence` (high when deterministic checks fully cover solution path)
- `validation_notes`

### Release Gate

- Block publication if `sufficient = false`.
- Route to review if `minimal = false` or `validation_confidence < threshold`.

## Procedural Question Generation

Each question is generated from a constrained template + data binding, not free-form only.

Generation pipeline:
1. Select target skill bundle within `allowedSkills`.
2. Select suitable molecule/spectrum examples from mapped coverage.
3. Instantiate archetype constraints (difficulty, distractors, ambiguity level).
4. Run LLM-assisted draft generation.
5. Validate against rule checks before publication.

Hard validation gates:
- Scope check: all required skills are in `allowedSkills`.
- Data check: referenced peaks/patterns exist in bound spectrum data.
- Answerability check: question has at least one defensible solution path.
- Difficulty check: predicted difficulty falls in requested band.
- Safety check: no hallucinated lab procedure claims for unavailable datasets.

## Difficulty and Progression

Difficulty should be multidimensional:
- Signal complexity (noise, overlap, isotope clarity).
- Reasoning depth (single-step vs multi-step evidence chains).
- Ambiguity load (clear answer vs competing plausible hypotheses).
- Error diagnosis complexity (obvious vs subtle flawed rationale).

Progression policy:
- Default target: `70-85%` expected success probability.
- If student is consistently above target, increase ambiguity/reasoning depth.
- If consistently below target, reduce cognitive load but keep core skill focus.
- Use spaced reinforcement for recently improved but unstable skills.

## Track-Specific Recommendation Strategy

### Track A (Now): Spectrum-Only Personalization

Recommended archetypes:
- Unknown identification challenge.
- Competing hypotheses adjudication.
- Deliberate error detection.
- Pairwise isomer discrimination.
- Escalation decision (is spectrum evidence sufficient?).

Key adaptation signals:
- Misread isotope clusters -> increase isotope-focused variants.
- Formula-rule mistakes -> increase constrained formula elimination tasks.
- Overconfident wrong answers -> add evidence-citation requirements.
- Slow but correct -> reduce noise, keep complexity, improve fluency.

## Track A Archetype Specifications

Each archetype must define:
- `required_skills[]`, `optional_skills[]`, `forbidden_skills[]`
- `evidence_operations[]`
- canonical answer and rubric
- sufficiency/necessity validation results

### Shared Pre-Answer Workspace Pattern (Q1/Q2)

Before final answer submission, learner gets a tap-only "working memory" layer (scratch mode).
This is analogous to Sudoku pencil marks and is required for Q1 and Q2.

Behavior:
- Learner can create, edit, and remove tentative evidence links:
- tap peak/pattern region
- tap functional group on an option/hypothesis card
- assign tentative tag (`supports`, `rules out`, `shared`, `uncertain`)
- System does not force final answer first.
- Learner commits final answer only after workspace stage.

Telemetry and scoring policy:
- Workspace actions are logged as process signals (reasoning path), not direct correctness labels.
- Mastery updates from workspace are lower-weight than committed-answer outcomes.
- Product should infer "on-track/off-track" tendencies from workspace structure:
- alignment between tentative links and solver-derived discriminating operations
- correction behavior (self-repair before submit)
- contradiction density in tentative graph
- Do not penalize exploratory edits heavily; penalize persistent unresolved contradictions at commit.

### 1) Unknown Identification Challenge

Generation:
1. Select one spectrum with adequate signal quality and known molecule label.
2. Sample `2-4` candidate identities/classes including at least one plausible distractor.
3. Build two-phase delivery for this archetype:
4. Phase 1 (Iteration 1): pure multiple choice selection only.
5. Phase 2 (Iteration 2): pre-answer tap-only workspace, then answer commit with scaffolded exclusion workflow.
6. Bind expected evidence anchors (molecular ion region, isotope cluster, key fragments).
7. Derive Phase 2 prompt obligations from solver outputs:
8. Run sufficiency/necessity solver to compute `minimal_required_skills` and discriminating operations.
9. Generate required evidence chips only for those discriminating operations (no static chip defaults).

Phase 2 UX contract (mobile-first, no text input):
- Input mode is tap/click only; no free-text field.
- Shared interaction primitive (Q1/Q2): learner creates an evidence link by:
- tap peak/pattern region in spectrum
- tap functional group on a molecule option card
- assign one structured claim tag to that link
- Supported tags:
- `supports selected candidate`
- `rules out candidate A/B/C`
- `unique-to-selected`
- `non-diagnostic`
- System provides guided prompts in sequence:
- Step 1: "Tap a peak and functional group that support your chosen answer."
- Step 2: "For each alternative, tap a peak and functional group that rules it out."
- Step 3: "Tap a peak-group pair that is diagnostic/unique for your chosen answer (or mark none)."
- Completion rule: learner must satisfy solver-derived required chips before final answer submit.

Sufficiency and necessity checks:
- Phase 1 sufficiency pass condition: using declared skills, solver can identify the best-supported candidate from options.
- Phase 2 sufficiency pass condition: solver can complete scaffold steps with correctly tagged evidence chips that support selection and rule out distractors.
- Necessity test: ablate each required skill; if Phase 1 selection or Phase 2 evidence-tag quality remains unchanged, the ablated skill is not necessary.
- Reject if any successful path requires Track B operations.
- Prompt-obligation check: every required learner prompt must map to a solver-identified differentiating operation.
- Over-constraint check: reject required prompts not implied by `minimal_required_skills`.

Phase 2 validation rules:
- Peak-anchor validity: tapped coordinates must map to known local peak windows or approved isotope clusters.
- Functional-group validity: tapped group must belong to the selected molecule option and map to known structural annotation.
- Claim-tag validity: each tag must satisfy archetype-specific checks (for example, `rules out candidate B` requires mismatch with B's expected pattern set).
- Peak-group coherence validity: linked peak-group pair must be chemically plausible under the template rule set.
- Coverage validity: required support/rule-out coverage must match solver-derived obligations.
- Contradiction validity: reject mutually inconsistent tag sets (for example, same chip tagged as both `unique-to-selected` and `non-diagnostic`).

Augmentations:
- Phase 1: confidence calibration prompt after option commit.
- Phase 2: counterfactual distractor critique ("what would be true if candidate B were correct?").
- Phase 2: evidence-coverage scoring that penalizes missing rule-out coverage.
- Phase 2: adaptive scaffolding that unlocks guided hint chips when learner stalls.
- Phase 2: workspace health indicator (missing evidence type, unresolved contradiction, or balanced coverage).

### 2) Competing Hypotheses Adjudication

Generation:
1. Choose one spectrum and `2-3` explicit structural/formula hypotheses.
2. Ensure at least two hypotheses are reasonably plausible at first glance.
3. Require pre-answer tap-based workspace, followed by ranked adjudication (no free-text required).
4. Include optional "insufficient evidence" outcome when ambiguity is genuine.
5. Derive required prompt obligations from solver outputs:
6. Run sufficiency/necessity solver to compute `minimal_required_skills` and discriminating operations.
7. Generate scaffold prompts only for those discriminating operations (no static prompt defaults).

UX contract (tap-first):
- Learner orders hypotheses (`best supported`, `second`, `least supported`).
- Shared interaction primitive (Q1/Q2): learner creates an evidence link by:
- tap peak/pattern region in spectrum
- tap functional group on a hypothesis molecule card
- assign a structured adjudication tag to that link
- Supported tags:
- `supports hypothesis X`
- `contradicts hypothesis X`
- `shared/non-discriminating`
- `insufficient-to-decide`
- Prompt flow is guided per hypothesis:
- "Tap a peak-group pair that supports this hypothesis."
- "Tap a peak-group pair that contradicts this hypothesis."
- If solver marks case ambiguous, require at least one `insufficient-to-decide` tag before submit.

Sufficiency and necessity checks:
- Sufficiency pass condition: declared skills yield a stable top-ranked hypothesis or justified ambiguity decision.
- Necessity test: remove each required skill and recompute ranking confidence; if confidence/ordering is unaffected, mark that skill extraneous.
- Validate that adjudication does not depend on unavailable metadata.
- Prompt-obligation check: every required learner prompt must correspond to a solver-identified discriminating operation.
- Over-constraint check: reject prompts that require operations not in `minimal_required_skills`.

Validation rules:
- Ranking validity: learner must provide complete ordering or explicit ambiguity decision.
- Functional-group validity: selected group must belong to the hypothesis molecule card being evaluated.
- Evidence-tag validity: each support/contradict tag must map to hypothesis-specific expected/forbidden pattern sets.
- Peak-group coherence validity: linked peak-group pair must pass archetype rule checks.
- Discrimination validity: at least one chip must distinguish between top two hypotheses unless canonical outcome is ambiguous.
- Consistency validity: reject self-contradictory adjudication graphs.

Augmentations:
- LLM adversarial argument for each hypothesis before final verdict.
- Automatic contradiction checker between evidence tags and cited peaks.
- "One extra test" suggestion generation (kept conceptual, no Track B dependency).
- Adaptive prompting: if learner stalls, surface next best solver-derived discriminating check.
- Workspace graph replay for learner review before final ranking commit.

### 3) Deliberate Error Detection

Generation:
1. Start from a valid worked solution trace.
2. Inject `1-3` controlled errors (peak misassignment, rule misuse, overclaim, logic leap).
3. Render trace as a non-linear evidence graph (not strict linear bullets):
4. Nodes: peaks/patterns, functional groups, candidate claims.
5. Edges: `supports`, `rules_out`, `shared`, `uncertain`.
6. Ask learner to locate wrong edges/nodes and repair them via tap interactions only.

UX contract (interactive, mobile-first):
- Tapping a peak opens a local discussion panel ("peak thread") for that peak.
- Panel shows all claim edges currently linked to that peak (correct or incorrect).
- Learner can react per edge with:
- `thumbs_up` (looks correct)
- `thumbs_down` (looks wrong)
- If `thumbs_down`, learner chooses error chip:
- `wrong_peak`
- `wrong_functional_group`
- `wrong_rule`
- `overclaim`
- `contradiction`
- After `thumbs_down`, UI presents multiple-choice interpretation options for repair.
- Options must include:
- plausible corrected interpretations from solver candidate set
- `no valid interpretation` (peak should not be identified as diagnostic)
- Learner repairs by selecting one option, then re-linking peak/group/claim tag if needed, or deleting edge.
- No free-text required; optional quick-reason chips only.
- Input parity requirement: all actions must support tap, mouse click, and keyboard navigation/selection.

Failure modes to inject and test:
- FM1 `missing_peak_identification`: sample solution omits a peak/pattern that is solver-required; learner must add it.
- FM2 `incorrect_peak_attribution`: sample solution identifies a peak but maps it to wrong interpretation; learner must downvote and correct or choose `no valid interpretation`.

Sufficiency and necessity checks:
- Sufficiency pass condition: declared skills are enough to detect all injected errors and reconstruct a correct minimal reasoning chain.
- Necessity test: for each required skill, verify at least one injected error is undetectable when that skill is removed.
- Reject if any injected error requires external lab-method knowledge.
- Graph-completeness check: repaired graph must satisfy solver-required discriminating operations.
- Non-linearity check: do not require a single fixed step order if multiple valid correction paths exist.

Validation rules:
- Error-detection validity: learner flags all injected errors (or meets threshold in partial-credit mode).
- Diagnosis validity: chosen error chips match injected error classes within allowed equivalence map.
- Repair validity: corrected edge set must be chemically plausible and solver-consistent.
- Contradiction validity: final graph contains no unresolved conflicting edges.

Augmentations:
- Hint ladder by error type (light cue -> targeted cue -> near-explicit cue).
- Error-taxonomy feedback with personalized remediation mapping.
- Partial-credit scoring across detect/diagnose/repair stages.
- Peak-thread prioritization: system highlights highest-impact ambiguous peak when learner stalls.

### 4) Pairwise Isomer Discrimination

Generation:
1. Select a known isomer pair with two corresponding spectra.
2. Control similarity level to tune difficulty (easy: clear diagnostic ions; hard: subtle intensity pattern differences).
3. Ask for decision: assign which spectrum matches which structure or conclude inconclusive.
4. Require explicit comparison table of discriminating vs shared evidence.

Sufficiency and necessity checks:
- Sufficiency pass condition: declared skills support either a correct assignment or a justified inconclusive judgment.
- Necessity test: remove each required skill; if discriminating evidence remains fully actionable, downscope required skills.
- Enforce ambiguity guardrail: if no robust discriminator exists, canonical answer must allow "inconclusive."

Augmentations:
- LLM-generated counterargument to stress-test overconfident assignments.
- Structured "difference-first" scaffold (what differs, why it matters, confidence).
- Uncertainty-aware scoring that rewards justified restraint.

### 5) Escalation Decision (Is Spectrum Alone Enough?)

Generation:
1. Select borderline cases where multiple interpretations remain plausible.
2. Ask learner to decide whether current evidence is sufficient for claim strength level.
3. Require explicit claim boundary: what is supported vs not supported.
4. Require an escalation recommendation framed as information need (not hardware procedure).

Sufficiency and necessity checks:
- Sufficiency pass condition: declared skills allow correct boundary-setting between supported and unsupported claims.
- Necessity test: remove each required skill and check whether boundary-setting quality collapses.
- Reject if canonical answer implies certainty where data supports only ambiguity.

Augmentations:
- LLM-generated "review board" critique focused on overreach and missing caveats.
- Claim-strength rubric (`confirmed`, `probable`, `tentative`, `unsupported`).
- Meta-cognitive prompt: confidence vs evidence gap reflection.

## LLM Augmentation Policy

LLMs are used as constrained generators and critics, not unrestricted answer engines.

Allowed LLM roles:
- Generate candidate distractors and rationale stubs.
- Produce targeted feedback after student commitment.
- Simulate adversarial counter-hypotheses for robustness.
- Classify error types and recommend next practice focus.

Guardrails:
- Require explicit evidence citation by peak/pattern.
- Keep hidden-answer mode on by default.
- Log model version and prompt template version for every generated question.
- Route low-confidence generations to review queue.

## Molecule and Spectrum Similarity (Track A)

Goal:
- Retrieve "nearby" molecules/spectra for distractor generation, contrastive practice, and novelty control.
- Support both chemical-space similarity and observed-spectrum similarity.

### Input Standardization

Molecule identity normalization:
- Use canonical SMILES as primary structure key.
- Keep InChI/InChIKey as stable cross-dataset identifier.
- Deduplicate records by InChIKey (same molecule with multiple aliases).

Spectrum normalization (single spectrum per item in Track A):
- Convert peaks to centroid list: `(m/z, intensity)`.
- Remove very low-intensity noise peaks (for example `<1%` base peak).
- Normalize intensities to base-peak `100` (or L2-normalize for cosine scoring).
- Keep instrument/context metadata when available (ionization mode, nominal resolution); do not compare across incompatible contexts by default.

### Molecule Similarity

Primary chemical similarity (MVP default):
- Compute Morgan/ECFP fingerprints from canonical SMILES.
- Use Tanimoto similarity between bit vectors.
- Default fingerprint config for MVP: radius `2`, `2048` bits.

Secondary molecule similarity signals:
- Exact mass proximity score.
- Formula overlap score (shared elements + count distance).
- Bemis-Murcko scaffold match flag (same scaffold vs different scaffold).
- Functional-group overlap score (from deterministic substructure tags).

Recommended molecule similarity blend:
- `mol_sim = 0.70*tanimoto_ecfp + 0.15*scaffold_flag + 0.10*formula_score + 0.05*mass_proximity`
- Keep blend weights configurable per archetype.

### Spectrum Similarity

Peak alignment:
- Match peaks using m/z tolerance windows:
- high-resolution context: ppm tolerance (for example `10 ppm`)
- unit-resolution context: absolute tolerance (for example `0.3 Da`)

Primary spectral similarity (MVP default):
- Weighted cosine similarity on aligned peak vectors.
- Use intensity transform to reduce base-peak dominance (`sqrt(intensity)` or log-scale).

Secondary spectral similarity signals:
- Shared top-N peak overlap (precision/recall style).
- Spectral entropy similarity for distribution-shape robustness.
- Isotope-cluster agreement score (presence/spacing/intensity-ratio consistency).

Recommended spectral similarity blend:
- `spec_sim = 0.65*cosine + 0.20*topN_overlap + 0.10*entropy_sim + 0.05*isotope_agreement`

### Spectrum Embeddings (Vector Representation)

Yes, Track A can use learned vector embeddings for spectral similarity.

MVP representation (non-neural embedding):
- Convert each spectrum to a fixed-length vector via m/z binning + intensity transform + normalization.
- Use this vector directly in ANN retrieval with cosine similarity.
- Treat this as baseline embedding even before model training.

Learned embedding path (phase 2+):
- Train an encoder that maps peak lists to dense vectors where similar spectra are nearby.
- Candidate training objectives:
- contrastive loss (positive/negative spectrum pairs)
- triplet loss (anchor/positive/negative)
- autoencoder pretraining followed by metric-learning fine-tuning
- Optional multimodal objective: align spectrum embeddings with molecule fingerprint/SMILES embeddings.

Practical scoring integration:
- `spec_sim_final = w1*spec_sim + w2*embedding_cosine`
- Start with `w1=1.0, w2=0.0` (no learned model), then gradually increase `w2` after offline validation.
- Keep deterministic checks (isotope/peak-rule coherence) as hard guards even when embedding score is high.

Operational guardrails:
- Train and evaluate per compatible instrument context when possible (polarity/resolution/adduct regime).
- Version embedding model and preprocessing pipeline; log both in recommendation decisions.
- Monitor nearest-neighbor quality with instructor spot checks and false-neighbor rate.

### Joint Similarity for Retrieval

For Track A question generation:
- Retrieve candidates with joint score:
- `joint_sim = alpha*mol_sim + (1-alpha)*spec_sim`
- Default `alpha = 0.5`; tune by archetype:
- Q1/Q2 distractors: bias toward spectral confusion (`alpha 0.3-0.5`)
- Q4 isomer discrimination: bias toward structural closeness (`alpha 0.6-0.8`)

Hard constraints before ranking:
- Same ionization polarity/adduct family when required.
- Similar resolution regime.
- Exclude exact duplicate spectrum unless intentionally used for control items.

### "Similar Molecules" vs "Similar Schemas"

Definitions:
- Similar molecules: high `mol_sim` under structural features.
- Similar schemas: similar reasoning pattern requirements, even if molecules differ.

Schema similarity representation:
- Represent each question as a sparse schema vector:
- required skills
- evidence operation pattern (for example isotope-ratio check + neutral-loss check)
- ambiguity level
- distractor type pattern
- expected reasoning depth

Schema similarity scoring:
- Jaccard/cosine over schema vectors.
- Use schema similarity to avoid repetitive cognitive patterning, even when molecules are new.

Recommendation policy use:
- Enforce novelty by penalizing high similarity to recent attempts:
- `novelty_penalty = max(joint_sim_recent, schema_sim_recent)`
- Prefer questions that are close on target skill but not near-duplicate on both molecule and schema.

### Retrieval Architecture (Practical MVP)

Indexes:
- Molecule ANN index over fingerprint vectors (or exact search if dataset size allows).
- Spectrum ANN index over embedded/aligned spectral vectors.
- Schema index over sparse template features.

Two-stage retrieval:
1. Candidate fetch by target skills + hard metadata filters.
2. Re-rank by joint similarity + schema novelty + difficulty fit.

### Quality Gates

- Calibration checks:
- Top-K neighbor sanity review by instructor for each archetype.
- Track false-neighbor rate (high score but clearly different interpretation behavior).
- Drift checks:
- Recompute similarities when preprocessing pipeline/version changes.
- Log versioned similarity parameters in `RecommendationDecisionLog`.

## Minimal Data Model (Proposed)

`QuestionTemplate`
- `id`, `archetype`, `required_skills`, `optional_skills`
- `input_requirements` (track A or B dataset types)
- `difficulty_profile`
- `rubric_schema`

`GeneratedQuestion`
- `id`, `template_id`, `session_id`
- `bound_data_refs` (molecule/spectrum IDs)
- `target_skills`
- `predicted_difficulty`
- `llm_model_version`, `prompt_version`
- `validation_status`

`StudentSkillState`
- `student_id`, `skill_id`
- `mastery`, `uncertainty`
- `last_practiced_at`
- `recent_error_modes`

`RecommendationDecisionLog`
- `student_id`, `session_id`, `timestamp`
- `candidate_ids`
- `selected_id`
- `score_breakdown`
- `policy_version`

`WorkspaceEvidenceEvent`
- `student_id`, `session_id`, `question_id`, `timestamp`
- `phase` (`workspace` or `commit`)
- `peak_ref`, `molecule_option_id`, `functional_group_id`
- `tag` (`supports`, `rules_out`, `shared`, `uncertain`)
- `event_type` (`add_link`, `edit_tag`, `remove_link`, `commit`)
- `alignment_to_solver_ops` (bool/score)

## Evaluation Metrics

Learning quality:
- Mastery lift per skill over time.
- Time-to-stable-mastery.
- Retention after delay.

Recommendation quality:
- Target-band hit rate (actual success vs expected).
- Productive struggle rate (not too easy, not repeated failure).
- Skill coverage completeness within session goals.

Generation quality:
- Validation pass rate.
- Instructor override/reject rate.
- Student-reported clarity and fairness.

## Rollout Plan

1. Implement Track A recommender with template-constrained generation.
2. Add decision logging + telemetry for policy tuning.
3. Run shadow evaluation against non-personalized baseline.
4. Enable progressive rollout by session.
5. Publish a separate Track B spec once required datasets are integrated.

## Open Questions

- Should instructors set policy weights (difficulty vs coverage vs novelty), or use global defaults only?
- What is the max acceptable LLM-generated question rate before mandatory human review?
- Should recommendation include explicit confidence explanations to students, instructors, or both?
