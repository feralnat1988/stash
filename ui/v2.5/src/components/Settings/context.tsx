import { ApolloError } from "@apollo/client/errors";
import { debounce } from "lodash";
import React, { useState, useEffect, useMemo } from "react";
import { useIntl } from "react-intl";
import * as GQL from "src/core/generated-graphql";
import {
  useConfiguration,
  useConfigureDefaults,
  useConfigureGeneral,
  useConfigureInterface,
} from "src/core/StashService";
import { useToast } from "src/hooks";
import { withoutTypename } from "src/utils";

export interface ISettingsContextState {
  loading: boolean;
  error: ApolloError | undefined;
  general: GQL.ConfigGeneralInput;
  interface: GQL.ConfigInterfaceInput;
  defaults: GQL.ConfigDefaultSettingsInput;

  // apikey isn't directly settable, so expose it here
  apiKey: string;

  saveGeneral: (input: Partial<GQL.ConfigGeneralInput>) => void;
  saveInterface: (input: Partial<GQL.ConfigInterfaceInput>) => void;
  saveDefaults: (input: Partial<GQL.ConfigDefaultSettingsInput>) => void;
}

export const SettingStateContext = React.createContext<ISettingsContextState>({
  loading: false,
  error: undefined,
  general: {},
  interface: {},
  defaults: {},
  apiKey: "",
  saveGeneral: () => {},
  saveInterface: () => {},
  saveDefaults: () => {},
});

export const SettingsContext: React.FC = ({ children }) => {
  const intl = useIntl();
  const Toast = useToast();

  const { data, error, loading } = useConfiguration();

  const [general, setGeneral] = useState<GQL.ConfigGeneralInput>({});
  const [pendingGeneral, setPendingGeneral] = useState<
    GQL.ConfigGeneralInput | undefined
  >();
  const [updateGeneralConfig] = useConfigureGeneral();

  const [iface, setIface] = useState<GQL.ConfigInterfaceInput>({});
  const [pendingInterface, setPendingInterface] = useState<
    GQL.ConfigInterfaceInput | undefined
  >();
  const [updateInterfaceConfig] = useConfigureInterface();

  const [defaults, setDefaults] = useState<GQL.ConfigDefaultSettingsInput>({});
  const [pendingDefaults, setPendingDefaults] = useState<
    GQL.ConfigDefaultSettingsInput | undefined
  >();
  const [updateDefaultsConfig] = useConfigureDefaults();

  const [apiKey, setApiKey] = useState("");

  useEffect(() => {
    if (!data?.configuration || error) return;

    setGeneral({ ...withoutTypename(data.configuration.general) });
    setIface({ ...withoutTypename(data.configuration.interface) });
    setDefaults({ ...withoutTypename(data.configuration.defaults) });
    setApiKey(data.configuration.general.apiKey);
  }, [data, error]);

  // saves the configuration if no further changes are made after a half second
  const saveGeneralConfig = useMemo(
    () =>
      debounce(async (input: GQL.ConfigGeneralInput) => {
        try {
          await updateGeneralConfig({
            variables: {
              input,
            },
          });

          // TODO - use different notification method
          Toast.success({
            content: intl.formatMessage(
              { id: "toast.updated_entity" },
              {
                entity: intl
                  .formatMessage({ id: "configuration" })
                  .toLocaleLowerCase(),
              }
            ),
          });
          setPendingGeneral(undefined);
        } catch (e) {
          Toast.error(e);
        }
      }, 500),
    [Toast, intl, updateGeneralConfig]
  );

  useEffect(() => {
    if (!pendingGeneral) {
      return;
    }

    saveGeneralConfig(pendingGeneral);
  }, [pendingGeneral, saveGeneralConfig]);

  function saveGeneral(input: Partial<GQL.ConfigGeneralInput>) {
    if (!general) {
      return;
    }

    setGeneral({
      ...general,
      ...input,
    });

    setPendingGeneral((current) => {
      if (!current) {
        return input;
      }
      return {
        ...current,
        ...input,
      };
    });
  }

  // saves the configuration if no further changes are made after a half second
  const saveInterfaceConfig = useMemo(
    () =>
      debounce(async (input: GQL.ConfigInterfaceInput) => {
        try {
          await updateInterfaceConfig({
            variables: {
              input,
            },
          });

          // TODO - use different notification method
          Toast.success({
            content: intl.formatMessage(
              { id: "toast.updated_entity" },
              {
                entity: intl
                  .formatMessage({ id: "configuration" })
                  .toLocaleLowerCase(),
              }
            ),
          });
          setPendingInterface(undefined);
        } catch (e) {
          Toast.error(e);
        }
      }, 500),
    [Toast, intl, updateInterfaceConfig]
  );

  useEffect(() => {
    if (!pendingInterface) {
      return;
    }

    saveInterfaceConfig(pendingInterface);
  }, [pendingInterface, saveInterfaceConfig]);

  function saveInterface(input: Partial<GQL.ConfigInterfaceInput>) {
    if (!iface) {
      return;
    }

    setIface({
      ...iface,
      ...input,
    });

    setPendingInterface((current) => {
      if (!current) {
        return input;
      }
      return {
        ...current,
        ...input,
      };
    });
  }

  // saves the configuration if no further changes are made after a half second
  const saveDefaultsConfig = useMemo(
    () =>
      debounce(async (input: GQL.ConfigDefaultSettingsInput) => {
        try {
          await updateDefaultsConfig({
            variables: {
              input,
            },
          });

          // TODO - use different notification method
          Toast.success({
            content: intl.formatMessage(
              { id: "toast.updated_entity" },
              {
                entity: intl
                  .formatMessage({ id: "configuration" })
                  .toLocaleLowerCase(),
              }
            ),
          });
          setPendingDefaults(undefined);
        } catch (e) {
          Toast.error(e);
        }
      }, 500),
    [Toast, intl, updateDefaultsConfig]
  );

  useEffect(() => {
    if (!pendingDefaults) {
      return;
    }

    saveDefaultsConfig(pendingDefaults);
  }, [pendingDefaults, saveDefaultsConfig]);

  function saveDefaults(input: Partial<GQL.ConfigDefaultSettingsInput>) {
    if (!defaults) {
      return;
    }

    setDefaults({
      ...defaults,
      ...input,
    });

    setPendingDefaults((current) => {
      if (!current) {
        return input;
      }
      return {
        ...current,
        ...input,
      };
    });
  }

  return (
    <SettingStateContext.Provider
      value={{
        loading,
        error,
        apiKey,
        general,
        interface: iface,
        defaults,
        saveGeneral,
        saveInterface,
        saveDefaults,
      }}
    >
      {children}
    </SettingStateContext.Provider>
  );
};
