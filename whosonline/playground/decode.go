package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

func main() {
	// Read armored private key into type EntityList
	// An EntityList contains one or more Entities.
	// This assumes there is only one Entity involved
	entitylist, err := openpgp.ReadArmoredKeyRing(bytes.NewBufferString(privateKey))
	if err != nil {
		log.Fatal(err)
	}
	entity := entitylist[0]
	fmt.Println("Private key from armored string:", entity.Identities)

	// Decrypt armor encrypted message using private key
	decbuf := bytes.NewBuffer([]byte(encryptedMessage))
	result, err := armor.Decode(decbuf)
	if err != nil {
		log.Fatal(err)
	}

	md, err := openpgp.ReadMessage(result.Body, entitylist, nil, nil)
	if err != nil {
		fmt.Println("error reading message", err)
	}

	bytes, err := ioutil.ReadAll(md.UnverifiedBody)
	fmt.Println("md:", string(bytes))

}

// sec   1024R/1D489768 2014-01-08
// uid                  Golang Test (Private key has no password) <golangtest@test.com>
// ssb   1024R/A46346C7 2014-01-08
const privateKey = `-----BEGIN PGP PRIVATE KEY BLOCK-----
Version: GnuPG v1

lQHYBFLMpyYBBACaEwtBnjtwzygaZadrx2ar2NiYc9FtsWtX0tHIA58UjtrrFVgr
SVhxIfFzsk2uIB/aTdZVtEgCJPw4NSwgSjQmzPGf5I47Q+pDWfcJSpiXYZfHSa+K
MV5Pmckw/IY6L0VGjgP/sPI8nPXG3WdQZBoXY0iR0r5I4cdMPGD3RDHlmwARAQAB
AAP8DcSy5k++DxQsCMPZxLtidOvhWypKZkuLBKOdxIpyTlN/w686HFAnG6EFYynl
YDCVuBvpkeoU30Av29QAl/gZldSNRGFVxUSox0W+OIZELICc+j2pWyD5fnv3Ezs1
TNIjtFsiQGtyWZ1QFQsTSLIVtT8UZWxbvDrJw46ndJlGo7UCAML6BkGkh5sqjHEv
9HS61t+QxI8+9OrphEFKtYYzpRHGxcZAWNue7q3czuXXNpP3emRHHjJ9ooUkBlS1
/X6KlZ0CAMpL1huNPwQ25geusz+pqEHLFnVOCbAf+QbQvWlRHeo3/t//eCl1naLK
GK29nXrKoM7wHqN1XOF6q9heOl5UXpcB/1E0jwJ+myuqUzGjWeWW51heh9iiQeBv
nHxqq7VHoXDE/TmNyc2/ByUBuYvuJCKN5Uge0u+2sa9UBHdKQiCARvSghLQ/R29s
YW5nIFRlc3QgKFByaXZhdGUga2V5IGhhcyBubyBwYXNzd29yZCkgPGdvbGFuZ3Rl
c3RAdGVzdC5jb20+iLgEEwECACIFAlLMpyYCGwMGCwkIBwMCBhUIAgkKCwQWAgMB
Ah4BAheAAAoJEMRDjVMdSJdoo3EEAJXzk+gvNTL3mADd2MhoYxI6V3S5x2yWQlum
o1afl4J9u+pyLZgSifX4OO23EFnfrRYivR4AGbmYqUZ/Gi0J0KBmpZpY5IULKf2v
V7RH49QFmp1twNTC3JG3ovOIID5aZlHKKfSHCz7xt9jngHbqjnSEX/BHR0K1RsC9
yOUc2qjinQHYBFLMpyYBBADaEokbf4iYagyagJTAwJSflmu2ihfslHOwVSjr+SOV
573kMQ03N9U9kESLm9CCC+b1G4vgBxuWBmvy5O0TR2R0Z9PYyr+81/Rr2xitq/yB
kf+b8WuhdXXc4mGm/V2WiP4zAT8ibtydn+NFITi8SNvGAZRasjRObzdoiJSAwE6A
6wARAQABAAP8C4iiJoBi7uMGuTIUNSnspLUDeY2XN+lDipRXnoUmAuQMOf2t4J37
/6CjGb7mnA/rzm8RnOXFoGIOIyHXN9l1cGA2LzzRKoNcSO5sk8c5WSbpsKQewFEM
cLMmC1KsTXAIjhdg54MgcKzBWySoCe/2mxlyHq++/YKew+ZgzBOTFkkCAN4ybSJG
9I0x5UvwRzLtKz3IPTaagdhZ8EHWDzvI9WS5vwENDFYnFmSQnjYEPLvntJCmpdku
d+TKWlEqEG6ASEkCAPs/ecZlUskW26GD/4hYEAvkOGLg4BH/oTgVqbwMRzlLtUn2
StPHXlfodc3ClaeIc0kivwivBHiLCML7v4NTB5MCAImcKRUrPGE8ZUM+C0WzTGTd
BJIeM9yRtF+H+MkmAnKJsFbDJx145C0KxANVpagEu7sz7txysWgkd7Mt3mPkKgWX
WoifBBgBAgAJBQJSzKcmAhsMAAoJEMRDjVMdSJdofBQD+gKh9E+kbugDIAejjE0I
1woNp/08ib4nRd7ZshEdykOZuvUEr7L7F6of197/x7IiL3pd+b/GTuv4h5N9rT92
U+ue0CwwhmmxYzu6Jkzet2fvmOdGFrKglnVBsreJkLXyDxAOHmFhulg7ABHVr3Vo
ybM0WVIalA8ZbghCuGrhSb0T
=HP3h
-----END PGP PRIVATE KEY BLOCK-----
`

// Encrypted by public key message for
// 1024R/A46346C7 2014-01-08 "Golang Test (Private key has no password) <golangtest@test.com>"
const encryptedMessage = `-----BEGIN PGP MESSAGE-----                  
Version: GnuPG v1

hIwDa78K16RjRscBA/9rEXnQ06mAOwrhTXwkps0cYhOBc/dX7BtdEPoRLBd0zNpr
hlU2YPcmQjBaL2Zi2E64K9Ud+R3D4RHt78a1145jjPFBMpWg1z2SMgERaRtwT2++
AinkXavDTuwoyPD7X3z4Jyp+aPcvReUQhq7idb2Kl5cjnXe1Z31a9wy4wMHEItJh
ARDEtteoYHsszGqTdIQOuGrI5P/V4biteh5JpxWqeTWDiY9HR2O31kwRhh9bV3J/
0kFgWpTRqfaSmu8ZiHtfQzTGCHpbLhz2IQgOSjBevzZACbn0z7h5Ro+EfAsonGfH
3w==
=xUUr
-----END PGP MESSAGE-----`
