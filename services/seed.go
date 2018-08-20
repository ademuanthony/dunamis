package services

import (
	"github.com/ademuanthony/dunamis/app"
	"github.com/ademuanthony/dunamis/models"
	"time"
	"github.com/gocolly/colly"
	"fmt"
	"strings"
	"github.com/ademuanthony/dunamis/util"
)

// seedDao specifies the interface of the seed DAO needed by SeedService.
type seedDao interface {
	// Get returns the seed with the specified seed ID.
	Get(rs app.RequestScope, id int) (*models.Seed, error)
	// Get returns the seed for the specified day, month and year.
	GetByDay(rs app.RequestScope, day int, month time.Month, year int) (*models.Seed, error)
	// Count returns the number of seed.
	Count(rs app.RequestScope) (int, error)
	// Query returns the list of seeds with the given offset and limit.
	Query(rs app.RequestScope, offset, limit int) ([]models.Seed, error)
	// Create saves a new seed in the storage.
	Create(rs app.RequestScope, artist *models.Seed) error
	// Update updates the seed with given ID in the storage.
	Update(rs app.RequestScope, id int, artist *models.Seed) error
}

// SeedService provides services related with seed.
type SeedService struct {
	dao seedDao
}

// NewSeedService creates a new SeedService with the given seed DAO.
func NewSeedService(dao seedDao) *SeedService {
	return &SeedService{dao}
}

// Get returns the artist with the specified the artist ID.
func (s *SeedService) Get(rs app.RequestScope, id int) (*models.Seed, error) {
	return s.dao.Get(rs, id)
}


// GetByDay returns the seed for the specified day, month and year from the db or
// the http://www.dunamisgospel.org website.
func (s *SeedService) GetByDay(rs app.RequestScope, day int, month time.Month, year int) (*models.Seed, error) {
	seed, err := s.dao.GetByDay(rs, day,month,year)
	if err != nil {
		err = s.scrapSeeds(rs)
		if err != nil {
			return seed, err
		}
		seed, err = s.dao.GetByDay(rs, day,month,year)
	}
	return seed, err
}

// GetForNextThreeDays returns the seeds for the next three days starting from the specified day, month and year from the db or
// the http://www.dunamisgospel.org website.
func (s *SeedService) GetForNextThreeDays(rs app.RequestScope, day int, month time.Month, year int) ([]*models.Seed, error) {
	var seeds []*models.Seed
	date := time.Date(year, month, day, 0, 0,0, 0, util.DefaultLocation)

	//first day
	seed1, err := s.GetByDay(rs, date.Day(), date.Month(), date.Year())
	if err != nil {
		return seeds, err
	}
	seeds = append(seeds, seed1)

	//second day
	date = date.AddDate(0,0,1)
	seed2, err := s.GetByDay(rs, date.Day(), date.Month(), date.Year())
	if err != nil {
		return seeds, err
	}
	seeds = append(seeds, seed2)

	//third day
	date = date.AddDate(0,0,1)
	seed3, err := s.GetByDay(rs, date.Day(), date.Month(), date.Year())
	if err != nil {
		return seeds, err
	}
	seeds = append(seeds, seed3)

	return seeds, err
}

//scrapSeeds scraps latest seed of destiny and save to the storage
func (s *SeedService) scrapSeeds(rs app.RequestScope,) error {
	links, lastErr := s.getSodLinks()
	if lastErr != nil {
		fmt.Print(lastErr)
		return lastErr
	}
	for title, link := range links {
		seed, err := s.getSodFromLink("http://www.dunamisgospel.org"+link, title)
		if err != nil {
			fmt.Println(err)
			lastErr = err
		} else {
			_, lastErr = s.Create(rs, &seed)

		}
	}


	return lastErr
}

//getSodLinks retrieves latest seed of destiny links from the site
func (s *SeedService) getSodLinks() (map[string]string, error){
	links := map[string]string{}

	var err error
	categoryLink := "http://www.dunamisgospel.org/index.php/component/k2/itemlist/category/3"
	c := colly.NewCollector()

	c.OnHTML("#itemListLeading .itemContainer", func(element *colly.HTMLElement) {
		hrefs := element.ChildAttrs("h3.catItemTitle a", "href")
		title := element.ChildText("h3.catItemTitle a")
		if len(hrefs) > 0 {
			links[title] = hrefs[0]
		}
	})

	c.OnError(func(response *colly.Response, e error) {
		err = e
	})

	c.Visit(categoryLink)

	return links, err
}

// getSodFromLink scraps a single sod fro the specified url
func (s *SeedService) getSodFromLink(link string, title string) (models.Seed, error) {
	seed := models.Seed{Title:title, Paragraphs:[]models.Paragraph{}}
	var err error

	c := colly.NewCollector()

	selector := "#k2Container.row.itemView.clearfix div.itemViewContent.span8.pull-right div.itemBody div.itemFullText div"

	index := -1
	var paragraphs []models.Paragraph
	c.OnHTML(selector, func(element *colly.HTMLElement) {
		text := strings.TrimSpace(element.Text)

		if text != "&nbsp;" && text != ""{
			index++
			switch index {
			case 1:
				dateString := strings.Join(strings.Split(text, " ")[1:], " ")
				date, err := util.ParseOrdinalDate("31 JANUARY 2006",
					dateString)
				if err == nil {
					seed.Date = text
					seed.Day = date.Day()
					seed.Month = int(date.Month())
					seed.Year = date.Year()
					return
				} else {
					fmt.Println(err, text)
				}
				break
			case 2:
				seed.Scripture = text
				break
			case 3:
				seed.Thought = text
				break
			case 10:
				seed.RememberThis = text
				break
			case 11:
				seed.Assignment = text
				break
			case 12:
				seed.Prayer = text
				break
			case 13:
				seed.Resource = text
				break
			case 14:
				seed.DailyReading = text
				break
			}

			paragraphs = append(paragraphs, models.Paragraph{Content: text, Type:models.ParagraphTypes.PlainText})
		}
		//seed.Content += text

	})

	//start scrapping
	c.Visit(link)

	for _, paragraph := range paragraphs {
		seed.Paragraphs = append(seed.Paragraphs, paragraph)
	}
	//spew.Dump(seed)

	return seed,err
}

// Create creates a new seed.
func (s *SeedService) Create(rs app.RequestScope, model *models.Seed) (*models.Seed, error) {
	if err := model.Validate(); err != nil {
		return nil, err
	}
	if err := s.dao.Create(rs, model); err != nil {
		return nil, err
	}
	return s.dao.Get(rs, model.Id)
}

// Update updates the seed with the specified ID.
func (s *SeedService) Update(rs app.RequestScope, id int, model *models.Seed) (*models.Seed, error) {
	if err := model.Validate(); err != nil {
		return nil, err
	}
	if err := s.dao.Update(rs, id, model); err != nil {
		return nil, err
	}
	return s.dao.Get(rs, id)
}

// Count returns the number of seed.
func (s *SeedService) Count(rs app.RequestScope) (int, error) {
	return s.dao.Count(rs)
}

// Query returns the seeds with the specified offset and limit.
func (s *SeedService) Query(rs app.RequestScope, offset, limit int) ([]models.Seed, error) {
	return s.dao.Query(rs, offset, limit)
}
