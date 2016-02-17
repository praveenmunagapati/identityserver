package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

	resource "github.com/itsyouonline/identityserver/http"
	companyModel "github.com/itsyouonline/identityserver/identityservice/company/models"
)

type CompanyDetails struct {
	companyModel.Company
	Uri string `json:"uri"`
}

type CompanyResource struct {
	resource.Resource
}

func NewCompanyResource() *CompanyResource {
	c := &CompanyResource{}
	c.ResourceHandler = c

	return c
}

func (c *CompanyResource) GetRoutes() resource.Routes {
	routes := resource.Routes{
		resource.Route{
			Name: "CompanyList",
			Methods: resource.RouteMethods{
				resource.POST,
			},
			Path:        "/companies/",
			HandlerFunc: c.DispatchList,
		},
		resource.Route{
			Name: "CompanyDetail",
			Methods: resource.RouteMethods{
				resource.GET,
				resource.PUT,
			},
			Path:        "/companies/{globalId}/",
			HandlerFunc: c.DispatchDetail,
		},
	}

	return routes
}

func (c *CompanyResource) PostList(w http.ResponseWriter, r *http.Request) {
	company, err := c.deserialize(r)
	if err != nil {
		log.Debug(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = company.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := c.serialize(company)

	c.Respond(w, response)
}

func (c *CompanyResource) GetDetail(w http.ResponseWriter, r *http.Request) {
	companyMgr := companyModel.NewCompanyManager(r)

	globalId := mux.Vars(r)["globalId"]

	company, err := companyMgr.GetByName(globalId)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	response := c.serialize(company)

	c.Respond(w, response)
}

func (c *CompanyResource) PutDetail(w http.ResponseWriter, r *http.Request) {
	globalId := mux.Vars(r)["globalId"]

	company, err := c.deserialize(r)
	if err != nil {
		log.Debug(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	companyMgr := companyModel.NewCompanyManager(r)

	oldCompany, cerr := companyMgr.GetByName(globalId)
	if cerr != nil {
		log.Debug(cerr)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if company.GlobalId != globalId || company.GetId() != oldCompany.GetId() {
		http.Error(w, "Changing globalId or id is Forbidden!", http.StatusForbidden)
		return
	}

	err = company.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := c.serialize(company)

	c.Respond(w, response)
}

func (c *CompanyResource) serialize(company *companyModel.Company) *CompanyDetails {
	uri := c.resourceUri(company)

	return &CompanyDetails{
		*company,
		uri,
	}
}

func (c *CompanyResource) deserialize(r *http.Request) (*companyModel.Company, error) {
	company := companyModel.NewCompany(r)

	data, err := ioutil.ReadAll(io.LimitReader(r.Body, 4096))
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	if err := json.Unmarshal(data, company); err != nil {
		return nil, err
	}

	return company, nil
}

func (c *CompanyResource) resourceUri(company *companyModel.Company) string {
	return c.BuildUri("/companies/", company.GlobalId)
}
