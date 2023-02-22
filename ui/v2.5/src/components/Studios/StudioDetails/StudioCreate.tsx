import React, { useMemo, useState } from "react";
import { useHistory, useLocation } from "react-router-dom";
import { useIntl } from "react-intl";

import * as GQL from "src/core/generated-graphql";
import { useStudioCreate } from "src/core/StashService";
import ImageUtils from "src/utils/image";
import { LoadingIndicator } from "src/components/Shared/LoadingIndicator";
import { useToast } from "src/hooks/Toast";
import { StudioEditPanel } from "./StudioEditPanel";

const StudioCreate: React.FC = () => {
  const history = useHistory();
  const location = useLocation();
  const Toast = useToast();

  const query = useMemo(() => new URLSearchParams(location.search), [location]);
  const studio = {
    name: query.get("q") ?? undefined,
  };

  const intl = useIntl();

  // Studio state
  const [image, setImage] = useState<string | null>();

  const [createStudio] = useStudioCreate();

  function onImageLoad(imageData: string) {
    setImage(imageData);
  }

  const imageEncoding = ImageUtils.usePasteImage(onImageLoad, true);

  async function onSave(
    input: Partial<GQL.StudioCreateInput | GQL.StudioUpdateInput>
  ) {
    try {
      const result = await createStudio({
        variables: {
          input: input as GQL.StudioCreateInput,
        },
      });
      if (result.data?.studioCreate?.id) {
        history.push(`/studios/${result.data.studioCreate.id}`);
      }
    } catch (e) {
      Toast.error(e);
    }
  }

  function renderImage() {
    if (image) {
      return <img className="logo" alt="" src={image} />;
    }
  }

  return (
    <div className="row">
      <div className="studio-details col-md-8">
        <h2>
          {intl.formatMessage(
            { id: "actions.add_entity" },
            { entityType: intl.formatMessage({ id: "studio" }) }
          )}
        </h2>
        <div className="text-center">
          {imageEncoding ? (
            <LoadingIndicator message="Encoding image..." />
          ) : (
            renderImage()
          )}
        </div>
        <StudioEditPanel
          studio={studio}
          onSubmit={onSave}
          onImageChange={setImage}
          onCancel={() => history.push("/studios")}
          onDelete={() => {}}
        />
      </div>
    </div>
  );
};

export default StudioCreate;
