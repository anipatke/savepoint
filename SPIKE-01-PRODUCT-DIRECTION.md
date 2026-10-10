# Spike 01: Product Direction & Audience

Date: 2026-10-10 · Type: product discovery spike (no code changed) · Repo state: `de822fa` (v2.2.1)

**Evidence labels used throughout**

- **[Verified]**: checked directly in this repository, by running Savepoint, or against a cited primary source.
- **[Inference]**: a reasonable conclusion drawn from verified facts.
- **[Hypothesis]**: plausible but untested. Needs real users to confirm.

---

## 1. Executive summary

Savepoint's strongest asset is **independent checking**: a separate agent session tests the work against what was asked, and the owner decides whether it's done. On Savepoint's own development, 33 of 63 Checks came back NEEDS WORK, and they found real defects. No other part of the product has evidence that strong.

Its biggest problem isn't the terminal. **Savepoint has no outside users we can see:** 1 GitHub star, 0 forks, and no issues ever filed by anyone outside the project. Thirty-one Objectives and 140 Issues have all come from the owner's own use. The product has been tuned for one user, and it has grown about sevenfold since May while nothing tested whether anyone else wants it.

The stated audience (vibe coders, "bored dads") does not match the product. The product adds about 13,000 words of agent instructions to a repo, and its Check reports are dense engineering documents.

**Recommendation: Rework.** Narrow the audience to solo developers who already use a terminal coding agent. Make "did the agent actually do what I asked?" the headline. Freeze new capability work until real people outside the project have tried a lighter version. Do not build a visual app or orchestration layer now.

---

## 2. Current product assessment

### What Savepoint really is today [Verified]

- A Go CLI plus a terminal board (Bubble Tea): about 38,000 lines of non-test Go and 61,000 lines of tests. A May 2026 audit counted about 5,600 production lines ([consolidated audit](project-audit/consolidated-audit-report.md)).
- A set of Markdown/YAML records under `.savepoint/`, plus four agent "skills" and three references that do most of the actual workflow work. **Most of the product is the instructions an agent reads, not the binary.** The binary loads records, works out the next step (`resume`), enforces gates, and shows the board.
- 404 commits since 2026-04-27. Of those, 309 landed in September and October.

### What it genuinely does well [Verified unless marked]

1. **Separates building from checking.** A Check must come from a different session than the one that built the work. It writes a permanent record, and only the owner marks a Task done. On Savepoint itself, 33 of 63 Check records say `NEEDS WORK`. A sample showed real defects, for example C-963: "Three reproducible defects and one unverified native-platform gate", and C-952 found the health display "does not yet meet the Objective" in six specific ways.
2. **Keeps state outside the chat.** `savepoint resume` gives a single copyable Next line, so a fresh agent session can pick up the work. The selection logic is deterministic and covered by tests.
3. **Gives REPLAN REQUIRED as a real exit.** An agent is allowed to stop and say the plan is wrong instead of improvising.
4. **Is honest about its limits.** For example, Code Health shows "Unknown" rather than "fine" when a report is missing.

### What differentiates it from task tracking plus AGENTS.md [Inference]

A plain task list and an instruction file can say "write tests" and "plan first". They do not, by default:
- require a *separate* session to verify the work against written acceptance criteria,
- keep "the checker cleared it" separate from "the owner accepted it",
- stop at a recorded gate when the plan breaks.

That combination is the real difference. Planning documents and a kanban board are not, because many tools already provide them (see §4).

### Where users hit friction [Verified from a fresh `savepoint init` in a scratch repo]

- **The first Next line is wrong for a new project.** On a fresh project, the router is in the `idea` state with no Objectives, but `resume` says *"Select an Objective: press p on the board"*. No Objective exists yet. (The clipboard prompt from `init` does point to the Idea skill, so the two messages disagree.)
- **Commands are inconsistent.** `savepoint init .` and `savepoint resume .` accept a directory, but `savepoint doctor .` fails with "takes no positional arguments".
- **Heavy setup in the user's repo.** The scaffold adds about 13,100 words of agent guidance: the scaffold `AGENTS.md` (2,205 words) plus the skills and references. Every agent session pays some of that context cost. The current Goal G-002 exists to reduce it.
- **Many concepts to learn.** Goal → Objective → Task, plus Check, Issue, waiver, `requires: clear` versus `requires: accepted`, R-### versus G-### identities, lanes, Code Health. `AGENTS.md` spends paragraphs on Goal migration rules a new user never needs.
- **Check records are hard to read for the stated audience.** They use coverage matrices, "frozen scope locks" and test names. That works for an engineer and is impenetrable to a vibe coder.

### What requires terminal familiarity [Verified]

- Install and run: `npx`, Node, Git.
- **Owner-only decisions happen on the TUI board**: marking a Task done, accepting an Issue, choosing a Goal, Advanced Options. The alternative is hand-editing YAML that includes actor and timestamp fields.
- Code Health setup asks the user to install `lizard`, `jscpd`, and `osv-scanner` and to configure report output.
- Everything else is done by the agent. **The terminal burden falls on owner decisions, not daily use.**

### Real versus aspirational [Verified]

| Claim | Reality |
|---|---|
| Website: "verifies work with deterministic checks separate from the builder agent" | The Check *verdict* is an AI session following `check-method.md`. Only the gates (build/test commands, `doctor`, the rules for which records must exist) are deterministic. Overstated. |
| README: for "a bored dad… building something cool on the weekend" | The workflow, vocabulary, and records are engineer-grade. Nothing has been done to make Check output readable for a non-engineer. |
| Idea.md: "a vibe coder can go from init to a merged epic in one weekend" | Never measured. No external user data exists. |
| "Works with Claude Code, Cursor, Codex, Gemini, Aider" | Plausible because everything is plain files. Evidence in the repo is mostly Claude/Codex-style sessions. Other agents are untested here. |
| Token efficiency | The README says: "No goddamn idea." That is honest, and also unmeasured. |

### Where the effort has gone [Verified]

Of the 31 V2 Objectives, about **9 are Code Health**, about **8 are internal identity, renaming, or migration work** (hyphenated IDs, Release→Goal rename, V1 removal, Goal requirements), **5 are skills and guidance optimisation**, and **2 are parallel worktree planning**. **None are about onboarding, first-run, demos, or learning from outside users.** [Inference] Development is driven by what the owner runs into while using Savepoint on Savepoint.

---

## 3. Audience analysis

### A. Experienced developers using AI coding agents

- **Problems:** "Almost right" output. In the 2025 Stack Overflow survey, 66% named this as their top AI frustration, and trust in AI accuracy fell to roughly 33% ([VentureBeat](https://venturebeat.com/ai/stack-overflow-data-reveals-the-hidden-productivity-tax-of-almost-right-ai-code), [InfoWorld](https://infoworld.com/article/4031673/ai-use-among-software-developers-grows-but-trust-remains-an-issue-stack-overflow-survey.html)). They also lose context between sessions, and they misjudge their own speed: METR found experienced developers were 19% slower with AI while believing they were faster ([eWeek](https://eweek.com/news/news-ai-tools-slow-developer-productivity-study)). [Verified, secondary sources]
- **Workarounds:** Plan mode, their own AGENTS.md or CLAUDE.md, PR review bots (Claude Code Review, Codex review, Bugbot), skill packs such as Superpowers, and personal discipline.
- **Motivation:** Medium. They feel the pain, but they already have partial fixes and strong opinions.
- **Barriers:** Prescriptive process, context cost, and "I already have a system." The README currently tells them, jokingly, to go away.
- **Value of a visual UI:** Low. They live in the terminal and the editor.
- **Likely ongoing use:** Medium if Savepoint is light and optional in parts. Low if they have to adopt the full hierarchy.

### B. Technical hobbyists and indie builders (comfortable with the terminal, not professional-grade reviewers)

- **Problems:** Scope drift across weekends, losing the thread between sessions, agents that "finish" work that doesn't work, and projects that decay.
- **Workarounds:** Notes files, TODO.md, starting over, and Spec Kit or Superpowers.
- **Motivation:** Highest. Persistent state plus a second pair of eyes speaks directly to their pain. This is the owner's own profile.
- **Barriers:** Setup weight and vocabulary. They need to see the value in minutes, not after learning about Objectives and waivers.
- **Value of a visual UI:** Moderate for *reading* progress and Check results. Low for driving the workflow.
- **Likely ongoing use:** Medium to high, **if** the first loop (idea → task → check) completes in one sitting. [Hypothesis]

### C. Less experienced developers and vibe coders

- **Problems:** They can't judge generated code, they don't know when something is unsafe, and they ship exposed data. Coverage of RedAccess research reports thousands of vibe-coded apps leaking sensitive data ([RTInsights](https://www.rtinsights.com/can-vibe-coding-survive-the-new-era-of-security/), [Android Headlines](https://androidheadlines.com/2026/05/vibe-coding-security-risks-data-leaks-ai-apps.html)). [Verified, secondary sources]
- **Workarounds:** Hosted builders (Lovable, Replit, Bolt) with built-in previews and security scans, asking the AI "is this OK?", and hoping.
- **Motivation:** High pain, but **low awareness that process is the fix.** They want outcomes, not discipline.
- **Barriers:** Node/npx/Git, terminal-only owner decisions, engineering-grade Check reports, and the vocabulary. A GUI removes only the first of these.
- **Value of a visual UI:** High in principle, but only if the *content* changes too. Checks would need to answer in plain language: "this does what you asked / these 3 things don't / this part is risky".
- **Likely ongoing use:** Low today. [Inference] Serving this group is a different product: plain-language verification, likely a hosted or desktop experience. That conflicts with the current depth.

### Recommended primary audience

**Solo developers and technical indie builders who already run a terminal coding agent** (Claude Code, Codex CLI, Gemini CLI). This merges A's more pragmatic members with B.

**Trade-offs**
- *Gain:* the product's real complexity matches the people it serves. They already work in a terminal, so the CLI stops being a barrier. They can read and act on a Check. They are reachable through the same channels where Superpowers and Spec Kit spread.
- *Lose:* the larger, more emotionally appealing vibe-coder market, and the current playful "not for elite developers" voice.
- *Risk:* this is the audience best served by free, agent-native alternatives (§4). Savepoint must win on verification, not on planning.

Vibe coders should be a **later, separate bet**, revisited only if verification proves valuable for the primary audience *and* can be rewritten in plain language.

---

## 4. Competitive landscape

GitHub star counts were taken with the GitHub API on 2026-10-10 [Verified]. Stars measure attention, not use.

| Tool | Stars |
|---|---|
| [obra/superpowers](https://github.com/obra/superpowers) | ~297k |
| [github/spec-kit](https://github.com/github/spec-kit) | ~141k |
| [Fission-AI/OpenSpec](https://github.com/Fission-AI/OpenSpec) | ~71k |
| [BMAD-METHOD](https://github.com/bmad-code-org/BMAD-METHOD) | ~54k |
| [vibe-kanban](https://github.com/BloopAI/vibe-kanban) | ~28k |
| [claude-task-master](https://github.com/eyaltoledano/claude-task-master) | ~28k (last push April) |
| [Backlog.md](https://github.com/MrLesk/Backlog.md) | ~7k |
| [agent-os](https://github.com/buildermethods/agent-os) | ~5.5k |
| Savepoint | 1 |

### Agent-native workflow packs: Superpowers

- **Problem and audience:** Brainstorm → plan → execute with review checkpoints, TDD, and a reviewer subagent, delivered as a Claude Code plugin of Markdown skills ([guide](https://trevorlasn.com/blog/superpowers-claude-code-skills)).
- **Better than Savepoint:** Zero extra tool, installs inside the agent, works within seconds, and has huge distribution.
- **Savepoint's possible edge:** Durable, agent-agnostic records; a checker in a *separate session*; owner-only completion. [Inference]
- **Lesson:** **The most successful product in this space is pure Markdown skills with no binary and no UI.** This is the strongest evidence that neither the CLI nor the lack of a GUI is what limits reach.

### Spec-driven development kits: Spec Kit, OpenSpec, BMAD, Kiro, Tessl

- **Problem:** Write a spec first, then plan, then tasks, then implement. Spec Kit is a CLI added to your existing agent. Kiro is a paid IDE with review gates at each phase. OpenSpec tracks requirement changes in existing codebases ([comparison](https://ssojet.com/blog/best-spec-driven-development-tools); [Böckeler, martinfowler.com](https://www.martinfowler.com/articles/exploring-gen-ai/sdd-3-tools.html)).
- **Better than Savepoint:** Brand (GitHub, AWS), clearer one-line value, and large communities.
- **Known weaknesses:** Commentators report Markdown verbosity, agents not following the spec, and documentation overload ([Böckeler](https://www.martinfowler.com/articles/exploring-gen-ai/sdd-3-tools.html)). Savepoint shares these risks.
- **Savepoint's possible edge:** The closing half of the loop: independent Checks and REPLAN REQUIRED. Spec kits mostly cover the opening half. [Inference] Not verified per tool; Spec Kit has consistency-analysis commands that may partly overlap.
- **Standard expectation now:** "Plan before code" is table stakes, not a differentiator.

### Markdown task boards for agents: Backlog.md, Task Master

- **Backlog.md is Savepoint's closest neighbour:** local, Markdown in Git, a terminal board, **and already a local browser UI** (`backlog browser`), with a demo GIF, screenshots, and conference talks. It also follows "one task = one context window = one PR" ([README](https://cdn.jsdelivr.net/gh/MrLesk/Backlog.md@main/README.md)). [Verified]
- **Lesson:** A local web board did not make Backlog.md a breakout (~7k stars against Spec Kit's ~141k without a GUI). **Option B has already been built by a direct competitor**, so it would not differentiate Savepoint. [Inference]

### Parallel-agent orchestrators: Conductor, Vibe Kanban, the Codex app, Superset

- **Problem:** Run many agents in separate worktrees, then review the diffs ([Conductor](https://www.conductor.build/docs); [Vibe Kanban](https://aiidelist.com/ide/vibe-kanban)).
- **Better than Savepoint:** They actually run agents, show diffs and previews, and create PRs. Savepoint deliberately does none of this.
- **Warning sign:** The company behind Vibe Kanban **shut down in April 2026**, and the project is now community-maintained ([aiidelist](https://aiidelist.com/ide/vibe-kanban)). Agent vendors are building orchestration into their own products. [Verified, secondary source]

### Built-in review: Claude Code Review, Codex review, Cursor Bugbot

- **Problem:** Find bugs in a PR before a human reviews it. Claude Code Review uses several agents ([InfoQ](https://infoq.com/news/2026/04/claude-code-review)). Codex reviews on `@codex review` and reads review rules from AGENTS.md ([OpenAI docs](https://developers.openai.com/codex/use-cases/github-code-reviews)). Bugbot runs at large scale ([report](https://nevercodealone.de/vibe-coding/vibe-coding-modelle/cursor-bugbot-ki-debugging-vibe-coding)).
- **Better than Savepoint:** Automatic, zero setup, and in the PR where developers already look.
- **Savepoint's possible edge:** These tools ask "does this diff have bugs?". Savepoint asks **"does this work meet the outcome and acceptance criteria we agreed before building?"** It does this locally, works with any agent, and leaves a permanent record. [Inference] That is a real but narrow gap, and vendors could close it.

### "Software factory" approaches: StrongDM

- **Approach:** Specs and *holdout scenarios* that the building agent never sees decide whether a change passes; no human reviews the code ([Simon Willison](https://simonwillison.net/2026/Feb/7/software-factory/)).
- **Lesson:** Serious teams treat **independent verification as the hard problem**. They hide what the work will be measured against so agents can't game it. That validates direction C and suggests a stronger version of it: acceptance scenarios written before building and kept away from the builder. [Inference]

### Capabilities that are now standard expectations [Inference]

- Plan mode and planning documents.
- Instruction files following the AGENTS.md standard. It is now governed by the Linux Foundation's Agentic AI Foundation ([OpenAI](https://openai.com/index/agentic-ai-foundation/)).
- Automated PR review.
- Worktree-based parallel agents.
- Installing as an agent plugin or skill pack.
- A 30-second demo.

**Savepoint should not compete on any of these.**

---

## 5. Evaluation of product directions

| | Audience value | Differentiation | Complexity | Adoption friction | Maintenance | Main risk |
|---|---|---|---|---|---|---|
| **A: CLI-first discipline tool** (as is) | Medium for the primary audience | Low on planning, medium on checks | Already high | Medium to high (concept load) | High: 38k lines, migrations, skill parity | Superpowers and Spec Kit cover "discipline" for free |
| **B: Visual companion** | Medium for readers, low for drivers | **Low**: Backlog.md already has one | High: a second UI, local server, cross-platform | Lower to *look at*, unchanged to *set up* | High: two UIs to keep consistent with one runtime | Builds reach for an audience whose real barrier is elsewhere |
| **C: Independent verification** | **High**: hits the #1 pain ("almost right") | **Medium to high**: checks against requirements, not just bugs | Medium if built as a *subset* of today's product | Low *if* usable without the full hierarchy | Medium | Vendors add requirement-aware review |
| **D: Agent control layer** | Medium for power users | Low: Conductor, Codex app, agent teams | Very high | High | Very high | Vendors bundle it; Vibe Kanban's company already shut down |

**E: Proposed alternative, "Savepoint as a verification loop": C at the centre, a slimmed A as delivery.**
- Keep the files, `resume`, and owner authority.
- Make the **Check** the headline: an outcome and acceptance criteria written *before* building, an independent session that checks against them, and a short plain-language verdict at the top of each record.
- Offer a **light path**: one Objective with a few Tasks, no Goal ceremony up front, Code Health and lanes off by default.
- Consider shipping the skills as an **agent plugin** alongside the CLI (the Superpowers route) [Hypothesis].
- A tiny read-only output, such as a static HTML or Markdown summary of the latest Check, is a far cheaper test of "visual" than option B. [Hypothesis]

These directions overlap. E keeps everything in A that still pays its way and borrows B's only cheap piece. It rejects D.

---

## 6. Adoption and distribution findings

- **Is the CLI the barrier?** For the recommended audience, no [Inference]. They already use terminal agents. The most-adopted tools in this space (Superpowers, Spec Kit, OpenSpec) are CLI or skill packs with no GUI. Backlog.md has a GUI and less reach. For vibe coders, the CLI *is* a barrier, but it is not the only one, and a GUI alone would not fix their inability to judge a Check.
- **Is the value obvious?** No [Verified]. Savepoint currently has four different pitches:
  - README: "A little terminal app that helps you keep track of what you're building with AI"
  - Website: "Keep the project outside the chat"
  - The brief: "A little discipline for AI coding"
  - package.json: "cinematic Terminal UI … force you … to slow down"

  None of them names the one thing it does that others don't: independent checking.
- **Demos:** The website has no GIF, video, or real screenshot, only a text mock-up of the board [Verified via [getsavepoint.dev](https://www.getsavepoint.dev/)]. The README has a banner but no demo. Backlog.md leads with a GIF, a screenshot, and talks.
- **Trust signals:** The README FAQ jokes "Is it secure? Sounds like a question for someone who's actually read the code" and "Is it production ready? What's a production?" That is charming, but it undercuts a product whose core promise is *verification*. [Inference]
- **Usage signals:** npm shows about 3,300 downloads in the last month and about 8,600 in the year. But monthly downloads follow release activity, not discovery: 3,111 in the May launch month, 447 in August (1 commit), and 1,466/1,878 in Sep/Oct (309 commits, 45 published versions overall). [Verified numbers; the interpretation that mirrors, CI, and the owner's own `npx` runs make up much of this is an Inference.] GitHub has 1 star and zero external issues. **There is no evidence of retained outside users.**
- **What would make someone try it and keep using it** [Hypothesis]:
  - A 60-second demo of an agent claiming "done", a Savepoint Check finding what's missing, and the fix.
  - A first loop that finishes in one sitting.
  - Light context cost.
  - A visible "caught N issues" record over time.
  - Spreading the way Superpowers did: a plugin people can install inside the agent they already use.

---

## 7. Assumptions challenged

**Does Savepoint solve a painful enough problem?**
The *problem* is painful and well evidenced ("almost right" code, eroding trust, self-assessment that can't be trusted). Whether *Savepoint's* answer is worth its cost to anyone besides the owner is **unproven**. The internal 33/63 NEEDS WORK rate shows the Check catches things. It also shows the Check takes several rounds, which costs time and tokens.

**Is the workflow too prescriptive?**
For anyone but the owner, yes [Inference]. Evidence:
- the project's own Goal G-002 exists to cut duplicated guidance;
- the scaffold carries about 13k words;
- a mandatory Goal → Objective → Task hierarchy even for tiny projects;
- waiver records with actor and timestamp to *skip* an optional check.

The owner's own memory note says they prefer "simple V1-like designs", yet the product kept getting more complex.

**Could AGENTS.md or agent-native features make Savepoint unnecessary?**
- **Largely for planning and tracking**: plan mode, skills/plugins, Superpowers, Spec Kit.
- **Not yet for the verification loop with owner sign-off.** Even there, a skill pack plus a review subagent covers maybe 70% of it [Hypothesis]. Savepoint's lasting value rests on the remaining 30%: an independent session, checking against criteria written beforehand, a permanent record, and owner-only completion.

**Building what users want, or what we enjoy engineering?**
The evidence leans toward the second, and that should be said plainly:
- no external user has ever filed an issue;
- 9 Objectives on Code Health, an extensive internal measurement system, without one outside request;
- several Objectives on identity formats and renaming;
- zero Objectives on onboarding.

It is a sign that the feedback loop is closed. [Inference]

**Is the "terminal limits our audience" premise right?**
Mostly no. Audience fit, a fuzzy pitch, and missing proof come first. A GUI would make Savepoint *look* approachable while keeping the hard parts: the concepts and the dense Checks.

**Is local-first and deterministic a strategic cost?**
It rules out hosted onboarding and the vibe-coder market. That is acceptable for the primary audience, who value it. "Deterministic" should be claimed only for what really is deterministic: the gates and `resume`. Checks are judgment by a separate AI session.

---

## 8. Recommended direction and rationale

**Direction E: a narrower Savepoint centred on verification, for solo developers using terminal coding agents.**

1. **Reposition around one job.** Something like "Your agent says it's done. Savepoint checks." Keep the planning, but present it as the *input* to the check, not as the product. *Rationale:* this is the only capability with strong evidence (33/63) and a gap competitors haven't filled.
2. **Make a light path the default.** One Objective and a few Tasks, created without a Goal ceremony. Code Health, lanes, and Advanced Options become opt-in and hidden. Aim to cut the scaffold's agent guidance by more than half. *Rationale:* the concept load is the real adoption barrier for the primary audience.
3. **Put a plain-language verdict at the top of every Check.** Three to five lines: what was asked, what passed, what didn't, and what's risky. The matrix stays below. *Rationale:* this is useful now and is the prerequisite for any future less-technical audience.
4. **Fix the first five minutes.** The wrong first Next line on a new project, `doctor` refusing a directory, and a demo GIF or video.
5. **Get outside evidence before building any new capability.** Put 5 to 10 real users from the primary audience through the light path and watch them. Track whether they reach a second Check unprompted.
6. **Explore distribution as an agent plugin** (Spike 02 candidate). Package the skills so people can try them inside Claude Code or Codex before installing the binary. [Hypothesis]

This keeps local-first, file-based, deterministic gates, agent-agnostic records, and owner authority unchanged.

---

## 9. What NOT to build

- **No agent orchestration (D).** No launching agents, managing worktrees, scheduling, or monitoring. Vendors and well-funded tools own this, and one has already folded.
- **No full visual app (B) now.** No local web board, drag-and-drop, or desktop app. Backlog.md already has this, and it doesn't address the real barrier. Revisit only if real users ask to *read* Checks visually, and start with a static, read-only report.
- **No more Code Health scope.** Freeze it at its current shape until an outside user asks for it.
- **No hosted service, accounts, telemetry, or embedded AI.** These would break the principles that are Savepoint's reason to exist for this audience.
- **No PR bug-hunting reviewer.** Do not compete with Claude Code Review, Codex review, or Bugbot. Check against the agreed outcome, not the diff in general.
- **No vibe-coder product yet.** Don't build it until the plain-language verdict exists and has been tested.
- **No more internal-only identity, migration, or vocabulary work** unless it removes concepts for users.

---

## 10. Open questions and evidence gaps

- **No user research was possible in this spike.** Every audience claim beyond the owner is inference or secondary data. *This is the largest gap.*
- Is anyone using Savepoint besides the owner? npm numbers can't tell us. A short opt-in feedback request in the README, or direct outreach, could.
- How often does a Savepoint Check catch something that a vendor PR reviewer, or a simple "review against the plan" prompt, would miss? This could be tested by replaying past NEEDS WORK Checks against Claude Code Review or Codex review on the same diffs. **It is the most important experiment for direction C.**
- What does a full Check cost in time and tokens? Never measured (README: "No goddamn idea").
- Does the light path keep enough discipline to still catch problems, or does the value depend on the full ceremony?
- Plugin versus binary: how much of Savepoint works with skills alone, and what truly needs the runtime gates?
- Competitor details here come mostly from secondary sources and change fast. Spec Kit's analysis and checklist features and Superpowers' reviewer flow should be checked hands-on before claiming the verification gap.
- **Process note:** In Savepoint's own workflow, a spike is a Task with `complexity_tier: spike` under an Objective. This document was written as a standalone file, as the brief asked. No Objective, Task, or router record was created. If the owner wants it on record, `savepoint-design` can open a product-direction Objective that cites this file.

---

## 11. Go / Rework / Stop

**REWORK.**

- **Not Go:** the evidence does not support expanding (B or D). It also doesn't support continuing A unchanged, because there are no outside users and the guidance and concept load is high.
- **Not Stop:** the independent-check loop shows real, measured value on a real project, and the gap against competitors is real if narrow.
- **Rework means:**
  - narrow the audience;
  - lead with verification;
  - ship a light path;
  - fix first-run problems;
  - pause new capability work, including the rest of G-002 if it doesn't serve the light path, until outside users have tried it.

**Suggested kill criterion** [Hypothesis]: if 5 to 10 people from the primary audience try the light path and fewer than 3 run a second Check on their own initiative, stop expanding Savepoint. Keep it as a personal tool.

---

## The three decisions to make before Spike 02

1. **Who is it for?** Commit to *solo developers and indie builders already using terminal coding agents*, and retire the "bored dad / vibe coder" positioning. The alternative is to commit to vibe coders, which means building a different, plain-language product.
2. **What is the one job, and what gets frozen?** Accept "independent check against what was agreed" as the headline. Decide which of Code Health, lanes, the Goal hierarchy, and the remaining G-002 Objectives move off the default path or pause.
3. **What evidence earns the next round of building?** Set the outside-user test, the "Savepoint Check versus vendor reviewer" replay experiment, and the kill criterion *before* writing more code. Then let Spike 02 take on distribution (plugin versus binary versus both) on the basis of that evidence.

---

### Sources

Internal [Verified]: [`README.md`](README.md), [`.savepoint/Idea.md`](.savepoint/Idea.md), [`AGENTS.md`](AGENTS.md), [`templates/project-v2/`](templates/project-v2/), [`.savepoint/checks/`](.savepoint/checks/) (63 records), [`.savepoint/objectives/`](.savepoint/objectives/), [`.savepoint/releases/G-002-skills-optimisation/Release.md`](.savepoint/releases/G-002-skills-optimisation/Release.md), [`project-audit/consolidated-audit-report.md`](project-audit/consolidated-audit-report.md), a fresh `savepoint init` / `resume` / `doctor` / `board` run in a scratch repo, GitHub API and npm download API (2026-10-10).

External:
- Savepoint site: https://www.getsavepoint.dev/
- Stack Overflow 2025 survey coverage: https://venturebeat.com/ai/stack-overflow-data-reveals-the-hidden-productivity-tax-of-almost-right-ai-code · https://infoworld.com/article/4031673/ai-use-among-software-developers-grows-but-trust-remains-an-issue-stack-overflow-survey.html
- METR RCT coverage: https://eweek.com/news/news-ai-tools-slow-developer-productivity-study
- Spec-driven development: https://www.martinfowler.com/articles/exploring-gen-ai/sdd-3-tools.html · https://ssojet.com/blog/best-spec-driven-development-tools
- Superpowers: https://github.com/obra/superpowers · https://trevorlasn.com/blog/superpowers-claude-code-skills
- Backlog.md: https://github.com/MrLesk/Backlog.md
- Vibe Kanban: https://aiidelist.com/ide/vibe-kanban
- Conductor: https://www.conductor.build/docs
- Claude Code Review: https://infoq.com/news/2026/04/claude-code-review
- Codex GitHub review: https://developers.openai.com/codex/use-cases/github-code-reviews
- Cursor Bugbot: https://nevercodealone.de/vibe-coding/vibe-coding-modelle/cursor-bugbot-ki-debugging-vibe-coding
- StrongDM software factory: https://simonwillison.net/2026/Feb/7/software-factory/
- AGENTS.md / AAIF: https://openai.com/index/agentic-ai-foundation/
- Vibe-coding data exposure: https://www.rtinsights.com/can-vibe-coding-survive-the-new-era-of-security/ · https://androidheadlines.com/2026/05/vibe-coding-security-risks-data-leaks-ai-apps.html
