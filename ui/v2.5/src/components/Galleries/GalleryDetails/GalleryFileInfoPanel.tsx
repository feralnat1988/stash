import React, { useMemo, useState } from "react";
import { Accordion, Button, Card } from "react-bootstrap";
import { FormattedMessage } from "react-intl";
import { TruncatedText } from "src/components/Shared";
import DeleteFilesDialog from "src/components/Shared/DeleteFilesDialog";
import ReassignFilesDialog from "src/components/Shared/ReassignFilesDialog";
import * as GQL from "src/core/generated-graphql";
import { mutateGallerySetPrimaryFile } from "src/core/StashService";
import { useToast } from "src/hooks";
import { TextUtils } from "src/utils";
import { TextField, URLField } from "src/utils/field";

interface IFileInfoPanelProps {
  folder?: Pick<GQL.Folder, "id" | "path">;
  file?: GQL.GalleryFileDataFragment;
  primary?: boolean;
  ofMany?: boolean;
  onSetPrimaryFile?: () => void;
  onDeleteFile?: () => void;
  onReassign?: () => void;
  loading?: boolean;
}

const FileInfoPanel: React.FC<IFileInfoPanelProps> = (
  props: IFileInfoPanelProps
) => {
  const checksum = props.file?.fingerprints.find((f) => f.type === "md5");
  const path = props.folder ? props.folder.path : props.file?.path ?? "";
  const id = props.folder ? "folder" : "path";

  return (
    <div>
      <dl className="container gallery-file-info details-list">
        {props.primary && (
          <>
            <dt></dt>
            <dd className="primary-file">
              <FormattedMessage id="primary_file" />
            </dd>
          </>
        )}
        <TextField id="media_info.checksum" value={checksum?.value} truncate />
        <URLField
          id={id}
          url={`file://${path}`}
          value={`file://${path}`}
          truncate
        />
      </dl>
      {props.ofMany && props.onSetPrimaryFile && !props.primary && (
        <div>
          <Button
            className="edit-button"
            disabled={props.loading}
            onClick={props.onSetPrimaryFile}
          >
            <FormattedMessage id="actions.make_primary" />
          </Button>
          <Button
            className="edit-button"
            disabled={props.loading}
            onClick={props.onReassign}
          >
            <FormattedMessage id="actions.reassign" />
          </Button>
          <Button
            variant="danger"
            disabled={props.loading}
            onClick={props.onDeleteFile}
          >
            <FormattedMessage id="actions.delete_file" />
          </Button>
        </div>
      )}
    </div>
  );
};
interface IGalleryFileInfoPanelProps {
  gallery: GQL.GalleryDataFragment;
}

export const GalleryFileInfoPanel: React.FC<IGalleryFileInfoPanelProps> = (
  props: IGalleryFileInfoPanelProps
) => {
  const Toast = useToast();

  const [loading, setLoading] = useState(false);
  const [deletingFile, setDeletingFile] = useState<
    GQL.GalleryFileDataFragment | undefined
  >();
  const [reassigningFile, setReassigningFile] = useState<
    GQL.GalleryFileDataFragment | undefined
  >();

  const filesPanel = useMemo(() => {
    if (props.gallery.folder) {
      return <FileInfoPanel folder={props.gallery.folder} />;
    }

    if (props.gallery.files.length === 0) {
      return <></>;
    }

    if (props.gallery.files.length === 1) {
      return <FileInfoPanel file={props.gallery.files[0]} />;
    }

    async function onSetPrimaryFile(fileID: string) {
      try {
        setLoading(true);
        await mutateGallerySetPrimaryFile(props.gallery.id, fileID);
      } catch (e) {
        Toast.error(e);
      } finally {
        setLoading(false);
      }
    }

    return (
      <Accordion defaultActiveKey={props.gallery.files[0].id}>
        {deletingFile && (
          <DeleteFilesDialog
            onClose={() => setDeletingFile(undefined)}
            selected={[deletingFile]}
          />
        )}
        {reassigningFile && (
          <ReassignFilesDialog
            type="galleries"
            onClose={() => setReassigningFile(undefined)}
            selected={[reassigningFile]}
            reassign={() => {}}
          />
        )}
        {props.gallery.files.map((file, index) => (
          <Card key={file.id} className="gallery-file-card">
            <Accordion.Toggle as={Card.Header} eventKey={file.id}>
              <TruncatedText text={TextUtils.fileNameFromPath(file.path)} />
            </Accordion.Toggle>
            <Accordion.Collapse eventKey={file.id}>
              <Card.Body>
                <FileInfoPanel
                  file={file}
                  primary={index === 0}
                  ofMany
                  onSetPrimaryFile={() => onSetPrimaryFile(file.id)}
                  loading={loading}
                  onDeleteFile={() => setDeletingFile(file)}
                  onReassign={() => setReassigningFile(file)}
                />
              </Card.Body>
            </Accordion.Collapse>
          </Card>
        ))}
      </Accordion>
    );
  }, [props.gallery, loading, Toast, deletingFile, reassigningFile]);

  return (
    <>
      <dl className="container gallery-file-info details-list">
        <URLField
          id="media_info.downloaded_from"
          url={props.gallery.url}
          value={props.gallery.url}
          truncate
        />
      </dl>

      {filesPanel}
    </>
  );
};
