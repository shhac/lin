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

type EstimateConfig = { type: string; allowZero: boolean; extended: boolean };

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
