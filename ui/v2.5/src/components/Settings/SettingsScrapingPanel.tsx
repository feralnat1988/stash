import React, { PropsWithChildren, useMemo, useState } from "react";
import { FormattedMessage, useIntl } from "react-intl";
import { Button } from "react-bootstrap";
import {
  mutateReloadScrapers,
  useListGroupScrapers,
  useListPerformerScrapers,
  useListSceneScrapers,
  useListGalleryScrapers,
} from "src/core/StashService";
import { useToast } from "src/hooks/Toast";
import TextUtils from "src/utils/text";
import { CollapseButton } from "../Shared/CollapseButton";
import { Icon } from "../Shared/Icon";
import { LoadingIndicator } from "../Shared/LoadingIndicator";
import { ScrapeType } from "src/core/generated-graphql";
import { SettingSection } from "./SettingSection";
import { BooleanSetting, StringListSetting, StringSetting } from "./Inputs";
import { useSettings } from "./context";
import { StashBoxSetting } from "./StashBoxConfiguration";
import { faSyncAlt } from "@fortawesome/free-solid-svg-icons";
import {
  AvailableScraperPackages,
  InstalledScraperPackages,
} from "./ScraperPackageManager";
import { ExternalLink } from "../Shared/ExternalLink";

const ScraperTable: React.FC<
  PropsWithChildren<{
    entityType: string;
  }>
> = ({ entityType, children }) => {
  const intl = useIntl();
  const title = intl.formatMessage(
    { id: "config.scraping.entity_scrapers" },
    { entityType: intl.formatMessage({ id: entityType }) }
  );

  return (
    <CollapseButton text={title}>
      <table className="scraper-table">
        <thead>
          <tr>
            <th>
              <FormattedMessage id="name" />
            </th>
            <th>
              <FormattedMessage id="config.scraping.supported_types" />
            </th>
            <th>
              <FormattedMessage id="config.scraping.supported_urls" />
            </th>
          </tr>
        </thead>
        <tbody>{children}</tbody>
      </table>
    </CollapseButton>
  );
};

const ScrapeTypeList: React.FC<{
  types: ScrapeType[];
  entityType: string;
}> = ({ types, entityType }) => {
  const intl = useIntl();

  const typeStrings = useMemo(
    () =>
      types.map((t) => {
        switch (t) {
          case ScrapeType.Fragment:
            return intl.formatMessage(
              { id: "config.scraping.entity_metadata" },
              { entityType: intl.formatMessage({ id: entityType }) }
            );
          default:
            return t;
        }
      }),
    [types, entityType, intl]
  );

  return (
    <ul>
      {typeStrings.map((t) => (
        <li key={t}>{t}</li>
      ))}
    </ul>
  );
};

interface IURLList {
  urls: string[];
}

const URLList: React.FC<IURLList> = ({ urls }) => {
  const maxCollapsedItems = 5;
  const [expanded, setExpanded] = useState<boolean>(false);

  const items = useMemo(() => {
    function linkSite(url: string) {
      const u = new URL(url);
      return `${u.protocol}//${u.host}`;
    }

    const ret = urls.map((u) => {
      const sanitised = TextUtils.sanitiseURL(u);
      const siteURL = linkSite(sanitised!);

      return (
        <li key={u}>
          <ExternalLink href={siteURL}>{sanitised}</ExternalLink>
        </li>
      );
    });

    if (ret.length > maxCollapsedItems) {
      if (!expanded) {
        ret.length = maxCollapsedItems;
      }

      ret.push(
        <li>
          <Button onClick={() => setExpanded(!expanded)} variant="link">
            {expanded ? "less" : "more"}
          </Button>
        </li>
      );
    }

    return ret;
  }, [urls, expanded]);

  return <ul>{items}</ul>;
};

const ScraperTableRow: React.FC<{
  name: string;
  entityType: string;
  supportedScrapes: ScrapeType[];
  urls: string[];
}> = ({ name, entityType, supportedScrapes, urls }) => {
  return (
    <tr>
      <td>{name}</td>
      <td>
        <ScrapeTypeList types={supportedScrapes} entityType={entityType} />
      </td>
      <td>
        <URLList urls={urls} />
      </td>
    </tr>
  );
};

export const SettingsScrapingPanel: React.FC = () => {
  const Toast = useToast();
  const { data: performerScrapers, loading: loadingPerformers } =
    useListPerformerScrapers();
  const { data: sceneScrapers, loading: loadingScenes } =
    useListSceneScrapers();
  const { data: galleryScrapers, loading: loadingGalleries } =
    useListGalleryScrapers();
  const { data: groupScrapers, loading: loadingGroups } =
    useListGroupScrapers();

  const { general, scraping, loading, error, saveGeneral, saveScraping } =
    useSettings();

  async function onReloadScrapers() {
    try {
      await mutateReloadScrapers();
    } catch (e) {
      Toast.error(e);
    }
  }

  if (error) return <h1>{error.message}</h1>;
  if (
    loading ||
    loadingScenes ||
    loadingGalleries ||
    loadingPerformers ||
    loadingGroups
  )
    return <LoadingIndicator />;

  return (
    <>
      <StashBoxSetting
        value={general.stashBoxes ?? []}
        onChange={(v) => saveGeneral({ stashBoxes: v })}
      />

      <SettingSection headingID="config.general.scraping">
        <StringSetting
          id="scraperUserAgent"
          headingID="config.general.scraper_user_agent"
          subHeadingID="config.general.scraper_user_agent_desc"
          value={scraping.scraperUserAgent ?? undefined}
          onChange={(v) => saveScraping({ scraperUserAgent: v })}
        />

        <StringSetting
          id="scraperCDPPath"
          headingID="config.general.chrome_cdp_path"
          subHeadingID="config.general.chrome_cdp_path_desc"
          value={scraping.scraperCDPPath ?? undefined}
          onChange={(v) => saveScraping({ scraperCDPPath: v })}
        />

        <BooleanSetting
          id="scraper-cert-check"
          headingID="config.general.check_for_insecure_certificates"
          subHeadingID="config.general.check_for_insecure_certificates_desc"
          checked={scraping.scraperCertCheck ?? undefined}
          onChange={(v) => saveScraping({ scraperCertCheck: v })}
        />

        <StringListSetting
          id="excluded-tag-patterns"
          headingID="config.scraping.excluded_tag_patterns_head"
          subHeadingID="config.scraping.excluded_tag_patterns_desc"
          value={scraping.excludeTagPatterns ?? undefined}
          onChange={(v) => saveScraping({ excludeTagPatterns: v })}
        />
      </SettingSection>

      <InstalledScraperPackages />
      <AvailableScraperPackages />

      <SettingSection headingID="config.scraping.scrapers">
        <div className="content">
          <Button onClick={() => onReloadScrapers()}>
            <span className="fa-icon">
              <Icon icon={faSyncAlt} />
            </span>
            <span>
              <FormattedMessage id="actions.reload_scrapers" />
            </span>
          </Button>
        </div>

        <div className="content">
          <ScraperTable entityType="scene">
            {sceneScrapers?.listScrapers.map((scraper) => (
              <ScraperTableRow
                key={scraper.id}
                name={scraper.name}
                entityType="scene"
                supportedScrapes={scraper.scene?.supported_scrapes ?? []}
                urls={scraper.scene?.urls ?? []}
              />
            ))}
          </ScraperTable>

          <ScraperTable entityType="gallery">
            {galleryScrapers?.listScrapers.map((scraper) => (
              <ScraperTableRow
                key={scraper.id}
                name={scraper.name}
                entityType="gallery"
                supportedScrapes={scraper.gallery?.supported_scrapes ?? []}
                urls={scraper.gallery?.urls ?? []}
              />
            ))}
          </ScraperTable>

          <ScraperTable entityType="performer">
            {performerScrapers?.listScrapers.map((scraper) => (
              <ScraperTableRow
                key={scraper.id}
                name={scraper.name}
                entityType="performer"
                supportedScrapes={scraper.performer?.supported_scrapes ?? []}
                urls={scraper.performer?.urls ?? []}
              />
            ))}
          </ScraperTable>

          <ScraperTable entityType="group">
            {groupScrapers?.listScrapers.map((scraper) => (
              <ScraperTableRow
                key={scraper.id}
                name={scraper.name}
                entityType="group"
                supportedScrapes={scraper.group?.supported_scrapes ?? []}
                urls={scraper.group?.urls ?? []}
              />
            ))}
          </ScraperTable>
        </div>
      </SettingSection>
    </>
  );
};
