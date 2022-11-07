import {
  createMandatoryNumberCriterionOption,
  createStringCriterionOption,
  NullNumberCriterionOption,
} from "./criteria/criterion";
import { MovieIsMissingCriterionOption } from "./criteria/is-missing";
import { StudiosCriterionOption } from "./criteria/studios";
import { PerformersCriterionOption } from "./criteria/performers";
import { ListFilterOptions } from "./filter-options";
import { DisplayMode } from "./types";

const defaultSortBy = "name";

const sortByOptions = ["name", "random", "date", "duration", "rating"]
  .map(ListFilterOptions.createSortBy)
  .concat([
    {
      messageID: "scene_count",
      value: "scenes_count",
    },
  ]);
const displayModeOptions = [DisplayMode.Grid];
const criterionOptions = [
  StudiosCriterionOption,
  MovieIsMissingCriterionOption,
  createStringCriterionOption("url"),
  createStringCriterionOption("name"),
  createStringCriterionOption("director"),
  createStringCriterionOption("synopsis"),
  createMandatoryNumberCriterionOption("duration"),
  new NullNumberCriterionOption("rating", "rating100"),
  PerformersCriterionOption,
];

export const MovieListFilterOptions = new ListFilterOptions(
  defaultSortBy,
  sortByOptions,
  displayModeOptions,
  criterionOptions
);
