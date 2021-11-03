import React, { useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { Button, Form } from "react-bootstrap";
import {
  mutateMetadataImport,
  mutateMetadataExport,
  mutateMigrateHashNaming,
  usePlugins,
  mutateRunPluginTask,
  mutateBackupDatabase,
} from "src/core/StashService";
import { useToast } from "src/hooks";
import * as GQL from "src/core/generated-graphql";
import { LoadingIndicator, Modal } from "src/components/Shared";
import { downloadFile } from "src/utils";
import IdentifyDialog from "src/components/Dialogs/IdentifyDialog/IdentifyDialog";
import { ImportDialog } from "./ImportDialog";
import { JobTable } from "./JobTable";
import ScanDialog from "src/components/Dialogs/ScanDialog/ScanDialog";
import AutoTagDialog from "src/components/Dialogs/AutoTagDialog";
import { GenerateDialog } from "src/components/Dialogs/GenerateDialog";
import CleanDialog from "src/components/Dialogs/CleanDialog";

type Plugin = Pick<GQL.Plugin, "id">;
type PluginTask = Pick<GQL.PluginTask, "name" | "description">;

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
        <div key={o.name}>
          <Button
            onClick={() => onPluginTaskClicked(plugin, o)}
            className="mt-3"
            variant="secondary"
            size="sm"
          >
            {o.name}
          </Button>
          {o.description ? (
            <Form.Text className="text-muted">{o.description}</Form.Text>
          ) : undefined}
        </div>
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
        <h5>{intl.formatMessage({ id: "config.tasks.plugin_tasks" })}</h5>
        {taskPlugins.map((o) => {
          return (
            <div key={`${o.id}`} className="mb-3">
              <h6>{o.name}</h6>
              {renderPluginTasks(o, o.tasks ?? [])}
              <hr />
            </div>
          );
        })}
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
        <Form.Group>
          <Button
            className="mr-2"
            variant="secondary"
            type="submit"
            onClick={() => setDialogOpen({ scan: true })}
          >
            <FormattedMessage id="actions.scan" />…
          </Button>
          <Form.Text className="text-muted">
            {intl.formatMessage({ id: "config.tasks.scan_for_content_desc" })}
          </Form.Text>
        </Form.Group>

        <Form.Group>
          <Button
            className="mr-2"
            variant="secondary"
            type="submit"
            onClick={() => setDialogOpen({ identify: true })}
          >
            <FormattedMessage id="actions.identify" />…
          </Button>
          <Form.Text className="text-muted">
            <FormattedMessage id="config.tasks.identify.description" />
          </Form.Text>
        </Form.Group>

        <Form.Group>
          <Button
            variant="secondary"
            type="submit"
            className="mr-2"
            onClick={() => setDialogOpen({ autoTag: true })}
          >
            <FormattedMessage id="actions.auto_tag" />…
          </Button>
          <Form.Text className="text-muted">
            {intl.formatMessage({
              id: "config.tasks.auto_tag_based_on_filenames",
            })}
          </Form.Text>
        </Form.Group>

        <Form.Group>
          <Button
            id="clean"
            variant="danger"
            onClick={() => setDialogOpen({ clean: true })}
          >
            <FormattedMessage id="actions.clean" />…
          </Button>
          <Form.Text className="text-muted">
            {intl.formatMessage({ id: "config.tasks.cleanup_desc" })}
          </Form.Text>
        </Form.Group>
      </Form.Group>

      <hr />

      <Form.Group>
        <h5>{intl.formatMessage({ id: "config.tasks.generated_content" })}</h5>
        <Button
          id="generate"
          variant="secondary"
          type="submit"
          onClick={() => setDialogOpen({ generate: true })}
        >
          <FormattedMessage id="actions.generate" />…
        </Button>
        <Form.Text className="text-muted">
          {intl.formatMessage({ id: "config.tasks.generate_desc" })}
        </Form.Text>
      </Form.Group>

      <hr />

      <h5>{intl.formatMessage({ id: "metadata" })}</h5>
      <Form.Group>
        <Button
          id="export"
          variant="secondary"
          type="submit"
          onClick={() => onExport()}
        >
          <FormattedMessage id="actions.full_export" />
        </Button>
        <Form.Text className="text-muted">
          {intl.formatMessage({ id: "config.tasks.export_to_json" })}
        </Form.Text>
      </Form.Group>

      <Form.Group>
        <Button
          id="import"
          variant="danger"
          onClick={() => setDialogOpen({ importAlert: true })}
        >
          <FormattedMessage id="actions.full_import" />
        </Button>
        <Form.Text className="text-muted">
          {intl.formatMessage({ id: "config.tasks.import_from_exported_json" })}
        </Form.Text>
      </Form.Group>

      <Form.Group>
        <Button
          id="partial-import"
          variant="danger"
          onClick={() => setDialogOpen({ import: true })}
        >
          <FormattedMessage id="actions.import_from_file" />
        </Button>
        <Form.Text className="text-muted">
          {intl.formatMessage({ id: "config.tasks.incremental_import" })}
        </Form.Text>
      </Form.Group>

      <hr />

      <h5>{intl.formatMessage({ id: "actions.backup" })}</h5>
      <Form.Group>
        <Button
          id="backup"
          variant="secondary"
          type="submit"
          onClick={() => onBackup()}
        >
          <FormattedMessage id="actions.backup" />
        </Button>
        <Form.Text className="text-muted">
          {intl.formatMessage(
            { id: "config.tasks.backup_database" },
            {
              filename_format: (
                <code>
                  [origFilename].sqlite.[schemaVersion].[YYYYMMDD_HHMMSS]
                </code>
              ),
            }
          )}
        </Form.Text>
      </Form.Group>

      <Form.Group>
        <Button
          id="backupDownload"
          variant="secondary"
          type="submit"
          onClick={() => onBackup(true)}
        >
          <FormattedMessage id="actions.download_backup" />
        </Button>
        <Form.Text className="text-muted">
          {intl.formatMessage({ id: "config.tasks.backup_and_download" })}
        </Form.Text>
      </Form.Group>

      {renderPlugins()}

      <hr />

      <h5>{intl.formatMessage({ id: "config.tasks.migrations" })}</h5>

      <Form.Group>
        <Button
          id="migrateHashNaming"
          variant="danger"
          onClick={() => onMigrateHashNaming()}
        >
          <FormattedMessage id="actions.rename_gen_files" />
        </Button>
        <Form.Text className="text-muted">
          {intl.formatMessage({ id: "config.tasks.migrate_hash_files" })}
        </Form.Text>
      </Form.Group>
    </>
  );
};
