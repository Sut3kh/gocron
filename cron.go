package main

import (
  // "fmt"
  "log"
  "io/ioutil"
  "os/exec"
  "github.com/robfig/cron"
  "gopkg.in/yaml.v2"
)

type CronTask struct {
  Timer string
  Commands []string
}

/**
* Run an array of shell commands.
**/
func run_cmds(cmds []string) {
  for _, cmd := range cmds {
    out, err := exec.Command("sh", "-c", cmd).Output()
    if err != nil {
      log.Fatalf("error: failed to run %s: %v", cmd, err)
    }
    log.Printf("%s: %s\n", cmd, out)
  }
}

func main() {
  // look for cronfiles in /etc/gocron.d
  crondir := "/etc/gocron.d"
  files, _ := ioutil.ReadDir(crondir)
  if files == nil {
    log.Fatalf("error: no cron files found in %s", crondir)
  }

  // create new cron schedule
  c := cron.New()

  // loop cronfiles
  for _, f := range files {
    // try and read file
    crontask := CronTask{}
    file, err := ioutil.ReadFile(crondir + "/" + f.Name())
    if err != nil {
      log.Fatalf("error: cannot read %s: %v", f.Name(), err)
    }

    // try to parse yaml
    err = yaml.Unmarshal(file, &crontask)
    if err != nil {
      log.Fatalf("error: invalid cronfile %s: %v", f.Name(), err)
    }

    // add tasks
    c.AddFunc(crontask.Timer, func() {
      log.Printf("\n------Running %s------\n", f.Name())
      run_cmds(crontask.Commands)
    })
  }

  // run crons
  c.Start()
  defer c.Stop()
  log.Println("Started gocron")

  // run forever
  select {}
}
