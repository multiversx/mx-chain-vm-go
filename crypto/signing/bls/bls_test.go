package bls

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	crypto "github.com/multiversx/mx-chain-crypto-go"
	"github.com/multiversx/mx-chain-crypto-go/signing"
	"github.com/multiversx/mx-chain-crypto-go/signing/mcl"
	llsig "github.com/multiversx/mx-chain-crypto-go/signing/mcl/multisig"
	"github.com/multiversx/mx-chain-crypto-go/signing/multisig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type multiSignerSetup struct {
	privKeys          [][]byte
	pubKeys           [][]byte
	partialSignatures [][][]byte
	messages          []string
	aggSignatures     [][]byte
}

const checkOK = "3e886a4c6e109a151f4105aee65a5192d150ef1fa68d3cd76964a0b086006dbe4324c989deb0e4416c6d6706db1b1910eb2732f08842fb4886067b9ed191109ac2188d76002d2e11da80a3f0ea89fee6b59c834cc478a6bd49cb8a193b1abb16@e96bd0f36b70c5ccc0c4396343bd7d8255b8a526c55fa1e218511fafe6539b8e@04725db195e37aa237cdbbda76270d4a229b6e7a3651104dc58c4349c0388e8546976fe54a04240530b99064e434c90f"
const checkNOK = "2c9a358953f61d34401d7ee4175eec105c476b18baacab371e2f47270035b539d84ad79ba587552b7e38802be00ff7148fc2a9c7a7034ff1e63ee24602ee952235ad14ca7d36e2be617fb2c99ed22a7a2729d86ae9fbb4df06f957ba07fec50e@1e46d9cbb995e30b82485525c29f80ac78aca295a6e88a11c3df8f9a445494bb@be8c460db180d6254c712ead3aa81935bc9be15b919dd45cb152b3dece04762569778c5e70e7af03fa1c66409d4f4711"

func TestBls_VerifyBLS(t *testing.T) {
	t.Parallel()

	b, _ := NewBLS()
	assert.Nil(t, b.VerifyBLS(splitString(t, checkOK)))
	assert.NotNil(t, b.VerifyBLS(splitString(t, checkNOK)))
}

func TestBls_VerifyBLSSigShare(t *testing.T) {
	t.Parallel()

	b, _ := NewBLS()
	assert.Nil(t, b.VerifySignatureShare(splitString(t, checkOK)))
	assert.NotNil(t, b.VerifySignatureShare(splitString(t, checkNOK)))
}

func TestBls_VerifyBLSMultiSig(t *testing.T) {
	t.Parallel()

	b, _ := NewBLS()

	numMessages := 5
	setupKOSK, multiSignerKOSK := createMultiSigSetupKOSK(uint16(numMessages), numMessages)
	setupKOSK.aggSignatures = aggregateSignatures(setupKOSK, multiSignerKOSK)

	for i := 0; i < len(setupKOSK.pubKeys); i++ {
		fmt.Println(hex.EncodeToString(setupKOSK.pubKeys[i]))
	}

	for i := 0; i < numMessages; i++ {
		fmt.Println(setupKOSK.messages[i])
		fmt.Println(hex.EncodeToString(setupKOSK.aggSignatures[i]))

		assert.Nil(t, b.VerifyAggregatedSig(setupKOSK.pubKeys, []byte(setupKOSK.messages[i]), setupKOSK.aggSignatures[i]))
		changedSig := make([]byte, len(setupKOSK.aggSignatures[i]))
		copy(changedSig, setupKOSK.aggSignatures[i])
		changedSig[0] += 1
		assert.NotNil(t, b.VerifyAggregatedSig(setupKOSK.pubKeys, []byte(setupKOSK.messages[i]), changedSig))
	}
}

func splitString(t testing.TB, str string) ([]byte, []byte, []byte) {
	split := strings.Split(str, "@")
	pkBuff, err := hex.DecodeString(split[0])
	require.Nil(t, err)

	msgBuff, err := hex.DecodeString(split[1])
	require.Nil(t, err)

	sigBuff, err := hex.DecodeString(split[2])
	require.Nil(t, err)

	return pkBuff, msgBuff, sigBuff
}

func createKeysAndMultiSignerBlsKOSK(
	grSize uint16,
	suite crypto.Suite,
) ([][]byte, [][]byte, crypto.MultiSigner) {

	kg, privKeys, pubKeys := createMultiSignerSetup(grSize, suite)
	llSigner := &llsig.BlsMultiSignerKOSK{}
	multiSigner, _ := multisig.NewBLSMultisig(llSigner, kg)

	return privKeys, pubKeys, multiSigner
}

func createMultiSignerSetup(grSize uint16, suite crypto.Suite) (crypto.KeyGenerator, [][]byte, [][]byte) {
	kg := signing.NewKeyGenerator(suite)
	privKeys := make([][]byte, grSize)
	pubKeys := make([][]byte, grSize)

	for i := uint16(0); i < grSize; i++ {
		sk, pk := kg.GeneratePair()
		privKeys[i], _ = sk.ToByteArray()
		pubKeys[i], _ = pk.ToByteArray()
	}
	return kg, privKeys, pubKeys
}

func createSignaturesShares(privKeys [][]byte, multiSigner crypto.MultiSigner, message []byte) [][]byte {
	sigShares := make([][]byte, len(privKeys))
	for i := uint16(0); i < uint16(len(privKeys)); i++ {
		sigShares[i], _ = multiSigner.CreateSignatureShare(privKeys[i], message)
	}

	return sigShares
}

func createMultiSigSetupKOSK(numSigners uint16, numMessages int) (*multiSignerSetup, crypto.MultiSigner) {
	var multiSigner crypto.MultiSigner
	setup := &multiSignerSetup{}
	suite := mcl.NewSuiteBLS12()
	setup.privKeys, setup.pubKeys, multiSigner = createKeysAndMultiSignerBlsKOSK(numSigners, suite)
	setup.messages, setup.partialSignatures = createMessagesAndPartialSignatures(numMessages, setup.privKeys, multiSigner)

	return setup, multiSigner
}

func createMessagesAndPartialSignatures(numMessages int, privKeys [][]byte, multiSigner crypto.MultiSigner) ([]string, [][][]byte) {
	partialSignatures := make([][][]byte, numMessages)
	messages := make([]string, numMessages)

	for i := 0; i < numMessages; i++ {
		messages[i] = fmt.Sprintf("message%d", i)
		signatures := createSignaturesShares(privKeys, multiSigner, []byte(messages[i]))
		partialSignatures[i] = signatures
	}

	return messages, partialSignatures
}

func aggregateSignatures(
	setup *multiSignerSetup,
	multiSigner crypto.MultiSigner,
) [][]byte {
	aggSignatures := make([][]byte, len(setup.messages))
	for i := 0; i < len(setup.messages); i++ {
		aggSignatures[i], _ = multiSigner.AggregateSigs(setup.pubKeys, setup.partialSignatures[i])
	}

	return aggSignatures
}
