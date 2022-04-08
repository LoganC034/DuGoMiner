package miner

import (
	"DuGoMiner/config"
	"DuGoMiner/errors"
	job "DuGoMiner/job"
	"DuGoMiner/pool"
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"time"
)

// Miner TODO Pointers?
type Miner struct {
	Pool pool.Pool
	Job  job.Job
}

type MineResult struct {
	Nonce    string
	HashRate string
}

func New() *Miner {
	log.Println("Building Miner")
	c := config.New()
	c.GetConfig()
	log.Printf("Using The Following Config :\n \t\t\t\t USERNAME : %s\n \t\t\t\t DIFFICULTY : %s\n \t\t\t\t URL : %s\n", c.UserName, c.Difficulty, c.GetPoolUrl)
	p := pool.New(c)
	p.GetPool()
	log.Printf("Pool Data Received :\n \t\t\t\t SERVER: %s\n \t\t\t\t POOL : %s\n \t\t\t\t IP : %s\n \t\t\t\t PORT : %d", p.Server, p.Name, p.Ip, p.Port)
	return &Miner{
		Pool: *p,
	}
}

func (m *Miner) Mine() {
	defer m.Pool.CloseConnection()
	var err error
	for {
		log.Println("Requesting Job")
		m.Job = m.Pool.GetJob()
		log.Printf("Job Received :\n \t\t\t\t BASE: %s\n \t\t\t\t TARGET : %s\n \t\t\t\t DIFFICULTY : %s\n", m.Job.Base, m.Job.Target, m.Job.Difficulty)

		//MINE
		start := time.Now().Unix()
		//TODO Possible refactor of what type values are stored as...
		diff, _ := strconv.Atoi(m.Job.Difficulty)
		//h := HexToBytes(m.Job.Base)
		var mineResult MineResult

		for i := 0; i < 100*diff; i++ {
			strNonce := strconv.Itoa(i)
			h := m.Job.Base + strNonce

			result := sha1.Sum([]byte(h))
			resultHex := BytesToHex(result[:])
			if resultHex == m.Job.Target {
				fmt.Println("Match")
				stop := time.Now().Unix()
				seconds := stop - start
				fmt.Println(seconds)
				mineResult.HashRate = fmt.Sprintf("%d", int64(i)/seconds)
				mineResult.Nonce = strNonce
				//return mineResult
				_, err = m.Pool.Connection.Write([]byte(mineResult.Nonce + "," + mineResult.HashRate + "," + "DuGoMiner" + "," + "DevMiner"))
				errors.CheckErr(err)
				reply := make([]byte, 1024)
				m.Pool.Connection.Read(reply)
				fmt.Println(string(reply))

			}
		}
	}
}

func BytesToHex(b []byte) string {
	hexString := hex.EncodeToString(b)
	return hexString
}

func HexToBytes(s string) []byte {
	byteArray, _ := hex.DecodeString(s)
	return byteArray
}

func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)

	}

	return buff.Bytes()
}
