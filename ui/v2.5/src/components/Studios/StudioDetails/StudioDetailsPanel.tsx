import React from "react";
import { useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import { TextUtils } from "src/utils";
import { RatingStars } from "src/components/Scenes/SceneDetails/RatingStars";
import { TextField, URLField } from "src/utils/field";

interface IStudioDetailsPanel {
  studio: Partial<GQL.StudioDataFragment>;
}

export const StudioDetailsPanel: React.FC<IStudioDetailsPanel> = ({
  studio,
}) => {
  const intl = useIntl();

  function renderRatingField() {
    if (!studio.rating) {
      return;
    }

    return (
      <>
        <dt>{intl.formatMessage({ id: "rating" })}</dt>
        <dd>
          <RatingStars value={studio.rating} disabled />
        </dd>
      </>
    );
  }

  return (
    <div className="studio-details">
      <div>
        <h2>{studio.name}</h2>
      </div>

      <dl className="details-list">
        <URLField
          id="url"
          value={studio.url}
          url={TextUtils.sanitiseURL(studio.url ?? "")}
        />

        <TextField id="details" value={studio.details} />

        <URLField
          id="parent_studios"
          value={studio.parent_studio?.name}
          url={`/studios/${studio.parent_studio?.id}`}
          trusted
          target="_self"
        />

        {renderRatingField()}
      </dl>
    </div>
  );
};
