# 2026-02-13 - Mass Spectrometry Skills

This document defines the skill model for Hydragen V2 mass spectrometry learning.
It is the curriculum layer referenced by session authorization (`allowedSkills`).

## Goals

- Represent classroom-relevant learning outcomes as explicit, reusable skills.
- Enable instructors to control session scope by selecting allowed skills.
- Provide predictable mapping from skills to practice content and assessment coverage.

## Non-Goals

- Replacing full chemistry ontologies.
- Supporting free-form, unreviewed skill definitions in production.
- Defining authorization rules (covered in Auth & Permission Model).
- Teaching hardware internals (vacuum systems, ion optics design, detector electronics, instrument assembly, source construction, laser alignment).

## Scope Constraint

Hydragen V2 skill scope is restricted to:

- Data-facing interpretation from spectrum only, or chromatogram + spectrum.
- Instrument-facing decisions that are visible in data outcomes.

Hydragen V2 skill scope excludes:

- Hardware theory and instrument construction details.

## Skill Definition

A skill is a teachable, testable learning outcome tied to mass spectrometry interpretation.

Each skill should be:
- Atomic: one main concept per skill.
- Observable: can be evaluated from student response behavior.
- Mappable: linked to molecules/spectra that exercise the concept.
- Versioned: changes are explicit and traceable.

## Data Model (Proposed)

Core `Skill` fields:
- `id` (stable identifier, immutable)
- `code` (human-readable short code, unique)
- `title`
- `description`
- `domain` (e.g., isotope-patterns, fragmentation, rules)
- `difficulty` (introductory, intermediate, advanced)
- `status` (draft, active, deprecated)
- `version`
- `createdAt`, `updatedAt`

Supporting relationships:
- `Skill -> MoleculeSpectrum` (many-to-many coverage mapping)
- `Skill -> Skill` (optional prerequisite links)
- `Session -> allowedSkills` (selected scope list)

## Canonical Session Scope Rule

- Session curriculum scope is defined only by `allowedSkills`.
- Session-visible spectra are derived from the union of spectra mapped to those skills.
- If `allowedSkills` is empty, the session returns no spectra.

## Operational Core (Compressed)

All spectrum-facing competence compresses to:

1. Extract correct information from peaks.
2. Apply chemical logic to extracted information.
3. Evaluate plausibility of structural claims.
4. Choose appropriate ions for identification/quantification tasks.
5. Recognize ambiguity and limitations.

## Curation and Change Management

- Skill creation/update is instructor/admin controlled.
- Changes to active skills require version increments.
- Deprecating a skill does not remove historical analytics; it only blocks future selection by default.
- Coverage mappings (`Skill -> MoleculeSpectrum`) are reviewed before activation.

## Instructor Experience Requirements

- Instructors can search and select skills while configuring sessions.
- Instructors can progressively add skills over time.
- Instructors should see estimated content impact before saving:
  - number of affected spectra
  - major concept distribution

## Quality Gates

Before a skill is marked `active`:
- Has clear description and examples.
- Has at least one mapped spectrum.
- Has no duplicate meaning with existing active skills.
- Has reviewer sign-off.

## Open Questions

- Should prerequisites be strictly enforced or advisory only?
- Should difficulty levels gate recommendation logic automatically?
- Should deprecated skills stay selectable for legacy sessions?

## Data Availability Tracks

Track A: Current dataset (available now)
- Data available: molecule identity + corresponding mass spectrum.
- Practical scope: spectrum-only interpretation tasks without chromatographic or tandem-MS dependencies.

Track B: Expanded datasets (future)
- Data needed: chromatograms, MS/MS acquisitions, calibration workflows, proteomics search context, and matched cross-instrument runs.
- Practical scope: GC-MS, LC-MS/MS, quantitative reporting, proteomics sequencing, and instrument-comparison tasks.

## Cross-Section Practice Design

Practice should target integrated reasoning across multiple skill areas, not isolated section drills.
Questions should require students to combine peak reading, chemical logic, and evidence evaluation in one workflow.

### Track A Archetypes (Current Dataset)

1. Unknown Identification Challenge
Prompt: Given one spectrum, identify the most plausible compound class or candidate and justify exclusions.
Data: molecular ion/precursor region, isotope cluster, fragment peaks, molecule identity label for grading.
Skill blend: Basic Spectrum Literacy + Isotope Pattern Interpretation + Formula/Rule Application + EI Fragmentation + Critical Evaluation.
LLM augmentation: structured critique on missing evidence links, top alternative hypotheses, and one additional discriminating check request.

2. Competing Hypotheses Adjudication (Spectrum-Only)
Prompt: Two or more structure hypotheses are provided; decide which is best supported from a single spectrum.
Data: one spectrum plus candidate formulas/structures.
Skill blend: Isomer Differentiation + Formula/Rule Application + Fragmentation Interpretation + Critical Evaluation.
LLM augmentation: adversarial argumentation for each hypothesis, then calibrated verdict with confidence and unresolved ambiguity.

3. Deliberate Error Detection
Prompt: Review a worked assignment with subtle interpretation mistakes and identify exactly what is wrong.
Data: annotated peak table and rationale text with intentional errors.
Skill blend: Critical Evaluation + Isotope Interpretation + Formula/Rule Application + Fragmentation Logic.
LLM augmentation: error-type classification (assignment, rule misuse, overclaim, unsupported inference), hint ladder, corrected reasoning chain.

4. Pairwise Isomer Discrimination
Prompt: Compare two spectra and decide whether they support isomer differentiation or are inconclusive.
Data: paired single spectra from known isomeric molecules.
Skill blend: Isomer Differentiation + Fragmentation Interpretation + Critical Evaluation.
LLM augmentation: force evidence citation by peak, then generate a counterargument to stress-test the student conclusion.

5. Escalation Decision (Is Spectrum Alone Enough?)
Prompt: Decide whether evidence from current data is sufficient or whether escalation is required.
Data: borderline spectrum cases with plausible competing interpretations.
Skill blend: Critical Evaluation + Formula/Rule Checks + Isomer Reasoning.
LLM augmentation: suggest escalation actions ranked by expected information gain and explain why current evidence is insufficient.

### Track B Archetypes (Expanded Datasets)

1. Co-Elution and Interference Triage
Prompt: Determine if observed signal is one analyte, mixed analytes, or contamination.
Data: TIC/XIC traces, retention-time slices, spectral snapshots, blanks/controls.
Skill blend: GC-MS/LC-MS Interpretation + Quantitative Reasoning + Critical Evaluation.
LLM augmentation: recommend interference diagnostics and classify likely failure mode.

2. Quant-Identity Coupled Decision
Prompt: Decide whether a result is reportable for both identity and concentration.
Data: calibration series, internal-standard ratios, qualifier/quantifier ions, confirmation transitions.
Skill blend: Quantitative Reasoning + LC-MS/MS Interpretation + Critical Evaluation.
LLM augmentation: rubric-driven pass/fail rationale, failed criteria report, minimum corrective rerun plan.

3. Method-Shift Comparison (EI vs ESI, unit vs high resolution)
Prompt: Same analyte across methods; identify claims that are robust vs method-dependent.
Data: matched spectra across ionization, resolution, and collision settings.
Skill blend: Instrument-Dependent Awareness + Fragmentation Interpretation + Critical Evaluation.
LLM augmentation: side-by-side claim diffing and transferability warnings.

4. MS/MS Pathway Justification
Prompt: Build and defend precursor-to-product interpretation with mechanism plausibility checks.
Data: precursor scan plus product-ion spectra across collision energies.
Skill blend: LC-MS/MS Interpretation + Fragmentation Interpretation + Critical Evaluation.
LLM augmentation: challenge mechanism assumptions and suggest discriminating transitions.

5. Proteomics Identification Challenge
Prompt: Assign likely peptide identity and explain confidence limits.
Data: peptide MS/MS, candidate sequences, database search outputs.
Skill blend: Proteomics Spectrum Skills + Critical Evaluation.
LLM augmentation: explain which ions support the match, flag missing evidence, and separate confident from ambiguous regions.

Design constraints for high-value practice:
- Use realistic ambiguity with at least one plausible distractor explanation.
- Score reasoning quality, not only final-answer correctness.
- Require explicit evidence citation (which peaks/patterns support each claim).
- Include "insufficient evidence" as a valid high-quality outcome when justified.
- Prefer multi-step tasks that mirror real analyst workflows over single-fact recall.
- Keep LLM in feedback/critique mode by default; hide final answers until student commitment unless tutoring mode is enabled.

## Appendix: Spectrum-Based Skill Catalog (Reference)

Informed by: https://doi.org/10.1177/14690667241237431

### A. Skills Supported by Current Dataset (Molecule Identity + Single Mass Spectrum)

This group is feasible now with the existing dataset and no chromatographic or MS/MS dependency.

1. Basic Spectrum Literacy (spectrum-only subset)
- Identify molecular ion (or most likely precursor).
- Distinguish base peak vs molecular ion.
- Recognize presence/absence of M+.
- Determine nominal mass from spectrum.
- Estimate monoisotopic mass from isotope cluster.
- Recognize multiply charged ions from isotope spacing.
- Deconvolute multiply charged ions to neutral mass.
- Identify adducts only when adduct labels/mode context are present in source records.

2. Isotope Pattern Interpretation
- Identify Cl/Br from 3:1 or 1:1 patterns.
- Detect sulfur from M+2 enrichment.
- Count halogens from isotope ratios.
- Determine charge state from peak spacing.
- Predict plausible elemental composition from isotope cluster.
- Distinguish isotope peak vs fragment peak.
- Evaluate whether resolution appears sufficient for isotope-level claims.

3. Formula and Rule Application (Using Observed Mass)
- Apply nitrogen rule.
- Calculate DBE from candidate formula.
- Apply rule of 13.
- Eliminate implausible formulas.
- Compare candidate formulas to isotope-pattern consistency.
- Check hydrogen/carbon ratio plausibility.
- Evaluate whether observed m/z supports a proposed structure.

4. Fragmentation Interpretation (EI-focused subset)
- Identify likely neutral losses from fragment differences.
- Recognize McLafferty rearrangement in EI-like cases.
- Distinguish even-electron vs radical fragment behavior when inferable.
- Identify diagnostic fragments.
- Evaluate whether a proposed pathway is plausible or impossible.
- Distinguish likely primary vs secondary fragments.

5. Isomer Differentiation from Spectra (spectrum-only subset)
- Identify differences in fragmentation patterns.
- Use relative intensities diagnostically.
- Identify unique fragment ions.
- Determine when spectra are insufficient to differentiate isomers.

6. Critical Evaluation Skills (spectrum-only)
- Identify flawed fragmentation rationale.
- Detect misassignment of product/fragment ions.
- Recognize ambiguous interpretation.
- Evaluate whether evidence supports structural claims.
- Compare competing structure proposals against observed data.
- Identify overinterpretation.

### B. Skills Requiring Additional Datasets (Not Fully Supported Yet)

This group requires data types not currently available in the single-spectrum molecule dataset.

1. GC-MS Interpretation Skills
Required datasets:
- Full chromatograms (TIC/XIC) with retention times.
- Spectrum slices across retention windows.
- Blanks/controls and contamination examples.
Skills unlocked:
- Peak selection from chromatogram.
- Spectrum extraction at correct retention window.
- SIM ion strategy and co-elution/background analysis.
- Library-match critique with chromatographic context.

2. LC-MS / MS/MS Interpretation Skills
Required datasets:
- LC-MS full-scan runs with retention information.
- MS/MS acquisitions (precursor/product ions).
- Collision-energy sweeps and transition metadata.
Skills unlocked:
- Precursor/product ion selection.
- In-source vs true MS/MS fragment discrimination.
- Transition choice for confirmation and quantification.
- CE-dependent interpretation.

3. Quantitative Reasoning from Spectra
Required datasets:
- Calibration series with known concentrations.
- Internal-standard metadata and replicate injections.
- Quantifier/qualifier ion traces and acceptance thresholds.
Skills unlocked:
- Calibration construction and linearity checks.
- Ratio-based quantification.
- Interference diagnosis in ion-ratio behavior.
- Reportability decisions for quantitative output.

4. Proteomics-Specific Spectrum Skills
Required datasets:
- Peptide MS/MS spectra with precursor context.
- Sequence database (FASTA) and search-engine outputs.
- Ground-truth or benchmark identifications for evaluation.
Skills unlocked:
- b/y ion interpretation and charge-state handling.
- Peptide assignment and plausibility checks.
- Missing-fragment limitation analysis.
- Peptide-level ambiguity reporting.

5. Instrument-Dependent Spectrum Awareness (full coverage)
Required datasets:
- Matched analytes acquired across EI/ESI and instrument classes.
- Unit-resolution and high-resolution matched runs.
- Controlled acquisition parameter sweeps.
Skills unlocked:
- Method-dependent behavior comparisons.
- Resolution-driven claim limits.
- Transferability judgments across instruments.

6. MS/MS-Dependent Subskills from Other Sections
Required datasets:
- Tandem-MS datasets linking precursor and product ions.
Skills unlocked:
- Fragmentation tree construction from precursor to product ions.
- MS/MS-based isomer differentiation.
- Collision-energy optimization reasoning.
