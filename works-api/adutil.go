package main

import (
	"crypto/tls"
	"fmt"
	ldap "github.com/go-ldap/ldap/v3"
	auth "github.com/korylprince/go-ad-auth/v3"
	//"github.com/sirupsen/logrus"

	//"github.com/sirupsen/logrus"
	"strings"
)

//var log = logrus.New().WithField("who", "AD")

var ADusername = "Administrator"
var ADpassword = "Ablecloud1!"
var domain = "dc1.local"
var ADserver = "dc1.local"
var ADport = 636
var ADbasedn = "DC=dc1,DC=local"


var username = "user"
var password = "Ablecloud1!"

var config = &auth.Config{
Server:   ADserver,
Port:     ADport,
BaseDN:   ADbasedn,
Security: auth.SecurityInsecureTLS,
}

func ConnectAD() (conn *auth.Conn, status bool, err error) {
	//connect
	conn, err = config.Connect()
	if err != nil {
		log.Errorln(err)
		return conn, false, err
	}

	//bind
	upn, err := config.UPN(ADusername)
	if err != nil {
		log.Errorln(err)
		upn = ADusername
	}
	status, err = conn.Bind(upn, ADpassword)
	if err != nil {
		log.Errorln(err)
		return conn, false, err
	}
	if !status {
		return conn, false, nil
	}
	return conn, status, err
}

func Auth(conn *auth.Conn, username string, password string) (status bool, err error){
	status, err = auth.Authenticate(config, username, password)
	return status, err
}
func listGroups(conn*auth.Conn, username string) (status bool, entry *ldap.Entry, userGroups []string, err error){

	setLog()
	upn, err := config.UPN(username)
	if err != nil {
		log.Errorln(err)
		upn = username
	}
	log.Errorln(upn)
	entry, err = conn.GetAttributes("userPrincipalName", upn, []string{"cn"})
	if err != nil {
		log.Errorln(entry)
		log.Errorln(err)
		return false, nil, nil, err
	}
	log.Infof("entry")
	for _, attrs := range entry.Attributes{
		attrs.PrettyPrint(4)
	}
	log.Errorln(err)
	const LDAPMatchingRuleInChain = "1.2.840.113556.1.4.1941"

	foundGroups, err := conn.Search(fmt.Sprintf("(member:%s:=%s)", LDAPMatchingRuleInChain, ldap.EscapeFilter(entry.DN)), []string{""}, 1000)
	if err != nil {
		return false, nil, nil, err
	}


	log.Infoln("foundGroups")
	for _, group := range foundGroups {
		userGroups=append(userGroups, group.DN)
	}
	log.Errorln(err)
	return status, entry, userGroups, nil
}
func inGroup(conn *auth.Conn, username string, groupname string) (bool, error) {
	_, _, groups, err := listGroups(conn, username)
	if err != nil {
		return false, err
	}
	for _, group := range groups{
		if strings.Compare(strings.TrimSpace(strings.ToLower(groupname)), strings.TrimSpace(strings.ToLower(group))) == 0{
			return true, nil
		}
	}

	return false, nil
}
/*
func InGroup(username string, groupname string){
	status, entry, groups, err := auth.AuthenticateExtended(config, fmt.Sprintf("%v", username), password, []string{"cn"}, []string{groupname})

	if err != nil {
		fmt.Println("handle error")
		fmt.Println(err)
	}

	if !status {
		fmt.Println("auth failed")
		return
	}

	fmt.Println(groups)
	if len(groups) == 0 {
		//handle user not being in any groups
		return
	}

	//get attributes
	cn := entry.GetAttributeValue("cn")

	fmt.Println(cn)
}
*/

func testLDAP(){
	setLog()
	// https://cybernetist.com/2020/05/18/getting-started-with-go-ldap/
	var filter = []string{
		"(cn=user)",
		"(cn=user2)",
		"(&(owner=*)(cn=user))",
		"(&(objectclass=rfc822mailgroup)(cn=*Computer*))",
		"(&(objectclass=rfc822mailgroup)(cn=*Mathematics*))"}
	var attributes = []string{
		"dn",
		"cn",
		"description"}

	ldapsServer := fmt.Sprintf("ldaps://%v:%v", ADserver, ADport)
	//ldapServer := fmt.Sprintf("ldap://%v:%v", ADserver, ADport)
	baseDN := ADbasedn
	l, err := ldap.DialURL(ldapsServer, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: true}))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	upn, err := config.UPN(ADusername)
	if err != nil {
		log.Errorln(err)
		upn = ADusername
	}
	err = l.Bind(upn, ADpassword)
	if err != nil {
		log.Fatal(err)
	}

	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, 0, false,
		filter[0],
		attributes,
		nil)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("TestSearch: %s -> num of entries = %d", searchRequest.Filter, len(sr.Entries))
	sr.PrettyPrint(4)



	groupname:="dev4"


	//ou add
	addReq := ldap.NewAddRequest(fmt.Sprintf("ou=%v,%v", groupname, ADbasedn), []ldap.Control{})
	addReq.Attribute("objectClass", []string{"top", "organizationalUnit"})
	addReq.Attribute("name", []string{groupname})
	addReq.Attribute("ou", []string{groupname})
	addReq.Attribute("instanceType", []string{fmt.Sprintf("%d", 0x00000004)})
	if err := l.Add(addReq); err != nil {
		log.Error("error adding OU:", addReq, err)
	}


	//group add
	addReq = ldap.NewAddRequest(fmt.Sprintf("cn=%v,ou=%v,%v",groupname, groupname, ADbasedn), []ldap.Control{})
	addReq.Attribute("objectClass", []string{"top", "group"})
	addReq.Attribute("name", []string{groupname})
	addReq.Attribute("sAMAccountName", []string{groupname})
	addReq.Attribute("instanceType", []string{fmt.Sprintf("%d", 0x00000004)})
	addReq.Attribute("groupType", []string{fmt.Sprintf("%d", 0x00000004 | 0x80000000)})
	if err := l.Add(addReq); err != nil {
		log.Error("error adding group:", addReq, err)
	}

	//computer add
	comname:="com4"
	addReq = ldap.NewAddRequest(fmt.Sprintf("cn=%v,cn=%v,%v",comname, "Computers", ADbasedn), []ldap.Control{})
	addReq.Attribute("objectClass", []string{"top", "computer", "organizationalPerson", "person", "user"})
	addReq.Attribute("name", []string{comname})
	addReq.Attribute("sAMAccountName", []string{fmt.Sprintf("%v$",comname)})
	addReq.Attribute("dNSHostName", []string{fmt.Sprintf("%v.%v",comname,domain)})
	addReq.Attribute("instanceType", []string{fmt.Sprintf("%d", 0x00000004)})
	if err := l.Add(addReq); err != nil {
		log.Error("error adding computer:", addReq, err)
	}

	//user add
	sn:="새"
	givenName:="사용자"
	initials:="이니셜"
	username:=fmt.Sprintf("%v %v %v.", sn, givenName, initials)
	accountname:="newuser"

	//telephoneNumber:="010"//일반->전화
	//pager:="011"//전화 -> 호출기
	//mobile:="016"//전화->휴대폰
	//facsimileTelephoneNumber:="팩스"//전화->팩스
	//homePhone:="042"//전화->집
	//ipPhone:="070"//전화->ip전화
	//postalCode:="35200"//주소->우편번호
	//countryCode :=410//주소->국가(한국?)
	//o:="Able" // ldap 회사명
	//manager:="CN=User3,CN=Users,DC=dc1,DC=local"//AD 상사이름
	//wWWHomePage:="https://www."//홈페이지 주소
	//c:="KR"//국가명
	//sAMAccountName:=accountname
	//mail:="ycyun@ablecloud.io"//메일주소
	userPrincipalName:=fmt.Sprintf("%v@%v",accountname, domain)
	//department:="개발팀"
	//st:=
	//	streetAddress:=
	//		physicalDeliveryOfficeName:=
	//			postOfficeBox:=
	//				I:=
	//					description:=
	//						company:=
	//


	addReq = ldap.NewAddRequest(fmt.Sprintf("cn=%v,cn=%v,%v",username, "Users", ADbasedn), []ldap.Control{})
	addReq.Attribute("objectClass", []string{"top", "organizationalPerson", "person", "user"})
	addReq.Attribute("name", []string{username})
	addReq.Attribute("sn", []string{sn})
	addReq.Attribute("givenName", []string{givenName})
	addReq.Attribute("initials", []string{initials})
	addReq.Attribute("displayName", []string{username})
	addReq.Attribute("cn", []string{username})
	addReq.Attribute("sAMAccountName", []string{accountname})
	addReq.Attribute("userPrincipalName", []string{userPrincipalName})
	addReq.Attribute("instanceType", []string{fmt.Sprintf("%d", 0x00000004)})
	if err := l.Add(addReq); err != nil {
		log.Error("error adding user:", addReq, err)
	}
}

