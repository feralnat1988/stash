import { ButtonGroup } from "react-bootstrap";
import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import * as GQL from "src/core/generated-graphql";
import NavUtils from "src/utils/navigation";
import { FormattedMessage } from "react-intl";
import { TruncatedText } from "../Shared/TruncatedText";
import { GridCard, calculateCardWidth } from "../Shared/GridCard/GridCard";
import { PopoverCountButton } from "../Shared/PopoverCountButton";
import ScreenUtils from "src/utils/screen";

interface IProps {
  tag: GQL.TagDataFragment;
  containerWidth?: number;
  zoomIndex: number;
  selecting?: boolean;
  selected?: boolean;
  onSelectedChanged?: (selected: boolean, shiftKey: boolean) => void;
}

export const TagCard: React.FC<IProps> = (props: IProps) => {
  const [cardWidth, setCardWidth] = useState<number>();

  useEffect(() => {
    if (
      !props.containerWidth ||
      props.zoomIndex === undefined ||
      ScreenUtils.isMobile()
    )
      return;

    let zoomValue = props.zoomIndex;
    let preferredCardWidth: number;
    switch (zoomValue) {
      case 0:
        preferredCardWidth = 280;
        break;
      case 1:
        preferredCardWidth = 340;
        break;
      case 2:
        preferredCardWidth = 480;
        break;
      case 3:
        preferredCardWidth = 640;
    }
    let fittedCardWidth = calculateCardWidth(
      props.containerWidth,
      preferredCardWidth!
    );
    setCardWidth(fittedCardWidth);
  }, [props.containerWidth, props.zoomIndex]);

  function maybeRenderDescription() {
    if (props.tag.description) {
      return (
        <TruncatedText
          className="tag-description"
          text={props.tag.description}
          lineCount={3}
        />
      );
    }
  }

  function maybeRenderParents() {
    if (props.tag.parents.length === 1) {
      const parent = props.tag.parents[0];
      return (
        <div className="tag-parent-tags">
          <FormattedMessage
            id="sub_tag_of"
            values={{
              parent: <Link to={`/tags/${parent.id}`}>{parent.name}</Link>,
            }}
          />
        </div>
      );
    }

    if (props.tag.parents.length > 1) {
      return (
        <div className="tag-parent-tags">
          <FormattedMessage
            id="sub_tag_of"
            values={{
              parent: (
                <Link to={NavUtils.makeParentTagsUrl(props.tag)}>
                  {props.tag.parents.length}&nbsp;
                  <FormattedMessage
                    id="countables.tags"
                    values={{ count: props.tag.parents.length }}
                  />
                </Link>
              ),
            }}
          />
        </div>
      );
    }
  }

  function maybeRenderChildren() {
    if (props.tag.children.length > 0) {
      return (
        <div className="tag-sub-tags">
          <FormattedMessage
            id="parent_of"
            values={{
              children: (
                <Link to={NavUtils.makeChildTagsUrl(props.tag)}>
                  {props.tag.children.length}&nbsp;
                  <FormattedMessage
                    id="countables.tags"
                    values={{ count: props.tag.children.length }}
                  />
                </Link>
              ),
            }}
          />
        </div>
      );
    }
  }

  function maybeRenderScenesPopoverButton() {
    if (!props.tag.scene_count) return;

    return (
      <PopoverCountButton
        className="scene-count"
        type="scene"
        count={props.tag.scene_count}
        url={NavUtils.makeTagScenesUrl(props.tag)}
      />
    );
  }

  function maybeRenderSceneMarkersPopoverButton() {
    if (!props.tag.scene_marker_count) return;

    return (
      <PopoverCountButton
        className="marker-count"
        type="marker"
        count={props.tag.scene_marker_count}
        url={NavUtils.makeTagSceneMarkersUrl(props.tag)}
      />
    );
  }

  function maybeRenderImagesPopoverButton() {
    if (!props.tag.image_count) return;

    return (
      <PopoverCountButton
        className="image-count"
        type="image"
        count={props.tag.image_count}
        url={NavUtils.makeTagImagesUrl(props.tag)}
      />
    );
  }

  function maybeRenderGalleriesPopoverButton() {
    if (!props.tag.gallery_count) return;

    return (
      <PopoverCountButton
        className="gallery-count"
        type="gallery"
        count={props.tag.gallery_count}
        url={NavUtils.makeTagGalleriesUrl(props.tag)}
      />
    );
  }

  function maybeRenderPerformersPopoverButton() {
    if (!props.tag.performer_count) return;

    return (
      <PopoverCountButton
        className="performer-count"
        type="performer"
        count={props.tag.performer_count}
        url={NavUtils.makeTagPerformersUrl(props.tag)}
      />
    );
  }

  function maybeRenderPopoverButtonGroup() {
    if (props.tag) {
      return (
        <>
          <hr />
          <ButtonGroup className="card-popovers">
            {maybeRenderScenesPopoverButton()}
            {maybeRenderImagesPopoverButton()}
            {maybeRenderGalleriesPopoverButton()}
            {maybeRenderSceneMarkersPopoverButton()}
            {maybeRenderPerformersPopoverButton()}
          </ButtonGroup>
        </>
      );
    }
  }

  return (
    <GridCard
      className={`tag-card zoom-${props.zoomIndex}`}
      url={`/tags/${props.tag.id}`}
      width={cardWidth}
      title={props.tag.name ?? ""}
      linkClassName="tag-card-header"
      image={
        <img
          loading="lazy"
          className="tag-card-image"
          alt={props.tag.name}
          src={props.tag.image_path ?? ""}
        />
      }
      details={
        <>
          {maybeRenderDescription()}
          {maybeRenderParents()}
          {maybeRenderChildren()}
        </>
      }
      popovers={maybeRenderPopoverButtonGroup()}
      selected={props.selected}
      selecting={props.selecting}
      onSelectedChanged={props.onSelectedChanged}
    />
  );
};
