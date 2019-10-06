
package main

import (
    "fmt"
    "sync"
    "time"
    "os"
	"encoding/json"
    
)
var Tick,Tock,Bong int=0,0,0 

var quit = make(chan bool)
type Progress struct {
    current string
    rwlock  sync.RWMutex
}

func (p *Progress) Set(value string) {
    p.rwlock.Lock()
    defer p.rwlock.Unlock()
    p.current = value
}
func (p *Progress) Get() string {
    p.rwlock.RLock()
    defer p.rwlock.RUnlock()
    return p.current
}

func CronJob(progress *Progress) { 
    for { 
        time.Sleep(1* time.Second)
        unitSound := LoadConfiguration()
        Tick++
        if Tick >= 59 { 
            Tick =0
            if Tock >=59{
                Tock =0
                Bong++
                if Bong >=3{
                    progress.Set(fmt.Sprintf("%v %v",unitSound.BongSound,Bong))
                    fmt.Println("*******End********")
                    quit <- true
                }else{
                      fmt.Print("Hr is Over---->:",Bong) 
                     } 
                } else{
                    Tock++
                    progress.Set(fmt.Sprintf("%v --> %v",unitSound.TockSound,Tock))

            }
        }else { 
            progress.Set(fmt.Sprintf("%v %v",unitSound.TickSound,Tick))   
        }
    }
}
func LoadConfiguration() ticktokUnit {
    var config ticktokUnit
    configFile, err := os.Open("ticktockConfig.json")
    defer configFile.Close()
    if err != nil {
        fmt.Println(err.Error())
    }
    jsonParser := json.NewDecoder(configFile)
    jsonParser.Decode(&config)
    return config
}

type ticktokUnit struct {
	TickSound string `json:"TickSound"`
	TockSound string `json:"TockSound"`
	BongSound string  `json:"BongSound"`
}
func main() {
    
    fmt.Println("*******Starting********")
    c := time.Tick(1 * time.Second)
    progress := &Progress{}
    
    go CronJob(progress)
    for {
        select {
        case <-c:
            fmt.Println(progress.Get())
        case <- quit:
            return
        }

    }
}