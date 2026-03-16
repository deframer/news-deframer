type EmotionVector = {
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

type CoreStateCode =
  | "threatened_activation"
  | "hostile_activation"
  | "uneasy_low_control"
  | "depleted_negative"
  | "cold_negative_control"
  | "empowered_positive_activation"
  | "confident_positive"
  | "calm_positive_control"
  | "soft_positive_passive"
  | "focused_activation"
  | "agitated_uncertain"
  | "flat_neutral"
  | "mixed_state";

type PrimaryEmotionCode =
  | EmotionName
  | "mixed"
  | "none_salient";

type InterpretationCode =
  | "panic_threat_powerlessness"
  | "fear_powerlessness_uncertainty"
  | "alarm_tension_worry"
  | "vulnerable_despair"
  | "anger_aggressive_control"
  | "hostile_control_confrontation"
  | "frustration_irritability"
  | "contempt_harsh_rejection"
  | "quiet_grief_resignation"
  | "sadness_loss_grief"
  | "agitated_grief_despair"
  | "helpless_despair_negative_vulnerability"
  | "euphoria_triumph_enthusiasm"
  | "joy_confidence_positive_security"
  | "contentment_warmth_calm_wellbeing"
  | "positive_agency_selfconfidence"
  | "contempt_devaluing_rejection"
  | "defensive_aversion_threat"
  | "disgust_distancing_aversion"
  | "hope_fear_ambivalence"
  | "bitter_hurt_negativity"
  | "alarmed_aggression_defensive_tension"
  | "bittersweet_ambivalence"
  | "ambiguous_tension"
  | "mixed_emotional_state"
  | "positive_colored_mood"
  | "negative_colored_mood"
  | "neutral_to_mixed";

type TensionCode =
  | "high_negative_tension"
  | "high_positive_tension"
  | "strong_activation"
  | "low_tension"
  | "medium_tension";

type ControlCode =
  | "high_control"
  | "low_control"
  | "medium_control";

type MoodCode =
  | "positive_tone"
  | "negative_tone"
  | "ambivalent_tone";

type ClarityCode =
  | "clear_emotional_state"
  | "mixed_emotional_state"
  | "ambiguous_emotional_state";

type AnalysisOutput = {
  core_state: CoreStateCode;
  primary_emotion: PrimaryEmotionCode;
  secondary_emotion: EmotionName;
  interpretation: InterpretationCode;
  tension_label: TensionCode;
  control_label: ControlCode;
  mood_label: MoodCode;
  clarity_label: ClarityCode;
};

type RuleContext = {
  v: Level3;
  a: Level3;
  d: Level3;
  be5: Record<EmotionName, BE5Level>;
  primaryEmotion: PrimaryEmotionCode;
  secondaryEmotion: EmotionName;
  clarity: "clear" | "mixed" | "ambiguous";
  scores: EmotionVector;
  emotionsSorted: Array<{ name: EmotionName; value: number }>;
};

type DecisionRule = {
  id: string;
  when: (ctx: RuleContext) => boolean;
  then: InterpretationCode;
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

export const LABELS_DE = {
  core_state: {
    threatened_activation: "bedrohlich, angespannt, kontrollarm",
    hostile_activation: "negativ, hochaktiviert, durchsetzungsstark",
    uneasy_low_control: "unangenehm, unsicher, kontrollarm",
    depleted_negative: "leer, resignativ, energiereduziert",
    cold_negative_control: "negativ, gedämpft, kontrolliert",
    empowered_positive_activation: "euphorisch, kraftvoll, kontrolliert",
    confident_positive: "positiv, stabil, selbstsicher",
    calm_positive_control: "gelassen, zufrieden, souverän",
    soft_positive_passive: "sanft, warm, passiv positiv",
    focused_activation: "angespannt fokussiert, handlungsbereit",
    agitated_uncertain: "nervös, instabil, reaktiv",
    flat_neutral: "flach, neutral, wenig aktiviert",
    mixed_state: "gemischtes Profil",
  } satisfies Record<CoreStateCode, string>,

  primary_emotion: {
    joy: "Freude",
    anger: "Ärger",
    sadness: "Traurigkeit",
    fear: "Angst",
    disgust: "Ekel",
    mixed: "gemischt",
    none_salient: "keine dominante diskrete Emotion",
  } satisfies Record<PrimaryEmotionCode, string>,

  interpretation: {
    panic_threat_powerlessness: "Panik, Bedrohung, Ausgeliefertsein",
    fear_powerlessness_uncertainty: "Angst, Ohnmacht, Unsicherheit",
    alarm_tension_worry: "Alarm, Anspannung, Sorge",
    vulnerable_despair: "verletzliche Verzweiflung",
    anger_aggressive_control: "Wut, aggressive Kontrolle, Kampfbereitschaft",
    hostile_control_confrontation: "Durchsetzung, Konfrontation, feindselige Kontrolle",
    frustration_irritability: "Frustration, Gereiztheit, reaktive Negativität",
    contempt_harsh_rejection: "Verachtung, harte Ablehnung",
    quiet_grief_resignation: "stille Trauer, Resignation, Melancholie",
    sadness_loss_grief: "Traurigkeit, Kummer, Verlust",
    agitated_grief_despair: "Verzweiflung, aufgewühlte Trauer",
    helpless_despair_negative_vulnerability: "Hilflosigkeit, Verzweiflung, verletzliche Negativität",
    euphoria_triumph_enthusiasm: "Euphorie, Triumph, Begeisterung",
    joy_confidence_positive_security: "Freude, Zuversicht, positive Sicherheit",
    contentment_warmth_calm_wellbeing: "Zufriedenheit, Wärme, ruhiges Wohlgefühl",
    positive_agency_selfconfidence: "Selbstsicherheit, positive Handlungsfähigkeit",
    contempt_devaluing_rejection: "Verachtung, abwertende Ablehnung",
    defensive_aversion_threat: "Abwehr, Abstoßung, aversive Bedrohung",
    disgust_distancing_aversion: "Ekel, Distanzierung, Abstoßung",
    hope_fear_ambivalence: "ambivalente Spannung zwischen Hoffnung und Angst",
    bitter_hurt_negativity: "bitterer Schmerz, kränkende Negativität",
    alarmed_aggression_defensive_tension: "alarmierte Aggression, defensive Kampfspannung",
    bittersweet_ambivalence: "bittersüße Ambivalenz",
    ambiguous_tension: "uneindeutige Spannung",
    mixed_emotional_state: "gemischte emotionale Lage",
    positive_colored_mood: "positiv gefärbte Stimmung",
    negative_colored_mood: "negativ gefärbte Stimmung",
    neutral_to_mixed: "neutral bis gemischt",
  } satisfies Record<InterpretationCode, string>,

  tension_label: {
    high_negative_tension: "hohe negative Spannung",
    high_positive_tension: "hohe positive Spannung",
    strong_activation: "starke Aktivierung",
    low_tension: "niedrige Spannung",
    medium_tension: "mittlere Spannung",
  } satisfies Record<TensionCode, string>,

  control_label: {
    high_control: "hohes Kontrollgefühl",
    low_control: "niedriges Kontrollgefühl",
    medium_control: "mittleres Kontrollgefühl",
  } satisfies Record<ControlCode, string>,

  mood_label: {
    positive_tone: "positiver Grundton",
    negative_tone: "negativer Grundton",
    ambivalent_tone: "ambivalenter Grundton",
  } satisfies Record<MoodCode, string>,

  clarity_label: {
    clear_emotional_state: "klare Emotionslage",
    mixed_emotional_state: "gemischte Emotionslage",
    ambiguous_emotional_state: "uneindeutige Emotionslage",
  } satisfies Record<ClarityCode, string>,
} as const;

export type AnalysisOutputTranslated = AnalysisOutput & {
  core_state_de: string;
  primary_emotion_de: string;
  secondary_emotion_de: string;
  interpretation_de: string;
  tension_label_de: string;
  control_label_de: string;
  mood_label_de: string;
  clarity_label_de: string;
};

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

function deriveCoreState(v: Level3, a: Level3, d: Level3): CoreStateCode {
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
): PrimaryEmotionCode {
  const top = emotionsSorted[0];

  if (top.value < BE5_THRESHOLDS.lowMax) return "none_salient";
  if (top.value >= BE5_THRESHOLDS.highMin) return top.name;
  if (clarity === "ambiguous") return "mixed";
  return top.name;
}

const INTERPRETATION_RULES: DecisionRule[] = [
  {
    id: "H1",
    when: (c) =>
      c.be5.fear === "high" && c.v === "low" && c.a === "high" && c.d === "low",
    then: "panic_threat_powerlessness",
  },
  {
    id: "H2",
    when: (c) =>
      c.be5.fear === "high" && c.v === "low" && c.d === "low",
    then: "fear_powerlessness_uncertainty",
  },
  {
    id: "H3",
    when: (c) =>
      c.be5.fear === "high" && c.a === "high" && (c.d === "mid" || c.d === "high"),
    then: "alarm_tension_worry",
  },
  {
    id: "H4",
    when: (c) =>
      (c.be5.fear === "mid" || c.be5.fear === "high") &&
      c.be5.sadness === "high" &&
      c.d === "low",
    then: "vulnerable_despair",
  },
  {
    id: "H5",
    when: (c) =>
      c.be5.anger === "high" && c.v === "low" && c.a === "high" && c.d === "high",
    then: "anger_aggressive_control",
  },
  {
    id: "H6",
    when: (c) =>
      c.be5.anger === "high" && c.d === "high",
    then: "hostile_control_confrontation",
  },
  {
    id: "H7",
    when: (c) =>
      c.be5.anger === "high" && c.d === "low",
    then: "frustration_irritability",
  },
  {
    id: "H8",
    when: (c) =>
      (c.be5.anger === "mid" || c.be5.anger === "high") &&
      c.be5.disgust === "high",
    then: "contempt_harsh_rejection",
  },
  {
    id: "H9",
    when: (c) =>
      c.be5.sadness === "high" && c.v === "low" && c.a === "low",
    then: "quiet_grief_resignation",
  },
  {
    id: "H10",
    when: (c) =>
      c.be5.sadness === "high" && c.v === "low" && c.a === "mid",
    then: "sadness_loss_grief",
  },
  {
    id: "H11",
    when: (c) =>
      c.be5.sadness === "high" && c.a === "high",
    then: "agitated_grief_despair",
  },
  {
    id: "H12",
    when: (c) =>
      c.be5.sadness === "high" && c.be5.fear === "high" && c.d === "low",
    then: "helpless_despair_negative_vulnerability",
  },
  {
    id: "H13",
    when: (c) =>
      c.be5.joy === "high" && c.v === "high" && c.a === "high" && c.d === "high",
    then: "euphoria_triumph_enthusiasm",
  },
  {
    id: "H14",
    when: (c) =>
      c.be5.joy === "high" && c.v === "high" && c.a === "mid" && c.d === "high",
    then: "joy_confidence_positive_security",
  },
  {
    id: "H15",
    when: (c) =>
      c.be5.joy === "high" && c.v === "high" && c.a === "low",
    then: "contentment_warmth_calm_wellbeing",
  },
  {
    id: "H16",
    when: (c) =>
      (c.be5.joy === "mid" || c.be5.joy === "high") &&
      c.d === "high" &&
      c.v === "high",
    then: "positive_agency_selfconfidence",
  },
  {
    id: "H17",
    when: (c) =>
      c.be5.disgust === "high" && c.be5.anger === "high",
    then: "contempt_devaluing_rejection",
  },
  {
    id: "H18",
    when: (c) =>
      c.be5.disgust === "high" && c.be5.fear === "high",
    then: "defensive_aversion_threat",
  },
  {
    id: "H19",
    when: (c) =>
      c.be5.disgust === "high",
    then: "disgust_distancing_aversion",
  },
  {
    id: "H20",
    when: (c) =>
      c.be5.joy === "high" && c.be5.fear === "high",
    then: "hope_fear_ambivalence",
  },
  {
    id: "H21",
    when: (c) =>
      c.be5.sadness === "high" && c.be5.anger === "high",
    then: "bitter_hurt_negativity",
  },
  {
    id: "H22",
    when: (c) =>
      c.be5.fear === "high" && c.be5.anger === "high",
    then: "alarmed_aggression_defensive_tension",
  },
  {
    id: "H23",
    when: (c) =>
      c.be5.joy === "high" && c.be5.sadness === "high",
    then: "bittersweet_ambivalence",
  },
  {
    id: "H24",
    when: (c) =>
      c.clarity === "ambiguous" && c.a === "high" && c.v === "low",
    then: "ambiguous_tension",
  },
  {
    id: "H25",
    when: (c) =>
      c.clarity === "ambiguous" && c.v === "mid",
    then: "mixed_emotional_state",
  },
  {
    id: "H26",
    when: (c) => c.v === "high",
    then: "positive_colored_mood",
  },
  {
    id: "H27",
    when: (c) => c.v === "low",
    then: "negative_colored_mood",
  },
  {
    id: "H28",
    when: () => true,
    then: "neutral_to_mixed",
  },
];

function deriveInterpretation(ctx: RuleContext): InterpretationCode {
  const matched = INTERPRETATION_RULES.find((rule) => rule.when(ctx));
  return matched ? matched.then : "neutral_to_mixed";
}

function deriveTensionLabel(v: Level3, a: Level3): TensionCode {
  if (a === "high" && v === "low") return "high_negative_tension";
  if (a === "high" && v === "high") return "high_positive_tension";
  if (a === "high" && v === "mid") return "strong_activation";
  if (a === "low") return "low_tension";
  return "medium_tension";
}

function deriveControlLabel(d: Level3): ControlCode {
  if (d === "high") return "high_control";
  if (d === "low") return "low_control";
  return "medium_control";
}

function deriveMoodLabel(v: Level3): MoodCode {
  if (v === "high") return "positive_tone";
  if (v === "low") return "negative_tone";
  return "ambivalent_tone";
}

function deriveClarityLabel(clarity: "clear" | "mixed" | "ambiguous"): ClarityCode {
  if (clarity === "clear") return "clear_emotional_state";
  if (clarity === "mixed") return "mixed_emotional_state";
  return "ambiguous_emotional_state";
}

export function analyzeEmotionVector(input: EmotionVector): AnalysisOutput {
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
  const primaryEmotion = derivePrimaryEmotion(emotionsSorted, clarity);
  const secondaryEmotion = second.name;

  const ctx: RuleContext = {
    v,
    a,
    d,
    be5,
    primaryEmotion,
    secondaryEmotion,
    clarity,
    scores: input,
    emotionsSorted,
  };

  return {
    core_state: deriveCoreState(v, a, d),
    primary_emotion: primaryEmotion,
    secondary_emotion: secondaryEmotion,
    interpretation: deriveInterpretation(ctx),
    tension_label: deriveTensionLabel(v, a),
    control_label: deriveControlLabel(d),
    mood_label: deriveMoodLabel(v),
    clarity_label: deriveClarityLabel(clarity),
  };
}

export function translateAnalysisToGerman(result: AnalysisOutput): AnalysisOutputTranslated {
  return {
    ...result,
    core_state_de: LABELS_DE.core_state[result.core_state],
    primary_emotion_de: LABELS_DE.primary_emotion[result.primary_emotion],
    secondary_emotion_de: LABELS_DE.primary_emotion[result.secondary_emotion],
    interpretation_de: LABELS_DE.interpretation[result.interpretation],
    tension_label_de: LABELS_DE.tension_label[result.tension_label],
    control_label_de: LABELS_DE.control_label[result.control_label],
    mood_label_de: LABELS_DE.mood_label[result.mood_label],
    clarity_label_de: LABELS_DE.clarity_label[result.clarity_label],
  };
}

// Beispiel
const input: EmotionVector = {
  valence: 7.72,
  arousal: 4.77,
  dominance: 6.92,
  joy: 3.96,
  anger: 1.24,
  sadness: 1.29,
  fear: 1.37,
  disgust: 1.15,
};

const analysis = analyzeEmotionVector(input);
const analysisDe = translateAnalysisToGerman(analysis);

console.log(JSON.stringify(analysis, null, 2));
console.log(JSON.stringify(analysisDe, null, 2));

/*
analysis:
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

analysisDe:
{
  "core_state": "confident_positive",
  "primary_emotion": "joy",
  "secondary_emotion": "fear",
  "interpretation": "joy_confidence_positive_security",
  "tension_label": "medium_tension",
  "control_label": "high_control",
  "mood_label": "positive_tone",
  "clarity_label": "clear_emotional_state",
  "core_state_de": "positiv, stabil, selbstsicher",
  "primary_emotion_de": "Freude",
  "secondary_emotion_de": "Angst",
  "interpretation_de": "Freude, Zuversicht, positive Sicherheit",
  "tension_label_de": "mittlere Spannung",
  "control_label_de": "hohes Kontrollgefühl",
  "mood_label_de": "positiver Grundton",
  "clarity_label_de": "klare Emotionslage"
}
*/