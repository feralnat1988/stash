import React, {
  useState,
  useCallback,
  useMemo,
  useEffect,
  useContext,
} from "react";
import * as GQL from "src/core/generated-graphql";
import { getFilterOptions } from "src/models/list-filter/factory";
import { CriterionOption } from "src/models/list-filter/criteria/criterion";
import { ListFilterModel } from "src/models/list-filter/filter";
import { useHistory, useLocation } from "react-router-dom";
import isEqual from "lodash-es/isEqual";
import { useConfigureUI, useFindDefaultFilter } from "src/core/StashService";
import { ConfigurationContext } from "src/hooks/Config";
import { IUIConfig } from "src/core/config";
import { useMemoOnce } from "src/hooks/state";
import { useInterfaceLocalForage } from "src/hooks/LocalForage";
import Mousetrap from "mousetrap";
import ScreenUtils from "src/utils/screen";

export interface ICriterionOption {
  option: CriterionOption;
  showInSidebar: boolean;
}

export function useFilterConfig(mode: GQL.FilterMode) {
  const { configuration } = useContext(ConfigurationContext);
  const [saveUI] = useConfigureUI();

  const ui = (configuration?.ui ?? {}) as IUIConfig;

  const savedOrder: string[] = useMemo(
    () => ui.criterionOrder?.[mode.toLowerCase()] ?? [],
    [mode, ui.criterionOrder]
  );

  const savedSidebar: string[] | undefined =
    ui.sidebarCriteria?.[mode.toLocaleLowerCase()];

  const defaultOptions = useMemo(() => {
    const options = getFilterOptions(mode);

    return options.criterionOptions.map((o) => {
      return {
        option: o,
        showInSidebar: !options.defaultHiddenOptions.some(
          (c) => c.type === o.type
        ),
      } as ICriterionOption;
    });
  }, [mode]);

  const [criterionOptions, setCriterionOptionsState] = useState(defaultOptions);

  useEffect(() => {
    const newOrder: ICriterionOption[] = [];
    savedOrder.forEach((o) => {
      const option = defaultOptions.find((d) => d.option.type === o);
      if (option) {
        newOrder.push({ ...option });
      }
    });

    // insert any missing options at the index they would be in the default order
    defaultOptions.forEach((o, i) => {
      if (!newOrder.some((n) => n.option.type === o.option.type)) {
        newOrder.splice(i, 0, { ...o });
      }
    });

    // override sidebar options
    if (savedSidebar) {
      newOrder.forEach((o) => {
        o.showInSidebar = savedSidebar.includes(o.option.type);
      });
    }

    setCriterionOptionsState(newOrder);
  }, [defaultOptions, savedOrder, savedSidebar]);

  function saveCriterionOptions(newOptions: ICriterionOption[]) {
    const criteriaOrder = newOptions.map((o) => o.option.type);
    const sidebarCriteria = newOptions
      .filter((o) => o.showInSidebar)
      .map((o) => o.option.type);

    saveUI({
      variables: {
        partial: {
          criterionOrder: {
            [mode.toLowerCase()]: criteriaOrder,
          },
          sidebarCriteria: {
            [mode.toLowerCase()]: sidebarCriteria,
          },
        },
      },
    });
  }

  function setCriterionOptions(newOptions: ICriterionOption[]) {
    setCriterionOptionsState(newOptions);
    saveCriterionOptions(newOptions);
  }

  const sidebarOptions = useMemo(
    () => criterionOptions.filter((o) => o.showInSidebar).map((o) => o.option),
    [criterionOptions]
  );
  const hiddenOptions = useMemo(
    () => criterionOptions.filter((o) => !o.showInSidebar).map((o) => o.option),
    [criterionOptions]
  );

  return {
    criterionOptions,
    sidebarOptions,
    hiddenOptions,
    setCriterionOptions,
  };
}

export function useFilterURL(
  filter: ListFilterModel,
  setFilter: React.Dispatch<React.SetStateAction<ListFilterModel>>,
  defaultFilter: ListFilterModel | undefined
) {
  const history = useHistory();
  const location = useLocation();

  // this hook causes the initial render to update the URL, losing
  // the existing URL params.
  // useEffect(() => {
  //   const newParams = filter.makeQueryParameters();
  //   history.replace({ ...history.location, search: newParams });
  // }, [filter, history]);

  // when the filter changes, update the URL
  const updateFilter = useCallback(
    (newFilter: ListFilterModel) => {
      const newParams = newFilter.makeQueryParameters();
      history.replace({ ...history.location, search: newParams });
    },
    [history]
  );

  // This hook runs on every page location change (ie navigation),
  // and updates the filter accordingly.
  useEffect(() => {
    // re-init to load default filter on empty new query params
    if (!location.search) {
      if (defaultFilter) updateFilter(defaultFilter.clone());
      return;
    }

    // the query has changed, update filter if necessary
    setFilter((prevFilter) => {
      let newFilter = prevFilter.empty();
      newFilter.configureFromQueryString(location.search);
      if (!isEqual(newFilter, prevFilter)) {
        return newFilter;
      } else {
        return prevFilter;
      }
    });
  }, [location.search, defaultFilter, setFilter, updateFilter]);

  return { setFilter: updateFilter };
}

// returns true if the filter has changed in a way that impacts the total count
function totalCountImpacted(
  oldFilter: ListFilterModel,
  newFilter: ListFilterModel
) {
  return (
    oldFilter.criteria.length !== newFilter.criteria.length ||
    oldFilter.criteria.some((c) => {
      const newCriterion = newFilter.criteria.find(
        (nc) => nc.getId() === c.getId()
      );
      return !newCriterion || !isEqual(c, newCriterion);
    })
  );
}

// this hook caches the total count of results, and only updates it when the filter changes
export function useResultCount(
  filter: ListFilterModel,
  loading: boolean,
  count: number
) {
  const [resultCount, setResultCount] = useState(count);
  const [lastFilter, setLastFilter] = useState(filter);

  // if we are only changing the page or sort, don't update the result count
  useEffect(() => {
    if (!loading) {
      setResultCount(count);
    } else {
      if (totalCountImpacted(lastFilter, filter)) {
        setResultCount(count);
      }
    }

    setLastFilter(filter);
  }, [loading, filter, count, lastFilter]);

  return resultCount;
}

export function useDefaultFilter(mode: GQL.FilterMode) {
  const emptyFilter = useMemo(() => new ListFilterModel(mode), [mode]);

  const { data, loading } = useFindDefaultFilter(mode);

  const defaultFilter = useMemo(() => {
    if (data?.findDefaultFilter) {
      const newFilter = emptyFilter.clone();

      newFilter.currentPage = 1;
      try {
        newFilter.configureFromSavedFilter(data.findDefaultFilter);
      } catch (err) {
        console.log(err);
        // ignore
      }
      // #1507 - reset random seed when loaded
      newFilter.randomSeed = -1;
      return newFilter;
    }
  }, [data?.findDefaultFilter, emptyFilter]);

  const retFilter = loading ? undefined : defaultFilter ?? emptyFilter;

  return { defaultFilter: retFilter, loading };
}

export function useInitialFilter(mode: GQL.FilterMode) {
  const { defaultFilter } = useDefaultFilter(mode);

  // load the default filter on first render
  const initialFilter = useMemoOnce(() => {
    if (!defaultFilter) return [undefined, false];

    if (!location.search) {
      return [defaultFilter, true];
    }

    const newFilter = new ListFilterModel(mode);
    newFilter.configureFromQueryString(location.search);
    return [newFilter, true];
  }, [defaultFilter, location.search]);

  return initialFilter;
}

export function useLocalFilterState(pageView: string, mode: GQL.FilterMode) {
  const [interfaceData] = useInterfaceLocalForage();

  const { loading } = interfaceData;

  const localFilterState = useMemo(() => {
    if (!pageView) return null;
    if (loading) {
      return undefined;
    }

    const { data: existing } = interfaceData;

    return existing?.queryConfig?.[pageView] ?? null;
  }, [interfaceData, loading, pageView]);

  const localState = useMemoOnce<
    | {
        filter: ListFilterModel | undefined;
        sidebarCollapsed: boolean;
      }
    | undefined
  >(() => {
    if (loading) return [undefined, false];

    if (localFilterState) {
      // TODO - set the filter state from local storage
      const storedQuery = localFilterState?.filter;

      const newFilter = new ListFilterModel(mode);
      newFilter.configureFromQueryString(storedQuery);

      return [
        {
          filter: newFilter,
          sidebarCollapsed: localFilterState.sidebarCollapsed,
        },
        true,
      ];
    }

    return [{ filter: undefined, sidebarCollapsed: false }, true];
  }, [localFilterState, mode, loading]);

  return localState;
}

export function useSaveLocalFilterState(
  pageView: string | undefined,
  filter: ListFilterModel,
  sidebarCollapsed: boolean
) {
  const [, setInterfaceState] = useInterfaceLocalForage();

  const setLocalFilterState = useCallback(
    (updatedFilter: ListFilterModel, updatedSidebarCollapsed: boolean) => {
      if (!pageView) return;

      setInterfaceState((prevState) => {
        return {
          ...prevState,
          queryConfig: {
            ...(prevState?.queryConfig ?? {}),
            [pageView]: {
              filter: updatedFilter.makeQueryParameters(),
              itemsPerPage: updatedFilter.itemsPerPage,
              currentPage: updatedFilter.currentPage,
              sidebarCollapsed: updatedSidebarCollapsed,
            },
          },
        };
      });
    },
    [pageView, setInterfaceState]
  );

  // set the filter and sidebar state when changed
  useEffect(() => {
    if (!pageView) return;
    setLocalFilterState(filter, sidebarCollapsed);
  }, [pageView, filter, setLocalFilterState, sidebarCollapsed]);
}

export function initialSidebarState(defaultSidebarCollapsed: boolean) {
  const isMobile = ScreenUtils.isMobile();
  const initialSidebarCollapsed = isMobile ? true : defaultSidebarCollapsed;
  return initialSidebarCollapsed;
}

export function useListKeyboardShortcuts(props: {
  filter?: ListFilterModel;
  setFilter?: (filter: ListFilterModel) => void;
  showEditFilter?: () => void;
  totalCount?: number;
  toggleSidebarCollapsed?: () => void;
  onSelectAll?: () => void;
  onSelectNone?: () => void;
}) {
  const {
    filter,
    setFilter,
    showEditFilter,
    totalCount = 0,
    toggleSidebarCollapsed,
    onSelectAll,
    onSelectNone,
  } = props;

  // set up hotkeys
  useEffect(() => {
    if (showEditFilter) {
      Mousetrap.bind("f", (e) => {
        showEditFilter();
        // prevent default behavior of typing f in a text field
        // otherwise the filter dialog closes, the query field is focused and
        // f is typed.
        e.preventDefault();
      });

      return () => {
        Mousetrap.unbind("f");
      };
    }
  }, [showEditFilter]);

  useEffect(() => {
    if (!filter || !setFilter || !totalCount) return;

    const pages = Math.ceil(totalCount / filter.itemsPerPage);

    function onChangePage(page: number) {
      if (!filter || !setFilter || !totalCount) return;
      if (page >= 1 && page <= pages) {
        setFilter(filter.changePage(page));
      }
    }

    Mousetrap.bind("right", () => {
      onChangePage(filter.currentPage + 1);
    });
    Mousetrap.bind("left", () => {
      onChangePage(filter.currentPage - 1);
    });
    Mousetrap.bind("shift+right", () => {
      onChangePage(Math.min(pages, filter.currentPage + 10));
    });
    Mousetrap.bind("shift+left", () => {
      onChangePage(Math.max(1, filter.currentPage - 10));
    });
    Mousetrap.bind("ctrl+end", () => {
      onChangePage(pages);
    });
    Mousetrap.bind("ctrl+home", () => {
      onChangePage(1);
    });

    return () => {
      Mousetrap.unbind("right");
      Mousetrap.unbind("left");
      Mousetrap.unbind("shift+right");
      Mousetrap.unbind("shift+left");
      Mousetrap.unbind("ctrl+end");
      Mousetrap.unbind("ctrl+home");
    };
  }, [filter, setFilter, totalCount]);

  useEffect(() => {
    if (toggleSidebarCollapsed) {
      Mousetrap.bind(",", () => {
        toggleSidebarCollapsed();
      });

      return () => {
        Mousetrap.unbind(",");
      };
    }
  }, [toggleSidebarCollapsed]);

  useEffect(() => {
    Mousetrap.bind("s a", () => onSelectAll?.());
    Mousetrap.bind("s n", () => onSelectNone?.());

    return () => {
      Mousetrap.unbind("s a");
      Mousetrap.unbind("s n");
    };
  }, [onSelectAll, onSelectNone]);
}
