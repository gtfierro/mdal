package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/jinzhu/configor"
	"github.com/op/go-logging"
)

// logger
var log *logging.Logger

func init() {
	log = logging.MustGetLogger("mdal")
	var format = "%{color}%{level} %{shortfile} %{time:Jan 02 15:04:05} %{color:reset} ▶ %{message}"
	var logBackend = logging.NewLogBackend(os.Stderr, "", 0)
	logBackendLeveled := logging.AddModuleLevel(logBackend)
	logging.SetBackend(logBackendLeveled)
	logging.SetFormatter(logging.MustStringFormatter(format))
}

func main() {

	configor.Load(&Config, "config.yml")

	c := newCore()

	t0, _ := time.Parse("2006-01-02 15:04:05", "2017-10-23 00:00:30")
	t1, _ := time.Parse("2006-01-02 15:04:05", "2017-10-13 00:00:00")
	q := Query{
		Composition: []string{"temp", "4d6e251a-48e1-3bc0-907d-7d5440c34bb9"},
		Selectors:   []Selector{MEAN, MEAN},
		Variables: []VarParams{
			VarParams{
				Name:       "temp",
				Definition: "SELECT ?temp_uuid WHERE { ?temp rdf:type/rdfs:subClassOf* brick:Temperature_Sensor . ?temp bf:uuid ?temp_uuid . };",
				Units:      "na",
			},
		},
		Time: TimeParams{
			T0:         t0,
			T1:         t1,
			WindowSize: 36000000000,
			Aligned:    true,
		},
		Params: Params{
			Statistical: false,
			Window:      true,
		},
	}
	go func() {
		log.Fatal(http.ListenAndServe("localhost:6060", nil))
	}()

	ts, err := c.HandleQuery(q)
	if err != nil {
		log.Error(err)
	} else {
		log.Debugf("%+v", ts.Info())
	}
	b1, err := ts.msg.Marshal()
	if err != nil {
		log.Error(err)
	} else {
		log.Info(len(b1))
	}
	b2, err := ts.msg.MarshalPacked()
	if err != nil {
		log.Error(err)
	} else {
		log.Info(len(b2))
	}
	fmt.Println(q)

	log.Fatal(RunBosswave(c))
}
