---
type: visual-identity
status: active
last_audited: never
---

# Visual Identity — Atari-Noir

> The brand identity for {{PROJECT_NAME}}. **Loaded conditionally** — only when working on TUI rendering, theme, or design-system tasks. Non-visual tasks must skip this file to honor the token budget.

## Vibe

A serious digital system that loves old arcade hardware. Dark, cinematic, playful. Crisp, not noisy. Technical, not cold.

**Avoid:** neon cyberpunk chaos, SaaS minimalism, fake terminal gimmicks.

## Palette

| Element       | Hex       | Role                          |
| ------------- | --------- | ----------------------------- |
| Background    | `#121212` | Main screen background        |
| Surface       | `#0D0D0D` | Cards and panels              |
| Surface 2     | `#0F0F0F` | Secondary panels / title bars |
| Border        | `#1A1A1A` | Quiet structural edges        |
| Border Subtle | `#222222` | Slightly stronger separators  |
| Primary Text  | `#F0E6DA` | Warm off-white terminal text  |
| Atari Orange  | `#FC6323` | Primary CTA, active highlight |
| NPP Green     | `#A4C639` | Success, live systems         |
| Vibe Purple   | `#B1A1DF` | AI, reflection                |

**Color rules:**

- **Intentional accents.** One accent color per major section/type. Use for labels, hover, glows, active text — never giant background fills.
- **Dark backgrounds.** Keep them dark so accents pop.
- **Visual encoding.** Color semantically encodes categories or states; reinforce with minimal text.

## Typography

- **System heading font:** `Chakra Petch` (web only — terminal cannot control fonts)
- **System body/UI font:** `Space Mono` (web only — terminal cannot control fonts)
- **Accent retro fonts:** `Silkscreen` or `Press Start 2P`, used rarely for score counters or deliberately extreme moments

**Rules:**

- Headings uppercase, with deliberate letter-spacing where the medium allows.
- Body readable and restrained; abandon monospace for long text if readability suffers.
- Render magnitude/comparison numbers as proportional visuals (bars, circles), not raw values.

## Spacing rhythm

- Sections breathe. Default spacing feels generous.
- Cards have enough internal padding to feel like panels, not chips.
- Whitespace creates hierarchy before borders or color do.

## Signature UI patterns

- **Scanlines** — low-opacity CRT atmosphere; never obscure readability.
- **Glows** — radiate from perimeter ("lit from within"). Transparent section accent colors; no neon spam.
- **Panels** — flat, dark, structured with quiet borders. Depth from contrast/glow, not heavy shadows.
- **Search** — simple, single-line; integrated with section accent color.
- **MDX prose** — integrated using standard typography, palette, border logic.

## Interaction principles ("The Playable Dashboard")

**Content is the interface.** Avoid traditional UI patterns (dropdowns, accordions) if content can express the information directly.

- **Show, don't explain.** Visuals (scaled circles, markers) before text.
- **Motion** — authored, not generic. Ease like a system booting up. No bouncy toy motion.
- **Hover & focus** — internal light (underglow/surface tint), not thick outlines or loud transforms.
- **Expansion** — breathe open, don't snap. Inline reveals, shared-element transitions.

## Replication brief

If recreating this look-and-feel, preserve:

- dark charcoal background with warm off-white text
- three-color accent system (one per major section)
- `Chakra Petch` for headings, `Space Mono` for body/UI (where typography can be controlled)
- uppercase, tracked headings
- quiet borders, dark surfaces, selective glow
- one strong interactive hero element
- copy that sounds human, competent, intentionally non-corporate
- **Young Explorer Baseline:** all content and visuals understandable by a 7-year-old. Intuitive visual metaphors, no unexplained jargon, discoverable through play.

**Family resemblance, not exact duplication.** Preserve the underlying feel, hierarchy, restraint.

## Flex & constraints

- **Layout & components:** can adapt to content needs, provided palette discipline and tonal restraint remain.
- **Fonts:** supporting fonts can vary if heading/body contrast remains.
- **Hero:** interaction can vary, but requires one strong, ownable interactive element.
- **Visual-first:** show relationships visually before explaining them in text (gradient bar > numbers). Text reinforces; visuals carry comprehension.
- **Cognitive accessibility:** anchor complex concepts to physical/visual analogies a 7-year-old can grasp.

## What survives in the terminal

| Web rule                            | Terminal feasibility | Adaptation                                                       |
| ----------------------------------- | -------------------- | ---------------------------------------------------------------- |
| Dark bg + warm off-white text       | ✓                    | 24-bit color with 256/16-color fallbacks                         |
| 3-color accent system               | ✓                    | per status, per epic, per section                                |
| `Chakra Petch` / `Space Mono` fonts | ✗                    | terminal owns the font; README discloses                         |
| Uppercase tracked headings          | ⚠                    | uppercase yes; letter-spacing no (fixed-width cells)             |
| Scanlines                           | ✗                    | flicker/ugly in text — skip                                      |
| Glows / underglow                   | ⚠                    | substitute with subtle bg tint on focused row + accent border    |
| Quiet borders, dark surfaces        | ✓                    | box-drawing chars (`─ │ ┌ ┐`) in border-subtle gray              |
| Inline reveal cards                 | ✓                    | Ink state-driven expand/collapse                                 |
| Visual encoding before text         | ✓                    | colored glyphs (`▣ ◇ ◆ ✓`) with text reinforcement               |
| "System booting up" motion          | ⚠                    | 200ms init sequence on launch acceptable; running animation skip |

## When `savepoint init` ships

This file is the canonical default that `savepoint init` writes into a user's `.savepoint/visual-identity.md`. Users replace it with their own design system. The file's existence (not its contents) is what `savepoint` cares about.
