package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/file"
	"gopkg.in/guregu/null.v4"
)

const (
	fileTable              = "files"
	videoFileTable         = "video_files"
	imageFileTable         = "image_files"
	fileIDColumn           = "file_id"
	filesFingerprintsTable = "files_fingerprints"
)

type basicFileRow struct {
	ID             file.ID       `db:"id" goqu:"skipinsert"`
	Basename       string        `db:"basename"`
	ZipFileID      null.Int      `db:"zip_file_id"`
	ParentFolderID file.FolderID `db:"parent_folder_id"`
	Size           int64         `db:"size"`
	ModTime        time.Time     `db:"mod_time"`
	MissingSince   null.Time     `db:"missing_since"`
	LastScanned    time.Time     `db:"last_scanned"`
	CreatedAt      time.Time     `db:"created_at"`
	UpdatedAt      time.Time     `db:"updated_at"`
}

func (r *basicFileRow) fromBasicFile(o file.BaseFile) {
	r.ID = o.ID
	r.Basename = o.Basename
	r.ZipFileID = nullIntFromFileIDPtr(o.ZipFileID)
	r.ParentFolderID = o.ParentFolderID
	r.Size = o.Size
	r.ModTime = o.ModTime
	r.MissingSince = null.TimeFromPtr(o.MissingSince)
	r.LastScanned = o.LastScanned
	r.CreatedAt = o.CreatedAt
	r.UpdatedAt = o.UpdatedAt
}

type videoFileRow struct {
	FileID     file.ID `db:"file_id"`
	Format     string  `db:"format"`
	Width      int     `db:"width"`
	Height     int     `db:"height"`
	Duration   float64 `db:"duration"`
	VideoCodec string  `db:"video_codec"`
	AudioCodec string  `db:"audio_codec"`
	FrameRate  float64 `db:"frame_rate"`
	BitRate    int64   `db:"bit_rate"`
}

func (f *videoFileRow) fromVideoFile(ff file.VideoFile) {
	f.FileID = ff.ID
	f.Format = ff.Format
	f.Width = ff.Width
	f.Height = ff.Height
	f.Duration = ff.Duration
	f.VideoCodec = ff.VideoCodec
	f.AudioCodec = ff.AudioCodec
	f.FrameRate = ff.FrameRate
	f.BitRate = ff.BitRate
}

type imageFileRow struct {
	FileID file.ID `db:"file_id"`
	Format string  `db:"format"`
	Width  int     `db:"width"`
	Height int     `db:"height"`
}

func (f *imageFileRow) fromImageFile(ff file.ImageFile) {
	f.FileID = ff.ID
	f.Format = ff.Format
	f.Width = ff.Width
	f.Height = ff.Height
}

// we redefine this to change the columns around
// otherwise, we collide with the image file columns
type videoFileQueryRow struct {
	FileID     null.Int    `db:"file_id_video"`
	Format     null.String `db:"video_format"`
	Width      null.Int    `db:"video_width"`
	Height     null.Int    `db:"video_height"`
	Duration   null.Float  `db:"duration"`
	VideoCodec null.String `db:"video_codec"`
	AudioCodec null.String `db:"audio_codec"`
	FrameRate  null.Float  `db:"frame_rate"`
	BitRate    null.Int    `db:"bit_rate"`
}

func (f *videoFileQueryRow) resolve() *file.VideoFile {
	return &file.VideoFile{
		Format:     f.Format.String,
		Width:      int(f.Width.Int64),
		Height:     int(f.Height.Int64),
		Duration:   f.Duration.Float64,
		VideoCodec: f.VideoCodec.String,
		AudioCodec: f.AudioCodec.String,
		FrameRate:  f.FrameRate.Float64,
		BitRate:    f.BitRate.Int64,
	}
}

func videoFileQueryColumns() []interface{} {
	table := videoFileTableMgr.table
	return []interface{}{
		table.Col("file_id").As("file_id_video"),
		table.Col("format").As("video_format"),
		table.Col("width").As("video_width"),
		table.Col("height").As("video_height"),
		table.Col("duration"),
		table.Col("video_codec"),
		table.Col("audio_codec"),
		table.Col("frame_rate"),
		table.Col("bit_rate"),
	}
}

// we redefine this to change the columns around
// otherwise, we collide with the video file columns
type imageFileQueryRow struct {
	Format null.String `db:"image_format"`
	Width  null.Int    `db:"image_width"`
	Height null.Int    `db:"image_height"`
}

func (imageFileQueryRow) columns(table *table) []interface{} {
	ex := table.table
	return []interface{}{
		ex.Col("format").As("image_format"),
		ex.Col("width").As("image_width"),
		ex.Col("height").As("image_height"),
	}
}

func (f *imageFileQueryRow) resolve() *file.ImageFile {
	return &file.ImageFile{
		Format: f.Format.String,
		Width:  int(f.Width.Int64),
		Height: int(f.Height.Int64),
	}
}

type fileQueryRow struct {
	FileID         null.Int    `db:"file_id"`
	Basename       null.String `db:"basename"`
	ZipFileID      null.Int    `db:"zip_file_id"`
	ParentFolderID null.Int    `db:"parent_folder_id"`
	Size           null.Int    `db:"size"`
	ModTime        null.Time   `db:"mod_time"`
	MissingSince   null.Time   `db:"missing_since"`
	LastScanned    null.Time   `db:"last_scanned"`
	CreatedAt      null.Time   `db:"created_at"`
	UpdatedAt      null.Time   `db:"updated_at"`

	ZipBasename   null.String `db:"zip_basename"`
	ZipFolderPath null.String `db:"zip_folder_path"`

	FolderPath null.String `db:"folder_path"`
	fingerprintQueryRow
	videoFileQueryRow
	imageFileQueryRow
}

func (r *fileQueryRow) resolve() file.File {
	basic := &file.BaseFile{
		ID: file.ID(r.FileID.Int64),
		DirEntry: file.DirEntry{
			ZipFileID:    nullIntFileIDPtr(r.ZipFileID),
			ModTime:      r.ModTime.Time,
			MissingSince: r.MissingSince.Ptr(),
			LastScanned:  r.LastScanned.Time,
		},
		Path:           filepath.Join(r.FolderPath.String, r.Basename.String),
		ParentFolderID: file.FolderID(r.ParentFolderID.Int64),
		Basename:       r.Basename.String,
		Size:           r.Size.Int64,
		CreatedAt:      r.CreatedAt.Time,
		UpdatedAt:      r.UpdatedAt.Time,
	}

	if basic.ZipFileID != nil && r.ZipFolderPath.Valid && r.ZipBasename.Valid {
		basic.ZipFile = &file.BaseFile{
			ID:       *basic.ZipFileID,
			Path:     filepath.Join(r.ZipFolderPath.String, r.ZipBasename.String),
			Basename: r.ZipBasename.String,
		}
	}

	var ret file.File = basic

	if r.videoFileQueryRow.FileID.Valid {
		vf := r.videoFileQueryRow.resolve()
		vf.BaseFile = basic
		ret = vf
	}

	if r.imageFileQueryRow.Format.Valid {
		imf := r.imageFileQueryRow.resolve()
		imf.BaseFile = basic
		ret = imf
	}

	r.appendRelationships(basic)

	return ret
}

func appendFingerprintsUnique(vs []file.Fingerprint, v ...file.Fingerprint) []file.Fingerprint {
	for _, vv := range v {
		found := false
		for _, vsv := range vs {
			if vsv.Type == vv.Type {
				found = true
				break
			}
		}

		if !found {
			vs = append(vs, vv)
		}
	}
	return vs
}

func (r *fileQueryRow) appendRelationships(i *file.BaseFile) {
	if r.fingerprintQueryRow.valid() {
		i.Fingerprints = appendFingerprintsUnique(i.Fingerprints, r.fingerprintQueryRow.resolve())
	}
}

func mergeFiles(dest file.File, src file.File) {
	if src.Base().Fingerprints != nil {
		dest.Base().Fingerprints = appendFingerprintsUnique(dest.Base().Fingerprints, src.Base().Fingerprints...)
	}
}

type fileQueryRows []fileQueryRow

func (r fileQueryRows) resolve() []file.File {
	var ret []file.File
	var last file.File
	var lastID file.ID

	for _, row := range r {
		if last == nil || lastID != file.ID(row.FileID.Int64) {
			f := row.resolve()
			last = f
			lastID = file.ID(row.FileID.Int64)
			ret = append(ret, last)
			continue
		}

		// must be merging with previous row
		row.appendRelationships(last.Base())
	}

	return ret
}

type relatedFileQueryRow struct {
	fileQueryRow
	Primary null.Bool `db:"primary"`
}

type FileStore struct {
	repository

	tableMgr *table
}

func NewFileStore() *FileStore {
	return &FileStore{
		repository: repository{
			tableName: sceneTable,
			idColumn:  idColumn,
		},

		tableMgr: fileTableMgr,
	}
}

func (qb *FileStore) table() exp.IdentifierExpression {
	return qb.tableMgr.table
}

func (qb *FileStore) Create(ctx context.Context, f file.File) error {
	var r basicFileRow
	r.fromBasicFile(*f.Base())

	id, err := qb.tableMgr.insertID(ctx, r)
	if err != nil {
		return err
	}

	fileID := file.ID(id)

	// create extended stuff here
	switch ef := f.(type) {
	case *file.VideoFile:
		if err := qb.createVideoFile(ctx, fileID, *ef); err != nil {
			return err
		}
	case *file.ImageFile:
		if err := qb.createImageFile(ctx, fileID, *ef); err != nil {
			return err
		}
	}

	fpIDs, err := qb.getOrCreateFingerprintIDs(ctx, f.Base())
	if err != nil {
		return err
	}

	if err := filesFingerprintsTableMgr.insertJoins(ctx, id, fpIDs); err != nil {
		return err
	}

	// only assign id once we are successful
	f.Base().ID = fileID

	return nil
}

func (qb *FileStore) Update(ctx context.Context, f file.File) error {
	var r basicFileRow
	r.fromBasicFile(*f.Base())

	id := f.Base().ID

	if err := qb.tableMgr.updateByID(ctx, id, r); err != nil {
		return err
	}

	// create extended stuff here
	switch ef := f.(type) {
	case *file.VideoFile:
		if err := qb.updateOrCreateVideoFile(ctx, id, *ef); err != nil {
			return err
		}
	case *file.ImageFile:
		if err := qb.updateOrCreateImageFile(ctx, id, *ef); err != nil {
			return err
		}
	}

	fpIDs, err := qb.getOrCreateFingerprintIDs(ctx, f.Base())
	if err != nil {
		return err
	}

	if err := filesFingerprintsTableMgr.replaceJoins(ctx, int(id), fpIDs); err != nil {
		return err
	}

	// TODO - delete unused fingerprints

	return nil
}

func (qb *FileStore) createVideoFile(ctx context.Context, id file.ID, f file.VideoFile) error {
	var r videoFileRow
	r.fromVideoFile(f)
	r.FileID = id
	if _, err := videoFileTableMgr.insert(ctx, r); err != nil {
		return err
	}

	return nil
}

func (qb *FileStore) updateOrCreateVideoFile(ctx context.Context, id file.ID, f file.VideoFile) error {
	exists, err := videoFileTableMgr.idExists(ctx, id)
	if err != nil {
		return err
	}

	if !exists {
		return qb.createVideoFile(ctx, id, f)
	}

	var r videoFileRow
	r.fromVideoFile(f)
	r.FileID = id
	if err := videoFileTableMgr.updateByID(ctx, id, r); err != nil {
		return err
	}

	return nil
}

func (qb *FileStore) createImageFile(ctx context.Context, id file.ID, f file.ImageFile) error {
	var r imageFileRow
	r.fromImageFile(f)
	r.FileID = id
	if _, err := imageFileTableMgr.insert(ctx, r); err != nil {
		return err
	}

	return nil
}

func (qb *FileStore) updateOrCreateImageFile(ctx context.Context, id file.ID, f file.ImageFile) error {
	exists, err := imageFileTableMgr.idExists(ctx, id)
	if err != nil {
		return err
	}

	if !exists {
		return qb.createImageFile(ctx, id, f)
	}

	var r imageFileRow
	r.fromImageFile(f)
	r.FileID = id
	if err := imageFileTableMgr.updateByID(ctx, id, r); err != nil {
		return err
	}

	return nil
}

func (qb *FileStore) getOrCreateFingerprintIDs(ctx context.Context, f *file.BaseFile) ([]int, error) {
	fpqb := FingerprintReaderWriter
	var ids []int
	for _, fp := range f.Fingerprints {
		id, err := fpqb.getOrCreate(ctx, fp)
		if err != nil {
			return nil, err
		}

		if id != nil {
			ids = append(ids, *id)
		}
	}

	return ids, nil
}

func (qb *FileStore) selectDataset() *goqu.SelectDataset {
	table := qb.table()

	folderTable := folderTableMgr.table
	fingerprintTable := fingerprintTableMgr.table
	videoFileTable := videoFileTableMgr.table
	imageFileTable := imageFileTableMgr.table

	cols := []interface{}{
		table.Col("id").As("file_id"),
		table.Col("basename"),
		table.Col("zip_file_id"),
		table.Col("parent_folder_id"),
		table.Col("size"),
		table.Col("mod_time"),
		table.Col("missing_since"),
		table.Col("last_scanned"),
		table.Col("created_at"),
		table.Col("updated_at"),
		folderTable.Col("path").As("folder_path"),
		fingerprintTable.Col("type").As("fingerprint_type"),
		fingerprintTable.Col("fingerprint"),
	}

	cols = append(cols, videoFileQueryColumns()...)
	cols = append(cols, imageFileQueryRow{}.columns(imageFileTableMgr)...)

	ret := dialect.From(table).Select(cols...)

	return ret.InnerJoin(
		folderTable,
		goqu.On(table.Col("parent_folder_id").Eq(folderTable.Col(idColumn))),
	).LeftJoin(
		filesFingerprintsJoinTable,
		goqu.On(table.Col(idColumn).Eq(filesFingerprintsJoinTable.Col(fileIDColumn))),
	).LeftJoin(
		fingerprintTable,
		goqu.On(filesFingerprintsJoinTable.Col(fingerprintIDColumn).Eq(fingerprintTable.Col(idColumn))),
	).LeftJoin(
		videoFileTable,
		goqu.On(table.Col(idColumn).Eq(videoFileTable.Col(fileIDColumn))),
	).LeftJoin(
		imageFileTable,
		goqu.On(table.Col(idColumn).Eq(imageFileTable.Col(fileIDColumn))),
	)
}

func (qb *FileStore) get(ctx context.Context, q *goqu.SelectDataset) (file.File, error) {
	ret, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, err
	}

	if len(ret) == 0 {
		return nil, sql.ErrNoRows
	}

	return ret[0], nil
}

func (qb *FileStore) getMany(ctx context.Context, q *goqu.SelectDataset) ([]file.File, error) {
	const single = false
	var rows fileQueryRows
	if err := queryFunc(ctx, q, single, func(r *sqlx.Rows) error {
		var f fileQueryRow
		if err := r.StructScan(&f); err != nil {
			return err
		}

		rows = append(rows, f)
		return nil
	}); err != nil {
		return nil, err
	}

	return rows.resolve(), nil
}

func (qb *FileStore) Find(ctx context.Context, id file.ID) (file.File, error) {
	q := qb.selectDataset().Where(qb.tableMgr.byID(id))

	ret, err := qb.get(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("getting file by id %d: %w", id, err)
	}

	return ret, nil
}

func (qb *FileStore) FindByPath(ctx context.Context, p string) (file.File, error) {
	// separate basename from path
	basename := filepath.Base(p)
	dir, _ := path(filepath.Dir(p)).Value()

	table := qb.table()
	folderTable := folderTableMgr.table

	q := qb.selectDataset().Prepared(true).Where(
		folderTable.Col("path").Eq(dir),
		table.Col("basename").Eq(basename),
	)

	ret, err := qb.get(ctx, q)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("getting folder by path %s: %w", p, err)
	}

	return ret, nil
}

func (qb *FileStore) findBySubquery(ctx context.Context, sq *goqu.SelectDataset) ([]file.File, error) {
	table := qb.table()

	q := qb.selectDataset().Prepared(true).Where(
		table.Col(idColumn).Eq(
			sq,
		),
	)

	return qb.getMany(ctx, q)
}

func (qb *FileStore) FindByFingerprint(ctx context.Context, fp file.Fingerprint) ([]file.File, error) {
	fingerprintTable := fingerprintTableMgr.table

	filesFingerprints := filesFingerprintsJoinTable.As("ff")
	fingerprints := fingerprintTable.As("fp")

	sq := dialect.From(filesFingerprints).Select(filesFingerprints.Col(fileIDColumn)).LeftJoin(
		fingerprints,
		goqu.On(filesFingerprints.Col(fingerprintIDColumn).Eq(fingerprints.Col(idColumn))),
	).Where(
		fingerprints.Col("type").Eq(fp.Type),
		fingerprints.Col("fingerprint").Eq(fp.Fingerprint),
	)

	return qb.findBySubquery(ctx, sq)
}

func (qb *FileStore) MarkMissing(ctx context.Context, scanStartTime time.Time, scanPaths []string) (int, error) {
	now := time.Now()
	table := qb.table()
	folderTable := folderTableMgr.table

	var pathEx []exp.Expression
	for _, p := range scanPaths {
		pathEx = append(pathEx, folderTable.Col("path").Like(p+"%"))
	}

	q := dialect.Update(table).Prepared(true).Set(exp.Record{
		"missing_since": now,
	}).Where(
		table.Col("last_scanned").Lt(scanStartTime),
		table.Col("missing_since").IsNull(),
	)

	if len(pathEx) > 0 {
		sq := dialect.From(table).Select(table.Col(idColumn)).InnerJoin(
			folderTable,
			goqu.On(table.Col("parent_folder_id").Eq(folderTable.Col(idColumn))),
		).Where(
			goqu.Or(pathEx...),
		)

		q = q.Where(table.Col(idColumn).Eq(sq))
	}

	r, err := exec(ctx, q)
	if err != nil {
		return 0, fmt.Errorf("marking files as missing: %w", err)
	}

	n, _ := r.RowsAffected()
	return int(n), nil
}
