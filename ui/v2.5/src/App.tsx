import React from "react";
import { Route, Switch, useRouteMatch } from "react-router-dom";
import { IntlProvider } from "react-intl";
import { ToastProvider } from "src/hooks/Toast";
import LightboxProvider from "src/hooks/Lightbox/context";
import { library } from "@fortawesome/fontawesome-svg-core";
import { fas } from "@fortawesome/free-solid-svg-icons";
import { initPolyfills } from "src/polyfills";

import locales from "src/locale";
import { useConfiguration } from "src/core/StashService";
import { flattenMessages } from "src/utils";
import Mousetrap from "mousetrap";
import MousetrapPause from "mousetrap-pause";
import { ErrorBoundary } from "./components/ErrorBoundary";
import Galleries from "./components/Galleries/Galleries";
import { MainNavbar } from "./components/MainNavbar";
import { PageNotFound } from "./components/PageNotFound";
import Performers from "./components/Performers/Performers";
import Scenes from "./components/Scenes/Scenes";
import { Settings } from "./components/Settings/Settings";
import { Stats } from "./components/Stats";
import Studios from "./components/Studios/Studios";
import { SceneFilenameParser } from "./components/SceneFilenameParser/SceneFilenameParser";
import Movies from "./components/Movies/Movies";
import Tags from "./components/Tags/Tags";
import Images from "./components/Images/Images";
import { Setup } from "./components/Setup/Setup";

initPolyfills();

MousetrapPause(Mousetrap);

// Set fontawesome/free-solid-svg as default fontawesome icons
library.add(fas);

const intlFormats = {
  date: {
    long: { year: "numeric", month: "long", day: "numeric" },
  },
};

export const App: React.FC = () => {
  const config = useConfiguration();
  const language = config.data?.configuration?.interface?.language ?? "en-GB";
  const messageLanguage = language.replace(/-/, "");
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const messages = flattenMessages((locales as any)[messageLanguage]);
  const setupMatch = useRouteMatch("/setup");

  function maybeRenderNavbar() {
    // don't render navbar for setup views
    if (!setupMatch) {
      return <MainNavbar />;
    }
  }

  return (
    <ErrorBoundary>
      <IntlProvider locale={language} messages={messages} formats={intlFormats}>
        <ToastProvider>
          <LightboxProvider>
            {maybeRenderNavbar()}
            <div className="main container-fluid">
              <Switch>
                <Route exact path="/" component={Stats} />
                <Route path="/scenes" component={Scenes} />
                <Route path="/images" component={Images} />
                <Route path="/galleries" component={Galleries} />
                <Route path="/performers" component={Performers} />
                <Route path="/tags" component={Tags} />
                <Route path="/studios" component={Studios} />
                <Route path="/movies" component={Movies} />
                <Route path="/settings" component={Settings} />
                <Route
                  path="/sceneFilenameParser"
                  component={SceneFilenameParser}
                />
                <Route path="/setup" component={Setup} />
                <Route component={PageNotFound} />
              </Switch>
            </div>
          </LightboxProvider>
        </ToastProvider>
      </IntlProvider>
    </ErrorBoundary>
  );
};
