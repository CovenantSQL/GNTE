package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

type tcParams struct {
	Delay     string `yaml:"delay"`
	Loss      string `yaml:"loss"`
	Duplicate string `yaml:"duplicate"`
	Rate      string `yaml:"rate"`
}

type group struct {
	Name      string   `yaml:"name"`
	Delay     string   `yaml:"delay"`
	Loss      string   `yaml:"loss"`
	Duplicate string   `yaml:"duplicate"`
	Rate      string   `yaml:"rate"`
	Nodes     []string `yaml:"nodes"`
}

type network struct {
	Groups    []string `yaml:"groups"`
	Delay     string   `yaml:"delay"`
	Loss      string   `yaml:"loss"`
	Duplicate string   `yaml:"duplicate"`
	Rate      string   `yaml:"rate"`
}

type root struct {
	Group   []group   `yaml:"group"`
	Network []network `yaml:"network"`
}

func keyInArray(key string, array []string) bool {
	for _, data := range array {
		if data == key {
			return true
		}
	}
	return false
}

func processOneNode(node string, groupName string, r root) []string {

	var tcRules []string

	//add init rules
	tcRules = append(tcRules, "#!/bin/sh\n")
	tcRules = append(tcRules, "#"+node)
	tcRules = append(tcRules, "tc qdisc del dev eth0 root")
	tcRules = append(tcRules, "tc qdisc add dev eth0 root handle 1: htb default 2")
	tcRules = append(tcRules, "tc class add dev eth0 parent 1: classid 1:2 htb rate 10gbps")
	tcRules = append(tcRules, "tc qdisc add dev eth0 parent 1:2 handle 2: sfq")

	tcIndex := 10
	//3. build tc tree for this node
	for _, group := range r.Group {
		rule := "tc class add dev eth0 parent 1: classid 1:" + strconv.Itoa(tcIndex) + " htb rate 10gbps"
		tcRules = append(tcRules, rule)
		if group.Name == groupName { //local group
			//	3.1 parse other node in same group

			rule = "tc qdisc add dev eth0 parent 1:" + strconv.Itoa(tcIndex) + " handle " + strconv.Itoa(tcIndex) + ": netem"
			rule = rule + " delay " + group.Delay
			tcRules = append(tcRules, rule)
			for _, otherNode := range group.Nodes {
				if otherNode == node {
					continue
				}

				//3.2 build tc leaf for inner-group
				rule = "tc filter add dev eth0 protocol ip parent 1:0 prio 1 u32 match ip dst " + otherNode + " flowid 1:" + strconv.Itoa(tcIndex)
				tcRules = append(tcRules, rule)
			}

		} else { //other group
			//3.3 parse node in other group
			// find network first
			for _, network := range r.Network {
				if keyInArray(group.Name, network.Groups) && keyInArray(groupName, network.Groups) {
					rule = "tc qdisc add dev eth0 parent 1:" + strconv.Itoa(tcIndex) + " handle " + strconv.Itoa(tcIndex) + ": netem"
					rule = rule + " delay " + network.Delay
					tcRules = append(tcRules, rule)
				}
			}

			//3.4 build tc leaf for group-connection
			for _, otherNode := range group.Nodes {
				rule = "tc filter add dev eth0 protocol ip parent 1:0 prio 1 u32 match ip dst " + otherNode + " flowid 1:" + strconv.Itoa(tcIndex)
				tcRules = append(tcRules, rule)
			}
		}

		tcIndex = tcIndex + 1
	}

	tcRules = append(tcRules, "tail -f /dev/null")
	return tcRules
}

func printRules(rules []string) {
	fmt.Println()
	for _, rule := range rules {
		fmt.Println(rule)
	}
	fmt.Println()
}

func printTcScript(rules []string, node string) {
	ip := strings.Split(node, "/")[0]
	fmt.Println(ip)

	var data []byte
	rulestr := strings.Join(rules, "\n")
	data = []byte(rulestr + "\n")

	//	fmt.Println(rulestr)
	err := ioutil.WriteFile("/home/laodouya/thunderdb/ns/scripts/"+ip+".sh", data, 0777)
	if err != nil {
		fmt.Println(err)
	}
}

func printDockerScript(r root) {

	initFile, err := os.OpenFile(
		"/home/laodouya/thunderdb/ns/init.sh",
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		0777,
	)
	if err != nil {
		fmt.Println(err)
	}
	defer initFile.Close()

	lanuchFile, err := os.OpenFile(
		"/home/laodouya/thunderdb/ns/launch.sh",
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		0777,
	)
	if err != nil {
		fmt.Println(err)
	}
	defer lanuchFile.Close()

	var initFileData, lanuchFileData []string
	initFileData = append(initFileData, "#!/bin/bash\n")
	lanuchFileData = append(lanuchFileData, "#!/bin/bash\n")
	lanuchFileData = append(lanuchFileData, `DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"`)

	for _, group := range r.Group {
		for _, node := range group.Nodes {
			ip := strings.Split(node, "/")[0]
			initFileData = append(initFileData, "docker network create --subnet="+node+" "+ip+"\n")

			lanuchFileData = append(lanuchFileData, "sudo docker run -dit --rm --net "+ip+" --ip  "+ip+" -v $DIR/scripts:/scripts --cap-add=NET_ADMIN ns /scripts/" + ip + ".sh")
		}
	}

	initFileByte := []byte(strings.Join(initFileData, "\n") + "\n")
	lanuchFileByte := []byte(strings.Join(lanuchFileData, "\n") + "\n")
	_, err = initFile.Write(initFileByte)
	if err != nil {
		fmt.Println(err)
	}
	_, err = lanuchFile.Write(lanuchFileByte)
	if err != nil {
		fmt.Println(err)
	}

}

/*
1. read yaml from file
2. select one node
3. build tc tree for this node
	3.1 parse other node in same group
	3.2 build tc leaf for inner-group
	3.3 parse node in other group
	3.4 build tc leaf for group-connection
4. print tc tree
*/
//TODO 1. read yaml from specific file
// 2. write to running dir
func main() {
	r := root{}
	//1. read yaml from file
	data, err := ioutil.ReadFile("/home/laodouya/thunderdb/ns/usage/example.yaml")
	if err != nil {
		fmt.Println(err)
	}

	err = yaml.Unmarshal(data, &r)
	if err != nil {
		fmt.Println(err)
	}

	//2. select one node
	for _, group := range r.Group {
		for _, node := range group.Nodes {
			tcRules := processOneNode(node, group.Name, r)
			//4. print tc tree
			printTcScript(tcRules, node)
			//	printRules(tcRules)
		}
	}
	printDockerScript(r)
}
