// Code generated by github.com/Yamashou/gqlgenc, DO NOT EDIT.

package graphql

import (
	"context"
	"net/http"

	"github.com/Yamashou/gqlgenc/client"
)

type Client struct {
	Client *client.Client
}

func NewClient(cli *http.Client, baseURL string, options ...client.HTTPRequestOption) *Client {
	return &Client{Client: client.NewClient(cli, baseURL, options...)}
}

type Query struct {
	FindPerformer            *Performer                "json:\"findPerformer\" graphql:\"findPerformer\""
	QueryPerformers          QueryPerformersResultType "json:\"queryPerformers\" graphql:\"queryPerformers\""
	FindStudio               *Studio                   "json:\"findStudio\" graphql:\"findStudio\""
	QueryStudios             QueryStudiosResultType    "json:\"queryStudios\" graphql:\"queryStudios\""
	FindTag                  *Tag                      "json:\"findTag\" graphql:\"findTag\""
	QueryTags                QueryTagsResultType       "json:\"queryTags\" graphql:\"queryTags\""
	FindScene                *Scene                    "json:\"findScene\" graphql:\"findScene\""
	FindSceneByFingerprint   []*Scene                  "json:\"findSceneByFingerprint\" graphql:\"findSceneByFingerprint\""
	FindScenesByFingerprints []*Scene                  "json:\"findScenesByFingerprints\" graphql:\"findScenesByFingerprints\""
	QueryScenes              QueryScenesResultType     "json:\"queryScenes\" graphql:\"queryScenes\""
	FindEdit                 *Edit                     "json:\"findEdit\" graphql:\"findEdit\""
	QueryEdits               QueryEditsResultType      "json:\"queryEdits\" graphql:\"queryEdits\""
	FindUser                 *User                     "json:\"findUser\" graphql:\"findUser\""
	QueryUsers               QueryUsersResultType      "json:\"queryUsers\" graphql:\"queryUsers\""
	Me                       *User                     "json:\"me\" graphql:\"me\""
	SearchPerformer          []*Performer              "json:\"searchPerformer\" graphql:\"searchPerformer\""
	SearchScene              []*Scene                  "json:\"searchScene\" graphql:\"searchScene\""
	Version                  Version                   "json:\"version\" graphql:\"version\""
}

type Mutation struct {
	SceneCreate       *Scene     "json:\"sceneCreate\" graphql:\"sceneCreate\""
	SceneUpdate       *Scene     "json:\"sceneUpdate\" graphql:\"sceneUpdate\""
	SceneDestroy      bool       "json:\"sceneDestroy\" graphql:\"sceneDestroy\""
	PerformerCreate   *Performer "json:\"performerCreate\" graphql:\"performerCreate\""
	PerformerUpdate   *Performer "json:\"performerUpdate\" graphql:\"performerUpdate\""
	PerformerDestroy  bool       "json:\"performerDestroy\" graphql:\"performerDestroy\""
	StudioCreate      *Studio    "json:\"studioCreate\" graphql:\"studioCreate\""
	StudioUpdate      *Studio    "json:\"studioUpdate\" graphql:\"studioUpdate\""
	StudioDestroy     bool       "json:\"studioDestroy\" graphql:\"studioDestroy\""
	TagCreate         *Tag       "json:\"tagCreate\" graphql:\"tagCreate\""
	TagUpdate         *Tag       "json:\"tagUpdate\" graphql:\"tagUpdate\""
	TagDestroy        bool       "json:\"tagDestroy\" graphql:\"tagDestroy\""
	UserCreate        *User      "json:\"userCreate\" graphql:\"userCreate\""
	UserUpdate        *User      "json:\"userUpdate\" graphql:\"userUpdate\""
	UserDestroy       bool       "json:\"userDestroy\" graphql:\"userDestroy\""
	ImageCreate       *Image     "json:\"imageCreate\" graphql:\"imageCreate\""
	ImageUpdate       *Image     "json:\"imageUpdate\" graphql:\"imageUpdate\""
	ImageDestroy      bool       "json:\"imageDestroy\" graphql:\"imageDestroy\""
	RegenerateAPIKey  string     "json:\"regenerateAPIKey\" graphql:\"regenerateAPIKey\""
	ChangePassword    bool       "json:\"changePassword\" graphql:\"changePassword\""
	SceneEdit         Edit       "json:\"sceneEdit\" graphql:\"sceneEdit\""
	PerformerEdit     Edit       "json:\"performerEdit\" graphql:\"performerEdit\""
	StudioEdit        Edit       "json:\"studioEdit\" graphql:\"studioEdit\""
	TagEdit           Edit       "json:\"tagEdit\" graphql:\"tagEdit\""
	EditVote          Edit       "json:\"editVote\" graphql:\"editVote\""
	EditComment       Edit       "json:\"editComment\" graphql:\"editComment\""
	ApplyEdit         Edit       "json:\"applyEdit\" graphql:\"applyEdit\""
	CancelEdit        Edit       "json:\"cancelEdit\" graphql:\"cancelEdit\""
	SubmitFingerprint bool       "json:\"submitFingerprint\" graphql:\"submitFingerprint\""
}
type URLFragment struct {
	URL  string "json:\"url\" graphql:\"url\""
	Type string "json:\"type\" graphql:\"type\""
}
type ImageFragment struct {
	ID     string "json:\"id\" graphql:\"id\""
	URL    string "json:\"url\" graphql:\"url\""
	Width  *int   "json:\"width\" graphql:\"width\""
	Height *int   "json:\"height\" graphql:\"height\""
}
type StudioFragment struct {
	Name   string           "json:\"name\" graphql:\"name\""
	ID     string           "json:\"id\" graphql:\"id\""
	Urls   []*URLFragment   "json:\"urls\" graphql:\"urls\""
	Images []*ImageFragment "json:\"images\" graphql:\"images\""
}
type TagFragment struct {
	Name string "json:\"name\" graphql:\"name\""
	ID   string "json:\"id\" graphql:\"id\""
}
type FuzzyDateFragment struct {
	Date     string           "json:\"date\" graphql:\"date\""
	Accuracy DateAccuracyEnum "json:\"accuracy\" graphql:\"accuracy\""
}
type MeasurementsFragment struct {
	BandSize *int    "json:\"band_size\" graphql:\"band_size\""
	CupSize  *string "json:\"cup_size\" graphql:\"cup_size\""
	Waist    *int    "json:\"waist\" graphql:\"waist\""
	Hip      *int    "json:\"hip\" graphql:\"hip\""
}
type BodyModificationFragment struct {
	Location    string  "json:\"location\" graphql:\"location\""
	Description *string "json:\"description\" graphql:\"description\""
}
type PerformerFragment struct {
	ID              string                      "json:\"id\" graphql:\"id\""
	Name            string                      "json:\"name\" graphql:\"name\""
	Disambiguation  *string                     "json:\"disambiguation\" graphql:\"disambiguation\""
	Aliases         []string                    "json:\"aliases\" graphql:\"aliases\""
	Gender          *GenderEnum                 "json:\"gender\" graphql:\"gender\""
	Urls            []*URLFragment              "json:\"urls\" graphql:\"urls\""
	Images          []*ImageFragment            "json:\"images\" graphql:\"images\""
	Birthdate       *FuzzyDateFragment          "json:\"birthdate\" graphql:\"birthdate\""
	Ethnicity       *EthnicityEnum              "json:\"ethnicity\" graphql:\"ethnicity\""
	Country         *string                     "json:\"country\" graphql:\"country\""
	EyeColor        *EyeColorEnum               "json:\"eye_color\" graphql:\"eye_color\""
	HairColor       *HairColorEnum              "json:\"hair_color\" graphql:\"hair_color\""
	Height          *int                        "json:\"height\" graphql:\"height\""
	Measurements    MeasurementsFragment        "json:\"measurements\" graphql:\"measurements\""
	BreastType      *BreastTypeEnum             "json:\"breast_type\" graphql:\"breast_type\""
	CareerStartYear *int                        "json:\"career_start_year\" graphql:\"career_start_year\""
	CareerEndYear   *int                        "json:\"career_end_year\" graphql:\"career_end_year\""
	Tattoos         []*BodyModificationFragment "json:\"tattoos\" graphql:\"tattoos\""
	Piercings       []*BodyModificationFragment "json:\"piercings\" graphql:\"piercings\""
}
type PerformerAppearanceFragment struct {
	As        *string           "json:\"as\" graphql:\"as\""
	Performer PerformerFragment "json:\"performer\" graphql:\"performer\""
}
type FingerprintFragment struct {
	Algorithm FingerprintAlgorithm "json:\"algorithm\" graphql:\"algorithm\""
	Hash      string               "json:\"hash\" graphql:\"hash\""
	Duration  int                  "json:\"duration\" graphql:\"duration\""
}
type SceneFragment struct {
	ID           string                         "json:\"id\" graphql:\"id\""
	Title        *string                        "json:\"title\" graphql:\"title\""
	Details      *string                        "json:\"details\" graphql:\"details\""
	Duration     *int                           "json:\"duration\" graphql:\"duration\""
	Date         *string                        "json:\"date\" graphql:\"date\""
	Urls         []*URLFragment                 "json:\"urls\" graphql:\"urls\""
	Images       []*ImageFragment               "json:\"images\" graphql:\"images\""
	Studio       *StudioFragment                "json:\"studio\" graphql:\"studio\""
	Tags         []*TagFragment                 "json:\"tags\" graphql:\"tags\""
	Performers   []*PerformerAppearanceFragment "json:\"performers\" graphql:\"performers\""
	Fingerprints []*FingerprintFragment         "json:\"fingerprints\" graphql:\"fingerprints\""
}
type FindSceneByFingerprint struct {
	FindSceneByFingerprint []*SceneFragment "json:\"findSceneByFingerprint\" graphql:\"findSceneByFingerprint\""
}
type FindScenesByFingerprints struct {
	FindScenesByFingerprints []*SceneFragment "json:\"findScenesByFingerprints\" graphql:\"findScenesByFingerprints\""
}
type SearchScene struct {
	SearchScene []*SceneFragment "json:\"searchScene\" graphql:\"searchScene\""
}
type SubmitFingerprintPayload struct {
	SubmitFingerprint bool "json:\"submitFingerprint\" graphql:\"submitFingerprint\""
}

const FindSceneByFingerprintQuery = `query FindSceneByFingerprint ($fingerprint: FingerprintQueryInput!) {
	findSceneByFingerprint(fingerprint: $fingerprint) {
		... SceneFragment
	}
}
fragment SceneFragment on Scene {
	id
	title
	details
	duration
	date
	urls {
		... URLFragment
	}
	images {
		... ImageFragment
	}
	studio {
		... StudioFragment
	}
	tags {
		... TagFragment
	}
	performers {
		... PerformerAppearanceFragment
	}
	fingerprints {
		... FingerprintFragment
	}
}
fragment URLFragment on URL {
	url
	type
}
fragment PerformerAppearanceFragment on PerformerAppearance {
	as
	performer {
		... PerformerFragment
	}
}
fragment PerformerFragment on Performer {
	id
	name
	disambiguation
	aliases
	gender
	urls {
		... URLFragment
	}
	images {
		... ImageFragment
	}
	birthdate {
		... FuzzyDateFragment
	}
	ethnicity
	country
	eye_color
	hair_color
	height
	measurements {
		... MeasurementsFragment
	}
	breast_type
	career_start_year
	career_end_year
	tattoos {
		... BodyModificationFragment
	}
	piercings {
		... BodyModificationFragment
	}
}
fragment FuzzyDateFragment on FuzzyDate {
	date
	accuracy
}
fragment FingerprintFragment on Fingerprint {
	algorithm
	hash
	duration
}
fragment ImageFragment on Image {
	id
	url
	width
	height
}
fragment StudioFragment on Studio {
	name
	id
	urls {
		... URLFragment
	}
	images {
		... ImageFragment
	}
}
fragment TagFragment on Tag {
	name
	id
}
fragment MeasurementsFragment on Measurements {
	band_size
	cup_size
	waist
	hip
}
fragment BodyModificationFragment on BodyModification {
	location
	description
}
`

func (c *Client) FindSceneByFingerprint(ctx context.Context, fingerprint FingerprintQueryInput, httpRequestOptions ...client.HTTPRequestOption) (*FindSceneByFingerprint, error) {
	vars := map[string]interface{}{
		"fingerprint": fingerprint,
	}

	var res FindSceneByFingerprint
	if err := c.Client.Post(ctx, FindSceneByFingerprintQuery, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const FindScenesByFingerprintsQuery = `query FindScenesByFingerprints ($fingerprints: [String!]!) {
	findScenesByFingerprints(fingerprints: $fingerprints) {
		... SceneFragment
	}
}
fragment PerformerAppearanceFragment on PerformerAppearance {
	as
	performer {
		... PerformerFragment
	}
}
fragment PerformerFragment on Performer {
	id
	name
	disambiguation
	aliases
	gender
	urls {
		... URLFragment
	}
	images {
		... ImageFragment
	}
	birthdate {
		... FuzzyDateFragment
	}
	ethnicity
	country
	eye_color
	hair_color
	height
	measurements {
		... MeasurementsFragment
	}
	breast_type
	career_start_year
	career_end_year
	tattoos {
		... BodyModificationFragment
	}
	piercings {
		... BodyModificationFragment
	}
}
fragment MeasurementsFragment on Measurements {
	band_size
	cup_size
	waist
	hip
}
fragment FingerprintFragment on Fingerprint {
	algorithm
	hash
	duration
}
fragment SceneFragment on Scene {
	id
	title
	details
	duration
	date
	urls {
		... URLFragment
	}
	images {
		... ImageFragment
	}
	studio {
		... StudioFragment
	}
	tags {
		... TagFragment
	}
	performers {
		... PerformerAppearanceFragment
	}
	fingerprints {
		... FingerprintFragment
	}
}
fragment URLFragment on URL {
	url
	type
}
fragment ImageFragment on Image {
	id
	url
	width
	height
}
fragment TagFragment on Tag {
	name
	id
}
fragment StudioFragment on Studio {
	name
	id
	urls {
		... URLFragment
	}
	images {
		... ImageFragment
	}
}
fragment FuzzyDateFragment on FuzzyDate {
	date
	accuracy
}
fragment BodyModificationFragment on BodyModification {
	location
	description
}
`

func (c *Client) FindScenesByFingerprints(ctx context.Context, fingerprints []string, httpRequestOptions ...client.HTTPRequestOption) (*FindScenesByFingerprints, error) {
	vars := map[string]interface{}{
		"fingerprints": fingerprints,
	}

	var res FindScenesByFingerprints
	if err := c.Client.Post(ctx, FindScenesByFingerprintsQuery, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const SearchSceneQuery = `query SearchScene ($term: String!) {
	searchScene(term: $term) {
		... SceneFragment
	}
}
fragment URLFragment on URL {
	url
	type
}
fragment ImageFragment on Image {
	id
	url
	width
	height
}
fragment TagFragment on Tag {
	name
	id
}
fragment PerformerFragment on Performer {
	id
	name
	disambiguation
	aliases
	gender
	urls {
		... URLFragment
	}
	images {
		... ImageFragment
	}
	birthdate {
		... FuzzyDateFragment
	}
	ethnicity
	country
	eye_color
	hair_color
	height
	measurements {
		... MeasurementsFragment
	}
	breast_type
	career_start_year
	career_end_year
	tattoos {
		... BodyModificationFragment
	}
	piercings {
		... BodyModificationFragment
	}
}
fragment FuzzyDateFragment on FuzzyDate {
	date
	accuracy
}
fragment MeasurementsFragment on Measurements {
	band_size
	cup_size
	waist
	hip
}
fragment SceneFragment on Scene {
	id
	title
	details
	duration
	date
	urls {
		... URLFragment
	}
	images {
		... ImageFragment
	}
	studio {
		... StudioFragment
	}
	tags {
		... TagFragment
	}
	performers {
		... PerformerAppearanceFragment
	}
	fingerprints {
		... FingerprintFragment
	}
}
fragment StudioFragment on Studio {
	name
	id
	urls {
		... URLFragment
	}
	images {
		... ImageFragment
	}
}
fragment PerformerAppearanceFragment on PerformerAppearance {
	as
	performer {
		... PerformerFragment
	}
}
fragment BodyModificationFragment on BodyModification {
	location
	description
}
fragment FingerprintFragment on Fingerprint {
	algorithm
	hash
	duration
}
`

func (c *Client) SearchScene(ctx context.Context, term string, httpRequestOptions ...client.HTTPRequestOption) (*SearchScene, error) {
	vars := map[string]interface{}{
		"term": term,
	}

	var res SearchScene
	if err := c.Client.Post(ctx, SearchSceneQuery, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const SubmitFingerprintQuery = `mutation SubmitFingerprint ($input: FingerprintSubmission!) {
	submitFingerprint(input: $input)
}
`

func (c *Client) SubmitFingerprint(ctx context.Context, input FingerprintSubmission, httpRequestOptions ...client.HTTPRequestOption) (*SubmitFingerprintPayload, error) {
	vars := map[string]interface{}{
		"input": input,
	}

	var res SubmitFingerprintPayload
	if err := c.Client.Post(ctx, SubmitFingerprintQuery, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}
