---
type: visual-identity
status: active
last_audited: never
---

# Visual Identity — Design System

> This file captures the design system for {{PROJECT_NAME}}: palette, typography, spacing, visual patterns, interaction principles, and the constraints that keep everything coherent. It is the single source of truth for visual decisions — consulted for any task that touches UI, rendering, or branding.

The sections below are written generically so they work for any project type (web, CLI, desktop, API). **Atari-Noir values shown in examples are the default shipped example theme** — replace them with your project's own design system.

## Palette

| Element       | Hex (example) | Role (example)                          |
| ------------- | ------------- | --------------------------------------- |
| Background    | `#000000`     | Base background                         |
| Surface       | `#000000`     | Cards and panels                        |
| Surface 2     | `#000000`     | Secondary panels / title bars           |
| Border        | `#1A1A1A`     | Quiet structural edges                  |
| Border Subtle | `#222222`     | Slightly stronger separators            |
| Primary Text  | `#F0E6DA`     | Warm off-white body text                |
| Accent 1      | `#FC6323`     | Primary CTA, active highlight           |
| Accent 2      | `#A4C639`     | Success, live systems                   |
| Accent 3      | `#B1A1DF`     | AI, reflection, secondary accent        |

**Color rules:**

- **One accent per section.** Use accent colors for labels, hover, active text, and highlights — never for large background fills.
- **Dark backgrounds.** Keep them dark so accents pop.
- **Visual encoding.** Color semantically encodes categories or states; reinforce with minimal text or glyphs.
- **Accessibility.** Ensure text/background contrast meets WCAG AA minimum.

## Typography

**Rules:**

- Headings: uppercase with deliberate letter-spacing where the medium supports it.
- Body: readable and restrained; abandon monospace for long text if readability suffers.
- Visual hierarchy: use size and weight before color. Magnitude and comparison values render better as proportional visuals (bars, circles) than raw numbers.

*Example (Atari-Noir theme): heading font is `Chakra Petch`, body/UI font is `Space Mono`, accent retro font is `Silkscreen` or `Press Start 2P` used rarely for score counters or deliberately extreme moments.*

## Spacing & Layout Rhythm

- Sections breathe. Default spacing feels generous.
- Cards and panels have enough internal padding to feel like panels, not chips.
- Whitespace creates hierarchy before borders or color do.

## Visual Patterns

- **Panels** — flat, structured with quiet borders. Depth comes from contrast and selective accent, not heavy shadows.
- **Search** — simple, single-line; integrated with the section accent color.
- **Data visualization** — prefer proportional visuals (scaled circles, gradient bars, markers) before raw values. Text reinforces; visuals carry comprehension.
- **Content blocks** — prose, code, and media integrated using standard typography, palette, and border logic.

## Interaction Principles ("The Playable Dashboard")

> **Content is the interface.** Avoid traditional UI patterns (dropdowns, accordions) if content can express the information directly.

- **Show, don't explain.** Visuals before text.
- **Motion** — authored, not generic. Ease like a system booting up. No bouncy toy motion.
- **Hover & focus** — internal light (underglow, surface tint, accent border), not thick outlines or loud transforms.
- **Expansion** — breathe open, don't snap. Inline reveals, shared-element transitions.

## Replication Brief

If recreating this look-and-feel, preserve:

- dark background with warm off-white text
- a small accent-color system (one per major section)
- uppercase headings with deliberate tracking
- quiet borders, dark surfaces, selective accent
- one strong interactive hero element
- copy that sounds human, competent, intentionally non-corporate
- **Audience baseline:** all content and visuals understandable by your least-expert intended user. Intuitive visual metaphors, no unexplained jargon, discoverable through use.

**Family resemblance, not exact duplication.** Preserve the underlying feel, hierarchy, and restraint.

*Example (Atari-Noir): `Chakra Petch` headings, `Space Mono` body, three-color accent system (orange/green/purple), uniform black backgrounds.*

## Flex & Constraints

- **Layout & components:** can adapt to content needs, provided palette discipline and tonal restraint remain.
- **Fonts:** supporting fonts can vary if heading/body contrast remains.
- **Hero:** interaction can vary, but requires one strong, ownable interactive element.
- **Visual-first:** show relationships visually before explaining them in text. Text reinforces; visuals carry comprehension.
- **Cognitive accessibility:** anchor complex concepts to concrete analogies your audience already understands.

---

## Appendix: TUI / Terminal Adaptation

The sections above apply broadly. When your project targets a terminal UI, the table below and the guardrails that follow describe how the design system translates.

### What survives in the terminal

| Design rule                          | Terminal feasibility | Adaptation                                                       |
| ------------------------------------ | -------------------- | ---------------------------------------------------------------- |
| Dark bg + warm off-white text        | ✓                    | 24-bit color with 256/16-color fallbacks                         |
| Small accent-color system            | ✓                    | Per status, per epic, per section                                |
| Custom heading / body fonts          | ✗                    | Terminal owns the font; README discloses                         |
| Uppercase tracked headings           | ⚠                    | Uppercase yes; letter-spacing no (fixed-width cells)             |
| Scanlines                            | ✗                    | Flicker/ugly in text — skip                                      |
| Glows / underglow                    | ⚠                    | Substitute with subtle bg tint on focused row + accent border    |
| Quiet borders, dark surfaces         | ✓                    | Box-drawing chars (`─ │ ┌ ┐`) in border-subtle gray              |
| Inline reveal cards                  | ✓                    | Application state-driven expand/collapse                         |
| Visual encoding before text          | ✓                    | Colored glyphs (`▣ ◇ ◆ ✓`) with text reinforcement               |
| Authored motion                      | ⚠                    | 200ms init sequence on launch acceptable; running animation skip |

### Terminal UI guardrails

- Keep UI state explicit and local to the smallest useful surface.
- Use project data files as the source of truth; do not invent hidden UI state.
- Keep layout compact and readable on narrow terminals first.
- Treat accidental line wrapping as a bug.
- Make selection, focus, and status changes obvious without relying on color alone.
- Make every keyboard action predictable and idempotent.
- Test branching input handling, navigation, state transitions, render output, and non-TTY fallbacks.

## When `savepoint init` ships

This file is the canonical default that `savepoint init` writes into a user's `.savepoint/visual-identity.md`. Users replace it with their own design system. The file's existence (not its contents) is what `savepoint` cares about.
