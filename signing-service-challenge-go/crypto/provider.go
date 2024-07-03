package crypto

type KeyPair struct {
	Public  []byte
	Private []byte
}

// Generator defines a contract for different types of signing implementations.
type Provider interface {
	Provide() (KeyPair, error)
}

// RSAGenerator generates a RSA key pair.
type RSAProvider struct {
	RSAGenerator
	RSAMarshaler
}

func (g *RSAProvider) Provide() (KeyPair, error) {
	pair, err := g.Generate()
	if err != nil {
		return KeyPair{}, err
	}

	publicKey, privateKey, err := g.Marshal(*pair)
	if err != nil {
		return KeyPair{}, err
	}

	return KeyPair{
		Public:  publicKey,
		Private: privateKey,
	}, nil

}

// RSAGenerator generates a RSA key pair.
type ECDSAProvider struct {
	ECCGenerator
	ECCMarshaler
}

func (g *ECDSAProvider) Provide() (KeyPair, error) {
	pair, err := g.Generate()
	if err != nil {
		return KeyPair{}, err
	}

	publicKey, privateKey, err := g.Encode(*pair)
	if err != nil {
		return KeyPair{}, err
	}

	return KeyPair{
		Public:  publicKey,
		Private: privateKey,
	}, nil
}
