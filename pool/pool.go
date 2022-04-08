package pool

import (
	"DuGoMiner/config"
	"DuGoMiner/errors"
	"DuGoMiner/job"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

type Pool struct {
	Ip         string `json:"ip"`
	Name       string `json:"name"`
	Port       uint16 `json:"port"`
	Server     string `json:"server"`
	Success    bool   `json:"success"`
	socket     string
	config     *config.Config
	Connection net.Conn
}

func New(config *config.Config) *Pool {
	return &Pool{
		config: config,
	}
}

func (p *Pool) GetPool() {
	//Get a mining pool server
	resp, err := http.Get(p.config.GetPoolUrl)
	_ = errors.CheckErr(err)
	body, err := ioutil.ReadAll(resp.Body)
	_ = errors.CheckErr(err)
	if err := json.Unmarshal(body, &p); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}
	p.socket = fmt.Sprintf("%s:%d", p.Ip, p.Port)
	_ = errors.CheckErr(err)

}

// GetJob TODO should this return a pointer?
func (p *Pool) GetJob() job.Job {

	var err error
	p.Connection, err = net.Dial("tcp", p.socket)
	_ = errors.CheckErr(err)
	reply := make([]byte, 1024)
	_, err = p.Connection.Read(reply)
	_ = errors.CheckErr(err)
	log.Printf("Version : %s", string(reply))
	_, err = p.Connection.Write([]byte(fmt.Sprintf("JOB,%s,%s", p.config.UserName, p.config.Difficulty)))
	_ = errors.CheckErr(err)
	_, err = p.Connection.Read(reply)
	_ = errors.CheckErr(err)
	replySplit := bytes.Split(reply, []byte(","))
	for i := 0; i < len(replySplit); i++ {
		replySplit[i] = bytes.Replace(replySplit[i], []byte("\x00"), []byte(""), -1)
		replySplit[i] = bytes.Replace(replySplit[i], []byte("\n"), []byte(""), -1)
	}

	return job.Job{
		Base:       string(replySplit[0]),
		Target:     string(replySplit[1]),
		Difficulty: string(replySplit[2]),
	}

}

func (p *Pool) CloseConnection() {
	p.Connection.Close()
}
