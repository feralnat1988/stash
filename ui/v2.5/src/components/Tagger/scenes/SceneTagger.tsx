import React, { useContext, useState } from "react";
import * as GQL from "src/core/generated-graphql";
import { SceneQueue } from "src/models/sceneQueue";
import { Button, Form } from "react-bootstrap";
import { FormattedMessage, useIntl } from "react-intl";

import { Icon, LoadingIndicator } from "src/components/Shared";
import { OperationButton } from "src/components/Shared/OperationButton";
import { IScrapedScene, TaggerStateContext } from "../context";
import Config from "./Config";
import { TaggerScene } from "./TaggerScene";
import { SceneTaggerModals } from "./sceneTaggerModals";
import { SceneSearchResults } from "./StashSearchResult";
import { ConfigurationContext } from "src/hooks/Config";
import { faCog, faUnderline } from "@fortawesome/free-solid-svg-icons";
import { distance } from "src/utils/hamming";

interface ITaggerProps {
  scenes: GQL.SlimSceneDataFragment[];
  queue?: SceneQueue;
}

export const Tagger: React.FC<ITaggerProps> = ({ scenes, queue }) => {
  const {
    sources,
    setCurrentSource,
    currentSource,
    doSceneQuery,
    doSceneFragmentScrape,
    doMultiSceneFragmentScrape,
    stopMultiScrape,
    searchResults,
    loading,
    loadingMulti,
    multiError,
    submitFingerprints,
    pendingFingerprints,
  } = useContext(TaggerStateContext);
  const { configuration } = React.useContext(ConfigurationContext);

  const [showConfig, setShowConfig] = useState(false);
  const [hideUnmatched, setHideUnmatched] = useState(false);

  const intl = useIntl();

  const cont = configuration?.interface.continuePlaylistDefault ?? false;

  function generateSceneLink(scene: GQL.SlimSceneDataFragment, index: number) {
    return queue
      ? queue.makeLink(scene.id, { sceneIndex: index, continue: cont })
      : `/scenes/${scene.id}`;
  }

  function handleSourceSelect(e: React.ChangeEvent<HTMLSelectElement>) {
    setCurrentSource(sources!.find((s) => s.id === e.currentTarget.value));
  }

  function renderSourceSelector() {
    return (
      <Form.Group controlId="scraper">
        <Form.Label>
          <FormattedMessage id="component_tagger.config.source" />
        </Form.Label>
        <div>
          <Form.Control
            as="select"
            value={currentSource?.id}
            className="input-control"
            disabled={loading || !sources.length}
            onChange={handleSourceSelect}
          >
            {!sources.length && <option>No scraper sources</option>}
            {sources.map((i) => (
              <option value={i.id} key={i.id}>
                {i.displayName}
              </option>
            ))}
          </Form.Control>
        </div>
      </Form.Group>
    );
  }

  function renderConfigButton() {
    return (
      <div className="ml-2">
        <Button onClick={() => setShowConfig(!showConfig)}>
          <Icon className="fa-fw" icon={faCog} />
        </Button>
      </div>
    );
  }

  function calculatePhashComparisonScore(stashScene: GQL.SlimSceneDataFragment, scrapedScene: IScrapedScene){
    if ( stashScene.phash == undefined)
      return 0

    const phashFingerprints = scrapedScene.fingerprints?.filter((f) => f.algorithm === "PHASH") ?? [];
    const filteredFingerprints = phashFingerprints.filter((f) =>  distance(f.hash, stashScene.phash) <= 8 )

    if (phashFingerprints.length == 0)
      return 0

    return Math.round(filteredFingerprints.length * 100 / phashFingerprints.length)
  }

  function calculateDurationComparisonScore(stashScene: GQL.SlimSceneDataFragment, scrapedScene: IScrapedScene){
    if ( scrapedScene.fingerprints && scrapedScene.fingerprints.length > 0){
      const stashDuration = stashScene.file.duration ?? 0;
      const durations = scrapedScene.fingerprints.map((f) => f.duration);
      const filteredDurations = durations.map((d) => Math.abs(d - stashDuration)).filter((duration) => duration <= 5);

      return Math.round(filteredDurations.length * 100 / durations.length)
    }
    return 0
  }

  function renderScenes() {
    const filteredScenes = !hideUnmatched
      ? scenes
      : scenes.filter((s) => searchResults[s.id]?.results?.length);

    return filteredScenes.map((scene, index) => {
      const sceneLink = generateSceneLink(scene, index);
      let errorMessage: string | undefined;
      const searchResult = searchResults[scene.id];
      if (searchResult?.error) {
        errorMessage = searchResult.error;
      } else if (searchResult && searchResult.results?.length === 0) {
        errorMessage = intl.formatMessage({
          id: "component_tagger.results.match_failed_no_result",
        });
      }
      else if (searchResult && searchResult.results && searchResult.results?.length >= 2) {
        searchResult.results?.sort((scrapedSceneA, scrapedSceneB) => {
          // first order by phash match
          const phashScoreMatchesA = calculatePhashComparisonScore(scene, scrapedSceneA)
          const phashScoreMatchesB = calculatePhashComparisonScore(scene, scrapedSceneB)

          if ( phashScoreMatchesA == phashScoreMatchesB ){
            // if phash scores are identical, compare based on duration matches
            const durationScoreA = calculateDurationComparisonScore(scene, scrapedSceneA)
            const durationScoreB = calculateDurationComparisonScore(scene, scrapedSceneB)

            if ( durationScoreA == 0){
              return -1
            } else if (durationScoreB == 0 ){
              return 1
            }
            
            return durationScoreB - durationScoreA
          } else {
            return phashScoreMatchesB - phashScoreMatchesA
          }
        });
      }

      

      return (
        <TaggerScene
          key={scene.id}
          loading={loading}
          scene={scene}
          url={sceneLink}
          errorMessage={errorMessage}
          doSceneQuery={
            currentSource?.supportSceneQuery
              ? async (v) => {
                  await doSceneQuery(scene.id, v);
                }
              : undefined
          }
          scrapeSceneFragment={
            currentSource?.supportSceneFragment
              ? async () => {
                  await doSceneFragmentScrape(scene.id);
                }
              : undefined
          }
        >
          {searchResult && searchResult.results?.length ? (
            <SceneSearchResults scenes={searchResult.results} target={scene} />
          ) : undefined}
        </TaggerScene>
      );
    });
  }

  const toggleHideUnmatchedScenes = () => {
    setHideUnmatched(!hideUnmatched);
  };

  function maybeRenderShowHideUnmatchedButton() {
    if (Object.keys(searchResults).length) {
      return (
        <Button onClick={toggleHideUnmatchedScenes}>
          <FormattedMessage
            id="component_tagger.verb_toggle_unmatched"
            values={{
              toggle: (
                <FormattedMessage
                  id={`actions.${!hideUnmatched ? "hide" : "show"}`}
                />
              ),
            }}
          />
        </Button>
      );
    }
  }

  function maybeRenderSubmitFingerprintsButton() {
    if (pendingFingerprints.length) {
      return (
        <OperationButton
          className="ml-1"
          operation={submitFingerprints}
          disabled={loading || loadingMulti}
        >
          <span>
            <FormattedMessage
              id="component_tagger.verb_submit_fp"
              values={{ fpCount: pendingFingerprints.length }}
            />
          </span>
        </OperationButton>
      );
    }
  }

  function renderFragmentScrapeButton() {
    if (!currentSource?.supportSceneFragment) {
      return;
    }

    if (loadingMulti) {
      return (
        <Button
          className="ml-1"
          variant="danger"
          onClick={() => {
            stopMultiScrape();
          }}
        >
          <LoadingIndicator message="" inline small />
          <span className="ml-2">
            {intl.formatMessage({ id: "actions.stop" })}
          </span>
        </Button>
      );
    }

    return (
      <div className="ml-1">
        <OperationButton
          disabled={loading}
          operation={async () => {
            await doMultiSceneFragmentScrape(scenes.map((s) => s.id));
          }}
        >
          {intl.formatMessage({ id: "component_tagger.verb_scrape_all" })}
        </OperationButton>
        {multiError && (
          <>
            <br />
            <b className="text-danger">{multiError}</b>
          </>
        )}
      </div>
    );
  }

  return (
    <SceneTaggerModals>
      <div className="tagger-container mx-md-auto">
        <div className="tagger-container-header">
          <div className="d-flex justify-content-between align-items-center flex-wrap">
            <div className="w-auto">{renderSourceSelector()}</div>
            <div className="d-flex">
              {maybeRenderShowHideUnmatchedButton()}
              {maybeRenderSubmitFingerprintsButton()}
              {renderFragmentScrapeButton()}
              {renderConfigButton()}
            </div>
          </div>
          <Config show={showConfig} />
        </div>
        <div>{renderScenes()}</div>
      </div>
    </SceneTaggerModals>
  );
};
