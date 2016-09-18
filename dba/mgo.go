package dba

import (
	"errors"
	"fmt"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoDB describes a mongodb database object.
type MongoDB struct {
	mongodb *mgo.Database
	name    string
	user    string
	pass    string
	addr    string
}

// NewMongoDB creates a new mongodb database object.
func NewMongoDB(name, user, pass, addr string) (*MongoDB, error) {
	if user == "" {
		return nil, fmt.Errorf("missing mongodb username for %s", addr)
	}
	if pass == "" {
		return nil, fmt.Errorf("missing mongodb password for %s", addr)
	}
	if addr == "" {
		return nil, errors.New("missing mongodb address")
	}
	return &MongoDB{name: name, user: user, pass: pass, addr: addr}, nil
}

// Name satisfies the Database interface.  The mongodb database object returns
// the name of the last database the Dial method attempted to connect with.
func (m *MongoDB) Name() string {
	return m.name
}

// Dial satisfies the Database interface. The mongodb database object attempts
// to dial and start a session with a local or remote MongoDB instance.
func (m *MongoDB) Dial(name string) error {

	// catch any potential errors
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("error connecting to mongodb database %s: %s", name, r)
		}
	}()

	// parse connection string
	dial := fmt.Sprintf("mongodb://%s:%s@%s", m.user, m.pass, m.addr)
	dialInfo, err := mgo.ParseURL(dial)
	if err != nil {
		return err
	}

	// dial and establish database session
	s, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return err
	}

	// set db session and check the connection is not nil
	m.mongodb = s.DB(name)
	if m.mongodb == nil {
		return fmt.Errorf("error connection to mongodb database %s", name)
	}

	return nil
}

// Close satisfies the Database interface.  The mongodb database object closes
// the connection the Dial method created.
func (m *MongoDB) Close() {
	defer func() { // catch any potential errors
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()
	m.mongodb.Session.Close()
}

// BSONID takes a string ID and converts it to a bson.ObjectId. If the string
// cannot be converted an error is returned.
func BSONID(id string) (bson.ObjectId, error) {
	if !bson.IsObjectIdHex(id) {
		return "", errors.New("invalid id")
	}
	return bson.ObjectIdHex(id), nil
}
