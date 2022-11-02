import {
  convertFromRatingFormat,
  convertToRatingFormat,
} from "../../../components/Scenes/SceneDetails/RatingSystem";
import {
  CriterionModifier,
  IntCriterionInput,
  RatingSystem,
} from "../../../core/generated-graphql";
import { INumberValue } from "../types";
import { Criterion, CriterionOption } from "./criterion";

export class RatingCriterion extends Criterion<INumberValue> {
  ratingSystem: RatingSystem;

  public get value(): INumberValue {
    return this._value;
  }
  public set value(newValue: number | INumberValue) {
    // backwards compatibility - if this.value is a number, use that
    if (typeof newValue !== "object") {
      this._value = {
        value: convertFromRatingFormat(newValue, this.ratingSystem),
        value2: undefined,
      };
    } else {
      this._value = newValue;
    }
  }

  protected toCriterionInput(): IntCriterionInput {
    return {
      modifier: this.modifier,
      value: this.value.value,
      value2: this.value.value2,
    };
  }

  public getLabelValue() {
    const { value, value2 } = this.value;
    if (
      this.modifier === CriterionModifier.Between ||
      this.modifier === CriterionModifier.NotBetween
    ) {
      return `${convertToRatingFormat(value, this.ratingSystem)}, ${
        convertToRatingFormat(value2, this.ratingSystem) ?? 0
      }`;
    } else {
      return `${convertToRatingFormat(value, this.ratingSystem)}`;
    }
  }

  constructor(type: CriterionOption, ratingSystem: RatingSystem) {
    super(type, { value: 0, value2: undefined });
    this.ratingSystem = ratingSystem;
  }
}
