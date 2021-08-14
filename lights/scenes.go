package lights

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/2tvenom/golifx"
	"gopkg.in/yaml.v3"
)

// Types

type Scene struct {
	ID               string        `yaml:"id"`
	SceneName        string        `yaml:"name"`
	SceneDescription string        `yaml:"description"`
	Actions          []lightAction `yaml:"actions"`
	Order            int           `yaml:"order"`
}

type SceneResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Order       int      `json:"order"`
	Actions     []string `json:"actions"`
}

type SceneSaveRequest struct {
	ID          string   `json:"ID"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Actions     []string `json:"actions"`
	Order       int      `json:"order"`
}

// TODO, maybe expose this too so we can have actions come straight from the api?
type lightAction struct {
	Mac        string `yaml:"mac"`
	State      bool   `yaml:"state"`
	Brightness uint16 `yaml:"brightness"`
	Hue        uint16 `yaml:"hue"`
	Saturation uint16 `yaml:"saturation"`
	Kelvin     uint16 `yaml:"kelvin"`
}

type Scenes []*Scene

// package global scene vars

var allScenes Scenes

var scenesPath string

// init

func InitScenes(configPath string) error {
	scenesPath = configPath + "/scenes.yaml"
	return LoadScenes()
}

// type functions

func (s Scene) runScene() error {
	log.Printf("Running scene %v \n", s.SceneName)

	for _, a := range s.Actions {
		light := allLights.FindLight(a.Mac)
		if light != nil {
			err := light.light.SetPowerState(a.State)
			if err != nil {
				log.Printf("Unable to set '%v' state on light '%v'", a.State, a.Mac)
			}

			newColour := &golifx.HSBK{
				Hue:        a.Hue,
				Saturation: a.Saturation,
				Brightness: a.Brightness,
				Kelvin:     a.Kelvin,
			}
			err = light.light.SetColorState(newColour, 0)
			if err != nil {
				log.Printf("Unable to set '%v' colour on light '%v'", newColour, a.Mac)
			}
		} else {
			log.Printf("Unable to find light '%v'", a.Mac)
		}
	}

	return nil

}

// search

func findScene(id string) *Scene {
	for _, s := range allScenes {
		if s.ID == id {
			return s
		}
	}
	return nil
}

func TriggerScene(id string) error {
	for _, s := range allScenes {
		if s.ID == id {
			return s.runScene()
		}
	}
	return nil
}

// Load and list

func LoadScenes() error {

	log.Printf("Loading scenes from %v\n", scenesPath)
	scenes, err := ioutil.ReadFile(scenesPath)
	if err != nil {
		return err
	}

	var data []Scene

	err = yaml.Unmarshal(scenes, &data)
	if err != nil {
		return err
	}

	for _, s := range data {
		allScenes = append(allScenes, &Scene{
			ID:               s.ID,
			SceneName:        s.SceneName,
			SceneDescription: s.SceneDescription,
			Actions:          s.Actions,
		})

	}

	log.Println("Loaded Scenes: ")
	for _, scene := range allScenes {
		log.Printf("%v", *scene)
	}
	return nil
}

func ListAllScenes() []SceneResponse {
	scenes := []SceneResponse{}
	for _, s := range allScenes {
		var actionList []string
		for _, a := range s.Actions {
			actionList = append(actionList, a.Mac)
		}
		scenes = append(scenes, SceneResponse{ID: s.ID, Name: s.SceneName, Description: s.SceneDescription, Order: s.Order, Actions: actionList})
	}
	return scenes
}

// CRUD

// TODO Rename this to New
func SaveScene(req SceneSaveRequest) error {
	newScene := Scene{
		ID:               generateID(),
		SceneName:        req.Name,
		SceneDescription: req.Description,
		Actions:          []lightAction{},
		Order:            req.Order,
	}

	if newScene.SceneName == "" {
		return fmt.Errorf("scene name must not be empty")
	}

	for _, a := range req.Actions {
		if a == "" {
			continue
		}
		light := allLights.FindLight(a)
		if light == nil {
			return fmt.Errorf("unable to find light with mac: '%v'", a)
		}

		err := light.refreshLight()
		if err != nil {
			log.Printf("Unable to refresh light '%v' , using old values", light)
		}

		newAction := lightAction{
			Mac:        a,
			State:      light.state,
			Brightness: light.colour.Brightness,
			Hue:        light.colour.Hue,
			Saturation: light.colour.Saturation,
			Kelvin:     light.colour.Kelvin,
		}

		newScene.Actions = append(newScene.Actions, newAction)

	}

	allScenes = append(allScenes, &newScene)

	saveScenes()
	return nil
}

func ModifyScene(req SceneSaveRequest) error {

	scene := findScene(req.ID)

	log.Printf("Trying to modify: $v", scene)

	if scene == nil {
		return fmt.Errorf("unable to find scene with name: %v", req.Name)
	}

	scene.Order = req.Order
	scene.SceneName = req.Name
	scene.SceneDescription = req.Description
	// Wipe out the actions and refill them
	scene.Actions = nil

	for _, a := range req.Actions {
		if a == "" {
			continue
		}
		light := allLights.FindLight(a)
		if light == nil {
			log.Printf("unable to find light with mac: '%v'", a)
			continue
		}

		err := light.refreshLight()
		if err != nil {
			log.Printf("Unable to refresh light %v , using old values", light)
		}

		newAction := lightAction{
			Mac:        a,
			State:      light.state,
			Brightness: light.colour.Brightness,
			Hue:        light.colour.Hue,
			Saturation: light.colour.Saturation,
			Kelvin:     light.colour.Kelvin,
		}

		scene.Actions = append(scene.Actions, newAction)

	}
	saveScenes()
	return nil
}

func DeleteScene(id string) error {
	newScenes := Scenes{}

	for i := 0; i < len(allScenes); i++ {
		if allScenes[i].ID != id {
			newScenes = append(newScenes, allScenes[i])
		}
	}
	allScenes = newScenes
	saveScenes()
	return fmt.Errorf("Scene id:'%v' Not found", id)
}

// Save

func saveScenes() error {
	log.Printf("Attempting to save scenes")
	data, err := yaml.Marshal(allScenes)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(scenesPath, data, 0644)
	if err != nil {
		return err
	}

	log.Printf("Saved Scenes: %v ", allScenes)
	return nil
}
