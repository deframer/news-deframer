// sentimentLabels.ts
//
// Why these heuristic thresholds are used
// --------------------------------------
// This file turns continuous VAD / BE5-like scores into discrete labels for
// rule-based interpretation. The thresholds below are not claimed to be
// canonical literature cutoffs. They are an engineering discretization built
// on top of well-established continuous affect models.
//
// Research basis for the model assumptions:
//
// - Russell (1977) provides evidence for three core affective dimensions:
//   pleasure-displeasure (valence), arousal, and dominance-submissiveness. ([sciencedirect.com](https://www.sciencedirect.com/science/article/pii/009265667790037X?utm_source=chatgpt.com))
//
// - Mohammad (2018) provides reliable human VAD ratings for 20,000+ English
//   words, supporting VAD as an operational lexical representation. ([aclanthology.org](https://aclanthology.org/P18-1017/?utm_source=chatgpt.com))
//
// - Buechel & Hahn (2018) show that dimensional representations such as VAD
//   can be mapped to discrete emotion formats and vice versa, which supports
//   the general rule-based linkage used here between VAD and BE5-like signals. ([aclanthology.org](https://aclanthology.org/C18-1245/?utm_source=chatgpt.com))
//
// Important:
// The affect model itself is research-backed.
// The exact thresholds and decision rules in this file are implementation
// heuristics, not canonical cutoffs from those papers.
//
// What is literature-backed vs. what is heuristic
// -----------------------------------------------
// Literature-backed:
// - Using Valence, Arousal, Dominance as continuous affect dimensions.
// - Using discrete emotion intensities as continuous scores.
// - Combining dimensional and discrete emotion representations.
//
// Heuristic / implementation choice:
// - The exact thresholds used here to create low / mid / high buckets.
// - The exact delta cutoffs for "clear", "mixed", "ambiguous" profiles.
// - The exact rule table mapping combinations to interpretation codes.
//
// Why the VAD thresholds are reasonable
// -------------------------------------
// VAD is on a 1..9 scale. The natural midpoint is 5.
// We use:
//
//   low  < 4.0
//   mid  4.0 .. < 6.0
//   high >= 6.0
//
// Rationale:
// - This creates a protected middle zone around the scale midpoint.
// - Values near the center are not over-interpreted as clearly positive /
//   negative, calm / activated, or weak / dominant.
// - It is a conservative discretization of a continuous scale.
// - In practice this is more stable for short texts and sentence-level input.
//
// In other words: 4 and 6 are not "published standard cutoffs"; they are a
// defensible symmetric partition around the midpoint of a 1..9 scale.
//
// Why the BE5 thresholds are reasonable
// -------------------------------------
// BE5-like emotion intensities here are treated as 1..5 continuous scores.
// We use:
//
//   very_low < 1.5
//   low      < 2.3
//   mid      < 3.2
//   high     >= 3.2
//
// Rationale:
// - 1.5 captures values that are close to absent / negligible.
// - 2.3 separates weak intensity from meaningfully present intensity.
// - 3.2 requires a value to be clearly above the midpoint before it is treated
//   as "high".
// - This avoids over-triggering strong labels such as joy/fear/anger for mildly
//   elevated values.
//
// Again: these are not canonical paper-defined cutoffs. They are practical
// working boundaries for stable rule-based interpretation.
//
// Why the delta thresholds are reasonable
// ---------------------------------------
// We compute the gap between the strongest and second-strongest discrete emotion:
//
//   delta >= 0.8  -> clear
//   delta >= 0.4  -> mixed
//   else          -> ambiguous
//
// Rationale:
// - On a 1..5 scale, a gap of 0.8 is already substantial.
// - A gap of 0.4 indicates some preference, but not a sharply dominant profile.
// - Smaller gaps are too close for a confident single-emotion reading.
//
// Why the rules use Dominance so heavily
// --------------------------------------
// Dominance operationalizes perceived control / power.
// This makes rules such as:
//
// - fear + low dominance -> threat / powerlessness / uncertainty
// - anger + high dominance -> hostile control / confrontation
// - joy + high dominance -> positive agency / confidence
//
// theoretically coherent within the PAD/VAD tradition.

export type EmotionVector = {
  valence: number;   // 1..9
  arousal: number;   // 1..9
  dominance: number; // 1..9
  joy: number;       // 1..5
  anger: number;     // 1..5
  sadness: number;   // 1..5
  fear: number;      // 1..5
  disgust: number;   // 1..5
};

type Level3 = "low" | "mid" | "high";
type BE5Level = "very_low" | "low" | "mid" | "high";
type EmotionName = "joy" | "anger" | "sadness" | "fear" | "disgust";

export type CodeMap = {
  core_state: string;
  primary_emotion: string;
  secondary_emotion: string;
  interpretation: string;
  tension_label: string;
  control_label: string;
  mood_label: string;
  clarity_label: string;
};

type RuleContext = {
  v: Level3;
  a: Level3;
  d: Level3;
  be5: Record<EmotionName, BE5Level>;
  clarity: "clear" | "mixed" | "ambiguous";
};

const VAD_THRESHOLDS = {
  lowMax: 4.0,
  highMin: 6.0,
} as const;

const BE5_THRESHOLDS = {
  veryLowMax: 1.5,
  lowMax: 2.3,
  highMin: 3.2,
} as const;

function classifyVAD(value: number): Level3 {
  if (value < VAD_THRESHOLDS.lowMax) return "low";
  if (value >= VAD_THRESHOLDS.highMin) return "high";
  return "mid";
}

function classifyBE5(value: number): BE5Level {
  if (value < BE5_THRESHOLDS.veryLowMax) return "very_low";
  if (value < BE5_THRESHOLDS.lowMax) return "low";
  if (value < BE5_THRESHOLDS.highMin) return "mid";
  return "high";
}

function sortEmotions(scores: EmotionVector): Array<{ name: EmotionName; value: number }> {
  const emotions: Array<{ name: EmotionName; value: number }> = [
    { name: "joy", value: scores.joy },
    { name: "anger", value: scores.anger },
    { name: "sadness", value: scores.sadness },
    { name: "fear", value: scores.fear },
    { name: "disgust", value: scores.disgust },
  ];
  emotions.sort((a, b) => b.value - a.value);
  return emotions;
}

function getClarity(delta: number): "clear" | "mixed" | "ambiguous" {
  if (delta >= 0.8) return "clear";
  if (delta >= 0.4) return "mixed";
  return "ambiguous";
}

function deriveCoreState(v: Level3, a: Level3, d: Level3): string {
  if (v === "low" && a === "high" && d === "low") return "threatened_activation";
  if (v === "low" && a === "high" && d === "high") return "hostile_activation";
  if (v === "low" && a === "mid" && d === "low") return "uneasy_low_control";
  if (v === "low" && a === "low" && d === "low") return "depleted_negative";
  if (v === "low" && a === "low" && d === "high") return "cold_negative_control";
  if (v === "high" && a === "high" && d === "high") return "empowered_positive_activation";
  if (v === "high" && a === "mid" && d === "high") return "confident_positive";
  if (v === "high" && a === "low" && d === "high") return "calm_positive_control";
  if (v === "high" && a === "low" && d === "low") return "soft_positive_passive";
  if (v === "mid" && a === "high" && d === "high") return "focused_activation";
  if (v === "mid" && a === "high" && d === "low") return "agitated_uncertain";
  if (v === "mid" && a === "low" && d === "mid") return "flat_neutral";
  return "mixed_state";
}

function derivePrimaryEmotion(
  emotionsSorted: Array<{ name: EmotionName; value: number }>,
  clarity: "clear" | "mixed" | "ambiguous"
): string {
  const top = emotionsSorted[0];
  if (top.value < BE5_THRESHOLDS.lowMax) return "none_salient";
  if (top.value >= BE5_THRESHOLDS.highMin) return top.name;
  if (clarity === "ambiguous") return "mixed";
  return top.name;
}

function deriveInterpretation(ctx: RuleContext): string {
  if (ctx.be5.fear === "high" && ctx.v === "low" && ctx.a === "high" && ctx.d === "low") {
    return "panic_threat_powerlessness";
  }
  if (ctx.be5.fear === "high" && ctx.v === "low" && ctx.d === "low") {
    return "fear_powerlessness_uncertainty";
  }
  if (ctx.be5.fear === "high" && ctx.a === "high" && (ctx.d === "mid" || ctx.d === "high")) {
    return "alarm_tension_worry";
  }
  if ((ctx.be5.fear === "mid" || ctx.be5.fear === "high") && ctx.be5.sadness === "high" && ctx.d === "low") {
    return "vulnerable_despair";
  }

  if (ctx.be5.anger === "high" && ctx.v === "low" && ctx.a === "high" && ctx.d === "high") {
    return "anger_aggressive_control";
  }
  if (ctx.be5.anger === "high" && ctx.d === "high") {
    return "hostile_control_confrontation";
  }
  if (ctx.be5.anger === "high" && ctx.d === "low") {
    return "frustration_irritability";
  }
  if ((ctx.be5.anger === "mid" || ctx.be5.anger === "high") && ctx.be5.disgust === "high") {
    return "contempt_harsh_rejection";
  }

  if (ctx.be5.sadness === "high" && ctx.v === "low" && ctx.a === "low") {
    return "quiet_grief_resignation";
  }
  if (ctx.be5.sadness === "high" && ctx.v === "low" && ctx.a === "mid") {
    return "sadness_loss_grief";
  }
  if (ctx.be5.sadness === "high" && ctx.a === "high") {
    return "agitated_grief_despair";
  }
  if (ctx.be5.sadness === "high" && ctx.be5.fear === "high" && ctx.d === "low") {
    return "helpless_despair_negative_vulnerability";
  }

  if (ctx.be5.joy === "high" && ctx.v === "high" && ctx.a === "high" && ctx.d === "high") {
    return "euphoria_triumph_enthusiasm";
  }
  if (ctx.be5.joy === "high" && ctx.v === "high" && ctx.a === "mid" && ctx.d === "high") {
    return "joy_confidence_positive_security";
  }
  if (ctx.be5.joy === "high" && ctx.v === "high" && ctx.a === "low") {
    return "contentment_warmth_calm_wellbeing";
  }
  if ((ctx.be5.joy === "mid" || ctx.be5.joy === "high") && ctx.d === "high" && ctx.v === "high") {
    return "positive_agency_selfconfidence";
  }

  if (ctx.be5.disgust === "high" && ctx.be5.anger === "high") {
    return "contempt_devaluing_rejection";
  }
  if (ctx.be5.disgust === "high" && ctx.be5.fear === "high") {
    return "defensive_aversion_threat";
  }
  if (ctx.be5.disgust === "high") {
    return "disgust_distancing_aversion";
  }

  if (ctx.clarity === "ambiguous" && ctx.a === "high" && ctx.v === "low") {
    return "ambiguous_tension";
  }
  if (ctx.clarity === "ambiguous" && ctx.v === "mid") {
    return "mixed_emotional_state";
  }
  if (ctx.v === "high") {
    return "positive_colored_mood";
  }
  if (ctx.v === "low") {
    return "negative_colored_mood";
  }
  return "neutral_to_mixed";
}

function deriveTensionLabel(v: Level3, a: Level3): string {
  if (a === "high" && v === "low") return "high_negative_tension";
  if (a === "high" && v === "high") return "high_positive_tension";
  if (a === "high" && v === "mid") return "strong_activation";
  if (a === "low") return "low_tension";
  return "medium_tension";
}

function deriveControlLabel(d: Level3): string {
  if (d === "high") return "high_control";
  if (d === "low") return "low_control";
  return "medium_control";
}

function deriveMoodLabel(v: Level3): string {
  if (v === "high") return "positive_tone";
  if (v === "low") return "negative_tone";
  return "ambivalent_tone";
}

function deriveClarityLabel(clarity: "clear" | "mixed" | "ambiguous"): string {
  if (clarity === "clear") return "clear_emotional_state";
  if (clarity === "mixed") return "mixed_emotional_state";
  return "ambiguous_emotional_state";
}

export function analyzeEmotionVectorToCodes(input: EmotionVector): CodeMap {
  const v = classifyVAD(input.valence);
  const a = classifyVAD(input.arousal);
  const d = classifyVAD(input.dominance);

  const be5: Record<EmotionName, BE5Level> = {
    joy: classifyBE5(input.joy),
    anger: classifyBE5(input.anger),
    sadness: classifyBE5(input.sadness),
    fear: classifyBE5(input.fear),
    disgust: classifyBE5(input.disgust),
  };

  const emotionsSorted = sortEmotions(input);
  const top = emotionsSorted[0];
  const second = emotionsSorted[1];
  const delta = top.value - second.value;
  const clarity = getClarity(delta);

  return {
    core_state: deriveCoreState(v, a, d),
    primary_emotion: derivePrimaryEmotion(emotionsSorted, clarity),
    secondary_emotion: second.name,
    interpretation: deriveInterpretation({ v, a, d, be5, clarity }),
    tension_label: deriveTensionLabel(v, a),
    control_label: deriveControlLabel(d),
    mood_label: deriveMoodLabel(v),
    clarity_label: deriveClarityLabel(clarity),
  };
}
