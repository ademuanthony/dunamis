package apis

import (
	"github.com/ademuanthony/dunamis/app"
	"github.com/ademuanthony/dunamis/models"
	"time"
	"github.com/go-ozzo/ozzo-routing"
	"strconv"
)

type (
	// seedService specifies the interface for the artist service needed by seedResource.
	seedService interface {
		Get(rs app.RequestScope, id int) (*models.Seed, error)
		GetByDay(rs app.RequestScope, day int, month time.Month, year int) (*models.Seed, error)
		GetForNextThreeDays(rs app.RequestScope, day int, month time.Month, year int) ([]*models.Seed, error)
		Query(rs app.RequestScope, offset, limit int) ([]models.Seed, error)
		Count(rs app.RequestScope) (int, error)
		Create(rs app.RequestScope, model *models.Seed) (*models.Seed, error)
		Update(rs app.RequestScope, id int, model *models.Seed) (*models.Seed, error)
	}

	// seedResource defines the handlers for the CRUD APIs.
	seedResource struct {
		service seedService
	}
)

// ServeSeed sets up the routing of seed endpoints and the corresponding handlers.
func ServeSeedResource(rg *routing.RouteGroup, service seedService) {
	r := &seedResource{service}
	rg.Get("/seeds/<id>", r.get)
	rg.Get("/seeds/<day>/<month>/<year>", r.getByDay)
	rg.Get("/seeds/next-three/<day>/<month>/<year>", r.getForNextThreeDays)
	rg.Get("/seeds", r.query)
	rg.Put("/seeds/<id>", r.update)
}

func (r *seedResource) get(c *routing.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	response, err := r.service.Get(app.GetRequestScope(c), id)
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *seedResource) getByDay(c *routing.Context) error {
	day, err := strconv.Atoi(c.Param("day"))
	if err != nil {
		return err
	}

	month, err := strconv.Atoi(c.Param("month"))
	if err != nil {
		return err
	}
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		return err
	}

	response, err := r.service.GetByDay(app.GetRequestScope(c), day, time.Month(month), year)
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *seedResource) getForNextThreeDays(c *routing.Context) error {
	day, err := strconv.Atoi(c.Param("day"))
	if err != nil {
		return err
	}

	month, err := strconv.Atoi(c.Param("month"))
	if err != nil {
		return err
	}
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		return err
	}

	response, err := r.service.GetForNextThreeDays(app.GetRequestScope(c), day, time.Month(month), year)
	if err != nil && len(response) == 0{
		return err
	}

	return c.Write(response)
}

func (r *seedResource) query(c *routing.Context) error {
	rs := app.GetRequestScope(c)
	count, err := r.service.Count(rs)
	if err != nil {
		return err
	}
	paginatedList := getPaginatedListFromRequest(c, count)
	items, err := r.service.Query(app.GetRequestScope(c), paginatedList.Offset(), paginatedList.Limit())
	if err != nil {
		return err
	}
	paginatedList.Items = items
	return c.Write(paginatedList)
}


func (r *seedResource) update(c *routing.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	rs := app.GetRequestScope(c)

	model, err := r.service.Get(rs, id)
	if err != nil {
		return err
	}

	if err := c.Read(model); err != nil {
		return err
	}

	response, err := r.service.Update(rs, id, model)
	if err != nil {
		return err
	}

	return c.Write(response)
}

