package lights

import (
	"fmt"
	"log"

	"github.com/2tvenom/golifx"
)

// Types
type Light struct {
	light  *golifx.Bulb
	name   string
	state  bool
	colour golifx.HSBK
	group  string
}
type LightResponse struct {
	Mac    string      `json:"mac"`
	Name   string      `json:"name"`
	State  bool        `json:"state"`
	Colour golifx.HSBK `json:"colour"`
	Group  string      `json:"group"`
}

type Lights []*Light

// package Global vars
var allLights Lights

// type functions

func (l *Light) refreshLight() error {
	state, err := l.light.GetColorState()
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("error getting bulb state: %v", err))
	}
	if l.colour != *state.Color {
		log.Printf("Light '%v' has new colour:\n %v. \n", l.name, state.Color.String())
		l.colour = *state.Color
	}
	if l.state != state.Power {
		log.Printf("Light '%v' has new power state: %v. \n", l.name, state.Power)
		l.state = state.Power
	}
	return nil
}

func (l Lights) FindLight(mac string) *Light {
	for i := 0; i < len(allLights); i++ {
		if mac == l[i].light.MacAddress() {
			return l[i]
		}
	}
	return nil
}

// load and list

func LoadLights() error {

	bulbs, err := golifx.LookupBulbs()
	if err != nil {
		return err
	}

	for _, b := range bulbs {
		state, err := b.GetColorState()
		if err != nil {
			return fmt.Errorf(fmt.Sprintf("error getting bulb state: %v", err))
		}
		group, err := b.GetGroup()
		if err != nil {
			return fmt.Errorf(fmt.Sprintf("error getting bulb group: %v", err))
		}
		newLight := Light{light: b, name: state.Label, state: state.Power, colour: *state.Color, group: group.Label}

		log.Printf("Found Light: %v", newLight)

		allLights = append(allLights, &newLight)
	}
	return nil
}

func ListAllLights() []LightResponse {
	list := []LightResponse{}
	for _, light := range allLights {
		list = append(list, LightResponse{Mac: light.light.MacAddress(), Name: light.name, State: light.state, Colour: light.colour})
	}
	return list
}

// refresh

func RefreshLights() error {

	log.Println("Refreshing lights.")
	bulbs, err := golifx.LookupBulbs()
	if err != nil {
		return err
	}
	/* Check for new lights */
	for _, b := range bulbs {
		state, err := b.GetColorState()
		if err != nil {
			log.Printf("error getting bulb state: %v, Skipping update. ", err)
			continue // Skip bulb for now
		}
		group, groupErr := b.GetGroup()
		if groupErr != nil {
			log.Printf("error getting bulb group: %v", err)
		}
		if light := allLights.FindLight(b.MacAddress()); light != nil {
			if groupErr == nil && light.group != group.Label {
				log.Printf("Light '%v' has new group name: %v. \n", light.name, group.Label)
				light.group = group.Label
			}
			if light.colour != *state.Color {
				log.Printf("Light '%v' has new colour: %v. \n", light.name, state.Color.String())
				light.colour = *state.Color
			}
			if light.state != state.Power {
				log.Printf("Light '%v' has new power state: %v. \n", light.name, state.Power)
				light.state = state.Power
			}

		} else {
			log.Printf("Found new light: %v. \n", b)
			newLight := Light{light: b, name: state.Label, state: state.Power, colour: *state.Color, group: group.Label}
			allLights = append(allLights, &newLight)
		}

	}

	return nil
}

// CRUD

func LightOff(mac string) error {
	light := allLights.FindLight(mac)
	if light != nil {
		return light.light.SetPowerState(false)
	}
	return fmt.Errorf("unable to find light for mac: %v", mac)
}

func AllLightsOff() error {
	var ret error = nil
	for _, b := range allLights {
		err := b.light.SetPowerState(false)
		if err != nil {
			ret = fmt.Errorf("error setting light power stage")
			// TODO log better
		}
	}
	return ret
}
func AllLightsOn() error {
	var ret error = nil
	for _, b := range allLights {
		err := b.light.SetPowerState(true)
		if err != nil {
			ret = fmt.Errorf("error setting light power stage")
			// TODO log better
		}
	}
	return ret
}

func LightOn(mac string) error {
	light := allLights.FindLight(mac)
	if light != nil {
		return light.light.SetPowerState(true)
	}
	return fmt.Errorf("unable to find light for mac: %v", mac)
}

func SetColour(mac string, hue int, saturation int, brightness int, kelvin int, duration int) error {
	light := allLights.FindLight(mac)

	if light != nil {
		return light.light.SetColorState(&golifx.HSBK{
			Hue:        uint16(hue),
			Saturation: uint16(saturation),
			Kelvin:     uint16(kelvin),
			Brightness: uint16(brightness),
		}, uint32(duration))
	}
	return fmt.Errorf("unable to find light for mac: %v", mac)
}
