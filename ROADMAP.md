# ProjectVem Roadmap

Planning document derived from the initial product description for tracking milestones and ensuring smart sequencing of work.

**Scheduling assumption:** Week numbering begins at project kickoff; revisit after the Phase 1 retrospective before locking dates beyond Phase 2. Later targets below are provisional and should be reconfirmed once M1 completes.

### Milestone Overview
| Milestone | Target Window | Owner | Notes |
| --- | --- | --- | --- |
| M1 Foundations & Architecture | Weeks 1-4 | Ava (Core Platform Lead) | Lock architecture, rendering decision, plugin ABI skeleton |
| M2 Core Editing Experience | Weeks 5-10 | Ben (Editing Lead) | Ship Vim-parity editing loop and macro flow |
| M3 Language Intelligence & Extensibility | Weeks 11-18 | Casey (Language & Extensibility Lead) | Deliver LSP, syntax pipeline, safe plugin loader |
| M4 User Fluency & Ergonomics | Weeks 19-24 | Diya (Product Experience Lead) | Bundle dependencies and polish UX accelerators |
| M5 Packaging & Community Launch | Weeks 25-28 | Riley (Release & Community Lead) | Prepare installers, docs, community programs |

Use the milestone names verbatim when creating GitHub milestones; open the issues listed in each phase section under the corresponding milestone.

## Phase 1 – Foundations & Architecture
- **Objectives:** Confirm editor scope, define core data structures, and ensure the project is OS-agnostic from day one.
- **Key Tasks:** Draft architecture spec, spike buffer/window models in Go, decide rendering backend, outline plugin API surface, prototype scripting host for Lua/Python/Carrion, and document cross-platform build pipeline.
- **Exit Criteria:** Architecture doc approved, minimal renderer proof validated on macOS/Linux/Windows, build scripts reproducible on CI, and plugin ABI skeleton merged.

### GitHub Milestone Plan
- **Milestone:** `M1 Foundations & Architecture`
- **Target Window:** Weeks 1-4
- **Milestone Owner:** Ava (Core Platform Lead)
- **Issues to create:**
  1. `Draft architecture spec & scope lock` – Owner: Ava – Due Week 2.
  2. `Buffer + window model spike in Go` – Owner: Ben (Editing Lead) – Due Week 3.
  3. `Rendering backend evaluation & prototype` – Owner: Mina (Rendering Engineer) – Due Week 3.
  4. `Cross-platform build pipeline scripts` – Owner: Mina (Rendering Engineer) – Due Week 4.
  5. `Plugin ABI skeleton & host contracts` – Owner: Casey (Language & Extensibility Lead) – Due Week 4.
  6. `Lua/Python/Carrion scripting host prototype` – Owner: Casey – Due Week 4.
- **Exit criteria ownership:**

| Exit Criterion | Owner | Target |
| --- | --- | --- |
| Architecture doc approved | Ava (Core Platform Lead) | Week 4 |
| Minimal renderer proof validated on macOS/Linux/Windows | Mina (Rendering Engineer) | Week 4 |
| Build scripts reproducible on CI | Mina (Rendering Engineer) | Week 4 |
| Plugin ABI skeleton merged | Casey (Language & Extensibility Lead) | Week 4 |

## Phase 2 – Core Editing Experience
- **Objectives:** Deliver a NeoVim-like editing loop that feels native to Vim users.
- **Key Tasks:** Implement normal/insert/visual modes, motion/text-object parity, registers, undo-tree, and macro recording/playback; add basic file I/O, multi-buffer handling, and window splits.
- **Exit Criteria:** Users can open/edit/save multiple files, run macros, and navigate entirely via Vim motions with latency comparable to NeoVim on sample projects.

### GitHub Milestone Plan
- **Milestone:** `M2 Core Editing Experience`
- **Target Window:** Weeks 5-10 (reconfirm after the M1 retrospective)
- **Milestone Owner:** Ben (Editing Lead)
- **Issues to create:**
  1. `Mode engine (normal/insert/visual) state machine` – Owner: Ben – Due Week 6.
  2. `Motion + text object parity suite` – Owner: Ben – Due Week 7.
  3. `Registers and undo-tree subsystem` – Owner: Ava – Due Week 7.
  4. `Macro recording and playback pipeline` – Owner: Ben – Due Week 8.
  5. `File I/O + buffer manager improvements` – Owner: Ava – Due Week 8.
  6. `Window/split layout manager` – Owner: Ben – Due Week 9.
- **Exit criteria ownership:**

| Exit Criterion | Owner | Target |
| --- | --- | --- |
| Open/edit/save multiple files reliably | Ava (Core Platform Lead) | Week 9 |
| Macros runnable end-to-end | Ben (Editing Lead) | Week 9 |
| Navigation parity + latency checks vs NeoVim | Ben (Editing Lead) | Week 10 |

## Phase 3 – Language Intelligence & Extensibility
- **Objectives:** Ship first-class LSP integration, syntax highlighting, and a safe extension system.
- **Key Tasks:** Build LSP client manager (multi-server), diagnostics/completion UI, Treesitter-style syntax pipeline, plugin loader with sandboxing, scripting ergonomics (APIs, events, config reloading), and extension samples.
- **Exit Criteria:** At least two languages demonstrate completions/diagnostics/highlighting; third-party scripts can add commands/keymaps without restart; documentation for plugin authors published.

### GitHub Milestone Plan
- **Milestone:** `M3 Language Intelligence & Extensibility`
- **Target Window:** Weeks 11-18 (adjust based on M1/M2 learnings)
- **Milestone Owner:** Casey (Language & Extensibility Lead)
- **Issues to create:**
  1. `LSP client manager with multi-server support` – Owner: Casey – Due Week 13.
  2. `Diagnostics & completion UI surfaces` – Owner: Diya (Product Experience Lead) – Due Week 15.
  3. `Syntax highlighting pipeline (Treesitter-like)` – Owner: Casey – Due Week 14.
  4. `Plugin loader + sandboxing enforcement` – Owner: Casey – Due Week 16.
  5. `Scripting ergonomics (APIs, events, hot reload)` – Owner: Casey – Due Week 17.
  6. `Sample extensions + author docs` – Owner: Riley (Release & Community Lead) – Due Week 18.
- **Exit criteria ownership:**

| Exit Criterion | Owner | Target |
| --- | --- | --- |
| Two languages demo completions/diagnostics/highlighting | Casey (Language & Extensibility Lead) | Week 18 |
| Scripts can add commands/keymaps without restart | Casey (Language & Extensibility Lead) | Week 17 |
| Plugin author documentation published | Riley (Release & Community Lead) | Week 18 |

## Phase 4 – User Fluency & Ergonomics
- **Objectives:** Remove friction so users can “fly” without manual setup.
- **Key Tasks:** Bundle required fonts/assets, auto-handle dependencies, add command palette, fuzzy file switcher, session restore, theme system, and polish default keybindings.
- **Exit Criteria:** Fresh install runs identically across supported OSes with zero external font/package steps; power-user workflows (jumping files, running commands) are optimized and tested with usability feedback.

### GitHub Milestone Plan
- **Milestone:** `M4 User Fluency & Ergonomics`
- **Target Window:** Weeks 19-24 (tentative)
- **Milestone Owner:** Diya (Product Experience Lead)
- **Issues to create:**
  1. `Bundle fonts/assets + dependency bootstrapper` – Owner: Ava – Due Week 20.
  2. `Command palette implementation` – Owner: Diya – Due Week 21.
  3. `Fuzzy file switcher + recent files memory` – Owner: Diya – Due Week 21.
  4. `Session restore + workspace persistence` – Owner: Ava – Due Week 22.
  5. `Theme system + defaults polish` – Owner: Diya – Due Week 23.
  6. `Keybinding review + usability playtests` – Owner: Diya – Due Week 24.
- **Exit criteria ownership:**

| Exit Criterion | Owner | Target |
| --- | --- | --- |
| Fresh install parity across OSes without extra setup | Ava (Core Platform Lead) | Week 23 |
| Power-user workflows optimized/tested | Diya (Product Experience Lead) | Week 24 |

## Phase 5 – Packaging & Community Launch
- **Objectives:** Prepare for public release and sustainable growth.
- **Key Tasks:** Finalize GPLv2 notices, create installers/packages per OS, set up plugin registry or discovery mechanism, write migration guides, publish roadmap updates, and plan community support channels.
- **Exit Criteria:** Signed release candidates per platform, public documentation site, extension submission guidelines, and community feedback loop established.

### GitHub Milestone Plan
- **Milestone:** `M5 Packaging & Community Launch`
- **Target Window:** Weeks 25-28 (refine once earlier milestones stabilize)
- **Milestone Owner:** Riley (Release & Community Lead)
- **Issues to create:**
  1. `Finalize GPLv2 notices + compliance audit` – Owner: Riley – Due Week 26.
  2. `Installer/packaging scripts per OS` – Owner: Mina (Rendering Engineer) – Due Week 27.
  3. `Plugin registry / discovery portal` – Owner: Casey – Due Week 27.
  4. `Migration + onboarding guides` – Owner: Riley – Due Week 27.
  5. `Public roadmap & release comms plan` – Owner: Riley – Due Week 28.
  6. `Community support channels + moderation process` – Owner: Riley – Due Week 28.
- **Exit criteria ownership:**

| Exit Criterion | Owner | Target |
| --- | --- | --- |
| Signed release candidates per platform | Riley (Release & Community Lead) | Week 28 |
| Public documentation site live | Riley (Release & Community Lead) | Week 27 |
| Extension submission guidelines posted | Casey (Language & Extensibility Lead) | Week 27 |
| Community feedback loop established | Riley (Release & Community Lead) | Week 28 |

## Tracking & Next Actions
- Maintain this file as the single source of truth for milestone status; update exit criteria dates as phases complete.
- Open corresponding GitHub milestones/issues per phase to track fine-grained tasks.
- Revisit sequencing after Phase 1 to incorporate learnings before locking later phases.
