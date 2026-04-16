const ESTIMATE_SCALES: Record<string, { base: number[]; extended: number[] }> = {
  exponential: { base: [1, 2, 4, 8, 16], extended: [32, 64] },
  fibonacci: { base: [1, 2, 3, 5, 8, 13], extended: [21, 34] },
  linear: { base: [1, 2, 3, 4, 5], extended: [6, 7, 8, 9, 10] },
  tShirt: { base: [1, 2, 3, 4, 5], extended: [6] },
};

const TSHIRT_LABELS: Record<number, string> = {
  0: "None",
  1: "XS",
  2: "S",
  3: "M",
  4: "L",
  5: "XL",
  6: "XXL",
};

export type EstimateConfig = { type: string; allowZero: boolean; extended: boolean };

/** Build an EstimateConfig from a team object's estimation properties. */
export function buildEstimateConfig(team: {
  issueEstimationType: string;
  issueEstimationAllowZero: boolean;
  issueEstimationExtended: boolean;
}): EstimateConfig {
  return {
    type: team.issueEstimationType,
    allowZero: team.issueEstimationAllowZero,
    extended: team.issueEstimationExtended,
  };
}

/**
 * Validate an estimate value against a team's configuration.
 * Returns null if valid, or an error message string if invalid.
 */
export function validateEstimate(
  team: {
    key: string;
    issueEstimationType: string;
    issueEstimationAllowZero: boolean;
    issueEstimationExtended: boolean;
  },
  estimate: number,
): string | null {
  if (team.issueEstimationType === "notUsed") {
    return `Team "${team.key}" does not use estimates.`;
  }
  const config = buildEstimateConfig(team);
  const valid = getValidEstimates(config);
  if (!valid.includes(estimate)) {
    const scale = formatEstimateScale(team.issueEstimationType, valid);
    return `Invalid estimate: ${estimate}. Team "${team.key}" uses ${team.issueEstimationType} scale. Valid values: ${scale}`;
  }
  return null;
}

/** Returns the valid numeric estimate values for a team's configuration. */
export function getValidEstimates(config: EstimateConfig): number[] {
  const scale = ESTIMATE_SCALES[config.type];
  if (!scale) {
    return [];
  }
  const values = config.extended ? [...scale.base, ...scale.extended] : [...scale.base];
  return config.allowZero ? [0, ...values] : values;
}

/** Formats valid estimate values for display in error messages. */
export function formatEstimateScale(type: string, values: number[]): string {
  if (type === "tShirt") {
    return values.map((v) => `${v} (${TSHIRT_LABELS[v] ?? v})`).join(" | ");
  }
  return values.join(" | ");
}
