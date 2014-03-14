# ansible_influxdb

A [golang](http://golang.org/) module for the [Ansible](http://www.ansible.com/) deployment program.
It can manage databases and users.

## Release

You can download the lastest version for [GNU/Linux amd64 here](https://github.com/novaquark/ansible_influxdb/releases/download/0.1.0/ansible_influxdb).

## Synopsis

Add or remove [InfluxDB](http://influxdb.org/) databases, cluster admins and database users.

Without any parameter, returns some facts.

## Options

<table>
  <tr>
    <th>parameter</th>
    <th>required</th>
    <th>default</th>
    <th>choices</th>
    <th>comments</th>
  </tr>
  <tr>
    <td>database</td>
    <td>yes with user_type = user or to create/remove</td>
    <td></td>
    <td></td>
    <td>name of the database to add or remove or change users</td>
  </tr>
  <tr>
    <td>login_host</td>
    <td>no</td>
    <td>localhost:8086</td>
    <td></td>
    <td>Host running the database and API port</td>
  </tr>
  <tr>
    <td>login_password</td>
    <td>no</td>
    <td></td>
    <td></td>
    <td>The password used to authenticate with</td>
  </tr>
  <tr>
    <td>login_username</td>
    <td>no</td>
    <td>root</td>
    <td></td>
    <td>The username used to authenticate with</td>
  </tr>
  <tr>
    <td>password</td>
    <td>no</td>
    <td></td>
    <td></td>
    <td>Set the user's password</td>
  </tr>
  <tr>
    <td>state</td>
    <td>no</td>
    <td>present</td>
    <td>present or absent</td>
    <td>The database, cluster admin or user state</td>
  </tr>
  <tr>
    <td>user_type</td>
    <td>no</td>
    <td></td>
    <td>cluster_admin or user</td>
    <td>User type to manage</td>
  </tr>
  <tr>
    <td>username</td>
    <td>yes when managing users</td>
    <td></td>
    <td></td>
    <td>name of the user to add or remove</td>
  </tr>
</table>

## Examples

### Create a new database with name 'koala'

	- ansible_influxdb: database=koala state=present

### Create database user with name 'bob' and password '12345' on database "koala"

	- ansible_influxdb: user_type=user username=bob password=12345 database=koala state=present

## Building

	cd $GOPATH
	mkdir -p src/github.com/novaquark/
	cd src/github.com/novaquark/
	git clone https://github.com/novaquark/ansible_influxdb.git
	cd ansible_influxdb
	go get
	go install
