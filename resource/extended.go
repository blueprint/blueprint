package resource

import (
	"errors"
	"log"
	"net/http"

	"github.com/blueprint/blueprint"
)

// ExtendedResource represents a resource controller with CRUD operations.
type ExtendedResource struct {
	model    blueprint.Model
	resource *Resource
	ID       func(req *http.Request) (string, error)
}

// Extend extends a base resource.
//
// Optional way to get an ID from request.
// This is useful for '/user' routes based on session identification.
func (r *Resource) Extend(m blueprint.Model, optID ...func(req *http.Request) (string, error)) *ExtendedResource {
	e := &ExtendedResource{
		model: m,
		ID: func(req *http.Request) (string, error) {
			return PathID(req, m.PathID())
		},
		resource: r,
	}

	if len(optID) > 0 {
		e.ID = optID[0]
	}

	return e
}

// Index is a GET request for returning a list of items.
func (e *ExtendedResource) Index(resp http.ResponseWriter, req *http.Request) {
	baseID, err := e.resource.ID(req)
	if err != nil {
		resp.WriteErrs(req, err, errors.New("unable to identify resource id"))
		return
	}

	base := e.resource.model.New()
	if err := base.FindByID(baseID); err != nil {
		resp.WriteErrs(req, err)
		return
	}

	items, err := e.model.FindAllByOwner(base)
	if err != nil {
		log.Fatal(err)
	}

	resp.WriteFormatList(req, items)
}

// Show is a GET request for showing an item.
func (e *ExtendedResource) Show(resp http.ResponseWriter, req *http.Request) {
	baseID, err := e.resource.ID(req)
	if err != nil {
		resp.WriteErrs(req, err)
		return
	}

	base := e.resource.model.New()
	if err := base.FindByID(baseID); err != nil {
		resp.WriteErrs(req, err)
		return
	}

	id, err := e.ID(req)
	if err != nil {
		resp.WriteErrs(req, err)
		return
	}

	item := e.model.New()
	if err := item.FindByID(id); err != nil {
		resp.WriteErrs(req, err)
		return
	}

	// items might not store its relationships in a database, so BelongsTo() will
	// ensure an item and its relationships are show in the reponse
	if err := item.BelongsTo(base); err != nil {
		resp.WriteErrs(req, err)
		return
	}

	resp.WriteFormat(req, item)
}

// Store is a POST request for creating a new item.
func (e *ExtendedResource) Store(resp http.ResponseWriter, req *http.Request) {
	baseID, err := e.resource.ID(req)
	if err != nil {
		resp.WriteErrs(req, err)
		return
	}

	base := e.resource.model.New()
	if err := base.FindByID(baseID); err != nil {
		resp.WriteErrs(req, err)
		return
	}

	item := e.model.New()
	if err := req.Unmarshal(item); err != nil {
		resp.WriteErrs(req, err)
		return
	}

	if err := item.BelongsTo(base); err != nil {
		resp.WriteErrs(req, err)
		return
	}

	if err := item.Validate(); err != nil {
		resp.WriteErrs(req, err)
		return
	}

	if err := item.Save(); err != nil {
		log.Fatal(err)
	}

	resp.Status(http.StatusCreated)
	resp.WriteFormat(req, item)
}

// Apply is a PATCH request for updating a single item.
func (e *ExtendedResource) Apply(resp http.ResponseWriter, req *http.Request) {
	baseID, err := e.resource.ID(req)
	if err != nil {
		resp.WriteErrs(req, err)
		return
	}

	base := e.resource.model.New()
	if err := base.FindByID(baseID); err != nil {
		resp.WriteErrs(req, err)
		return
	}

	id, err := e.ID(req)
	if err != nil {
		resp.WriteErrs(req, err)
		return
	}

	item := e.model.New()
	if err := item.FindByID(id); err != nil {
		resp.WriteErrs(req, err)
		return
	}

	// set the values specified in the PATCH body
	if err := req.Unmarshal(item); err != nil {
		resp.WriteErrs(req, err)
		return
	}

	// ensure all relationships are valid before validation
	if err := item.BelongsTo(base); err != nil {
		resp.WriteErrs(req, err)
		return
	}

	if err := item.Validate(); err != nil {
		resp.WriteErrs(req, err)
		return
	}

	if err := item.Save(); err != nil {
		resp.WriteErrs(req, err)
		return
	}

	resp.WriteFormat(req, item)
}

// Update is a PUT request for replacing a single item.
func (e *ExtendedResource) Update(resp http.ResponseWriter, req *http.Request) {
	baseID, err := e.resource.ID(req)
	if err != nil {
		resp.WriteErrs(req, err)
		return
	}

	base := e.resource.model.New()
	if err := base.FindByID(baseID); err != nil {
		resp.WriteErrs(req, err)
		return
	}

	// get ID in the URL path
	id, err := e.ID(req)
	if err != nil {
		resp.WriteErrs(req, err)
		return
	}

	item := e.model.New()
	if err := item.FindByID(id); err != nil {
		resp.WriteErrs(req, err)
		return
	}

	if err := req.Unmarshal(item); err != nil {
		resp.WriteErrs(req, err)
		return
	}

	if err := item.SetID(id); err != nil {
		resp.WriteErrs(req, err)
		return
	}

	// items might not store its relationships in a database, so BelongsTo() will
	// ensure an item and its relationships are show in the reponse
	if err := item.BelongsTo(base); err != nil {
		resp.WriteErrs(req, err)
		return
	}

	if err := item.Validate(); err != nil {
		resp.WriteErrs(req, err)
		return
	}

	if err := item.Save(); err != nil {
		resp.WriteErrs(req, err)
		return
	}

	resp.WriteFormat(req, item)
}

// Destroy is a DELETE request for deleting a single item.
func (e *ExtendedResource) Destroy(resp http.ResponseWriter, req *http.Request) {
	baseID, err := e.resource.ID(req)
	if err != nil {
		resp.WriteErrs(req, err)
		return
	}

	base := e.resource.model.New()
	if err := base.FindByID(baseID); err != nil {
		resp.WriteErrs(req, err)
		return
	}

	// get ID in URL path
	id, err := e.ID(req)
	if err != nil {
		resp.WriteErrs(req, err)
		return
	}

	item := e.model.New()
	if err := item.FindByID(id); err != nil {
		resp.WriteErrs(req, err)
		return
	}

	// we must to be aware of an item's relationships when deleting it
	if err := item.BelongsTo(base); err != nil {
		resp.WriteErrs(req, err)
		return
	}

	if err := item.Delete(); err != nil {
		resp.WriteErrs(req, err)
		return
	}

	resp.Status(http.StatusNoContent)
	resp.WriteFormat(req, "")
}
