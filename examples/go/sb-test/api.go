package main

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"os"

	googlebigtable "cloud.google.com/go/bigtable"
	"cloud.google.com/go/pubsub"
	googlespanner "cloud.google.com/go/spanner/admin/instance/apiv1"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/pivotal-golang/lager"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	googlebigquery "google.golang.org/api/bigquery/v2"
	"google.golang.org/api/option"
	storage "google.golang.org/api/storage/v1"
	instancepb "google.golang.org/genproto/googleapis/spanner/admin/instance/v1"
)

var projectId string

func AttachRoutes(router *mux.Router) {
	router.HandleFunc("/put-picture", putPicture).Methods("GET")
	router.HandleFunc("/test-storage", listBuckets).Methods("GET")
	router.HandleFunc("/test-pubsub", testPubSub).Methods("GET")
	router.HandleFunc("/test-bigquery", testBigquery).Methods("GET")
	router.HandleFunc("/test-bigtable", testBigtable).Methods("GET")
	router.HandleFunc("/test-cloudsql", testCloudSQL).Methods("GET")
	router.HandleFunc("/test-spanner", testSpanner).Methods("GET")
}

func NewAppRouter() http.Handler {
	r := mux.NewRouter()
	AttachRoutes(r)
	return r
}

// gets an authenticated config object
func getConfig(serviceName string) *jwt.Config {
	vcap := os.Getenv("VCAP_SERVICES")

	var asMap map[string][]map[string]interface{}
	if err := json.Unmarshal([]byte(vcap), &asMap); err != nil {
		println("ERROR setting config from json")
		println(err.Error())
		println(vcap)
	}

	credsInterface := asMap[serviceName][0]["credentials"]

	pkData := credsInterface.(map[string]interface{})["PrivateKeyData"].(string)

	projectId = credsInterface.(map[string]interface{})["ProjectId"].(string)

	decodedbytes, _ := base64.StdEncoding.DecodeString(pkData)

	conf, err := google.JWTConfigFromJSON(decodedbytes, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		println("ERROR setting config from json")
		println(err.Error())
	}
	return conf
}

// writes the service account credential to a local file
func writeSAJsonToFile(serviceName string) string {
	vcap := os.Getenv("VCAP_SERVICES")

	var asMap map[string][]map[string]interface{}
	if err := json.Unmarshal([]byte(vcap), &asMap); err != nil {
		println("ERROR setting config from json")
		println(err.Error())
		println(vcap)
	}

	credsInterface := asMap[serviceName][0]["credentials"]

	pkData := credsInterface.(map[string]interface{})["PrivateKeyData"].(string)

	decodedbytes, _ := base64.StdEncoding.DecodeString(pkData)

	file, err := os.Create(serviceName + ".json") // For read access.
	if err != nil {
		panic(err.Error())
	}
	file.Write(decodedbytes)
	file.Close()
	return serviceName + ".json"
}

// Puts a picture of the cloud foundry bunny in a given bucket under the key name test-img
// takes bucket_name as a get parameter
func putPicture(w http.ResponseWriter, req *http.Request) {

	reqParams := req.URL.Query()
	bucketName := reqParams["bucket_name"][0]

	conf := getConfig("google-storage")

	if conf == nil {
		respond(w, http.StatusOK, "No credentials found")
		return
	}
	storageService, _ := storage.New(conf.Client(context.Background()))

	object := &storage.Object{Name: "test-img"}

	out, err := os.Create("/tmp/test-img.png")
	defer out.Close()
	resp, err := http.Get("https://www.cloudfoundry.org/wp-content/uploads/2015/11/CF_rabbit_Blacksmith_rgb_trans_back-269x300.png")
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)

	file, err := os.Open("/tmp/test-img.png")
	if err != nil {
		println("Error opening file")
	}
	_, err = storageService.Objects.Insert(bucketName, object).Media(file).Do()
	if err != nil {
		println("error inserting object")
		println(err.Error())
	}

	respond(w, http.StatusOK, "success!")
}

// lists all buckets in the given project
func listBuckets(w http.ResponseWriter, req *http.Request) {

	conf := getConfig("google-storage")

	if conf == nil {
		respond(w, http.StatusOK, "No credentials found")
		return
	}
	storageService, err := storage.New(conf.Client(context.Background()))
	if err != nil {
		println(err.Error())
		respond(w, http.StatusInternalServerError, err.Error())
		return
	}

	buckets, err := storageService.Buckets.List(projectId).Do()

	if err != nil {
		println(err.Error())
		respond(w, http.StatusInternalServerError, err.Error())
		return
	}

	respond(w, http.StatusOK, buckets)
}

// returns the name of the first topic in the project
func testPubSub(w http.ResponseWriter, req *http.Request) {
	conf := getConfig("google-pubsub")
	filename := writeSAJsonToFile("google-pubsub")

	if conf == nil {
		respond(w, http.StatusOK, "No credentials found")
		return
	}
	pubsubService, err := pubsub.NewClient(context.Background(), projectId, option.WithServiceAccountFile(filename))
	if err != nil {
		println(err.Error())
		respond(w, http.StatusInternalServerError, err.Error())
		return
	}
	topic, err := pubsubService.Topics(context.Background()).Next()
	if err != nil {
		println(err.Error())
		respond(w, http.StatusInternalServerError, err.Error())
		return
	}

	respond(w, http.StatusOK, topic.String())
}

// lists all of the bigtable instances in the project
func testBigtable(w http.ResponseWriter, req *http.Request) {
	conf := getConfig("google-bigtable")
	if conf == nil {
		respond(w, http.StatusOK, "No credentials found")
		return
	}
	filename := writeSAJsonToFile("google-bigtable")
	service, err := googlebigtable.NewInstanceAdminClient(context.Background(), projectId, option.WithServiceAccountFile(filename))
	if err != nil {
		println(err.Error())
		respond(w, http.StatusInternalServerError, err.Error())
		return
	}
	instances, err := service.Instances(context.Background())
	if err != nil {
		println(err.Error())
		respond(w, http.StatusInternalServerError, err.Error())
		return
	}
	respond(w, http.StatusOK, instances)
}

// lists all of the bigquery datasets in the project
func testBigquery(w http.ResponseWriter, req *http.Request) {
	conf := getConfig("google-bigquery")
	service, err := googlebigquery.New(conf.Client(context.Background()))
	if err != nil {
		println(err.Error())
		respond(w, http.StatusInternalServerError, err.Error())
		return
	}
	list, err := service.Datasets.List(projectId).Do()
	if err != nil {
		println(err.Error())
		respond(w, http.StatusInternalServerError, err.Error())
		return
	}

	respond(w, http.StatusOK, list)
}

// connects to the given cloudsql database
func testCloudSQL(w http.ResponseWriter, req *http.Request) {
	vcap := os.Getenv("VCAP_SERVICES")

	var asMap map[string][]map[string]interface{}
	if err := json.Unmarshal([]byte(vcap), &asMap); err != nil {
		println("ERROR setting config from json")
		println(err.Error())
		println(vcap)
	}

	credsInterface := asMap["google-cloudsql"][0]["credentials"]
	caCert := credsInterface.(map[string]interface{})["CaCert"].(string)
	clientCertStr := credsInterface.(map[string]interface{})["ClientCert"].(string)
	clientKeyStr := credsInterface.(map[string]interface{})["ClientKey"].(string)

	dbHost := credsInterface.(map[string]interface{})["host"].(string)
	dbUsername := credsInterface.(map[string]interface{})["Username"].(string)
	dbPassword := credsInterface.(map[string]interface{})["Password"].(string)
	databaseName := credsInterface.(map[string]interface{})["database_name"].(string)
	dbPort := "3306"

	rootCertPool := x509.NewCertPool()

	if ok := rootCertPool.AppendCertsFromPEM([]byte(caCert)); !ok {
		respond(w, http.StatusInternalServerError, "error appending certs")
		return
	}
	clientCert := make([]tls.Certificate, 0, 1)

	certs, err := tls.X509KeyPair([]byte(clientCertStr), []byte(clientKeyStr))
	if err != nil {
		respond(w, http.StatusInternalServerError, err.Error())
		return
	}
	clientCert = append(clientCert, certs)
	mysql.RegisterTLSConfig("custom", &tls.Config{
		RootCAs:            rootCertPool,
		Certificates:       clientCert,
		InsecureSkipVerify: true,
	})

	tlsStr := "&tls=custom"
	if err != nil {
		respond(w, http.StatusInternalServerError, err.Error())
		return
	}

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local%s", dbUsername, dbPassword, dbHost, dbPort, databaseName, tlsStr)
	_, err = gorm.Open("mysql", connStr)
	if err != nil {
		respond(w, http.StatusInternalServerError, err.Error())
		return
	}

	respond(w, http.StatusOK, fmt.Sprintf("I connected to the %s database, yay!", databaseName))
}

func testSpanner(w http.ResponseWriter, req *http.Request) {
	conf := getConfig("google-spanner")
	if conf == nil {
		respond(w, http.StatusOK, "No credentials found")
		return
	}
	filename := writeSAJsonToFile("google-spanner")
	client, err := googlespanner.NewInstanceAdminClient(context.Background(), option.WithServiceAccountFile(filename))
	if err != nil {
		println(err.Error())
		respond(w, http.StatusInternalServerError, err.Error())
		return
	}
	listIter := client.ListInstances(context.Background(), &instancepb.ListInstancesRequest{
		Parent: "projects/" + projectId,
	})
	first, err := listIter.Next()
	if err != nil {
		println(err.Error())
		respond(w, http.StatusInternalServerError, err.Error())
		return
	}

	respond(w, http.StatusOK, first)
}

// reads from the given subscription name
func pullFromPubSub(w http.ResponseWriter, req *http.Request) {

	reqParams := req.URL.Query()
	subscriptionName := reqParams["subscription_name"][0]
	println(subscriptionName)

	conf := getConfig("google-pubsub")

	if conf == nil {
		respond(w, http.StatusOK, "No credentials found")
		return
	}
	pubsubService, _ := pubsub.NewClient(context.Background(), projectId, option.WithHTTPClient(conf.Client(context.Background())))

	cctx, cancel := context.WithCancel(context.Background())

	err := pubsubService.Subscription("projects/"+projectId+"/subscriptions/"+subscriptionName).Receive(cctx, func(fctx context.Context, m *pubsub.Message) {
		strMessage, _ := base64.StdEncoding.DecodeString(string(m.Data))
		cancel()
		respond(w, http.StatusOK, string(strMessage))

	})
	if err != nil {
		println(err.Error())
		respond(w, http.StatusOK, err.Error())
	}

}

func respond(w http.ResponseWriter, status int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	err := encoder.Encode(response)
	if err != nil {
		println("encoding response", err, lager.Data{"status": status, "response": response})
	}
}
