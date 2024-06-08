import React, { useCallback, useEffect } from "react";
import { Pagination } from "../List/Pagination";
import { DisplayModeSelect } from "../List/ListViewOptions";
import { DisplayMode } from "src/models/list-filter/types";
import {
  PageSizeSelect,
  SavedFilterSelect,
  SearchField,
  SortBySelect,
} from "../List/ListFilter";
import { SortDirectionEnum } from "src/core/generated-graphql";
import { getFilterOptions } from "src/models/list-filter/factory";
import { ListFilterModel } from "src/models/list-filter/filter";
import { Button, ButtonGroup, Dropdown } from "react-bootstrap";
import { Icon } from "../Shared/Icon";
import { FormattedMessage } from "react-intl";
import {
  faChevronRight,
  faEllipsisH,
  faTimes,
} from "@fortawesome/free-solid-svg-icons";
import useFocus from "src/utils/focus";
import { FilterButton } from "./Filters/FilterButton";
import Mousetrap from "mousetrap";

interface IDefaultListHeaderProps {
  filter: ListFilterModel;
  setFilter: (filter: ListFilterModel) => void;
  showFilterDialog?: () => void;
  totalItems: number;
  actionButtons?: React.ReactNode;
  sidebarCollapsed?: boolean;
  showSidebar?: () => void;
}

const DefaultListHeader: React.FC<IDefaultListHeaderProps> = ({
  filter,
  setFilter,
  showFilterDialog,
  totalItems,
  actionButtons,
  sidebarCollapsed,
  showSidebar,
}) => {
  const [queryRef, setQueryFocus] = useFocus();

  const filterOptions = getFilterOptions(filter.mode);

  function onChangeZoom(newZoomIndex: number) {
    const newFilter = filter.clone();
    newFilter.zoomIndex = newZoomIndex;
    setFilter(newFilter);
  }

  const searchQueryUpdated = useCallback(
    (value: string) => {
      const newFilter = filter.clone();
      newFilter.searchTerm = value;
      newFilter.currentPage = 1;
      setFilter(newFilter);
    },
    [filter, setFilter]
  );

  function onChangeDisplayMode(displayMode: DisplayMode) {
    const newFilter = filter.clone();
    newFilter.displayMode = displayMode;
    setFilter(newFilter);
  }

  function onChangePageSize(val: number) {
    const newFilter = filter.clone();
    newFilter.itemsPerPage = val;
    newFilter.currentPage = 1;
    setFilter(newFilter);
  }

  function onChangeSortDirection(dir: SortDirectionEnum) {
    const newFilter = filter.clone();
    newFilter.sortDirection = dir;
    setFilter(newFilter);
  }

  function onChangeSortBy(eventKey: string | null) {
    const newFilter = filter.clone();
    newFilter.sortBy = eventKey ?? undefined;
    newFilter.currentPage = 1;
    setFilter(newFilter);
  }

  const onReshuffleRandomSort = useCallback(() => {
    const newFilter = filter.clone();
    newFilter.currentPage = 1;
    newFilter.randomSeed = -1;
    setFilter(newFilter);
  }, [filter, setFilter]);

  const onChangePage = useCallback(
    (page: number) => {
      const newFilter = filter.clone();
      newFilter.currentPage = page;
      setFilter(newFilter);

      // if the current page has a detail-header, then
      // scroll up relative to that rather than 0, 0
      const detailHeader = document.querySelector(".detail-header");
      if (detailHeader) {
        window.scrollTo(0, detailHeader.scrollHeight - 50);
      } else {
        window.scrollTo(0, 0);
      }
    },
    [filter, setFilter]
  );

  useEffect(() => {
    Mousetrap.bind("/", (e) => {
      setQueryFocus();
      e.preventDefault();
    });

    return () => {
      Mousetrap.unbind("/");
    };
  }, [setQueryFocus]);

  useEffect(() => {
    Mousetrap.bind("r", () => onReshuffleRandomSort());

    return () => {
      Mousetrap.unbind("r");
    };
  }, [onReshuffleRandomSort]);

  const zoomSelectProps =
    filter.displayMode === DisplayMode.Grid
      ? {
          minZoom: 0,
          maxZoom: 3,
          zoomIndex: filter.zoomIndex,
          onChangeZoom,
        }
      : undefined;

  return (
    <div className="list-header">
      <div className="list-header-left">
        {sidebarCollapsed && showSidebar && (
          <Button
            className="show-siderbar-button-xs"
            variant="secondary"
            onClick={() => showSidebar()}
          >
            <Icon size="sm" icon={faChevronRight} />
          </Button>
        )}
        <SearchField
          searchTerm={filter.searchTerm}
          setSearchTerm={searchQueryUpdated}
          queryRef={queryRef}
          setQueryFocus={setQueryFocus}
        />
        {sidebarCollapsed && (
          <ButtonGroup className="filter-button-group">
            <SavedFilterSelect filter={filter} onFilterUpdate={setFilter} />
            {showFilterDialog && (
              <FilterButton
                filter={filter}
                onClick={() => showFilterDialog()}
              />
            )}
          </ButtonGroup>
        )}
      </div>
      <div className="list-header-center">
        <Pagination
          currentPage={filter.currentPage}
          itemsPerPage={filter.itemsPerPage}
          totalItems={totalItems}
          onChangePage={onChangePage}
        />
      </div>
      <div className="list-header-right">
        {actionButtons && (
          <Dropdown>
            <Dropdown.Toggle variant="secondary" id="more-menu">
              <Icon icon={faEllipsisH} />
            </Dropdown.Toggle>
            <Dropdown.Menu className="bg-secondary text-white">
              {actionButtons}
            </Dropdown.Menu>
          </Dropdown>
        )}
        <SortBySelect
          sortBy={filter.sortBy}
          direction={filter.sortDirection}
          options={filterOptions.sortByOptions}
          setSortBy={onChangeSortBy}
          setDirection={onChangeSortDirection}
          onReshuffleRandomSort={onReshuffleRandomSort}
        />
        <PageSizeSelect
          pageSize={filter.itemsPerPage}
          setPageSize={onChangePageSize}
        />
        <DisplayModeSelect
          displayMode={filter.displayMode}
          displayModeOptions={filterOptions.displayModeOptions}
          onSetDisplayMode={onChangeDisplayMode}
          zoomSelectProps={zoomSelectProps}
        />
      </div>
    </div>
  );
};

interface ISelectedListHeader {
  selectedIds: Set<string>;
  onSelectAll: () => void;
  onSelectNone: () => void;
  selectedButtons?: (selectedIds: Set<string>) => React.ReactNode;
}

export const SelectedListHeader: React.FC<ISelectedListHeader> = ({
  selectedIds,
  onSelectAll,
  onSelectNone,
  selectedButtons = () => null,
}) => {
  return (
    <div className="list-header selected-list-header">
      <div className="list-header-left">
        <span>{selectedIds.size} items selected</span>
        <Button variant="link" onClick={() => onSelectAll()}>
          <FormattedMessage id="actions.select_all" />
        </Button>
      </div>
      <div className="list-header-center">{selectedButtons(selectedIds)}</div>
      <div className="list-header-right">
        <Button className="minimal select-none" onClick={() => onSelectNone()}>
          <Icon icon={faTimes} />
        </Button>
      </div>
    </div>
  );
};

export interface IListHeaderProps
  extends IDefaultListHeaderProps,
    ISelectedListHeader {}

export const ListHeader: React.FC<IListHeaderProps> = (props) => {
  if (props.selectedIds.size === 0) {
    return <DefaultListHeader {...props} />;
  } else {
    return <SelectedListHeader {...props} />;
  }
};
