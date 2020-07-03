// +build plugin_example

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/shurcooL/graphql"
	"github.com/stashapp/stash/pkg/plugin/common"
	"github.com/stashapp/stash/pkg/plugin/util"
)

type api struct {
	stopping bool
}

const tagName = "Hawwwwt"

// graphql inputs and returns
type TagCreate struct {
	ID graphql.ID `graphql:"id"`
}

type TagCreateInput struct {
	Name graphql.String `graphql:"name" json:"name"`
}

type TagDestroyInput struct {
	ID graphql.ID `graphql:"id" json:"id"`
}

type FindScenesResultType struct {
	Count  graphql.Int
	Scenes []Scene
}

type Tag struct {
	ID   graphql.ID     `graphql:"id"`
	Name graphql.String `graphql:"name"`
}

type Scene struct {
	ID   graphql.ID
	Tags []Tag
}

func (s Scene) getTagIds() []graphql.ID {
	ret := []graphql.ID{}

	for _, t := range s.Tags {
		ret = append(ret, t.ID)
	}

	return ret
}

type FindFilterType struct {
	PerPage *graphql.Int    `graphql:"per_page" json:"per_page"`
	Sort    *graphql.String `graphql:"sort" json:"sort"`
}

type SceneUpdate struct {
	ID graphql.ID `graphql:"id"`
}

type SceneUpdateInput struct {
	ID     graphql.ID   `graphql:"id" json:"id"`
	TagIds []graphql.ID `graphql:"tag_ids" json:"tag_ids"`
}

func (a *api) Stop(input struct{}, output *bool) error {
	log.Println("Stopping...")
	a.stopping = true
	*output = true
	return nil
}

func (a *api) Run(input common.PluginInput, output *common.PluginOutput) error {
	client := util.NewClient(input)

	modeArg := common.GetValue(input.Args, "mode")

	var err error
	if modeArg == nil || modeArg.String() == "add" {
		err = addTag(client)
	} else if modeArg.String() == "remove" {
		err = removeTag(client)
	} else if modeArg.String() == "long" {
		err = a.doLongTask()
	}

	if err != nil {
		errStr := err.Error()
		*output = common.PluginOutput{
			Error: &errStr,
		}
		return nil
	}

	outputStr := "ok"
	*output = common.PluginOutput{
		Output: &outputStr,
	}

	return nil
}

func (a *api) doLongTask() error {
	log.Println("Sleeping indefinitely")
	for {
		time.Sleep(time.Second)
		if a.stopping {
			return nil
		}
	}

	return nil
}

func getTagID(client *graphql.Client, create bool) (*graphql.ID, error) {
	log.Println("Checking if tag exists already")

	// see if tag exists already
	var q struct {
		AllTags []Tag `graphql:"allTags"`
	}

	err := client.Query(context.Background(), &q, nil)
	if err != nil {
		return nil, fmt.Errorf("Error getting tags: %s\n", err.Error())
	}

	for _, t := range q.AllTags {
		if t.Name == tagName {
			id := t.ID
			return &id, nil
		}
	}

	if !create {
		log.Println("Not found and not creating")
		return nil, nil
	}

	// create the tag
	var m struct {
		TagCreate TagCreate `graphql:"tagCreate(input: $s)"`
	}

	input := TagCreateInput{
		Name: tagName,
	}

	vars := map[string]interface{}{
		"s": input,
	}

	log.Println("Creating new tag")

	err = client.Mutate(context.Background(), &m, vars)
	if err != nil {
		return nil, fmt.Errorf("Error mutating scene: %s\n", err.Error())
	}

	return &m.TagCreate.ID, nil
}

func findRandomScene(client *graphql.Client) (*Scene, error) {
	// get a random scene
	var q struct {
		FindScenes FindScenesResultType `graphql:"findScenes(filter: $c)"`
	}

	pp := graphql.Int(1)
	sort := graphql.String("random")
	filterInput := &FindFilterType{
		PerPage: &pp,
		Sort:    &sort,
	}

	vars := map[string]interface{}{
		"c": filterInput,
	}

	log.Println("Finding a random scene")
	err := client.Query(context.Background(), &q, vars)
	if err != nil {
		return nil, fmt.Errorf("Error getting random scene: %s\n", err.Error())
	}

	if q.FindScenes.Count == 0 {
		return nil, nil
	}

	return &q.FindScenes.Scenes[0], nil
}

func addTagId(tagIds []graphql.ID, tagId graphql.ID) []graphql.ID {
	for _, t := range tagIds {
		if t == tagId {
			return tagIds
		}
	}

	tagIds = append(tagIds, tagId)
	return tagIds
}

func addTag(client *graphql.Client) error {
	tagID, err := getTagID(client, true)

	if err != nil {
		return err
	}

	scene, err := findRandomScene(client)

	if err != nil {
		return err
	}

	if scene == nil {
		return errors.New("no scenes to add tag to")
	}

	var m struct {
		SceneUpdate SceneUpdate `graphql:"sceneUpdate(input: $s)"`
	}

	input := SceneUpdateInput{
		ID:     scene.ID,
		TagIds: scene.getTagIds(),
	}

	input.TagIds = addTagId(input.TagIds, *tagID)

	vars := map[string]interface{}{
		"s": input,
	}

	log.Printf("Adding tag to scene %v\n", scene.ID)
	err = client.Mutate(context.Background(), &m, vars)
	if err != nil {
		return fmt.Errorf("Error mutating scene: %s", err.Error())
	}

	return nil
}

func removeTag(client *graphql.Client) error {
	tagID, err := getTagID(client, false)

	if err != nil {
		return err
	}

	if tagID == nil {
		log.Println("Tag does not exist. Nothing to remove")
		return nil
	}

	// destroy the tag
	var m struct {
		TagDestroy bool `graphql:"tagDestroy(input: $s)"`
	}

	input := TagDestroyInput{
		ID: *tagID,
	}

	vars := map[string]interface{}{
		"s": input,
	}

	log.Println("Destroying tag")

	err = client.Mutate(context.Background(), &m, vars)
	if err != nil {
		return fmt.Errorf("Error destroying tag: %s", err.Error())
	}

	return nil
}

func main() {
	log.SetPrefix("[Hawwwwt plugin]")

	debug := false

	if len(os.Args) >= 2 && os.Args[1] == "debug" {
		debug = true
	}

	if debug {
		api := api{}
		output := common.PluginOutput{}
		err := api.Run(common.PluginInput{
			ServerConnection: common.StashServerConnection{
				Scheme: "http",
				Port:   9999,
			},
		}, &output)

		if err != nil {
			panic(err)
		}

		if output.Error != nil {
			panic(*output.Error)
		}

		return
	}

	err := common.ServePlugin(&api{})
	if err != nil {
		panic(err)
	}
}
