## Simple api for interacting with lifx lights over the local network 

Intended to be used as a dockerised service and probably interacted with from a web app.

# Running

Running from docker is tested as working on a linux host. 
Example in `docker-compose.yml`

A scenes.yaml must exist alongside the binary or at a path given by env `CONFIG_PATH`

`cp ./scenes.yaml.example ./scenes.yaml`



# Api 

Examples of use

```typescript
/**
 * shared scene item across the things
 */

interface sceneItem {
  name: string
  description: string
  order: number
  actions: string[]
  id: string
}
/**
 * shared light item across the things
 */
interface lightItem {
  mac: string
  name: string
  state: boolean
  group: string
  colour: {
    kelvin: number
    brightness: number
    hue: number
    saturation: number
  }
}

//const RecieverUrl = "http://home.phcyso.com/lights/"
const RecieverUrl = "/lights"
function refreshScenes() {
  helpers
    .fetchWithTimeout(`${RecieverUrl}/scenes/list`, { timeout: 5000 })
    .then((response) => response.json())
    .then((rawData) => {
      let castData = rawData as sceneItem[]
      castData.sort((a, b) => (a.order > b.order ? 1 : 0))
      scenes.value = castData
      console.log(scenes)
    })
}
helpers
  .fetchWithTimeout(`${RecieverUrl}/lights/list`, { timeout: 5000 })
  .then((response) => response.json())
  .then((data) => {
    lights.value = data
  })

function powerChange(state: boolean) {
  let url = `${RecieverUrl}/lights/all/off`
  if (state) {
    url = `${RecieverUrl}/lights/all/on`
  }
  fetch(url)
    .then((response) => response.json())
    .then((data) => {
      console.log(data)
    })
}

function lightOff(mac: string) {
  fetch(`${RecieverUrl}/light/off/${mac}`)
    .then((response) => response.json())
    .then((data) => {
      console.log(data)
    })
}
function lightOn(mac: string) {
  fetch(`${RecieverUrl}/light/on/${mac}`)
    .then((response) => response.json())
    .then((data) => {
      console.log(data)
    })
}
function runScene(id: string) {
  fetch(`${RecieverUrl}/scene/run/${id}`)
    .then((response) => response.json())
    .then((data) => {
      console.log(data)
    })
}

function deleteScene(id: string) {  
  console.log(`trying to delete scene: '${id}'`)
  if (id === "") {
    return
  }
  fetch(`${RecieverUrl}/scene/delete/${id}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      // 'Content-Type': 'application/x-www-form-urlencoded',
    },
  })
    .then((response) => response.json())
    .then((data) => {
      console.log(data)
      refreshScenes()
    })
}
```