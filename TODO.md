# Lights 
 
 * [x] - Load a list of lights into memory 
 * [x] - api for turning lights off and on
 * [x] - refresh the lights details every so often
 * [] - cache lights on disk for faster startup
 * [] - thread all on/all off light operations

# Scenes

* [x] - scene object and scene response object 
* [x] - hard code scenes
* [x] - list scenes. 
* [x] - ~~print api to get the code to create a new scene. ~~ Just take in a json blob of a scene
* [x] - api to create new scenes
* [x] - api to delete scenes
* [x] - api to delete scenes
* [x] - save and load scenes from file on change and boot
* [] - colour/style property on scenes for customizing them further on client side
* [] - error check everything coming in
* [] - thread scene light operations
* [] - expose `lightActions`
  * [] - api actions to set a light with exact lightAction values

# misc 

[] - add log levels to be able to hide update messages


# Far future 

* [] - Talk to lifx's api
  * [] - allow activating remote scenes
  * [] - parse remote scenes so local lan actions can still be run
