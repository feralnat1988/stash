import React from "react";
import * as GQL from "src/core/generated-graphql";
import { SceneList } from "src/components/Scenes/SceneList";
import { usePerformerFilterHook } from "src/core/performers";
import { PersistanceLevel } from "src/components/List/ItemList";

interface IPerformerDetailsProps {
  active: boolean;
  performer: GQL.PerformerDataFragment;
}

export const PerformerScenesPanel: React.FC<IPerformerDetailsProps> = ({
  active,
  performer,
}) => {
  const filterHook = usePerformerFilterHook(performer);
  return (
    <SceneList
      filterHook={filterHook}
      alterQuery={active}
      persistState={PersistanceLevel.SAVEDLINKFILTER}
    />
  );
};
