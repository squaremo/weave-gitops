package server_test

import (
	"context"
	"log"
	"net"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	wego "github.com/weaveworks/weave-gitops/api/v1alpha1"
	pb "github.com/weaveworks/weave-gitops/pkg/api/applications"
	"github.com/weaveworks/weave-gitops/pkg/kube"
	"github.com/weaveworks/weave-gitops/pkg/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	apiruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	corev1 "k8s.io/api/core/v1"
)

func TestServer(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Server")

}

const bufSize = 1024 * 1024

var lis *bufconn.Listener

var s *grpc.Server
var apps pb.ApplicationsServer
var appsClient pb.ApplicationsClient
var conn *grpc.ClientConn
var err error
var k8sClient client.Client
var testEnv *envtest.Environment
var testClustername = "test-cluster"
var cfg *rest.Config
var scheme *apiruntime.Scheme
var k kube.Kube
var k8sManager ctrl.Manager

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

var _ = BeforeSuite(func(done Done) {
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{
			"../../manifests/crds",
			"../../tools/testcrds",
		},
	}

	var err error
	cfg, err = testEnv.Start()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())

	scheme = kube.CreateScheme()

	k8sManager, err = ctrl.NewManager(cfg, ctrl.Options{
		ClientDisableCacheFor: []client.Object{
			&wego.Application{},
			&corev1.Namespace{},
			&corev1.Secret{},
		},
		Scheme: scheme,
	})
	Expect(err).ToNot(HaveOccurred())

	go func() {
		err = k8sManager.Start(ctrl.SetupSignalHandler())
		Expect(err).ToNot(HaveOccurred())
	}()

	k8sClient = k8sManager.GetClient()
	Expect(k8sClient).ToNot(BeNil())
	close(done)
}, 60)

var _ = BeforeEach(func() {
	lis = bufconn.Listen(bufSize)
	s = grpc.NewServer()

	k = &kube.KubeHTTP{Client: k8sClient, ClusterName: testClustername}

	apps = server.NewApplicationsServer(k)
	pb.RegisterApplicationsServer(s, apps)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf(err.Error())
		}
	}()

	ctx := context.Background()
	conn, err = grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())

	Expect(err).NotTo(HaveOccurred())

	appsClient = pb.NewApplicationsClient(conn)
})

var _ = AfterEach(func() {
	conn.Close()
})
