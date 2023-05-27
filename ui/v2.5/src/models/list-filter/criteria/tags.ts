import { CriterionModifier } from "src/core/generated-graphql";
import { CriterionType } from "../types";
import { CriterionOption, IHierarchicalLabeledIdCriterion } from "./criterion";

export class TagsCriterion extends IHierarchicalLabeledIdCriterion {}

class tagsCriterionOption extends CriterionOption {
  constructor(messageID: string, value: CriterionType, parameterName: string) {
    const modifierOptions = [
      CriterionModifier.Includes,
      CriterionModifier.IncludesAll,
      CriterionModifier.Equals,
    ];

    let defaultModifier = CriterionModifier.IncludesAll;

    super({
      messageID,
      type: value,
      parameterName,
      modifierOptions,
      defaultModifier,
    });
  }
}

export const TagsCriterionOption = new tagsCriterionOption(
  "tags",
  "tags",
  "tags"
);
export const SceneTagsCriterionOption = new tagsCriterionOption(
  "sceneTags",
  "sceneTags",
  "scene_tags"
);
export const PerformerTagsCriterionOption = new tagsCriterionOption(
  "performerTags",
  "performerTags",
  "performer_tags"
);
export const ParentTagsCriterionOption = new tagsCriterionOption(
  "parent_tags",
  "parentTags",
  "parents"
);
export const ChildTagsCriterionOption = new tagsCriterionOption(
  "sub_tags",
  "childTags",
  "children"
);
