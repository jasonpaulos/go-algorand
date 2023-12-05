// Copyright (C) 2019-2023 Algorand, Inc.
// This file is part of go-algorand
//
// go-algorand is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// go-algorand is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with go-algorand.  If not, see <https://www.gnu.org/licenses/>.

// Package node is the Algorand node itself, with functions exposed to the frontend

package p2p

import (
	"crypto/rand"
	"fmt"
	"os"
	"path"

	"github.com/algorand/go-algorand/config"
	algocrypto "github.com/algorand/go-algorand/crypto"
	"github.com/algorand/go-algorand/util"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
)

// DefaultPrivKeyPath is the default path inside the node's root directory at which the private key
// for p2p identity is found and persisted to when a new one is generated.
const DefaultPrivKeyPath = "peerIDPrivKey.pem"

// PeerID is a string representation of a peer's public key, primarily used to avoid importing libp2p into packages that shouldn't need it
type PeerID string // This seems unnecessary? Importing libp2p does not seem like an issue, unless we want to abstract away/hide the library

// GetPrivKey manages loading and creation of private keys for network PeerIDs
// It prioritizes, in this order:
//  1. user supplied path to privKey
//  2. default path to privKey,
//  3. generating a new privKey.
//
// If a new privKey is generated it will be saved to default path if cfg.P2PPersistPeerID.
func GetPrivKey(cfg config.Local, dataDir string) (crypto.PrivKey, error) {
	// if user-supplied, try to load it from there
	if cfg.P2PPrivateKeyLocation != "" {
		return loadPrivateKeyFromFile(cfg.P2PPrivateKeyLocation)
	}
	// if a default path key exists load it
	var defaultPrivKeyPath string
	if dataDir != "" {
		defaultPrivKeyPath = path.Join(dataDir, DefaultPrivKeyPath)
		if util.FileExists(defaultPrivKeyPath) {
			return loadPrivateKeyFromFile(defaultPrivKeyPath)
		}
	}
	// generate a new key
	privKey, err := generatePrivKey()
	if err != nil {
		return privKey, fmt.Errorf("failed to generate private key %w", err)
	}
	// if we want persistent PeerID, save the generated PrivKey
	if cfg.P2PPersistPeerID && defaultPrivKeyPath != "" {
		return privKey, writePrivateKeyToFile(defaultPrivKeyPath, privKey)
	}
	return privKey, nil
}

// PeerIDFromPublicKey returns a PeerID from a public key, thin wrapper over libp2p function doing the same
func PeerIDFromPublicKey(pubKey crypto.PubKey) (PeerID, error) {
	peerID, err := peer.IDFromPublicKey(pubKey)
	if err != nil {
		return "", err
	}
	return PeerID(peerID), nil
}

// loadPrivateKeyFromFile attempts to read raw privKey bytes from path
// It only supports Ed25519 keys.
func loadPrivateKeyFromFile(path string) (crypto.PrivKey, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// We only support Ed25519 keys
	return crypto.UnmarshalEd25519PrivateKey(bytes)
}

// writePrivateKeyToFile attempts to write raw privKey bytes to path
func writePrivateKeyToFile(path string, privKey crypto.PrivKey) error {
	bytes, err := privKey.Raw()
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes, 0600)
}

// generatePrivKey creates a new Ed25519 key
func generatePrivKey() (crypto.PrivKey, error) {
	priv, _, err := crypto.GenerateEd25519Key(rand.Reader)
	return priv, err
}

// PeerIDChallengeSigner implements the identityChallengeSigner interface in the network package.
type PeerIDChallengeSigner struct {
	key crypto.PrivKey
}

// Sign implements the identityChallengeSigner interface.
func (p *PeerIDChallengeSigner) Sign(message algocrypto.Hashable) algocrypto.Signature {
	return p.SignBytes(algocrypto.HashRep(message))
}

// SignBytes implements the identityChallengeSigner interface.
func (p *PeerIDChallengeSigner) SignBytes(message []byte) algocrypto.Signature {
	// libp2p Ed25519PrivateKey.Sign() returns a signature with a length of 64 bytes and no error
	sig, err := p.key.Sign(message)
	if len(sig) != len(algocrypto.Signature{}) {
		panic(fmt.Sprintf("invalid signature length: %d", len(sig)))
	}
	if err != nil {
		panic(err)
	}
	return algocrypto.Signature(sig)
}

// Verify implements the identityChallengeSigner interface.
func (p *PeerIDChallengeSigner) Verify(message algocrypto.Hashable, sig algocrypto.Signature) bool {
	// libp2p Ed25519PublicKey.Verify() returns a bool and no error
	ret, err := p.key.GetPublic().Verify(algocrypto.HashRep(message), sig[:])
	if err != nil {
		panic(err)
	}
	return ret
}

// PublicKey implements the identityChallengeSigner interface.
func (p *PeerIDChallengeSigner) PublicKey() algocrypto.PublicKey {
	// libp2p Ed25519PublicKey.Raw() returns a 32-byte public key and no error
	pub, err := p.key.GetPublic().Raw()
	if len(pub) != len(algocrypto.PublicKey{}) {
		panic(fmt.Sprintf("invalid public key length: %d", len(pub)))
	}
	if err != nil {
		panic(err)
	}
	return algocrypto.PublicKey(pub)
}
