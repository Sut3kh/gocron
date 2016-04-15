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
  name string
  Timer string
  Commands []string
}

/**
* Execute a CronTask.
**/
func run_cmds(crontask CronTask) {
  log.Printf("------Running %s------\n", crontask.name)
  for _, cmd := range crontask.Commands {
    out, err := exec.Command("sh", "-c", cmd).Output()
    log.Printf("%s: %s\n", cmd, out)
    if err != nil {
      log.Printf("error: command returned a non 0 exit status: %v %s", cmd, err)
    }
  }
}

func main() {
  log.Println("Starting gocron")

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
    crontask.name = f.Name()
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
    c.AddFunc(crontask.Timer, func() { run_cmds(crontask) })
  }

  // run crons
  c.Start()
  defer c.Stop()
  log.Printf("loaded %d cron task(s)\n\n", len(files))

  // run forever
  select {}
}
