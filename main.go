// ansible_influxdb by Novaquark
//
// To the extent possible under law, the person who associated CC0 with
// sysinfo_influxdb has waived all copyright and related or neighboring rights
// to sysinfo_influxdb.
//
// You should have received a copy of the CC0 legalcode along with this
// work.  If not, see <http://creativecommons.org/publicdomain/zero/1.0/>.

package main

import (
	"errors"
	"fmt"
	"os"
	"io/ioutil"
	"strings"
	"encoding/json"
	influxClient "github.com/influxdb/influxdb-go"
)

type ReturnCh struct {
	Changed bool			`json:"changed"`
}

type ReturnEr struct {
	Failed	bool			`json:"failed"`
	Msg     string			`json:"msg"`
}

type ReturnFa struct {
	Changed bool			`json:"changed"`
	Facts   map[string]interface{}	`json:"ansible_facts"`
}

func parseConf(path string) map[string]string {
	conf := make(map[string]string)

	data, err:= ioutil.ReadFile(path)
	if err != nil {
		fail(err)
	}
	for _, arg := range(strings.Split(string(data), " ")) {
		if strings.Contains(arg, "="){
			line := strings.Split(arg, "=")
			conf[line[0]] = line[1];
		}
	}

	return conf
}

func fail(err error) {
	ret := ReturnEr{ true, err.Error() }
	s, _ := json.Marshal(ret)
	fmt.Printf("%s\n", s)
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		fail(errors.New("No argument file receive"))
	}

	args := parseConf(os.Args[1])

	config := influxClient.ClientConfig{}

	if v, present := args["login_host"]; present {
		config.Host = v
	}
	if v, present := args["login_username"]; present {
		config.Username = v
	}
	if v, present := args["login_password"]; present {
		config.Password = v
	}

	client, err := influxClient.NewClient(&config)
	if err != nil {
		fail(err)
	}

	if uType, present := args["user_type"]; present {
		// Manage users

		var username string
		if username, present = args["username"]; ! present {
			fail(errors.New("You have not given username parameter"))
		}
		password, _ := args["password"];

		user_exists := false
		ret := ReturnCh{ false }

		if uType == "cluster_admin" {
			u, _ := client.GetClusterAdminList()
			for _, v := range(u) {
				if v["username"] == username {
					user_exists = true
					break
				}
			}

			if state, present := args["state"]; present && state == "absent" && user_exists {
				err := client.DeleteClusterAdmin(username)

				if err != nil {
					fail(err)
				} else {
					ret.Changed = true
				}
			} else if state, present := args["state"]; present && state == "present" {
				if ! user_exists {
					err := client.CreateClusterAdmin(username, password)

					if err != nil {
						fail(err)
					} else {
						ret.Changed = true
					}
				} else if password != "" {
					client.UpdateClusterAdmin(username, password)
				}
			}
			
		} else if dbname, present := args["database"]; present && uType == "user" {
			u, _ := client.GetDatabaseUserList(dbname)
			for _, v := range(u) {
				if v["name"] == username {
					user_exists = true
					break
				}
			}

			if state, present := args["state"]; present && state == "absent" && user_exists {
				err := client.DeleteDatabaseUser(dbname, username)

				if err != nil {
					fail(err)
				} else {
					ret.Changed = true
				}
			} else if state, present := args["state"]; present && state == "present" {
				if ! user_exists {
					err := client.CreateDatabaseUser(dbname, username, password)

					if err != nil {
						fail(err)
					} else {
						ret.Changed = true
					}
				} else if password != "" {
					client.UpdateDatabaseUser(dbname, username, password)
				}
			}
		}

		s, _ := json.Marshal(ret)
		fmt.Printf("%s\n", s)
	} else if dbname, present := args["database"]; present {
		// Manage databases

		db_exists := false
		ret := ReturnCh{ false }

		dbs, _ := client.GetDatabaseList()
		for _, v := range(dbs) {
			if v["name"] == dbname {
				db_exists = true
				break
			}
		}

		if state, present := args["state"]; present && state == "absent" && db_exists {
			err := client.DeleteDatabase(dbname)

			if err != nil {
				fail(err)
			} else {
				ret.Changed = true
			}
		} else if state, present := args["state"]; present && state == "present" && ! db_exists {
			err := client.CreateDatabase(dbname)

			if err != nil {
				fail(err)
			} else {
				ret.Changed = true
			}
		}

		s, _ := json.Marshal(ret)
		fmt.Printf("%s\n", s)
	} else {
		// Gathering facts

		gprefix := "influxdb"
		if v, present := args["facts_prefix"]; present {
			gprefix = v
		}

		ret := ReturnFa{ false, make(map[string]interface{}) }

		dbs, _ := client.GetDatabaseList()
		ret.Facts[gprefix + "_dbs"] = dbs

		adms, _ := client.GetClusterAdminList()
		ret.Facts[gprefix + "_cluster_admins"] = adms

		for _, db := range(dbs) {
			dbname := db["name"].(string);
			users, _ := client.GetDatabaseUserList(dbname)
			ret.Facts[gprefix + "_users_" + dbname] = users
		}

		s, _ := json.Marshal(ret)
		fmt.Printf("%s\n", s)
	}

//	s, err := json.Marshal(m) // produce JSON from Time struct
//	if err != nil {
//		fmt.Println(err)
//	}else{
//		fmt.Printf("%s", s)
//	}
}
