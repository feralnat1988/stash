import React, { useState } from "react";
import { Form } from "react-bootstrap";
import { mutateExportObjects } from "src/core/StashService";
import { Modal } from "src/components/Shared";
import { useToast } from "src/hooks";
import { downloadFile } from "src/utils";

interface IMovieExportDialogProps {
  selectedIds?: string[];
  all?: boolean;
  onClose: () => void;
}

export const MovieExportDialog: React.FC<IMovieExportDialogProps> = (
  props: IMovieExportDialogProps
) => {
  const [includeDependencies, setIncludeDependencies] = useState(true);

  // Network state
  const [isRunning, setIsRunning] = useState(false);

  const Toast = useToast();

  async function onExport() {
    try {
      setIsRunning(true);
      const ret = await mutateExportObjects({
        movies: {
          ids: props.selectedIds,
          all: props.all,
        },
        includeDependencies,
      });

      // download the result
      if (ret.data && ret.data.exportObjects) {
        const link = ret.data.exportObjects;
        downloadFile(link);
      }
    } catch (e) {
      Toast.error(e);
    } finally {
      setIsRunning(false);
      props.onClose();
    }
  }

  return (
    <Modal
      show
      icon="cogs"
      header="Export"
      accept={{ onClick: onExport, text: "Export" }}
      cancel={{
        onClick: () => props.onClose(),
        text: "Cancel",
        variant: "secondary",
      }}
      isRunning={isRunning}
    >
      <Form>
        <Form.Group>
          <Form.Check
            id="include-dependencies"
            checked={includeDependencies}
            label="Include related studio in export"
            onChange={() => setIncludeDependencies(!includeDependencies)}
          />
        </Form.Group>
      </Form>
    </Modal>
  );
};
