import React, { useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { Button, ButtonGroup, Card, Form } from "react-bootstrap";
import {
  mutateMetadataImport,
  mutateMetadataExport,
  mutateMigrateHashNaming,
  usePlugins,
  mutateRunPluginTask,
  mutateBackupDatabase,
  mutateMetadataScan,
  mutateMetadataIdentify,
  mutateMetadataAutoTag,
  mutateMetadataGenerate,
} from "src/core/StashService";
import { useToast } from "src/hooks";
import * as GQL from "src/core/generated-graphql";
import { LoadingIndicator, Modal } from "src/components/Shared";
import { downloadFile, withoutTypename } from "src/utils";
import IdentifyDialog from "src/components/Dialogs/IdentifyDialog/IdentifyDialog";
import { ImportDialog } from "./ImportDialog";
import { JobTable } from "./JobTable";
import ScanDialog from "src/components/Dialogs/ScanDialog/ScanDialog";
import AutoTagDialog from "src/components/Dialogs/AutoTagDialog";
import { GenerateDialog } from "src/components/Dialogs/GenerateDialog";
import CleanDialog from "src/components/Dialogs/CleanDialog";
import { ConfigurationContext } from "src/hooks/Config";
import { PropsWithChildren } from "react-router/node_modules/@types/react";

type Plugin = Pick<GQL.Plugin, "id">;
type PluginTask = Pick<GQL.PluginTask, "name" | "description">;

interface ITask {
  description?: React.ReactNode;
}

const Task: React.FC<PropsWithChildren<ITask>> = ({
  children,
  description,
}) => (
  <div className="task">
    {children}
    {description ? (
      <Form.Text className="text-muted">{description}</Form.Text>
    ) : undefined}
  </div>
);

export const SettingsTasksPanel: React.FC = () => {
  const intl = useIntl();
  const Toast = useToast();
  const [dialogOpen, setDialogOpenState] = useState({
    importAlert: false,
    import: false,
    clean: false,
    scan: false,
    autoTag: false,
    identify: false,
    generate: false,
  });

  type DialogOpenState = typeof dialogOpen;

  const [isBackupRunning, setIsBackupRunning] = useState<boolean>(false);

  const { configuration } = React.useContext(ConfigurationContext);

  const plugins = usePlugins();

  function setDialogOpen(s: Partial<DialogOpenState>) {
    setDialogOpenState((v) => {
      return { ...v, ...s };
    });
  }

  async function onImport() {
    setDialogOpen({ importAlert: false });
    try {
      await mutateMetadataImport();
      Toast.success({
        content: intl.formatMessage(
          { id: "config.tasks.added_job_to_queue" },
          { operation_name: intl.formatMessage({ id: "actions.import" }) }
        ),
      });
    } catch (e) {
      Toast.error(e);
    }
  }

  function renderImportAlert() {
    return (
      <Modal
        show={dialogOpen.importAlert}
        icon="trash-alt"
        accept={{
          text: intl.formatMessage({ id: "actions.import" }),
          variant: "danger",
          onClick: onImport,
        }}
        cancel={{ onClick: () => setDialogOpen({ importAlert: false }) }}
      >
        <p>{intl.formatMessage({ id: "actions.tasks.import_warning" })}</p>
      </Modal>
    );
  }

  function renderCleanDialog() {
    if (!dialogOpen.clean) {
      return;
    }

    return <CleanDialog onClose={() => setDialogOpen({ clean: false })} />;
  }

  function renderImportDialog() {
    if (!dialogOpen.import) {
      return;
    }

    return <ImportDialog onClose={() => setDialogOpen({ import: false })} />;
  }

  function renderScanDialog() {
    if (!dialogOpen.scan) {
      return;
    }

    return <ScanDialog onClose={() => setDialogOpen({ scan: false })} />;
  }

  function renderAutoTagDialog() {
    if (!dialogOpen.autoTag) {
      return;
    }

    return <AutoTagDialog onClose={() => setDialogOpen({ autoTag: false })} />;
  }

  function maybeRenderIdentifyDialog() {
    if (!dialogOpen.identify) return;

    return (
      <IdentifyDialog onClose={() => setDialogOpen({ identify: false })} />
    );
  }

  function maybeRenderGenerateDialog() {
    if (!dialogOpen.generate) return;

    return (
      <GenerateDialog onClose={() => setDialogOpen({ generate: false })} />
    );
  }

  async function onPluginTaskClicked(plugin: Plugin, operation: PluginTask) {
    await mutateRunPluginTask(plugin.id, operation.name);
    Toast.success({
      content: intl.formatMessage(
        { id: "config.tasks.added_job_to_queue" },
        { operation_name: operation.name }
      ),
    });
  }

  function renderPluginTasks(plugin: Plugin, pluginTasks: PluginTask[]) {
    if (!pluginTasks) {
      return;
    }

    return pluginTasks.map((o) => {
      return (
        <Task description={o.description} key={o.name}>
          <Button
            onClick={() => onPluginTaskClicked(plugin, o)}
            variant="secondary"
            size="sm"
          >
            {o.name}
          </Button>
        </Task>
      );
    });
  }

  async function onBackup(download?: boolean) {
    try {
      setIsBackupRunning(true);
      const ret = await mutateBackupDatabase({
        download,
      });

      // download the result
      if (download && ret.data && ret.data.backupDatabase) {
        const link = ret.data.backupDatabase;
        downloadFile(link);
      }
    } catch (e) {
      Toast.error(e);
    } finally {
      setIsBackupRunning(false);
    }
  }

  function renderPlugins() {
    if (!plugins.data || !plugins.data.plugins) {
      return;
    }

    const taskPlugins = plugins.data.plugins.filter(
      (p) => p.tasks && p.tasks.length > 0
    );

    return (
      <>
        <hr />

        <Form.Group>
          <h5>{intl.formatMessage({ id: "config.tasks.plugin_tasks" })}</h5>
          {taskPlugins.map((o) => {
            return (
              <Form.Group key={`${o.id}`}>
                <h6>{o.name}</h6>
                <Card className="task-group">
                  {renderPluginTasks(o, o.tasks ?? [])}
                </Card>
              </Form.Group>
            );
          })}
        </Form.Group>
      </>
    );
  }

  async function onMigrateHashNaming() {
    try {
      await mutateMigrateHashNaming();
      Toast.success({
        content: intl.formatMessage(
          { id: "config.tasks.added_job_to_queue" },
          {
            operation_name: intl.formatMessage({
              id: "actions.hash_migration",
            }),
          }
        ),
      });
    } catch (err) {
      Toast.error(err);
    }
  }

  async function onExport() {
    try {
      await mutateMetadataExport();
      Toast.success({
        content: intl.formatMessage(
          { id: "config.tasks.added_job_to_queue" },
          { operation_name: intl.formatMessage({ id: "actions.backup" }) }
        ),
      });
    } catch (err) {
      Toast.error(err);
    }
  }

  async function onScanClicked() {
    // check if defaults are set for scan
    // if not, then open the dialog
    if (!configuration) {
      return;
    }

    const { scan } = configuration?.defaults;
    if (!scan) {
      setDialogOpen({ scan: true });
    } else {
      mutateMetadataScan(withoutTypename(scan));
    }
  }

  async function onIdentifyClicked() {
    // check if defaults are set for identify
    // if not, then open the dialog
    if (!configuration) {
      return;
    }

    const { identify } = configuration?.defaults;
    if (!identify) {
      setDialogOpen({ identify: true });
    } else {
      mutateMetadataIdentify(withoutTypename(identify));
    }
  }

  async function onAutoTagClicked() {
    // check if defaults are set for auto tag
    // if not, then open the dialog
    if (!configuration) {
      return;
    }

    const { autoTag } = configuration?.defaults;
    if (!autoTag) {
      setDialogOpen({ autoTag: true });
    } else {
      mutateMetadataAutoTag(withoutTypename(autoTag));
    }
  }

  async function onGenerateClicked() {
    // check if defaults are set for generate
    // if not, then open the dialog
    if (!configuration) {
      return;
    }

    const { generate } = configuration?.defaults;
    if (!generate) {
      setDialogOpen({ generate: true });
    } else {
      mutateMetadataGenerate(withoutTypename(generate));
    }
  }

  if (isBackupRunning) {
    return (
      <LoadingIndicator
        message={intl.formatMessage({ id: "config.tasks.backing_up_database" })}
      />
    );
  }

  return (
    <>
      {renderImportAlert()}
      {renderCleanDialog()}
      {renderImportDialog()}
      {renderScanDialog()}
      {renderAutoTagDialog()}
      {maybeRenderIdentifyDialog()}
      {maybeRenderGenerateDialog()}

      <h4>{intl.formatMessage({ id: "config.tasks.job_queue" })}</h4>

      <JobTable />

      <hr />

      <Form.Group>
        <h5>{intl.formatMessage({ id: "library" })}</h5>
        <Card className="task-group">
          <Task
            description={intl.formatMessage({
              id: "config.tasks.scan_for_content_desc",
            })}
          >
            <ButtonGroup className="ellipsis-button">
              <Button
                variant="secondary"
                type="submit"
                onClick={() => onScanClicked()}
              >
                <FormattedMessage id="actions.scan" />
              </Button>
              <Button
                variant="secondary"
                onClick={() => setDialogOpen({ scan: true })}
              >
                …
              </Button>
            </ButtonGroup>
          </Task>

          <Task
            description={intl.formatMessage({
              id: "config.tasks.identify.description",
            })}
          >
            <ButtonGroup className="ellipsis-button">
              <Button
                variant="secondary"
                type="submit"
                onClick={() => onIdentifyClicked()}
              >
                <FormattedMessage id="actions.identify" />
              </Button>
              <Button
                variant="secondary"
                onClick={() => setDialogOpen({ identify: true })}
              >
                …
              </Button>
            </ButtonGroup>
          </Task>

          <Task
            description={intl.formatMessage({
              id: "config.tasks.auto_tag_based_on_filenames",
            })}
          >
            <ButtonGroup className="ellipsis-button">
              <Button
                variant="secondary"
                type="submit"
                onClick={() => onAutoTagClicked()}
              >
                <FormattedMessage id="actions.auto_tag" />
              </Button>
              <Button
                variant="secondary"
                onClick={() => setDialogOpen({ autoTag: true })}
              >
                …
              </Button>
            </ButtonGroup>
          </Task>

          <Task
            description={intl.formatMessage({
              id: "config.tasks.cleanup_desc",
            })}
          >
            <Button
              variant="danger"
              type="submit"
              onClick={() => setDialogOpen({ clean: true })}
            >
              <FormattedMessage id="actions.clean" />…
            </Button>
          </Task>
        </Card>
      </Form.Group>

      <hr />

      <Form.Group>
        <h5>{intl.formatMessage({ id: "config.tasks.generated_content" })}</h5>

        <Card className="task-group">
          <Task
            description={intl.formatMessage({
              id: "config.tasks.generate_desc",
            })}
          >
            <ButtonGroup className="ellipsis-button">
              <Button
                variant="secondary"
                type="submit"
                onClick={() => onGenerateClicked()}
              >
                <FormattedMessage id="actions.generate" />
              </Button>
              <Button
                variant="secondary"
                onClick={() => setDialogOpen({ generate: true })}
              >
                …
              </Button>
            </ButtonGroup>
          </Task>
        </Card>
      </Form.Group>

      <hr />

      <Form.Group>
        <h5>{intl.formatMessage({ id: "metadata" })}</h5>
        <Card className="task-group">
          <Task
            description={intl.formatMessage({
              id: "config.tasks.export_to_json",
            })}
          >
            <Button
              id="export"
              variant="secondary"
              type="submit"
              onClick={() => onExport()}
            >
              <FormattedMessage id="actions.full_export" />
            </Button>
          </Task>

          <Task
            description={intl.formatMessage({
              id: "config.tasks.import_from_exported_json",
            })}
          >
            <Button
              id="import"
              variant="danger"
              type="submit"
              onClick={() => setDialogOpen({ importAlert: true })}
            >
              <FormattedMessage id="actions.full_import" />
            </Button>
          </Task>

          <Task
            description={intl.formatMessage({
              id: "config.tasks.incremental_import",
            })}
          >
            <Button
              id="partial-import"
              variant="danger"
              type="submit"
              onClick={() => setDialogOpen({ import: true })}
            >
              <FormattedMessage id="actions.import_from_file" />
            </Button>
          </Task>
        </Card>
      </Form.Group>

      <hr />

      <Form.Group>
        <h5>{intl.formatMessage({ id: "actions.backup" })}</h5>
        <Card className="task-group">
          <Task
            description={intl.formatMessage(
              { id: "config.tasks.backup_database" },
              {
                filename_format: (
                  <code>
                    [origFilename].sqlite.[schemaVersion].[YYYYMMDD_HHMMSS]
                  </code>
                ),
              }
            )}
          >
            <Button
              id="backup"
              variant="secondary"
              type="submit"
              onClick={() => onBackup()}
            >
              <FormattedMessage id="actions.backup" />
            </Button>
          </Task>

          <Task
            description={intl.formatMessage({
              id: "config.tasks.backup_and_download",
            })}
          >
            <Button
              id="backupDownload"
              variant="secondary"
              type="submit"
              onClick={() => onBackup(true)}
            >
              <FormattedMessage id="actions.download_backup" />
            </Button>
          </Task>
        </Card>
      </Form.Group>

      {renderPlugins()}

      <hr />

      <Form.Group>
        <h5>{intl.formatMessage({ id: "config.tasks.migrations" })}</h5>

        <Card className="task-group">
          <Task
            description={intl.formatMessage({
              id: "config.tasks.migrate_hash_files",
            })}
          >
            <Button
              id="migrateHashNaming"
              variant="danger"
              onClick={() => onMigrateHashNaming()}
            >
              <FormattedMessage id="actions.rename_gen_files" />
            </Button>
          </Task>
        </Card>
      </Form.Group>
    </>
  );
};
