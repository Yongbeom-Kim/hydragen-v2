# 2026-02-13 - Auth & Permission Model

This document explains the **Authentication** and **Authorization** model of **Hydragen V2**: in short, who can do what.

## **Core Policy**

> [!IMPORTANT]
> **About Classroom Adoption**
> - Introducing a new classroom tool may be difficult if the tool does not adhere to classroom syllabus, or the stated learning outcomes.
> - Therefore, we treat strict **Syllabus Adherence** is a core product requirement, not a nice-to-have.
> - The **Instructor** will control session scope so the tool matches the actual teaching plan. 
> - Instructors first **limit** the freedoms and scope for the problems within Hydragen. Then, we will take over and **personalize instruction** within this defined space.

### **Roles**

Everyone belongs to one of three **Roles**:

1. **Student**
2. **Instructor**
3. **Admin**

Everyone starts as a **Student** (except the SenpaiLearn/Hydragen team, which starts as **Admin**).

- **Admins** can promote Students -> Instructors and Instructors -> Admins.
- **Instructors** can manage Sessions and invite Students.
- **Students** can be invited to Sessions.

> [!NOTE]
> **Technical Details**
> - **Authentication** is handled by **[`authentik`](https://goauthentik.io/)** using **[`OIDC`](https://openid.net/connect/)**.
> - **Authorization** is split between:
>   - **Global Role Claims** (provided by the **Identity Provider (IdP)**)
>   - **Session-Level Rules** (stored in the **Application Database**)

### **Sessions**

A **Session** is essentially a semester-long class.
When an **Instructor** (or **Admin**) creates a **Session**, they define a **Syllabus**.

A **Syllabus** is a collection of **Skills**.

A **Skill** is a specific, actionable learning outcome the student should learn.
In **Mass Spectrometry (M/S)**, some examples of skills are:

1. Identification of isotope patterns
2. Identification of common fragments with mass X
3. Application of the Nitrogen rule

By selecting a collection of **Skills** in a **Session**, instructors can narrow the question set to what is relevant for that course.

Over time, the **Instructor** can progressively broaden the **Skills** included in that **Session**.

### **Session-Scoped Peak Curriculum**

Each **Instructor** can choose which **Skills** are in a particular **Session**.
This selection restricts which molecules and their spectra are available in that **Session**.

Practical effect:
- If a peak is included in the **Session Curriculum**, related spectra can appear.
- If a peak is not included, related spectra are excluded from that **Session**.

This gives instructors direct control over content scope and progression.

> [!NOTE]
> **Technical Details**
> - We maintain an index: **`Skill Usage -> Molecule Spectrum`**.
> - Each **Session** stores its own set of allowed **Skills**.
> - Session-visible spectra are computed from that allowed peak set.
> - Recommended default policy: if no skills are selected, show no spectra (**restricted mode**).


### **Permission Endpoints**

- Only **Instructors** and **Admins** can create and configure sessions.
- **Students** can participate, but cannot change membership or curriculum scope.
- **Admins** handle instructor role promotion/demotion.

> [!NOTE]
> **Technical Details**
> - **`POST /sessions`** -> **Instructor/Admin**
> - **`POST /sessions/{id}/members`** -> **Instructor-of-session/Admin**
> - **`DELETE /sessions/{id}/members/{userId}`** -> **Instructor-of-session/Admin**
> - **`PUT /sessions/{id}/allowed-peaks`** -> **Instructor-of-session/Admin**
> - **`GET /sessions/{id}/allowed-peaks`** -> **Session Member/Admin**
> - **`GET /sessions/{id}/mass-spectra`** -> Session-filtered spectra only
> - **`POST /admin/roles/instructors/{userId}`** -> **Admin only**
> - **`DELETE /admin/roles/instructors/{userId}`** -> **Admin only**
