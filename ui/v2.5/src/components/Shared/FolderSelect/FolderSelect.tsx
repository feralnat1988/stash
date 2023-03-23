import React, { useEffect, useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { Button, InputGroup, Form, Collapse } from "react-bootstrap";
import { Icon } from "../Icon";
import { LoadingIndicator } from "../LoadingIndicator";
import { useDirectory } from "src/core/StashService";
import { faEllipsis, faTimes } from "@fortawesome/free-solid-svg-icons";
import { useDebouncedSetState } from "src/hooks/debounce";

interface IProps {
  currentDirectory: string;
  setCurrentDirectory: (value: string) => void;
  defaultDirectories?: string[];
  appendButton?: JSX.Element;
  collapsible?: boolean;
}

export const FolderSelect: React.FC<IProps> = ({
  currentDirectory,
  setCurrentDirectory,
  defaultDirectories,
  appendButton,
  collapsible = false,
}) => {
  const [showBrowser, setShowBrowser] = React.useState(false);
  const [directory, setDirectory] = useState(currentDirectory);
  const { data, error, loading } = useDirectory(directory);
  const intl = useIntl();

  const selectableDirectories: string[] = currentDirectory
    ? data?.directory.directories ?? defaultDirectories ?? []
    : defaultDirectories ?? [];

  const debouncedSetDirectory = useDebouncedSetState(setDirectory, 250);

  useEffect(() => {
    if (currentDirectory !== directory) {
      debouncedSetDirectory(currentDirectory);
    }
  }, [currentDirectory, directory, debouncedSetDirectory]);

  function setInstant(value: string) {
    setCurrentDirectory(value);
    setDirectory(value);
  }

  function setDebounced(value: string) {
    setCurrentDirectory(value);
    debouncedSetDirectory(value);
  }

  function goUp() {
    if (defaultDirectories?.includes(currentDirectory)) {
      setInstant("");
    } else if (data?.directory.parent) {
      setInstant(data.directory.parent);
    }
  }

  const topDirectory =
    currentDirectory && data?.directory?.parent ? (
      <li className="folder-list-parent folder-list-item">
        <Button variant="link" onClick={() => goUp()}>
          <FormattedMessage id="setup.folder.up_dir" />
        </Button>
      </li>
    ) : null;

  return (
    <>
      <InputGroup>
        <Form.Control
          className="btn-secondary"
          placeholder={intl.formatMessage({ id: "setup.folder.file_path" })}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
            setDebounced(e.currentTarget.value);
          }}
          value={currentDirectory}
          spellCheck={false}
        />
        {appendButton ? (
          <InputGroup.Append>{appendButton}</InputGroup.Append>
        ) : undefined}
        {collapsible ? (
          <InputGroup.Append>
            <Button
              variant="secondary"
              onClick={() => setShowBrowser(!showBrowser)}
            >
              <Icon icon={faEllipsis} />
            </Button>
          </InputGroup.Append>
        ) : undefined}
        {!data || !data.directory || loading ? (
          <InputGroup.Append className="align-self-center">
            {loading ? (
              <LoadingIndicator inline small message="" />
            ) : (
              <Icon icon={faTimes} color="red" className="ml-3" />
            )}
          </InputGroup.Append>
        ) : undefined}
      </InputGroup>
      {error !== undefined && (
        <h5 className="mt-4 text-break">Error: {error.message}</h5>
      )}
      <Collapse in={!collapsible || showBrowser}>
        <ul className="folder-list">
          {topDirectory}
          {selectableDirectories.map((path) => {
            return (
              <li key={path} className="folder-list-item">
                <Button variant="link" onClick={() => setInstant(path)}>
                  {path}
                </Button>
              </li>
            );
          })}
        </ul>
      </Collapse>
    </>
  );
};
