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

type AnalysisOutput = {
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
  primaryEmotion: EmotionName | "mixed" | "none_salient";
  secondaryEmotion: EmotionName;
  clarity: "clear" | "mixed" | "ambiguous";
  scores: EmotionVector;
  emotionsSorted: Array<{ name: EmotionName; value: number }>;
};

type DecisionRule = {
  id: string;
  when: (ctx: RuleContext) => boolean;
  then: string;
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
): EmotionName | "mixed" | "none_salient" {
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
    then: "Panik, Bedrohung, Ausgeliefertsein",
  },
  {
    id: "H2",
    when: (c) =>
      c.be5.fear === "high" && c.v === "low" && c.d === "low",
    then: "Angst, Ohnmacht, Unsicherheit",
  },
  {
    id: "H3",
    when: (c) =>
      c.be5.fear === "high" && c.a === "high" && (c.d === "mid" || c.d === "high"),
    then: "Alarm, Anspannung, Sorge",
  },
  {
    id: "H4",
    when: (c) =>
      (c.be5.fear === "mid" || c.be5.fear === "high") &&
      c.be5.sadness === "high" &&
      c.d === "low",
    then: "verletzliche Verzweiflung",
  },
  {
    id: "H5",
    when: (c) =>
      c.be5.anger === "high" && c.v === "low" && c.a === "high" && c.d === "high",
    then: "Wut, aggressive Kontrolle, Kampfbereitschaft",
  },
  {
    id: "H6",
    when: (c) =>
      c.be5.anger === "high" && c.d === "high",
    then: "Durchsetzung, Konfrontation, feindselige Kontrolle",
  },
  {
    id: "H7",
    when: (c) =>
      c.be5.anger === "high" && c.d === "low",
    then: "Frustration, Gereiztheit, reaktive Negativität",
  },
  {
    id: "H8",
    when: (c) =>
      (c.be5.anger === "mid" || c.be5.anger === "high") &&
      c.be5.disgust === "high",
    then: "Verachtung, harte Ablehnung",
  },
  {
    id: "H9",
    when: (c) =>
      c.be5.sadness === "high" && c.v === "low" && c.a === "low",
    then: "stille Trauer, Resignation, Melancholie",
  },
  {
    id: "H10",
    when: (c) =>
      c.be5.sadness === "high" && c.v === "low" && c.a === "mid",
    then: "Traurigkeit, Kummer, Verlust",
  },
  {
    id: "H11",
    when: (c) =>
      c.be5.sadness === "high" && c.a === "high",
    then: "Verzweiflung, aufgewühlte Trauer",
  },
  {
    id: "H12",
    when: (c) =>
      c.be5.sadness === "high" && c.be5.fear === "high" && c.d === "low",
    then: "Hilflosigkeit, Verzweiflung, verletzliche Negativität",
  },
  {
    id: "H13",
    when: (c) =>
      c.be5.joy === "high" && c.v === "high" && c.a === "high" && c.d === "high",
    then: "Euphorie, Triumph, Begeisterung",
  },
  {
    id: "H14",
    when: (c) =>
      c.be5.joy === "high" && c.v === "high" && c.a === "mid" && c.d === "high",
    then: "Freude, Zuversicht, positive Sicherheit",
  },
  {
    id: "H15",
    when: (c) =>
      c.be5.joy === "high" && c.v === "high" && c.a === "low",
    then: "Zufriedenheit, Wärme, ruhiges Wohlgefühl",
  },
  {
    id: "H16",
    when: (c) =>
      (c.be5.joy === "mid" || c.be5.joy === "high") &&
      c.d === "high" &&
      c.v === "high",
    then: "Selbstsicherheit, positive Handlungsfähigkeit",
  },
  {
    id: "H17",
    when: (c) =>
      c.be5.disgust === "high" && c.be5.anger === "high",
    then: "Verachtung, abwertende Ablehnung",
  },
  {
    id: "H18",
    when: (c) =>
      c.be5.disgust === "high" && c.be5.fear === "high",
    then: "Abwehr, Abstoßung, aversive Bedrohung",
  },
  {
    id: "H19",
    when: (c) =>
      c.be5.disgust === "high",
    then: "Ekel, Distanzierung, Abstoßung",
  },
  {
    id: "H20",
    when: (c) =>
      c.be5.joy === "high" && c.be5.fear === "high",
    then: "ambivalente Spannung zwischen Hoffnung und Angst",
  },
  {
    id: "H21",
    when: (c) =>
      c.be5.sadness === "high" && c.be5.anger === "high",
    then: "bitterer Schmerz, kränkende Negativität",
  },
  {
    id: "H22",
    when: (c) =>
      c.be5.fear === "high" && c.be5.anger === "high",
    then: "alarmierte Aggression, defensive Kampfspannung",
  },
  {
    id: "H23",
    when: (c) =>
      c.be5.joy === "high" && c.be5.sadness === "high",
    then: "bittersüße Ambivalenz",
  },
  {
    id: "H24",
    when: (c) =>
      c.clarity === "ambiguous" && c.a === "high" && c.v === "low",
    then: "uneindeutige Spannung",
  },
  {
    id: "H25",
    when: (c) =>
      c.clarity === "ambiguous" && c.v === "mid",
    then: "gemischte emotionale Lage",
  },
  {
    id: "H26",
    when: (c) => c.v === "high",
    then: "positiv gefärbte Stimmung",
  },
  {
    id: "H27",
    when: (c) => c.v === "low",
    then: "negativ gefärbte Stimmung",
  },
  {
    id: "H28",
    when: () => true,
    then: "neutral bis gemischt",
  },
];

function deriveInterpretation(ctx: RuleContext): string {
  const matched = INTERPRETATION_RULES.find((rule) => rule.when(ctx));
  return matched ? matched.then : "neutral bis gemischt";
}

function deriveTensionLabel(v: Level3, a: Level3): string {
  if (a === "high" && v === "low") return "hohe negative Spannung";
  if (a === "high" && v === "high") return "hohe positive Spannung";
  if (a === "high" && v === "mid") return "starke Aktivierung";
  if (a === "low") return "niedrige Spannung";
  return "mittlere Spannung";
}

function deriveControlLabel(d: Level3): string {
  if (d === "high") return "hohes Kontrollgefühl";
  if (d === "low") return "niedriges Kontrollgefühl";
  return "mittleres Kontrollgefühl";
}

function deriveMoodLabel(v: Level3): string {
  if (v === "high") return "positiver Grundton";
  if (v === "low") return "negativer Grundton";
  return "ambivalenter Grundton";
}

function deriveClarityLabel(clarity: "clear" | "mixed" | "ambiguous"): string {
  if (clarity === "clear") return "klare Emotionslage";
  if (clarity === "mixed") return "gemischte Emotionslage";
  return "uneindeutige Emotionslage";
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

console.log(JSON.stringify(analyzeEmotionVector(input), null, 2));

/*
{
  "core_state": "confident_positive",
  "primary_emotion": "joy",
  "secondary_emotion": "fear",
  "interpretation": "Freude, Zuversicht, positive Sicherheit",
  "tension_label": "mittlere Spannung",
  "control_label": "hohes Kontrollgefühl",
  "mood_label": "positiver Grundton",
  "clarity_label": "klare Emotionslage"
}
*/