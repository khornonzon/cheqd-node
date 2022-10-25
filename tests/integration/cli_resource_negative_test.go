//go:build integration

package integration

import (
	"crypto/ed25519"
	"fmt"

	"github.com/cheqd/cheqd-node/tests/integration/cli"
	"github.com/cheqd/cheqd-node/tests/integration/network"
	"github.com/cheqd/cheqd-node/tests/integration/testdata"
	cli_types "github.com/cheqd/cheqd-node/x/cheqd/client/cli"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/google/uuid"
	"github.com/multiformats/go-multibase"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tendermint/tendermint/libs/rand"
)

var _ = Describe("cheqd cli - negative resource", func() {
	var collectionId string
	var did string
	var signInputs []cli_types.SignInput
	var resourceId string
	var resourceName string

	BeforeEach(func() {
		collectionId = uuid.NewString()
		did := "did:cheqd:" + network.DID_NAMESPACE + ":" + collectionId
		keyId := did + "#key1"

		pubKey, privKey, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())

		pubKeyMultibase58, err := multibase.Encode(multibase.Base58BTC, pubKey)
		Expect(err).To(BeNil())

		payload := types.MsgCreateDidPayload{
			Id: did,
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                 keyId,
					Type:               "Ed25519VerificationKey2020",
					Controller:         did,
					PublicKeyMultibase: string(pubKeyMultibase58),
				},
			},
			Authentication: []string{keyId},
		}

		signInputs = []cli_types.SignInput{
			{
				VerificationMethodId: keyId,
				PrivKey:              privKey,
			},
		}

		res, err := cli.CreateDid(payload, signInputs, testdata.BASE_ACCOUNT_1)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))

		// Initialize shared resourceId
		resourceId = uuid.NewString()
		resourceName = "TestName"
	})

	It("cannot create resource with missing cli arguments, sign inputs mismatch", func() {
		// Generate a new DID Doc
		collectionId2 := uuid.NewString()

		// Generate extra sign inputs
		keyId2 := did + "#key2"
		_, privKey2, err := ed25519.GenerateKey(nil)
		Expect(err).To(BeNil())
		signInputs2 := []cli_types.SignInput{
			{
				VerificationMethodId: keyId2,
				PrivKey:              privKey2,
			},
		}

		// Fail to create a resource in non-existing collection
		resourceName = "TestResource"
		resourceType := "TestType"
		resourceFile, err := testdata.CreateTestJson(GinkgoT().TempDir())
		Expect(err).To(BeNil())

		_, err = cli.CreateResource(collectionId2, resourceId, resourceName, resourceType, resourceFile, signInputs, testdata.BASE_ACCOUNT_1)
		Expect(err).To(HaveOccurred())

		AddReportEntry("Integration", fmt.Sprintf("%sNegative: %s", cli.PURPLE, "cannot create resource with missing cli arguments"))
		// Fail to create a resource with missing cli arguments
		//   a. missing collection id
		_, err = cli.CreateResource("", resourceId, resourceName, resourceType, resourceFile, signInputs, testdata.BASE_ACCOUNT_1)
		Expect(err).To(HaveOccurred())

		//  b. missing resource id
		_, err = cli.CreateResource(collectionId, "", resourceName, resourceType, resourceFile, signInputs, testdata.BASE_ACCOUNT_1)
		Expect(err).To(HaveOccurred())

		// c. missing resource name
		_, err = cli.CreateResource(collectionId, resourceId, "", resourceType, resourceFile, signInputs, testdata.BASE_ACCOUNT_1)
		Expect(err).To(HaveOccurred())

		// d. missing resource type
		_, err = cli.CreateResource(collectionId, resourceId, resourceName, "", resourceFile, signInputs, testdata.BASE_ACCOUNT_1)
		Expect(err).To(HaveOccurred())

		// e. missing resource file
		_, err = cli.CreateResource(collectionId, resourceId, resourceName, resourceType, "", signInputs, testdata.BASE_ACCOUNT_1)
		Expect(err).To(HaveOccurred())

		// f. missing sign inputs
		_, err = cli.CreateResource(collectionId, resourceId, resourceName, resourceType, resourceFile, []cli_types.SignInput{}, testdata.BASE_ACCOUNT_1)
		Expect(err).To(HaveOccurred())

		// g. missing account
		_, err = cli.CreateResource(collectionId, resourceId, resourceName, resourceType, resourceFile, signInputs, "")
		Expect(err).To(HaveOccurred())

		AddReportEntry("Integration", fmt.Sprintf("%sNegative: %s", cli.PURPLE, "cannot create resource with sign inputs mismatch"))
		// Fail to create a resource with sign inputs mismatch
		//   a. sign inputs mismatch
		_, err = cli.CreateResource(collectionId, resourceId, resourceName, resourceType, resourceFile, signInputs2, testdata.BASE_ACCOUNT_1)
		Expect(err).To(HaveOccurred())

		//   b. non-existing key id
		_, err = cli.CreateResource(collectionId, resourceId, resourceName, resourceType, resourceFile, []cli_types.SignInput{
			{
				VerificationMethodId: "non-existing-key-id",
				PrivKey:              signInputs[0].PrivKey,
			},
		}, testdata.BASE_ACCOUNT_1)
		Expect(err).To(HaveOccurred())

		//   c. non-matching private key
		_, err = cli.CreateResource(collectionId, resourceId, resourceName, resourceType, resourceFile, []cli_types.SignInput{
			{
				VerificationMethodId: signInputs[0].VerificationMethodId,
				PrivKey:              privKey2,
			},
		}, testdata.BASE_ACCOUNT_1)
		Expect(err).To(HaveOccurred())

		//   d. invalid private key
		_, err = cli.CreateResource(collectionId, resourceId, resourceName, resourceType, resourceFile, []cli_types.SignInput{
			{
				VerificationMethodId: signInputs[0].VerificationMethodId,
				PrivKey:              testdata.GenerateByteEntropy(),
			},
		}, testdata.BASE_ACCOUNT_1)
		Expect(err).To(HaveOccurred())

		// Finally, create the resource
		res, err := cli.CreateResource(collectionId, resourceId, resourceName, resourceType, resourceFile, signInputs, testdata.BASE_ACCOUNT_1)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})

	It("cannot query a resource with missing cli arguments, non-existing collection, non-existing resource", func() {
		collectionId2 := uuid.NewString()
		resourceId2 := uuid.NewString()

		AddReportEntry("Integration", fmt.Sprintf("%sNegative: %s", cli.PURPLE, "cannot query a resource with missing cli arguments"))
		// Fail to query a resource with missing cli arguments
		//   a. missing collection id, resource id
		_, err := cli.QueryResource("", "")
		Expect(err).To(HaveOccurred())

		//   b. missing collection id
		_, err = cli.QueryResource("", resourceId2)
		Expect(err).To(HaveOccurred())

		//   c. missing resource id
		_, err = cli.QueryResource(collectionId2, "")
		Expect(err).To(HaveOccurred())

		AddReportEntry("Integration", fmt.Sprintf("%sNegative: %s", cli.PURPLE, "cannot query a resource with non-existing collection id"))
		// Fail to query a resource with non-existing collection id
		_, err = cli.QueryResource(collectionId2, resourceId)
		Expect(err).To(HaveOccurred())

		AddReportEntry("Integration", fmt.Sprintf("%sNegative: %s", cli.PURPLE, "cannot query a resource with non-existing resource id"))
		// Fail to query a resource with non-existing resource id
		_, err = cli.QueryResource(collectionId, resourceId2)
		Expect(err).To(HaveOccurred())
	})

	It("cannot query all resource versions with missing cli arguments, non-existing collection, non-existing resource", func() {
		collectionId2 := uuid.NewString()
		resourceName2 := rand.Str(10)
		resourceType := "TestType"

		AddReportEntry("Integration", fmt.Sprintf("%sNegative: %s", cli.PURPLE, "cannot query all resource versions with missing cli arguments"))
		// Fail to query all resource versions with missing cli arguments
		//   a. missing collection id, resource name, resource type
		_, err := cli.QueryAllResourceVersions("", "", "")
		Expect(err).To(HaveOccurred())

		//   b. missing collection id
		_, err = cli.QueryAllResourceVersions("", resourceName, resourceType)
		Expect(err).To(HaveOccurred())

		//   c. missing resource name
		_, err = cli.QueryAllResourceVersions(collectionId2, "", resourceType)
		Expect(err).To(HaveOccurred())

		AddReportEntry("Integration", fmt.Sprintf("%sNegative: %s", cli.PURPLE, "cannot query all resource versions with non-existing collection id"))
		// Fail to query all resource versions with non-existing collection id
		_, err = cli.QueryAllResourceVersions(collectionId2, resourceName, resourceType)
		Expect(err).To(HaveOccurred())

		AddReportEntry("Integration", fmt.Sprintf("%sNegative: %s", cli.PURPLE, "cannot query all resource versions with non-existing resource name"))
		// Fail to query all resource versions with non-existing resource name
		res, err := cli.QueryAllResourceVersions(collectionId, resourceName2, resourceType)
		Expect(err).To(BeNil())
		Expect(len(res.Resources)).To(BeEquivalentTo(0))
	})

	It("cannot query resource collection with missing cli arguments, non-existing collection id", func() {
		collectionId2 := uuid.NewString()

		AddReportEntry("Integration", fmt.Sprintf("%sNegative: %s", cli.PURPLE, "cannot query resource collection with missing cli arguments"))
		// Fail to query resource collection with missing cli arguments
		//   a. missing collection id
		_, err := cli.QueryResourceCollection("")
		Expect(err).To(HaveOccurred())

		AddReportEntry("Integration", fmt.Sprintf("%sNegative: %s", cli.PURPLE, "cannot query resource collection with non-existing collection id"))
		// Fail to query resource collection with non-existing collection id
		_, err = cli.QueryResourceCollection(collectionId2)
		Expect(err).To(HaveOccurred())
	})
})