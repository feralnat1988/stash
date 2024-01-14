import React, {
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";
import { CriterionType } from "src/models/list-filter/types";
import { SearchField } from "../List/ListFilter";
import { ListFilterModel } from "src/models/list-filter/filter";
import useFocus from "src/utils/focus";
import {
  Criterion,
  CriterionOption,
  CriterionValue,
} from "src/models/list-filter/criteria/criterion";
import { Button, ButtonGroup } from "react-bootstrap";
import { Icon } from "../Shared/Icon";
import { FormattedMessage, useIntl } from "react-intl";
import {
  faChevronLeft,
  faFilter,
  faTimes,
} from "@fortawesome/free-solid-svg-icons";
import { CriterionEditor } from "../List/CriterionEditor";
import { CollapseButton } from "../Shared/CollapseButton";
import cx from "classnames";
import { EditFilterDialog } from "../List/EditFilterDialog";
import { SavedFilterList } from "../List/SavedFilterList";
import { useFilterConfig } from "./util";

const FilterCriteriaList: React.FC<{
  filter: ListFilterModel;
  hiddenOptions: CriterionOption[];
  onRemoveCriterion: (c: Criterion<CriterionValue>) => void;
  onEditCriterion: (c: Criterion<CriterionValue>) => void;
}> = ({ filter, hiddenOptions, onRemoveCriterion, onEditCriterion }) => {
  const intl = useIntl();

  const criteria = useMemo(
    () =>
      filter.criteria.filter((c) => {
        return hiddenOptions.some((h) => h.type === c.criterionOption.type);
      }),
    [filter.criteria, hiddenOptions]
  );

  if (criteria.length === 0) return null;

  function onClickRemoveCriterion(
    criterion: Criterion<CriterionValue>,
    $event: React.MouseEvent<HTMLElement, MouseEvent>
  ) {
    if (!criterion) {
      return;
    }
    onRemoveCriterion(criterion);
    $event.stopPropagation();
  }

  function onClickCriterionTag(criterion: Criterion<CriterionValue>) {
    onEditCriterion(criterion);
  }

  return (
    <div className="filter-criteria-list">
      <ul>
        {criteria.map((c) => {
          return (
            <li className="filter-criteria-list-item" key={c.getId()}>
              <a onClick={() => onClickCriterionTag(c)}>
                <span>{c.getLabel(intl)}</span>
                <Button
                  className="remove-criterion-button"
                  variant="minimal"
                  onClick={($event) => onClickRemoveCriterion(c, $event)}
                >
                  <Icon icon={faTimes} />
                </Button>
              </a>
            </li>
          );
        })}
      </ul>
      <hr />
    </div>
  );
};

interface ICriterionList {
  filter: ListFilterModel;
  currentCriterion?: Criterion<CriterionValue>;
  setCriterion: (c: Criterion<CriterionValue>) => void;
  criterionOptions: CriterionOption[];
  onRemoveCriterion: (c: string) => void;
}

const CriterionOptionList: React.FC<ICriterionList> = ({
  filter,
  currentCriterion,
  setCriterion,
  criterionOptions,
  onRemoveCriterion,
}) => {
  const intl = useIntl();

  const scrolled = useRef(false);

  const type = currentCriterion?.criterionOption.type;

  const criteriaRefs = useMemo(() => {
    const refs: Record<string, React.RefObject<HTMLDivElement>> = {};
    criterionOptions.forEach((c) => {
      refs[c.type] = React.createRef();
    });
    return refs;
  }, [criterionOptions]);

  useEffect(() => {
    // scrolling to the current criterion doesn't work well when the
    // dialog is already open, so limit to when we click on the
    // criterion from the external tags
    if (!scrolled.current && type && criteriaRefs[type]?.current) {
      criteriaRefs[type].current!.scrollIntoView({
        behavior: "smooth",
        block: "start",
      });
      scrolled.current = true;
    }
  }, [currentCriterion, criteriaRefs, type]);

  function getReleventCriterion(t: CriterionType) {
    // find the existing criterion if present
    const existing = filter.criteria.find((c) => c.criterionOption.type === t);
    if (existing) {
      return existing;
    } else {
      const newCriterion = filter.makeCriterion(t);
      return newCriterion;
    }
  }

  function removeClicked(ev: React.MouseEvent, t: string) {
    // needed to prevent the nav item from being selected
    ev.stopPropagation();
    ev.preventDefault();
    onRemoveCriterion(t);
  }

  function renderCard(c: CriterionOption) {
    return (
      <CollapseButton
        text={intl.formatMessage({ id: c.messageID })}
        rightControls={
          <span>
            <Button
              className={cx("remove-criterion-button", {
                invisible: !filter.criteria.some(
                  (cc) => c.type === cc.criterionOption.type
                ),
              })}
              variant="minimal"
              onClick={(e) => removeClicked(e, c.type)}
            >
              <Icon icon={faTimes} />
            </Button>
          </span>
        }
      >
        <CriterionEditor
          criterion={getReleventCriterion(c.type)!}
          setCriterion={setCriterion}
        />
      </CollapseButton>
    );
  }

  return (
    <div className="criterion-list">
      {criterionOptions.map((c) => renderCard(c))}
    </div>
  );
};

export const FilterSidebar: React.FC<{
  filter: ListFilterModel;
  setFilter: (filter: ListFilterModel) => void;
  onHide: () => void;
}> = ({ filter, setFilter, onHide }) => {
  const intl = useIntl();

  const [queryRef, setQueryFocus] = useFocus();

  const [criterion, setCriterion] = useState<Criterion<CriterionValue>>();
  const {
    criterionOptions,
    sidebarOptions,
    hiddenOptions,
    setCriterionOptions,
  } = useFilterConfig(filter.mode);

  const [editingCriterion, setEditingCriterion] = useState<string>();
  const [showEditFilter, setShowEditFilter] = useState(false);

  const { criteria } = filter;

  const searchQueryUpdated = useCallback(
    (value: string) => {
      const newFilter = filter.clone();
      newFilter.searchTerm = value;
      newFilter.currentPage = 1;
      setFilter(newFilter);
    },
    [filter, setFilter]
  );

  const optionSelected = useCallback(
    (option?: CriterionOption) => {
      if (!option) {
        setCriterion(undefined);
        return;
      }

      // find the existing criterion if present
      const existing = criteria.find(
        (c) => c.criterionOption.type === option.type
      );
      if (existing) {
        setCriterion(existing);
      } else {
        const newCriterion = filter.makeCriterion(option.type);
        setCriterion(newCriterion);
      }
    },
    [filter, criteria]
  );

  function removeCriterion(c: Criterion<CriterionValue>) {
    const newFilter = filter.clone();

    const newCriteria = criteria.filter((cc) => {
      return cc.getId() !== c.getId();
    });

    newFilter.criteria = newCriteria;

    setFilter(newFilter);
    if (criterion?.getId() === c.getId()) {
      optionSelected(undefined);
    }
  }

  function removeCriterionString(c: string) {
    const cc = criteria.find((ccc) => ccc.criterionOption.type === c);
    if (cc) {
      removeCriterion(cc);
    }
  }

  function replaceCriterion(c: Criterion<CriterionValue>) {
    const newFilter = filter.clone();

    if (!c.isValid()) {
      // remove from the filter if present
      const newCriteria = criteria.filter((cc) => {
        return cc.criterionOption.type !== c.criterionOption.type;
      });

      newFilter.criteria = newCriteria;
    } else {
      let found = false;

      const newCriteria = criteria.map((cc) => {
        if (cc.criterionOption.type === c.criterionOption.type) {
          found = true;
          return c;
        }

        return cc;
      });

      if (!found) {
        newCriteria.push(c);
      }

      newFilter.criteria = newCriteria;
    }

    setFilter(newFilter);
  }

  function onApplyEditFilter(f?: ListFilterModel) {
    setShowEditFilter(false);
    setEditingCriterion(undefined);

    if (!f) return;
    setFilter(f);
  }

  return (
    <div className="filter-sidebar">
      <ButtonGroup className="search-field-group">
        <SearchField
          searchTerm={filter.searchTerm}
          setSearchTerm={searchQueryUpdated}
          queryRef={queryRef}
          setQueryFocus={setQueryFocus}
        />
        <Button
          onClick={() => onHide()}
          variant="secondary"
          className="collapse-filter-button"
        >
          <Icon icon={faChevronLeft} />
        </Button>
      </ButtonGroup>
      <div>
        <FilterCriteriaList
          filter={filter}
          hiddenOptions={hiddenOptions}
          onRemoveCriterion={(c) =>
            removeCriterionString(c.criterionOption.type)
          }
          onEditCriterion={(c) => setEditingCriterion(c.criterionOption.type)}
        />
      </div>
      <div className="saved-filters">
        <CollapseButton
          text={intl.formatMessage({ id: "search_filter.saved_filters" })}
        >
          <SavedFilterList filter={filter} onSetFilter={setFilter} />
        </CollapseButton>
      </div>
      <CriterionOptionList
        filter={filter}
        currentCriterion={criterion}
        setCriterion={replaceCriterion}
        criterionOptions={sidebarOptions}
        onRemoveCriterion={(c) => removeCriterionString(c)}
      />
      <div>
        <Button
          variant="secondary"
          className="edit-filter-button"
          onClick={() => setShowEditFilter(true)}
        >
          <Icon icon={faFilter} />{" "}
          <FormattedMessage id="search_filter.edit_filter" />
        </Button>
      </div>
      {(showEditFilter || editingCriterion) && (
        <EditFilterDialog
          filter={filter}
          criterionOptions={criterionOptions}
          setCriterionOptions={(o) => setCriterionOptions(o)}
          onClose={onApplyEditFilter}
          editingCriterion={editingCriterion}
        />
      )}
    </div>
  );
};
