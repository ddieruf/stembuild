package package_test

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/onsi/gomega/gexec"

	"github.com/cloudfoundry-incubator/stembuild/test/helpers"

	"github.com/vmware/govmomi/govc/cli"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	_ "github.com/vmware/govmomi/govc/vm"
)

var _ = Describe("Package", func() {
	const (
		baseVMNameEnvVar                  = "PACKAGE_TEST_BASE_VM_NAME"
		vcenterURLVariable                = "VCENTER_BASE_URL"
		vcenterAdminCredentialUrlVariable = "VCENTER_ADMIN_CREDENTIAL_URL"
		vcenterFolderVariable             = "VM_FOLDER"
		existingVMVariable                = "EXISTING_SOURCE_VM"
		vcenterStembuildUsernameVariable  = "VCENTER_STEMBUILD_USER"
		vcenterStembuildPasswordVariable  = "VCENTER_STEMBUILD_PASSWORD"
		stembuildVersionVariable          = "STEMBUILD_VERSION"
	)

	var (
		workingDir                string
		baseVMName                string
		sourceVMName              string
		vmPath                    string
		vcenterURL                string
		vcenterAdminCredentialUrl string
		vcenterStembuildUsername  string
		vcenterStembuildPassword  string
		stembuildVersion          string
		executable                string
		err                       error
	)

	BeforeSuite(func() {

		stembuildVersion = helpers.EnvMustExist(stembuildVersionVariable)
		executable, err = helpers.BuildStembuild(stembuildVersion)
		Expect(err).NotTo(HaveOccurred())

		baseVMName = os.Getenv(baseVMNameEnvVar)
		Expect(baseVMName).NotTo(Equal(""), fmt.Sprintf("%s required", baseVMNameEnvVar))
	})

	BeforeEach(func() {
		existingVM := os.Getenv(existingVMVariable)
		vcenterFolder := helpers.EnvMustExist(vcenterFolderVariable)

		rand.Seed(time.Now().UnixNano())
		if existingVM == "" {
			sourceVMName = fmt.Sprintf("stembuild-package-test-%d", rand.Int())
		} else {
			sourceVMName = fmt.Sprintf("%s-%d", existingVM, rand.Int())
		}

		baseVMWithPath := fmt.Sprintf(vcenterFolder + "/" + baseVMName)
		vmPath = strings.Join([]string{vcenterFolder, sourceVMName}, "/")

		vcenterAdminCredentialUrl = helpers.EnvMustExist(vcenterAdminCredentialUrlVariable)

		cli.Run([]string{
			"vm.clone",
			"-vm", baseVMWithPath,
			"-on=false",
			"-u", vcenterAdminCredentialUrl,
			sourceVMName,
		})

		vcenterURL = helpers.EnvMustExist(vcenterURLVariable)
		vcenterStembuildUsername = helpers.EnvMustExist(vcenterStembuildUsernameVariable)
		vcenterStembuildPassword = helpers.EnvMustExist(vcenterStembuildPasswordVariable)

		workingDir, err = ioutil.TempDir(os.TempDir(), "stembuild-package-test")
		Expect(err).NotTo(HaveOccurred())
	})

	It("generates a stemcell with the correct shasum", func() {
		session := helpers.RunCommandInDir(
			workingDir,
			executable, "package",
			"-vcenter-url", vcenterURL,
			"-vcenter-username", vcenterStembuildUsername,
			"-vcenter-password", vcenterStembuildPassword,
			"-vm-inventory-path", vmPath,
		)

		Eventually(session, 60*time.Minute, 5*time.Second).Should(gexec.Exit(0))
		var out []byte
		session.Out.Write(out)
		fmt.Print(string(out))

		expectedOSVersion := strings.Split(stembuildVersion, ".")[0]
		expectedStemcellVersion := strings.Split(stembuildVersion, ".")[:2]

		expectedFilename := fmt.Sprintf(
			"bosh-stemcell-%s-vsphere-esxi-windows%s-go_agent.tgz", strings.Join(expectedStemcellVersion, "."), expectedOSVersion)

		stemcellPath := filepath.Join(workingDir, expectedFilename)

		image, err := os.Create(filepath.Join(workingDir, "image"))
		Expect(err).NotTo(HaveOccurred())
		copyFileFromTar(stemcellPath, "image", image)

		vmdkSha := sha1.New()
		ovfSha := sha1.New()

		imageMF, err := os.Create(filepath.Join(workingDir, "image.mf"))
		Expect(err).NotTo(HaveOccurred())

		copyFileFromTar(filepath.Join(workingDir, "image"), ".mf", imageMF)
		copyFileFromTar(filepath.Join(workingDir, "image"), ".vmdk", vmdkSha)
		copyFileFromTar(filepath.Join(workingDir, "image"), ".ovf", ovfSha)

		imageMFFile, err := helpers.ReadFile(filepath.Join(workingDir, "image.mf"))
		Expect(err).NotTo(HaveOccurred())
		Expect(imageMFFile).To(ContainSubstring(fmt.Sprintf("%x", vmdkSha.Sum(nil))))
		Expect(imageMFFile).To(ContainSubstring(fmt.Sprintf("%x", ovfSha.Sum(nil))))

	})

	AfterEach(func() {
		os.RemoveAll(workingDir)
		if vmPath != "" {
			cli.Run([]string{"vm.destroy", "-vm.ipath", vmPath, "-u", vcenterAdminCredentialUrl})
		}
	})
})

func copyFileFromTar(t string, f string, w io.Writer) {
	z, err := os.OpenFile(t, os.O_RDONLY, 0777)
	Expect(err).NotTo(HaveOccurred())
	gzr, err := gzip.NewReader(z)
	Expect(err).NotTo(HaveOccurred())
	defer gzr.Close()

	r := tar.NewReader(gzr)
	for {
		header, err := r.Next()
		if err == io.EOF {
			break
		}

		if strings.Contains(header.Name, f) {
			_, err = io.Copy(w, r)
			Expect(err).NotTo(HaveOccurred())
		}
	}
}
