/* eslint-disable @typescript-eslint/no-explicit-any */
import React, { useContext, useEffect, useRef, useState } from "react";
import VideoJS, { VideoJsPlayer, VideoJsPlayerOptions } from "video.js";
import "videojs-vtt-thumbnails-freetube";
import "videojs-seek-buttons";
import "videojs-landscape-fullscreen";
import "./live";
import "./PlaylistButtons";
import cx from "classnames";

import * as GQL from "src/core/generated-graphql";
import { ScenePlayerScrubber } from "./ScenePlayerScrubber";
import { ConfigurationContext } from "src/hooks/Config";
import { Interactive } from "src/utils/interactive";

export const VIDEO_PLAYER_ID = "VideoJsPlayer";

interface IScenePlayerProps {
  className?: string;
  scene: GQL.SceneDataFragment | undefined | null;
  timestamp: number;
  autoplay?: boolean;
  onComplete?: () => void;
  onNext?: () => void;
  onPrevious?: () => void;
}

export const ScenePlayer: React.FC<IScenePlayerProps> = ({
  className,
  autoplay,
  scene,
  timestamp,
  onComplete,
  onNext,
  onPrevious,
}) => {
  const { configuration } = useContext(ConfigurationContext);
  const config = configuration?.interface;
  const videoRef = useRef<HTMLVideoElement>(null);
  const playerRef = useRef<VideoJsPlayer | undefined>();
  const skipButtonsRef = useRef<any>();

  const [time, setTime] = useState(0);

  const [interactiveClient] = useState(
    new Interactive(config?.handyKey || "", config?.funscriptOffset || 0)
  );

  const [initialTimestamp] = useState(timestamp);

  const maxLoopDuration = config?.maximumLoopDuration ?? 0;

  useEffect(() => {
    if (playerRef.current && timestamp >= 0) {
      const player = playerRef.current;
      player.play()?.then(() => {
        player.currentTime(timestamp);
      });
    }
  }, [timestamp]);

  useEffect(() => {
    const videoElement = videoRef.current;
    if (!videoElement) return;

    const options: VideoJsPlayerOptions = {
      controls: true,
      controlBar: {
        pictureInPictureToggle: false,
        volumePanel: {
          inline: false,
        },
      },
      nativeControlsForTouch: false,
      playbackRates: [0.75, 1, 1.5, 2, 3, 4],
      inactivityTimeout: 2000,
      preload: "none",
      userActions: {
        hotkeys: true,
      },
    };

    const player = VideoJS(videoElement, options);

    (player as any).landscapeFullscreen({
      fullscreen: {
        enterOnRotate: true,
        exitOnRotate: true,
        alwaysInLandscapeMode: true,
        iOS: true,
      },
    });

    (player as any).offset();

    player.focus();
    playerRef.current = player;
  }, []);

  useEffect(() => {
    if (scene?.interactive) {
      interactiveClient.uploadScript(scene.paths.funscript || "");
    }
  }, [interactiveClient, scene?.interactive, scene?.paths.funscript]);

  useEffect(() => {
    if (skipButtonsRef.current) {
      skipButtonsRef.current.setForwardHandler(onNext);
      skipButtonsRef.current.setBackwardHandler(onPrevious);
    }
  }, [onNext, onPrevious]);

  useEffect(() => {
    const player = playerRef.current;
    if (player) {
      player.seekButtons({
        forward: 10,
        back: 10,
      });

      skipButtonsRef.current = player.skipButtons() ?? undefined;

      player.focus();
    }

    // Video player destructor
    return () => {
      if (playerRef.current) {
        playerRef.current.dispose();
        playerRef.current = undefined;
      }
    };
  }, []);

  useEffect(() => {
    function handleOffset(player: VideoJsPlayer) {
      if (!scene) return;

      const currentSrc = player.currentSrc();

      const isDirect =
        currentSrc.endsWith("/stream") || currentSrc.endsWith("/stream.m3u8");
      if (!isDirect) {
        (player as any).setOffsetDuration(scene.file.duration);
      } else {
        (player as any).clearOffsetDuration();
      }
    }

    function handleError(play: boolean) {
      const player = playerRef.current;
      if (!player) return;

      const currentFile = player.currentSource();
      if (currentFile) {
        // eslint-disable-next-line no-console
        console.log(`Source failed: ${currentFile.src}`);
        player.focus();
      }

      if (tryNextStream()) {
        // eslint-disable-next-line no-console
        console.log(`Trying next source in playlist: ${player.currentSrc()}`);
        player.load();
        if (play) {
          player.play();
        }
      }
    }

    function tryNextStream() {
      const player = playerRef.current;
      if (!player) return;

      const sources = player.currentSources();

      if (sources.length > 1) {
        sources.shift();
        player.src(sources);
        return true;
      }

      return false;
    }

    if (!scene) return;

    const player = playerRef.current;
    if (!player) return;

    const auto =
      autoplay || (config?.autostartVideo ?? false) || initialTimestamp > 0;
    if (!auto && scene.paths?.screenshot) player.poster(scene.paths.screenshot);
    else player.poster("");
    player.src(
      scene.sceneStreams.map((stream) => ({
        src: stream.url,
        type: stream.mime_type ?? undefined,
        label: stream.label ?? undefined,
      }))
    );
    player.currentTime(0);

    player.loop(
      !!scene.file.duration &&
        maxLoopDuration !== 0 &&
        scene.file.duration < maxLoopDuration
    );

    player.on("loadstart", function (this: VideoJsPlayer) {
      // handle offset after loading so that we get the correct current source
      handleOffset(this);
    });

    player.on("play", function (this: VideoJsPlayer) {
      if (scene.interactive) {
        interactiveClient.play(this.currentTime());
      }
    });

    player.on("pause", () => {
      if (scene.interactive) {
        interactiveClient.pause();
      }
    });

    player.on("timeupdate", function (this: VideoJsPlayer) {
      if (scene.interactive) {
        interactiveClient.ensurePlaying(this.currentTime());
      }

      setTime(this.currentTime());
    });

    player.on("seeking", function (this: VideoJsPlayer) {
      // backwards compatibility - may want to remove this in future
      this.play();
    });

    player.on("error", () => {
      handleError(true);
    });

    player.on("loadedmetadata", () => {
      if (!player.videoWidth() && !player.videoHeight()) {
        // Occurs during preload when videos with supported audio/unsupported video are preloaded.
        // Treat this as a decoding error and try the next source without playing.
        // However on Safari we get an media event when m3u8 is loaded which needs to be ignored.
        const currentFile = player.currentSrc();
        if (currentFile != null && !currentFile.includes("m3u8")) {
          const play = !player.paused();
          handleError(play);
        }
      }
    });

    if (auto) {
      player
        .play()
        ?.then(() => {
          if (initialTimestamp > 0) {
            player.currentTime(initialTimestamp);
          }
        })
        .catch(() => {
          if (scene.paths.screenshot) player.poster(scene.paths.screenshot);
        });
    }

    if ((player as any).vttThumbnails?.src)
      (player as any).vttThumbnails?.src(scene?.paths.vtt);
    else
      (player as any).vttThumbnails({
        src: scene?.paths.vtt,
        showTimestamp: true,
      });
  }, [
    scene,
    config?.autostartVideo,
    maxLoopDuration,
    initialTimestamp,
    autoplay,
    interactiveClient,
  ]);

  useEffect(() => {
    // Attach handler for onComplete event
    const player = playerRef.current;
    if (!player) return;

    player.on("ended", () => {
      onComplete?.();
    });

    return () => player.off("ended");
  }, [onComplete]);

  const onScrubberScrolled = () => {
    playerRef.current?.pause();
  };
  const onScrubberSeek = (seconds: number) => {
    playerRef.current?.currentTime(seconds);
  };

  const isPortrait =
    scene &&
    scene.file.height &&
    scene.file.width &&
    scene.file.height > scene.file.width;

  return (
    <div className={cx("VideoPlayer", { portrait: isPortrait })}>
      <div data-vjs-player className={cx("video-wrapper", className)}>
        <video
          ref={videoRef}
          id={VIDEO_PLAYER_ID}
          className="video-js vjs-big-play-centered"
        />
      </div>
      {scene && (
        <ScenePlayerScrubber
          scene={scene}
          position={time}
          onSeek={onScrubberSeek}
          onScrolled={onScrubberScrolled}
        />
      )}
    </div>
  );
};

export const getPlayerPosition = () =>
  VideoJS.getPlayer(VIDEO_PLAYER_ID).currentTime();
