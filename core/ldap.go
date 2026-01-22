package core

import (
	"crypto/tls"
	"encoding/binary"
	"fmt"

	"github.com/go-ldap/ldap/v3"

	"github.com/charmbracelet/log"
)

type AHLdap struct {
	IsLDAPS       bool
	IP            string
	BindUsername  string
	BindPassword  string
	Domain        string
	SkipVerifySSL bool
}

func InitLdap(dc_ip string, username string, password string, domain string, isLDAPS bool, skipVerifySSL bool) AHLdap {

	ldap := AHLdap{
		IP:            dc_ip,
		BindUsername:  username,
		BindPassword:  password,
		Domain:        domain,
		IsLDAPS:       isLDAPS,
		SkipVerifySSL: skipVerifySSL,
	}

	return ldap
}

func Connect(ldapObject AHLdap) (*ldap.Conn, error) {

	scheme := "ldap"
	port := 389
	var dialOpts []ldap.DialOpt

	if ldapObject.IsLDAPS {
		scheme = "ldaps"
		port = 636
		dialOpts = append(dialOpts, ldap.DialWithTLSConfig(&tls.Config{
			InsecureSkipVerify: ldapObject.SkipVerifySSL,
		}))
	}

	conn, err := ldap.DialURL(fmt.Sprintf("%s://%s:%d", scheme, ldapObject.IP, port), dialOpts...)

	if err != nil {
		log.Error("Failed to connect on the domain controller.")
		log.Error(err)
		return nil, err
	}

	dn := ldapObject.Domain + "\\" + ldapObject.BindUsername

	err = conn.Bind(dn, ldapObject.BindPassword)

	if err != nil {
		log.Fatalf("Authentication on the domain controller failed: %v", err)
		log.Error(err)
		return nil, err
	}

	return conn, err
}

func Search(conn *ldap.Conn, objectDN string) (string, error) {

	searchReq := ldap.NewSearchRequest(
		objectDN,
		ldap.ScopeBaseObject, ldap.NeverDerefAliases, 0, 0, false,
		"(objectClass=*)",
		[]string{"objectSid"},
		nil,
	)

	sr, err := conn.Search(searchReq)
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}

	if len(sr.Entries) != 1 {
		log.Fatalf("Expected 1 entry, got %d", len(sr.Entries))
	}

	entry := sr.Entries[0]
	objectSid := sidBytesToString(entry.GetRawAttributeValue("objectSid"))

	return objectSid, err
}

func sidBytesToString(sidBytes []byte) string {

	if len(sidBytes) < 8 {
		return "<invalid SID>"
	}

	revision := sidBytes[0]

	numSubAuthorities := sidBytes[1]

	var authority uint64
	for i := 0; i < 6; i++ {
		authority = authority<<8 | uint64(sidBytes[2+i])
	}

	sid := fmt.Sprintf("S-%d-%d", revision, authority)

	offset := 8
	for i := uint8(0); i < numSubAuthorities; i++ {
		subAuthority := binary.LittleEndian.Uint32(sidBytes[offset : offset+4])
		sid += fmt.Sprintf("-%d", subAuthority)
		offset += 4
	}

	return sid
}
