package daos

import (
	"github.com/ademuanthony/dunamis/app"
	"github.com/ademuanthony/dunamis/models"
	"time"
	"github.com/go-ozzo/ozzo-dbx"
)

// SeedDAO persists artist data in database
type SeedDAO struct{}

// NewSeedDAO creates a new SeedDAO
func NewSeedDAO() *SeedDAO {
	return &SeedDAO{}
}

// Get reads the seed with the specified ID from the database.
func (dao *SeedDAO) Get(rs app.RequestScope, id int) (*models.Seed, error) {
	var seed models.Seed
	err := rs.Tx().Select().Model(id, &seed)

	if err != nil {
		return &seed, err
	}
	var paragraphs []models.Paragraph
	err = rs.Tx().Select().Where(dbx.HashExp{"seed_id": seed.Id}).All(&paragraphs)
	seed.Paragraphs = paragraphs

	return &seed, err
}

// GetByDay reads the seed for the specified day, month and year from the database.
func (dao *SeedDAO) GetByDay(rs app.RequestScope, day int, month time.Month, year int) (*models.Seed, error) {
	var seed models.Seed
	err := rs.Tx().Select().Where(dbx.HashExp{"day": day, "month": int(month), "year": year}).One(&seed)

	if err != nil {
		return &seed, err
	}
	var paragraphs []models.Paragraph
	err = rs.Tx().Select().Where(dbx.HashExp{"seed_id": seed.Id}).All(&paragraphs)
	seed.Paragraphs = paragraphs

	return &seed, err
}

// Create saves a new seed record in the database. or updates if the record already exists
// The Seed.Id field will be populated with an automatically generated ID upon successful saving.
func (dao *SeedDAO) Create(rs app.RequestScope, seed *models.Seed) error {
	oldSeed, err := dao.GetByDay(rs, seed.Day, time.Month(seed.Month), seed.Year)
	if err == nil {
		return dao.Update(rs, oldSeed.Id, seed)
	}
	seed.Id = 0
	err = rs.Tx().Model(seed).Insert()
	if err != nil {
		return err
	}
	for _, paragraph := range seed.Paragraphs{
		paragraph.Id = 0
		paragraph.SeedId = seed.Id
		err = rs.Tx().Model(&paragraph).Insert()
		if err != nil {
			return err
		}
	}
	return err
}

// Update saves the changes to an seed in the database.
func (dao *SeedDAO) Update(rs app.RequestScope, id int, seed *models.Seed) error {
	if _, err := dao.Get(rs, id); err != nil {
		return err
	}
	seed.Id = id
	return rs.Tx().Model(seed).Exclude("id").Update()
}

// Count returns the number of the seed records in the database.
func (dao *SeedDAO) Count(rs app.RequestScope) (int, error) {
	var count int
	err := rs.Tx().Select("COUNT(*)").From("seed").Row(&count)
	return count, err
}

// Query retrieves the seed records with the specified offset and limit from the database.
func (dao *SeedDAO) Query(rs app.RequestScope, offset, limit int) ([]models.Seed, error) {
	seeds := []models.Seed{}
	err := rs.Tx().Select().OrderBy("id").Offset(int64(offset)).Limit(int64(limit)).All(&seeds)
	return seeds, err
}
