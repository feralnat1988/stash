import React from "react";
import { useIntl } from "react-intl";
import { Button, Form } from "react-bootstrap";
import { mutateRunPluginTask, usePlugins } from "src/core/StashService";
import { useToast } from "src/hooks";
import * as GQL from "src/core/generated-graphql";
import { Task } from "./Task";

type Plugin = Pick<GQL.Plugin, "id">;
type PluginTask = Pick<GQL.PluginTask, "name" | "description">;

export const PluginTasks: React.FC = () => {
  const intl = useIntl();
  const Toast = useToast();

  const plugins = usePlugins();

  function renderPlugins() {
    if (!plugins.data || !plugins.data.plugins) {
      return;
    }

    const taskPlugins = plugins.data.plugins.filter(
      (p) => p.tasks && p.tasks.length > 0
    );

    return (
      <Form.Group>
        <h5>{intl.formatMessage({ id: "config.tasks.plugin_tasks" })}</h5>
        {taskPlugins.map((o) => {
          return (
            <Form.Group key={`${o.id}`}>
              <h6>{o.name}</h6>
              <div className="task-group">
                {renderPluginTasks(o, o.tasks ?? [])}
              </div>
            </Form.Group>
          );
        })}
      </Form.Group>
    );
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

  async function onPluginTaskClicked(plugin: Plugin, operation: PluginTask) {
    await mutateRunPluginTask(plugin.id, operation.name);
    Toast.success({
      content: intl.formatMessage(
        { id: "config.tasks.added_job_to_queue" },
        { operation_name: operation.name }
      ),
    });
  }

  return <Form.Group>{renderPlugins()}</Form.Group>;
};
