package core

import (
	"encoding/binary"
	"fmt"

	"github.com/go-ldap/ldap/v3"

	"github.com/charmbracelet/log"
)

func InitLdap(dc_ip string, username string, password string, domain string, isLDAPS bool) AHLdap {

	ldap := AHLdap{
		IP:           dc_ip,
		BindUsername: username,
		BindPassword: password,
		Domain:       domain,
		IsLDAPS:      isLDAPS,
	}

	return ldap
}

func Connect(ipAddress string, isLDAPS bool, username string, password string, domain string) (*ldap.Conn, error) {

	// TODO : LDAPS : conn, err := ldap.DialURL(fmt.Sprintf("ldaps://%s:636", ip))

	conn, err := ldap.DialURL(fmt.Sprintf("ldap://%s:389", ipAddress))

	if err != nil {
		log.Error("Failed to connect on the domain controller.")
		log.Error(err)
		return nil, err
	}

	// Besoin d'un format du type : "CN=svc_awx,OU=Service Accounts,DC=sleekboycompany,DC=com"
	// E.g. sur <http://10.77.101.192:8080/api/v2/settings/ldap/> - AUTH_LDAP_BIND_DN
	// svc_awx -> KO
	// svc_awx@sleekboycompany.com -> KO
	// SLEEKBOYCOMPANY\\svc_awx -> FONCTIONNE
	// SLEEKBOYCOMPANY.COM\\svc_awx -> KO

	dn := domain + "\\" + username

	err = conn.Bind(dn, password)

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

	// Byte 0: Revision (always 1)
	revision := sidBytes[0]

	// Byte 1: Number of sub-authorities (after the identifier authority)
	numSubAuthorities := sidBytes[1]

	// Bytes 2–7: Identifier Authority (big-endian, 48-bit)
	var authority uint64
	for i := 0; i < 6; i++ {
		authority = authority<<8 | uint64(sidBytes[2+i])
	}

	// Build base SID: S-revision-authority
	sid := fmt.Sprintf("S-%d-%d", revision, authority)

	// Sub-authorities (32-bit each, little-endian in the byte array!)
	offset := 8
	for i := uint8(0); i < numSubAuthorities; i++ {
		subAuthority := binary.LittleEndian.Uint32(sidBytes[offset : offset+4])
		sid += fmt.Sprintf("-%d", subAuthority)
		offset += 4
	}

	return sid
}
