## Purpose

`sentimentLabels.ts` converts a continuous affect vector into a compact set of rule-based sentiment labels.

Input:

```json
{
  "valence": 7.72,
  "arousal": 4.77,
  "dominance": 6.92,
  "joy": 3.96,
  "anger": 1.24,
  "sadness": 1.29,
  "fear": 1.37,
  "disgust": 1.15
}
````

Output:

```json
{
  "core_state": "confident_positive",
  "primary_emotion": "joy",
  "secondary_emotion": "fear",
  "interpretation": "joy_confidence_positive_security",
  "tension_label": "medium_tension",
  "control_label": "high_control",
  "mood_label": "positive_tone",
  "clarity_label": "clear_emotional_state"
}
```

The file does **not** try to be a scientific emotion model by itself. It is a **heuristic interpretation layer** built on top of research-backed continuous affect models. The underlying dimensions are research-based; the thresholds and rules are engineering choices designed to produce stable sentence-level labels.

---

## Research basis

### 1. VAD / PAD as a dimensional affect model

The use of **Valence**, **Arousal**, and **Dominance** comes from the dimensional emotion tradition. Russell (1977) provides evidence for three core affective dimensions: pleasure-displeasure, arousal, and dominance-submissiveness.

### 2. VAD as an operational lexical representation

Mohammad (2018) provides reliable human VAD ratings for more than 20,000 English words. This supports using VAD as a structured numeric affect representation.

### 3. Mapping between dimensional and discrete emotion formats

Buechel and Hahn (2018) show that dimensional emotion representations such as VAD can be mapped to discrete emotion formats and back. This supports the general architecture used here: combine VAD-like dimensions with BE5-like discrete emotion scores and derive rule-based labels from both.

---

## What is research-backed and what is heuristic

### Research-backed

* using **Valence**, **Arousal**, and **Dominance** as continuous affect dimensions
* using continuous discrete-emotion intensities as an input representation
* linking dimensional and discrete emotion spaces

### Heuristic / implementation choice

* the exact thresholds used to create `low`, `mid`, `high`
* the exact cutoffs for `clear`, `mixed`, `ambiguous`
* the exact interpretation table
* the exact output vector design

This means the model assumptions are literature-supported, but the final labeling system is a practical operationalization.

---

## Why these threshold values are used

### VAD thresholds

The VAD dimensions are treated as a **1..9** scale.

The implementation uses:

* `low < 4.0`
* `mid >= 4.0 && < 6.0`
* `high >= 6.0`

#### Why this is reasonable

The natural midpoint of a 1..9 scale is **5**. Using `4` and `6` creates a protected middle zone around the center instead of forcing values near the midpoint into strong categories. This is a conservative discretization of a continuous scale. It avoids over-interpreting values that are only slightly above or below neutral. The literature supports the use of continuous VAD dimensions, but not these exact cutoffs; these are practical thresholds derived from the scale structure itself.

#### Practical rationale

* values near 5 stay in a neutral-ish buffer
* only clearly shifted values become `high` or `low`
* this is more stable for sentence-level analysis

---

### BE5-style thresholds

The discrete emotion scores are treated as a **1..5** intensity scale.

The implementation uses:

* `very_low < 1.5`
* `low < 2.3`
* `mid < 3.2`
* `high >= 3.2`

#### Why this is reasonable

These are not canonical paper-defined cutoffs. They are working boundaries.

The logic is:

* `1.5` is close to near-absence
* `2.3` separates weak presence from more meaningful presence
* `3.2` requires a score to be clearly above the midpoint before it is treated as high

This keeps the rule system from over-triggering strong labels such as joy, fear, or anger for only mildly elevated values. The literature supports continuous discrete emotion representations and mappings, but these exact numeric boundaries are an engineering decision.

---

### Clarity thresholds

The implementation compares the top discrete emotion score with the second-highest score.

* `delta >= 0.8` → `clear`
* `delta >= 0.4` → `mixed`
* otherwise → `ambiguous`

#### Why this is reasonable

On a 1..5 scale, a gap of `0.8` is already substantial. A gap of `0.4` indicates some preference but not a sharply dominant profile. Smaller gaps are too close for a confident single-emotion reading. These numbers are heuristic, but they behave well on compact emotion scales.

---

## Why Dominance matters so much

Dominance operationalizes perceived control, agency, or power. That makes it especially important for interpretation rules such as:

* `fear + low dominance` → threat / powerlessness / uncertainty
* `anger + high dominance` → hostile control / confrontation
* `joy + high dominance` → positive agency / confidence

This is consistent with the PAD/VAD tradition, where dominance is not just another number but a psychologically meaningful control dimension.

---

## Why the output uses exactly these 8 labels

The result vector is intentionally split into eight fields because a single emotion label is too coarse for sentence-level affect analysis. The goal is to separate different analytic layers that are related, but not identical.

### 1. `core_state`

Example: `confident_positive`

**Purpose**

* captures the global VAD configuration
* expresses the broad affective state derived from valence, arousal, and dominance alone
* gives a state-space reading before discrete emotions are applied

**Why it exists**

* VAD describes overall affective position better than any single discrete emotion
* two sentences can both express joy but differ strongly in activation or control
* `core_state` preserves that distinction

**In the example**

* valence is high
* arousal is mid
* dominance is high

So the global state is `confident_positive`.

---

### 2. `primary_emotion`

Example: `joy`

**Purpose**

* identifies the strongest discrete emotion channel

**Why it exists**

* VAD alone does not tell which named emotion is dominant
* the strongest BE5-like score provides the most salient discrete emotion

**In the example**

* joy is the highest of joy/anger/sadness/fear/disgust

So the primary discrete emotion is `joy`.

---

### 3. `secondary_emotion`

Example: `fear`

**Purpose**

* preserves the second strongest discrete emotion

**Why it exists**

* emotional readings are often not purely single-label
* the runner-up signal is useful for mixed, ambivalent, or contrastive states
* this helps explain borderline or layered interpretations

**In the example**

* fear is the second-highest discrete score

So the secondary emotion is `fear`.

---

### 4. `interpretation`

Example: `joy_confidence_positive_security`

**Purpose**

* produces the main semantic reading from the combination of VAD state and discrete emotion profile
* this is the main human-facing interpretation code

**Why it exists**

* `primary_emotion = joy` is too thin
* the interpretation combines:

  * joy as dominant emotion
  * high valence
  * medium arousal
  * high dominance
* this yields not just joy, but a more specific reading: confidence, positive security, stable positive affect

**In the example**

* joy high
* valence high
* arousal mid
* dominance high

So the interpretation is `joy_confidence_positive_security`.

---

### 5. `tension_label`

Example: `medium_tension`

**Purpose**

* separates activation/tension from general positivity or negativity

**Why it exists**

* a sentence can be positive but calm, or positive and highly energized
* tension is mainly a simplified reading of arousal combined with valence

**In the example**

* arousal is mid

So the tension label is `medium_tension`.

---

### 6. `control_label`

Example: `high_control`

**Purpose**

* makes the dominance dimension explicit

**Why it exists**

* dominance is too important to leave buried inside `core_state`
* it often separates fear-like helplessness from assertive or confident states
* this dimension is especially useful in narrative analysis because it tracks agency, power, vulnerability, and control

**In the example**

* dominance is high

So the control label is `high_control`.

---

### 7. `mood_label`

Example: `positive_tone`

**Purpose**

* gives a simplified polarity reading

**Why it exists**

* sometimes a downstream system only needs broad tone
* this is a compressed reading of valence alone
* it is intentionally simpler than `core_state`

**In the example**

* valence is high

So the broad tone is `positive_tone`.

---

### 8. `clarity_label`

Example: `clear_emotional_state`

**Purpose**

* indicates how sharply one discrete emotion dominates the others

**Why it exists**

* not every sentence has a clean emotional profile
* this label tells whether the result is clear, mixed, or ambiguous
* it is derived from the gap between the strongest and second strongest discrete emotion scores

**In the example**

* joy is clearly above fear and the other channels

So the clarity label is `clear_emotional_state`.

---

## Why these 8 are enough

These eight fields cover distinct analytic questions:

* What is the global affective state? → `core_state`
* Which named emotion is strongest? → `primary_emotion`
* What is the secondary pull? → `secondary_emotion`
* What is the combined semantic reading? → `interpretation`
* How activated or tense is it? → `tension_label`
* How much control or agency is present? → `control_label`
* Is the tone broadly positive or negative? → `mood_label`
* How sharp or mixed is the profile? → `clarity_label`

This is why the vector is not collapsed into one label. A single label would destroy important distinctions such as:

* joy with calm vs joy with excitement
* fear with helplessness vs fear with partial control
* positive tone with mixed secondary emotion
* clear emotion vs ambiguous profile

---

## Why the example result is correct

Input:

```json
{
  "valence": 7.72,
  "arousal": 4.77,
  "dominance": 6.92,
  "joy": 3.96,
  "anger": 1.24,
  "sadness": 1.29,
  "fear": 1.37,
  "disgust": 1.15
}
```

Step by step:

* `valence = 7.72` → `high`
* `arousal = 4.77` → `mid`
* `dominance = 6.92` → `high`

Therefore:

* `core_state = confident_positive`

Then the discrete emotions:

* `joy = 3.96` is highest
* `fear = 1.37` is second-highest

Therefore:

* `primary_emotion = joy`
* `secondary_emotion = fear`

Then the fused interpretation:

* joy is high
* valence is high
* arousal is mid
* dominance is high

Therefore:

* `interpretation = joy_confidence_positive_security`

The supporting summary labels follow from the same profile:

* arousal is mid → `tension_label = medium_tension`
* dominance is high → `control_label = high_control`
* valence is high → `mood_label = positive_tone`
* the top discrete emotion is clearly above the second → `clarity_label = clear_emotional_state`

---

## Design principle

The result vector is deliberately multi-axis:

* VAD contributes state, tension, control, and tone
* discrete emotions contribute salience and mixture
* the interpretation field fuses both layers

That is why these eight labels exist.

---

## Short summary

`sentimentLabels.ts` is a rule-based interpretation layer over continuous affect scores.

It is based on:

* dimensional affect theory for VAD
* mapping between dimensional and discrete emotion spaces

It uses:

* conservative thresholds for stable discretization
* a rule table for semantic interpretation
* an 8-field output vector to preserve distinctions that a single label would lose

