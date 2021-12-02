import React, { useState } from "react";
import { Button, Collapse, Form, Modal } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";
import { PropsWithChildren } from "react-router/node_modules/@types/react";
import { Icon } from "../Shared";
import { StringListInput } from "../Shared/StringListInput";

interface ISetting {
  id?: string;
  heading?: string;
  headingID?: string;
  subHeadingID?: string;
  subHeading?: React.ReactNode;
}

export const Setting: React.FC<PropsWithChildren<ISetting>> = ({
  id,
  heading,
  headingID,
  subHeadingID,
  subHeading,
  children,
}) => {
  const intl = useIntl();

  function renderHeading() {
    if (headingID) {
      return intl.formatMessage({ id: headingID });
    }
    return heading;
  }

  function renderSubHeading() {
    if (subHeadingID) {
      return (
        <div className="sub-heading">
          {intl.formatMessage({ id: subHeadingID })}
        </div>
      );
    }
    if (subHeading) {
      return <div className="sub-heading">{subHeading}</div>;
    }
  }

  return (
    <div className="setting" id={id}>
      <div>
        <h3>{renderHeading()}</h3>
        {renderSubHeading()}
      </div>
      <div>{children}</div>
    </div>
  );
};

interface ISettingGroup {
  settingProps: ISetting;
  topLevel?: JSX.Element;
  collapsible?: boolean;
  collapsedDefault?: boolean;
}

export const SettingGroup: React.FC<PropsWithChildren<ISettingGroup>> = ({
  settingProps,
  topLevel,
  collapsible,
  collapsedDefault,
  children,
}) => {
  const [open, setOpen] = useState(!collapsedDefault);

  function renderCollapseButton() {
    if (!collapsible) return;

    return (
      <Button
        className="collapse-button"
        variant="minimal"
        onClick={() => setOpen(!open)}
      >
        <Icon
          className="fa-fw"
          icon={open ? "chevron-circle-down" : "chevron-circle-left"}
        />
      </Button>
    );
  }

  return (
    <div className="setting-group">
      <Setting {...settingProps}>
        {topLevel}
        {renderCollapseButton()}
      </Setting>
      <Collapse in={open}>
        <div className="collapsible-section">{children}</div>
      </Collapse>
    </div>
  );
};

interface IBooleanSetting extends ISetting {
  id: string;
  disabled?: boolean;
  checked?: boolean;
  onChange: (v: boolean) => void;
}

export const BooleanSetting: React.FC<IBooleanSetting> = ({
  id,
  headingID,
  subHeadingID,
  disabled,
  checked,
  onChange,
}) => {
  return (
    <Setting headingID={headingID} subHeadingID={subHeadingID}>
      <Form.Switch
        id={id}
        disabled={disabled}
        checked={checked}
        onChange={() => onChange(!checked)}
      />
    </Setting>
  );
};

interface ISelectSetting extends ISetting {
  value?: string | number | string[] | undefined;
  onChange: (v: string) => void;
}

export const SelectSetting: React.FC<PropsWithChildren<ISelectSetting>> = ({
  id,
  headingID,
  subHeadingID,
  value,
  children,
  onChange,
}) => {
  return (
    <Setting headingID={headingID} subHeadingID={subHeadingID} id={id}>
      <Form.Control
        className="input-control"
        as="select"
        value={value}
        onChange={(e) => onChange(e.currentTarget.value)}
      >
        {children}
      </Form.Control>
    </Setting>
  );
};

interface IDialogSetting<T> extends ISetting {
  buttonTextID?: string;
  value?: T;
  renderValue?: (v: T | undefined) => JSX.Element;
  onChange: () => void;
}

export const ChangeButtonSetting = <T extends {}>(props: IDialogSetting<T>) => {
  const {
    id,
    headingID,
    subHeadingID,
    value,
    onChange,
    renderValue,
    buttonTextID,
  } = props;
  const intl = useIntl();

  return (
    <div className="setting" id={id}>
      <div>
        <h3>{intl.formatMessage({ id: headingID })}</h3>

        <div className="value">
          {renderValue ? renderValue(value) : undefined}
        </div>

        {subHeadingID ? (
          <div className="sub-heading">
            {intl.formatMessage({ id: subHeadingID })}
          </div>
        ) : undefined}
      </div>
      <div>
        <Button onClick={() => onChange()}>
          <FormattedMessage id={buttonTextID ?? "actions.edit"} />
        </Button>
      </div>
    </div>
  );
};

export interface ISettingModal<T> {
  heading?: string;
  headingID?: string;
  subHeadingID?: string;
  value: T | undefined;
  close: (v?: T) => void;
  renderField: (value: T | undefined, setValue: (v?: T) => void) => JSX.Element;
}

export const SettingModal = <T extends {}>(props: ISettingModal<T>) => {
  const { heading, headingID, subHeadingID, value, close, renderField } = props;

  const intl = useIntl();
  const [currentValue, setCurrentValue] = useState<T | undefined>(value);

  return (
    <Modal show onHide={() => close()} id="setting-dialog">
      <Form
        onSubmit={(e) => {
          close(currentValue);
          e.preventDefault();
        }}
      >
        <Modal.Header closeButton>
          {headingID ? <FormattedMessage id={headingID} /> : heading}
        </Modal.Header>
        <Modal.Body>
          {renderField(currentValue, setCurrentValue)}
          {subHeadingID ? (
            <div className="sub-heading">
              {intl.formatMessage({ id: subHeadingID })}
            </div>
          ) : undefined}
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={() => close()}>
            Cancel
          </Button>
          <Button
            type="submit"
            variant="primary"
            onClick={() => close(currentValue)}
          >
            Confirm
          </Button>
        </Modal.Footer>
      </Form>
    </Modal>
  );
};

interface IModalSetting<T> extends ISetting {
  value: T | undefined;
  buttonTextID?: string;
  onChange: (v: T) => void;
  renderField: (value: T | undefined, setValue: (v?: T) => void) => JSX.Element;
  renderValue?: (v: T | undefined) => JSX.Element;
}

export const ModalSetting = <T extends {}>(props: IModalSetting<T>) => {
  const {
    id,
    value,
    headingID,
    subHeadingID,
    onChange,
    renderField,
    renderValue,
    buttonTextID,
  } = props;
  const [showModal, setShowModal] = useState(false);

  return (
    <>
      {showModal ? (
        <SettingModal<T>
          headingID={headingID}
          subHeadingID={subHeadingID}
          value={value}
          renderField={renderField}
          close={(v) => {
            if (v !== undefined) onChange(v);
            setShowModal(false);
          }}
        />
      ) : undefined}

      <ChangeButtonSetting<T>
        id={id}
        buttonTextID={buttonTextID}
        headingID={headingID}
        subHeadingID={subHeadingID}
        value={value}
        onChange={() => setShowModal(true)}
        renderValue={renderValue}
      />
    </>
  );
};

interface IStringSetting extends ISetting {
  value: string | undefined;
  onChange: (v: string) => void;
}

export const StringSetting: React.FC<IStringSetting> = (props) => {
  return (
    <ModalSetting<string>
      {...props}
      renderField={(value, setValue) => (
        <Form.Control
          className="text-input"
          value={value}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            setValue(e.currentTarget.value)
          }
        />
      )}
      renderValue={(value) => <span>{value}</span>}
    />
  );
};

interface INumberSetting extends ISetting {
  value: number | undefined;
  onChange: (v: number) => void;
}

export const NumberSetting: React.FC<INumberSetting> = (props) => {
  return (
    <ModalSetting<number>
      {...props}
      renderField={(value, setValue) => (
        <Form.Control
          className="text-input"
          type="number"
          value={value}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            setValue(Number.parseInt(e.currentTarget.value || "0", 10))
          }
        />
      )}
      renderValue={(value) => <span>{value}</span>}
    />
  );
};

interface IStringListSetting extends ISetting {
  value: string[] | undefined;
  defaultNewValue?: string;
  onChange: (v: string[]) => void;
}

export const StringListSetting: React.FC<IStringListSetting> = (props) => {
  return (
    <ModalSetting<string[]>
      {...props}
      renderField={(value, setValue) => (
        <StringListInput
          value={value ?? []}
          setValue={setValue}
          defaultNewValue={props.defaultNewValue}
        />
      )}
      renderValue={(value) => (
        <div>
          {value?.map((v, i) => (
            // eslint-disable-next-line react/no-array-index-key
            <div key={i}>{v}</div>
          ))}
        </div>
      )}
    />
  );
};
