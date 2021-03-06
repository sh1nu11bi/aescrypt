package stoutbox

import "bytes"
import "crypto/rand"
import "fmt"
import "io/ioutil"
import "math/big"
import "testing"

var testMessages = []string{
	"Hello, world.",
	"Yes... yes. This is a fertile land, and we will thrive. We will rule over all this land, and we will call it... This Land.",
	"Ah! Curse your sudden but inevitable betrayal!",
	"And I'm thinkin' you weren't burdened with an overabundance of schooling. So why don't we just ignore each other until we go away?",
	"Sir, I think you have a problem with your brain being missing.",
	"It's the way of life in my findings that journeys end when and where they want to; and that's where you make your home.",
	"I get confused. I remember everything. I remember too much. And... some of it's made up, and... some of it can't be quantified, and... there's secrets... and...",
	"Yeah, we're pretty much just giving each other significant glances and laughing incessantly.",
	"Jayne, go play with your rainstick.",
}

var (
	testBoxes   = make([]string, len(testMessages))
	testBoxFile []byte
	testGoodKey = PrivateKey{
		0x01, 0xd2, 0x21, 0x46, 0x28, 0x73, 0x25, 0x7f,
		0xf9, 0x74, 0x70, 0xe0, 0x65, 0xdd, 0xff, 0x18,
		0x03, 0x0b, 0x41, 0x30, 0x67, 0x24, 0x50, 0x47,
		0x6a, 0xe0, 0x44, 0xc7, 0x30, 0x05, 0xc1, 0x29,
		0xf8, 0x56, 0x5e, 0x6f, 0x15, 0xac, 0x20, 0xc7,
		0xe9, 0x32, 0xe1, 0x90, 0x15, 0xcc, 0xeb, 0x0a,
		0x86, 0x02, 0x9a, 0x9f, 0x8d, 0x22, 0x37, 0xc6,
		0x4f, 0x46, 0xde, 0xf0, 0x17, 0x26, 0x7f, 0x8a,
		0xa2, 0x64,
	}
	testGoodPub = PublicKey{
		0x04, 0x00, 0xbc, 0xcd, 0x20, 0x7d, 0xba, 0xb0,
		0x42, 0xd1, 0x41, 0x35, 0x8c, 0xd8, 0x26, 0x5d,
		0x3d, 0x92, 0x2c, 0x28, 0x58, 0xbb, 0xfa, 0xbe,
		0x67, 0xf8, 0xae, 0xd4, 0x01, 0xa8, 0xb7, 0x74,
		0xe0, 0x5f, 0xfe, 0x7b, 0x8b, 0x0b, 0x8a, 0x31,
		0x8e, 0xd8, 0x6d, 0xe6, 0xfe, 0x9e, 0xbb, 0x9f,
		0x1a, 0x40, 0xe5, 0x91, 0x62, 0x0a, 0xe4, 0x12,
		0xc8, 0xe2, 0x59, 0xa8, 0x14, 0xf8, 0xdc, 0x70,
		0x4f, 0xb3, 0x73, 0x00, 0xc0, 0x29, 0x1e, 0x60,
		0xd5, 0xa3, 0xe4, 0xcd, 0x90, 0x60, 0xae, 0x9d,
		0x45, 0xfa, 0xa8, 0x0f, 0x7f, 0xab, 0xf8, 0x06,
		0x54, 0xc4, 0xe4, 0x7e, 0x48, 0xe6, 0xa4, 0xea,
		0x2b, 0xf5, 0x6a, 0xfd, 0x82, 0x35, 0x63, 0xe1,
		0xe2, 0x04, 0xd5, 0xd7, 0xa1, 0x2e, 0x35, 0xcc,
		0xdf, 0x8a, 0xa1, 0x54, 0xa9, 0xef, 0x01, 0xae,
		0xde, 0x97, 0x1f, 0x3b, 0xa5, 0x68, 0xde, 0xd0,
		0x58, 0xf8, 0xa3, 0x1e, 0xd9,
	}
	testPeerKey = PrivateKey{
		0x00, 0x6f, 0x1b, 0x51, 0x3e, 0x67, 0x24, 0x0f,
		0xf7, 0x19, 0x77, 0x8d, 0x00, 0xe2, 0x94, 0xf7,
		0x27, 0x5f, 0x67, 0x9b, 0x60, 0x8e, 0x08, 0xae,
		0xc1, 0x0a, 0x79, 0x13, 0xc0, 0x21, 0x42, 0x52,
		0xe4, 0xf0, 0x3c, 0xe1, 0x8a, 0xc2, 0x6b, 0x1d,
		0x0a, 0x91, 0xb3, 0x19, 0xc2, 0xd6, 0x7b, 0x04,
		0x25, 0xab, 0x5b, 0x9a, 0x1d, 0xd3, 0x76, 0xca,
		0xd2, 0x9b, 0x4a, 0xdd, 0x60, 0x97, 0x31, 0x53,
		0x8b, 0x8e,
	}
	testPeerPub = PublicKey{
		0x04, 0x00, 0x1a, 0x94, 0x55, 0x77, 0x4f, 0x1e,
		0x47, 0x5a, 0x34, 0x81, 0xf7, 0xa0, 0x4a, 0xc9,
		0x73, 0x61, 0x23, 0x8d, 0xda, 0x72, 0x2a, 0xd9,
		0x0c, 0xd4, 0x14, 0x18, 0x29, 0x99, 0x23, 0xa3,
		0x0c, 0x03, 0xde, 0x93, 0x63, 0x50, 0x55, 0x99,
		0xa6, 0x6b, 0xa3, 0x1e, 0x45, 0xbb, 0xc8, 0xd1,
		0x04, 0x5a, 0xe9, 0x9f, 0x28, 0xc1, 0xc6, 0x6e,
		0xfc, 0x5e, 0xb3, 0x4a, 0x81, 0x9f, 0x61, 0x75,
		0xdf, 0xa9, 0x89, 0x01, 0xd6, 0x29, 0x76, 0xe4,
		0x3c, 0xf3, 0x2e, 0x43, 0xdb, 0x26, 0x0e, 0xf5,
		0xab, 0xd4, 0x58, 0xc4, 0x5d, 0x8c, 0xe5, 0x0c,
		0x5b, 0x02, 0xb7, 0xe2, 0x4c, 0xa3, 0x65, 0x88,
		0x59, 0x75, 0xcf, 0x22, 0x1a, 0x94, 0x0f, 0x5d,
		0xbb, 0x9f, 0xe0, 0x18, 0x01, 0xe9, 0x23, 0x4f,
		0x72, 0xa9, 0x9a, 0x1b, 0x47, 0xa0, 0x9e, 0x55,
		0x9f, 0x8a, 0xfb, 0x1c, 0xf0, 0x4f, 0xf8, 0xd4,
		0x2e, 0xda, 0xa5, 0x23, 0x96,
	}
	testBadKey = PrivateKey{
		0x01, 0xe2, 0xc5, 0x9c, 0x6e, 0x6b, 0xf0, 0x19,
		0x84, 0xfa, 0x0a, 0x2a, 0x2e, 0xbe, 0x9c, 0x13,
		0xc8, 0x5e, 0x35, 0xc7, 0x85, 0x5f, 0x3b, 0xe9,
		0x95, 0xfc, 0x03, 0x42, 0xe3, 0xd7, 0xa7, 0x49,
		0xde, 0x45, 0xf1, 0xc5, 0x09, 0xbe, 0x3b, 0x22,
		0x1e, 0x96, 0x8d, 0xa4, 0x71, 0x91, 0x3e, 0x89,
		0x77, 0x9e, 0xcf, 0x3e, 0x7d, 0xf3, 0xf9, 0x37,
		0x69, 0xe0, 0xd5, 0x51, 0x72, 0xcc, 0x84, 0x03,
		0x39, 0x9b,
	}
	testBadPub = PublicKey{
		0x04, 0x00, 0x1e, 0x00, 0xa8, 0xcb, 0x1a, 0xd8,
		0x7c, 0x9c, 0xcf, 0xd8, 0xdb, 0xf4, 0x27, 0xc4,
		0x8a, 0x70, 0x01, 0x23, 0xbe, 0xd0, 0x37, 0x0e,
		0x10, 0x8e, 0x07, 0x3c, 0x28, 0x7c, 0x45, 0x6f,
		0xa8, 0x61, 0x3d, 0xf9, 0x64, 0xca, 0x20, 0xfb,
		0xe9, 0x60, 0xea, 0x2a, 0xa2, 0xeb, 0x5a, 0x69,
		0x98, 0x6d, 0x7e, 0x9f, 0x11, 0x10, 0x06, 0x27,
		0x57, 0xcb, 0x4c, 0xf9, 0x72, 0x0e, 0x0f, 0x9b,
		0xa2, 0xe8, 0x54, 0x00, 0x1f, 0x71, 0x78, 0xdd,
		0x0a, 0xea, 0x5d, 0xe7, 0xda, 0xa8, 0x70, 0x3e,
		0xa4, 0x42, 0x49, 0xa6, 0xb4, 0x94, 0x20, 0x36,
		0x83, 0x35, 0xc8, 0xf7, 0x1b, 0xae, 0x76, 0x17,
		0xfc, 0x02, 0x10, 0xa0, 0x65, 0x66, 0x19, 0x1c,
		0x28, 0xc0, 0xed, 0x05, 0x00, 0x0c, 0xf2, 0x0b,
		0x72, 0xe1, 0x4c, 0x45, 0x39, 0xb4, 0x88, 0xf1,
		0xd8, 0xbd, 0xa7, 0xcd, 0xd3, 0x15, 0x14, 0xaf,
		0xbf, 0x6b, 0x1a, 0xd9, 0x5f,
	}
)

func randInt(max int64) int64 {
	maxBig := big.NewInt(max)
	n, err := rand.Int(PRNG, maxBig)
	if err != nil {
		return -1
	}
	return n.Int64()
}

func mutate(in []byte) (out []byte) {
	out = make([]byte, len(in))
	copy(out, in)

	iterations := (randInt(int64(len(out))) / 2) + 1
	if iterations == -1 {
		panic("mutate failed")
	}
	for i := 0; i < int(iterations); i++ {
		mByte := randInt(int64(len(out)))
		mBit := randInt(7)
		if mBit == -1 || mByte == -1 {
			panic("mutate failed")
		}
		out[mByte] ^= (1 << uint(mBit))
	}
	if bytes.Equal(out, in) {
		panic("mutate failed")
	}
	return out
}

/*
func TestKeyGeneration(t *testing.T) {
	var ok bool
	testGoodKey, testGoodPub, ok = GenerateKey()
	if !ok {
		fmt.Println("key generation failed")
		t.FailNow()
	}
	testPeerKey, testPeerPub, ok = GenerateKey()
	if !ok {
		fmt.Println("key generation failed")
		t.FailNow()
	}
	testBadKey, testBadPub, ok = GenerateKey()
	if !ok {
		fmt.Println("key generation failed")
		t.FailNow()
	}
	ioutil.WriteFile("testvectors/good.key", testGoodKey, 0644)
	ioutil.WriteFile("testvectors/good.pub", testGoodPub, 0644)
	ioutil.WriteFile("testvectors/peer.key", testPeerKey, 0644)
	ioutil.WriteFile("testvectors/peer.pub", testPeerPub, 0644)
	ioutil.WriteFile("testvectors/bad.key", testBadKey, 0644)
	ioutil.WriteFile("testvectors/bad.pub", testBadPub, 0644)
}
*/

func TestBoxing(t *testing.T) {
	for i := 0; i < len(testMessages); i++ {
		box, ok := Seal([]byte(testMessages[i]), testPeerPub)
		if !ok {
			fmt.Println("Boxing failed: message", i)
			t.FailNow()
		} else if len(box) != len(testMessages[i])+Overhead {
			fmt.Println("The box length is invalid.")
			t.FailNow()
		}
		testBoxes[i] = string(box)
		/*
			fileName := fmt.Sprintf("testvectors/test_vector-%d.bin", i+1)
			ioutil.WriteFile(fileName, []byte(testMessages[i]), 0644)
			fileName = fmt.Sprintf("testvectors/test_box-%d.bin", i+1)
			ioutil.WriteFile(fileName, box, 0644)
		*/
	}
}

func TestUnboxing(t *testing.T) {
	for i := 0; i < len(testMessages); i++ {
		message, ok := Open([]byte(testBoxes[i]), testPeerKey)
		if !ok {
			fmt.Println("Unboxing failed: message", i)
			t.FailNow()
		} else if string(message) != testMessages[i] {
			fmt.Printf("Unboxing failed: expected '%s', got '%s'\n",
				testMessages[i], string(message))
			t.FailNow()
		}
	}
}

func TestBadUnboxing(t *testing.T) {
	for i := 0; i < len(testMessages); i++ {
		_, ok := Open([]byte(testBoxes[i]), testBadKey)
		if ok {
			fmt.Println("Unboxing should have failed: message", i)
			t.FailNow()
		}
		_, ok = Open(mutate([]byte(testBoxes[i])), testGoodKey)
		if ok {
			fmt.Println("Unboxing should have failed: message", i)
			t.FailNow()
		}
	}
}

func TestSignedBoxing(t *testing.T) {
	for i := 0; i < len(testMessages); i++ {
		box, ok := SignAndSeal([]byte(testMessages[i]), testGoodKey, testGoodPub, testPeerPub)
		if !ok {
			fmt.Println("Signed boxing failed: message", i)
			t.FailNow()
		} else if len(box) < len(testMessages[i])+Overhead+sigSize {
			fmt.Println("The box length is invalid.")
			t.FailNow()
		}
		testBoxes[i] = string(box)
		/*
			fileName := fmt.Sprintf("testvectors/test_signed_vector-%d.bin", i+1)
			ioutil.WriteFile(fileName, []byte(testMessages[i]), 0644)
			fileName = fmt.Sprintf("testvectors/test_signed_box-%d.bin", i+1)
			ioutil.WriteFile(fileName, box, 0644)
		*/
	}
}

func TestSignedUnboxing(t *testing.T) {
	for i := 0; i < len(testMessages); i++ {
		message, ok := OpenAndVerify([]byte(testBoxes[i]), testPeerKey, testGoodPub)
		if !ok {
			fmt.Println("Signed unboxing failed: message", i)
			t.FailNow()
		} else if string(message) != testMessages[i] {
			fmt.Printf("Signed unboxing failed: expected '%s', got '%s'\n",
				testMessages[i], string(message))
			t.FailNow()
		}
	}
}

func TestSignedBadUnboxing(t *testing.T) {
	for i := 0; i < len(testMessages); i++ {
		_, ok := OpenAndVerify([]byte(testBoxes[i]), testPeerKey, testBadPub)
		if ok {
			fmt.Println("Unboxing should have failed: message", i)
			t.FailNow()
		} else if _, ok = OpenAndVerify(mutate([]byte(testBoxes[i])), testPeerKey, testGoodPub); ok {
			fmt.Println("Unboxing should have failed: message", i)
			t.FailNow()
		} else if _, ok = OpenAndVerify(mutate([]byte(testBoxes[i])), testPeerKey, testGoodPub); ok {
			fmt.Println("Unboxing should have failed: message", i)
			t.FailNow()
		}
	}
}

// TestLargerBox tests the encryption of a 4,026 byte test file.
func TestLargerBox(t *testing.T) {
	var err error
	testBoxFile, err = ioutil.ReadFile("testdata/TEST.txt")
	if err != nil {
		fmt.Println("Failed to read test data:", err.Error())
		t.FailNow()
	}

	box, ok := Seal(testBoxFile, testPeerPub)
	if !ok {
		fmt.Println("Failed to box message.")
		t.FailNow()
	}

	message, ok := Open(box, testPeerKey)
	if !ok {
		fmt.Println("Failed to unbox message.")
		t.FailNow()
	}

	if !bytes.Equal(message, testBoxFile) {
		fmt.Println("Recovered message is invalid.")
		t.FailNow()
	}
}

func BenchmarkUnsignedSeal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, ok := Seal(testBoxFile, testPeerPub)
		if !ok {
			fmt.Println("Couldn't seal message: benchmark aborted.")
			b.FailNow()
		}
	}
}

func BenchmarkSignAndSeal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, ok := SignAndSeal(testBoxFile, testGoodKey, testGoodPub, testPeerPub)
		if !ok {
			fmt.Println("Couldn't seal message: benchmark aborted.")
			b.FailNow()
		}
	}
}

// Benchmark the Open function, which retrieves a message from a box.
func BenchmarkUnsignedOpen(b *testing.B) {
	box, ok := Seal(testBoxFile, testPeerPub)
	if !ok {
		fmt.Println("Can't seal message: benchmark aborted.")
		b.FailNow()
	}
	for i := 0; i < b.N; i++ {
		_, ok := Open(box, testPeerKey)
		if !ok {
			fmt.Println("Couldn't open message: benchmark aborted.")
			b.FailNow()
		}
	}
}

// Benchmark the OpenSigned function, which retrieves a message from a box and verifies a
// signature on it.
func BenchmarkOpenSigned(b *testing.B) {
	box, ok := SignAndSeal(testBoxFile, testGoodKey, testGoodPub, testPeerPub)
	if !ok {
		fmt.Println("Can't seal message: benchmark aborted.")
		b.FailNow()
	}
	for i := 0; i < b.N; i++ {
		_, ok := OpenAndVerify(box, testPeerKey, testGoodPub)
		if !ok {
			fmt.Println("Couldn't open message: benchmark aborted.")
			b.FailNow()
		}
	}
}

// Benchmark the SharedKey function.
func BenchmarkSharedKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, ok := SharedKey(testGoodKey, testPeerPub)
		if !ok {
			fmt.Println("Computing shared key failed: benchmark aborted.")
			b.FailNow()
		}
	}
}

/*
func TestSharedKeyPairs(t *testing.T) {
	for i := 0; i < 4; i++ {
		p_priv, p_pub, ok := GenerateKey()
		if !ok {
			fmt.Println("Failed to generate peer key.")
			t.FailNow()
		}
		privFN := fmt.Sprintf("testvectors/peer_%d.key", i)
		pubFN := fmt.Sprintf("testvectors/peer_%d.pub", i)
		ioutil.WriteFile(privFN, p_priv, 0644)
		ioutil.WriteFile(pubFN, p_pub, 0644)
	}
}
*/

var peerPublicList = []PublicKey{
	PublicKey{
		0x04, 0x00, 0x13, 0xd9, 0xcc, 0x7c, 0x96, 0x06,
		0x6f, 0x2d, 0x45, 0x5b, 0xda, 0x81, 0x36, 0xaf,
		0x26, 0xa0, 0x6b, 0x31, 0x88, 0x90, 0x69, 0x30,
		0x0a, 0x9d, 0x7d, 0xd5, 0x35, 0x9c, 0xd4, 0x94,
		0x85, 0xf9, 0x34, 0x16, 0x73, 0x2c, 0xda, 0x16,
		0x35, 0x4d, 0xdd, 0x45, 0xf4, 0x52, 0xdf, 0x9b,
		0x57, 0xfa, 0x51, 0x00, 0x11, 0xcb, 0x33, 0x47,
		0x66, 0xb4, 0x79, 0xe3, 0x68, 0x4b, 0xf6, 0x1f,
		0xf3, 0x9b, 0xdd, 0x00, 0xd3, 0x54, 0x86, 0xb3,
		0xa1, 0x0b, 0x59, 0x6f, 0x21, 0x54, 0xe7, 0x5b,
		0x04, 0xab, 0x3a, 0xbf, 0xc7, 0xee, 0x16, 0x52,
		0xc4, 0xf7, 0x57, 0xbd, 0xae, 0x23, 0x5b, 0xa1,
		0x6d, 0xf2, 0x25, 0x14, 0xcc, 0xe5, 0x75, 0x9f,
		0x70, 0xd6, 0x79, 0x66, 0xa6, 0x2a, 0x72, 0xe8,
		0x4f, 0x9c, 0x1e, 0xdc, 0x8a, 0xaa, 0x10, 0x8a,
		0xb9, 0x49, 0xd1, 0xdb, 0x1f, 0xad, 0xd8, 0xf1,
		0x82, 0xe9, 0xf5, 0x64, 0x9f,
	},
	PublicKey{
		0x04, 0x00, 0x71, 0x45, 0x3b, 0x48, 0x56, 0x13,
		0xad, 0x7a, 0xa4, 0xb9, 0xd4, 0xed, 0x74, 0x16,
		0x45, 0x54, 0x38, 0xdf, 0x94, 0xba, 0xf7, 0x27,
		0x1b, 0x2f, 0x52, 0xe0, 0x31, 0xf9, 0xa0, 0x1d,
		0x19, 0xd5, 0x0b, 0x5d, 0xda, 0x21, 0xc6, 0x7a,
		0x01, 0xf0, 0xdf, 0x1e, 0xbf, 0xf2, 0xe6, 0x48,
		0xde, 0x71, 0xd9, 0x60, 0x12, 0x0d, 0x0f, 0xdc,
		0x40, 0x08, 0x54, 0xd6, 0x29, 0xac, 0x9b, 0xd5,
		0xc5, 0x97, 0xf6, 0x00, 0xf4, 0xff, 0xe6, 0x9f,
		0xa3, 0x41, 0x9e, 0x5e, 0x2c, 0xbd, 0xb4, 0x40,
		0x08, 0xb4, 0x24, 0xa5, 0x48, 0x74, 0xe8, 0xa4,
		0x9f, 0x65, 0x42, 0xc0, 0x5e, 0x45, 0x11, 0xf7,
		0x74, 0xf9, 0x89, 0x26, 0xf3, 0x8e, 0x42, 0x82,
		0x6b, 0x4c, 0xe4, 0xe5, 0xaa, 0x26, 0x23, 0x24,
		0x77, 0x53, 0x1d, 0x2f, 0xa2, 0x56, 0x67, 0x6e,
		0x0b, 0xb5, 0x22, 0xd1, 0x86, 0x81, 0x33, 0xfa,
		0x6d, 0xe3, 0xbf, 0xeb, 0x94,
	},
	PublicKey{
		0x04, 0x01, 0xd6, 0x5b, 0xe9, 0x99, 0xda, 0x7e,
		0x02, 0x25, 0x54, 0x98, 0xb4, 0x49, 0xf4, 0x0f,
		0x45, 0x3b, 0xcf, 0xa2, 0x21, 0x75, 0x58, 0x1a,
		0xbf, 0x45, 0xaf, 0x64, 0xa9, 0xef, 0x19, 0x22,
		0xb8, 0xb2, 0x5e, 0x4e, 0xb2, 0x0a, 0x62, 0x97,
		0x20, 0x62, 0xa4, 0x6f, 0x10, 0x34, 0xc7, 0x27,
		0xc1, 0x88, 0x39, 0xdf, 0x40, 0x78, 0xd3, 0x65,
		0x43, 0x3c, 0x68, 0xe3, 0x30, 0xb6, 0x1b, 0x06,
		0x8b, 0x55, 0x8b, 0x01, 0x71, 0xd4, 0xa6, 0x07,
		0x00, 0xb9, 0xa4, 0x2d, 0xbe, 0x26, 0x97, 0x2f,
		0xe7, 0xf1, 0x0b, 0x1a, 0xf8, 0x43, 0x93, 0x0c,
		0xde, 0xdf, 0x25, 0x84, 0x22, 0x69, 0xe9, 0x6e,
		0xcc, 0x5a, 0x6c, 0x1e, 0x8b, 0x34, 0xdb, 0x6f,
		0xa5, 0x52, 0xd8, 0x59, 0x39, 0x00, 0xec, 0x2e,
		0xc0, 0x49, 0x92, 0xfb, 0xe9, 0x4b, 0x67, 0x97,
		0xd6, 0xbd, 0x1d, 0x47, 0xdd, 0xff, 0xa6, 0xee,
		0xd7, 0x22, 0xac, 0x57, 0xb6,
	},
	PublicKey{
		0x04, 0x01, 0xf6, 0x2c, 0x3a, 0x95, 0x72, 0x4a,
		0x42, 0x0c, 0x18, 0xbc, 0x3b, 0xa5, 0x70, 0x6f,
		0xfa, 0x8c, 0xde, 0x41, 0xf3, 0x41, 0x49, 0x7c,
		0xff, 0x2b, 0x26, 0x00, 0xef, 0x48, 0x63, 0xa6,
		0x08, 0x64, 0x93, 0xc3, 0x10, 0x0b, 0x07, 0x79,
		0xf9, 0x1b, 0x60, 0x60, 0xc1, 0xbf, 0x9b, 0x20,
		0xaa, 0x0e, 0xc7, 0xc2, 0x64, 0x9d, 0xab, 0xfb,
		0x69, 0xbf, 0xc1, 0x68, 0xa5, 0xb7, 0x14, 0xef,
		0xc3, 0x64, 0x5c, 0x00, 0x29, 0x07, 0xaf, 0x35,
		0x48, 0x07, 0x6d, 0x6a, 0xe4, 0x10, 0xc3, 0x4a,
		0xae, 0x2c, 0xb7, 0x79, 0x06, 0x62, 0xd2, 0x98,
		0xcc, 0xb1, 0x3c, 0x24, 0x05, 0xc5, 0xbf, 0xe3,
		0xcb, 0xcb, 0xd5, 0xf2, 0x5e, 0xf7, 0x62, 0x30,
		0x8c, 0x75, 0x0f, 0xe2, 0xce, 0xee, 0x61, 0x20,
		0x4d, 0xf5, 0x3e, 0xe7, 0xce, 0xbf, 0x09, 0x68,
		0x9c, 0x0e, 0xf5, 0x76, 0x64, 0xc3, 0x4e, 0x35,
		0x54, 0x4d, 0x13, 0xd5, 0x58,
	},
}

var peerPrivList = []PrivateKey{
	PrivateKey{
		0x00, 0x83, 0x32, 0x4d, 0x6c, 0x29, 0x61, 0xc0,
		0xe1, 0x5c, 0xd3, 0x32, 0x1f, 0x37, 0x28, 0xc1,
		0x9b, 0x6a, 0x3e, 0x81, 0xe6, 0xfb, 0xb6, 0x19,
		0x7d, 0x9a, 0x6d, 0x7e, 0x16, 0x09, 0xff, 0x19,
		0xbc, 0x54, 0xc3, 0x8a, 0x13, 0xae, 0x5c, 0xb9,
		0x29, 0xc5, 0x76, 0x3e, 0x41, 0xa1, 0xaa, 0xbf,
		0x25, 0x5a, 0x81, 0x56, 0x07, 0x2b, 0xe9, 0x97,
		0x7a, 0xad, 0xb5, 0x25, 0x80, 0x65, 0x10, 0x74,
		0x1b, 0x60,
	},
	PrivateKey{
		0x01, 0x39, 0xe3, 0x43, 0x2c, 0x6c, 0x19, 0x12,
		0x51, 0xa9, 0xc7, 0x5c, 0x01, 0xe1, 0xe5, 0x0c,
		0x4a, 0x8f, 0xc5, 0xf8, 0x02, 0x75, 0x2b, 0x94,
		0x1b, 0xc8, 0x70, 0x07, 0x22, 0x58, 0xdc, 0x2f,
		0x18, 0xd7, 0xef, 0xd3, 0x7a, 0x53, 0x91, 0xf7,
		0xaa, 0x0d, 0xa3, 0xae, 0x8b, 0xdc, 0x4b, 0xc1,
		0xdb, 0x42, 0xa3, 0x3c, 0x9c, 0x38, 0x6f, 0xf8,
		0x64, 0xdf, 0x4f, 0x80, 0x72, 0x1f, 0xa6, 0x9f,
		0x66, 0x5e,
	},
	PrivateKey{
		0x01, 0x8b, 0x48, 0x3b, 0x24, 0xbb, 0x77, 0x65,
		0x2e, 0xa1, 0x47, 0x15, 0xc1, 0x3d, 0x50, 0x19,
		0xb7, 0x14, 0x4b, 0x3d, 0x7c, 0x21, 0x4c, 0x4c,
		0x97, 0x80, 0x14, 0x07, 0x52, 0x27, 0xb8, 0xe3,
		0x1a, 0x3c, 0x04, 0x3b, 0x0a, 0x43, 0xfd, 0x7f,
		0xd5, 0xa6, 0xd1, 0x0e, 0x5a, 0x55, 0x17, 0x7b,
		0x03, 0xf4, 0x5d, 0x42, 0x91, 0x63, 0xca, 0x32,
		0x53, 0xe6, 0x92, 0xe7, 0x98, 0x27, 0x4c, 0x4d,
		0x47, 0x8c,
	},
	PrivateKey{
		0x00, 0x0b, 0x32, 0xf4, 0xe0, 0x1c, 0x77, 0xc9,
		0x21, 0x99, 0xfe, 0xba, 0xac, 0xe5, 0xb8, 0xfd,
		0xd2, 0x93, 0xa0, 0x66, 0xba, 0x7a, 0x58, 0x5d,
		0x2e, 0xe0, 0x29, 0x4d, 0x77, 0xb5, 0x24, 0xb1,
		0x70, 0xda, 0x4d, 0xb4, 0xa7, 0x54, 0x3f, 0x46,
		0x66, 0x1d, 0x46, 0x65, 0xb8, 0x23, 0x49, 0x76,
		0xad, 0x17, 0x10, 0xbe, 0x40, 0xee, 0x95, 0x32,
		0xb0, 0x97, 0x74, 0xf9, 0x3e, 0x8d, 0x6e, 0xc2,
		0xc9, 0xd0,
	},
}

func TestSharedBoxing(t *testing.T) {
	for i := 0; i < len(testMessages); i++ {
		box, ok := SealShared([]byte(testMessages[i]), peerPublicList)
		if !ok {
			fmt.Println("Shared boxing failed: message", i)
			t.FailNow()
		}
		testBoxes[i] = string(box)
	}
}

func TestSharedUnboxing(t *testing.T) {
	for i := 0; i < len(testMessages); i++ {
		for kn := 0; kn < 4; kn++ {
			m, ok := OpenShared([]byte(testBoxes[i]),
				peerPrivList[kn],
				peerPublicList[kn])
			if !ok {
				fmt.Println("Shared unboxing failed: message", i)
				fmt.Printf("box: %x\n", testBoxes[i])
				t.FailNow()
			} else if string(m) != testMessages[i] {
				fmt.Println("Shared unboxing did not return same plaintext.")
				t.FailNow()
			}
			_, ok = OpenShared([]byte(testBoxes[i]),
				testPeerKey, testPeerPub)
			if ok {
				fmt.Println("Shared unboxing should have failed!")
				t.FailNow()
			}
		}
		_, ok := OpenShared(mutate([]byte(testBoxes[i])),
			peerPrivList[0], peerPublicList[0])
		if ok {
			fmt.Println("Unboxing should have failed: message", i)
			t.FailNow()
		}
	}
}

func TestSharedSignedBoxing(t *testing.T) {
	for i := 0; i < len(testMessages); i++ {
		box, ok := SignAndSealShared([]byte(testMessages[i]), peerPublicList, testGoodKey,
			testGoodPub)
		if !ok {
			fmt.Println("Shared boxing failed: message", i)
			t.FailNow()
		}
		testBoxes[i] = string(box)
	}
}

func TestSharedSignedUnboxing(t *testing.T) {
	for i := 0; i < len(testMessages); i++ {
		for kn := 0; kn < 4; kn++ {
			m, ok := OpenSharedAndVerify([]byte(testBoxes[i]),
				peerPrivList[kn],
				peerPublicList[kn],
				testGoodPub)
			if !ok {
				fmt.Println("Shared unboxing failed: message", i)
				fmt.Printf("box: %x\n", testBoxes[i])
				t.FailNow()
			} else if string(m) != testMessages[i] {
				fmt.Println("Shared unboxing did not return same plaintext.")
				t.FailNow()
			}
			_, ok = OpenSharedAndVerify([]byte(testBoxes[i]),
				testPeerKey, testPeerPub, testGoodPub)
			if ok {
				fmt.Println("Shared unboxing should have failed!")
				t.FailNow()
			}
			_, ok = OpenSharedAndVerify([]byte(testBoxes[i]),
				peerPrivList[kn],
				peerPublicList[kn],
				testPeerPub)
			if ok {
				fmt.Println("Signature verification should have failed!")
				t.FailNow()
			}
		}
		_, ok := OpenSharedAndVerify(mutate([]byte(testBoxes[i])),
			peerPrivList[0],
			peerPublicList[0],
			testPeerPub)
		if ok {
			fmt.Println("Signature verification should have failed!")
			t.FailNow()
		}
	}
}

func BenchmarkSharedUnsignedSeal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, ok := SealShared(testBoxFile, peerPublicList)
		if !ok {
			fmt.Println("Couldn't seal message: benchmark aborted.")
			b.FailNow()
		}
	}
}

func BenchmarkSharedSignAndSeal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, ok := SignAndSealShared(testBoxFile, peerPublicList, testGoodKey,
			testGoodPub)
		if !ok {
			fmt.Println("Couldn't seal message: benchmark aborted.")
			b.FailNow()
		}
	}
}

// Benchmark the Open function, which retrieves a message from a box.
func BenchmarkSharedUnsignedOpen(b *testing.B) {
	box, ok := SealShared(testBoxFile, peerPublicList)
	if !ok {
		fmt.Println("Can't seal message: benchmark aborted.")
		b.FailNow()
	}
	for i := 0; i < b.N; i++ {
		_, ok := OpenShared(box, peerPrivList[3], peerPublicList[3])
		if !ok {
			fmt.Println("Couldn't open message: benchmark aborted.")
			b.FailNow()
		}
	}
}

// Benchmark the OpenSigned function, which retrieves a message from a box and verifies a
// signature on it.
func BenchmarkOpenSharedSigned(b *testing.B) {
	box, ok := SignAndSealShared(testBoxFile, peerPublicList, testGoodKey, testGoodPub)
	if !ok {
		fmt.Println("Can't seal message: benchmark aborted.")
		b.FailNow()
	}
	for i := 0; i < b.N; i++ {
		_, ok := OpenSharedAndVerify(box, peerPrivList[3], peerPublicList[3], testGoodPub)
		if !ok {
			fmt.Println("Couldn't open message: benchmark aborted.")
			b.FailNow()
		}
	}
}
