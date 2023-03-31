import * as GQL from "../core/generated-graphql";

export const stringCircumMap = new Map<string, GQL.CircumEnum>([
  ["Uncut", GQL.CircumEnum.Uncut],
  ["Cut", GQL.CircumEnum.Cut],
]);

export const circumcisedToString = (value?: GQL.CircumEnum | String | null) => {
  if (!value) {
    return undefined;
  }

  const foundEntry = Array.from(stringCircumMap.entries()).find((e) => {
    return e[1] === value;
  });

  if (foundEntry) {
    return foundEntry[0];
  }
};

export const stringToCircumcised = (
  value?: string | null,
  caseInsensitive?: boolean
): GQL.CircumEnum | undefined => {
  if (!value) {
    return undefined;
  }

  const existing = Object.entries(GQL.CircumEnum).find((e) => e[1] === value);
  if (existing) return existing[1];

  const ret = stringCircumMap.get(value);
  if (ret || !caseInsensitive) {
    return ret;
  }
  const asUpper = value.toUpperCase();
  const foundEntry = Array.from(stringCircumMap.entries()).find((e) => {
    return e[0].toUpperCase() === asUpper;
  });

  if (foundEntry) {
    return foundEntry[1];
  }
};

export const circumcisedStrings = Array.from(stringCircumMap.keys());
