import React, { useState, useMemo, useEffect } from "react";
import { Button, Form } from "react-bootstrap";
import {
  mutateMetadataScan,
  useConfigureDefaults,
} from "src/core/StashService";
import { Icon, Modal, OperationButton } from "src/components/Shared";
import { useToast } from "src/hooks";
import * as GQL from "src/core/generated-graphql";
import { FormattedMessage, useIntl } from "react-intl";
import { DirectorySelectionDialog } from "src/components/Settings/SettingsTasksPanel/DirectorySelectionDialog";
import { Manual } from "src/components/Help/Manual";
import { ScanOptions } from "./Options";
import { withoutTypename } from "src/utils";
import { ConfigurationContext } from "src/hooks/Config";

interface IScanDialogProps {
  onClose: () => void;
}

export const ScanDialog: React.FC<IScanDialogProps> = ({ onClose }) => {
  const [configureDefaults] = useConfigureDefaults();

  const [options, setOptions] = useState<GQL.ScanMetadataInput>({});
  const [paths, setPaths] = useState<string[]>([]);
  const [showManual, setShowManual] = useState(false);
  const [settingPaths, setSettingPaths] = useState(false);
  const [animation, setAnimation] = useState(true);
  const [savingDefaults, setSavingDefaults] = useState(false);

  const intl = useIntl();
  const Toast = useToast();

  const { configuration } = React.useContext(ConfigurationContext);

  useEffect(() => {
    if (!configuration?.defaults) {
      return;
    }

    const { scan } = configuration.defaults;

    if (scan) {
      setOptions(withoutTypename(scan));
    }
  }, [configuration]);

  const selectionStatus = useMemo(() => {
    const message = paths.length ? (
      <div>
        <FormattedMessage id="config.tasks.scan.scanning_paths" />:
        <ul>
          {paths.map((p) => (
            <li key={p}>{p}</li>
          ))}
        </ul>
      </div>
    ) : (
      <span>
        <FormattedMessage id="config.tasks.scan.scanning_all_paths" />.
      </span>
    );

    function onClick() {
      setAnimation(false);
      setSettingPaths(true);
    }

    return (
      <Form.Group className="dialog-selected-folders">
        <div>
          {message}
          <div>
            <Button
              title={intl.formatMessage({ id: "actions.select_folders" })}
              onClick={() => onClick()}
            >
              <Icon icon="folder-open" />
            </Button>
          </div>
        </div>
      </Form.Group>
    );
  }, [intl, paths]);

  if (!configuration) return <div />;

  function makeDefaultScanInput() {
    const ret = options;
    const { paths: _paths, ...withoutSpecifics } = ret;
    return withoutSpecifics;
  }

  async function onScan() {
    try {
      await mutateMetadataScan({
        ...options,
        paths: paths.length ? paths : undefined,
      });

      Toast.success({
        content: intl.formatMessage(
          { id: "config.tasks.added_job_to_queue" },
          { operation_name: intl.formatMessage({ id: "actions.scan" }) }
        ),
      });
    } catch (e) {
      Toast.error(e);
    } finally {
      onClose();
    }
  }

  function onShowManual() {
    setAnimation(false);
    setShowManual(true);
  }

  async function setAsDefault() {
    try {
      setSavingDefaults(true);
      await configureDefaults({
        variables: {
          input: {
            scan: makeDefaultScanInput(),
          },
        },
      });

      Toast.success({
        content: intl.formatMessage(
          { id: "config.tasks.defaults_set" },
          { action: intl.formatMessage({ id: "actions.scan" }) }
        ),
      });
    } catch (e) {
      Toast.error(e);
    } finally {
      setSavingDefaults(false);
    }
  }

  if (settingPaths) {
    return (
      <DirectorySelectionDialog
        animation={false}
        allowEmpty
        initialPaths={paths}
        onClose={(p) => {
          if (p) {
            setPaths(p);
          }
          setSettingPaths(false);
        }}
      />
    );
  }

  if (showManual) {
    return (
      <Manual
        animation={false}
        show
        onClose={() => setShowManual(false)}
        defaultActiveTab="Tasks.md"
      />
    );
  }

  return (
    <Modal
      modalProps={{ animation, size: "lg" }}
      show
      icon="cogs"
      header={intl.formatMessage({ id: "actions.scan" })}
      accept={{
        onClick: onScan,
        text: intl.formatMessage({ id: "actions.scan" }),
      }}
      cancel={{
        onClick: () => onClose(),
        text: intl.formatMessage({ id: "actions.cancel" }),
        variant: "secondary",
      }}
      disabled={savingDefaults}
      footerButtons={
        <OperationButton variant="secondary" operation={setAsDefault}>
          <FormattedMessage id="actions.set_as_default" />
        </OperationButton>
      }
      leftFooterButtons={
        <Button
          title="Help"
          className="minimal help-button"
          onClick={() => onShowManual()}
        >
          <Icon icon="question-circle" />
        </Button>
      }
    >
      <Form>
        {selectionStatus}
        <ScanOptions options={options} setOptions={(o) => setOptions(o)} />
      </Form>
    </Modal>
  );
};

export default ScanDialog;
