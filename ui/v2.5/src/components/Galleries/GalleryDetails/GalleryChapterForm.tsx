import React from "react";
import { Button, Form } from "react-bootstrap";
import { FormattedMessage } from "react-intl";
import { Field, FieldProps, Form as FormikForm, Formik } from "formik";
import * as GQL from "src/core/generated-graphql";
import {
  useGalleryChapterCreate,
  useGalleryChapterUpdate,
  useGalleryChapterDestroy,
} from "src/core/StashService";
import useToast from "src/hooks/Toast";

interface IFormFields {
  title: string;
  pageNumber: string;
}

interface IGalleryChapterForm {
  galleryID: string;
  editingChapter?: GQL.GalleryChapterDataFragment;
  onClose: () => void;
}

export const GalleryChapterForm: React.FC<IGalleryChapterForm> = ({
  galleryID,
  editingChapter,
  onClose,
}) => {

  const [galleryChapterCreate] = useGalleryChapterCreate();
  const [galleryChapterUpdate] = useGalleryChapterUpdate();
  const [galleryChapterDestroy] = useGalleryChapterDestroy();
  const Toast = useToast();

  const onSubmit = (values: IFormFields) => {
    const variables: GQL.GalleryChapterUpdateInput | GQL.GalleryChapterCreateInput = {
      title: values.title,
      page_number: parseInt(values.pageNumber),
      gallery_id: galleryID,
    };
    if (!editingChapter) {
      galleryChapterCreate({ variables })
        .then(onClose)
        .catch((err) => Toast.error(err));
    } else {
      const updateVariables = variables as GQL.GalleryChapterUpdateInput;
      updateVariables.id = editingChapter!.id;
      galleryChapterUpdate({ variables: updateVariables })
        .then(onClose)
        .catch((err) => Toast.error(err));
    }
  };

  const onDelete = () => {
    if (!editingChapter) return;

    galleryChapterDestroy({ variables: { id: editingChapter.id } })
      .then(onClose)
      .catch((err) => Toast.error(err));
  };
  const renderTitleField = (fieldProps: FieldProps<string>) => (
    <input
      className="text-input"
      value={fieldProps.field.value}
      onChange={(query: string) => { fieldProps.form.setFieldValue("title", query.target.value) } }
      mandatory
    />
  );

  const renderPageNumberField = (fieldProps: FieldProps<string>) => (
    <input
      className="text-input"
      value={fieldProps.field.value}
      onChange={(query: string) => { fieldProps.form.setFieldValue("pageNumber", query.target.value)} }
      mandatory
    />
  );

  const values: IFormFields = {
    title: editingChapter?.title ?? "",
    pageNumber: editingChapter?.page_number ?? 1,
  };

  return (
    <Formik initialValues={values} onSubmit={onSubmit}>
      <FormikForm>
        <div>
          <Form.Group>
            <Form.Label
              htmlFor="title"
              className="col-sm-3 col-md-2 col-xl-12 col-form-label"
            >
              Chapter Title
            </Form.Label>
            <div className="col-sm-9 col-md-10 col-xl-12">
              <Field name="title">{renderTitleField}</Field>
            </div>
            <Form.Label
              htmlFor="pageNumber"
              className="col-sm-4 col-md-4 col-xl-12 col-form-label text-sm-right text-xl-left"
            >
              Page
            </Form.Label>
            <div className="col-sm-8 col-xl-12">
              <Field name="pageNumber">{renderPageNumberField}</Field>
          </div>
          </Form.Group>
        </div>
        <div className="buttons-container row">
          <div className="col d-flex">
            <Button variant="primary" type="submit">
              Submit
            </Button>
            <Button
              variant="secondary"
              type="button"
              onClick={onClose}
              className="ml-2"
            >
              <FormattedMessage id="actions.cancel" />
            </Button>
            {editingChapter && (
              <Button
                variant="danger"
                className="ml-auto"
                onClick={() => onDelete()}
              >
                <FormattedMessage id="actions.delete" />
              </Button>
            )}
          </div>
        </div>
      </FormikForm>
    </Formik>
  );
};
