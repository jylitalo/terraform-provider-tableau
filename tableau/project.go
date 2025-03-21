package tableau

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Project struct {
	ID                 string `json:"id,omitempty"`
	Name               string `json:"name,omitempty"`
	ParentProjectID    string `json:"parentProjectId,omitempty"`
	Description        string `json:"description,omitempty"`
	ContentPermissions string `json:"contentPermissions,omitempty"`
	Owner              Owner  `json:"owner,omitempty"`
}

type ProjectRequest struct {
	Project Project `json:"project"`
}

type ProjectResponse struct {
	Project Project `json:"project"`
}

type ProjectsResponse struct {
	Projects []Project `json:"project"`
}

type ProjectListResponse struct {
	ProjectsResponse ProjectsResponse  `json:"projects"`
	Pagination       PaginationDetails `json:"pagination"`
}

func (c *Client) GetProjects() ([]Project, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/projects", c.ApiUrl), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	projectListResponse := ProjectListResponse{}
	err = json.Unmarshal(body, &projectListResponse)
	if err != nil {
		return nil, err
	}

	pageNumber, totalPageCount, totalAvailable, err := GetPaginationNumbers(projectListResponse.Pagination)
	if err != nil {
		return nil, err
	}

	allProjects := make([]Project, 0, totalAvailable)
	allProjects = append(allProjects, projectListResponse.ProjectsResponse.Projects...)

	for page := pageNumber + 1; page <= totalPageCount; page++ {
		fmt.Printf("Searching page %d", page)
		req, err = http.NewRequest("GET", fmt.Sprintf("%s/projects?pageNumber=%d", c.ApiUrl, page), nil)
		if err != nil {
			return nil, err
		}
		body, err = c.doRequest(req)
		if err != nil {
			return nil, err
		}
		projectListResponse = ProjectListResponse{}
		err = json.Unmarshal(body, &projectListResponse)
		if err != nil {
			return nil, err
		}
		allProjects = append(allProjects, projectListResponse.ProjectsResponse.Projects...)
	}

	return allProjects, nil
}

func (c *Client) GetProject(projectID string) (*Project, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/projects", c.ApiUrl), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	projectListResponse := ProjectListResponse{}
	err = json.Unmarshal(body, &projectListResponse)
	if err != nil {
		return nil, err
	}

	// TODO: Generalise pagination handling and use elsewhere
	pageNumber, totalPageCount, _, err := GetPaginationNumbers(projectListResponse.Pagination)
	if err != nil {
		return nil, err
	}
	for i, project := range projectListResponse.ProjectsResponse.Projects {
		if project.ID == projectID {
			return &projectListResponse.ProjectsResponse.Projects[i], nil
		}
	}

	for page := pageNumber + 1; page <= totalPageCount; page++ {
		fmt.Printf("Searching page %d", page)
		req, err = http.NewRequest("GET", fmt.Sprintf("%s/projects?pageNumber=%d", c.ApiUrl, page), nil)
		if err != nil {
			return nil, err
		}
		body, err = c.doRequest(req)
		if err != nil {
			return nil, err
		}
		projectListResponse = ProjectListResponse{}
		err = json.Unmarshal(body, &projectListResponse)
		if err != nil {
			return nil, err
		}
		// Check page for project match
		for i, project := range projectListResponse.ProjectsResponse.Projects {
			if project.ID == projectID {
				return &projectListResponse.ProjectsResponse.Projects[i], nil
			}
		}
	}

	return nil, fmt.Errorf("Did not find project ID %s", projectID)
}

func (c *Client) CreateProject(name, parentProjectId, description, contentPermissions, ownerId string) (*Project, error) {

	newProject := Project{
		Name:               name,
		ParentProjectID:    parentProjectId,
		Description:        description,
		ContentPermissions: contentPermissions,
	}
	if ownerId != "" {
		newOwner := Owner{
			ID: ownerId,
		}
		newProject.Owner = newOwner
	}
	projectRequest := ProjectRequest{
		Project: newProject,
	}

	newProjectJson, err := json.Marshal(projectRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/projects", c.ApiUrl), strings.NewReader(string(newProjectJson)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	projectResponse := ProjectResponse{}
	err = json.Unmarshal(body, &projectResponse)
	if err != nil {
		return nil, err
	}

	time.Sleep(1 * time.Second)
	return &projectResponse.Project, nil
}

func (c *Client) UpdateProject(projectID, name, parentProjectId, description, contentPermissions, ownerId string) (*Project, error) {
	newOwner := Owner{
		ID: ownerId,
	}
	newProject := Project{
		Name:               name,
		ParentProjectID:    parentProjectId,
		Description:        description,
		ContentPermissions: contentPermissions,
		Owner:              newOwner,
	}
	projectRequest := ProjectRequest{
		Project: newProject,
	}

	newProjectJson, err := json.Marshal(projectRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/projects/%s", c.ApiUrl, projectID), strings.NewReader(string(newProjectJson)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	projectResponse := ProjectResponse{}
	err = json.Unmarshal(body, &projectResponse)
	if err != nil {
		return nil, err
	}

	return &projectResponse.Project, nil
}

func (c *Client) DeleteProject(projectID string) error {

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/projects/%s", c.ApiUrl, projectID), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
