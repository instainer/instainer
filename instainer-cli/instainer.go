package main

import (
  "os"
  "bytes"
  "fmt"
  "net/http"
  "flag"
  "io/ioutil"
  "github.com/codegangsta/cli"
  "github.com/fsouza/go-dockerclient"
  "github.com/fgrehm/go-dockerpty"
  "encoding/json"
  "github.com/fatih/color"
  "github.com/antonholmquist/jason"
  "github.com/CrowdSurge/banner"
  "github.com/rakyll/globalconf"
)

const (
  backend     = "http://beta.instainer.com/backend/api"
  execBackend = "http://beta.instainer.com:80/docker"
)

var (
    flagApiKey    = flag.String("apikey", "", "Instainer API Key")
)
type flagValue struct {
	str string
}

func (f *flagValue) String() string {
	return f.str
}

func (f *flagValue) Set(value string) error {
	f.str = value
	return nil
}

func newFlagValue(val string) *flagValue {
	return &flagValue{str: val}
}

func main() {
  conf, err := globalconf.New("instainer")
  conf.ParseAll()
  if err != nil {
      fmt.Printf("%s", err)
      os.Exit(1)
  }

  app := cli.NewApp()
  app.Name = "instainer-cli"
  app.Usage = "Instant Docker containers on the cloud"
  app.Version = "1.0.0"

  app.Commands = []cli.Command{
    {
        Name: "config",
        Usage:  "Configure instainer client",
        Action: func(c *cli.Context) {
          if len(c.Args()) == 0 {
            fmt.Println("Please enter a valid API key")
            os.Exit(1)
          }
          configureInstainerClient(c.Args()[0])
        },
    },
    {
        Name: "run",
        Usage:  "Run a container",

        Flags : []cli.Flag {
          cli.StringSliceFlag{
            Name: "volume,v",
            Usage: "volume to mount",
          },
          cli.StringSliceFlag{
            Name: "env,e",
            Usage: "environment variable to add",
          },
        },
        Action: func(c *cli.Context) {
          createContainer(c.Args()[0], c.StringSlice("volume"), c.StringSlice("env"))
        },
    },
    {
        Name: "bash",
        Usage:  "Bash into container",
        Action: func(c *cli.Context) {
          bashContainer(c.Args()[0])
        },
    },
    {
        Name: "ps",
        Usage:  "List containers",
        Action: func(c *cli.Context) {
          listContainers()
        },
    },
    {
        Name: "exec",
        Usage:  "Exec into container",
        Action: func(c *cli.Context) {
          execCommand(c.Args()[0],c.Args()[1])
        },
    },
    {
        Name: "logs",
        Usage:  "Get logs from container",
        Action: func(c *cli.Context) {
          getLogs(c.Args()[0])
        },
    },
    {
        Name: "compose",
        Usage:  "Up compose file",
        Subcommands: []cli.Command{
          {
            Name:  "up",
            Usage: "up new docker-compose file",
            Action: func(c *cli.Context) {
                upCompose(c.Args().First())
            },
          },
        },
    },
}
  app.Run(os.Args)
}

func check(e error) {
  if e != nil {
      fmt.Printf("%s", e)
      os.Exit(1)
  }
}
func configureInstainerClient(apikey string) {
  conf, err := globalconf.New("instainer")
  conf.ParseAll()
  if err != nil {
      fmt.Printf("%s", err)
      os.Exit(1)
  }

  f := &flag.Flag{Name: "apikey", Value: newFlagValue(apikey)}
  conf.Set("", f)
  conf.ParseAll()

  instainerGet("/instainer-cli")

}

func instainerPost(url string,data string) (*http.Response) {
  client := &http.Client{}
  req, err := http.NewRequest("POST", backend+url,bytes.NewBufferString(data))
  req.Header.Add("API-KEY", *flagApiKey)
  req.Header.Add("Content-Type", "application/json")

  response, err := client.Do(req)
  if err != nil {
      fmt.Printf("Error occured")
      os.Exit(1)
  }

  if (response.StatusCode!=200){
    contents, err := ioutil.ReadAll(response.Body)
    check(err)
    fmt.Println("Error occured ", string(contents))
    os.Exit(1)
  }

  return response
}

func instainerGet(url string) (*http.Response) {
  client := &http.Client{}
  req, err := http.NewRequest("GET", backend+url,nil)
  req.Header.Add("API-KEY", *flagApiKey)
  req.Header.Add("Content-Type", "application/json")

  response, err := client.Do(req)
  if err != nil {
      fmt.Printf("Error occured")
      os.Exit(1)
  }

  if (response.StatusCode!=200){
    contents, err := ioutil.ReadAll(response.Body)
    check(err)
    fmt.Println("Error occured ", string(contents))
    os.Exit(1)
  }
  return response
}

type RunParams struct {
    VolumeRequests []string
    EnvVariables []string
}

func createContainer(dockerName string, volumeRequest []string, envVariables []string) {

  m := &RunParams{VolumeRequests:volumeRequest,EnvVariables:envVariables}
  b, err := json.Marshal(m)

  response := instainerPost("/container/run?image="+dockerName,string(b))

  defer response.Body.Close()

  contents, err := ioutil.ReadAll(response.Body)
  check(err)
  v, err := jason.NewObjectFromBytes(contents)
  data, err := v.GetObject("data")

  username, err :=  data.GetString("gitUser")
  password, err :=  data.GetString("gitPassword")

  banner.Print("instainer")
  fmt.Println("")
  fmt.Println("")

  color.Green("------------Git Information------------")
  color.Yellow("Git User      = %s", username)
  color.Yellow("Git Password  = %s", password)
  fmt.Println("")
  fmt.Println("")

  color.Green("----------Volumes Information----------")
  volumes, err := data.GetObjectArray("volumes")
  for _, volume := range volumes {
    mntDir, err := volume.GetString("mntDir")
    gitUrl, err := volume.GetString("gitUrl")
    check(err)
    color.Blue("    %s", mntDir)
    color.Yellow("    Git URL  = %s", gitUrl)
    fmt.Println("")
  }
  fmt.Println("")

  color.Green("------------Port Information------------")
  ports, err := data.GetObjectArray("ports")

  for _, port := range ports {
    dockerPort, err := port.GetString("dockerPort")
    instainerPort, err := port.GetString("instainerPort")
    check(err)
    color.Blue("  Container Port    %s", dockerPort)
    color.Yellow("  Instainer Port    instainer.io:%s", instainerPort)
    fmt.Println("")
  }
  fmt.Println("")

  envVariablesData, err := data.GetStringArray("envVariables")

  if (len(envVariablesData)>0){
    color.Green("---------Environment Variables----------")

    for _, variable := range envVariablesData {
      color.Yellow("    Variable=Value      = %s", variable)
    }
  }
  fmt.Println("")
  fmt.Println("")
  fmt.Println("Successfully deployed!")
  fmt.Println("")
  fmt.Println("")

  check(err)
}

func upCompose(dockerComposeFile string) {
  dat, err := ioutil.ReadFile(dockerComposeFile)
  check(err)
  response := instainerPost("/compose/up",string(dat))

  defer response.Body.Close()
  contents, err := ioutil.ReadAll(response.Body)
  check(err)

  fmt.Println(string(contents))

  v, err := jason.NewObjectFromBytes(contents)
  data, err := v.GetObjectArray("data")

  banner.Print("instainer")
  fmt.Println("")
  fmt.Println("")


  for _, container := range data {
    username, err :=  container.GetString("gitUser")
    password, err :=  container.GetString("gitPassword")

    name, err :=  container.GetString("name")
    containerId, err :=  container.GetString("id")

    fmt.Println("Container Name ", name)
    fmt.Println("Container Id  ", containerId)

    fmt.Println("")
    color.Green("------------Git Information------------")
    color.Yellow("Git User      = %s", username)
    color.Yellow("Git Password  = %s", password)
    fmt.Println("")
    fmt.Println("")

    color.Green("----------Volumes Information----------")
    volumes, err := container.GetObjectArray("volumes")
    for _, volume := range volumes {
      mntDir, err := volume.GetString("mntDir")
      gitUrl, err := volume.GetString("gitUrl")
      check(err)
      color.Blue("    %s", mntDir)
      color.Yellow("    Git URL  = %s", gitUrl)
      fmt.Println("")
    }
    fmt.Println("")

    color.Green("------------Port Information------------")
    ports, err := container.GetObjectArray("ports")

    for _, port := range ports {
      dockerPort, err := port.GetString("dockerPort")
      instainerPort, err := port.GetString("instainerPort")
      check(err)
      color.Blue("  Container Port    %s", dockerPort)
      color.Yellow("  Instainer Port    instainer.io:%s", instainerPort)
      fmt.Println("")
    }
    fmt.Println("")

    envVariablesData, err := container.GetStringArray("envVariables")

    if (len(envVariablesData)>0){
      color.Green("---------Environment Variables----------")

      for _, variable := range envVariablesData {
        color.Yellow("    Variable=Value      = %s", variable)
      }
    }
    fmt.Println("")
    fmt.Println("")
    check(err)
  }

  fmt.Println("Successfully deployed!")
  fmt.Println("")
  fmt.Println("")


}

func getLogs(dockerId string) {

  response := instainerGet("/container/logs/"+dockerId)

  defer response.Body.Close()
  contents, err := ioutil.ReadAll(response.Body)
  if err != nil {
      fmt.Printf("%s", err)
      os.Exit(1)
  }
  fmt.Printf("%s\n", string(contents))
}

func listContainers() {

  response := instainerGet("/containers")

  defer response.Body.Close()
  contents, err := ioutil.ReadAll(response.Body)
  check(err)
  v, err := jason.NewObjectFromBytes(contents)
  data, err := v.GetObject("data")

  permanents, err := data.GetObjectArray("permanents")
  nonpermanents, err := data.GetObjectArray("nonPermanents")

  banner.Print("instainer")
  fmt.Println("")
  fmt.Println("")

  fmt.Println("Permanent Containers")
  fmt.Println("")

  fmt.Printf("%-32s %-52s %-32s %-24s\n", "CONTAINER ID", "NAME","IMAGE NAME","CREATED")

  for _, container := range permanents {

    containerId, err :=  container.GetString("id")
    name, err :=  container.GetString("name")
    imageName, err :=  container.GetString("imageName")
    createdTime, err :=  container.GetString("createdTime")

    check(err)

    fmt.Printf("%-32s %-52s %-32s %-24s\n", containerId, name, imageName,createdTime)

  }
  fmt.Println("")
  fmt.Println("")

  fmt.Println("Non-Permanent Containers")
  fmt.Println("")

  fmt.Printf("%-32s %-52s %-32s %-24s\n", "CONTAINER ID", "NAME","IMAGE NAME","CREATED")

  for _, container := range nonpermanents {

    containerId, err :=  container.GetString("id")
    name, err :=  container.GetString("name")
    imageName, err :=  container.GetString("imageName")
    createdTime, err :=  container.GetString("createdTime")

    check(err)

    fmt.Printf("%-32s %-52s %-32s %-24s\n", containerId, name, imageName,createdTime)

  }
}

type ExecResponse struct {
    Data string `json:"data"`
    Success bool `json:"success"`
}

type ExecRequest struct {
    Commands []string `json:"commands"`
}

func bashContainer(containerId string) {
  m := &ExecRequest{Commands:[]string{"bash"}}
  b, err := json.Marshal(m)

  response := instainerPost("/container/"+containerId+"/exec",string(b))
  defer response.Body.Close()
  contents, err := ioutil.ReadAll(response.Body)
  if err != nil {
      fmt.Printf("%s", err)
      os.Exit(1)
  }
  fmt.Println(string(contents))

  var execResponse ExecResponse
  err = json.Unmarshal(contents, &execResponse)
  fmt.Println(execResponse.Data)

  client, _ := docker.NewClient(execBackend)

  if err != nil {
      fmt.Println(err)
      os.Exit(1)
  }

  // Fire up the console
  if err = dockerpty.StartExecWithId(client, execResponse.Data); err != nil {
      fmt.Println(err)
      os.Exit(1)
  }
}

func execCommand(containerId string,command string) {

  m := &ExecRequest{Commands:[]string{command}}
  b, err := json.Marshal(m)

  response := instainerPost("/container/"+containerId+"/exec",string(b))
  defer response.Body.Close()
  contents, err := ioutil.ReadAll(response.Body)
  if err != nil {
      fmt.Printf("%s", err)
      os.Exit(1)
  }
  fmt.Println(string(contents))

  var execResponse ExecResponse
  err = json.Unmarshal(contents, &execResponse)
  fmt.Println(execResponse.Data)

  client, _ := docker.NewClient(execBackend)

  if err != nil {
      fmt.Println(err)
      os.Exit(1)
  }

  if err = dockerpty.StartExecWithId(client, execResponse.Data); err != nil {
      fmt.Println(err)
      os.Exit(1)
  }
}

