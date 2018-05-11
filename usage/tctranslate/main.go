package main

import (
	"fmt"
	"io/ioutil"
	"strconv"

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

	//TODO add init rules
	tcRules = append(tcRules, node)
	tcRules = append(tcRules, "tc qdisc del dev eth0 root")
	tcRules = append(tcRules, "tc qdisc add dev eth0 root handle 1: htb default 1")
	tcRules = append(tcRules, "tc class add dev eth0 parent 1: classid 1:1 htb")
	tcRules = append(tcRules, "tc qdisc add dev eth0 parent 1:1 handle 1: sfq")

	tcIndex := 10
	//3. build tc tree for this node
	for _, group := range r.Group {
		rule := "tc class add dev eth0 parent 1: classid 1:" + strconv.Itoa(tcIndex) + " htb"
		tcRules = append(tcRules, rule)
		if group.Name == groupName { //local group
			//	3.1 parse other node in same group

			rule = "tc qdisc add dev eth0 parent 1:" + strconv.Itoa(tcIndex) + " handle " + strconv.Itoa(tcIndex) + ": tbf"
			rule = rule + " delay " + group.Delay
			tcRules = append(tcRules, rule)
			for _, otherNode := range group.Nodes {
				if otherNode == node {
					continue
				}

				rule = "tc filter add dev eth0 protocol ip parent 1:0 prio 1 u32 match ip dst " + otherNode + " flowid 1:" + strconv.Itoa(tcIndex)
				tcRules = append(tcRules, rule)
			}

		} else { //other group
			// find network first
			for _, network := range r.Network {
				if keyInArray(group.Name, network.Groups) && keyInArray(groupName, network.Groups) {
					rule = "tc qdisc add dev eth0 parent 1:" + strconv.Itoa(tcIndex) + " handle " + strconv.Itoa(tcIndex) + ": tbf"
					rule = rule + " delay " + network.Delay
					tcRules = append(tcRules, rule)
				}
			}

			for _, otherNode := range group.Nodes {
				rule = "tc filter add dev eth0 protocol ip parent 1:0 prio 1 u32 match ip dst " + otherNode + " flowid 1:" + strconv.Itoa(tcIndex)
				tcRules = append(tcRules, rule)
			}
		}

		tcIndex = tcIndex + 1
	}

	return tcRules
}

func printRules(rules []string) {
	fmt.Println()
	for _, rule := range rules {
		fmt.Println(rule)
	}
	fmt.Println()
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
			printRules(tcRules)
		}
	}
}
