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
type Lang = "de" | "en";

type CodeMap = {
  core_state: string;
  primary_emotion: string;
  secondary_emotion: string;
  interpretation: string;
  tension_label: string;
  control_label: string;
  mood_label: string;
  clarity_label: string;
};

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

const TEXTS: Record<Lang, Record<string, string>> = {
  de: {
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

    joy: "Freude",
    anger: "Ärger",
    sadness: "Traurigkeit",
    fear: "Angst",
    disgust: "Ekel",
    mixed: "gemischt",
    none_salient: "keine dominante diskrete Emotion",

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

    high_negative_tension: "hohe negative Spannung",
    high_positive_tension: "hohe positive Spannung",
    strong_activation: "starke Aktivierung",
    low_tension: "niedrige Spannung",
    medium_tension: "mittlere Spannung",

    high_control: "hohes Kontrollgefühl",
    low_control: "niedriges Kontrollgefühl",
    medium_control: "mittleres Kontrollgefühl",

    positive_tone: "positiver Grundton",
    negative_tone: "negativer Grundton",
    ambivalent_tone: "ambivalenter Grundton",

    clear_emotional_state: "klare Emotionslage",
    mixed_emotional_state: "gemischte Emotionslage",
    ambiguous_emotional_state: "uneindeutige Emotionslage",
  },

  en: {
    threatened_activation: "threatened, tense, low control",
    hostile_activation: "negative, highly activated, forceful",
    uneasy_low_control: "unpleasant, uncertain, low control",
    depleted_negative: "empty, resigned, low energy",
    cold_negative_control: "negative, subdued, controlled",
    empowered_positive_activation: "euphoric, powerful, controlled",
    confident_positive: "positive, stable, self-assured",
    calm_positive_control: "calm, content, composed",
    soft_positive_passive: "gentle, warm, passively positive",
    focused_activation: "tense, focused, action-ready",
    agitated_uncertain: "nervous, unstable, reactive",
    flat_neutral: "flat, neutral, low activation",
    mixed_state: "mixed profile",

    joy: "joy",
    anger: "anger",
    sadness: "sadness",
    fear: "fear",
    disgust: "disgust",
    mixed: "mixed",
    none_salient: "no dominant discrete emotion",

    panic_threat_powerlessness: "panic, threat, powerlessness",
    fear_powerlessness_uncertainty: "fear, powerlessness, uncertainty",
    alarm_tension_worry: "alarm, tension, worry",
    vulnerable_despair: "vulnerable despair",
    anger_aggressive_control: "anger, aggressive control, readiness to fight",
    hostile_control_confrontation: "assertion, confrontation, hostile control",
    frustration_irritability: "frustration, irritability, reactive negativity",
    contempt_harsh_rejection: "contempt, harsh rejection",
    quiet_grief_resignation: "quiet grief, resignation, melancholy",
    sadness_loss_grief: "sadness, sorrow, loss",
    agitated_grief_despair: "despair, agitated grief",
    helpless_despair_negative_vulnerability: "helplessness, despair, vulnerable negativity",
    euphoria_triumph_enthusiasm: "euphoria, triumph, enthusiasm",
    joy_confidence_positive_security: "joy, confidence, positive security",
    contentment_warmth_calm_wellbeing: "contentment, warmth, calm wellbeing",
    positive_agency_selfconfidence: "self-confidence, positive agency",
    contempt_devaluing_rejection: "contempt, devaluing rejection",
    defensive_aversion_threat: "defensiveness, aversion, threat",
    disgust_distancing_aversion: "disgust, distancing, aversion",
    hope_fear_ambivalence: "ambivalent tension between hope and fear",
    bitter_hurt_negativity: "bitter hurt, wounded negativity",
    alarmed_aggression_defensive_tension: "alarmed aggression, defensive combat tension",
    bittersweet_ambivalence: "bittersweet ambivalence",
    ambiguous_tension: "ambiguous tension",
    positive_colored_mood: "positively colored mood",
    negative_colored_mood: "negatively colored mood",
    neutral_to_mixed: "neutral to mixed",

    high_negative_tension: "high negative tension",
    high_positive_tension: "high positive tension",
    strong_activation: "strong activation",
    low_tension: "low tension",
    medium_tension: "medium tension",

    high_control: "high sense of control",
    low_control: "low sense of control",
    medium_control: "medium sense of control",

    positive_tone: "positive tone",
    negative_tone: "negative tone",
    ambivalent_tone: "ambivalent tone",

    clear_emotional_state: "clear emotional state",
    mixed_emotional_state: "mixed emotional state",
    ambiguous_emotional_state: "ambiguous emotional state",
  },
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

export function translateCodeMap(codes: CodeMap, lang: Lang = "de"): AnalysisOutput {
  const t = TEXTS[lang];

  return {
    core_state: t[codes.core_state] ?? codes.core_state,
    primary_emotion: t[codes.primary_emotion] ?? codes.primary_emotion,
    secondary_emotion: t[codes.secondary_emotion] ?? codes.secondary_emotion,
    interpretation: t[codes.interpretation] ?? codes.interpretation,
    tension_label: t[codes.tension_label] ?? codes.tension_label,
    control_label: t[codes.control_label] ?? codes.control_label,
    mood_label: t[codes.mood_label] ?? codes.mood_label,
    clarity_label: t[codes.clarity_label] ?? codes.clarity_label,
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

const codes = analyzeEmotionVectorToCodes(input);
const de = translateCodeMap(codes, "de");
const en = translateCodeMap(codes, "en");

console.log("codes", codes);
console.log("de", de);
console.log("en", en);

/*
codes:
{
  core_state: "confident_positive",
  primary_emotion: "joy",
  secondary_emotion: "fear",
  interpretation: "joy_confidence_positive_security",
  tension_label: "medium_tension",
  control_label: "high_control",
  mood_label: "positive_tone",
  clarity_label: "clear_emotional_state"
}

de:
{
  core_state: "positiv, stabil, selbstsicher",
  primary_emotion: "Freude",
  secondary_emotion: "Angst",
  interpretation: "Freude, Zuversicht, positive Sicherheit",
  tension_label: "mittlere Spannung",
  control_label: "hohes Kontrollgefühl",
  mood_label: "positiver Grundton",
  clarity_label: "klare Emotionslage"
}

en:
{
  core_state: "positive, stable, self-assured",
  primary_emotion: "joy",
  secondary_emotion: "fear",
  interpretation: "joy, confidence, positive security",
  tension_label: "medium tension",
  control_label: "high sense of control",
  mood_label: "positive tone",
  clarity_label: "clear emotional state"
}
*/
