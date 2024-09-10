import React, { useMemo, useState } from "react";
import * as GQL from "src/core/generated-graphql";
import TextUtils from "src/utils/text";
import { GalleryLink, TagLink } from "src/components/Shared/TagLink";
import { PerformerCard } from "src/components/Performers/PerformerCard";
import { sortPerformers } from "src/core/performers";
import { FormattedMessage, useIntl } from "react-intl";
import { PhotographerLink } from "src/components/Shared/Link";
import { useContainerDimensions } from "src/components/Shared/GridCard/GridCard";
import {
  DetailItem,
  maybeRenderShowMoreLess,
} from "src/components/Shared/DetailItem";
interface IImageDetailProps {
  image: GQL.ImageDataFragment;
}

export const ImageDetailPanel: React.FC<IImageDetailProps> = (props) => {
  const intl = useIntl();

  const [collapsedDetails, setCollapsedDetails] = useState<boolean>(true);
  const [collapsedPerformers, setCollapsedPerformers] = useState<boolean>(true);
  const [collapsedTags, setCollapsedTags] = useState<boolean>(true);

  const [detailsRef, { height: detailsHeight }] = useContainerDimensions();
  const [perfRef, { height: perfHeight }] = useContainerDimensions();
  const [tagRef, { height: tagHeight }] = useContainerDimensions();

  const details = useMemo(() => {
    const limit = 160;
    return props.image.details?.length ? (
      <>
        <div
          className={`details ${
            collapsedDetails && detailsHeight >= limit
              ? "collapsed-detail"
              : "expanded-detail"
          }`}
        >
          <p className="pre" ref={detailsRef}>
            {props.image.details}
          </p>
        </div>
        {maybeRenderShowMoreLess(
          detailsHeight,
          limit,
          detailsRef,
          setCollapsedDetails,
          collapsedDetails
        )}
      </>
    ) : undefined;
  }, [
    props.image.details,
    detailsRef,
    detailsHeight,
    setCollapsedDetails,
    collapsedDetails,
  ]);

  const tags = useMemo(() => {
    const limit = 160;
    if (props.image.tags.length === 0) return;
    const imageTags = props.image.tags.map((tag) => (
      <TagLink key={tag.id} tag={tag} />
    ));
    return (
      <>
        <div
          className={`image-tags ${
            collapsedTags && tagHeight >= limit
              ? "collapsed-detail"
              : "expanded-detail"
          }`}
          ref={tagRef}
        >
          {imageTags}
        </div>
        {maybeRenderShowMoreLess(
          tagHeight,
          limit,
          tagRef,
          setCollapsedTags,
          collapsedTags
        )}
      </>
    );
  }, [props.image.tags, tagRef, tagHeight, setCollapsedTags, collapsedTags]);

  const galleries = useMemo(() => {
    if (props.image.galleries.length === 0) return;
    const imageGalleries = props.image.galleries.map((gallery) => (
      <GalleryLink key={gallery.id} gallery={gallery} />
    ));
    return (
      <>
        <div className={`image-galleries`}>{imageGalleries}</div>
      </>
    );
  }, [props.image.galleries]);

  const performers = useMemo(() => {
    const limit = 365;
    const sorted = sortPerformers(props.image.performers);
    const cards = sorted.map((performer) => (
      <PerformerCard
        key={performer.id}
        performer={performer}
        ageFromDate={props.image.date ?? undefined}
      />
    ));

    return (
      <>
        <div
          className={`row justify-content-center image-performers ${
            collapsedPerformers && perfHeight >= limit
              ? "collapsed-detail"
              : "expanded-detail"
          }`}
          ref={perfRef}
        >
          {cards}
        </div>
        {maybeRenderShowMoreLess(
          perfHeight,
          limit,
          perfRef,
          setCollapsedPerformers,
          collapsedPerformers
        )}
      </>
    );
  }, [
    props.image.performers,
    props.image.date,
    perfRef,
    perfHeight,
    collapsedPerformers,
    setCollapsedPerformers,
  ]);

  // filename should use entire row if there is no studio
  const imageDetailsWidth = props.image.studio ? "col-9" : "col-12";

  return (
    <>
      <div className="row">
        <div className={`${imageDetailsWidth} col-12 image-details`}>
          <div className="detail-group">
            <DetailItem
              id="studio"
              value={props.image.studio?.name}
              fullWidth
            />
            <DetailItem id="scene_code" value={props.image.code} fullWidth />
            <DetailItem
              id="director"
              value={
                props.image.photographer ? (
                  <PhotographerLink
                    photographer={props.image.photographer}
                    linkType="gallery"
                  />
                ) : undefined
              }
              fullWidth
            />
            <DetailItem id="details" value={details} />
            <DetailItem
              id="tags"
              heading={<FormattedMessage id="tags" />}
              value={props.image.tags.length ? tags : undefined}
            />
            <DetailItem
              id="performers"
              heading={<FormattedMessage id="performers" />}
              value={props.image.performers.length ? performers : undefined}
            />
            <DetailItem
              id="galleries"
              heading={<FormattedMessage id="galleries" />}
              value={props.image.galleries.length ? galleries : undefined}
            />
            <DetailItem
              id="created_at"
              value={TextUtils.formatDateTime(intl, props.image.created_at)}
              fullWidth
            />
            <DetailItem
              id="updated_at"
              value={TextUtils.formatDateTime(intl, props.image.updated_at)}
              fullWidth
            />
          </div>
        </div>
      </div>
    </>
  );
};
