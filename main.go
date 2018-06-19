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
	Corrupt   string `yaml:"corrupt"`
	Reorder   string `yaml:"reorder"`
}

type group struct {
	Name  string   `yaml:"name"`
	Nodes []string `yaml:"nodes"`

	Delay     string `yaml:"delay"`
	Loss      string `yaml:"loss"`
	Duplicate string `yaml:"duplicate"`
	Rate      string `yaml:"rate"`
	Corrupt   string `yaml:"corrupt"`
	Reorder   string `yaml:"reorder"`
}

type network struct {
	Groups []string `yaml:"groups"`

	Delay     string `yaml:"delay"`
	Loss      string `yaml:"loss"`
	Duplicate string `yaml:"duplicate"`
	Rate      string `yaml:"rate"`
	Corrupt   string `yaml:"corrupt"`
	Reorder   string `yaml:"reorder"`
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
	err := ioutil.WriteFile("scripts/"+ip+".sh", data, 0777)
	if err != nil {
		fmt.Println(err)
	}
}

func printDockerScript(r root) {
	launchFile, err := os.OpenFile(
		"launch.sh",
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		0777,
	)
	if err != nil {
		fmt.Println(err)
	}
	defer launchFile.Close()

	cleanFile, err := os.OpenFile(
		"clean.sh",
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		0777,
	)
	if err != nil {
		fmt.Println(err)
	}
	defer cleanFile.Close()

	var launchFileData, cleanFileData []string
	launchFileData = append(launchFileData, "#!/bin/bash\n")
	launchFileData = append(launchFileData, "docker network create --subnet=10.0.0.1/5 thunderdb_testnet")
	launchFileData = append(launchFileData, `DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"`)
	cleanFileData = append(cleanFileData, "#!/bin/bash\n")

	for _, group := range r.Group {
		for _, node := range group.Nodes {
			ip := strings.Split(node, "/")[0]
			launchFileData = append(launchFileData, "docker run -dit --rm --net thunderdb_testnet --ip "+ip+
				" -v $DIR/scripts:/scripts --cap-add=NET_ADMIN --name "+ip+" ns /scripts/"+ip+".sh")

			cleanFileData = append(cleanFileData, "docker stop "+ip)
		}
	}
	cleanFileData = append(cleanFileData, "docker network rm thunderdb_testnet")

	launchFileByte := []byte(strings.Join(launchFileData, "\n") + "\n")
	_, err = launchFile.Write(launchFileByte)
	if err != nil {
		fmt.Println(err)
	}

	cleanFileByte := []byte(strings.Join(cleanFileData, "\n") + "\n")
	_, err = cleanFile.Write(cleanFileByte)
	if err != nil {
		fmt.Println(err)
	}
}

func printGraphScript(r root) {

	gvFile, err := os.OpenFile(
		"graph.gv",
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		0666,
	)
	if err != nil {
		fmt.Println(err)
	}
	defer gvFile.Close()

	var gvFileData []string
	gvFileData = append(gvFileData, "digraph G {")
	gvFileData = append(gvFileData, "    compound=true;")

	for _, group := range r.Group {
		gvFileData = append(gvFileData, "    subgraph cluster_"+group.Name+" {")
		gvFileData = append(gvFileData, "        label = "+group.Name+";")
		gvFileData = append(gvFileData, "        style = rounded;")
		for i := 0; i < len(group.Nodes); i++ {
			for j := i + 1; j < len(group.Nodes); j++ {
				nodei := group.Nodes[i]
				nodej := group.Nodes[j]
				//"10.2.1.1/16" -> "10.3.1.1/20" [arrowhead=none, arrowtail=none, label="delay\n 100ms ±10ms 30%"];
				gvFileData = append(gvFileData, "        \""+nodei+"\" -> \""+nodej+
					"\" [arrowhead=none, arrowtail=none, label=\"delay "+
					group.Delay+"\"];")
			}
		}
		gvFileData = append(gvFileData, "    }")
	}

	for _, network := range r.Network {
		//parse two group pair
		for i := 0; i < len(network.Groups); i++ {
			for j := i + 1; j < len(network.Groups); j++ {
				var groupNodei, groupNodej string
				// get group ip
				for _, group := range r.Group {
					if group.Name == network.Groups[i] {
						groupNodei = group.Nodes[0]
					} else if group.Name == network.Groups[j] {
						groupNodej = group.Nodes[0]
					}
				}
				gvFileData = append(gvFileData, "    \""+groupNodei+"\" -> \""+groupNodej+
					"\"\n        [ltail=cluster_"+network.Groups[i]+", lhead=cluster_"+network.Groups[j]+
					", arrowhead=none, arrowtail=none,\n        label=\"delay "+
					network.Delay+"\"];")
			}
		}
	}

	gvFileData = append(gvFileData, "}")
	gvFileByte := []byte(strings.Join(gvFileData, "\n") + "\n")
	_, err = gvFile.Write(gvFileByte)
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
//TODO add all tc param support
func main() {
	r := root{}
	//TODO 1. read yaml from specific file
	data, err := ioutil.ReadFile("example.yaml")
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
		}
	}
	printDockerScript(r)

	printGraphScript(r)
	//TODO 2. check dot commandline
	// dot -Tpng graph.gv -o graph.png

}
