package models

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/database"
)

type PerformerQueryBuilder struct{}

func NewPerformerQueryBuilder() PerformerQueryBuilder {
	return PerformerQueryBuilder{}
}

func (qb *PerformerQueryBuilder) Create(newPerformer Performer, tx *sqlx.Tx) (*Performer, error) {
	ensureTx(tx)
	result, err := tx.NamedExec(
		`INSERT INTO performers (checksum, name, url, gender, twitter, instagram, birthdate, ethnicity, country,
                        				eye_color, height, measurements, fake_tits, career_length, tattoos, piercings,
                        				aliases, favorite, created_at, updated_at)
				VALUES (:checksum, :name, :url, :gender, :twitter, :instagram, :birthdate, :ethnicity, :country,
                        :eye_color, :height, :measurements, :fake_tits, :career_length, :tattoos, :piercings,
                        :aliases, :favorite, :created_at, :updated_at)
		`,
		newPerformer,
	)
	if err != nil {
		return nil, err
	}
	performerID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	if err := tx.Get(&newPerformer, `SELECT * FROM performers WHERE id = ? LIMIT 1`, performerID); err != nil {
		return nil, err
	}
	return &newPerformer, nil
}

func (qb *PerformerQueryBuilder) Update(updatedPerformer PerformerPartial, tx *sqlx.Tx) (*Performer, error) {
	ensureTx(tx)
	_, err := tx.NamedExec(
		`UPDATE performers SET `+SQLGenKeysPartial(updatedPerformer)+` WHERE performers.id = :id`,
		updatedPerformer,
	)
	if err != nil {
		return nil, err
	}

	var ret Performer
	if err := tx.Get(&ret, `SELECT * FROM performers WHERE id = ? LIMIT 1`, updatedPerformer.ID); err != nil {
		return nil, err
	}
	return &ret, nil
}

func (qb *PerformerQueryBuilder) UpdateFull(updatedPerformer Performer, tx *sqlx.Tx) (*Performer, error) {
	ensureTx(tx)
	_, err := tx.NamedExec(
		`UPDATE performers SET `+SQLGenKeys(updatedPerformer)+` WHERE performers.id = :id`,
		updatedPerformer,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Get(&updatedPerformer, `SELECT * FROM performers WHERE id = ? LIMIT 1`, updatedPerformer.ID); err != nil {
		return nil, err
	}
	return &updatedPerformer, nil
}

func (qb *PerformerQueryBuilder) Destroy(id string, tx *sqlx.Tx) error {
	_, err := tx.Exec("DELETE FROM performers_scenes WHERE performer_id = ?", id)
	if err != nil {
		return err
	}

	return executeDeleteQuery("performers", id, tx)
}

func (qb *PerformerQueryBuilder) Find(id int) (*Performer, error) {
	query := "SELECT * FROM performers WHERE id = ? LIMIT 1"
	args := []interface{}{id}
	results, err := qb.queryPerformers(query, args, nil)
	if err != nil || len(results) < 1 {
		return nil, err
	}
	return results[0], nil
}

func (qb *PerformerQueryBuilder) FindMany(ids []int) ([]*Performer, error) {
	var performers []*Performer
	for _, id := range ids {
		performer, err := qb.Find(id)
		if err != nil {
			return nil, err
		}

		if performer == nil {
			return nil, fmt.Errorf("performer with id %d not found", id)
		}

		performers = append(performers, performer)
	}

	return performers, nil
}

func (qb *PerformerQueryBuilder) FindBySceneID(sceneID int, tx *sqlx.Tx) ([]*Performer, error) {
	query := selectAll("performers") + `
		LEFT JOIN performers_scenes as scenes_join on scenes_join.performer_id = performers.id
		WHERE scenes_join.scene_id = ?
	`
	args := []interface{}{sceneID}
	return qb.queryPerformers(query, args, tx)
}

func (qb *PerformerQueryBuilder) FindByImageID(imageID int, tx *sqlx.Tx) ([]*Performer, error) {
	query := selectAll("performers") + `
		LEFT JOIN performers_images as images_join on images_join.performer_id = performers.id
		WHERE images_join.image_id = ?
	`
	args := []interface{}{imageID}
	return qb.queryPerformers(query, args, tx)
}

func (qb *PerformerQueryBuilder) FindByGalleryID(galleryID int, tx *sqlx.Tx) ([]*Performer, error) {
	query := selectAll("performers") + `
		LEFT JOIN performers_galleries as galleries_join on galleries_join.performer_id = performers.id
		WHERE galleries_join.gallery_id = ?
	`
	args := []interface{}{galleryID}
	return qb.queryPerformers(query, args, tx)
}

func (qb *PerformerQueryBuilder) FindNameBySceneID(sceneID int, tx *sqlx.Tx) ([]*Performer, error) {
	query := `
		SELECT performers.name FROM performers
		LEFT JOIN performers_scenes as scenes_join on scenes_join.performer_id = performers.id
		WHERE scenes_join.scene_id = ?
	`
	args := []interface{}{sceneID}
	return qb.queryPerformers(query, args, tx)
}

func (qb *PerformerQueryBuilder) FindByNames(names []string, tx *sqlx.Tx, nocase bool) ([]*Performer, error) {
	query := "SELECT * FROM performers WHERE name"
	if nocase {
		query += " COLLATE NOCASE"
	}
	query += " IN " + getInBinding(len(names))

	var args []interface{}
	for _, name := range names {
		args = append(args, name)
	}
	return qb.queryPerformers(query, args, tx)
}

func (qb *PerformerQueryBuilder) Count() (int, error) {
	return runCountQuery(buildCountQuery("SELECT performers.id FROM performers"), nil)
}

func (qb *PerformerQueryBuilder) All() ([]*Performer, error) {
	return qb.queryPerformers(selectAll("performers")+qb.getPerformerSort(nil), nil, nil)
}

func (qb *PerformerQueryBuilder) AllSlim() ([]*Performer, error) {
	return qb.queryPerformers("SELECT performers.id, performers.name, performers.gender FROM performers "+qb.getPerformerSort(nil), nil, nil)
}

func (qb *PerformerQueryBuilder) Query(performerFilter *PerformerFilterType, findFilter *FindFilterType) ([]*Performer, int) {
	if performerFilter == nil {
		performerFilter = &PerformerFilterType{}
	}
	if findFilter == nil {
		findFilter = &FindFilterType{}
	}

	tableName := "performers"
	query := queryBuilder{
		tableName: tableName,
	}

	query.body = selectDistinctIDs(tableName)
	query.body += `
		left join performers_scenes as scenes_join on scenes_join.performer_id = performers.id
		left join scenes on scenes_join.scene_id = scenes.id
		left join performer_stash_ids on performer_stash_ids.performer_id = performers.id
	`

	if q := findFilter.Q; q != nil && *q != "" {
		searchColumns := []string{"performers.name", "performers.checksum", "performers.birthdate", "performers.ethnicity"}
		clause, thisArgs := getSearchBinding(searchColumns, *q, false)
		query.addWhere(clause)
		query.addArg(thisArgs...)
	}

	if favoritesFilter := performerFilter.FilterFavorites; favoritesFilter != nil {
		var favStr string
		if *favoritesFilter == true {
			favStr = "1"
		} else {
			favStr = "0"
		}
		query.addWhere("performers.favorite = " + favStr)
	}

	if birthYear := performerFilter.BirthYear; birthYear != nil {
		clauses, thisArgs := getBirthYearFilterClause(birthYear.Modifier, birthYear.Value)
		query.addWhere(clauses...)
		query.addArg(thisArgs...)
	}

	if age := performerFilter.Age; age != nil {
		clauses, thisArgs := getAgeFilterClause(age.Modifier, age.Value)
		query.addWhere(clauses...)
		query.addArg(thisArgs...)
	}

	if gender := performerFilter.Gender; gender != nil {
		query.addWhere("performers.gender = ?")
		query.addArg(gender.Value.String())
	}

	if isMissingFilter := performerFilter.IsMissing; isMissingFilter != nil && *isMissingFilter != "" {
		switch *isMissingFilter {
		case "scenes":
			query.addWhere("scenes_join.scene_id IS NULL")
		case "image":
			query.body += `left join performers_image on performers_image.performer_id = performers.id
			`
			query.addWhere("performers_image.performer_id IS NULL")
		case "stash_id":
			query.addWhere("performer_stash_ids.performer_id IS NULL")
		default:
			query.addWhere("performers." + *isMissingFilter + " IS NULL OR TRIM(performers." + *isMissingFilter + ") = ''")
		}
	}

	if stashIDFilter := performerFilter.StashID; stashIDFilter != nil {
		query.addWhere("performer_stash_ids.stash_id = ?")
		query.addArg(stashIDFilter)
	}

	query.handleStringCriterionInput(performerFilter.Ethnicity, tableName+".ethnicity")
	query.handleStringCriterionInput(performerFilter.Country, tableName+".country")
	query.handleStringCriterionInput(performerFilter.EyeColor, tableName+".eye_color")
	query.handleStringCriterionInput(performerFilter.Height, tableName+".height")
	query.handleStringCriterionInput(performerFilter.Measurements, tableName+".measurements")
	query.handleStringCriterionInput(performerFilter.FakeTits, tableName+".fake_tits")
	query.handleStringCriterionInput(performerFilter.CareerLength, tableName+".career_length")
	query.handleStringCriterionInput(performerFilter.Tattoos, tableName+".tattoos")
	query.handleStringCriterionInput(performerFilter.Piercings, tableName+".piercings")

	// TODO - need better handling of aliases
	query.handleStringCriterionInput(performerFilter.Aliases, tableName+".aliases")

	query.sortAndPagination = qb.getPerformerSort(findFilter) + getPagination(findFilter)
	idsResult, countResult := query.executeFind()

	var performers []*Performer
	for _, id := range idsResult {
		performer, _ := qb.Find(id)
		performers = append(performers, performer)
	}

	return performers, countResult
}

func getBirthYearFilterClause(criterionModifier CriterionModifier, value int) ([]string, []interface{}) {
	var clauses []string
	var args []interface{}

	yearStr := strconv.Itoa(value)
	startOfYear := yearStr + "-01-01"
	endOfYear := yearStr + "-12-31"

	if modifier := criterionModifier.String(); criterionModifier.IsValid() {
		switch modifier {
		case "EQUALS":
			// between yyyy-01-01 and yyyy-12-31
			clauses = append(clauses, "performers.birthdate >= ?")
			clauses = append(clauses, "performers.birthdate <= ?")
			args = append(args, startOfYear)
			args = append(args, endOfYear)
		case "NOT_EQUALS":
			// outside of yyyy-01-01 to yyyy-12-31
			clauses = append(clauses, "performers.birthdate < ? OR performers.birthdate > ?")
			args = append(args, startOfYear)
			args = append(args, endOfYear)
		case "GREATER_THAN":
			// > yyyy-12-31
			clauses = append(clauses, "performers.birthdate > ?")
			args = append(args, endOfYear)
		case "LESS_THAN":
			// < yyyy-01-01
			clauses = append(clauses, "performers.birthdate < ?")
			args = append(args, startOfYear)
		}
	}

	return clauses, args
}

func getAgeFilterClause(criterionModifier CriterionModifier, value int) ([]string, []interface{}) {
	var clauses []string
	var args []interface{}

	// get the date at which performer would turn the age specified
	dt := time.Now()
	birthDate := dt.AddDate(-value-1, 0, 0)
	yearAfter := birthDate.AddDate(1, 0, 0)

	if modifier := criterionModifier.String(); criterionModifier.IsValid() {
		switch modifier {
		case "EQUALS":
			// between birthDate and yearAfter
			clauses = append(clauses, "performers.birthdate >= ?")
			clauses = append(clauses, "performers.birthdate < ?")
			args = append(args, birthDate)
			args = append(args, yearAfter)
		case "NOT_EQUALS":
			// outside of birthDate and yearAfter
			clauses = append(clauses, "performers.birthdate < ? OR performers.birthdate >= ?")
			args = append(args, birthDate)
			args = append(args, yearAfter)
		case "GREATER_THAN":
			// < birthDate
			clauses = append(clauses, "performers.birthdate < ?")
			args = append(args, birthDate)
		case "LESS_THAN":
			// > yearAfter
			clauses = append(clauses, "performers.birthdate >= ?")
			args = append(args, yearAfter)
		}
	}

	return clauses, args
}

func (qb *PerformerQueryBuilder) getPerformerSort(findFilter *FindFilterType) string {
	var sort string
	var direction string
	if findFilter == nil {
		sort = "name"
		direction = "ASC"
	} else {
		sort = findFilter.GetSort("name")
		direction = findFilter.GetDirection()
	}
	return getSort(sort, direction, "performers")
}

func (qb *PerformerQueryBuilder) queryPerformers(query string, args []interface{}, tx *sqlx.Tx) ([]*Performer, error) {
	var rows *sqlx.Rows
	var err error
	if tx != nil {
		rows, err = tx.Queryx(query, args...)
	} else {
		rows, err = database.DB.Queryx(query, args...)
	}

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	performers := make([]*Performer, 0)
	for rows.Next() {
		performer := Performer{}
		if err := rows.StructScan(&performer); err != nil {
			return nil, err
		}
		performers = append(performers, &performer)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return performers, nil
}

func (qb *PerformerQueryBuilder) UpdatePerformerImage(performerID int, image []byte, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing cover and then create new
	if err := qb.DestroyPerformerImage(performerID, tx); err != nil {
		return err
	}

	_, err := tx.Exec(
		`INSERT INTO performers_image (performer_id, image) VALUES (?, ?)`,
		performerID,
		image,
	)

	return err
}

func (qb *PerformerQueryBuilder) DestroyPerformerImage(performerID int, tx *sqlx.Tx) error {
	ensureTx(tx)

	// Delete the existing joins
	_, err := tx.Exec("DELETE FROM performers_image WHERE performer_id = ?", performerID)
	if err != nil {
		return err
	}
	return err
}

func (qb *PerformerQueryBuilder) GetPerformerImage(performerID int, tx *sqlx.Tx) ([]byte, error) {
	query := `SELECT image from performers_image WHERE performer_id = ?`
	return getImage(tx, query, performerID)
}
