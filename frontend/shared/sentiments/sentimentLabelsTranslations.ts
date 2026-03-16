// sentimentLabelsTranslations.ts

export type Lang = "de" | "en";

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

export type AnalysisOutput = {
  core_state: string;
  primary_emotion: string;
  secondary_emotion: string;
  interpretation: string;
  tension_label: string;
  control_label: string;
  mood_label: string;
  clarity_label: string;
};

export const TEXTS: Record<Lang, Record<string, string>> = {
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
    mixed_emotional_state: "mixed emotional state",
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
    ambiguous_emotional_state: "ambiguous emotional state",
  },
};

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
