package service

import (
	"encoding/json"
	"net/http"

	"github.com/Mukam21/server_Golang/pkg/config"
	"github.com/Mukam21/server_Golang/pkg/model"
	"github.com/sirupsen/logrus"
)

type Repository interface {
	Create(person *model.Person) (int64, error)
	GetByID(id int64) (*model.Person, error)
	GetAll(page, limit int, filters map[string]string) ([]*model.Person, error)
	Update(person *model.Person) error
	Patch(id int64, patch *model.PersonPatchRequest) error
	Delete(id int64) error
}

type Service struct {
	repo Repository
	log  *logrus.Logger
	cfg  *config.Config
}

func NewService(repo Repository, log *logrus.Logger, cfg *config.Config) *Service {
	return &Service{repo: repo, log: log, cfg: cfg}
}

func (s *Service) CreatePerson(req *model.PersonRequest) (*model.Person, error) {
	person := &model.Person{
		Name:       req.Name,
		Surname:    req.Surname,
		Patronymic: req.Patronymic,
	}

	if age, err := s.getAge(person.Name); err == nil {
		person.Age = &age
	} else {
		s.log.Debugf("Failed to get age for %s: %v", person.Name, err)
	}
	if gender, err := s.getGender(person.Name); err == nil {
		person.Gender = &gender
	} else {
		s.log.Debugf("Failed to get gender for %s: %v", person.Name, err)
	}
	if nationality, err := s.getNationality(person.Name); err == nil {
		person.Nationality = &nationality
	} else {
		s.log.Debugf("Failed to get nationality for %s: %v", person.Name, err)
	}

	id, err := s.repo.Create(person)
	if err != nil {
		s.log.Errorf("Failed to create person: %v", err)
		return nil, err
	}
	person.ID = id

	s.log.Infof("Created person with ID: %d", id)
	return person, nil
}

func (s *Service) GetByID(id int64) (*model.Person, error) {
	person, err := s.repo.GetByID(id)
	if err != nil {
		s.log.Errorf("Failed to get person with ID %d: %v", id, err)
		return nil, err
	}
	return person, nil
}

func (s *Service) GetAll(page, limit int, filters map[string]string) ([]*model.Person, error) {
	persons, err := s.repo.GetAll(page, limit, filters)
	if err != nil {
		s.log.Errorf("Failed to get persons: %v", err)
		return nil, err
	}
	s.log.Infof("Retrieved %d persons", len(persons))
	return persons, nil
}

func (s *Service) Update(person *model.Person) error {
	if err := s.repo.Update(person); err != nil {
		s.log.Errorf("Failed to update person with ID %d: %v", person.ID, err)
		return err
	}
	s.log.Infof("Updated person with ID: %d", person.ID)
	return nil
}

func (s *Service) Patch(id int64, patch *model.PersonPatchRequest) error {
	if err := s.repo.Patch(id, patch); err != nil {
		s.log.Errorf("Failed to patch person with ID %d: %v", id, err)
		return err
	}
	s.log.Infof("Patched person with ID: %d", id)
	return nil
}

func (s *Service) Delete(id int64) error {
	if err := s.repo.Delete(id); err != nil {
		s.log.Errorf("Failed to delete person with ID %d: %v", id, err)
		return err
	}
	s.log.Infof("Deleted person with ID: %d", id)
	return nil
}

func (s *Service) getAge(name string) (int, error) {
	resp, err := http.Get(s.cfg.APIAgifyURL + "?name=" + name)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Age int `json:"age"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}
	return result.Age, nil
}

func (s *Service) getGender(name string) (string, error) {
	resp, err := http.Get(s.cfg.APIGenderizeURL + "?name=" + name)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Gender string `json:"gender"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.Gender != "male" && result.Gender != "female" {
		return "other", nil
	}
	return result.Gender, nil
}

func (s *Service) getNationality(name string) (string, error) {
	resp, err := http.Get(s.cfg.APINationalizeURL + "?name=" + name)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Country []struct {
			CountryID   string  `json:"country_id"`
			Probability float64 `json:"probability"`
		} `json:"country"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Country) == 0 {
		return "", nil
	}

	maxProb := 0.0
	var maxCountry string
	for _, country := range result.Country {
		if country.Probability > maxProb {
			maxProb = country.Probability
			maxCountry = country.CountryID
		}
	}
	return maxCountry, nil
}
